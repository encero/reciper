// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/encero/reciper-api/ent/recipe"
	"github.com/google/uuid"
)

// RecipeCreate is the builder for creating a Recipe entity.
type RecipeCreate struct {
	config
	mutation *RecipeMutation
	hooks    []Hook
}

// SetTitle sets the "title" field.
func (rc *RecipeCreate) SetTitle(s string) *RecipeCreate {
	rc.mutation.SetTitle(s)
	return rc
}

// SetPlanned sets the "planned" field.
func (rc *RecipeCreate) SetPlanned(b bool) *RecipeCreate {
	rc.mutation.SetPlanned(b)
	return rc
}

// SetNillablePlanned sets the "planned" field if the given value is not nil.
func (rc *RecipeCreate) SetNillablePlanned(b *bool) *RecipeCreate {
	if b != nil {
		rc.SetPlanned(*b)
	}
	return rc
}

// SetID sets the "id" field.
func (rc *RecipeCreate) SetID(u uuid.UUID) *RecipeCreate {
	rc.mutation.SetID(u)
	return rc
}

// Mutation returns the RecipeMutation object of the builder.
func (rc *RecipeCreate) Mutation() *RecipeMutation {
	return rc.mutation
}

// Save creates the Recipe in the database.
func (rc *RecipeCreate) Save(ctx context.Context) (*Recipe, error) {
	var (
		err  error
		node *Recipe
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RecipeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			if rc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RecipeCreate) SaveX(ctx context.Context) *Recipe {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *RecipeCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *RecipeCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *RecipeCreate) defaults() {
	if _, ok := rc.mutation.Planned(); !ok {
		v := recipe.DefaultPlanned
		rc.mutation.SetPlanned(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RecipeCreate) check() error {
	if _, ok := rc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Recipe.title"`)}
	}
	if _, ok := rc.mutation.Planned(); !ok {
		return &ValidationError{Name: "planned", err: errors.New(`ent: missing required field "Recipe.planned"`)}
	}
	return nil
}

func (rc *RecipeCreate) sqlSave(ctx context.Context) (*Recipe, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (rc *RecipeCreate) createSpec() (*Recipe, *sqlgraph.CreateSpec) {
	var (
		_node = &Recipe{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: recipe.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: recipe.FieldID,
			},
		}
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := rc.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: recipe.FieldTitle,
		})
		_node.Title = value
	}
	if value, ok := rc.mutation.Planned(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: recipe.FieldPlanned,
		})
		_node.Planned = value
	}
	return _node, _spec
}

// RecipeCreateBulk is the builder for creating many Recipe entities in bulk.
type RecipeCreateBulk struct {
	config
	builders []*RecipeCreate
}

// Save creates the Recipe entities in the database.
func (rcb *RecipeCreateBulk) Save(ctx context.Context) ([]*Recipe, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Recipe, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RecipeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RecipeCreateBulk) SaveX(ctx context.Context) []*Recipe {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *RecipeCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *RecipeCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}
