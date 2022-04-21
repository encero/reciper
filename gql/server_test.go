package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/encero/reciper/gql/configuration"
	"github.com/encero/reciper/pkg/tests"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/matryer/try"
	"go.uber.org/zap"
)

func setup(t *testing.T) (tests.IsT, func()) {
	is, conn, cleanup := tests.SetupAPI(t)

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())

	return is, func() {
		gqlCleanup()
		cleanup()
	}
}

func TestAddRecipe(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	_ = createRecipe(is, "the name")

	recipes := listRecipes(is)

	is.Equal(len(recipes), 1)

	recipe := recipes[0]
	is.Equal(recipe.Name, "the name")
}

func TestAddMoreRecipes(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	_ = createRecipe(is, "A the name")
	_ = createRecipe(is, "B the second name")

	recipes := listRecipes(is)

	is.Equal(len(recipes), 2)

	sort.Slice(recipes, func(i, j int) bool {
		return strings.Compare(recipes[i].Name, recipes[j].Name) < 0
	})

	recipe := recipes[0]
	is.Equal(recipe.Name, "A the name")

	recipe = recipes[1]
	is.Equal(recipe.Name, "B the second name")
}

func TestRecipePlanned(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	id := createRecipe(is, "the name")

	planRecipe(is, id)

	recipes := listRecipes(is)
	is.Equal(len(recipes), 1)          // count of recipes
	is.Equal(recipes[0].Planned, true) // recipe should be planned
}

func TestRecipePlanned_Validations(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	q := query{
		Query: `mutation {
            planRecipe(id: "not-uuid") {
                status
            }
        }`,
	}

	response, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	defer response.Body.Close()

	data := struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&data)
	is.NoErr(err)

	is.Equal(data.Errors[0].Message, "id must be a valid UUID")
}

func TestUpdateRecipe(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	id := createRecipe(is, "original name")

	updateRecipe(is, id, "new name")

	recipes := listRecipes(is)

	is.Equal(len(recipes), 1) // recipe count
	is.Equal(recipes[0].Name, "new name")
}

func TestDeleteRecipe(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	id := createRecipe(is, "the name")
	_ = createRecipe(is, "the second name")

	recipes := listRecipes(is)
	is.Equal(len(recipes), 2) // two recipes prepared

	deleteRecipe(is, id)

	recipes = listRecipes(is)
	is.Equal(len(recipes), 1)                    // one recipe left after deletion
	is.Equal(recipes[0].Name, "the second name") // left recipe is the correct one
}

func TestUnPlanRecipe(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	id := createRecipe(is, "the name")

	planRecipe(is, id)

	recipes := listRecipes(is)
	is.Equal(len(recipes), 1)   // expect one recipe
	is.True(recipes[0].Planned) // recipe is planned

	unPlanRecipe(is, id)

	recipes = listRecipes(is)
	is.Equal(len(recipes), 1)    // expect one recipe
	is.True(!recipes[0].Planned) // recipe is not planned
}

func TestCookRecipe(t *testing.T) {
	is, cleanup := setup(t)
	defer cleanup()

	id := createRecipe(is, "the name")

	planRecipe(is, id)

	recipes := listRecipes(is)
	is.Equal(len(recipes), 1)   // expect one recipe
	is.True(recipes[0].Planned) // recipe is planned

	cookRecipe(is, id)

	recipes = listRecipes(is)
	is.Equal(len(recipes), 1)    // expect one recipe
	is.True(!recipes[0].Planned) // recipe is not planned
}

func TestApiStatus(t *testing.T) {
	name := uuid.New().String()

	t.Setenv("SERVER_NAME", name)

	is, cleanup := setup(t)
	defer cleanup()

	q := query{
		Query: "query { apiStatus {name, ref}}",
	}

	data := struct {
		APIStatus struct {
			Name string `json:"name"`
			Ref  string `json:"ref"`
		} `json:"apiStatus"`
	}{}

	performQuery(is, q, &data)

	is.Equal(data.APIStatus.Ref, "gql-test") // we received some version information
}

////////////////////////////////////////////
// HELPERS
////////////////////////////////////////////

type recipe struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Planned bool      `json:"planned"`
}

type query struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (q query) Marshal() io.Reader {
	data := &bytes.Buffer{}

	err := json.NewEncoder(data).Encode(q)
	if err != nil {
		panic(err)
	}

	return data
}

