package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *httpServer) httpStart() error {
	port, err := strconv.Atoi(DefaultPortStr)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	routes := s.routes()

	fmt.Printf("[INFO] Starting HTTP server on :%d\n", port)

	httpServe := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  90 * time.Second,
	}

	return httpServe.ListenAndServe()
}
