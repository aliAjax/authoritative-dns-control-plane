package service

import (
	"context"
	"errors"
	"testing"

	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/zone_domain"
)

func TestCreateZoneHonorsCancellation(t *testing.T) {
	app := NewApplication(repository.NewMemoryStore())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := app.CreateZone(ctx, "example.com"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected create cancellation, got %v", err)
	}
}

func TestAddRecordHonorsCancellation(t *testing.T) {
	store := repository.NewMemoryStore()
	z := zone_domain.NewZone("example.com.")
	store.PutZone(z)
	app := NewApplication(store)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := app.AddRecord(ctx, z.ID, zone_domain.RecordSet{Name: "www.example.com.", Type: zone_domain.A, TTL: 60, Values: []zone_domain.RecordValue{{Value: "192.0.2.8"}}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected add record cancellation, got %v", err)
	}
}
