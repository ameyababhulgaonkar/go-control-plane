package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ameya/go-control-plane/internal/config"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func main() {

	config := config.Load()

	var driver string
	if strings.HasPrefix(config.DBUrl, "postgres://") {
		driver = "postgres"
	} else {
		driver = "sqlite"
	}

	fmt.Printf("Driver: %s\n", driver)
	db, err := sql.Open("sqlite", config.DBUrl)

	if err != nil {
		log.Fatalf("Failed to open DB, %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// store.New(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("HTTP server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("shutdown signal received")
	case <-ctx.Done():
		log.Println("context cancelled")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
