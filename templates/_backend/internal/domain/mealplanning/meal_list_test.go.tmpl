package mealplanning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMealListUpdate(t *testing.T) {
	t.Parallel()

	name := "new"
	desc := "new-desc"

	ml := &MealList{Name: "old", Description: "old-desc"}
	ml.Update(&MealListUpdateRequestInput{
		Name:        &name,
		Description: &desc,
	})

	require.Equal(t, name, ml.Name)
	require.Equal(t, desc, ml.Description)

	// nil input should be a no-op
	ml.Update(nil)
	require.Equal(t, name, ml.Name)
	require.Equal(t, desc, ml.Description)
}

func TestMealListValidation(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	require.NoError(t, (&MealListCreationRequestInput{Name: "name"}).ValidateWithContext(ctx))
	require.Error(t, (&MealListCreationRequestInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&MealListDatabaseCreationInput{
		ID:            "id",
		Name:          "name",
		Description:   "desc",
		BelongsToUser: "user",
	}).ValidateWithContext(ctx))
	require.Error(t, (&MealListDatabaseCreationInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&MealListUpdateRequestInput{Name: new("x")}).ValidateWithContext(ctx))
}
