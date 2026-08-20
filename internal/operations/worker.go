package operations

import (
	"context"
	"sync"
	"time"
)

type Worker struct {
	mu       sync.Mutex
	name     string
	interval time.Duration
	run      func(context.Context) error
	cancel   context.CancelFunc
	done     chan struct{}
	lastErr  error
	runs     uint64
}

func NewWorker(name string, interval time.Duration, fn func(context.Context) error) *Worker {
	if interval < time.Millisecond {
		interval = time.Second
	}
	return &Worker{name: name, interval: interval, run: fn}
}
func (w *Worker) Start(ctx context.Context) {
	w.mu.Lock()
	if w.cancel != nil {
		w.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	w.cancel = cancel
	w.done = done
	w.mu.Unlock()
	go func() {
		defer func() {
			w.mu.Lock()
			w.cancel = nil
			w.done = nil
			w.mu.Unlock()
			close(done)
		}()
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
	w.mu.Lock()
	cancel := w.cancel
	done := w.done
	w.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()
	<-done
}
func (w *Worker) Stats() (uint64, error) { w.mu.Lock(); defer w.mu.Unlock(); return w.runs, w.lastErr }
