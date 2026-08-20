package protocol

import "fmt"

func upstreamStatusError(attempt int, status int) error {
	return fmt.Errorf("upstream status %d on attempt %d", status, attempt)
}
