package tlsproxydialer

import (
	"crypto/tls"
	"net"
)

type TLSProxyDialer struct {
	Config *tls.Config
	Conn   *tls.Conn
}

func NewTLSProxyDialer(config *tls.Config) *TLSProxyDialer {
	self := &TLSProxyDialer{
		Config: config,
	}
	return self
}

func (self *TLSProxyDialer) Dial(network, addr string) (c net.Conn, err error) {
	t, err := tls.Dial(network, addr, self.Config)
	self.Conn = t
	c = t
	return
}
