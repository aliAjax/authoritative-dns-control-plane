package operations

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var labelRE = regexp.MustCompile(`^[A-Za-z0-9_*-]{1,63}$`)

func DomainName(s string) error {
	if len(s) > 253 || s == "" {
		return fmt.Errorf("domain length invalid")
	}
	for _, l := range strings.Split(strings.TrimSuffix(s, "."), ".") {
		if !labelRE.MatchString(l) {
			return fmt.Errorf("invalid label %q", l)
		}
	}
	return nil
}
func IPAddress(s string) error {
	if net.ParseIP(s) == nil {
		return fmt.Errorf("invalid ip")
	}
	return nil
}
func TTL(v uint32) error {
	if v == 0 || v > 604800 {
		return fmt.Errorf("ttl must be 1..604800")
	}
	return nil
}
func Port(v int) error {
	if v < 1 || v > 65535 {
		return fmt.Errorf("port out of range")
	}
	return nil
}
func ASN(v uint32) error {
	if v == 0 || v > 4294967294 {
		return fmt.Errorf("asn out of range")
	}
	return nil
}
func Weight(v int) error {
	if v < 0 || v > 10000 {
		return fmt.Errorf("weight out of range")
	}
	return nil
}
func NonEmpty(name, v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("%s required", name)
	}
	return nil
}
