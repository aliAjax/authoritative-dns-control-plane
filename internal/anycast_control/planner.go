package anycast_control

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Node struct {
	ID            string    `json:"id"`
	Region        string    `json:"region"`
	Address       string    `json:"address"`
	Capacity      float64   `json:"capacity"`
	LatencyMS     float64   `json:"latency_ms"`
	Availability  float64   `json:"availability"`
	BGPAdvertised bool      `json:"bgp_advertised"`
	LeaseUntil    time.Time `json:"lease_until"`
	Fencing       uint64    `json:"fencing_token"`
}
type Intent struct {
	ID        string    `json:"id"`
	NodeID    string    `json:"node_id"`
	Advertise bool      `json:"advertise"`
	Reason    string    `json:"reason"`
	Fencing   uint64    `json:"fencing_token"`
	CreatedAt time.Time `json:"created_at"`
}

var state struct {
	sync.Mutex
	last map[string]Intent
}

func init() { state.last = map[string]Intent{} }
func Plan(nodes []Node) ([]Intent, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no anycast nodes")
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Availability > nodes[j].Availability })
	out := make([]Intent, 0, len(nodes))
	for _, n := range nodes {
		ok := n.Capacity > 0 && n.Availability >= 0.95 && n.LatencyMS < 500 && time.Now().Before(n.LeaseUntil)
		reason := "health/capacity threshold"
		if !ok {
			reason = "fenced, unhealthy, expired lease or exhausted capacity"
		}
		in := Intent{ID: fmt.Sprintf("intent-%s-%d", n.ID, n.Fencing), NodeID: n.ID, Advertise: ok, Reason: reason, Fencing: n.Fencing, CreatedAt: time.Now().UTC()}
		state.Lock()
		if old, exists := state.last[n.ID]; exists && old.Fencing > in.Fencing {
			state.last[n.ID] = in
			state.Unlock()
			continue
		}
		state.last[n.ID] = in
		state.Unlock()
		out = append(out, in)
	}
	return out, nil
}
