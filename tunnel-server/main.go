package main

import (
	"log"
	"net/http"
)

func main() {
	connectionStore := NewInMemoryConnectionStore()
	requestManager := NewInMemoryRequestManager()
	s := NewIskndrServer(connectionStore, requestManager)
	log.Fatal(http.ListenAndServe(":8080", s))
}
