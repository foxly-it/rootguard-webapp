// =====================================================
// File: backend/cmd/rootguard-webapp/main.go
// Purpose: RootGuard WebApp Backend Entry Point
// Notes: Includes build metadata injection (version, commit)
// =====================================================

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/foxly-it/rootguard-webapp/backend/internal/httpapi"
)

// =====================================================
// Build metadata (werden via -ldflags injiziert)
// Docker build setzt:
// -X main.version=<version>
// -X main.commit=<commit>
// =====================================================

var version = "dev"
var commit = "unknown"

func main() {

	log.Printf("RootGuard WebApp starting (version=%s, commit=%s)", version, commit)

	port := getEnv("PORT", "8080")

	router := httpapi.NewRouter()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server stopped cleanly")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
