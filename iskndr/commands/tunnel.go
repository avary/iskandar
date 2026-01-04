package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/igneel64/iskandar/shared"
	"github.com/igneel64/iskandar/shared/protocol"
	"github.com/spf13/cobra"
)

func newTunnelCommand() *cobra.Command {
	tunnelCmd := &cobra.Command{
		Use: "tunnel <port>",

		Short:                 "Expose a local application to the internet",
		Long:                  "This command allows you to create a tunnel to your local application, making it accessible from the internet.",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("port must be a number: %w", err)
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("port must be between 1 and 65535")
			}
			destinationAddress := "http://localhost:" + strconv.Itoa(port)
			fmt.Printf("Starting iskndr on port %d...\n", port)
			fmt.Printf("Requests will be routed to %s\n", destinationAddress)

			c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/tunnel/connect", nil)
			if err != nil {
				return fmt.Errorf("failed to connect to websocket: %w", err)
			}
			defer c.Close()
			safeWriteConn := shared.NewSafeWebSocketConn(c)

			var regMsg protocol.RegisterTunnelMessage
			if err = c.ReadJSON(&regMsg); err != nil {
				return fmt.Errorf("failed to read register tunnel message: %w", err)
			}

			fmt.Printf("Tunnel url at %s\n", regMsg.Subdomain)

			for {
				var requestMsg protocol.Message
				if err := c.ReadJSON(&requestMsg); err != nil {
					return fmt.Errorf("failed to read request message: %w", err)
				}
				fmt.Printf("Received request: %+v\n", requestMsg)

				go sendResponse(safeWriteConn, &requestMsg, destinationAddress)
			}
		},
	}

	return tunnelCmd
}

func sendResponse(c *shared.SafeWebSocketConn, requestMsg *protocol.Message, destinationAddress string) {
	req, err := http.NewRequest(requestMsg.Method, destinationAddress+requestMsg.Path, bytes.NewReader(requestMsg.Body))

	if err != nil {
		c.WriteJSON(&protocol.Message{
			Type:   "response",
			Id:     requestMsg.Id,
			Status: http.StatusInternalServerError,
			Body:   []byte("Failed to create request"),
		})
		return
	}

	for k, v := range requestMsg.Headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {

		c.WriteJSON(&protocol.Message{
			Type:   "response",
			Id:     requestMsg.Id,
			Status: http.StatusBadGateway,
			Body:   []byte(fmt.Sprintf("Failed to reach local app: %v", err)),
		})
		return
	}

	/* Used for not re-sending extra data, mostly headers, which can be pretty big if response is not done. */
	firstChunk := true

	byteBuffer := make([]byte, 32*1024)
	for {
		byteCount, err := res.Body.Read(byteBuffer)

		if err != nil && err != io.EOF {
			if firstChunk {
				c.WriteJSON(&protocol.Message{
					Type:   "response",
					Id:     requestMsg.Id,
					Status: http.StatusBadGateway,
					Body:   []byte(fmt.Sprintf("Failed to read response body: %v", err)),
					Done:   true,
				})
			} else {
				// Already sent status - just log and abort
				fmt.Printf("error reading response body mid-stream: %v\n", err)
			}
			break
		}

		if firstChunk || byteCount > 0 {
			responseMsg := protocol.Message{
				Type: "response",
				Id:   requestMsg.Id,
				Body: byteBuffer[:byteCount],
				Done: err == io.EOF,
			}

			if firstChunk {
				responseMsg.Status = res.StatusCode
				responseMsg.Headers = shared.SerializeHeaders(res.Header)
				firstChunk = false
			}

			if err = c.WriteJSON(&responseMsg); err != nil {
				// log error
				fmt.Printf("failed to write response message: %v\n", err)
			}
		}

		if err == io.EOF {
			break
		}
	}

	res.Body.Close()
}
