package sniproxy

import (
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"
)

type readOnlyConn struct{ io.Reader }

func (readOnlyConn) Write([]byte) (int, error)        { return 0, net.ErrClosed }
func (readOnlyConn) Close() error                     { return nil }
func (readOnlyConn) LocalAddr() net.Addr              { return nil }
func (readOnlyConn) RemoteAddr() net.Addr             { return nil }
func (readOnlyConn) SetDeadline(time.Time) error      { return nil }
func (readOnlyConn) SetReadDeadline(time.Time) error  { return nil }
func (readOnlyConn) SetWriteDeadline(time.Time) error { return nil }

func copyConn(dst, src net.Conn, wg *sync.WaitGroup) {
	_, _ = io.Copy(dst, src)
	_ = dst.Close()
	wg.Done()
}

func readClientHello(r io.Reader) (fqdn string) {
	config := new(tls.Config)
	config.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
		fqdn = info.ServerName
		return nil, nil
	}
	_ = tls.Server(readOnlyConn{r}, config).Handshake()
	return
}
