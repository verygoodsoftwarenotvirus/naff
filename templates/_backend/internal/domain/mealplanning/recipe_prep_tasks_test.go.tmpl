package mealplanning

import (
	"testing"

	"github.com/primandproper/platform/numbers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipePrepTask_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTask{}
		input := &RecipePrepTaskUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Optional = new(true)
		input.BelongsToRecipe = new(t.Name())

		x.Update(input)
	})
}

func TestRecipePrepTaskCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskCreationRequestInput{
			BelongsToRecipe: t.Name(),
			Name:            t.Name(),
			StorageType:     t.Name(),
			StorageTemperatureInCelsius: numbers.OpenRange[float32]{
				Max: new(fake.Float32()),
				Min: new(fake.Float32()),
			},
			TimeBufferBeforeRecipeInSeconds: numbers.MinRange[uint32]{
				Max: new(fake.Uint32()),
				Min: fake.Uint32(),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipePrepTaskDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskDatabaseCreationInput{
			ID:              t.Name(),
			BelongsToRecipe: t.Name(),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskDatabaseCreationInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipePrepTaskUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskUpdateRequestInput{
			BelongsToRecipe: new(t.Name()),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
