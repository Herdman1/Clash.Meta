package tunnel

import (
	"fmt"
	"net"

	"github.com/Dreamacro/clash/adapter/inbound"
	N "github.com/Dreamacro/clash/common/net"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/transport/socks5"
)

type Listener struct {
	listener net.Listener
	addr     string
	target   socks5.Addr
	proxy    string
	closed   bool
}

// RawAddress implements C.Listener
func (l *Listener) RawAddress() string {
	return l.addr
}

// Address implements C.Listener
func (l *Listener) Address() string {
	return l.listener.Addr().String()
}

// Close implements C.Listener
func (l *Listener) Close() error {
	l.closed = true
	return l.listener.Close()
}

func (l *Listener) handleTCP(conn net.Conn, tunnel C.Tunnel, additions ...inbound.Addition) {
	N.TCPKeepAlive(conn)
	ctx := inbound.NewSocket(l.target, conn, C.TUNNEL, additions...)
	ctx.Metadata().SpecialProxy = l.proxy
	tunnel.HandleTCPConn(ctx)
}

func New(addr, target, proxy string, tunnel C.Tunnel, additions ...inbound.Addition) (*Listener, error) {
	l, err := inbound.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	targetAddr := socks5.ParseAddr(target)
	if targetAddr == nil {
		return nil, fmt.Errorf("invalid target address %s", target)
	}

	rl := &Listener{
		listener: l,
		target:   targetAddr,
		proxy:    proxy,
		addr:     addr,
	}

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				if rl.closed {
					break
				}
				continue
			}
			go rl.handleTCP(c, tunnel, additions...)
		}
	}()

	return rl, nil
}
