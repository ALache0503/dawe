package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginadapter "github.com/ALache0503/dawe/backend/internal/adapters/http/gin"
	"github.com/ALache0503/dawe/backend/internal/adapters/postgres"
	"github.com/ALache0503/dawe/backend/internal/application"
)

func main() {
	databaseURL := requiredEnv("DATABASE_URL")
	address := envOrDefault("HTTP_ADDRESS", ":8080")
	allowedOrigin := envOrDefault("CORS_ALLOWED_ORIGIN", "http://localhost:5173")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.OpenPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer pool.Close()

	repository := postgres.NewProteinRepository(pool)
	service := application.NewProteinService(repository)
	handler := ginadapter.NewProteinHandler(service)
	router := ginadapter.NewRouter(handler, ginadapter.RouterConfig{AllowedOrigin: allowedOrigin})

	serverDone := make(chan error, 1)
	go func() {
		log.Printf("API listening on %s", address)
		serverDone <- router.Run(address)
	}()

	select {
	case err := <-serverDone:
		log.Fatalf("server failed: %v", err)
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = shutdownCtx // Gin's Router.Run starts net/http internally; graceful Server shutdown is an optional next refinement.
}

func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
