package iprange

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
)

const (
	CidrMode   = iota // CIDR
	WideMode          // 1.1.1.1-1.1.2.3
	NarrowMode        // 1-3.1-5.4.1-7
)

// RangeParseMate Range Info
type RangeParseMate struct {
	s uint8 //start
	e uint8 //end
}

// RangeClassMate IP Range of every byte
type RangeClassMate []RangeParseMate

// Iter Iterator
type Iter struct {
	mode      int            // 模式
	isIpv6    bool           // 是否是ipv6
	isIpv4    bool           // 是否是ipv4
	ipStr     string         // 填充后的ip字符串
	lastIp    net.IP         // ip迭代空间
	classmate RangeClassMate // ip范围限制信息
	ipNet     *net.IPNet     // cidr 模式下的网段信息
	sip       net.IP         // 开始IP
	eip       net.IP         // 结束IP
	done      bool           // 结束
	totalNum  uint64         // IP总数	max 64bit(return 0xffffffffffffffff if overflow)
}

func NewIter(ipStr string) (it *Iter, startIp net.IP, err error) {

	it = &Iter{
		mode:      -1,
		isIpv6:    false,
		isIpv4:    false,
		ipStr:     "",
		lastIp:    nil,
		classmate: nil,
		ipNet:     nil,
		sip:       nil,
		eip:       nil,
		done:      false,
		totalNum:  0,
	}

	// IP判断和填充
	if strings.Contains(ipStr, ".") { // 分段生成IPv4
		it.isIpv4 = true
	} else if strings.Contains(ipStr, ":") {
		it.isIpv6 = true
		// 填充缩写
		if strings.Count(ipStr, "::") == 1 {
			// :: 扩展
			buf := ":0000"
			for i := strings.Count(ipStr, ":"); i < 7; i++ {
				buf += ":0000"
			}
			ipStr = strings.Replace(ipStr, "::", buf+":", 1)

			// 补零
			ipv6C := strings.Split(ipStr, ":")
			for i, v := range ipv6C {
				ipv6D := strings.Split(v, "-")
				for i2, v2 := range ipv6D {
					for len(v2) < 4 {
						v2 = "0" + v2
					}
					ipv6D[i2] = v2
				}
				ipv6C[i] = strings.Join(ipv6D, "-")
			}
			ipStr = strings.Join(ipv6C, ":")
		}
	}
	it.ipStr = ipStr
	if !it.isIpv4 && !it.isIpv6 {
		return nil, nil, fmt.Errorf("not is ip")
	}

	// CidrMode
	ip, ipNet, err := net.ParseCIDR(ipStr)
	if err == nil {
		it.sip = ip.Mask(ipNet.Mask)
		it.ipNet = ipNet
		it.mode = CidrMode
	}

	// WideMode
	if it.mode == -1 && strings.Count(ipStr, "-") == 1 {
		startIpStrList := strings.Split(ipStr, "-")
		if len(startIpStrList) == 2 {
			sip := net.ParseIP(startIpStrList[0])
			eip := net.ParseIP(startIpStrList[1])
			if sip == nil || eip == nil || len(sip) != len(eip) {
				err = fmt.Errorf("WideMode parse ip err: %s", ipStr)
			} else {
				it.mode = WideMode
				it.sip = sip
				it.eip = eip
			}
		}
	}

	// NarrowMode
	if it.mode == -1 {
		var ipClasses []string
		if it.isIpv4 { // 分段生成IPv4
			ipClasses = strings.Split(ipStr, ".")
		} else if it.isIpv6 { // 分段生成IPv6
			ipClassesV6 := strings.Split(ipStr, ":")
			if len(ipClassesV6) != 8 {
				err = fmt.Errorf("NarrowMode ipv6 parse err %s", ipStr)
				return
			}
			// 2001::1112-3334
			// to ipClasses
			// 20,01,...,11-33,12-34    (16个)
			for _, v := range ipClassesV6 {
				if len(v) == 4 {
					ipClasses = append(ipClasses, v[:2])
					ipClasses = append(ipClasses, v[2:])
				} else if len(v) == 9 && strings.Contains(v, "-") {
					ipClasses = append(ipClasses, v[:2]+"-"+v[5:7])
					ipClasses = append(ipClasses, v[2:4]+"-"+v[7:])
				}
			}
		}

		// ipClasses to RangeParseMate
		if len(ipClasses) == 4 || len(ipClasses) == 16 {
			for _, v := range ipClasses {
				l0 := strings.Split(v, "-") // range
				var l0s uint64
				if it.isIpv4 {
					l0s, err = strconv.ParseUint(l0[0], 10, 8)
				} else {
					// ipv6 is hex
					l0s, err = strconv.ParseUint(l0[0], 16, 8)
				}
				l0e := l0s // The default start and end are the same
				if len(l0) > 2 || err != nil {
					return nil, nil, err
				}
				if len(l0) == 2 {
					if it.isIpv4 {
						l0e, err = strconv.ParseUint(l0[1], 10, 8)
					} else {
						l0e, err = strconv.ParseUint(l0[1], 16, 8)
					}
					if err != nil {
						return nil, nil, err
					}
				}
				it.classmate = append(it.classmate, RangeParseMate{
					s: uint8(l0s),
					e: uint8(l0e),
				})
			}
			//
			_startIp := make(net.IP, len(it.classmate))
			endIp := make(net.IP, len(it.classmate))
			for i, v := range it.classmate {
				_startIp[i] = v.s
				endIp[i] = v.e
			}
			it.mode = NarrowMode
			it.sip = _startIp
			it.eip = endIp
		}
	}

	if it.mode == -1 {
		return nil, nil, fmt.Errorf("unknow mode")
	}

	// dup copy sip to lastIp
	dup := make(net.IP, len(it.sip))
	copy(dup, it.sip)
	it.lastIp = dup
	it.totalNum = it.getTotalNum()
	return it, it.sip, nil
}

