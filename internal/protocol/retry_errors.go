package protocol

import "fmt"

func wrapRequestError(attempt int, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("request attempt %d: %w", attempt, err)
}
