package resolver_runtime

import (
	"testing"
	"time"
)

func TestCacheReturnsIsolatedPacketUntilExpiry(t *testing.T) {
	c := NewCache(2)
	c.Put("www.example.com/A", []byte{1, 2, 3}, time.Second, false)
	p, ok, negative := c.Get("www.example.com/A")
	if !ok || negative {
		t.Fatalf("cache miss or unexpected negative result: ok=%v negative=%v", ok, negative)
	}
	p[0] = 9
	p2, ok, _ := c.Get("www.example.com/A")
	if !ok || p2[0] != 1 {
		t.Fatalf("cache packet was mutated: %v", p2)
	}
}
