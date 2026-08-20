package protocol

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestClientDoPreservesReadError(t *testing.T) {
	want := errors.New("response body truncated")
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(failingReader{err: want})}, nil
	})}, Retries: 1}
	_, _, err := c.Do(context.Background(), http.MethodGet, "http://resolver.invalid", nil)
	if !errors.Is(err, want) {
		t.Fatalf("read error was not preserved: %v", err)
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }
