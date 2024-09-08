package dns

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"

	"github.com/DryPeng/clashT/component/dialer"
	"github.com/DryPeng/clashT/component/resolver"

	D "github.com/miekg/dns"
	"github.com/lucas-clemente/quic-go"
)

type doqClient struct {
	url    string
	config *quic.Config
}

func (dc *doqClient) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return dc.ExchangeContext(context.Background(), m)
}

func (dc *doqClient) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	// Set ID to 0 for better cache friendliness
	newM := *m
	newM.Id = 0

	data, err := newM.Pack()
	if err != nil {
		return nil, err
	}

	conn, err := quic.DialAddrContext(ctx, dc.url, &tls.Config{NextProtos: []string{"doq"}}, dc.config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial QUIC: %w", err)
	}
	defer conn.CloseWithError(0, "")

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()

	_, err = stream.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to stream: %w", err)
	}

	resp := make([]byte, 4096)
	n, err := stream.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stream: %w", err)
	}

	msg = new(D.Msg)
	if err := msg.Unpack(resp[:n]); err != nil {
		return nil, fmt.Errorf("failed to unpack DNS response: %w", err)
	}

	// Restore original ID
	msg.Id = m.Id
	return msg, nil
}

func newDoQClient(url, iface string, r *Resolver) *doqClient {
	return &doqClient{
		url: url,
		config: &quic.Config{
			DialAddr: func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				ips, err := resolver.LookupIPWithResolver(ctx, host, r)
				if err != nil {
					return nil, err
				} else if len(ips) == 0 {
					return nil, fmt.Errorf("%w: %s", resolver.ErrIPNotFound, host)
				}
				ip := ips[rand.Intn(len(ips))]

				options := []dialer.Option{}
				if iface != "" {
					options = append(options, dialer.WithInterface(iface))
				}

				udpConn, err := dialer.ListenPacket(ctx, "udp", "", options...)
				if err != nil {
					return nil, err
				}

				return quic.DialEarlyContext(ctx, udpConn, &net.UDPAddr{IP: ip, Port: 853}, host, tlsCfg, cfg)
			},
		},
	}
}
