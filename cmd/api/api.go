package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/encero/reciper-api/api"
	"github.com/encero/reciper-api/ent"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	natsURL := "localhost:4222"

	sqldb, err := sql.Open("sqlite", "file:db.lite?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		return fmt.Errorf("cant open sql database: %w", err)
	}
	defer sqldb.Close()

	entc := ent.NewClient(ent.Driver(entsql.OpenDB("sqlite3", sqldb)))
	defer entc.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("setup logger %w", err)
	}

	return api.Run(context.Background(), entc, logger, natsURL)
}
