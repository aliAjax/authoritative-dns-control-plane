package protocol

import (
	"errors"
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

type flushErrorWriter struct{}

func (flushErrorWriter) Write(p []byte) (int, error) { return 0, errors.New("flush sink unavailable") }

func TestEncodeTransferReportsWriterError(t *testing.T) {
	err := EncodeTransfer(flushErrorWriter{}, Transfer{Zone: "example.com.", Serial: 7, Records: []zone_domain.RecordSet{{Name: "www.example.com.", TTL: 60, Type: zone_domain.A, Values: []zone_domain.RecordValue{{Value: "192.0.2.1"}}}}})
	if err == nil {
		t.Fatal("writer flush error was lost")
	}
}
