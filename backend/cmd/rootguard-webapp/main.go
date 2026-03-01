// =====================================================
// File: backend/cmd/rootguard-webapp/main.go
// Project: RootGuard WebApp
// Purpose: Application entrypoint
// Notes:
// - Injects build metadata via -ldflags
// - Registers graceful shutdown
// - Syncs version info with httpapi package
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
// Build Metadata (Injected at build time)
// These values are overwritten via:
//
// go build -ldflags "-X main.version=v0.1.0 -X main.commit=abc123"
//
// If not set, defaults remain.
// =====================================================

var version = "dev"
var commit = "unknown"

// =====================================================
// init()
// Sync build metadata into httpapi package
// This allows /api/version to expose real values
// =====================================================

func init() {
	httpapi.Version = version
	httpapi.Commit = commit
}

// =====================================================
// main()
// =====================================================

func main() {

	port := getEnv("PORT", "8080")

	router := httpapi.NewRouter()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Printf("RootGuard WebApp starting (version=%s, commit=%s)", version, commit)
		log.Printf("Listening on :%s", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown handling
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

// =====================================================
// getEnv()
// Helper to read environment variables with fallback
// =====================================================

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
