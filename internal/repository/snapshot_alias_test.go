package repository

import (
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

func TestSnapshotOwnsNestedRecords(t *testing.T) {
	s := NewMemoryStore()
	z := zone_domain.NewZone("snapshot.example.com.")
	s.PutZone(z)
	snap := zone_domain.Snapshot{ID: "snap", ZoneID: z.ID, Records: []zone_domain.RecordSet{{Values: []zone_domain.RecordValue{{Value: "192.0.2.4"}}}}}
	s.PutSnapshot(z.ID, snap)
	snap.Records[0].Values[0].Value = "203.0.113.4"
	got, ok := s.GetSnapshot(z.ID, "snap")
	if !ok || got.Records[0].Values[0].Value != "192.0.2.4" {
		t.Fatalf("snapshot stored caller-owned values: %#v", got)
	}
}
