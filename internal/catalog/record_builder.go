package catalog

import (
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Builder struct {
	record zone_domain.RecordSet
	errors []error
}

func NewBuilder(name string, t zone_domain.RecordType) *Builder {
	return &Builder{record: zone_domain.RecordSet{Name: name, Type: t, TTL: 300}}
}
func (b *Builder) TTL(v uint32) *Builder {
	if v == 0 || v > 604800 {
		b.errors = append(b.errors, fmt.Errorf("ttl out of range"))
	} else {
		b.record.TTL = v
	}
	return b
}
func (b *Builder) Value(v string) *Builder {
	if strings.TrimSpace(v) == "" {
		b.errors = append(b.errors, fmt.Errorf("empty value"))
	} else {
		b.record.Values = append(b.record.Values, zone_domain.RecordValue{Value: v})
	}
	return b
}
func (b *Builder) Address(v string) *Builder {
	if net.ParseIP(v) == nil {
		b.errors = append(b.errors, fmt.Errorf("invalid address"))
	} else {
		return b.Value(v)
	}
	return b
}
func (b *Builder) Target(v string) *Builder {
	if !strings.HasSuffix(v, ".") {
		v += "."
	}
	return b.Value(v)
}
func (b *Builder) Priority(v uint16) *Builder {
	b.record.Values = append(b.record.Values, zone_domain.RecordValue{Priority: v})
	return b
}
func (b *Builder) Weight(v uint16) *Builder {
	b.record.Values = append(b.record.Values, zone_domain.RecordValue{Weight: v})
	return b
}
func (b *Builder) Port(v uint16) *Builder {
	b.record.Values = append(b.record.Values, zone_domain.RecordValue{Port: v})
	return b
}
func (b *Builder) Policy(id string) *Builder { b.record.PolicyID = id; return b }
func (b *Builder) Build() (zone_domain.RecordSet, error) {
	if b.record.Name == "" {
		b.errors = append(b.errors, fmt.Errorf("name required"))
	}
	if len(b.record.Values) == 0 {
		b.errors = append(b.errors, fmt.Errorf("value required"))
	}
	if len(b.errors) > 0 {
		return zone_domain.RecordSet{}, fmt.Errorf("record invalid: %v", b.errors)
	}
	return b.record, nil
}
func ParseTXT(v string) []string {
	parts := strings.Fields(v)
	out := []string{}
	for _, p := range parts {
		if len(p) > 255 {
			for len(p) > 255 {
				out = append(out, p[:255])
				p = p[255:]
			}
		}
		out = append(out, p)
	}
	return out
}
func EncodeMX(priority uint16, target string) string {
	return strconv.Itoa(int(priority)) + " " + strings.TrimSuffix(target, ".") + "."
}
func EncodeSRV(priority, weight, port uint16, target string) string {
	return fmt.Sprintf("%d %d %d %s", priority, weight, port, strings.TrimSuffix(target, ".")+".")
}
func EncodeCAA(flags uint8, tag, value string) string {
	return fmt.Sprintf("%d %s %q", flags, tag, value)
}
func EncodeSOA(mname, rname string, serial, refresh, retry, expire, minimum uint32) string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", mname, rname, serial, refresh, retry, expire, minimum)
}
func DecodeMX(v string) (uint16, string, error) {
	p := strings.Fields(v)
	if len(p) != 2 {
		return 0, "", fmt.Errorf("invalid mx")
	}
	n, e := strconv.Atoi(p[0])
	if e != nil || n < 0 || n > 65535 {
		return 0, "", fmt.Errorf("invalid priority")
	}
	return uint16(n), p[1], nil
}
func DecodeSRV(v string) (uint16, uint16, uint16, string, error) {
	p := strings.Fields(v)
	if len(p) != 4 {
		return 0, 0, 0, "", fmt.Errorf("invalid srv")
	}
	n := make([]uint16, 3)
	for i := range n {
		x, e := strconv.Atoi(p[i])
		if e != nil || x < 0 || x > 65535 {
			return 0, 0, 0, "", fmt.Errorf("invalid srv number")
		}
		n[i] = uint16(x)
	}
	return n[0], n[1], n[2], p[3], nil
}
func CanonicalName(n string) string {
	return strings.ToLower(strings.TrimSuffix(strings.TrimSpace(n), ".")) + "."
}
func SameName(a, b string) bool { return CanonicalName(a) == CanonicalName(b) }
func IsSubdomain(name, zone string) bool {
	n, z := CanonicalName(name), CanonicalName(zone)
	return n == z || strings.HasSuffix(n, "."+strings.TrimSuffix(z, "."))
}
