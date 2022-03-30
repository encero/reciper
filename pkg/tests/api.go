package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/encero/reciper-api/api"
	"github.com/encero/reciper-api/ent"
	"github.com/matryer/is"
	"github.com/matryer/try"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	_ "modernc.org/sqlite" // intentional for tests
)

func SetupAPI(t *testing.T) (*is.I, *nats.Conn, func()) {
	is := is.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	conn, natsURL, natsCleanup := RunAndConnectNats(t)
	entc, entCleanup := TestDatabase(t)
	serverDone := make(chan struct{})

	go func() {
		lg := TestLogger(t).With(zap.String("system", "api"))
		err := api.Run(ctx, entc, lg, natsURL)

		if err != nil {
			fmt.Printf("api.RUN %v\n", err)
		}

		close(serverDone)
	}()

	start := time.Now()
	err := try.Do(func(_ int) (retry bool, err error) {
		_, err = conn.Request(api.HandlersRecipeList, nil, time.Second)

		if err != nil {
			time.Sleep(time.Millisecond * 10)
		}

		return time.Since(start) < time.Second, err
	})

	if err != nil {
		is.NoErr(err) // API not responding
	}

	return is, conn, func() {
		entCleanup()
		natsCleanup()
		cancel()

		<-serverDone
	}
}

func TestDatabase(t *testing.T) (*ent.Client, func()) {
	sqldb, err := sql.Open("sqlite", "file:ent?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}

	entc := ent.NewClient(ent.Driver(entsql.OpenDB("sqlite3", sqldb)))

	// Run the auto migration tool.
	if err := entc.Schema.Create(context.Background()); err != nil {
		sqldb.Close()
		t.Fatalf("failed creating schema resources: %v", err)
	}

	return entc, func() {
		entc.Close()
		sqldb.Close()
	}
}
