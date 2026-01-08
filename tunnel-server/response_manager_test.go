package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRequestManager(t *testing.T) {
	t.Run("registers request", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager()
		requestId := "test-request-1"

		ch := manager.RegisterRequest(requestId)

		assert.NotNil(t, ch)
		assert.Len(t, manager.requestChannelMap, 1)

		retrievedCh, ok := manager.GetRequestChannel(requestId)
		assert.True(t, ok)
		assert.Equal(t, ch, retrievedCh)
	})

	t.Run("gets registered request channel", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager()
		requestId := "test-request-2"

		ch := manager.RegisterRequest(requestId)
		require.NotNil(t, ch)

		retrievedCh, ok := manager.GetRequestChannel(requestId)
		assert.True(t, ok)
		assert.NotNil(t, retrievedCh)
		assert.Equal(t, ch, retrievedCh)
	})

	t.Run("returns false for non-existent request", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager()

		ch, ok := manager.GetRequestChannel("nonexistent")
		assert.False(t, ok)
		assert.Nil(t, ch)
	})

	t.Run("removes registered request and closes channel", func(t *testing.T) {
		t.Parallel()
		manager := NewInMemoryRequestManager()
		requestId := "test-request-3"

		ch := manager.RegisterRequest(requestId)
		require.NotNil(t, ch)

		manager.RemoveRequest(requestId)

		_, ok := <-ch
		assert.False(t, ok, "channel should be closed")

		retrievedCh, ok := manager.GetRequestChannel(requestId)
		assert.False(t, ok)
		assert.Nil(t, retrievedCh)
		assert.NotContains(t, manager.requestChannelMap, requestId)
	})
}
