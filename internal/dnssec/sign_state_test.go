package dnssec

import "testing"

func TestSignRejectsMissingKey(t *testing.T) {
	m := NewManager()
	if _, err := m.Sign("missing", "www.example.com.", "A", "192.0.2.5"); err == nil {
		t.Fatal("signing with a missing key must fail")
	}
}

func TestSignRejectsUnsupportedAlgorithm(t *testing.T) {
	m := NewManager()
	if err := m.Add(Key{ID: "bad-alg", ZoneID: "example.com.", Role: "zsk", Algorithm: 8, Public: "pub", State: Active}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Sign("bad-alg", "www.example.com.", "A", "192.0.2.5"); err == nil {
		t.Fatal("unsupported signing algorithm must fail")
	}
}
