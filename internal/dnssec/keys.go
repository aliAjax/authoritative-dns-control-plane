package dnssec

import (
	"fmt"
	"sync"
	"time"
)

type KeyState string

const (
	Prepublished KeyState = "prepublished"
	Active       KeyState = "active"
	Retired      KeyState = "retired"
	Revoked      KeyState = "revoked"
)

type Key struct {
	ID         string    `json:"id"`
	ZoneID     string    `json:"zone_id"`
	Role       string    `json:"role"`
	Algorithm  uint16    `json:"algorithm"`
	State      KeyState  `json:"state"`
	Public     string    `json:"public"`
	CreatedAt  time.Time `json:"created_at"`
	ActivateAt time.Time `json:"activate_at"`
	RetireAt   time.Time `json:"retire_at"`
}
type Signature struct {
	KeyID     string    `json:"key_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Digest    string    `json:"digest"`
	ExpiresAt time.Time `json:"expires_at"`
}
type Manager struct {
	mu   sync.Mutex
	keys map[string]Key
	sigs map[string]Signature
}

func NewManager() *Manager { return &Manager{keys: map[string]Key{}, sigs: map[string]Signature{}} }
func (m *Manager) Add(k Key) error {
	if k.ID == "" || k.ZoneID == "" {
		return fmt.Errorf("key id and zone required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keys[k.ID] = k
	return nil
}
func (m *Manager) Rotate(zone, role string, now time.Time) (Key, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range m.keys {
		if k.ZoneID == zone && k.Role == role && k.State == Active {
			k.State = Retired
			m.keys[k.ID] = k
		}
	}
	id := fmt.Sprintf("%s-%s-%d", zone, role, now.UnixNano())
	k := Key{ID: id, ZoneID: zone, Role: role, Algorithm: 13, State: Prepublished, CreatedAt: now, ActivateAt: now.Add(time.Hour)}
	m.keys[id] = k
	return k, nil
}
func (m *Manager) Activate(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k, ok := m.keys[id]
	if !ok {
		return fmt.Errorf("key not found")
	}
	k.State = Active
	m.keys[id] = k
	return nil
}
func (m *Manager) Sign(keyID, name, typ, data string) (Signature, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	k, ok := m.keys[keyID]
	if !ok || k.State != Active {
		return Signature{}, fmt.Errorf("active key required")
	}
	sig := Signature{KeyID: keyID, Name: name, Type: typ, Digest: fmt.Sprintf("sig-%x", []byte(data)), ExpiresAt: time.Now().Add(24 * time.Hour)}
	m.sigs[name+typ] = sig
	return sig, nil
}
