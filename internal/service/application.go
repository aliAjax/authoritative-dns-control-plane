package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"example.com/authoritativedns/internal/anycast_control"
	"example.com/authoritativedns/internal/dnssec"
	"example.com/authoritativedns/internal/policy_engine"
	"example.com/authoritativedns/internal/probe_domain"
	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/snapshot_release"
	"example.com/authoritativedns/internal/zone_domain"
)

type Application struct {
	Store    *repository.MemoryStore
	mu       sync.RWMutex
	policies map[string]policy_engine.Policy
	probes   map[string]probe_domain.Probe
	keys     map[string]dnssec.Key
}

func NewApplication(s *repository.MemoryStore) *Application {
	return &Application{Store: s, policies: map[string]policy_engine.Policy{}, probes: map[string]probe_domain.Probe{}, keys: map[string]dnssec.Key{}}
}
func (a *Application) CreateZone(ctx context.Context, name string) (zone_domain.Zone, error) {
	n := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(name)), ".")
	if n == "" || !strings.Contains(n, ".") {
		return zone_domain.Zone{}, fmt.Errorf("invalid zone name")
	}
	z := zone_domain.NewZone(n + ".")
	a.Store.PutZone(z)
	return z, nil
}
func (a *Application) GetZone(ctx context.Context, id string) (zone_domain.Zone, bool) {
	return a.Store.GetZone(id)
}
func (a *Application) ListZones(ctx context.Context) []zone_domain.Zone { return a.Store.ListZones() }
func (a *Application) AddRecord(ctx context.Context, zoneID string, rs zone_domain.RecordSet) (zone_domain.RecordSet, error) {
	z, ok := a.Store.GetZone(zoneID)
	if !ok {
		return rs, fmt.Errorf("zone not found")
	}
	if err := z.ValidateRecord(rs); err != nil {
		return rs, err
	}
	rs.Status = zone_domain.Draft
	rs.UpdatedAt = time.Now().UTC()
	a.Store.PutRecord(zoneID, rs)
	return rs, nil
}
func (a *Application) ListRecords(ctx context.Context, zoneID string) []zone_domain.RecordSet {
	return a.Store.ListRecords(zoneID)
}
func (a *Application) Validate(ctx context.Context, id string) (snapshot_release.Validation, error) {
	return snapshot_release.ValidateZone(a.Store, id)
}
func (a *Application) Publish(ctx context.Context, id string) (zone_domain.Snapshot, error) {
	return snapshot_release.Publish(ctx, a.Store, id)
}
func (a *Application) Rollback(ctx context.Context, id, snapshotID string) (zone_domain.Snapshot, error) {
	return snapshot_release.Rollback(ctx, a.Store, id, snapshotID)
}
func (a *Application) Resolve(ctx context.Context, q policy_engine.Query) (policy_engine.Decision, error) {
	return policy_engine.Resolve(a.Store, q)
}

// ResolveDNS selects the longest matching authoritative zone before evaluating records.
func (a *Application) ResolveDNS(ctx context.Context, q policy_engine.Query) (policy_engine.Decision, error) {
	best := ""
	for _, z := range a.Store.ListZones() {
		if strings.HasSuffix(strings.ToLower(q.Name), strings.ToLower(z.Name)) && len(z.Name) > len(best) {
			best = z.Name
			q.ZoneID = z.ID
		}
	}
	if q.ZoneID == "" {
		return policy_engine.Decision{Negative: true, Explanation: "no authoritative zone"}, nil
	}
	return a.Resolve(ctx, q)
}
func (a *Application) AddPolicy(p policy_engine.Policy) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.policies[p.ID] = p
}
func (a *Application) Policies() []policy_engine.Policy {
	a.mu.RLock()
	defer a.mu.RUnlock()
	r := make([]policy_engine.Policy, 0, len(a.policies))
	for _, p := range a.policies {
		r = append(r, p)
	}
	return r
}
func (a *Application) AddProbe(p probe_domain.Probe) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.probes[p.ID] = p
}
func (a *Application) AddKey(k dnssec.Key) { a.mu.Lock(); defer a.mu.Unlock(); a.keys[k.ID] = k }
func (a *Application) AnycastPlan(ctx context.Context, nodes []anycast_control.Node) ([]anycast_control.Intent, error) {
	return anycast_control.Plan(nodes)
}
