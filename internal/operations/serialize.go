package operations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

func JSON(v any) ([]byte, error)       { return json.Marshal(v) }
func PrettyJSON(v any) ([]byte, error) { return json.MarshalIndent(v, "", "  ") }
func Gzip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	z := gzip.NewWriter(&b)
	if _, e := z.Write(in); e != nil {
		return nil, e
	}
	if e := z.Close(); e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}
func Gunzip(in []byte) ([]byte, error) {
	z, e := gzip.NewReader(bytes.NewReader(in))
	if e != nil {
		return nil, e
	}
	defer z.Close()
	return io.ReadAll(z)
}
func Digest(in []byte) string      { h := sha256.Sum256(in); return fmt.Sprintf("sha256:%x", h) }
func EncodeCursor(v string) string { return base64.RawURLEncoding.EncodeToString([]byte(v)) }
func DecodeCursor(v string) (string, error) {
	b, e := base64.RawURLEncoding.DecodeString(v)
	return string(b), e
}
