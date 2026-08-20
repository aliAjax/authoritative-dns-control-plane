package protocol

import (
	"net"
	"testing"
)

func TestEDNSRoundTripPreservesClientSubnet(t *testing.T) {
	in := EDNS{UDPSize: 1232, DO: true, Options: []Option{ECS(net.ParseIP("192.0.2.15"), 24)}}
	out, err := ParseEDNS(in.Encode())
	if err != nil {
		t.Fatal(err)
	}
	ip, prefix, err := DecodeECS(out.Options[0])
	if err != nil || prefix != 24 || !ip.Equal(net.IPv4(192, 0, 2, 0)) {
		t.Fatalf("unexpected ECS result ip=%v prefix=%d err=%v", ip, prefix, err)
	}
}
