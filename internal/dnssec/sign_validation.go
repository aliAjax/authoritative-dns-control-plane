package dnssec

import "fmt"

func validateSigningKey(k Key) error {
	if k.State != Active {
		return fmt.Errorf("key %q is not active", k.ID)
	}
	if k.Algorithm == 0 {
		return fmt.Errorf("key %q has no signing algorithm", k.ID)
	}
	if k.Algorithm != 13 {
		return fmt.Errorf("key %q uses unsupported algorithm %d", k.ID, k.Algorithm)
	}
	return nil
}
