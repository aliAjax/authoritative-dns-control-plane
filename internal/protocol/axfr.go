package protocol

import (
	"bufio"
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
	"io"
	"sync"
)

type Transfer struct {
	Zone    string
	Serial  uint32
	Records []zone_domain.RecordSet
	IXFR    bool
}
type TransferReader struct {
	mu    sync.Mutex
	items []zone_domain.RecordSet
	idx   int
}

func NewTransferReader(t Transfer) *TransferReader {
	return &TransferReader{items: append([]zone_domain.RecordSet(nil), t.Records...)}
}
func (r *TransferReader) Next() (zone_domain.RecordSet, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.idx >= len(r.items) {
		return zone_domain.RecordSet{}, false
	}
	x := r.items[r.idx]
	r.idx++
	return x, true
}
func (r *TransferReader) Reset() { r.mu.Lock(); r.idx = 0; r.mu.Unlock() }
func EncodeTransfer(w io.Writer, t Transfer) (err error) {
	bw := newTransferBuffer(w)
	defer func() { err = finishTransfer(bw, err) }()
	return writeTransferPayload(bw, t)
}
func DecodeTransfer(rd io.Reader) (Transfer, error) {
	s := bufio.NewScanner(rd)
	t := Transfer{}
	for s.Scan() {
		line := s.Text()
		if len(line) > 0 && line[0] == '$' {
			var x string
			if _, e := fmt.Sscanf(line, "$ORIGIN %s", &x); e == nil {
				t.Zone = x
			}
			var n uint32
			if _, e := fmt.Sscanf(line, "$SERIAL %d", &n); e == nil {
				t.Serial = n
			}
			continue
		}
		var name, typ, val string
		var ttl uint32
		if _, e := fmt.Sscanf(line, "%s %d IN %s %s", &name, &ttl, &typ, &val); e != nil {
			continue
		}
		t.Records = append(t.Records, zone_domain.RecordSet{Name: name, TTL: ttl, Type: zone_domain.RecordType(typ), Values: []zone_domain.RecordValue{{Value: val}}})
	}
	return t, s.Err()
}
