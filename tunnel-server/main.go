package main

import (
	"log"
	"net/http"
	"os"

	"github.com/igneel64/iskandar/server/internal/logger"
)

func main() {
	logger.Initialize(true)

	connectionStore := NewInMemoryConnectionStore()
	requestManager := NewInMemoryRequestManager()

	// Read configuration from environment
	baseScheme := os.Getenv("ISKNDR_BASE_SCHEME")
	if baseScheme == "" {
		baseScheme = "http"
	}

	baseDomain := os.Getenv("ISKNDR_BASE_DOMAIN")
	if baseDomain == "" {
		baseDomain = "localhost.direct:8080"
	}

	port := os.Getenv("ISKNDR_PORT")
	if port == "" {
		port = "8080"
	}

	// Construct full public URL base
	publicURLBase := baseScheme + "://" + baseDomain

	s := NewIskndrServer(publicURLBase, connectionStore, requestManager)

	logger.ServerStarted(8080)
	log.Fatal(http.ListenAndServe(":"+port, s))
}
