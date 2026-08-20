package probe_domain

import (
	"fmt"
	"sync"
	"time"
)

type Kind string

const (
	HTTP Kind = "http"
	TCP  Kind = "tcp"
	TLS  Kind = "tls"
	ICMP Kind = "icmp"
)

type Probe struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Kind              Kind          `json:"kind"`
	Endpoint          string        `json:"endpoint"`
	Interval          time.Duration `json:"interval"`
	FailureThreshold  int           `json:"failure_threshold"`
	RecoveryThreshold int           `json:"recovery_threshold"`
	Jitter            time.Duration `json:"jitter"`
}
type Result struct {
	ProbeID string        `json:"probe_id"`
	Region  string        `json:"region"`
	Healthy bool          `json:"healthy"`
	Latency time.Duration `json:"latency"`
	At      time.Time     `json:"at"`
	Error   string        `json:"error,omitempty"`
}
type Window struct {
	mu                  sync.Mutex
	samples             []Result
	max                 int
	healthy             bool
	failures, successes int
}

func NewWindow(n int) *Window {
	if n < 1 {
		n = 5
	}
	return &Window{max: n}
}
func (w *Window) Add(r Result) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.samples = append(w.samples, r)
	if len(w.samples) > w.max {
		w.samples = w.samples[1:]
	}
	if r.Healthy {
		w.successes++
		w.failures = 0
	} else {
		w.failures++
		w.successes = 0
	}
	if w.failures >= 3 {
		w.healthy = false
	}
	if w.successes >= 2 {
		w.healthy = true
	}
	return w.healthy
}
func (w *Window) Majority() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.samples) == 0 {
		return false
	}
	n := 0
	for _, r := range w.samples {
		if r.Healthy {
			n++
		}
	}
	return n*2 > len(w.samples)
}
func ValidateProbe(p Probe) error {
	if p.ID == "" || p.Endpoint == "" {
		return fmt.Errorf("probe id and endpoint required")
	}
	switch p.Kind {
	case HTTP, TCP, TLS, ICMP:
	default:
		return fmt.Errorf("unsupported probe kind")
	}
	if p.FailureThreshold < 1 || p.RecoveryThreshold < 1 {
		return fmt.Errorf("threshold must be positive")
	}
	return nil
}
