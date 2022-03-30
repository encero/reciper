// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/encero/reciper-api/ent/recipe"
	"github.com/encero/reciper-api/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	recipeFields := schema.Recipe{}.Fields()
	_ = recipeFields
	// recipeDescPlanned is the schema descriptor for planned field.
	recipeDescPlanned := recipeFields[2].Descriptor()
	// recipe.DefaultPlanned holds the default value on creation for the planned field.
	recipe.DefaultPlanned = recipeDescPlanned.Default.(bool)
}
