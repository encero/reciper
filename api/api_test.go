package api_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/encero/reciper-api/api"
	"github.com/encero/reciper-api/pkg/tests"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/matryer/try"
	"github.com/nats-io/nats.go"
)

const reqTimeout = time.Millisecond * 5

func TestCreateRecipe(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	id := upsertRecipe(is, conn, api.Recipe{
		ID:   uuid.New(),
		Name: "The name",
	})

	r := getRecipe(is, conn, id) // fetch recipe

	is.Equal(r.Name, "The name") // r.Name
}

func TestListRecipes(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
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

func TestListRecipes_ReturnsEmptyArray(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	list := listRecipes(is, conn)

	is.True(list != nil) // empty List response Data should not be nil
}

func TestUpdateRecipe(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
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
	is, conn, cleanup := tests.SetupAPI(t)
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

func TestMarkRecipeAsPlanned(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	id := upsertRecipe(is, conn, api.Recipe{
		ID:   uuid.New(),
		Name: "The name",
	})

	recipes := listRecipes(is, conn)
	is.Equal(len(recipes), 1)
	is.Equal(recipes[0].Planned, false) // new recipe is unplanned

	markAsPlanned(is, conn, id, true) // mark recipe as planned

	recipes = listRecipes(is, conn)
	is.Equal(len(recipes), 1)
	is.True(recipes[0].Planned) // marked recipe is planned

	markAsPlanned(is, conn, id, false) // mark recipe as unplanned

	recipes = listRecipes(is, conn)
	is.Equal(len(recipes), 1)
	is.True(!recipes[0].Planned) // marked recipe is UNplanned
}

func markAsPlanned(is *is.I, conn *nats.Conn, id uuid.UUID, planned bool) {
	is.Helper()

	req := api.RequestPlanned{Planned: planned}

	payload, err := json.Marshal(req)
	is.NoErr(err) // marshalling recipes planned request

	msg, err := conn.Request(fmt.Sprintf("recipes.planned.%s", id.String()), payload, reqTimeout)
	is.NoErr(err) // request

	resp := api.Ack{}

	err = json.Unmarshal(msg.Data, &resp)
	is.NoErr(err) // response unmarshall

	is.Equal(resp.Status, api.StatusSuccess) // recipes.planned request status
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
