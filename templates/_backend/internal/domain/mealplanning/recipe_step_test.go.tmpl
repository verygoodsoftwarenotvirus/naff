package mealplanning

import (
	"testing"

	"github.com/primandproper/platform/numbers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func buildValidRecipeStepCreationRequestInput() *RecipeStepCreationRequestInput {
	return &RecipeStepCreationRequestInput{
		PreparationID: "PreparationID",
		EstimatedTimeInSeconds: numbers.OpenRange[uint32]{
			Max: new(fake.Uint32()),
			Min: new(fake.Uint32()),
		},
		TemperatureInCelsius: numbers.OpenRange[float32]{
			Max: nil,
			Min: new(float32(123.45)),
		},
		Notes:                "Notes",
		ExplicitInstructions: "ExplicitInstructions",
		Instruments: []*RecipeStepInstrumentCreationRequestInput{
			{
				ValidPreparationInstrumentID: new("ValidPreparationInstrumentID"),
				Name:                         "Name",
				Quantity:                     numbers.MinRange[uint32]{Min: fake.Uint32()},
			},
		},
		Products: []*RecipeStepProductCreationRequestInput{
			{
				Name:                "Name",
				Type:                RecipeStepProductIngredientType,
				MeasurementUnitID:   new("MeasurementUnitID"),
				MeasurementQuantity: numbers.OpenRange[float32]{Min: new(float32(1))},
				QuantityNotes:       "QuantityNotes",
			},
		},
		Ingredients: []*RecipeStepIngredientCreationRequestInput{
			{
				ValidIngredientPreparationID:     new("ValidIngredientPreparationID"),
				ValidIngredientMeasurementUnitID: new("ValidIngredientMeasurementUnitID"),
				QuantityNotes:                    "QuantityNotes",
				IngredientNotes:                  "IngredientNotes",
				Quantity:                         numbers.MinRange[float32]{Min: 1},
			},
		},
	}
}

func TestRecipeStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStep{
			TemperatureInCelsius: numbers.OpenRange[float32]{Min: new(float32(123.45))},
		}

		input := &RecipeStepUpdateRequestInput{}
		assert.NoError(t, fake.Struct(&input))
		input.Optional = new(true)
		input.StartTimerAutomatically = new(true)
		input.TemperatureInCelsius.Min = new(float32(543.21))
		input.TemperatureInCelsius.Max = new(float32(123.45))

		x.Update(input)
	})
}

func TestRecipeStepCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := buildValidRecipeStepCreationRequestInput()

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})

	T.Run("with too many ingredients", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationRequestInput{
			PreparationID: t.Name(),
			EstimatedTimeInSeconds: numbers.OpenRange[uint32]{
				Max: new(fake.Uint32()),
				Min: new(fake.Uint32()),
			},
			TemperatureInCelsius: numbers.OpenRange[float32]{
				Max: nil,
				Min: new(float32(123.45)),
			},
			Notes:                t.Name(),
			ExplicitInstructions: t.Name(),
			Products: []*RecipeStepProductCreationRequestInput{
				{
					Name: t.Name(),
				},
			},
			Ingredients: []*RecipeStepIngredientCreationRequestInput{},
		}

		for range maxIngredientsPerStep * 2 {
			x.Ingredients = append(x.Ingredients, &RecipeStepIngredientCreationRequestInput{
				ValidIngredientPreparationID:     new(t.Name()),
				ValidIngredientMeasurementUnitID: new(t.Name()),
				QuantityNotes:                    t.Name(),
				IngredientNotes:                  t.Name(),
				Quantity:                         numbers.MinRange[float32]{Min: 1},
			})
		}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepDatabaseCreationInput{
			ID:            t.Name(),
			PreparationID: t.Name(),
			Products: []*RecipeStepProductDatabaseCreationInput{
				{
					ID:                  t.Name(),
					Name:                t.Name(),
					Type:                RecipeStepProductIngredientType,
					BelongsToRecipeStep: t.Name(),
				},
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepDatabaseCreationInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{
			Index:       new(fake.Uint32()),
			Preparation: &ValidPreparation{},
			EstimatedTimeInSeconds: numbers.OpenRange[uint32]{
				Max: new(fake.Uint32()),
				Min: new(fake.Uint32()),
			},
			TemperatureInCelsius: numbers.OpenRange[float32]{
				Max: nil,
				Min: new(float32(123.45)),
			},
			Notes:                new(t.Name()),
			ExplicitInstructions: new(t.Name()),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
