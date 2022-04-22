package api

import (
	"context"
	"time"

	"github.com/encero/reciper/ent"
	"github.com/encero/reciper/ent/cookinghistory"
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

	LastCookedAt *time.Time `json:"lastCookedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func EntToRecipe(ctx context.Context, r *ent.Recipe) (Recipe, error) {
	out := Recipe{
		ID:      r.ID,
		Name:    r.Title,
		Planned: r.Planned,
	}

	lastHistory, err := r.QueryHistory().
		Order(ent.Desc(cookinghistory.FieldCookedAt)).
		First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return Recipe{}, err
	}

	if lastHistory != nil {
		out.LastCookedAt = &lastHistory.CookedAt
	}

	return out, nil
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
