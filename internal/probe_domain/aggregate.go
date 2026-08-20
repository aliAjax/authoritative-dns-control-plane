package probe_domain

import "time"

type Aggregator struct{ windows map[string]*Window }

func NewAggregator() *Aggregator { return &Aggregator{windows: map[string]*Window{}} }
func (a *Aggregator) Observe(region string, r Result) bool {
	w := a.windows[region]
	if w == nil {
		w = NewWindow(5)
		a.windows[region] = w
	}
	return w.Add(r)
}
func (a *Aggregator) HealthyRegions() []string {
	out := []string{}
	for k, w := range a.windows {
		if w.Majority() {
			out = append(out, k)
		}
	}
	return out
}
func (a *Aggregator) Availability() float64 {
	if len(a.windows) == 0 {
		return 0
	}
	n := 0
	for _, w := range a.windows {
		if w.Majority() {
			n++
		}
	}
	return float64(n) / float64(len(a.windows))
}
func (a *Aggregator) Stale(d time.Duration) []string {
	out := []string{}
	for k, w := range a.windows {
		for _, x := range w.samples {
			if time.Since(x.At) > d {
				out = append(out, k)
				break
			}
		}
	}
	return out
}
