package resolver_runtime

import (
	"context"
	"encoding/binary"
	"example.com/authoritativedns/internal/dns_wire"
	"example.com/authoritativedns/internal/policy_engine"
	"example.com/authoritativedns/internal/service"
	"example.com/authoritativedns/internal/zone_domain"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type CacheEntry struct {
	Packet   []byte
	Expires  time.Time
	Negative bool
}
type Cache struct {
	mu       sync.Mutex
	items    map[string]CacheEntry
	capacity int
}

func NewCache(n int) *Cache {
	if n < 1 {
		n = 1000
	}
	return &Cache{items: map[string]CacheEntry{}, capacity: n}
}
func (c *Cache) Get(k string) ([]byte, bool, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[k]
	if !ok {
		return nil, false, false
	}
	if time.Now().After(e.Expires) {
		delete(c.items, k)
		return nil, false, e.Negative
	}
	return e.Packet, true, e.Negative
}
func (c *Cache) Put(k string, p []byte, ttl time.Duration, neg bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.items) >= c.capacity {
		for x := range c.items {
			delete(c.items, x)
			break
		}
	}
	c.items[k] = CacheEntry{Packet: p, Expires: time.Now().Add(ttl), Negative: neg}
}

type Server struct {
	app     *service.Application
	cache   *Cache
	limiter chan struct{}
}

func NewServer(a *service.Application) *Server {
	return &Server{app: a, cache: NewCache(4096), limiter: make(chan struct{}, 1000)}
}
func (s *Server) Resolve(packet []byte) []byte {
	q, err := dns_wire.Parse(packet)
	if err != nil {
		return []byte{0, 0, 0, 0}
	}
	name := q.QuestionName()
	if len(q.Questions) == 0 {
		return packet
	}
	key := fmt.Sprintf("%s/%d", name, q.Questions[0].Type)
	if p, ok, _ := s.cache.Get(key); ok {
		return p
	}
	d, _ := s.app.ResolveDNS(context.Background(), policy_engine.Query{Name: name, Type: wireType(q.Questions[0].Type)})
	resp := dns_wire.Message{Header: dns_wire.Header{ID: q.Header.ID, Flags: 0x8000 | 0x0400}, Questions: q.Questions}
	if d.Negative {
		resp.Header.Flags |= 3
		p := dns_wire.Encode(resp)
		s.cache.Put(key, p, 30*time.Second, true)
		return p
	}
	for _, v := range d.Values {
		resp.Answers = append(resp.Answers, dns_wire.RR{Name: name, Type: q.Questions[0].Type, Class: 1, TTL: d.TTL, Data: []byte(v.Value)})
	}
	p := dns_wire.Encode(resp)
	s.cache.Put(key, p, time.Duration(d.TTL)*time.Second, false)
	return p
}
func (s *Server) ListenAndServeUDP(ctx context.Context, addr string) error {
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer pc.Close()
	go func() { <-ctx.Done(); pc.Close() }()
	buf := make([]byte, 4096)
	for {
		n, peer, e := pc.ReadFrom(buf)
		if e != nil {
			return e
		}
		select {
		case s.limiter <- struct{}{}:
		default:
			continue
		}
		go func(b []byte) { defer func() { <-s.limiter }(); _, _ = pc.WriteTo(s.Resolve(b), peer) }(append([]byte(nil), buf[:n]...))
	}
}
func (s *Server) DoH(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 405)
		return
	}
	b, e := io.ReadAll(io.LimitReader(r.Body, 65535))
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	w.Header().Set("content-type", "application/dns-message")
	w.Write(s.Resolve(b))
}
func (s *Server) TCPConn(conn net.Conn) {
	defer conn.Close()
	for {
		var n uint16
		if binary.Read(conn, binary.BigEndian, &n) != nil {
			return
		}
		if n > 65535 {
			return
		}
		b := make([]byte, n)
		if _, e := io.ReadFull(conn, b); e != nil {
			return
		}
		p := s.Resolve(b)
		if len(p) > 65535 {
			return
		}
		binary.Write(conn, binary.BigEndian, uint16(len(p)))
		conn.Write(p)
	}
}

var _ = http.MethodGet

func wireType(code uint16) zone_domain.RecordType {
	switch code {
	case 1:
		return zone_domain.A
	case 2:
		return zone_domain.NS
	case 5:
		return zone_domain.CNAME
	case 15:
		return zone_domain.MX
	case 16:
		return zone_domain.TXT
	case 28:
		return zone_domain.AAAA
	case 33:
		return zone_domain.SRV
	case 64:
		return zone_domain.SVCB
	case 65:
		return zone_domain.HTTPS
	case 257:
		return zone_domain.CAA
	case 6:
		return zone_domain.SOA
	default:
		return zone_domain.RecordType("ANY")
	}
}
