package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/igneel64/iskandar/shared"
	"github.com/igneel64/iskandar/shared/protocol"
)

type IskndrServer struct {
	http.Handler
	connStore      ConnectionStore
	requestManager RequestManager
}

func NewIskndrServer(connectionStore ConnectionStore, requestManager RequestManager) *IskndrServer {
	s := &IskndrServer{
		connStore:      connectionStore,
		requestManager: requestManager,
	}

	var upgrader = websocket.Upgrader{
		// CheckOrigin: func(r *http.Request) bool { return true },
	}

	router := http.NewServeMux()
	router.HandleFunc("/tunnel/connect", func(w http.ResponseWriter, r *http.Request) {
		con, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
			return
		}
		defer con.Close()

		subdomainKey, err := s.connStore.RegisterConnection(con)
		if err != nil {
			http.Error(w, "Failed to register connection", http.StatusInternalServerError)
			return
		}

		err = con.WriteJSON(&protocol.RegisterTunnelMessage{Subdomain: "http://" + subdomainKey + ".localhost.direct:8080"})
		if err != nil {
			http.Error(w, "Failed to send register tunnel message", http.StatusInternalServerError)
			return
		}

		for {
			var msg protocol.Message
			if err = con.ReadJSON(&msg); err != nil {
				fmt.Println("read error:", err)
				s.connStore.RemoveConnection(subdomainKey)
				return
			}
			fmt.Println("recv message:", msg)

			if ch, ok := s.requestManager.GetRequestChannel(msg.Id); ok {
				ch <- msg
			}
		}
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		subdomain := extractAssignedSubdomain(r.Host)
		conn, err := s.connStore.GetConnection(subdomain)
		fmt.Println(subdomain)
		if err != nil {
			http.Error(w, "No tunnel found for subdomain", http.StatusNotFound)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		requestId := uuid.New().String()

		ch := s.requestManager.RegisterRequest(requestId)
		defer s.requestManager.RemoveRequest(requestId)

		message := &protocol.Message{
			Type:    "request",
			Id:      requestId,
			Body:    bodyBytes,
			Method:  r.Method,
			Headers: shared.SerializeHeaders(r.Header),
			Path:    r.RequestURI,
		}

		if err = conn.WriteJSON(message); err != nil {
			http.Error(w, "Failed to forward request to tunnel", http.StatusInternalServerError)
			return
		}

		select {
		case response, ok := <-ch:
			if !ok {
				http.Error(w, "Failed to get response from tunnel", http.StatusInternalServerError)
				return
			}
			fmt.Println("resp: ", response)

			for k, v := range response.Headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(response.Status)
			w.Write(response.Body)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			for !response.Done {
				select {
				case response, ok = <-ch:
					if !ok {
						fmt.Println("channel closed mid-stream")
						return
					}

					w.Write(response.Body)
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}

				case <-time.After(30 * time.Second):
					fmt.Println("timeout waiting for chunk")
					return
				}
			}

		case <-time.After(30 * time.Second):
			http.Error(w, "Timeout waiting for response from tunnel", http.StatusGatewayTimeout)
			return
		}
	})

	s.Handler = router

	return s
}

func extractAssignedSubdomain(host string) string {
	parts := strings.Split(host, ".")
	return parts[0]
}
