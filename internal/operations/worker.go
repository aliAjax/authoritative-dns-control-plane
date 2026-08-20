package operations

import (
	"context"
	"sync"
	"time"
)

type Worker struct {
	mu        sync.Mutex
	lifecycle sync.Mutex
	name      string
	interval  time.Duration
	run       func(context.Context) error
	cancel    context.CancelFunc
	done      chan struct{}
	started   bool
	stopped   bool
	lastErr   error
	runs      uint64
}

func NewWorker(name string, interval time.Duration, fn func(context.Context) error) *Worker {
	if interval < time.Millisecond {
		interval = time.Second
	}
	return &Worker{name: name, interval: interval, run: fn, done: make(chan struct{})}
}
func (w *Worker) Start(ctx context.Context) {
	if !w.claimStart() {
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	w.setCancel(cancel)
	go func() {
		defer close(w.done)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.mu.Lock()
				w.runs++
				w.mu.Unlock()
				if e := w.run(ctx); e != nil {
					w.mu.Lock()
					w.lastErr = e
					w.mu.Unlock()
				}
			}
		}
	}()
}
func (w *Worker) Stop() {
	cancel, done, ok := w.claimStop()
	if !ok {
		return
	}
	if cancel != nil {
		cancel()
	}
	<-done
}
func (w *Worker) Stats() (uint64, error) { w.mu.Lock(); defer w.mu.Unlock(); return w.runs, w.lastErr }
