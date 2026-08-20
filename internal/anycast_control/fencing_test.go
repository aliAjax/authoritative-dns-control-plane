package anycast_control

import (
	"testing"
	"time"
)

func TestPlanNeverRegressesFencingToken(t *testing.T) {
	now := time.Now().Add(time.Minute)
	base := Node{ID: "edge-a", Capacity: 1, Availability: 0.99, LatencyMS: 10, LeaseUntil: now, Fencing: 2}
	start := make(chan struct{})
	results := make(chan error, 2)
	go func() { <-start; _, err := Plan([]Node{base}); results <- err }()
	go func() { <-start; _, err := Plan([]Node{base}); results <- err }()
	close(start)
	if err := <-results; err != nil {
		t.Fatal(err)
	}
	if err := <-results; err != nil {
		t.Fatal(err)
	}
	nodes := []Node{base}
	nodes[0].Fencing = 1
	if _, err := Plan(nodes); err != nil {
		t.Fatal(err)
	}
	got, err := Plan(nodes)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("stale fencing token produced intent: %#v", got)
	}
}
