package catalog

import "strings"

type TypeMetadata struct {
	Code            uint16
	Name            string
	RData           string
	SupportsRouting bool
	DNSSEC          bool
}

var metadata = map[string]TypeMetadata{
	"A": {1, "A", "ipv4", true, true}, "AAAA": {28, "AAAA", "ipv6", true, true}, "CNAME": {5, "CNAME", "name", false, true},
	"MX": {15, "MX", "priority target", false, true}, "TXT": {16, "TXT", "character-string", false, true}, "SRV": {33, "SRV", "priority weight port target", true, true},
	"CAA": {257, "CAA", "flags tag value", false, true}, "HTTPS": {65, "HTTPS", "priority target params", true, true}, "SVCB": {64, "SVCB", "priority target params", true, true},
	"NS": {2, "NS", "name", false, true}, "SOA": {6, "SOA", "authority tuple", false, true},
}

func Metadata(name string) (TypeMetadata, bool) {
	m, ok := metadata[strings.ToUpper(name)]
	return m, ok
}
func SupportedTypes() []string {
	r := make([]string, 0, len(metadata))
	for k := range metadata {
		r = append(r, k)
	}
	return r
}
func IsAddressType(name string) bool {
	return strings.EqualFold(name, "A") || strings.EqualFold(name, "AAAA")
}
func IsAliasType(name string) bool {
	return strings.EqualFold(name, "CNAME") || strings.EqualFold(name, "HTTPS") || strings.EqualFold(name, "SVCB")
}
func IsRoutingType(name string) bool  { m, _ := Metadata(name); return m.SupportsRouting }
func RequiresDNSSEC(name string) bool { m, _ := Metadata(name); return m.DNSSEC }
