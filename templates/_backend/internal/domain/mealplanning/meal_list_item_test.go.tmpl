package mealplanning

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMealListItemUpdate(t *testing.T) {
	t.Parallel()

	notes := "note2"
	mli := &MealListItem{Notes: "note1"}
	mli.Update(&MealListItemUpdateRequestInput{Notes: &notes})
	require.Equal(t, notes, mli.Notes)

	mli.Update(nil)
	require.Equal(t, notes, mli.Notes)
}

func TestMealListItemValidation(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	require.NoError(t, (&MealListItemCreationRequestInput{MealID: "meal"}).ValidateWithContext(ctx))
	require.Error(t, (&MealListItemCreationRequestInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&MealListItemDatabaseCreationInput{
		ID:                "id",
		MealID:            "meal",
		BelongsToMealList: "list",
	}).ValidateWithContext(ctx))
	require.Error(t, (&MealListItemDatabaseCreationInput{}).ValidateWithContext(ctx))

	require.NoError(t, (&MealListItemUpdateRequestInput{}).ValidateWithContext(ctx))
}
