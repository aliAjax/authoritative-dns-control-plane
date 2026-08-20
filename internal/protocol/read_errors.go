package protocol

import "fmt"

func wrapReadError(attempt int, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("response read attempt %d: %w", attempt, err)
}
