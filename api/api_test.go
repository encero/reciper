package api_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/encero/reciper-api/api"
	"github.com/encero/reciper-api/ent"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/matryer/try"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	natsserver "github.com/nats-io/nats-server/v2/test"
)

const reqTimeout = time.Millisecond * 5

func TestCreateRecipe(t *testing.T) {
	is, conn, cleanup := setup(t)
	defer cleanup()

	id := upsertRecipe(is, conn, api.Recipe{
		ID:   uuid.New(),
		Name: "The name",
	})

	r := getRecipe(is, conn, id) // fetch recipe

	is.Equal(r.Name, "The name") // r.Name
}

func TestListRecipes(t *testing.T) {
	is, conn, cleanup := setup(t)
	defer cleanup()

	var ids []uuid.UUID

	for i := 0; i < 10; i++ {
		id := upsertRecipe(is, conn, api.Recipe{
			ID:   uuid.New(),
			Name: fmt.Sprintf("The name #%d", i),
		}) // create recipe

		ids = append(ids, id)
	}

	list := listRecipes(is, conn) // list recipes

	is.Equal(len(list), 10) // count of recipes

	listIds := make([]uuid.UUID, 0, 10)
	for _, r := range list {
		listIds = append(listIds, r.ID)
	}

	sortUUIDs(ids)
	sortUUIDs(listIds)

	is.Equal(ids, listIds) // same recipec in list
}

func TestUpdateRecipe(t *testing.T) {
	is, conn, cleanup := setup(t)
	defer cleanup()

	id := upsertRecipe(is, conn, api.Recipe{
		ID:   uuid.New(),
		Name: "The name",
	})

	_ = upsertRecipe(is, conn, api.Recipe{
		ID:   id,
		Name: "Different name",
	})

	r := getRecipe(is, conn, id)

	is.Equal(r.Name, "Different name") // r.Name
}

func TestDeleteRecipe(t *testing.T) {
	is, conn, cleanup := setup(t)
	defer cleanup()

	id1 := upsertRecipe(is, conn, api.Recipe{
		ID:   uuid.New(),
		Name: "The name",
	})
	id2 := upsertRecipe(is, conn, api.Recipe{
		ID:   uuid.New(),
		Name: "The name",
	})

	list := listRecipes(is, conn)
	is.Equal(len(list), 2) // two recipes after upsert

	_, err := conn.Request(fmt.Sprintf("recipes.delete.%s", id2), nil, reqTimeout)
	is.NoErr(err) // delete recipe

	_, err = conn.Request(fmt.Sprintf("recipes.delete.%s", id2), nil, reqTimeout)
	is.NoErr(err) // delete recipe again, should be noop

	list = listRecipes(is, conn)
	is.Equal(len(list), 1)    // one recipe ater delter
	is.Equal(list[0].ID, id1) // correct recipe remains
}

func runNatsServer(t *testing.T) (*server.Server, string) {
	s := natsserver.RunRandClientPortServer()

	info := s.PortsInfo(time.Second)

	if len(info.Nats) == 0 {
		t.Fatalf("no nats ports")
	}

	return s, info.Nats[0]
}

func runAndConnectNats(t *testing.T) (*nats.Conn, string, func()) {
	s, url := runNatsServer(t)

	conn, err := nats.Connect(url)
	if err != nil {
		t.Fatalf("nats connect: %s", err)
	}

	return conn, url, func() {
		conn.Close()
		s.Shutdown()
	}
}

func setup(t *testing.T) (*is.I, *nats.Conn, func()) {
	is := is.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	conn, natsURL, natsCleanup := runAndConnectNats(t)
	entc, entCleanup := testDatabase(t)
	serverDone := make(chan struct{})

	go func() {
		err := api.Run(ctx, entc, natsURL)
		if err != nil {
			fmt.Printf("api.RUN %v\n", err)
		}

		close(serverDone)
	}()

	return is, conn, func() {
		entCleanup()
		natsCleanup()
		cancel()

		<-serverDone
	}
}

func upsertRecipe(is *is.I, conn *nats.Conn, r api.Recipe) uuid.UUID {
	is.Helper()

	data, err := json.Marshal(r)
	is.NoErr(err) // new recipe marshal

	var response *nats.Msg

	err = try.Do(func(attempt int) (bool, error) {
		var err error
		response, err = conn.Request("recipes.upsert", data, reqTimeout)

		if err != nil {
			time.Sleep(time.Millisecond * 10)
		}

		return attempt < 10, err
	})
	is.NoErr(err) // request recipes.upsert

	resp := api.Ack{}

	err = json.Unmarshal(response.Data, &resp)
	is.NoErr(err) // unsmarshal create response

	is.Equal(resp.Status, api.StatusSuccess) // recipes.upsert status

	return r.ID
}

func getRecipe(is *is.I, conn *nats.Conn, id uuid.UUID) api.Recipe {
	is.Helper()

	response, err := conn.Request(fmt.Sprintf("recipes.detail.%s", id), nil, reqTimeout)
	is.NoErr(err) // request detail

	envelope := api.Envelope[api.Recipe]{}
	err = json.Unmarshal(response.Data, &envelope)
	is.NoErr(err)                                // unmarshal detail
	is.Equal(envelope.Status, api.StatusSuccess) // recipes.detail status

	return envelope.Data
}

func listRecipes(is *is.I, conn *nats.Conn) api.List {
	is.Helper()

	resp, err := conn.Request("recipes.list", nil, reqTimeout)
	is.NoErr(err) // request recipes.list

	var envelope api.Envelope[api.List]
	err = json.Unmarshal(resp.Data, &envelope)
	is.NoErr(err)                                // unmarshaling recipes.list
	is.Equal(envelope.Status, api.StatusSuccess) // recipes.list status

	return envelope.Data
}

func sortUUIDs(ids []uuid.UUID) {
	sort.Slice(ids, func(i, j int) bool {
		return strings.Compare(ids[i].String(), ids[j].String()) > 0
	})
}

func testDatabase(t *testing.T) (*ent.Client, func()) {
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
