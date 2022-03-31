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

	"github.com/encero/reciper-api/pkg/tests"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/matryer/try"
	"go.uber.org/zap"
)

func TestAddRecipe(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())
	defer gqlCleanup()

	_ = createRecipe(t, "the name")

	recipes := listRecipes(t)

	is.Equal(len(recipes), 1)

	recipe := recipes[0]
	is.Equal(recipe.Name, "the name")
}

func TestAddMoreRecipes(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())
	defer gqlCleanup()

	_ = createRecipe(t, "A the name")
	_ = createRecipe(t, "B the second name")

	recipes := listRecipes(t)

	is.Equal(len(recipes), 2)

	sort.Slice(recipes, func(i, j int) bool {
		return strings.Compare(recipes[i].Name, recipes[j].Name) < 0
	})

	recipe := recipes[0]
	is.Equal(recipe.Name, "A the name")

	recipe = recipes[1]
	is.Equal(recipe.Name, "B the second name")
}

func TestRecipePlanned_Validations(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())
	defer gqlCleanup()

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

	read(t, response.Body, &data)

	is.Equal(response.StatusCode, http.StatusOK)
	is.Equal(data.Errors[0].Message, "id must be a valid UUID")
}

func TestRecipePlanned(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())
	defer gqlCleanup()

	id := createRecipe(t, "the name")

	planRecipe(t, is, id)

	recipes := listRecipes(t)
	is.Equal(len(recipes), 1)          // count of recipes
	is.Equal(recipes[0].Planned, true) // recipe should be planned
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
		err := run(ctx, tl.With(zap.String("system", "gql")), "8080", natsURL)
		if !errors.Is(err, http.ErrServerClosed) {
			is.NoErr(err)
		}
	}()

	start := time.Now()
	err := try.Do(func(_ int) (bool, error) {
		resp, err := http.Get("http://localhost:8080/query")
		if err != nil {
			time.Sleep(time.Millisecond)
			return time.Since(start) < time.Millisecond*10, err
		}
		resp.Body.Close()

		return false, nil
	})
	is.NoErr(err)

	return cancel
}

func createRecipe(t *testing.T, name string) uuid.UUID {
	is := is.New(t)

	q := query{
		Query: `mutation ($name: String!){
                    createRecipe( input: {
                        name: $name
                    }) {
                        id
                        name
                    }
                }`,
		Variables: map[string]interface{}{
			"name": name,
		},
	}

	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	defer resp.Body.Close()

	data := struct {
		Data struct {
			CreateRecipe struct {
				ID string `json:"id"`
			} `json:"createRecipe"`
		} `json:"data"`
	}{}

	read(t, resp.Body, &data)
	is.Equal(resp.StatusCode, http.StatusOK)

	id, err := uuid.Parse(data.Data.CreateRecipe.ID)
	is.NoErr(err)

	return id
}

func read(t *testing.T, body io.Reader, to interface{}) {
	is := is.New(t)

	data, err := io.ReadAll(body)
	is.NoErr(err)

	t.Log("response body", string(data))

	err = json.Unmarshal(data, &to)
	is.NoErr(err)
}

type recipe struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Planned bool   `json:"planned"`
}

func listRecipes(t *testing.T) []recipe {
	is := is.New(t)

	q := query{
		Query: `query {
            recipes {
                id
                name
                planned
            }
        }`,
	}

	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	body, err := io.ReadAll(resp.Body)
	is.NoErr(err)

	t.Log("response body", string(body))

	is.Equal(resp.StatusCode, http.StatusOK)

	list := struct {
		Data struct {
			Recipes []recipe `json:"recipes"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &list)
	is.NoErr(err)

	return list.Data.Recipes
}

func planRecipe(t *testing.T, is *is.I, id uuid.UUID) {
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

	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err) // planRecipe mutation request

	defer resp.Body.Close()

	data := struct {
		Data struct {
			PlanRecipe struct {
				Status string `json:"status"`
			} `json:"planRecipe"`
		} `json:"data"`
	}{}

	read(t, resp.Body, &data)

	is.Equal(data.Data.PlanRecipe.Status, "Success") // planRecipe mutation status
}
