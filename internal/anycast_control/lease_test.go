package anycast_control

import (
	"testing"
	"time"
)

func TestRegistryHeartbeatRenewsLeaseAndFencing(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(Node{ID: "tokyo-a", Address: "192.0.2.10"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Heartbeat("tokyo-a", 10, 15, 0.99, true); err != nil {
		t.Fatal(err)
	}
	n, ok := r.Get("tokyo-a")
	if !ok || n.Fencing != 1 || !n.LeaseUntil.After(time.Now()) || !n.BGPAdvertised {
		t.Fatalf("unexpected heartbeat state: %#v, exists=%v", n, ok)
	}
}
