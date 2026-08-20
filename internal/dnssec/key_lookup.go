package dnssec

import "fmt"

func lookupSigningKey(keys map[string]Key, keyID string) (Key, error) {
	k, ok := keys[keyID]
	if !ok {
		return Key{}, fmt.Errorf("key %q not found", keyID)
	}
	if k.ZoneID == "" || k.Role == "" {
		return Key{}, fmt.Errorf("key %q is incomplete", keyID)
	}
	return k, nil
}
