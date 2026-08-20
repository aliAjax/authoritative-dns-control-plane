package policy_engine

import "testing"

func TestEvaluateRejectsRuleWithoutID(t *testing.T) {
	_, err := Evaluate(Policy{Rules: []Rule{{Targets: []string{"192.0.2.1"}}}}, Query{})
	if err == nil {
		t.Fatal("unnamed rule must be rejected before selection")
	}
}
