package graph

import (
	"github.com/encero/reciper/gql/configuration"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ec  *nats.EncodedConn
	lg  *zap.Logger
	cfg configuration.Config
}

func NewResolver(ec *nats.EncodedConn, logger *zap.Logger, cfg configuration.Config) *Resolver {
	return &Resolver{
		ec:  ec,
		lg:  logger,
		cfg: cfg,
	}
}
