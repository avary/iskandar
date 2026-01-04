package shared

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/igneel64/iskandar/shared/protocol"
)

func SerializeHeaders(headers http.Header) map[string]string {
	serializedHeaders := make(map[string]string)
	for k, v := range headers {
		serializedHeaders[k] = strings.Join(v, ", ")
	}

	return serializedHeaders
}

type SafeWebSocketConn struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewSafeWebSocketConn(conn *websocket.Conn) *SafeWebSocketConn {
	return &SafeWebSocketConn{
		conn: conn,
	}
}

func (s *SafeWebSocketConn) WriteJSON(msg *protocol.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteJSON(msg)
}

func (s *SafeWebSocketConn) ReadJSON(msg *protocol.Message) error {
	return s.conn.ReadJSON(msg)
}

func (s *SafeWebSocketConn) Close() error {
	return s.conn.Close()
}
