package policy_engine

import (
	"crypto/sha256"
	"encoding/binary"
	"example.com/authoritativedns/internal/repository"
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
	"net"
	"sort"
	"strings"
)

type Rule struct {
	ID       string   `json:"id"`
	Region   string   `json:"region,omitempty"`
	ASN      uint32   `json:"asn,omitempty"`
	Provider string   `json:"provider,omitempty"`
	CIDR     string   `json:"cidr,omitempty"`
	Weight   int      `json:"weight"`
	Priority int      `json:"priority"`
	Targets  []string `json:"targets"`
	Failover bool     `json:"failover"`
}
type Policy struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Rules     []Rule `json:"rules"`
	Algorithm string `json:"algorithm"`
}
type Query struct {
	ZoneID   string
	Name     string
	Type     zone_domain.RecordType
	Region   string
	ASN      uint32
	Provider string
	ClientIP net.IP
	ECS      net.IP
	Healthy  map[string]bool
}
type Decision struct {
	Values      []zone_domain.RecordValue `json:"values"`
	RuleID      string                    `json:"rule_id"`
	Explanation string                    `json:"explanation"`
	TTL         uint32                    `json:"ttl"`
	Negative    bool                      `json:"negative"`
}

func Resolve(s *repository.MemoryStore, q Query) (Decision, error) {
	rs := s.ListRecords(q.ZoneID)
	var matches []zone_domain.RecordSet
	for _, r := range rs {
		if strings.EqualFold(r.Name, q.Name) && r.Type == q.Type {
			matches = append(matches, r)
		}
	}
	if len(matches) == 0 {
		return Decision{Negative: true, Explanation: "no record set matched"}, nil
	}
	r := matches[0]
	d := Decision{Values: r.Values, TTL: r.TTL, Explanation: "default record set"}
	return d, nil
}
func Evaluate(p Policy, q Query) (Decision, error) {
	c := append([]Rule(nil), p.Rules...)
	sort.SliceStable(c, func(i, j int) bool {
		if c[i].Priority == c[j].Priority {
			return c[i].Weight > c[j].Weight
		}
		return c[i].Priority < c[j].Priority
	})
	for _, r := range c {
		if err := validateRule(r); err != nil {
			return Decision{}, err
		}
		ok, why := match(r, q)
		if !ok {
			continue
		}
		targets := filterHealthy(r.Targets, q.Healthy)
		if len(targets) == 0 && r.Failover {
			continue
		}
		if len(targets) == 0 {
			return Decision{RuleID: r.ID, Explanation: why + "; no healthy targets", Negative: true}, nil
		}
		idx := pick(targets, q, p.Algorithm)
		return Decision{RuleID: r.ID, Values: []zone_domain.RecordValue{{Value: targets[idx]}}, Explanation: why + fmt.Sprintf("; selected target %s", targets[idx])}, nil
	}
	return Decision{Negative: true, Explanation: "no policy rule matched"}, nil
}
func validateRule(r Rule) error {
	if r.ID == "" {
		return fmt.Errorf("policy rule missing id")
	}
	if r.CIDR != "" {
		if _, _, err := net.ParseCIDR(r.CIDR); err != nil {
			return fmt.Errorf("rule %s: invalid cidr %q: %w", r.ID, r.CIDR, err)
		}
	}
	return nil
}
func match(r Rule, q Query) (bool, string) {
	if r.Region != "" && r.Region != q.Region {
		return false, "region mismatch"
	}
	if r.ASN != 0 && r.ASN != q.ASN {
		return false, "asn mismatch"
	}
	if r.Provider != "" && r.Provider != q.Provider {
		return false, "provider mismatch"
	}
	if r.CIDR != "" {
		_, n, e := net.ParseCIDR(r.CIDR)
		if e != nil || q.ClientIP == nil || !n.Contains(q.ClientIP) {
			return false, "cidr mismatch"
		}
	}
	return true, "matched " + r.ID
}
func filterHealthy(in []string, h map[string]bool) []string {
	if len(h) == 0 {
		return in
	}
	o := []string{}
	for _, x := range in {
		if h[x] {
			o = append(o, x)
		}
	}
	return o
}
func pick(t []string, q Query, algo string) int {
	if len(t) == 1 {
		return 0
	}
	if algo == "hash" {
		h := sha256.Sum256(q.ClientIP)
		return int(binary.BigEndian.Uint64(h[:8]) % uint64(len(t)))
	}
	return 0
}
