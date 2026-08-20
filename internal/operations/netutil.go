package operations

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func NormalizeAddress(s string) string {
	if h, p, e := net.SplitHostPort(s); e == nil {
		return net.JoinHostPort(strings.ToLower(h), p)
	}
	return strings.ToLower(strings.TrimSpace(s))
}
func IsPrivate(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 || ip4[0] == 127 || (ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || (ip4[0] == 192 && ip4[1] == 168)
	}
	return ip.IsLoopback() || ip.IsPrivate()
}
func ParseCIDRs(values []string) ([]*net.IPNet, error) {
	out := make([]*net.IPNet, 0, len(values))
	for _, v := range values {
		_, n, e := net.ParseCIDR(v)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", v, e)
		}
		out = append(out, n)
	}
	return out, nil
}
func ContainsAny(ns []*net.IPNet, ip net.IP) bool {
	for _, n := range ns {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
func RetryAfter(at time.Time) time.Duration {
	d := time.Until(at)
	if d < 0 {
		return 0
	}
	return d
}
func SameSubnet(a, b net.IP, bits int) bool {
	if a == nil || b == nil {
		return false
	}
	width := 128
	if a.To4() != nil && b.To4() != nil {
		width = 32
		a = a.To4()
		b = b.To4()
	}
	if bits < 0 || bits > width {
		return false
	}
	mask := net.CIDRMask(bits, width)
	return a.Mask(mask).Equal(b.Mask(mask))
}
