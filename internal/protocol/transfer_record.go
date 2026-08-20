package protocol

import (
	"bufio"
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
)

func writeTransferRecords(w *bufio.Writer, records []zone_domain.RecordSet) error {
	for _, r := range records {
		for _, v := range r.Values {
			if _, err := writeTransferRecord(w, r, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeTransferRecord(w *bufio.Writer, r zone_domain.RecordSet, v zone_domain.RecordValue) (int, error) {
	return fmt.Fprintf(w, "%s %d IN %s %s\n", r.Name, r.TTL, r.Type, v.Value)
}
