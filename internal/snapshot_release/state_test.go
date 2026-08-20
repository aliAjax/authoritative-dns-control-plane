package snapshot_release

import (
	"context"
	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

func TestRollbackThenPublishProducesPublishedRecords(t *testing.T) {
	s := repository.NewMemoryStore()
	z := zone_domain.NewZone("example.com.")
	s.PutZone(z)
	s.PutRecord(z.ID, zone_domain.RecordSet{ID: "www", Name: "www.example.com.", Type: zone_domain.A, TTL: 60, Values: []zone_domain.RecordValue{{Value: "192.0.2.1"}}})
	first, err := Publish(context.Background(), s, z.ID)
	if err != nil {
		t.Fatal(err)
	}
	rolled, err := Rollback(context.Background(), s, z.ID, first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rolled.Records[0].Status != zone_domain.RolledBack {
		t.Fatalf("rollback status lost: %#v", rolled)
	}
	second, err := Publish(context.Background(), s, z.ID)
	if err != nil || len(second.Records) != 1 || second.Records[0].Status != zone_domain.Published {
		t.Fatalf("rollback left stale state: %#v err=%v", second, err)
	}
}
