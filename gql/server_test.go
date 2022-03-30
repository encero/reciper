package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

	created := createRecipe(t, "the name")

	_, err := uuid.Parse(created.ID)
	is.NoErr(err)
	is.Equal(created.Title, "the name")

	recipes := listRecipes(t)

	is.Equal(len(recipes), 1)

	recipe := recipes[0]
	is.Equal(recipe.Title, "the name")
}

func TestAddMoreRecipes(t *testing.T) {
	is, conn, cleanup := tests.SetupAPI(t)
	defer cleanup()

	gqlCleanup := setupGQL(t, conn.ConnectedUrl())
	defer gqlCleanup()

	created := createRecipe(t, "the name")

	_, err := uuid.Parse(created.ID)
	is.NoErr(err)
	is.Equal(created.Title, "the name")

	created = createRecipe(t, "the second name")

	_, err = uuid.Parse(created.ID)
	is.NoErr(err)
	is.Equal(created.Title, "the second name")

	recipes := listRecipes(t)

	is.Equal(len(recipes), 2)

	recipe := recipes[0]
	is.Equal(recipe.Title, "the name")

	recipe = recipes[1]
	is.Equal(recipe.Title, "the second name")
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

type createdRecipe struct {
	ID    string `json:"id"`
	Title string `json:"name"`
}

func createRecipe(t *testing.T, name string) createdRecipe {
	is := is.New(t)
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
			"id":   uuid.New().String(),
			"name": name,
		},
	}

	resp, err := http.Post("http://localhost:8080/query", "application/json", q.Marshal())
	is.NoErr(err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	is.NoErr(err)

	t.Log("response body", string(body))

	data := struct {
		Data struct {
			CreatedRecipe createdRecipe `json:"createRecipe"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &data)
	is.NoErr(err)

	is.Equal(resp.StatusCode, http.StatusOK)

	return data.Data.CreatedRecipe
}

type recipe struct {
	ID    string `json:"id"`
	Title string `json:"name"`
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
