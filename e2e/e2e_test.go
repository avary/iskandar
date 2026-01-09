package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()

	composeFilePath := "docker-compose.e2e.yml"

	composeStack, err := compose.NewDockerCompose(composeFilePath)
	require.NoError(t, err, "Failed to create docker compose")

	t.Cleanup(func() {
		require.NoError(t, composeStack.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal))
	})

	err = composeStack.Up(ctx, compose.Wait(true))
	require.NoError(t, err, "Failed to start docker compose")

	// Start CLI process locally
	cliCmd, subdomain := startCLI(t, ctx)
	defer cliCmd.Process.Kill()

	t.Logf("Tunnel established with subdomain: %s", subdomain)

	t.Run("POST request proxied through tunnel", func(t *testing.T) {
		tunnelURL := fmt.Sprintf("http://%s/api/data", subdomain)

		resp, err := http.Post(tunnelURL, "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Contains(t, string(body), "Peter")
		assert.Contains(t, string(body), "Perlepes")
	})

	t.Run("GET request with query params", func(t *testing.T) {
		tunnelURL := fmt.Sprintf("http://%s/search?q=test&limit=5", subdomain)

		resp, err := http.Get(tunnelURL)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Contains(t, string(body), "test")
		assert.Contains(t, string(body), "limit")
	})

	t.Run("error response codes forwarded", func(t *testing.T) {
		// Test 404
		resp, err := http.Get(fmt.Sprintf("http://%s/not-found", subdomain))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		resp, err = http.Get(fmt.Sprintf("http://%s/server-error", subdomain))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func startCLI(t *testing.T, ctx context.Context) (*exec.Cmd, string) {
	t.Helper()

	// Build CLI binary if not exists
	cliBinary := "../iskndr/iskndr"
	if _, err := os.Stat(cliBinary); os.IsNotExist(err) {
		t.Log("Building CLI binary...")
		buildCmd := exec.Command("go", "build", "-o", "iskndr", "./cmd/iskndr")
		buildCmd.Dir = "../iskndr"
		output, err := buildCmd.CombinedOutput()
		if err != nil {
			t.Logf("Build output:\n%s", string(output))
			require.NoError(t, err, "Failed to build CLI")
		}
	}

	cmd := exec.CommandContext(ctx, cliBinary, "tunnel", "localhost:3003", "--server", "http://localhost:8080", "--logging")

	var outputBuf strings.Builder
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf

	err := cmd.Start()
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	cliOutput := outputBuf.String()

	t.Logf("CLI output:\n%s", cliOutput)

	subdomainRegex := regexp.MustCompile(`http://([a-z0-9]+\.localhost\.direct:8080)`)
	matches := subdomainRegex.FindStringSubmatch(cliOutput)
	require.NotEmpty(t, matches, "Failed to extract subdomain from CLI output")

	subdomain := matches[1]

	return cmd, subdomain
}
