package operations

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func RandomBytes(n int) ([]byte, error) {
	if n < 1 {
		return nil, fmt.Errorf("size must be positive")
	}
	b := make([]byte, n)
	_, e := rand.Read(b)
	return b, e
}
func RandomToken(n int) (string, error) {
	b, e := RandomBytes(n)
	if e != nil {
		return "", e
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func Fingerprint(secret string) string {
	h := sha256.Sum256([]byte(secret))
	return fmt.Sprintf("sha256:%x", h)
}
func Redact(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
func ConstantTimeEqual(a, b string) bool {
	ha := sha256.Sum256([]byte(a))
	hb := sha256.Sum256([]byte(b))
	var x byte
	for i := range ha {
		x |= ha[i] ^ hb[i]
	}
	return x == 0
}
