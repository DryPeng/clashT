package tunnel

import (
	"fmt"
	"net"

	"github.com/DryPeng/clashT/adapter/inbound"
	"github.com/DryPeng/clashT/common/pool"
	C "github.com/DryPeng/clashT/constant"
	"github.com/DryPeng/clashT/transport/socks5"
)

type PacketConn struct {
	conn   net.PacketConn
	addr   string
	target socks5.Addr
	proxy  string
	closed bool
}

// RawAddress implements C.Listener
func (l *PacketConn) RawAddress() string {
	return l.addr
}

// Address implements C.Listener
func (l *PacketConn) Address() string {
	return l.conn.LocalAddr().String()
}

// Close implements C.Listener
func (l *PacketConn) Close() error {
	l.closed = true
	return l.conn.Close()
}

func NewUDP(addr, target, proxy string, in chan<- *inbound.PacketAdapter) (*PacketConn, error) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}

	targetAddr := socks5.ParseAddr(target)
	if targetAddr == nil {
		return nil, fmt.Errorf("invalid target address %s", target)
	}

	sl := &PacketConn{
		conn:   l,
		target: targetAddr,
		proxy:  proxy,
		addr:   addr,
	}
	go func() {
		for {
			buf := pool.Get(pool.UDPBufferSize)
			n, remoteAddr, err := l.ReadFrom(buf)
			if err != nil {
				pool.Put(buf)
				if sl.closed {
					break
				}
				continue
			}
			sl.handleUDP(l, in, buf[:n], remoteAddr)
		}
	}()

	return sl, nil
}

func (l *PacketConn) handleUDP(pc net.PacketConn, in chan<- *inbound.PacketAdapter, buf []byte, addr net.Addr) {
	packet := &packet{
		pc:      pc,
		rAddr:   addr,
		payload: buf,
	}

	ctx := inbound.NewPacket(l.target, pc.LocalAddr(), packet, C.TUNNEL)
	ctx.Metadata().SpecialProxy = l.proxy
	select {
	case in <- ctx:
	default:
	}
}
