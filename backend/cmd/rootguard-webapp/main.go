// =====================================================
// File: backend/cmd/rootguard-webapp/main.go
// Project: RootGuard WebApp
// Purpose: Application entrypoint
// Notes:
// - Injects build metadata via -ldflags
// - Registers graceful shutdown
// - Syncs version info with httpapi package
// - Starts Docker Stats Collector
// =====================================================

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/foxly-it/rootguard-webapp/backend/internal/docker"
	"github.com/foxly-it/rootguard-webapp/backend/internal/httpapi"
)

// =====================================================
// Build Metadata (Injected at build time)
//
// Example:
//
// go build -ldflags "-X main.version=v0.1.0 -X main.commit=abc123"
// =====================================================

var version = "dev"
var commit = "unknown"

// =====================================================
// init()
// Sync build metadata into httpapi package
// =====================================================

func init() {
	httpapi.Version = version
	httpapi.Commit = commit
}

// =====================================================
// main()
// =====================================================

func main() {

	// -------------------------------------------------
	// Start Docker Stats Collector
	//
	// This goroutine continuously polls Docker stats
	// and updates the in-memory metrics cache.
	//
	// This ensures the dashboard does not need to call
	// the Docker API on every request.
	// -------------------------------------------------

	docker.StartStatsCollector()

	port := getEnv("PORT", "8080")

	router := httpapi.NewRouter()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// -------------------------------------------------
	// Start HTTP server
	// -------------------------------------------------

	go func() {

		log.Printf("RootGuard WebApp starting (version=%s, commit=%s)", version, commit)
		log.Printf("Listening on :%s", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}

	}()

	// -------------------------------------------------
	// Graceful Shutdown Handling
	// -------------------------------------------------

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
