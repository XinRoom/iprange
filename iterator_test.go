package iprange

import (
	"fmt"
	"net"
	"testing"
)

func TestName(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{input: "1.1.1.1", want: "1.1.1.1"},
		{input: "1.1.1.2/30", want: "1.1.1.3"},
		{input: "1.1.1.0-255", want: "1.1.1.200"},
		{input: "1.1-2.0-1.4", want: "1.2.0.4"},
		{input: "1.1.1.1-1.1.2.1", want: "1.1.2.0"},
		{input: "2001::59:63", want: "2001::59:63"},
		{input: "2001::59:63/126", want: "2001::59:62"},
		{input: "2001::59:63-f2", want: "2001::59:f0"},
		{input: "2001::59-60:63-f2", want: "2001::60:f0"},
		{input: "2001::59:63-2001::59:f2", want: "2001::59:f0"},
		{input: "0:0:0:0:0:ffff:aa51:0101", want: "170.81.1.1"},
	}

	for _, v := range testCases {
		t.Logf("Test %s", v.input)
		it, startIp, err := NewIter(v.input)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Size %d", it.TotalNum())
		index := it.TotalNum() - 1
		if index > 2 {
			index -= 2
		}

		index2 := it.TotalNum() - 1
		if index2 > 3 {
			index2 -= 3
		}

		// 根据序号，计算得到IP
		t.Logf("GetByIndex [%d]: %v", index, it.GetIpByIndex(index))
		t.Logf("GetByIndex2 [%d]: %v", index2, it.GetIpByIndex(index2))

		// 迭代
		it.GetIpByIndex(0) // rest index
		i := 0
		for itn := startIp; it.HasNext() && i <= 3; itn = it.Next() {
			t.Log(itn)
			i++
		}

		// 包含判断
		if !it.Contains(net.ParseIP(v.want)) {
			t.Error(fmt.Sprintf("[ERR] %s Contains %s?", v.input, v.want), it.Contains(net.ParseIP(v.want)))
		} else {
			t.Log(fmt.Sprintf("%s Contains %s?", v.input, v.want), it.Contains(net.ParseIP(v.want)))
		}
	}

	// 简单的获取IP序列
	ipSet, err := GenIpSet("1.1.1.1/30")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("1.1.1.1/30 GenIpSet is %s", ipSet)
}
