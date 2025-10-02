package main

import (
	"fmt"
	"net/http"
	"subscriptionplus/server/config"
	"time"
)

func (s *httpServer) httpStart() error {
	routes := s.routes()

	fmt.Println("\n" + `  _    _ _______ _______ _____
 | |  | |__   __|__   __|  __ \
 | |__| |  | |     | |  | |__) |
 |  __  |  | |     | |  |  ___/
 | |  | |  | |     | |  | |
 |_|  |_|  |_|     |_|  |_|

                                `)

	fmt.Printf("[INFO] Listening on :%d\n", config.HttpServerPort)

	httpServe := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HttpServerPort),
		Handler:      routes,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  90 * time.Second,
	}

	return httpServe.ListenAndServe()
}
