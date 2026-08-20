package operations

import (
	"container/list"
	"sync"
	"time"
)

type cacheItem struct {
	key     string
	value   any
	expires time.Time
}
type LRU struct {
	mu    sync.Mutex
	cap   int
	ll    *list.List
	items map[string]*list.Element
}

func NewLRU(capacity int) *LRU {
	if capacity < 1 {
		capacity = 1
	}
	return &LRU{cap: capacity, ll: list.New(), items: map[string]*list.Element{}}
}
func (c *LRU) Put(k string, v any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e := c.items[k]; e != nil {
		e.Value.(*cacheItem).value = v
		e.Value.(*cacheItem).expires = time.Now().Add(ttl)
		c.ll.MoveToFront(e)
		return
	}
	e := c.ll.PushFront(&cacheItem{key: k, value: v, expires: time.Now().Add(ttl)})
	c.items[k] = e
	for c.ll.Len() > c.cap {
		z := c.ll.Back()
		delete(c.items, z.Value.(*cacheItem).key)
		c.ll.Remove(z)
	}
}
func (c *LRU) Get(k string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.items[k]
	if e == nil {
		return nil, false
	}
	x := e.Value.(*cacheItem)
	if time.Now().After(x.expires) {
		delete(c.items, k)
		c.ll.Remove(e)
		return nil, false
	}
	c.ll.MoveToFront(e)
	return x.value, true
}
func (c *LRU) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e := c.items[k]; e != nil {
		delete(c.items, k)
		c.ll.Remove(e)
	}
}
func (c *LRU) Len() int { c.mu.Lock(); defer c.mu.Unlock(); return c.ll.Len() }
func (c *LRU) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	r := []string{}
	for e := c.ll.Front(); e != nil; e = e.Next() {
		r = append(r, e.Value.(*cacheItem).key)
	}
	return r
}
