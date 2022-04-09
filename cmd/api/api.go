package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/encero/reciper/api"
	"github.com/encero/reciper/ent"
	"github.com/encero/reciper/pkg/common"
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
	if url, ok := os.LookupEnv("NATS_URL"); ok {
		natsURL = url
	}

	sqldb, err := sql.Open("sqlite", "file:db.lite?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		return fmt.Errorf("cant open sql database: %w", err)
	}
	defer sqldb.Close()

	entc := ent.NewClient(ent.Driver(entsql.OpenDB("sqlite3", sqldb)))
	defer entc.Close()

	logger, err := common.LoggerFromEnv()
	if err != nil {
		return fmt.Errorf("setup logger %w", err)
	}

	return api.Run(context.Background(), entc, logger, natsURL)
}
