package mealplanning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecipeListUpdate(t *testing.T) {
	t.Parallel()

	name := "new"
	desc := "new-desc"

	rl := &RecipeList{Name: "old", Description: "old-desc"}
	rl.Update(&RecipeListUpdateRequestInput{
		Name:        &name,
		Description: &desc,
	})

	require.Equal(t, name, rl.Name)
	require.Equal(t, desc, rl.Description)

	rl.Update(nil)
	require.Equal(t, name, rl.Name)
	require.Equal(t, desc, rl.Description)
}

func TestRecipeListValidation(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	require.NoError(t, (&RecipeListCreationRequestInput{Name: "name"}).ValidateWithContext(ctx))
	require.Error(t, (&RecipeListCreationRequestInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&RecipeListDatabaseCreationInput{
		ID:            "id",
		Name:          "name",
		Description:   "desc",
		BelongsToUser: "user",
	}).ValidateWithContext(ctx))
	require.Error(t, (&RecipeListDatabaseCreationInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&RecipeListUpdateRequestInput{Name: new("x")}).ValidateWithContext(ctx))
}
