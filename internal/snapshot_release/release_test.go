package snapshot_release

import (
	"context"
	"testing"

	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/zone_domain"
)

func TestPublishCreatesImmutableSnapshot(t *testing.T) {
	s := repository.NewMemoryStore()
	z := zone_domain.NewZone("example.com.")
	s.PutZone(z)
	s.PutRecord(z.ID, zone_domain.RecordSet{Name: "www.example.com.", Type: zone_domain.A, TTL: 60, Values: []zone_domain.RecordValue{{Value: "192.0.2.20"}}})
	snap, err := Publish(context.Background(), s, z.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !snap.Immutable || len(snap.Records) != 1 || snap.Serial == 0 {
		t.Fatalf("unexpected snapshot: %#v", snap)
	}
}
