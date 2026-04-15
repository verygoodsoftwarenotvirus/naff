package mealplanning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecipeListItemUpdate(t *testing.T) {
	t.Parallel()

	notes := "note2"
	rli := &RecipeListItem{Notes: "note1"}
	rli.Update(&RecipeListItemUpdateRequestInput{Notes: &notes})
	require.Equal(t, notes, rli.Notes)

	rli.Update(nil)
	require.Equal(t, notes, rli.Notes)
}

func TestRecipeListItemValidation(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	require.NoError(t, (&RecipeListItemCreationRequestInput{RecipeID: "recipe"}).ValidateWithContext(ctx))
	require.Error(t, (&RecipeListItemCreationRequestInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&RecipeListItemDatabaseCreationInput{
		ID:                  "id",
		RecipeID:            "recipe",
		BelongsToRecipeList: "list",
	}).ValidateWithContext(ctx))
	require.Error(t, (&RecipeListItemDatabaseCreationInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&RecipeListItemUpdateRequestInput{}).ValidateWithContext(ctx))
}
