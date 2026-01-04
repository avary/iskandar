package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createWSServerConnection(t *testing.T) *websocket.Conn {
	var upgrader = websocket.Upgrader{
		// CheckOrigin: func(r *http.Request) bool { return true },
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		con, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer con.Close()
	}))

	t.Cleanup(s.Close)

	wsUrl := "ws" + s.URL[4:]
	c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		c.Close()
	})

	return c
}

func TestInMemoryConnectionStore(t *testing.T) {
	connectionStore := NewInMemoryConnectionStore()
	t.Run("registers connection", func(t *testing.T) {
		conn := createWSServerConnection(t)
		subdomain, err := connectionStore.RegisterConnection(conn)
		assert.NoError(t, err)
		assert.Len(t, connectionStore.connMap, 1)
		assert.Equal(t, conn, connectionStore.connMap[subdomain])
	})
}
