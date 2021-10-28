package iprange

import "testing"

func TestName(t *testing.T) {
	for _, v := range []string{"1.1.1.1", "1.1.1.1/30", "1.1.1.1-3", "1.1-2.0-1.4", "1.1.1.1-1.1.2.1", "2001::59:63", "2001::59:63/126", "2001::59:63-f2", "2001::59-60:63-f2", "2001::59:63-2001::59:f2"} {
		t.Logf("Test %s", v)
		it, err := NewIter(v)
		if err != nil {
			t.Fatal(err)
		}
		for itn := it.Next(); it.HasNext(); itn = it.Next() {
			t.Log(itn)
		}
	}

	ipSet, err := GenIpSet("1.1.1.1/30")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("1.1.1.1/30 GenIpSet is %s", ipSet)
}
