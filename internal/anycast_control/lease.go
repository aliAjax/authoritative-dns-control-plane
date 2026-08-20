package anycast_control

import (
	"fmt"
	"sync"
	"time"
)

type Registry struct {
	mu     sync.Mutex
	nodes  map[string]Node
	leases map[string]time.Time
}

func NewRegistry() *Registry {
	return &Registry{nodes: map[string]Node{}, leases: map[string]time.Time{}}
}
func (r *Registry) Register(n Node) error {
	if n.ID == "" || n.Address == "" {
		return fmt.Errorf("node id/address required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[n.ID] = n
	return nil
}
func (r *Registry) Heartbeat(id string, capacity, latency, availability float64, bgp bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.nodes[id]
	if !ok {
		return fmt.Errorf("node not registered")
	}
	n.Capacity = capacity
	n.LatencyMS = latency
	n.Availability = availability
	n.BGPAdvertised = bgp
	n.LeaseUntil = time.Now().Add(30 * time.Second)
	n.Fencing++
	r.nodes[id] = n
	return nil
}
func (r *Registry) Get(id string) (Node, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.nodes[id]
	return n, ok
}
func (r *Registry) List() []Node {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Node, 0, len(r.nodes))
	for _, n := range r.nodes {
		out = append(out, n)
	}
	return out
}
func (r *Registry) Remove(id string) { r.mu.Lock(); delete(r.nodes, id); r.mu.Unlock() }
func (r *Registry) Expired(now time.Time) []Node {
	out := []Node{}
	for _, n := range r.List() {
		if n.LeaseUntil.Before(now) {
			out = append(out, n)
		}
	}
	return out
}
