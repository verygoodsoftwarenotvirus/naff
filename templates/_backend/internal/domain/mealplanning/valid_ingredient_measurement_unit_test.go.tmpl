package mealplanning

import (
	"testing"

	"github.com/primandproper/platform/numbers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientMeasurementUnit_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnit{
			AllowableQuantity: numbers.MinRange[float32]{Max: new(float32(3.21))},
		}
		input := &ValidIngredientMeasurementUnitUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.AllowableQuantity.Max = new(float32(1.23))

		x.Update(input)
	})
}

func TestValidIngredientMeasurementUnitCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{
			Notes:                  t.Name(),
			ValidMeasurementUnitID: t.Name(),
			ValidIngredientID:      t.Name(),
			AllowableQuantity: numbers.MinRange[float32]{
				Max: new(fake.Float32()),
				Min: fake.Float32(),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestValidIngredientMeasurementUnitDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitDatabaseCreationInput{
			ID:                     t.Name(),
			ValidMeasurementUnitID: t.Name(),
			ValidIngredientID:      t.Name(),
			AllowableQuantity: numbers.MinRange[float32]{
				Min: fake.Float32(),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitDatabaseCreationInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestValidIngredientMeasurementUnitUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitUpdateRequestInput{
			Notes:                  new(t.Name()),
			ValidMeasurementUnitID: new(t.Name()),
			ValidIngredientID:      new(t.Name()),
			AllowableQuantity: numbers.OpenRangeUpdateRequestInput[float32]{
				Min: new(fake.Float32()),
				Max: new(fake.Float32()),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
