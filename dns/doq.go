package dns

import (
	"context"
	"fmt"
	"net"

	"github.com/miekg/dns"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

type doqClient struct {
	url       string
	transport *http3.RoundTripper
}

func newDoQClient(url string) *doqClient {
	return &doqClient{
		url: url,
		transport: &http3.RoundTripper{
			QuicConfig: &quic.Config{},
		},
	}
}

func (dc *doqClient) Exchange(m *dns.Msg) (msg *dns.Msg, err error) {
	return dc.ExchangeContext(context.Background(), m)
}

func (dc *doqClient) ExchangeContext(ctx context.Context, m *dns.Msg) (*dns.Msg, error) {
	// Pack DNS message to binary format
	data, err := m.Pack()
	if err != nil {
		return nil, err
	}

	// Dial QUIC connection
	conn, err := quic.DialAddrContext(ctx, dc.url, dc.transport.TLSClientConfig, dc.transport.QuicConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial QUIC: %w", err)
	}
	defer conn.CloseWithError(0, "")

	// Open a stream for the DNS query
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()

	// Write DNS query to the stream
	_, err = stream.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to stream: %w", err)
	}

	// Read the response from the stream
	resp := make([]byte, 4096)
	n, err := stream.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stream: %w", err)
	}

	// Unpack the response
	msg := new(dns.Msg)
	if err := msg.Unpack(resp[:n]); err != nil {
		return nil, fmt.Errorf("failed to unpack DNS response: %w", err)
	}

	return msg, nil
}
