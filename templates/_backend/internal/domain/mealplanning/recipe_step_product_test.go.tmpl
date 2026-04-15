package mealplanning

import (
	"testing"

	"github.com/primandproper/platform/numbers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProduct{
			MeasurementUnit:        &ValidMeasurementUnit{},
			ContainedInVesselIndex: new(uint16(3)),
		}
		input := &RecipeStepProductUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Compostable = new(true)
		input.ContainedInVesselIndex = new(uint16(1))
		input.IsLiquid = new(true)
		input.IsWaste = new(true)

		x.Update(input)
	})
}

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:              t.Name(),
			Type:              RecipeStepProductIngredientType,
			MeasurementUnitID: new(t.Name()),
			MeasurementQuantity: numbers.OpenRange[float32]{
				Max: nil,
				Min: new(fake.Float32()),
			},
			QuantityNotes:            t.Name(),
			Compostable:              fake.Bool(),
			StorageDurationInSeconds: numbers.OpenRange[uint32]{Max: new(fake.Uint32())},
			StorageTemperatureInCelsius: numbers.OpenRange[float32]{
				Max: new(fake.Float32()),
				Min: new(fake.Float32()),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{
			Name:                        new(t.Name()),
			Type:                        new(RecipeStepProductIngredientType),
			MeasurementUnitID:           new(t.Name()),
			MeasurementQuantity:         numbers.OpenRange[float32]{Max: new(fake.Float32()), Min: new(fake.Float32())},
			QuantityNotes:               new(t.Name()),
			Compostable:                 new(fake.Bool()),
			StorageTemperatureInCelsius: numbers.OpenRange[float32]{Max: new(fake.Float32()), Min: new(fake.Float32())},
			StorageDurationInSeconds:    numbers.OpenRange[uint32]{Max: new(fake.Uint32()), Min: new(fake.Uint32())},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
