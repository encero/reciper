// Code generated by entc, DO NOT EDIT.

package recipe

const (
	// Label holds the string label denoting the recipe type in the database.
	Label = "recipe"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldPlanned holds the string denoting the planned field in the database.
	FieldPlanned = "planned"
	// Table holds the table name of the recipe in the database.
	Table = "recipes"
)

// Columns holds all SQL columns for recipe fields.
var Columns = []string{
	FieldID,
	FieldTitle,
	FieldPlanned,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultPlanned holds the default value on creation for the "planned" field.
	DefaultPlanned bool
)
