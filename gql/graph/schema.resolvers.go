package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/encero/reciper-api/api"
	"github.com/encero/reciper-api/gql/graph/generated"
	"github.com/encero/reciper-api/gql/graph/model"
	"github.com/google/uuid"
)

func (r *mutationResolver) CreateRecipe(ctx context.Context, input model.NewRecipe) (*model.Recipe, error) {
	recipe := api.Recipe{
		ID:   uuid.New(),
		Name: input.Name,
	}

	resp := api.Envelope[api.Ack]{}

	err := r.ec.Request(api.HandlersRecipesUpsert, recipe, &resp, time.Second)
	if err != nil {
		return nil, fmt.Errorf("recipe upsert: %w", err)
	}

	return &model.Recipe{
		ID:   input.ID,
		Name: input.Name,
	}, nil
}

func (r *queryResolver) Recipes(ctx context.Context) ([]*model.Recipe, error) {
	resp := api.Envelope[api.List]{}

	err := r.ec.Request(api.HandlersRecipeList, nil, &resp, time.Second)
	if err != nil {
		return nil, fmt.Errorf("recipe list: %w", err)
	}

	out := []*model.Recipe{}
	for _, recipe := range resp.Data {
		out = append(out, &model.Recipe{
			ID:   recipe.ID.String(),
			Name: recipe.Name,
		})
	}

	return out, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }