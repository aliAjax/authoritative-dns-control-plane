package anycast_control

import (
	"testing"
	"time"
)

func TestPlanDoesNotRepeatSameFencingIntent(t *testing.T) {
	now := time.Now().Add(time.Minute)
	base := Node{ID: "edge-idempotent", Capacity: 1, Availability: 0.99, LatencyMS: 10, LeaseUntil: now, Fencing: 4}
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
	got, err := Plan(nodes)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("same fencing token produced duplicate intent: %#v", got)
	}
}
