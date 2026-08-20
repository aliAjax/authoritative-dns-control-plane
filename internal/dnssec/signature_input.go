package dnssec

import "fmt"

func validateSignatureInput(name, typ string) error {
	if name == "" || typ == "" {
		return fmt.Errorf("signature name and type required")
	}
	return nil
}
