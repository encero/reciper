// Code generated by entc, DO NOT EDIT.

package cookinghistory

const (
	// Label holds the string label denoting the cookinghistory type in the database.
	Label = "cooking_history"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCookedAt holds the string denoting the cookedat field in the database.
	FieldCookedAt = "cooked_at"
	// EdgeRecipe holds the string denoting the recipe edge name in mutations.
	EdgeRecipe = "recipe"
	// Table holds the table name of the cookinghistory in the database.
	Table = "cooking_histories"
	// RecipeTable is the table that holds the recipe relation/edge.
	RecipeTable = "cooking_histories"
	// RecipeInverseTable is the table name for the Recipe entity.
	// It exists in this package in order to avoid circular dependency with the "recipe" package.
	RecipeInverseTable = "recipes"
	// RecipeColumn is the table column denoting the recipe relation/edge.
	RecipeColumn = "recipe_history"
)

// Columns holds all SQL columns for cookinghistory fields.
var Columns = []string{
	FieldID,
	FieldCookedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "cooking_histories"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"recipe_history",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}
