package graph

import "github.com/nats-io/nats.go"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ec *nats.EncodedConn
}

func NewResolver(ec *nats.EncodedConn) *Resolver {
	return &Resolver{
		ec: ec,
	}
}
