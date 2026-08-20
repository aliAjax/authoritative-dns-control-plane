package operations

import (
	"context"
	"testing"
	"time"
)

func TestWorkerStartIsIdempotent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := NewWorker("reconcile", time.Millisecond, func(context.Context) error { return nil })
	start := make(chan struct{})
	started := make(chan struct{}, 2)
	go func() { <-start; w.Start(ctx); started <- struct{}{} }()
	go func() { <-start; w.Start(ctx); started <- struct{}{} }()
	close(start)
	<-started
	<-started
	cancel()
	done := make(chan struct{}, 2)
	stopStart := make(chan struct{})
	go func() { <-stopStart; w.Stop(); done <- struct{}{} }()
	go func() { <-stopStart; w.Stop(); done <- struct{}{} }()
	close(stopStart)
	select {
	case <-done:
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("second worker stop did not return")
		}
	case <-time.After(time.Second):
		t.Fatal("worker did not stop after a repeated Start")
	}
}

func TestWorkerStopBeforeStartReturns(t *testing.T) {
	w := NewWorker("reconcile", time.Millisecond, func(context.Context) error { return nil })
	done := make(chan struct{}, 2)
	start := make(chan struct{})
	go func() { <-start; w.Stop(); done <- struct{}{} }()
	go func() { <-start; w.Stop(); done <- struct{}{} }()
	close(start)
	select {
	case <-done:
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("second idle worker stop blocked")
		}
	case <-time.After(time.Second):
		t.Fatal("stopping an idle worker blocked")
	}
}
