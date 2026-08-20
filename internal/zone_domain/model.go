package zone_domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Status string

const (
	Draft      Status = "draft"
	Reviewing  Status = "reviewing"
	Published  Status = "published"
	Superseded Status = "superseded"
	RolledBack Status = "rolled_back"
)

type RecordType string

const (
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	CNAME RecordType = "CNAME"
	MX    RecordType = "MX"
	TXT   RecordType = "TXT"
	SRV   RecordType = "SRV"
	CAA   RecordType = "CAA"
	HTTPS RecordType = "HTTPS"
	SVCB  RecordType = "SVCB"
	SOA   RecordType = "SOA"
	NS    RecordType = "NS"
)

type RecordValue struct {
	Value    string `json:"value"`
	Priority uint16 `json:"priority,omitempty"`
	Weight   uint16 `json:"weight,omitempty"`
	Port     uint16 `json:"port,omitempty"`
	Target   string `json:"target,omitempty"`
}
type RecordSet struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Type      RecordType    `json:"type"`
	TTL       uint32        `json:"ttl"`
	Values    []RecordValue `json:"values"`
	Status    Status        `json:"status"`
	Version   int64         `json:"version"`
	UpdatedAt time.Time     `json:"updated_at"`
	PolicyID  string        `json:"policy_id,omitempty"`
}
type Zone struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Serial    uint32      `json:"serial"`
	SOA       RecordSet   `json:"soa"`
	NS        []string    `json:"ns"`
	DNSSEC    DNSSECState `json:"dnssec"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
type DNSSECState struct {
	KSKVersion string `json:"ksk_version"`
	ZSKVersion string `json:"zsk_version"`
	Signing    bool   `json:"signing"`
	NSEC3      bool   `json:"nsec3"`
}
type Snapshot struct {
	ID        string      `json:"id"`
	ZoneID    string      `json:"zone_id"`
	Serial    uint32      `json:"serial"`
	Records   []RecordSet `json:"records"`
	Digest    string      `json:"digest"`
	CreatedAt time.Time   `json:"created_at"`
	Immutable bool        `json:"immutable"`
}

func NewZone(name string) Zone {
	now := time.Now().UTC()
	id := hash(name + now.String())[:16]
	return Zone{ID: id, Name: name, Serial: uint32(now.Unix()), CreatedAt: now, UpdatedAt: now, NS: []string{"ns1." + name, "ns2." + name}}
}
func (z Zone) ValidateRecord(r RecordSet) error {
	if strings.TrimSpace(r.Name) == "" || r.TTL == 0 {
		return fmt.Errorf("record name and ttl required")
	}
	if !strings.HasSuffix(r.Name, ".") {
		r.Name += "."
	}
	switch r.Type {
	case A, AAAA, CNAME, MX, TXT, SRV, CAA, HTTPS, SVCB, SOA, NS:
	default:
		return fmt.Errorf("unsupported record type %s", r.Type)
	}
	if len(r.Values) == 0 {
		return fmt.Errorf("record values required")
	}
	return nil
}
func SnapshotDigest(records []RecordSet) string {
	cp := append([]RecordSet(nil), records...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].Name+string(cp[i].Type) < cp[j].Name+string(cp[j].Type) })
	h := sha256.New()
	for _, r := range cp {
		fmt.Fprintf(h, "%s|%s|%d|", r.Name, r.Type, r.TTL)
		for _, v := range r.Values {
			fmt.Fprintf(h, "%s:%d:%d:%d;", v.Value, v.Priority, v.Weight, v.Port)
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}
func hash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
