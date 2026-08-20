package operations

import "testing"

func TestLimiterSeparatesClientBudgets(t *testing.T) {
	l := NewLimiter(1, 1)
	if !l.Allow("edge-a") || l.Allow("edge-a") {
		t.Fatal("first client budget was not enforced")
	}
	if !l.Allow("edge-b") || l.Size() != 2 {
		t.Fatal("second client did not receive its own budget")
	}
}
