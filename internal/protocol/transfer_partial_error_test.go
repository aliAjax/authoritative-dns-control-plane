package protocol

import (
	"errors"
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

type partialTransferWriter struct{}

func (partialTransferWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) / 2, errors.New("partial transfer sink")
}

func TestEncodeTransferReportsPartialWriterError(t *testing.T) {
	err := EncodeTransfer(partialTransferWriter{}, Transfer{Zone: "example.com.", Serial: 7, Records: []zone_domain.RecordSet{{Name: "www.example.com.", TTL: 60, Type: zone_domain.A, Values: []zone_domain.RecordValue{{Value: "192.0.2.1"}}}}})
	if err == nil {
		t.Fatal("partial writer error was lost")
	}
}
