package dnssec

import "testing"

func TestSignRejectsMalformedRequest(t *testing.T) {
	m := NewManager()
	if err := m.Add(Key{ID: "valid", ZoneID: "example.com.", Role: "zsk", Algorithm: 13, Public: "pub", State: Active}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Sign("valid", "", "A", "192.0.2.5"); err == nil {
		t.Fatal("empty signature name must fail")
	}
}
