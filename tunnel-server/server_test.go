package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/igneel64/iskandar/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockConnectionStore struct {
	mock.Mock
}

func (m *MockConnectionStore) RegisterConnection(conn *websocket.Conn) (string, error) {
	args := m.Called(conn)
	return args.String(0), args.Error(1)
}

func (m *MockConnectionStore) GetConnection(subdomain string) (*shared.SafeWebSocketConn, error) {
	args := m.Called(subdomain)
	return args.Get(0).(*shared.SafeWebSocketConn), args.Error(1)
}

func (m *MockConnectionStore) RemoveConnection(subdomain string) {
	m.Called(subdomain)
}

type MockRequestManager struct {
	mock.Mock
}

func (m *MockRequestManager) GetRequestChannel(requestId string) (MessageChannel, bool) {
	args := m.Called(requestId)
	return args.Get(0).(MessageChannel), args.Bool(1)
}

func (m *MockRequestManager) RegisterRequest(requestId string) MessageChannel {
	args := m.Called(requestId)
	return args.Get(0).(MessageChannel)
}

func (m *MockRequestManager) RemoveRequest(requestId string) {
	m.Called(requestId)
}

func TestServer(t *testing.T) {
	server := NewIskndrServer("localhost.direct:8080", new(MockConnectionStore), new(MockRequestManager))

	t.Run("accepts websocket connection at /tunnel/connect", func(t *testing.T) {
		// Create test server
		ts := httptest.NewServer(server)
		defer ts.Close()

		// Convert http:// to ws://
		wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/tunnel/connect"

		// Connect via WebSocket
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "should connect to websocket")

		//nolint:errcheck
		defer conn.Close()

		// Verify connection is established
		assert.NotNil(t, conn)
	})
}