func setupGQL(t *testing.T, natsURL string) func() {
	is := is.New(t)
	tl := tests.TestLogger(t)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		cfg, err := configuration.FromEnvironment()
		if err != nil {
			is.NoErr(err) // environment load for gql server
		}

		cfg.VersionRef = "test"

		cfg.ServerPort = "8080"
		cfg.NatsURL = natsURL

		err = run(ctx, tl.With(zap.String("system", "gql")), cfg)
		if !errors.Is(err, http.ErrServerClosed) {
			is.NoErr(err)
		}

		t.Log("server start failure:", err)
	}()

	start := time.Now()
	err := try.Do(func(_ int) (bool, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()

		t.Log("waiting for gql server to start")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/query", nil)
		is.NoErr(err)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			time.Sleep(time.Millisecond)
			return time.Since(start) < time.Millisecond*100, err
		}
		resp.Body.Close()

		return false, nil
	})
	is.NoErr(err)

	return cancel
}

func performQuery(is tests.IsT, q query, out interface{}) {
	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	is.NoErr(err)

	is.T.Log("response body", string(data))

	envelope := struct {
		Data interface{}
	}{
		Data: out,
	}

	err = json.Unmarshal(data, &envelope)
	is.NoErr(err)
}

func createRecipe(is tests.IsT, name string) uuid.UUID {
	q := query{
		Query: `
        mutation($name: String!) {
            createRecipe(input: {name: $name}) {
                id
                name
                planned
            }
        }
        `,
		Variables: map[string]interface{}{
			"name": name,
		},
	}

	out := struct {
		CreateRecipe recipe `json:"createRecipe"`
	}{}

	performQuery(is, q, &out)

	return out.CreateRecipe.ID
}

func deleteRecipe(is tests.IsT, id uuid.UUID) uuid.UUID {
	q := query{
		Query: `mutation ($id: ID!){
            deleteRecipe( id: $id ) {
                status
            }
        }`,
		Variables: map[string]interface{}{
			"id": id,
		},
	}

	data := struct {
		RecipeDelete struct {
			Status string `json:"status"`
		} `json:"deleteRecipe"`
	}{}

	performQuery(is, q, &data)

	is.Equal(data.RecipeDelete.Status, "Success")

	return id
}

func updateRecipe(is tests.IsT, id uuid.UUID, name string) {
	q := query{
		Query: `mutation ($id: ID!, $name: String!){
                    updateRecipe( input: {
                        id: $id
                        name: $name
                    }) {
                        status
                    }
                }`,
		Variables: map[string]interface{}{
			"id":   id,
			"name": name,
		},
	}

	data := struct {
		RecipeUpdate struct {
			Status string `json:"status"`
		} `json:"updateRecipe"`
	}{}

	performQuery(is, q, &data)

	is.Equal(data.RecipeUpdate.Status, "Success")
}

func listRecipes(is tests.IsT) []recipe {
	q := query{
		Query: `query {
            recipes {
                id
                name
                planned
            }
        }`,
	}

	list := struct {
		Recipes []recipe `json:"recipes"`
	}{}

	performQuery(is, q, &list)

	return list.Recipes
}

func planRecipe(is tests.IsT, id uuid.UUID) {
	q := query{
		Query: `mutation ($id: ID!){
            planRecipe(id: $id) {
                status
            }
        }`,
		Variables: map[string]interface{}{
			"id": id,
		},
	}

	data := struct {
		PlanRecipe struct {
			Status string `json:"status"`
		} `json:"planRecipe"`
	}{}

	performQuery(is, q, &data)

	is.Equal(data.PlanRecipe.Status, "Success") // planRecipe mutation status
}

func unPlanRecipe(is tests.IsT, id uuid.UUID) {
	q := query{
		Query: `mutation ($id: ID!){
            unPlanRecipe(id: $id) {
                status
            }
        }`,
		Variables: map[string]interface{}{
			"id": id,
		},
	}

	data := struct {
		UnPlanRecipe struct {
			Status string `json:"status"`
		} `json:"unPlanRecipe"`
	}{}

	performQuery(is, q, &data)

	is.Equal(data.UnPlanRecipe.Status, "Success") // unPlanRecipe mutation status
}

func cookRecipe(is tests.IsT, id uuid.UUID) {
	q := query{
		Query: `mutation ($id: ID!){
            cookRecipe(id: $id) {
                status
            }
        }`,
		Variables: map[string]interface{}{
			"id": id,
		},
	}

	data := struct {
		CookRecipe struct {
			Status string `json:"status"`
		} `json:"cookRecipe"`
	}{}

	performQuery(is, q, &data)

	is.Equal(data.CookRecipe.Status, "Success") // cookRecipe mutation status
}
