package resolver_runtime

import (
	"testing"
	"time"
)

func TestCacheOwnsInputPacket(t *testing.T) {
	c := NewCache(2)
	p := []byte{1, 2, 3}
	c.Put("a", p, time.Minute, false)
	p[0] = 9
	got, ok, _ := c.Get("a")
	if !ok || got[0] != 1 {
		t.Fatalf("cached response aliases caller packet: %v", got)
	}
}
