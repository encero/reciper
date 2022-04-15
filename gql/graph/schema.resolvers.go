package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/encero/reciper/api"
	"github.com/encero/reciper/gql/graph/generated"
	"github.com/encero/reciper/gql/graph/model"
	"github.com/google/uuid"
)

func (r *mutationResolver) CreateRecipe(ctx context.Context, input model.NewRecipe) (*model.Recipe, error) {
	recipe := api.Recipe{
		ID:   uuid.New(),
		Name: input.Name,
	}

	resp := api.Envelope[api.Recipe]{}

	err := r.ec.Request(api.HandlersRecipesUpsert, recipe, &resp, time.Second)
	if err != nil {
		return nil, fmt.Errorf("recipe upsert: %w", err)
	}

	return &model.Recipe{
		ID:      resp.Data.ID.String(),
		Name:    resp.Data.Name,
		Planned: resp.Data.Planned,
	}, nil
}

func (r *mutationResolver) UpdateRecipe(ctx context.Context, input *model.UpdateRecipe) (*model.Result, error) {
	id := uuid.MustParse(input.ID)

	recipe := api.Recipe{
		ID:   id,
		Name: input.Name,
	}
	resp := api.Ack{}

	err := r.ec.Request(api.HandlersRecipesUpsert, recipe, &resp, time.Second)
	if err != nil {
		return nil, fmt.Errorf("recipe upsert %w", err)
	}

	return statusToResult(resp.Status)
}

func (r *mutationResolver) DeleteRecipe(ctx context.Context, id string) (*model.Result, error) {
	resp := api.Ack{}

	err := r.ec.Request(fmt.Sprintf("recipes.delete.%s", id), nil, &resp, time.Second)
	if err != nil {
		return nil, fmt.Errorf("recipe upsert %w", err)
	}

	return statusToResult(resp.Status)
}

func (r *mutationResolver) PlanRecipe(ctx context.Context, id string) (*model.Result, error) {
	resp := api.Ack{}

	err := r.ec.Request(fmt.Sprintf("recipes.planned.%s", id), api.RequestPlanned{Planned: true}, &resp, time.Second)
	if err != nil {
		return &model.Result{Status: model.StatusError}, nil
	}

	return statusToResult(resp.Status)
}

func (r *mutationResolver) UnPlanRecipe(ctx context.Context, id string) (*model.Result, error) {
	resp := api.Ack{}

	err := r.ec.Request(fmt.Sprintf("recipes.planned.%s", id), api.RequestPlanned{Planned: false}, &resp, time.Second)
	if err != nil {
		return &model.Result{Status: model.StatusError}, nil
	}

	return statusToResult(resp.Status)
}

func (r *mutationResolver) CookRecipe(ctx context.Context, id string) (*model.Result, error) {
	resp := api.Ack{}

	err := r.ec.Request(fmt.Sprintf("recipes.planned.%s", id), api.RequestPlanned{Planned: false}, &resp, time.Second)
	if err != nil {
		return &model.Result{Status: model.StatusError}, nil
	}

	return statusToResult(resp.Status)
}

func (r *queryResolver) APIStatus(ctx context.Context) (*model.APIStatus, error) {
	return &model.APIStatus{
		Name:    r.cfg.ServerName,
		Version: fmt.Sprintf("gql-%s", r.cfg.Version),
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
			ID:      recipe.ID.String(),
			Name:    recipe.Name,
			Planned: recipe.Planned,
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
