package protocol

import (
	"bytes"
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

func TestEncodeTransferFlushesHeaderWithoutRecords(t *testing.T) {
	var b bytes.Buffer
	if err := EncodeTransfer(&b, Transfer{Zone: "example.com.", Serial: 42, Records: []zone_domain.RecordSet{}}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(b.Bytes(), []byte("$ORIGIN example.com.")) {
		t.Fatalf("header was not flushed: %q", b.Bytes())
	}
}
