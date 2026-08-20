package operations

import (
	"sync"
	"time"
)

type HealthState struct {
	mu                              sync.Mutex
	healthy                         bool
	consecutiveGood, consecutiveBad int
	last                            time.Time
}

func (h *HealthState) Observe(ok bool, good, bad int) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.last = time.Now().UTC()
	if ok {
		h.consecutiveGood++
		h.consecutiveBad = 0
		if h.consecutiveGood >= good {
			h.healthy = true
		}
	} else {
		h.consecutiveBad++
		h.consecutiveGood = 0
		if h.consecutiveBad >= bad {
			h.healthy = false
		}
	}
	return h.healthy
}
func (h *HealthState) Healthy() bool   { h.mu.Lock(); defer h.mu.Unlock(); return h.healthy }
func (h *HealthState) Last() time.Time { h.mu.Lock(); defer h.mu.Unlock(); return h.last }
func Majority(values []bool) bool {
	if len(values) == 0 {
		return false
	}
	n := 0
	for _, x := range values {
		if x {
			n++
		}
	}
	return n*2 > len(values)
}
func Availability(values []bool) float64 {
	if len(values) == 0 {
		return 0
	}
	n := 0
	for _, x := range values {
		if x {
			n++
		}
	}
	return float64(n) / float64(len(values))
}
