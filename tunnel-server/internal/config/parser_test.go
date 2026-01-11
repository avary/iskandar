package config

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractSubdomainURL(t *testing.T) {
	t.Run("https URL with port", func(t *testing.T) {
		baseURL, err := url.Parse("https://example.com:8080")
		require.NoError(t, err)

		result := ExtractSubdomainURL(baseURL, "test")
		assert.Equal(t, "https://test.example.com:8080", result)
	})

	t.Run("http URL without port", func(t *testing.T) {
		baseURL, err := url.Parse("http://localhost.direct")
		require.NoError(t, err)

		result := ExtractSubdomainURL(baseURL, "abc123")
		assert.Equal(t, "http://abc123.localhost.direct", result)
	})
}

func TestExtractAssignedSubdomain(t *testing.T) {
	t.Run("extracts first part of hostname with port", func(t *testing.T) {
		result, err := ExtractAssignedSubdomain("abc123.localhost.direct:8080")
		require.NoError(t, err)
		assert.Equal(t, "abc123", result)
	})

	t.Run("extracts first part of hostname without port", func(t *testing.T) {
		result, err := ExtractAssignedSubdomain("abc123.example.com")
		require.NoError(t, err)
		assert.Equal(t, "abc123", result)
	})

	t.Run("returns error for single part hostname", func(t *testing.T) {
		_, err := ExtractAssignedSubdomain("localhost")
		assert.Error(t, err)
	})

	t.Run("returns error for hostname with port but no subdomain", func(t *testing.T) {
		_, err := ExtractAssignedSubdomain("localhost:8080")
		assert.Error(t, err)
	})
}
