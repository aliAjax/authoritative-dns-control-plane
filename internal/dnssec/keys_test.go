package dnssec

import (
	"testing"
	"time"
)

func TestRotationCanActivateAndSign(t *testing.T) {
	m := NewManager()
	now := time.Now().UTC()
	if err := m.Add(Key{ID: "old-zsk", ZoneID: "example.com.", Role: "zsk", State: Active}); err != nil {
		t.Fatal(err)
	}
	k, err := m.Rotate("example.com.", "zsk", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Activate(k.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Sign(k.ID, "www.example.com.", "A", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
}
