package policy_engine

import (
	"fmt"
	"net"
)

func validateRule(r Rule) error {
	if r.ID == "" {
		return fmt.Errorf("policy rule id required")
	}
	if r.CIDR == "" {
		return nil
	}
	if _, _, err := net.ParseCIDR(r.CIDR); err != nil {
		return fmt.Errorf("invalid rule %s CIDR: %w", r.ID, err)
	}
	return nil
}
