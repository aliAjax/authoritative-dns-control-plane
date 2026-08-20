package protocol

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestClientDoLabelsStatusRetry(t *testing.T) {
	c := &Client{HTTP: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadGateway, Body: http.NoBody}, nil
	})}, Retries: 2}
	_, _, err := c.Do(context.Background(), http.MethodGet, "http://resolver.invalid", nil)
	if err == nil || !strings.Contains(err.Error(), "attempt 2") {
		t.Fatalf("status retry lacked attempt context: %v", err)
	}
}
