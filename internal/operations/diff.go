package operations

import (
	"fmt"
	"reflect"
)

type ChangeKind string

const (
	Add    ChangeKind = "add"
	Remove ChangeKind = "remove"
	Update ChangeKind = "update"
)

type Change[T any] struct {
	Kind   ChangeKind
	Key    string
	Before T
	After  T
}

func Diff[T any](before, after map[string]T) []Change[T] {
	out := []Change[T]{}
	for k, v := range before {
		n, ok := after[k]
		if !ok {
			out = append(out, Change[T]{Kind: Remove, Key: k, Before: v})
		} else if !reflect.DeepEqual(v, n) {
			out = append(out, Change[T]{Kind: Update, Key: k, Before: v, After: n})
		}
	}
	for k, v := range after {
		if _, ok := before[k]; !ok {
			out = append(out, Change[T]{Kind: Add, Key: k, After: v})
		}
	}
	return out
}
func ChangeSummary[T any](c []Change[T]) string {
	a, r, u := 0, 0, 0
	for _, x := range c {
		switch x.Kind {
		case Add:
			a++
		case Remove:
			r++
		case Update:
			u++
		}
	}
	return fmt.Sprintf("add=%d remove=%d update=%d", a, r, u)
}
