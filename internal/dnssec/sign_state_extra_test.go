package dnssec

import "testing"

func TestSignRejectsIncompleteKey(t *testing.T) {
	m := NewManager()
	if err := m.Add(Key{ID: "incomplete", ZoneID: "example.com.", Role: "zsk", Algorithm: 0, State: Active}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Sign("incomplete", "www.example.com.", "A", "192.0.2.5"); err == nil {
		t.Fatal("incomplete signing key must fail")
	}
}
