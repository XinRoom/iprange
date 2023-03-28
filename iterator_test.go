package iprange

import (
	"net"
	"testing"
)

func TestName(t *testing.T) {
	for _, v := range []string{"1.1.1.1", "1.1.1.2/30", "1.1.1.0-255", "1.1-2.0-1.4", "1.1.1.1-1.1.2.1", "2001::59:63", "2001::59:63/126", "2001::59:63-f2", "2001::59-60:63-f2", "2001::59:63-2001::59:f2"} {
		t.Logf("Test %s", v)
		it, startIp, err := NewIter(v)
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
		for itn := startIp; it.HasNext(); itn = it.Next() {
			t.Log(itn)
		}

		// 包含判断
		t.Log("Contains 1.1.1.0?", it.Contains(net.ParseIP("1.1.1.0")))
		t.Log("Contains 1.1.1.1?", it.Contains(net.ParseIP("1.1.1.1")))
		t.Log("Contains 1.1.1.3?", it.Contains(net.ParseIP("1.1.1.3")))
		t.Log("Contains 2001::59:63?", it.Contains(net.ParseIP("2001::59:63")))
		t.Log("Contains 2001::59:f2?", it.Contains(net.ParseIP("2001::59:f2")))
		t.Log("Contains 2001::59:f3?", it.Contains(net.ParseIP("2001::59:f3")))
	}

	// 简单的获取IP序列
	ipSet, err := GenIpSet("1.1.1.1/30")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("1.1.1.1/30 GenIpSet is %s", ipSet)
}
