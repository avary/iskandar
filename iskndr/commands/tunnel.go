package commands

import (
	"fmt"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/igneel64/iskandar/iskndr/internal/client"
	"github.com/igneel64/iskandar/iskndr/internal/config"
	"github.com/igneel64/iskandar/iskndr/internal/logger"
	"github.com/igneel64/iskandar/iskndr/internal/ui"
	iskWS "github.com/igneel64/iskandar/iskndr/internal/websocket"
	"github.com/igneel64/iskandar/shared"
	"github.com/spf13/cobra"
)

func newTunnelCommand() *cobra.Command {
	var enableLogging bool
	var serverUrl string
	var allowInsecure bool

	tunnelCmd := &cobra.Command{
		Use:   "tunnel <destination>",
		Short: "Expose a local application to the internet",
		Long: `This command creates a tunnel to your local application, making it accessible from the internet.

The destination can be specified as:
  - port number only (e.g., '8080') - defaults to localhost:8080
  - host:port (e.g., 'foo.bar:80') - connects to the specified host and port`,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Initialize(enableLogging)

			destinationAddress, err := config.ParseDestination(args[0])
			if err != nil {
				return err
			}

			serverWSUrl, err := config.ParseServerURL(serverUrl)
			if err != nil {
				return err
			}

			logger.TunnelStarting(destinationAddress, serverWSUrl)

			dialer := iskWS.NewWriteSafeWSDialer(serverWSUrl, allowInsecure)
			c, err := dialer.Dial()
			if err != nil {
				logger.TunnelDisconnected(err)
				return fmt.Errorf("failed to connect to websocket: %w", err)
			}

			//nolint:errcheck
			defer c.Close()

			client := client.NewIskndrClient(c)

			regMsg, err := client.Register()
			if err != nil {
				logger.TunnelDisconnected(err)
				return fmt.Errorf("failed to read register tunnel message: %w", err)
			}
			logger.TunnelConnected(regMsg.Subdomain)

			var program *tea.Program
			if !enableLogging {
				program = ui.InitUi(destinationAddress, serverUrl, regMsg.Subdomain)
			}

			setupShutdownHandler(program, c)
			return client.AcceptRequests(destinationAddress)
		},
	}

	tunnelCmd.Flags().StringVar(&serverUrl, "server", "", "Tunnel server URL (e.g., localhost:8080, https://tunnel.example.com).")
	tunnelCmd.Flags().BoolVar(&enableLogging, "logging", false, "Enable structured logging to stdout")
	tunnelCmd.Flags().BoolVar(&allowInsecure, "allow-insecure", false, "Skip TLS certificate verification")
	if err := tunnelCmd.MarkFlagRequired("server"); err != nil {
		panic(err)
	}

	return tunnelCmd
}

func setupShutdownHandler(program *tea.Program, c *shared.SafeWebSocketConn) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		fmt.Println("\nShutting down tunnel...")
		if program != nil {
			program.Quit()
		}
		err := c.Close()
		if err != nil {
			fmt.Println("Error closing connection while shutting down:", err)
		}
	}()
}
