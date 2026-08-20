package probe_domain

import "testing"

func TestProbeWindowTracksHealthTransitions(t *testing.T) {
	w := NewWindow(5)
	if w.Add(Result{Healthy: true}) {
		t.Fatal("one success must not make a probe healthy")
	}
	if !w.Add(Result{Healthy: true}) {
		t.Fatal("two consecutive successes must make a probe healthy")
	}
	w.Add(Result{Healthy: false})
	w.Add(Result{Healthy: false})
	if w.Add(Result{Healthy: false}) {
		t.Fatal("three consecutive failures must make a probe unhealthy")
	}
}
