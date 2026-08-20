package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Option struct {
	Code uint16
	Data []byte
}
type EDNS struct {
	UDPSize uint16
	Version uint8
	DO      bool
	Options []Option
}

func ParseEDNS(data []byte) (EDNS, error) {
	if len(data) < 4 {
		return EDNS{}, fmt.Errorf("edns truncated")
	}
	e := EDNS{UDPSize: binary.BigEndian.Uint16(data), Version: data[2], DO: data[3]&0x80 != 0}
	off := 4
	for off < len(data) {
		if off+4 > len(data) {
			return EDNS{}, fmt.Errorf("option header truncated")
		}
		c := binary.BigEndian.Uint16(data[off:])
		n := int(binary.BigEndian.Uint16(data[off+2:]))
		off += 4
		if off+n > len(data) {
			return EDNS{}, fmt.Errorf("option body truncated")
		}
		e.Options = append(e.Options, Option{Code: c, Data: append([]byte(nil), data[off:off+n]...)})
		off += n
	}
	return e, nil
}
func (e EDNS) Encode() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b, e.UDPSize)
	b[2] = e.Version
	if e.DO {
		b[3] = 0x80
	}
	for _, o := range e.Options {
		x := make([]byte, 4)
		binary.BigEndian.PutUint16(x, o.Code)
		binary.BigEndian.PutUint16(x[2:], uint16(len(o.Data)))
		b = append(b, x...)
		b = append(b, o.Data...)
	}
	return b
}
func ECS(ip net.IP, prefix uint8) Option {
	v4 := ip.To4() != nil
	fam := uint16(1)
	bytes := ip.To4()
	if !v4 {
		fam = 2
		bytes = ip.To16()
	}
	n := (int(prefix) + 7) / 8
	data := make([]byte, 4+n)
	binary.BigEndian.PutUint16(data, fam)
	data[2] = prefix
	copy(data[4:], bytes[:n])
	return Option{Code: 8, Data: data}
}
func DecodeECS(o Option) (net.IP, uint8, error) {
	if o.Code != 8 || len(o.Data) < 4 {
		return nil, 0, fmt.Errorf("invalid ecs")
	}
	fam := binary.BigEndian.Uint16(o.Data)
	p := o.Data[2]
	if fam == 1 {
		if p > 32 {
			return nil, 0, fmt.Errorf("invalid ipv4 ecs prefix")
		}
		n := (int(p) + 7) / 8
		if len(o.Data) != 4+n {
			return nil, 0, fmt.Errorf("invalid ipv4 ecs address length")
		}
		ip := make(net.IP, net.IPv4len)
		copy(ip, o.Data[4:])
		return ip, p, nil
	}
	if fam == 2 {
		if p > 128 {
			return nil, 0, fmt.Errorf("invalid ipv6 ecs prefix")
		}
		n := (int(p) + 7) / 8
		if len(o.Data) != 4+n {
			return nil, 0, fmt.Errorf("invalid ipv6 ecs address length")
		}
		ip := make(net.IP, 16)
		copy(ip, o.Data[4:])
		return ip, p, nil
	}
	return nil, 0, fmt.Errorf("unknown family")
}
