package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type AuditEvent struct {
	ID           string         `json:"id"`
	Actor        string         `json:"actor"`
	Action       string         `json:"action"`
	Resource     string         `json:"resource"`
	Payload      map[string]any `json:"payload"`
	PreviousHash string         `json:"previous_hash"`
	Hash         string         `json:"hash"`
	At           time.Time      `json:"at"`
}
type AuditLog struct {
	mu     sync.Mutex
	events []AuditEvent
}

func NewAuditLog() *AuditLog { return &AuditLog{events: make([]AuditEvent, 0, 256)} }
func (l *AuditLog) Append(actor, action, resource string, p map[string]any) AuditEvent {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := AuditEvent{ID: fmt.Sprintf("audit-%d", time.Now().UnixNano()), Actor: actor, Action: action, Resource: resource, Payload: p, At: time.Now().UTC()}
	if len(l.events) > 0 {
		e.PreviousHash = l.events[len(l.events)-1].Hash
	}
	raw, _ := json.Marshal(e)
	h := sha256.Sum256(raw)
	e.Hash = hex.EncodeToString(h[:])
	l.events = append(l.events, e)
	return e
}
func (l *AuditLog) Verify() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	prev := ""
	for _, e := range l.events {
		if e.PreviousHash != prev {
			return fmt.Errorf("audit chain broken at %s", e.ID)
		}
		copy := e
		copy.Hash = ""
		raw, _ := json.Marshal(copy)
		h := sha256.Sum256(raw)
		if hex.EncodeToString(h[:]) != e.Hash {
			return fmt.Errorf("audit hash mismatch at %s", e.ID)
		}
		prev = e.Hash
	}
	return nil
}
func (l *AuditLog) List() []AuditEvent {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]AuditEvent(nil), l.events...)
}
func (l *AuditLog) Since(t time.Time) []AuditEvent {
	out := []AuditEvent{}
	for _, e := range l.List() {
		if e.At.After(t) {
			out = append(out, e)
		}
	}
	return out
}
func (l *AuditLog) Count() int { return len(l.List()) }
func (l *AuditLog) LastHash() string {
	e := l.List()
	if len(e) == 0 {
		return ""
	}
	return e[len(e)-1].Hash
}
