package policy_engine

import "testing"

func TestEvaluateRejectsInvalidCIDRRule(t *testing.T) {
	_, err := Evaluate(Policy{Rules: []Rule{{ID: "bad", CIDR: "not-a-cidr", Targets: []string{"192.0.2.1"}}}}, Query{ClientIP: nil})
	if err == nil {
		t.Fatal("invalid CIDR must be reported to caller")
	}
}
