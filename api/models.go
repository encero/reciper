package api

import (
	"time"

	"github.com/encero/reciper-api/ent"
	"github.com/google/uuid"
)

type Envelope[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type List []Recipe

type Recipe struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Planned bool      `json:"planned"`

	CreatedAt time.Time `json:"createdAt"`
}

func EntToRecipe(r *ent.Recipe) Recipe {
	return Recipe{
		ID:      r.ID,
		Name:    r.Title,
		Planned: r.Planned,
	}
}

type RequestPlanned struct {
	Planned bool `json:"planned"`
}

type Ack struct {
	Status string `json:"status"`
}

const StatusSuccess = "success"
const StatusError = "error"
const StatusNotFound = "notfound"
