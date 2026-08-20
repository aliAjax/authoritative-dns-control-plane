package service

import (
	"context"
	"errors"
	"testing"

	"example.com/authoritativedns/internal/policy_engine"
	"example.com/authoritativedns/internal/repository"
)

func TestResolveHonorsCancellation(t *testing.T) {
	app := NewApplication(repository.NewMemoryStore())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := app.Resolve(ctx, policy_engine.Query{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected direct resolve cancellation, got %v", err)
	}
}
