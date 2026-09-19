package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grip/authz"
	"grip/config"
	"grip/jobs"
	"grip/routes"
)

// @title Grip API
// @version 1.0
// @description Grip API with JWT auth and email verification.
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT. Example: "Bearer eyJhbGciOiJIUzI1NiJ9..."
func main() {
	config.InitLogger()
	config.LoadEnv()
	config.RunMigrations()
	config.ConnectDatabase()
	config.ConnectRedis()

	if err := authz.Init(config.DB); err != nil {
		slog.Error("Casbin setup failed", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := jobs.Setup(ctx, os.Getenv("DATABASE_URL")); err != nil {
		slog.Error("River setup failed", "error", err)
		os.Exit(1)
	}

	router := routes.Setup()

	// Timeouts bound slow-client (Slowloris) abuse. Port from APP_PORT.
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("Server listening", "addr", ":"+port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}
	if err := jobs.Shutdown(shutdownCtx); err != nil {
		slog.Error("River shutdown error", "error", err)
	}
	if err := config.CloseRedis(); err != nil {
		slog.Error("Redis shutdown error", "error", err)
	}
}
