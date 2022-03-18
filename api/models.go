package api

import (
	"time"

	"github.com/encero/reciper-api/ent"
	"github.com/google/uuid"
)

type List struct {
	Status  string   `json:"status"`
	Recipes []Recipe `json:"recipes"`
}

type Recipe struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
}

func EntToRecipe(r *ent.Recipe) Recipe {
	return Recipe{
		ID:   r.ID,
		Name: r.Title,
	}
}

type Ack struct {
	Status string `json:"status"`
}

const StatusSuccess = "success"
const StatusError = "error"
