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

func TestRecipePlanned(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())
	defer gqlCleanup()

	id := createRecipe(t, "the name")

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

	t.Log("plan recipe", id)

	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	defer resp.Body.Close()

	data := struct {
		Data struct {
			PlanRecipe struct {
				Status string `json:"status"`
			} `json:"planRecipe"`
		} `json:"data"`
	}{}

	read(t, resp.Body, &data)

	is.Equal(data.Data.PlanRecipe.Status, "Success")
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

	id := uuid.New()

	q := query{
		Query: `mutation ($id: ID!, $name: String!){
                    createRecipe( input: {
                        id: $id,
                        name: $name
                    }) {
                        id
                        name
                    }
                }`,
		Variables: map[string]interface{}{
			"id":   id,
			"name": name,
		},
	}

	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)

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
	ID   string `json:"id"`
	Name string `json:"name"`
}

func listRecipes(t *testing.T) []recipe {
	is := is.New(t)

	q := query{
		Query: `query {
            recipes {
                id
                name
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