func (it *Iter) Next() net.IP {
	if !it.HasNext() {
		return nil
	}
	switch it.mode {
	case CidrMode:
		inc(it.lastIp)
		if !it.ipNet.Contains(it.lastIp) {
			it.done = true
			return nil
		}
	case WideMode:
		inc(it.lastIp)
		if bytes.Compare(it.eip, it.lastIp) < 0 {
			it.done = true
			return nil
		}
	case NarrowMode:
		classInc(it.lastIp, it.classmate)
		// 自增后置为初始值，说明到上限了
		if bytes.Compare(it.sip, it.lastIp) == 0 {
			it.done = true
			return nil
		}
	default:
		it.done = true
		return nil
	}

	dup := make(net.IP, len(it.lastIp))
	copy(dup, it.lastIp)
	return dup
}

func (it *Iter) HasNext() bool {
	return !it.done
}

// TotalNum Calculating the Total NUMBER of IP addresses
func (it *Iter) TotalNum() uint64 {
	return it.totalNum
}

func (it *Iter) getTotalNum() uint64 {
	switch it.mode {
	case CidrMode:
		ones, bits := it.ipNet.Mask.Size()
		return uint64(math.Pow(2, float64(bits-ones)))
	case WideMode:
		if it.eip.To4() != nil {
			return uint64(binary.BigEndian.Uint32(it.eip.To4()) - binary.BigEndian.Uint32(it.lastIp.To4()) + 1)
		} else {
			ret := big.NewInt(1)
			ret = ret.Add(ret, new(big.Int).Sub(new(big.Int).SetBytes(it.eip), new(big.Int).SetBytes(it.lastIp)))
			if ret.IsUint64() {
				return ret.Uint64()
			} else {
				return 0xffffffffffffffff
			}
		}
	case NarrowMode:
		var ret = uint64(1)
		for _, v := range it.classmate {
			ret = ret * uint64(v.e-v.s+1)
		}
		return ret
	}
	return 0
}

func (it *Iter) GetIpByIndex(index uint64) net.IP {
	if index >= it.totalNum {
		return nil
	}
	it.incByIndex(index)
	return it.lastIp
}

// IP increment by index
func (it *Iter) incByIndex(index uint64) {
	if it.classmate == nil {
		if it.isIpv4 {
			it.lastIp = it.lastIp.To4()
			binary.BigEndian.PutUint32(it.lastIp, binary.BigEndian.Uint32(it.sip.To4())+uint32(index))
		} else {
			ret := new(big.Int).SetBytes(it.sip)
			ret = ret.Add(ret, new(big.Int).SetUint64(index))
			it.lastIp = ret.Bytes()
		}
	} else {
		length := len(it.classmate)
		ip := make([]byte, length)
		rangeSpace := make([]uint8, length)
		// 每一位的空间容量
		for i, rangeMate := range it.classmate {
			rangeSpace[i] = rangeMate.e - rangeMate.s + 1
		}
		// transform 进位除余
		carryBit := uint64(0) // 进位
		var noFirst bool
		for i := length - 1; i >= 0; i-- {
			if rangeSpace[i] == 1 { // 如果空间为1，则不变
				ip[i] = it.sip[i]
			} else {
				if !noFirst {
					noFirst = true
					carryBit = index
				}
				ip[i] = it.sip[i] + uint8(carryBit%uint64(rangeSpace[i]))
				carryBit = uint64(uint8(carryBit / uint64(rangeSpace[i])))
			}
		}
		it.lastIp = ip
	}
}

// GenIpSet simple generate a set of ip
func GenIpSet(ipStr string) (outs []net.IP, err error) {
	it, startIp, err := NewIter(ipStr)
	if err != nil {
		return
	}
	for nit := startIp; it.HasNext(); nit = it.Next() {
		outs = append(outs, nit)
	}
	return
}

// IP increment
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// IP Segmented(byte) increment, 有范围的IP自增 1-4.1-4.1-4.1-4 = 1.1.1.1-4.4.4.4
func classInc(ip net.IP, classMate RangeClassMate) {
	for j := len(ip) - 1; j >= 0; j-- {
		// 当前分段最大限制
		if ip[j] >= classMate[j].e {
			ip[j] = classMate[j].s // 归初始值
			continue
		}
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
