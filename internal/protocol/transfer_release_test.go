package protocol

import (
	"bytes"
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

func TestEncodeTransferFlushesHeaderAndRecords(t *testing.T) {
	var b bytes.Buffer
	err := EncodeTransfer(&b, Transfer{Zone: "example.com.", Serial: 42, Records: []zone_domain.RecordSet{{Name: "www.example.com.", TTL: 60, Type: zone_domain.A, Values: []zone_domain.RecordValue{{Value: "192.0.2.1"}}}}})
	if err != nil || !bytes.Contains(b.Bytes(), []byte("$SERIAL 42")) || !bytes.Contains(b.Bytes(), []byte("192.0.2.1")) {
		t.Fatalf("incomplete transfer output: %q err=%v", b.Bytes(), err)
	}
}
