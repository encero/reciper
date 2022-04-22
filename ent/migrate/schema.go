// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CookingHistoriesColumns holds the columns for the "cooking_histories" table.
	CookingHistoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "cooked_at", Type: field.TypeTime},
		{Name: "recipe_history", Type: field.TypeUUID},
	}
	// CookingHistoriesTable holds the schema information for the "cooking_histories" table.
	CookingHistoriesTable = &schema.Table{
		Name:       "cooking_histories",
		Columns:    CookingHistoriesColumns,
		PrimaryKey: []*schema.Column{CookingHistoriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "cooking_histories_recipes_history",
				Columns:    []*schema.Column{CookingHistoriesColumns[2]},
				RefColumns: []*schema.Column{RecipesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// RecipesColumns holds the columns for the "recipes" table.
	RecipesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "title", Type: field.TypeString},
		{Name: "planned", Type: field.TypeBool, Default: false},
	}
	// RecipesTable holds the schema information for the "recipes" table.
	RecipesTable = &schema.Table{
		Name:       "recipes",
		Columns:    RecipesColumns,
		PrimaryKey: []*schema.Column{RecipesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CookingHistoriesTable,
		RecipesTable,
	}
)

func init() {
	CookingHistoriesTable.ForeignKeys[0].RefTable = RecipesTable
}
