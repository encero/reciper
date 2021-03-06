// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/encero/reciper/ent/cookinghistory"
	"github.com/encero/reciper/ent/predicate"
	"github.com/encero/reciper/ent/recipe"
	"github.com/google/uuid"
)

// CookingHistoryUpdate is the builder for updating CookingHistory entities.
type CookingHistoryUpdate struct {
	config
	hooks    []Hook
	mutation *CookingHistoryMutation
}

// Where appends a list predicates to the CookingHistoryUpdate builder.
func (chu *CookingHistoryUpdate) Where(ps ...predicate.CookingHistory) *CookingHistoryUpdate {
	chu.mutation.Where(ps...)
	return chu
}

// SetCookedAt sets the "cookedAt" field.
func (chu *CookingHistoryUpdate) SetCookedAt(t time.Time) *CookingHistoryUpdate {
	chu.mutation.SetCookedAt(t)
	return chu
}

// SetRecipeID sets the "recipe" edge to the Recipe entity by ID.
func (chu *CookingHistoryUpdate) SetRecipeID(id uuid.UUID) *CookingHistoryUpdate {
	chu.mutation.SetRecipeID(id)
	return chu
}

// SetRecipe sets the "recipe" edge to the Recipe entity.
func (chu *CookingHistoryUpdate) SetRecipe(r *Recipe) *CookingHistoryUpdate {
	return chu.SetRecipeID(r.ID)
}

// Mutation returns the CookingHistoryMutation object of the builder.
func (chu *CookingHistoryUpdate) Mutation() *CookingHistoryMutation {
	return chu.mutation
}

// ClearRecipe clears the "recipe" edge to the Recipe entity.
func (chu *CookingHistoryUpdate) ClearRecipe() *CookingHistoryUpdate {
	chu.mutation.ClearRecipe()
	return chu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (chu *CookingHistoryUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(chu.hooks) == 0 {
		if err = chu.check(); err != nil {
			return 0, err
		}
		affected, err = chu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CookingHistoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = chu.check(); err != nil {
				return 0, err
			}
			chu.mutation = mutation
			affected, err = chu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(chu.hooks) - 1; i >= 0; i-- {
			if chu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = chu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, chu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (chu *CookingHistoryUpdate) SaveX(ctx context.Context) int {
	affected, err := chu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (chu *CookingHistoryUpdate) Exec(ctx context.Context) error {
	_, err := chu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (chu *CookingHistoryUpdate) ExecX(ctx context.Context) {
	if err := chu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (chu *CookingHistoryUpdate) check() error {
	if _, ok := chu.mutation.RecipeID(); chu.mutation.RecipeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "CookingHistory.recipe"`)
	}
	return nil
}

func (chu *CookingHistoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cookinghistory.Table,
			Columns: cookinghistory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: cookinghistory.FieldID,
			},
		},
	}
	if ps := chu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := chu.mutation.CookedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cookinghistory.FieldCookedAt,
		})
	}
	if chu.mutation.RecipeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cookinghistory.RecipeTable,
			Columns: []string{cookinghistory.RecipeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: recipe.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := chu.mutation.RecipeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cookinghistory.RecipeTable,
			Columns: []string{cookinghistory.RecipeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: recipe.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, chu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cookinghistory.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// CookingHistoryUpdateOne is the builder for updating a single CookingHistory entity.
type CookingHistoryUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CookingHistoryMutation
}

// SetCookedAt sets the "cookedAt" field.
func (chuo *CookingHistoryUpdateOne) SetCookedAt(t time.Time) *CookingHistoryUpdateOne {
	chuo.mutation.SetCookedAt(t)
	return chuo
}

// SetRecipeID sets the "recipe" edge to the Recipe entity by ID.
func (chuo *CookingHistoryUpdateOne) SetRecipeID(id uuid.UUID) *CookingHistoryUpdateOne {
	chuo.mutation.SetRecipeID(id)
	return chuo
}

// SetRecipe sets the "recipe" edge to the Recipe entity.
func (chuo *CookingHistoryUpdateOne) SetRecipe(r *Recipe) *CookingHistoryUpdateOne {
	return chuo.SetRecipeID(r.ID)
}

// Mutation returns the CookingHistoryMutation object of the builder.
func (chuo *CookingHistoryUpdateOne) Mutation() *CookingHistoryMutation {
	return chuo.mutation
}

// ClearRecipe clears the "recipe" edge to the Recipe entity.
func (chuo *CookingHistoryUpdateOne) ClearRecipe() *CookingHistoryUpdateOne {
	chuo.mutation.ClearRecipe()
	return chuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (chuo *CookingHistoryUpdateOne) Select(field string, fields ...string) *CookingHistoryUpdateOne {
	chuo.fields = append([]string{field}, fields...)
	return chuo
}

// Save executes the query and returns the updated CookingHistory entity.
func (chuo *CookingHistoryUpdateOne) Save(ctx context.Context) (*CookingHistory, error) {
	var (
		err  error
		node *CookingHistory
	)
	if len(chuo.hooks) == 0 {
		if err = chuo.check(); err != nil {
			return nil, err
		}
		node, err = chuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CookingHistoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = chuo.check(); err != nil {
				return nil, err
			}
			chuo.mutation = mutation
			node, err = chuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(chuo.hooks) - 1; i >= 0; i-- {
			if chuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = chuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, chuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (chuo *CookingHistoryUpdateOne) SaveX(ctx context.Context) *CookingHistory {
	node, err := chuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (chuo *CookingHistoryUpdateOne) Exec(ctx context.Context) error {
	_, err := chuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (chuo *CookingHistoryUpdateOne) ExecX(ctx context.Context) {
	if err := chuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (chuo *CookingHistoryUpdateOne) check() error {
	if _, ok := chuo.mutation.RecipeID(); chuo.mutation.RecipeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "CookingHistory.recipe"`)
	}
	return nil
}

func (chuo *CookingHistoryUpdateOne) sqlSave(ctx context.Context) (_node *CookingHistory, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   cookinghistory.Table,
			Columns: cookinghistory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: cookinghistory.FieldID,
			},
		},
	}
	id, ok := chuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "CookingHistory.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := chuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cookinghistory.FieldID)
		for _, f := range fields {
			if !cookinghistory.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cookinghistory.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := chuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := chuo.mutation.CookedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: cookinghistory.FieldCookedAt,
		})
	}
	if chuo.mutation.RecipeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cookinghistory.RecipeTable,
			Columns: []string{cookinghistory.RecipeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: recipe.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := chuo.mutation.RecipeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   cookinghistory.RecipeTable,
			Columns: []string{cookinghistory.RecipeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: recipe.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &CookingHistory{config: chuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, chuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cookinghistory.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
