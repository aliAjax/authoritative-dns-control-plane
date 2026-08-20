package catalog

import (
	"fmt"
	"sync"
	"time"
)

type SerialAllocator struct {
	mu     sync.Mutex
	values map[string]uint32
}

func NewSerialAllocator() *SerialAllocator { return &SerialAllocator{values: map[string]uint32{}} }
func (a *SerialAllocator) Next(zone string) uint32 {
	a.mu.Lock()
	defer a.mu.Unlock()
	x := a.values[zone]
	now := uint32(time.Now().Unix())
	if x < now {
		x = now
	}
	x++
	if x == 0 {
		x = 1
	}
	a.values[zone] = x
	return x
}
func (a *SerialAllocator) Reserve(zone string, n uint32) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if n <= a.values[zone] {
		return fmt.Errorf("serial regression")
	}
	a.values[zone] = n
	return nil
}
func (a *SerialAllocator) Current(zone string) uint32 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.values[zone]
}
func (a *SerialAllocator) Reset(zone string) { a.mu.Lock(); delete(a.values, zone); a.mu.Unlock() }
func (a *SerialAllocator) Snapshot() map[string]uint32 {
	a.mu.Lock()
	defer a.mu.Unlock()
	r := map[string]uint32{}
	for k, v := range a.values {
		r[k] = v
	}
	return r
}
