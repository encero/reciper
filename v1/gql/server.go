package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/encero/reciper/gql/configuration"
	"github.com/encero/reciper/gql/graph"
	"github.com/encero/reciper/gql/graph/generated"
	"github.com/encero/reciper/pkg/common"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var VersionRef = "development"
var VersionCommit = "dev"

func main() {
	if err := setupAndRun(); err != nil {
		fmt.Println("crashing:", err.Error())
		os.Exit(1)
	}
}

func setupAndRun() error {
	logger, err := common.LoggerFromEnv()
	if err != nil {
		return fmt.Errorf("setup logger: %w", err)
	}

	cfg, err := configuration.FromEnvironment()
	if err != nil {
		return fmt.Errorf("configuring server: %w", err)
	}

	cfg.VersionRef = VersionRef
	cfg.VersionCommit = VersionCommit

	err = run(context.Background(), logger, cfg)
	if err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context, lg *zap.Logger, cfg configuration.Config) error {
	conn, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return fmt.Errorf("connecting nats url: %q err: %w", cfg.NatsURL, err)
	}

	ec, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats encoded conn: %w", err)
	}

	resolver := graph.NewResolver(ec, lg, cfg)
	config := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			Validation: validations(),
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)

		lg.Debug("operation", zap.String("operation", oc.OperationName), zap.String("query", oc.RawQuery))

		return next(ctx)
	})

	mux := &http.ServeMux{}

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	lg.Sugar().Infof("connect to http://localhost:%s/ for GraphQL playground", cfg.ServerPort)

	server := http.Server{
		Addr:              ":" + cfg.ServerPort,
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

func setupValidations() (*validator.Validate, ut.Translator) {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)

	translator, ok := uni.GetTranslator("en")
	if !ok {
		panic(fmt.Errorf("setupValidations no translator for 'en'"))
	}

	err := en_translations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		panic(fmt.Errorf("RegisterDefaultTranslations: %w", err))
	}

	return validate, translator
}

func validations() func(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (res interface{}, err error) {
	validate, translator := setupValidations()

	return func(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (res interface{}, err error) {
		val, err := next(ctx)
		if err != nil {
			return nil, err
		}

		fieldName := *graphql.GetPathContext(ctx).Field

		err = validate.Var(val, constraint)
		if err != nil {
			validationErrors := validator.ValidationErrors{}
			if errors.As(err, &validationErrors) {
				return val, fmt.Errorf("%s%+v", fieldName, validationErrors[0].Translate(translator))
			}

			return val, err
		}

		return val, nil
	}
}
