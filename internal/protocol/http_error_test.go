package protocol

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestClientDoPreservesTransportError(t *testing.T) {
	want := errors.New("resolver connection refused")
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, want })}, Retries: 2}
	_, _, err := c.Do(context.Background(), http.MethodGet, "http://resolver.invalid", nil)
	if !errors.Is(err, want) {
		t.Fatalf("transport error was not preserved: %v", err)
	}
}

func TestClientDoStopsBeforeRetryWhenCanceled(t *testing.T) {
	called := 0
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		called++
		return nil, errors.New("should not run")
	})}, Retries: 2}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err := c.Do(ctx, http.MethodGet, "http://resolver.invalid", nil)
	if !errors.Is(err, context.Canceled) || called != 0 {
		t.Fatalf("canceled request was retried: err=%v calls=%d", err, called)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

var _ io.Reader
