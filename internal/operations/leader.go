package operations

import (
	"fmt"
	"sync"
	"time"
)

type Lease struct {
	Key     string
	Owner   string
	Token   uint64
	Expires time.Time
}
type LeaseManager struct {
	mu     sync.Mutex
	items  map[string]Lease
	tokens map[string]uint64
}

func NewLeaseManager() *LeaseManager {
	return &LeaseManager{items: map[string]Lease{}, tokens: map[string]uint64{}}
}
func (l *LeaseManager) Acquire(key, owner string, d time.Duration) (Lease, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if old, ok := l.items[key]; ok && time.Now().Before(old.Expires) && old.Owner != owner {
		return Lease{}, fmt.Errorf("lease held by %s", old.Owner)
	}
	l.tokens[key]++
	x := Lease{Key: key, Owner: owner, Token: l.tokens[key], Expires: time.Now().Add(d)}
	l.items[key] = x
	return x, nil
}
func (l *LeaseManager) Renew(x Lease, d time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	old, ok := l.items[x.Key]
	if !ok || old.Owner != x.Owner || old.Token != x.Token {
		return fmt.Errorf("lease fencing failed")
	}
	old.Expires = time.Now().Add(d)
	l.items[x.Key] = old
	return nil
}
func (l *LeaseManager) Release(x Lease) {
	l.mu.Lock()
	defer l.mu.Unlock()
	old := l.items[x.Key]
	if old.Owner == x.Owner && old.Token == x.Token {
		delete(l.items, x.Key)
	}
}
func (l *LeaseManager) Current(key string) (Lease, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	x, ok := l.items[key]
	if !ok || time.Now().After(x.Expires) {
		return Lease{}, false
	}
	return x, true
}
