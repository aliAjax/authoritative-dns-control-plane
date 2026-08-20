package catalog

import "testing"

func TestSerialAllocatorRejectsRegression(t *testing.T) {
	a := NewSerialAllocator()
	if err := a.Reserve("example.com.", 100); err != nil {
		t.Fatal(err)
	}
	if err := a.Reserve("example.com.", 100); err == nil {
		t.Fatal("equal serial must be rejected")
	}
	if got := a.Next("example.com."); got <= 100 {
		t.Fatalf("serial did not advance: %d", got)
	}
}
