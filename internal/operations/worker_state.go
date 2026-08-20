package operations

import "context"

func (w *Worker) claimStart() bool {
	w.lifecycle.Lock()
	defer w.lifecycle.Unlock()
	if w.started || w.stopped {
		return false
	}
	w.started = true
	w.stopped = false
	return true
}

func (w *Worker) setCancel(cancel context.CancelFunc) {
	w.lifecycle.Lock()
	w.cancel = cancel
	w.lifecycle.Unlock()
}

func (w *Worker) claimStop() (context.CancelFunc, <-chan struct{}, bool) {
	w.lifecycle.Lock()
	defer w.lifecycle.Unlock()
	if !w.started || w.stopped {
		return nil, nil, false
	}
	w.stopped = true
	return w.cancel, w.done, true
}
