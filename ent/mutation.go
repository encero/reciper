// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/encero/reciper/ent/cookinghistory"
	"github.com/encero/reciper/ent/predicate"
	"github.com/encero/reciper/ent/recipe"
	"github.com/google/uuid"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeCookingHistory = "CookingHistory"
	TypeRecipe         = "Recipe"
)

// CookingHistoryMutation represents an operation that mutates the CookingHistory nodes in the graph.
type CookingHistoryMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	cookedAt      *time.Time
	clearedFields map[string]struct{}
	recipe        *uuid.UUID
	clearedrecipe bool
	done          bool
	oldValue      func(context.Context) (*CookingHistory, error)
	predicates    []predicate.CookingHistory
}

var _ ent.Mutation = (*CookingHistoryMutation)(nil)

// cookinghistoryOption allows management of the mutation configuration using functional options.
type cookinghistoryOption func(*CookingHistoryMutation)

// newCookingHistoryMutation creates new mutation for the CookingHistory entity.
func newCookingHistoryMutation(c config, op Op, opts ...cookinghistoryOption) *CookingHistoryMutation {
	m := &CookingHistoryMutation{
		config:        c,
		op:            op,
		typ:           TypeCookingHistory,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withCookingHistoryID sets the ID field of the mutation.
func withCookingHistoryID(id uuid.UUID) cookinghistoryOption {
	return func(m *CookingHistoryMutation) {
		var (
			err   error
			once  sync.Once
			value *CookingHistory
		)
		m.oldValue = func(ctx context.Context) (*CookingHistory, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().CookingHistory.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withCookingHistory sets the old CookingHistory of the mutation.
func withCookingHistory(node *CookingHistory) cookinghistoryOption {
	return func(m *CookingHistoryMutation) {
		m.oldValue = func(context.Context) (*CookingHistory, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CookingHistoryMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CookingHistoryMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of CookingHistory entities.
func (m *CookingHistoryMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *CookingHistoryMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *CookingHistoryMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().CookingHistory.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCookedAt sets the "cookedAt" field.
func (m *CookingHistoryMutation) SetCookedAt(t time.Time) {
	m.cookedAt = &t
}

// CookedAt returns the value of the "cookedAt" field in the mutation.
func (m *CookingHistoryMutation) CookedAt() (r time.Time, exists bool) {
	v := m.cookedAt
	if v == nil {
		return
	}
	return *v, true
}

// OldCookedAt returns the old "cookedAt" field's value of the CookingHistory entity.
// If the CookingHistory object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CookingHistoryMutation) OldCookedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCookedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCookedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCookedAt: %w", err)
	}
	return oldValue.CookedAt, nil
}

// ResetCookedAt resets all changes to the "cookedAt" field.
func (m *CookingHistoryMutation) ResetCookedAt() {
	m.cookedAt = nil
}

// SetRecipeID sets the "recipe" edge to the Recipe entity by id.
func (m *CookingHistoryMutation) SetRecipeID(id uuid.UUID) {
	m.recipe = &id
}

// ClearRecipe clears the "recipe" edge to the Recipe entity.
func (m *CookingHistoryMutation) ClearRecipe() {
	m.clearedrecipe = true
}

// RecipeCleared reports if the "recipe" edge to the Recipe entity was cleared.
func (m *CookingHistoryMutation) RecipeCleared() bool {
	return m.clearedrecipe
}

// RecipeID returns the "recipe" edge ID in the mutation.
func (m *CookingHistoryMutation) RecipeID() (id uuid.UUID, exists bool) {
	if m.recipe != nil {
		return *m.recipe, true
	}
	return
}

// RecipeIDs returns the "recipe" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// RecipeID instead. It exists only for internal usage by the builders.
func (m *CookingHistoryMutation) RecipeIDs() (ids []uuid.UUID) {
	if id := m.recipe; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetRecipe resets all changes to the "recipe" edge.
func (m *CookingHistoryMutation) ResetRecipe() {
	m.recipe = nil
	m.clearedrecipe = false
}

// Where appends a list predicates to the CookingHistoryMutation builder.
func (m *CookingHistoryMutation) Where(ps ...predicate.CookingHistory) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *CookingHistoryMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (CookingHistory).
func (m *CookingHistoryMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *CookingHistoryMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.cookedAt != nil {
		fields = append(fields, cookinghistory.FieldCookedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *CookingHistoryMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case cookinghistory.FieldCookedAt:
		return m.CookedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *CookingHistoryMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case cookinghistory.FieldCookedAt:
		return m.OldCookedAt(ctx)
	}
	return nil, fmt.Errorf("unknown CookingHistory field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CookingHistoryMutation) SetField(name string, value ent.Value) error {
	switch name {
	case cookinghistory.FieldCookedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCookedAt(v)
		return nil
	}
	return fmt.Errorf("unknown CookingHistory field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *CookingHistoryMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *CookingHistoryMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CookingHistoryMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown CookingHistory numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *CookingHistoryMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *CookingHistoryMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *CookingHistoryMutation) ClearField(name string) error {
	return fmt.Errorf("unknown CookingHistory nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *CookingHistoryMutation) ResetField(name string) error {
	switch name {
	case cookinghistory.FieldCookedAt:
		m.ResetCookedAt()
		return nil
	}
	return fmt.Errorf("unknown CookingHistory field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *CookingHistoryMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.recipe != nil {
		edges = append(edges, cookinghistory.EdgeRecipe)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *CookingHistoryMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case cookinghistory.EdgeRecipe:
		if id := m.recipe; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *CookingHistoryMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *CookingHistoryMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *CookingHistoryMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedrecipe {
		edges = append(edges, cookinghistory.EdgeRecipe)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *CookingHistoryMutation) EdgeCleared(name string) bool {
	switch name {
	case cookinghistory.EdgeRecipe:
		return m.clearedrecipe
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *CookingHistoryMutation) ClearEdge(name string) error {
	switch name {
	case cookinghistory.EdgeRecipe:
		m.ClearRecipe()
		return nil
	}
	return fmt.Errorf("unknown CookingHistory unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *CookingHistoryMutation) ResetEdge(name string) error {
	switch name {
	case cookinghistory.EdgeRecipe:
		m.ResetRecipe()
		return nil
	}
	return fmt.Errorf("unknown CookingHistory edge %s", name)
}

// RecipeMutation represents an operation that mutates the Recipe nodes in the graph.
type RecipeMutation struct {
	config
	op             Op
	typ            string
	id             *uuid.UUID
	title          *string
	planned        *bool
	clearedFields  map[string]struct{}
	history        map[uuid.UUID]struct{}
	removedhistory map[uuid.UUID]struct{}
	clearedhistory bool
	done           bool
	oldValue       func(context.Context) (*Recipe, error)
	predicates     []predicate.Recipe
}

var _ ent.Mutation = (*RecipeMutation)(nil)

// recipeOption allows management of the mutation configuration using functional options.
type recipeOption func(*RecipeMutation)

// newRecipeMutation creates new mutation for the Recipe entity.
func newRecipeMutation(c config, op Op, opts ...recipeOption) *RecipeMutation {
	m := &RecipeMutation{
		config:        c,
		op:            op,
		typ:           TypeRecipe,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withRecipeID sets the ID field of the mutation.
func withRecipeID(id uuid.UUID) recipeOption {
	return func(m *RecipeMutation) {
		var (
			err   error
			once  sync.Once
			value *Recipe
		)
		m.oldValue = func(ctx context.Context) (*Recipe, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Recipe.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withRecipe sets the old Recipe of the mutation.
func withRecipe(node *Recipe) recipeOption {
	return func(m *RecipeMutation) {
		m.oldValue = func(context.Context) (*Recipe, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RecipeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RecipeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Recipe entities.
func (m *RecipeMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *RecipeMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *RecipeMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Recipe.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetTitle sets the "title" field.
func (m *RecipeMutation) SetTitle(s string) {
	m.title = &s
}

// Title returns the value of the "title" field in the mutation.
func (m *RecipeMutation) Title() (r string, exists bool) {
	v := m.title
	if v == nil {
		return
	}
	return *v, true
}

// OldTitle returns the old "title" field's value of the Recipe entity.
// If the Recipe object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RecipeMutation) OldTitle(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTitle is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTitle requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTitle: %w", err)
	}
	return oldValue.Title, nil
}

// ResetTitle resets all changes to the "title" field.
func (m *RecipeMutation) ResetTitle() {
	m.title = nil
}

// SetPlanned sets the "planned" field.
func (m *RecipeMutation) SetPlanned(b bool) {
	m.planned = &b
}

// Planned returns the value of the "planned" field in the mutation.
func (m *RecipeMutation) Planned() (r bool, exists bool) {
	v := m.planned
	if v == nil {
		return
	}
	return *v, true
}

// OldPlanned returns the old "planned" field's value of the Recipe entity.
// If the Recipe object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RecipeMutation) OldPlanned(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPlanned is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPlanned requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPlanned: %w", err)
	}
	return oldValue.Planned, nil
}

// ResetPlanned resets all changes to the "planned" field.
func (m *RecipeMutation) ResetPlanned() {
	m.planned = nil
}

// AddHistoryIDs adds the "history" edge to the CookingHistory entity by ids.
func (m *RecipeMutation) AddHistoryIDs(ids ...uuid.UUID) {
	if m.history == nil {
		m.history = make(map[uuid.UUID]struct{})
	}
	for i := range ids {
		m.history[ids[i]] = struct{}{}
	}
}

// ClearHistory clears the "history" edge to the CookingHistory entity.
func (m *RecipeMutation) ClearHistory() {
	m.clearedhistory = true
}

// HistoryCleared reports if the "history" edge to the CookingHistory entity was cleared.
func (m *RecipeMutation) HistoryCleared() bool {
	return m.clearedhistory
}

// RemoveHistoryIDs removes the "history" edge to the CookingHistory entity by IDs.
func (m *RecipeMutation) RemoveHistoryIDs(ids ...uuid.UUID) {
	if m.removedhistory == nil {
		m.removedhistory = make(map[uuid.UUID]struct{})
	}
	for i := range ids {
		delete(m.history, ids[i])
		m.removedhistory[ids[i]] = struct{}{}
	}
}

// RemovedHistory returns the removed IDs of the "history" edge to the CookingHistory entity.
func (m *RecipeMutation) RemovedHistoryIDs() (ids []uuid.UUID) {
	for id := range m.removedhistory {
		ids = append(ids, id)
	}
	return
}

// HistoryIDs returns the "history" edge IDs in the mutation.
func (m *RecipeMutation) HistoryIDs() (ids []uuid.UUID) {
	for id := range m.history {
		ids = append(ids, id)
	}
	return
}

// ResetHistory resets all changes to the "history" edge.
func (m *RecipeMutation) ResetHistory() {
	m.history = nil
	m.clearedhistory = false
	m.removedhistory = nil
}

// Where appends a list predicates to the RecipeMutation builder.
func (m *RecipeMutation) Where(ps ...predicate.Recipe) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *RecipeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Recipe).
func (m *RecipeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *RecipeMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.title != nil {
		fields = append(fields, recipe.FieldTitle)
	}
	if m.planned != nil {
		fields = append(fields, recipe.FieldPlanned)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *RecipeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case recipe.FieldTitle:
		return m.Title()
	case recipe.FieldPlanned:
		return m.Planned()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *RecipeMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case recipe.FieldTitle:
		return m.OldTitle(ctx)
	case recipe.FieldPlanned:
		return m.OldPlanned(ctx)
	}
	return nil, fmt.Errorf("unknown Recipe field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RecipeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case recipe.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	case recipe.FieldPlanned:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPlanned(v)
		return nil
	}
	return fmt.Errorf("unknown Recipe field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *RecipeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *RecipeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RecipeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Recipe numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *RecipeMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *RecipeMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *RecipeMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Recipe nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *RecipeMutation) ResetField(name string) error {
	switch name {
	case recipe.FieldTitle:
		m.ResetTitle()
		return nil
	case recipe.FieldPlanned:
		m.ResetPlanned()
		return nil
	}
	return fmt.Errorf("unknown Recipe field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *RecipeMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.history != nil {
		edges = append(edges, recipe.EdgeHistory)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *RecipeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case recipe.EdgeHistory:
		ids := make([]ent.Value, 0, len(m.history))
		for id := range m.history {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *RecipeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedhistory != nil {
		edges = append(edges, recipe.EdgeHistory)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *RecipeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case recipe.EdgeHistory:
		ids := make([]ent.Value, 0, len(m.removedhistory))
		for id := range m.removedhistory {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *RecipeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedhistory {
		edges = append(edges, recipe.EdgeHistory)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *RecipeMutation) EdgeCleared(name string) bool {
	switch name {
	case recipe.EdgeHistory:
		return m.clearedhistory
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *RecipeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Recipe unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *RecipeMutation) ResetEdge(name string) error {
	switch name {
	case recipe.EdgeHistory:
		m.ResetHistory()
		return nil
	}
	return fmt.Errorf("unknown Recipe edge %s", name)
}
