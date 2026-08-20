package catalog

import (
	"example.com/authoritativedns/internal/zone_domain"
	"net"
	"sort"
	"strings"
)

type Filter struct {
	Name   string
	Types  []zone_domain.RecordType
	Prefix string
	CIDR   *net.IPNet
	Status zone_domain.Status
}

func Apply(records []zone_domain.RecordSet, f Filter) []zone_domain.RecordSet {
	out := []zone_domain.RecordSet{}
	types := map[zone_domain.RecordType]bool{}
	for _, t := range f.Types {
		types[t] = true
	}
	for _, r := range records {
		if f.Name != "" && !SameName(r.Name, f.Name) {
			continue
		}
		if len(types) > 0 && !types[r.Type] {
			continue
		}
		if f.Prefix != "" && !strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(f.Prefix)) {
			continue
		}
		if f.Status != "" && r.Status != f.Status {
			continue
		}
		if f.CIDR != nil {
			ok := false
			for _, v := range r.Values {
				if ip := net.ParseIP(v.Value); ip != nil && f.CIDR.Contains(ip) {
					ok = true
				}
			}
			if !ok {
				continue
			}
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func GroupByType(rs []zone_domain.RecordSet) map[zone_domain.RecordType][]zone_domain.RecordSet {
	m := map[zone_domain.RecordType][]zone_domain.RecordSet{}
	for _, r := range rs {
		m[r.Type] = append(m[r.Type], r)
	}
	return m
}
func GroupByName(rs []zone_domain.RecordSet) map[string][]zone_domain.RecordSet {
	m := map[string][]zone_domain.RecordSet{}
	for _, r := range rs {
		m[CanonicalName(r.Name)] = append(m[CanonicalName(r.Name)], r)
	}
	return m
}
