package operations

import (
	"sync"
	"time"
)

type Event struct {
	ID       string    `json:"id"`
	Topic    string    `json:"topic"`
	Key      string    `json:"key"`
	Data     []byte    `json:"data"`
	At       time.Time `json:"at"`
	Attempts int       `json:"attempts"`
}
type Handler func(Event) error
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	closed   bool
	queue    chan Event
	done     chan struct{}
}

func NewBus(size int) *Bus {
	if size < 1 {
		size = 64
	}
	b := &Bus{handlers: map[string][]Handler{}, queue: make(chan Event, size), done: make(chan struct{})}
	go b.loop()
	return b
}
func (b *Bus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.closed {
		b.handlers[topic] = append(b.handlers[topic], h)
	}
}
func (b *Bus) Publish(e Event) error {
	b.mu.RLock()
	closed := b.closed
	b.mu.RUnlock()
	if closed {
		return contextClosed{}
	}
	if e.ID == "" {
		e.ID = time.Now().Format("20060102150405.000000000")
	}
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	select {
	case b.queue <- e:
		return nil
	default:
		return queueFull{}
	}
}
func (b *Bus) loop() {
	for e := range b.queue {
		b.mu.RLock()
		hs := append([]Handler(nil), b.handlers[e.Topic]...)
		b.mu.RUnlock()
		for _, h := range hs {
			_ = h(e)
		}
	}
	close(b.done)
}
func (b *Bus) Close() {
	b.mu.Lock()
	if !b.closed {
		b.closed = true
		close(b.queue)
	}
	b.mu.Unlock()
	<-b.done
}

type contextClosed struct{}

func (contextClosed) Error() string { return "event bus closed" }

type queueFull struct{}

func (queueFull) Error() string { return "event queue full" }
