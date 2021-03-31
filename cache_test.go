package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache(3)
	c.Set(1, 1, time.Hour)
	c.Set(1, 2, time.Hour)
	value, ok := c.Get(1)
	if ok {
		x := value.(int)
		if x == 2 {
			t.Logf("success, expect:2, got:%v", x)
		} else {
			t.Errorf("error, expect:2, got:%v", x)
		}
	}

}

func TestCache2(t *testing.T) {
	c := NewCache(3)
	for i := 1; i < 10; i++ {
		c.Set(i, i, time.Hour)
	}
	//fmt.Println(c.interalMap)
	//c.list.TraversalPrint()
	_, ok := c.Get(1)
	if !ok {
		t.Logf("1.success, expect:false, got:%v", ok)
	} else {
		t.Errorf("1.error, expect:true, got:%v", ok)
	}
	_, ok = c.Get(9)
	if ok {
		t.Logf("2.success, expect:true, got:%v", ok)
	} else {
		t.Errorf("2.error, expect:false, got:%v", ok)
	}
}
