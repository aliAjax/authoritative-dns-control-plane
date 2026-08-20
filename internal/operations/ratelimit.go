package operations

import (
	"sync"
	"time"
)

type Bucket struct {
	mu          sync.Mutex
	tokens      float64
	updated     time.Time
	rate, burst float64
}

func NewBucket(rate, burst float64) *Bucket {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = rate
	}
	return &Bucket{tokens: burst, updated: time.Now(), rate: rate, burst: burst}
}
func (b *Bucket) Allow(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.tokens += now.Sub(b.updated).Seconds() * b.rate
	if b.tokens > b.burst {
		b.tokens = b.burst
	}
	b.updated = now
	if b.tokens < n {
		return false
	}
	b.tokens -= n
	return true
}
func (b *Bucket) Remaining() float64 { b.mu.Lock(); defer b.mu.Unlock(); return b.tokens }

type Limiter struct {
	mu          sync.Mutex
	users       map[string]*Bucket
	rate, burst float64
}

func NewLimiter(rate, burst float64) *Limiter {
	return &Limiter{users: map[string]*Bucket{}, rate: rate, burst: burst}
}
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	b := l.users[key]
	if b == nil {
		b = NewBucket(l.rate, l.burst)
		l.users[key] = b
	}
	l.mu.Unlock()
	return b.Allow(1)
}
func (l *Limiter) Reset(key string) { l.mu.Lock(); defer l.mu.Unlock(); delete(l.users, key) }
func (l *Limiter) Size() int        { l.mu.Lock(); defer l.mu.Unlock(); return len(l.users) }
