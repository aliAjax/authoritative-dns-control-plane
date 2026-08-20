package dns_wire

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Header struct {
	ID          uint16
	Flags       uint16
	Questions   uint16
	Answers     uint16
	Authorities uint16
	Additionals uint16
}
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}
type RR struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []byte
}
type Message struct {
	Header      Header
	Questions   []Question
	Answers     []RR
	Authorities []RR
	Additionals []RR
}

func Parse(b []byte) (Message, error) {
	if len(b) < 12 {
		return Message{}, fmt.Errorf("dns header truncated")
	}
	m := Message{Header: Header{ID: binary.BigEndian.Uint16(b), Flags: binary.BigEndian.Uint16(b[2:]), Questions: binary.BigEndian.Uint16(b[4:]), Answers: binary.BigEndian.Uint16(b[6:]), Authorities: binary.BigEndian.Uint16(b[8:]), Additionals: binary.BigEndian.Uint16(b[10:])}}
	off := 12
	for i := 0; i < int(m.Header.Questions); i++ {
		name, n, err := decodeName(b, off)
		if err != nil {
			return Message{}, err
		}
		off = n
		if off+4 > len(b) {
			return Message{}, fmt.Errorf("question truncated")
		}
		m.Questions = append(m.Questions, Question{Name: name, Type: binary.BigEndian.Uint16(b[off:]), Class: binary.BigEndian.Uint16(b[off+2:])})
		off += 4
	}
	return m, nil
}
func Encode(m Message) []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf, m.Header.ID)
	binary.BigEndian.PutUint16(buf[2:], m.Header.Flags)
	binary.BigEndian.PutUint16(buf[4:], uint16(len(m.Questions)))
	binary.BigEndian.PutUint16(buf[6:], uint16(len(m.Answers)))
	binary.BigEndian.PutUint16(buf[8:], uint16(len(m.Authorities)))
	binary.BigEndian.PutUint16(buf[10:], uint16(len(m.Additionals)))
	for _, q := range m.Questions {
		buf = append(buf, encodeName(q.Name)...)
		x := make([]byte, 4)
		binary.BigEndian.PutUint16(x, q.Type)
		binary.BigEndian.PutUint16(x[2:], q.Class)
		buf = append(buf, x...)
	}
	for _, rr := range m.Answers {
		buf = append(buf, encodeName(rr.Name)...)
		x := make([]byte, 10)
		binary.BigEndian.PutUint16(x, rr.Type)
		binary.BigEndian.PutUint16(x[2:], rr.Class)
		binary.BigEndian.PutUint32(x[4:], rr.TTL)
		binary.BigEndian.PutUint16(x[8:], uint16(len(rr.Data)))
		buf = append(buf, x...)
		buf = append(buf, rr.Data...)
	}
	return buf
}
func decodeName(b []byte, off int) (string, int, error) {
	var labels []string
	seen := map[int]bool{}
	return decodeNameRec(b, off, &labels, seen, 0)
}
func decodeNameRec(b []byte, off int, labels *[]string, seen map[int]bool, depth int) (string, int, error) {
	if depth > 32 {
		return "", 0, fmt.Errorf("compression pointer depth exceeded")
	}
	start := off
	for {
		if off >= len(b) {
			return "", 0, fmt.Errorf("name truncated")
		}
		n := int(b[off])
		if n == 0 {
			off++
			return strings.Join(*labels, "."), off, nil
		}
		if n&0xc0 == 0xc0 {
			if off+1 >= len(b) {
				return "", 0, fmt.Errorf("pointer truncated")
			}
			p := int(binary.BigEndian.Uint16(b[off:off+2]) & 0x3fff)
			if seen[p] {
				return "", 0, fmt.Errorf("compression pointer loop")
			}
			seen[p] = true
			s, _, err := decodeNameRec(b, p, labels, seen, depth+1)
			return s, start + 2, err
		}
		if n > 63 || off+1+n > len(b) {
			return "", 0, fmt.Errorf("invalid label")
		}
		*labels = append(*labels, string(b[off+1:off+1+n]))
		off += 1 + n
	}
}
func encodeName(name string) []byte {
	out := []byte{}
	name = strings.TrimSuffix(name, ".")
	for _, l := range strings.Split(name, ".") {
		if l == "" {
			continue
		}
		out = append(out, byte(len(l)))
		out = append(out, []byte(l)...)
	}
	return append(out, 0)
}
func TypeCode(t string) uint16 {
	switch strings.ToUpper(t) {
	case "A":
		return 1
	case "NS":
		return 2
	case "CNAME":
		return 5
	case "SOA":
		return 6
	case "MX":
		return 15
	case "TXT":
		return 16
	case "AAAA":
		return 28
	case "SRV":
		return 33
	case "CAA":
		return 257
	case "HTTPS":
		return 65
	case "SVCB":
		return 64
	}
	return 255
}
func (m Message) QuestionName() string {
	if len(m.Questions) == 0 {
		return ""
	}
	return m.Questions[0].Name
}
