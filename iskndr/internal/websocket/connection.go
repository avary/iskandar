package websocket

import (
	"crypto/tls"

	"github.com/gorilla/websocket"
	"github.com/igneel64/iskandar/shared"
)

type Dialer interface {
	Dial() (*shared.SafeWebSocketConn, error)
}

func NewWriteSafeWSDialer(serverWSURL string, allowInsecure bool) *WriteSafeWSDialer {
	return &WriteSafeWSDialer{
		serverWSURL:   serverWSURL,
		allowInsecure: allowInsecure,
	}
}

type WriteSafeWSDialer struct {
	serverWSURL   string
	allowInsecure bool
}

func (d *WriteSafeWSDialer) Dial() (*shared.SafeWebSocketConn, error) {
	dialer := websocket.DefaultDialer
	if d.allowInsecure {
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c, _, err := dialer.Dial(d.serverWSURL, nil)
	if err != nil {
		return nil, err
	}
	safeWriteConn := shared.NewSafeWebSocketConn(c)

	return safeWriteConn, nil
}
