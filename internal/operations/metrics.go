package operations

import (
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	value uint64
}

func (c *Counter) Inc()          { c.mu.Lock(); c.value++; c.mu.Unlock() }
func (c *Counter) Add(n uint64)  { c.mu.Lock(); c.value += n; c.mu.Unlock() }
func (c *Counter) Value() uint64 { c.mu.Lock(); defer c.mu.Unlock(); return c.value }

type Gauge struct {
	mu    sync.Mutex
	value float64
}

func (g *Gauge) Set(v float64)  { g.mu.Lock(); g.value = v; g.mu.Unlock() }
func (g *Gauge) Add(v float64)  { g.mu.Lock(); g.value += v; g.mu.Unlock() }
func (g *Gauge) Value() float64 { g.mu.Lock(); defer g.mu.Unlock(); return g.value }

type Histogram struct {
	mu     sync.Mutex
	values []float64
}

func (h *Histogram) Observe(v float64) { h.mu.Lock(); h.values = append(h.values, v); h.mu.Unlock() }
func (h *Histogram) Count() int        { h.mu.Lock(); defer h.mu.Unlock(); return len(h.values) }
func (h *Histogram) Quantile(q float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.values) == 0 {
		return 0
	}
	cp := append([]float64(nil), h.values...)
	for i := range cp {
		for j := i + 1; j < len(cp); j++ {
			if cp[j] < cp[i] {
				cp[i], cp[j] = cp[j], cp[i]
			}
		}
	}
	idx := int(float64(len(cp)-1) * q)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(cp) {
		idx = len(cp) - 1
	}
	return cp[idx]
}

type Registry struct {
	Started   time.Time
	Queries   Counter
	CacheHits Counter
	Errors    Counter
	Latency   Histogram
}

func NewRegistry() *Registry { return &Registry{Started: time.Now().UTC()} }
