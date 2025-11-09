package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/trust-me-im-an-engineer/mini-reddit/graph"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/config"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/comment"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/post"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/service/subscription"
	"github.com/trust-me-im-an-engineer/mini-reddit/internal/storage/inmemory"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Set slog to output json to stdout
	{
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.LogLevel,
		})
		slog.SetDefault(slog.New(handler))
		slog.Info("set json logging to stdout", "level", cfg.LogLevel)
	}

	storage := inmemory.New()

	resolver := graph.NewResolver(
		post.NewService(storage),
		comment.NewService(storage),
		subscription.NewService(),
	)

	// Start server
	{
		srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

		srv.AddTransport(transport.Options{})
		srv.AddTransport(transport.GET{})
		srv.AddTransport(transport.POST{})

		srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

		srv.Use(extension.Introspection{})
		srv.Use(extension.AutomaticPersistedQuery{
			Cache: lru.New[string](100),
		})

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)

		slog.Info("server running", "address", cfg.Address)
		slog.Info("playground running", "address", cfg.Address)
		log.Fatal(http.ListenAndServe(cfg.Address, nil))
	}
}
