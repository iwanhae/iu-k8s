package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"iu-k8s.linecorp.com/server/internal/api"
	"iu-k8s.linecorp.com/server/internal/config"
	"iu-k8s.linecorp.com/server/internal/handlers"
	"iu-k8s.linecorp.com/server/internal/middleware"
	"iu-k8s.linecorp.com/server/internal/repository"
	"iu-k8s.linecorp.com/server/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository (in-memory for now)
	userRepo := repository.NewInMemoryUserRepository()

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handlers.NewUserHandler(userService)

	// Create router
	r := chi.NewRouter()

	// Configure middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RealIP())
	r.Use(middleware.Timeout(60 * time.Second))

	// Configure CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Configure JSON render
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {
			render.DefaultResponder(w, r, render.M{"error": err.Error()})
			return
		}
		render.DefaultResponder(w, r, v)
	}

	// Mount the generated API routes
	api.HandlerFromMux(userHandler, r)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
