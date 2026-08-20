package snapshot_release

import (
	"context"
	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
	"sync"
	"time"
)

var locks sync.Map

type Validation struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
	Digest   string   `json:"digest"`
}

func ValidateZone(s *repository.MemoryStore, id string) (Validation, error) {
	z, ok := s.GetZone(id)
	if !ok {
		return Validation{}, fmt.Errorf("zone not found")
	}
	rs := s.ListRecords(id)
	v := Validation{Valid: true, Digest: zone_domain.SnapshotDigest(rs)}
	for _, r := range rs {
		if err := z.ValidateRecord(r); err != nil {
			v.Valid = false
			v.Errors = append(v.Errors, fmt.Sprintf("%s: %v", r.Name, err))
		}
		if r.TTL < 30 {
			v.Warnings = append(v.Warnings, r.Name+" ttl below recommended minimum")
		}
	}
	if len(z.NS) < 2 {
		v.Warnings = append(v.Warnings, "less than two authoritative NS records")
	}
	return v, nil
}
func Publish(ctx context.Context, s *repository.MemoryStore, id string) (zone_domain.Snapshot, error) {
	mu := lockFor(id)
	mu.Lock()
	defer mu.Unlock()
	v, err := ValidateZone(s, id)
	if err != nil || !v.Valid {
		if err != nil {
			return zone_domain.Snapshot{}, err
		}
		return zone_domain.Snapshot{}, fmt.Errorf("validation failed: %v", v.Errors)
	}
	z, ok := s.GetZone(id)
	if !ok {
		return zone_domain.Snapshot{}, fmt.Errorf("zone not found")
	}
	serial := s.NextSerial(id)
	now := time.Now().UTC()
	snap := zone_domain.Snapshot{ID: fmt.Sprintf("snap-%d", now.UnixNano()), ZoneID: id, Serial: serial, Records: s.ListRecords(id), Digest: v.Digest, CreatedAt: now, Immutable: true}
	for i := range snap.Records {
		snap.Records[i].Status = zone_domain.Published
	}
	s.PutSnapshot(id, snap)
	z.Serial = serial
	s.PutZone(z)
	return snap, nil
}
func Rollback(ctx context.Context, s *repository.MemoryStore, id, snapID string) (zone_domain.Snapshot, error) {
	mu := lockFor(id)
	mu.Lock()
	defer mu.Unlock()
	snap, ok := s.GetSnapshot(id, snapID)
	if !ok {
		return zone_domain.Snapshot{}, fmt.Errorf("snapshot not found")
	}
	rs := rollbackRecords(snap.Records)
	s.ReplaceRecords(id, rs)
	v, _ := ValidateZone(s, id)
	serial := s.NextSerial(id)
	out := zone_domain.Snapshot{ID: fmt.Sprintf("rollback-%d", time.Now().UnixNano()), ZoneID: id, Serial: serial, Records: rs, Digest: v.Digest, CreatedAt: time.Now().UTC(), Immutable: true}
	s.PutSnapshot(id, out)
	return out, nil
}
func lockFor(id string) *sync.Mutex {
	v, _ := locks.LoadOrStore(id, &sync.Mutex{})
	return v.(*sync.Mutex)
}
