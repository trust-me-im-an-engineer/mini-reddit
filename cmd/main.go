package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/trust-me-im-an-engineer/mini-reddit/graph"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/config"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/comment"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/post"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/subscription"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage/inmemory"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage/postgres"
)

func main() {
	// --- Configuration ---
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// --- Logging ---
	{
		jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.LogLevel,
		})
		slog.SetDefault(slog.New(jsonHandler))
		slog.Info("set json logging to stdout", "level", cfg.LogLevel)
	}

	// --- Storage Initialization ---
	var storage storage.Storage
	if cfg.StorageType == "POSTGRES" {
		storage, err = postgres.New(context.Background(), *cfg.DB)
		if err != nil {
			slog.Error("failed to initialize postgres storage", "error", err)
			os.Exit(1)
		}

		defer func() {
			slog.Info("closing storage pool...")
			storage.Close()
			slog.Info("storage closed")
		}()
	} else {
		storage = inmemory.New()
	}

	// --- Services and GraphQL Resolver Setup ---
	resolver := graph.NewResolver(
		post.NewService(storage),
		comment.NewService(storage),
		subscription.NewService(),
	)

	// --- HTTP Server Setup ---
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](cfg.Graphql.QueryCache))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](cfg.Graphql.AutomaticPersistedQuery),
	})

	router := http.NewServeMux()
	router.Handle("/query", srv)

	if cfg.Graphql.Playground {
		router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
		slog.Info("playground running", "address", cfg.Address)
	}

	httpServer := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// --- Signal Handling Channel ---
	stopCh := make(chan os.Signal, 1)
	// Notify the stopCh for interrupt (Ctrl+C) and termination signals
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// --- Start Server in a Goroutine ---
	// Start the server in a goroutine so the main function can listen on stopCh
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server running", "address", cfg.Address)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// --- Block until signal or error ---
	select {
	case err := <-serverErrors:
		slog.Error("fatal error while serving http", "error", err)
		os.Exit(1)
	case <-stopCh:
		slog.Info("received stop signal, initiating graceful shutdown...")
	}

	// --- Graceful Shutdown Context ---
	// Create a context with a timeout for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// --- Shut Down HTTP Server ---
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("http server shutdown failed", "error", err)
	} else {
		slog.Info("http server gracefully stopped")
	}

	slog.Info("application stopped gracefully")
}
