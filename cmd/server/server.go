package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"subscriptionplus/server/config"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func (s *httpServer) httpStart() error {
	routes := s.routes()
	env := os.Getenv("GO_ENV")

	if env == "PROD" {
		fmt.Println("[INFO] Starting HTTPS server with autocert")

		manager := &autocert.Manager{
			Cache:      autocert.DirCache(filepath.Join(os.TempDir(), "cert")),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.ADomain, config.AWWWDomain),
		}

		httpServe := &http.Server{
			Addr:         fmt.Sprintf(":%d", config.HttpsServerPort),
			Handler:      routes,
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
			IdleTimeout:  90 * time.Second,
			TLSConfig:    manager.TLSConfig(),
		}

		go func() {
			fmt.Println("[INFO] Redirecting :80 → :443")
			http.ListenAndServe(fmt.Sprintf(":%d", config.HttpServerPort), manager.HTTPHandler(nil))
		}()

		return httpServe.ListenAndServeTLS("", "")
	}

	fmt.Printf("[INFO] Starting HTTP server on :%d\n", config.LocalHttpServerPort)

	httpServe := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.LocalHttpServerPort),
		Handler:      routes,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  90 * time.Second,
	}

	return httpServe.ListenAndServe()
}
