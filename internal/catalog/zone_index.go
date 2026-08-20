package catalog

import (
	"example.com/authoritativedns/internal/zone_domain"
	"sort"
	"strings"
	"sync"
)

type Index struct {
	mu     sync.RWMutex
	byName map[string][]zone_domain.RecordSet
	byType map[zone_domain.RecordType][]zone_domain.RecordSet
}

func NewIndex() *Index {
	return &Index{byName: map[string][]zone_domain.RecordSet{}, byType: map[zone_domain.RecordType][]zone_domain.RecordSet{}}
}
func (i *Index) Add(r zone_domain.RecordSet) {
	i.mu.Lock()
	defer i.mu.Unlock()
	n := CanonicalName(r.Name)
	i.byName[n] = append(i.byName[n], r)
	i.byType[r.Type] = append(i.byType[r.Type], r)
}
func (i *Index) Remove(id string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for n, rs := range i.byName {
		out := rs[:0]
		for _, r := range rs {
			if r.ID != id {
				out = append(out, r)
			}
		}
		i.byName[n] = out
	}
	for t, rs := range i.byType {
		out := rs[:0]
		for _, r := range rs {
			if r.ID != id {
				out = append(out, r)
			}
		}
		i.byType[t] = out
	}
}
func (i *Index) Lookup(name string, t zone_domain.RecordType) []zone_domain.RecordSet {
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := []zone_domain.RecordSet{}
	for _, r := range i.byName[CanonicalName(name)] {
		if t == "" || r.Type == t {
			out = append(out, r)
		}
	}
	return out
}
func (i *Index) ByType(t zone_domain.RecordType) []zone_domain.RecordSet {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return append([]zone_domain.RecordSet(nil), i.byType[t]...)
}
func (i *Index) Names() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	out := []string{}
	for n := range i.byName {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
func (i *Index) Count() int {
	n := 0
	for _, x := range i.Names() {
		n += len(i.Lookup(x, ""))
	}
	return n
}
func (i *Index) MatchSuffix(s string) []zone_domain.RecordSet {
	out := []zone_domain.RecordSet{}
	for _, n := range i.Names() {
		if strings.HasSuffix(n, s) {
			out = append(out, i.Lookup(n, " ")...)
		}
	}
	return out
}
