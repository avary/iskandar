package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRequestManager(t *testing.T) {
	t.Run("registers request", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager(10)
		requestId := "test-request-1"
		subdomain := "test-subdomain"

		ch, err := manager.RegisterRequest(requestId, subdomain)

		require.NoError(t, err)
		assert.NotNil(t, ch)
		assert.Len(t, manager.requestChannelMap, 1)
		assert.Equal(t, 1, manager.requestCounts[subdomain])

		retrievedCh, ok := manager.GetRequestChannel(requestId)
		assert.True(t, ok)
		assert.Equal(t, ch, retrievedCh)
	})

	t.Run("gets registered request channel", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager(10)
		requestId := "test-request-2"
		subdomain := "test-subdomain"

		ch, err := manager.RegisterRequest(requestId, subdomain)
		require.NoError(t, err)
		require.NotNil(t, ch)

		retrievedCh, ok := manager.GetRequestChannel(requestId)
		assert.True(t, ok)
		assert.NotNil(t, retrievedCh)
		assert.Equal(t, ch, retrievedCh)
	})

	t.Run("returns false for non-existent request", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager(10)

		ch, ok := manager.GetRequestChannel("nonexistent")
		assert.False(t, ok)
		assert.Nil(t, ch)
	})

	t.Run("removes registered request and closes channel", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager(10)
		requestId := "test-request-3"
		subdomain := "test-subdomain"

		ch, err := manager.RegisterRequest(requestId, subdomain)
		require.NoError(t, err)
		require.NotNil(t, ch)

		manager.RemoveRequest(requestId, subdomain)

		_, ok := <-ch
		assert.False(t, ok, "channel should be closed")

		retrievedCh, ok := manager.GetRequestChannel(requestId)
		assert.False(t, ok)
		assert.Nil(t, retrievedCh)
		assert.NotContains(t, manager.requestChannelMap, requestId)
		assert.Equal(t, 0, manager.requestCounts[subdomain])
	})

	t.Run("enforces maximum concurrent requests per tunnel", func(t *testing.T) {
		t.Parallel()
		maxPerTunnel := 2
		manager := NewInMemoryRequestManager(maxPerTunnel)
		subdomain := "test-subdomain"

		ch1, err := manager.RegisterRequest("req-1", subdomain)
		require.NoError(t, err)
		require.NotNil(t, ch1)

		ch2, err := manager.RegisterRequest("req-2", subdomain)
		require.NoError(t, err)
		require.NotNil(t, ch2)

		ch3, err := manager.RegisterRequest("req-3", subdomain)
		assert.ErrorIs(t, err, ErrMaxRequestsPerTunnel)
		assert.Nil(t, ch3)

		// After removing one, should be able to register again
		manager.RemoveRequest("req-1", subdomain)

		ch4, err := manager.RegisterRequest("req-4", subdomain)
		require.NoError(t, err)
		assert.NotNil(t, ch4)
	})

	t.Run("tracks requests per tunnel independently", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager(2)
		subdomain1 := "subdomain-1"
		subdomain2 := "subdomain-2"

		ch1, err := manager.RegisterRequest("req-1", subdomain1)
		require.NoError(t, err)
		require.NotNil(t, ch1)

		ch2, err := manager.RegisterRequest("req-2", subdomain1)
		require.NoError(t, err)
		require.NotNil(t, ch2)

		// subdomain1 at limit, but subdomain2 should still work
		ch3, err := manager.RegisterRequest("req-3", subdomain2)
		require.NoError(t, err)
		require.NotNil(t, ch3)

		assert.Equal(t, 2, manager.requestCounts[subdomain1])
		assert.Equal(t, 1, manager.requestCounts[subdomain2])
	})
}
