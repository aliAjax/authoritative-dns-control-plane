package protocol

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type TLSConfig struct {
	CertFile   string
	KeyFile    string
	MinVersion uint16
	ClientCA   string
	Timeout    time.Duration
}

func LoadTLS(c TLSConfig) (*tls.Config, error) {
	if c.MinVersion == 0 {
		c.MinVersion = tls.VersionTLS13
	}
	if c.CertFile == "" || c.KeyFile == "" {
		return nil, fmt.Errorf("certificate and key required")
	}
	cert, e := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if e != nil {
		return nil, e
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: c.MinVersion}, nil
}

type Listener struct {
	net.Listener
	Config *tls.Config
}

func Listen(addr string, c *tls.Config) (net.Listener, error) {
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return nil, e
	}
	return tls.NewListener(l, c), nil
}
func AcceptLoop(l net.Listener, handler func(net.Conn), stop <-chan struct{}) error {
	for {
		c, e := l.Accept()
		if e != nil {
			select {
			case <-stop:
				return nil
			default:
				return e
			}
		}
		go handler(c)
	}
}
