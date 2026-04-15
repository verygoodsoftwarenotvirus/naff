package mealplanning

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestMealCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Components: []*MealComponentCreationRequestInput{
				{
					RecipeID:      t.Name(),
					ComponentType: MealComponentTypesMain,
				},
			},
		}

		assert.NoError(t, x.ValidateWithContext(t.Context()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{}

		assert.Error(t, x.ValidateWithContext(t.Context()))
	})

	T.Run("with invalid component", func(t *testing.T) {
		t.Parallel()

		x := &MealCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Components: []*MealComponentCreationRequestInput{
				{},
			},
		}

		assert.Error(t, x.ValidateWithContext(t.Context()))
	})
}

func TestMealUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealUpdateRequestInput{
			Name:          new(t.Name()),
			Description:   new(t.Name()),
			CreatedByUser: new(t.Name()),
			Components: []*MealComponentUpdateRequestInput{
				{
					RecipeID:      new(t.Name()),
					RecipeScale:   new(float32(exampleQuantity)),
					ComponentType: new(MealComponentTypesAmuseBouche),
				},
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestMealComponentCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealComponentCreationRequestInput{
			RecipeID:      t.Name(),
			RecipeScale:   exampleQuantity,
			ComponentType: MealComponentTypesAmuseBouche,
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})
}

func TestMealDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealDatabaseCreationInput{
			Name: t.Name(),
			Components: []*MealComponentDatabaseCreationInput{
				{
					RecipeID: t.Name(),
				},
			},
			CreatedByUser: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealUpdateRequestInput{
			Name:        new(t.Name()),
			Description: new(t.Name()),
			Components: []*MealComponentUpdateRequestInput{
				{
					RecipeID: new(t.Name()),
				},
			},
			CreatedByUser: new(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealComponent_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealComponent{}
		input := &MealComponentUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestMeal_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Meal{}
		input := &MealUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.EligibleForMealPlans = new(true)

		x.Update(input)
	})
}
