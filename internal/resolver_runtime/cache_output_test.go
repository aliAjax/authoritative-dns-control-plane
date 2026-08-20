package resolver_runtime

import (
	"testing"
	"time"
)

func TestCacheOwnsReturnedPacket(t *testing.T) {
	c := NewCache(2)
	c.Put("edge/A", []byte{4, 5, 6}, time.Second, false)
	p, ok, _ := c.Get("edge/A")
	if !ok {
		t.Fatal("cache miss")
	}
	p[0] = 9
	again, ok, _ := c.Get("edge/A")
	if !ok || again[0] != 4 {
		t.Fatalf("returned packet changed cached state: %v", again)
	}
}
