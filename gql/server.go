package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/encero/reciper-api/gql/graph"
	"github.com/encero/reciper-api/gql/graph/generated"
	"github.com/encero/reciper-api/pkg/common"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const defaultPort = "8080"
const defaultNatsURL = "nats://localhost:4222"

func main() {
	if err := setupAndRun(); err != nil {
		fmt.Println("crashing:", err.Error())
		os.Exit(1)
	}
}

func setupAndRun() error {
	port := defaultPort
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	natsURL := defaultNatsURL
	if url, ok := os.LookupEnv("NATS_URL"); ok {
		natsURL = url
	}

	logger, err := common.LoggerFromEnv()
	if err != nil {
		return fmt.Errorf("setup logger: %w", err)
	}

	err = run(context.Background(), logger, port, natsURL)
	if err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context, lg *zap.Logger, port, natsURL string) error {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("connecting nats url: %q err: %w", natsURL, err)
	}

	ec, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats encoded conn: %w", err)
	}

	resolver := graph.NewResolver(ec)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)

		lg.Debug("operation", zap.String("operation", oc.OperationName), zap.String("query", oc.RawQuery))

		return next(ctx)
	})

	mux := &http.ServeMux{}

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	lg.Sugar().Infof("connect to http://localhost:%s/ for GraphQL playground", port)

	server := http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       time.Second * 30,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: time.Second * 30,
	}

	go func() {
		<-ctx.Done()

		lg.Info("GQL server shutting down")

		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}
