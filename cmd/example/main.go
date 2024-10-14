package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/delaneyj/geck/cmd/example/ecs"
	"github.com/go-chi/chi/v5"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Handle SIGINT and SIGTERM gracefully
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		log.Println("Shutting down...")
		cancel()
	}()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	w := ecs.NewWorld()

	e := w.NextEntity(
		ecs.WithName("Test"),
	)
	log.Printf("Entity: %v", e)

	r := chi.NewRouter()
	if err := ecs.SetupRoutes(ctx, w, r); err != nil {
		return fmt.Errorf("failed to setup routes: %w", err)
	}

	// Print out the routes
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Route: %s %s", method, route)
		return nil
	})

	port := 8080
	log.Printf("Hosting at http://localhost:%d", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown server: %v", err)
		}

	}()

	srv.ListenAndServe()
	<-ctx.Done()

	return nil
}
