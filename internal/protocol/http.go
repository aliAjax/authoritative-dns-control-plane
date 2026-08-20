package protocol

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HTTP    *http.Client
	Retries int
	Timeout time.Duration
}

func NewClient() *Client { return &Client{HTTP: &http.Client{Timeout: 10 * time.Second}, Retries: 3} }
func (c *Client) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, int, error) {
	if c.HTTP == nil {
		c.HTTP = &http.Client{Timeout: c.Timeout}
	}
	if c.Retries < 1 {
		c.Retries = 1
	}
	var last error
	for i := 0; i < c.Retries; i++ {
		if err := contextError(ctx); err != nil {
			return nil, 0, err
		}
		req, e := http.NewRequestWithContext(ctx, method, url, body)
		if e != nil {
			return nil, 0, e
		}
		resp, e := c.HTTP.Do(req)
		if e != nil {
			last = wrapRequestError(i+1, e)
			continue
		}
		b, e := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		resp.Body.Close()
		if e != nil {
			last = wrapReadError(i+1, e)
			continue
		}
		if resp.StatusCode >= 500 {
			last = upstreamStatusError(i+1, resp.StatusCode)
			continue
		}
		return b, resp.StatusCode, nil
	}
	return nil, 0, last
}
func IsRetryable(code int) bool { return code == 408 || code == 425 || code == 429 || code >= 500 }
