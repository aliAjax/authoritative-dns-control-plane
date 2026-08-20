package operations

import (
	"fmt"
	"sync"
	"time"
)

type State string

const (
	StatePending   State = "pending"
	StateRunning   State = "running"
	StateSucceeded State = "succeeded"
	StateFailed    State = "failed"
	StateCancelled State = "cancelled"
)

type Transition struct {
	From   State
	To     State
	At     time.Time
	Reason string
}
type StateMachine struct {
	mu      sync.Mutex
	state   State
	history []Transition
}

func NewStateMachine(initial State) *StateMachine { return &StateMachine{state: initial} }

var allowed = map[State]map[State]bool{StatePending: {StateRunning: true, StateCancelled: true}, StateRunning: {StateSucceeded: true, StateFailed: true, StateCancelled: true}, StateFailed: {StateRunning: true}, StateSucceeded: {}, StateCancelled: {}}

func (m *StateMachine) Move(to State, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !allowed[m.state][to] {
		return fmt.Errorf("invalid transition %s -> %s", m.state, to)
	}
	m.history = append(m.history, Transition{From: m.state, To: to, At: time.Now().UTC(), Reason: reason})
	m.state = to
	return nil
}
func (m *StateMachine) Current() State { m.mu.Lock(); defer m.mu.Unlock(); return m.state }
func (m *StateMachine) History() []Transition {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]Transition(nil), m.history...)
}
