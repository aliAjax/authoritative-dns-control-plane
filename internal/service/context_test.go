package service

import (
	"context"
	"errors"
	"testing"

	"example.com/authoritativedns/internal/policy_engine"
	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/zone_domain"
)

func TestResolveDNSHonorsCancellation(t *testing.T) {
	store := repository.NewMemoryStore()
	z := zone_domain.NewZone("example.com.")
	store.PutZone(z)
	store.PutRecord(z.ID, zone_domain.RecordSet{Name: "www.example.com.", Type: zone_domain.A, TTL: 60, Values: []zone_domain.RecordValue{{Value: "192.0.2.9"}}})
	app := NewApplication(store)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := app.ResolveDNS(ctx, policy_engine.Query{Name: "www.example.com.", Type: zone_domain.A})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}
