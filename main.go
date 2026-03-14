package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/peter-stratton/dark-factory-demo-go/internal/server"
	"github.com/peter-stratton/dark-factory-demo-go/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := store.New()
	srv := server.New(s)

	addr := fmt.Sprintf(":%s", port)
	log.Println("starting server on", addr)

	if err := http.ListenAndServe(addr, srv.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
