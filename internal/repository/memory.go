package repository

import (
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
	"sync"
	"time"
)

type MemoryStore struct {
	mu        sync.RWMutex
	zones     map[string]zone_domain.Zone
	records   map[string]map[string]zone_domain.RecordSet
	snapshots map[string][]zone_domain.Snapshot
	serial    map[string]uint32
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{zones: map[string]zone_domain.Zone{}, records: map[string]map[string]zone_domain.RecordSet{}, snapshots: map[string][]zone_domain.Snapshot{}, serial: map[string]uint32{}}
}
func (s *MemoryStore) PutZone(z zone_domain.Zone) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.zones[z.ID] = cloneZone(z)
	s.records[z.ID] = map[string]zone_domain.RecordSet{}
	s.serial[z.ID] = z.Serial
}
func (s *MemoryStore) GetZone(id string) (zone_domain.Zone, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	z, ok := s.zones[id]
	return cloneZone(z), ok
}
func (s *MemoryStore) ListZones() []zone_domain.Zone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := make([]zone_domain.Zone, 0, len(s.zones))
	for _, z := range s.zones {
		r = append(r, cloneZone(z))
	}
	return r
}
func (s *MemoryStore) PutRecord(zoneID string, r zone_domain.RecordSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.records[zoneID] == nil {
		s.records[zoneID] = map[string]zone_domain.RecordSet{}
	}
	if r.ID == "" {
		r.ID = fmt.Sprintf("r-%d", time.Now().UnixNano())
	}
	r = cloneRecordSet(r)
	r.Version++
	s.records[zoneID][r.ID] = r
}
func (s *MemoryStore) ListRecords(zoneID string) []zone_domain.RecordSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.records[zoneID]
	r := make([]zone_domain.RecordSet, 0, len(m))
	for _, v := range m {
		r = append(r, cloneRecordSet(v))
	}
	return r
}
func (s *MemoryStore) ReplaceRecords(zoneID string, rs []zone_domain.RecordSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := map[string]zone_domain.RecordSet{}
	for _, r := range rs {
		m[r.ID] = cloneRecordSet(r)
	}
	s.records[zoneID] = m
}

func (s *MemoryStore) NextSerial(zoneID string) uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serial[zoneID]++
	if s.serial[zoneID] == 0 {
		s.serial[zoneID] = 1
	}
	z := s.zones[zoneID]
	z.Serial = s.serial[zoneID]
	z.UpdatedAt = time.Now().UTC()
	s.zones[zoneID] = z
	return z.Serial
}
func (s *MemoryStore) PutSnapshot(zoneID string, snap zone_domain.Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap.Immutable = true
	s.snapshots[zoneID] = append(s.snapshots[zoneID], cloneSnapshot(snap))
}
func (s *MemoryStore) ListSnapshots(zoneID string) []zone_domain.Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.snapshots[zoneID]
	out := make([]zone_domain.Snapshot, 0, len(items))
	for _, snap := range items {
		out = append(out, cloneSnapshot(snap))
	}
	return out
}
func (s *MemoryStore) GetSnapshot(zoneID, id string) (zone_domain.Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.snapshots[zoneID] {
		if x.ID == id {
			return cloneSnapshot(x), true
		}
	}
	return zone_domain.Snapshot{}, false
}
