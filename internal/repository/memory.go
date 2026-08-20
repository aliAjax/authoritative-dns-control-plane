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
	z = cloneZone(z)
	s.zones[z.ID] = z
	s.records[z.ID] = map[string]zone_domain.RecordSet{}
	s.serial[z.ID] = z.Serial
}
func (s *MemoryStore) GetZone(id string) (zone_domain.Zone, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	z, ok := s.zones[id]
	if !ok {
		return z, ok
	}
	return cloneZone(z), true
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
	r.Version++
	s.records[zoneID][r.ID] = cloneRecordSet(r)
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
	snap = cloneSnapshot(snap)
	s.snapshots[zoneID] = append(s.snapshots[zoneID], snap)
}
func (s *MemoryStore) ListSnapshots(zoneID string) []zone_domain.Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.snapshots[zoneID]
	out := make([]zone_domain.Snapshot, len(src))
	for i, snap := range src {
		out[i] = cloneSnapshot(snap)
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

func cloneStrings(in []string) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func cloneRecordValues(in []zone_domain.RecordValue) []zone_domain.RecordValue {
	if in == nil {
		return nil
	}
	out := make([]zone_domain.RecordValue, len(in))
	copy(out, in)
	return out
}

func cloneRecordSet(r zone_domain.RecordSet) zone_domain.RecordSet {
	r.Values = cloneRecordValues(r.Values)
	return r
}

func cloneRecordSets(in []zone_domain.RecordSet) []zone_domain.RecordSet {
	if in == nil {
		return nil
	}
	out := make([]zone_domain.RecordSet, len(in))
	for i, r := range in {
		out[i] = cloneRecordSet(r)
	}
	return out
}

func cloneZone(z zone_domain.Zone) zone_domain.Zone {
	z.NS = cloneStrings(z.NS)
	z.SOA = cloneRecordSet(z.SOA)
	return z
}

func cloneSnapshot(snap zone_domain.Snapshot) zone_domain.Snapshot {
	snap.Records = cloneRecordSets(snap.Records)
	return snap
}
