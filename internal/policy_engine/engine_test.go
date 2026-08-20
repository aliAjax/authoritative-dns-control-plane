package policy_engine

import "testing"

func TestEvaluateSkipsUnhealthyFailoverRule(t *testing.T) {
	p := Policy{Algorithm: "first", Rules: []Rule{
		{ID: "primary", Priority: 1, Targets: []string{"192.0.2.1"}, Failover: true},
		{ID: "secondary", Priority: 2, Targets: []string{"192.0.2.2"}},
	}}
	d, err := Evaluate(p, Query{Healthy: map[string]bool{"192.0.2.1": false, "192.0.2.2": true}})
	if err != nil {
		t.Fatal(err)
	}
	if d.RuleID != "secondary" || len(d.Values) != 1 || d.Values[0].Value != "192.0.2.2" {
		t.Fatalf("unexpected decision: %#v", d)
	}
}
