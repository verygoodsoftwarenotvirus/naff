package mealplanning

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/primandproper/platform/database/filtering"
	"github.com/primandproper/platform/numbers"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

var (
	errValidIngredientPreparationIDRequired     = errors.New("validIngredientPreparationID is required when not referencing a recipe step product")
	errValidIngredientMeasurementUnitIDRequired = errors.New("validIngredientMeasurementUnitID is required when not referencing a recipe step product")
)

const (
	// RecipeStepIngredientCreatedServiceEventType indicates a recipe step ingredient was created.
	RecipeStepIngredientCreatedServiceEventType = "recipe_step_ingredient_created"
	// RecipeStepIngredientUpdatedServiceEventType indicates a recipe step ingredient was updated.
	RecipeStepIngredientUpdatedServiceEventType = "recipe_step_ingredient_updated"
	// RecipeStepIngredientArchivedServiceEventType indicates a recipe step ingredient was archived.
	RecipeStepIngredientArchivedServiceEventType = "recipe_step_ingredient_archived"
)

func init() {
	gob.Register(new(RecipeStepIngredient))
	gob.Register(new(RecipeStepIngredientCreationRequestInput))
	gob.Register(new(RecipeStepIngredientUpdateRequestInput))
}

type (
	// RecipeStepIngredient represents a recipe step ingredient.
	RecipeStepIngredient struct {
		_ struct{} `json:"-"`

		CreatedAt                 time.Time                 `json:"createdAt"`
		RecipeStepProductRecipeID *string                   `json:"productOfRecipeID"`
		ArchivedAt                *time.Time                `json:"archivedAt"`
		Ingredient                *ValidIngredient          `json:"ingredient"`
		LastUpdatedAt             *time.Time                `json:"lastUpdatedAt"`
		VesselIndex               *uint16                   `json:"vesselIndex"`
		ProductPercentageToUse    *float32                  `json:"productPercentageToUse"`
		RecipeStepProductID       *string                   `json:"recipeStepProductID"`
		BelongsToRecipeStep       string                    `json:"belongsToRecipeStep"`
		ID                        string                    `json:"id"`
		QuantityNotes             string                    `json:"quantityNotes"`
		IngredientNotes           string                    `json:"ingredientNotes"`
		Name                      string                    `json:"name"`
		Quantity                  numbers.MinRange[float32] `json:"quantity"`
		MeasurementUnit           ValidMeasurementUnit      `json:"measurementUnit"`
		Index                     uint16                    `json:"index"`
		OptionIndex               uint16                    `json:"optionIndex"`
		Optional                  bool                      `json:"optional"`
		ToTaste                   bool                      `json:"toTaste"`
		ScaleFactor               float32                   `json:"scaleFactor"`
	}

	// RecipeStepIngredientCreationRequestInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientCreationRequestInput struct {
		_ struct{} `json:"-"`

		Quantity                         numbers.MinRange[float32] `json:"quantity"`
		ValidIngredientMeasurementUnitID *string                   `json:"validIngredientMeasurementUnitID"`
		Index                            *uint16                   `json:"index,omitempty"`
		VesselIndex                      *uint16                   `json:"vesselIndex"`
		ProductPercentageToUse           *float32                  `json:"productPercentageToUse"`
		ProductOfRecipeStepIndex         *uint64                   `json:"productOfRecipeStepIndex"`
		ValidIngredientPreparationID     *string                   `json:"validIngredientPreparationID"`
		RecipeStepProductRecipeID        *string                   `json:"productOfRecipeID"`
		RecipeStepProductRecipeSlug      *string                   `json:"-"` // Used by resolveEmptyRecipeIDs when ID is empty; not persisted
		ProductOfRecipeStepProductIndex  *uint64                   `json:"productOfRecipeStepProductIndex"`
		Name                             string                    `json:"name"`
		QuantityNotes                    string                    `json:"quantityNotes"`
		IngredientNotes                  string                    `json:"ingredientNotes"`
		OptionIndex                      uint16                    `json:"optionIndex"`
		Optional                         bool                      `json:"optional"`
		ToTaste                          bool                      `json:"toTaste"`
		ScaleFactor                      float32                   `json:"scaleFactor"`
	}

	// RecipeStepIngredientDatabaseCreationInput represents what a user could set as input for creating recipe step ingredients.
	RecipeStepIngredientDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		Quantity                         numbers.MinRange[float32] `json:"-"`
		RecipeStepProductRecipeID        *string                   `json:"-"`
		IngredientID                     *string                   `json:"-"`
		RecipeStepProductID              *string                   `json:"-"`
		ProductOfRecipeStepIndex         *uint64                   `json:"-"`
		ProductOfRecipeStepProductIndex  *uint64                   `json:"-"`
		VesselIndex                      *uint16                   `json:"-"`
		ProductPercentageToUse           *float32                  `json:"-"`
		ValidIngredientPreparationID     *string                   `json:"-"`
		ValidIngredientMeasurementUnitID *string                   `json:"-"`
		ID                               string                    `json:"-"`
		BelongsToRecipeStep              string                    `json:"-"`
		Name                             string                    `json:"-"`
		IngredientNotes                  string                    `json:"-"`
		QuantityNotes                    string                    `json:"-"`
		MeasurementUnitID                string                    `json:"-"`
		Index                            uint16                    `json:"-"`
		OptionIndex                      uint16                    `json:"-"`
		Optional                         bool                      `json:"-"`
		ToTaste                          bool                      `json:"-"`
		ScaleFactor                      float32                   `json:"-"`
	}

	// RecipeStepIngredientUpdateRequestInput represents what a user could set as input for updating recipe step ingredients.
	RecipeStepIngredientUpdateRequestInput struct {
		_ struct{} `json:"-"`

		IngredientID              *string                                      `json:"ingredientID,omitempty"`
		RecipeStepProductID       *string                                      `json:"recipeStepProductID,omitempty"`
		Name                      *string                                      `json:"name,omitempty"`
		Optional                  *bool                                        `json:"optional,omitempty"`
		MeasurementUnitID         *string                                      `json:"measurementUnitID,omitempty"`
		QuantityNotes             *string                                      `json:"quantityNotes,omitempty"`
		IngredientNotes           *string                                      `json:"ingredientNotes,omitempty"`
		BelongsToRecipeStep       *string                                      `json:"belongsToRecipeStep,omitempty"`
		Quantity                  numbers.OpenRangeUpdateRequestInput[float32] `json:"quantity"`
		Index                     *uint16                                      `json:"index,omitempty"`
		OptionIndex               *uint16                                      `json:"optionIndex,omitempty"`
		VesselIndex               *uint16                                      `json:"vesselIndex,omitempty"`
		ToTaste                   *bool                                        `json:"toTaste,omitempty"`
		ProductPercentageToUse    *float32                                     `json:"productPercentageToUse,omitempty"`
		RecipeStepProductRecipeID *string                                      `json:"productOfRecipeID"`
		ScaleFactor               *float32                                     `json:"scaleFactor,omitempty"`
	}

	// RecipeStepIngredientDataManager describes a structure capable of storing recipe step ingredients permanently.
	RecipeStepIngredientDataManager interface {
		RecipeStepIngredientExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error)
		GetRecipeStepIngredient(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*RecipeStepIngredient, error)
		GetRecipeStepIngredients(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[RecipeStepIngredient], error)
		CreateRecipeStepIngredient(ctx context.Context, input *RecipeStepIngredientDatabaseCreationInput) (*RecipeStepIngredient, error)
		UpdateRecipeStepIngredient(ctx context.Context, updated *RecipeStepIngredient) error
		ArchiveRecipeStepIngredient(ctx context.Context, recipeStepID, recipeStepIngredientID string) error
	}

	// RecipeStepIngredientDataService describes a structure capable of serving traffic related to recipe step ingredients.
	RecipeStepIngredientDataService interface {
		ListRecipeStepIngredientsHandler(http.ResponseWriter, *http.Request)
		CreateRecipeStepIngredientHandler(http.ResponseWriter, *http.Request)
		ReadRecipeStepIngredientHandler(http.ResponseWriter, *http.Request)
		UpdateRecipeStepIngredientHandler(http.ResponseWriter, *http.Request)
		ArchiveRecipeStepIngredientHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepIngredientUpdateRequestInput with a recipe step ingredient.
func (x *RecipeStepIngredient) Update(input *RecipeStepIngredientUpdateRequestInput) {
	if input.IngredientID != nil && (x.Ingredient == nil || (*input.IngredientID != "" && *input.IngredientID != x.Ingredient.ID)) {
		x.Ingredient = &ValidIngredient{ID: *input.IngredientID}
	}

	if input.RecipeStepProductID != nil && (x.RecipeStepProductID == nil || (*input.RecipeStepProductID != "" && *input.RecipeStepProductID != *x.RecipeStepProductID)) {
		x.RecipeStepProductID = input.RecipeStepProductID
	}

	if input.RecipeStepProductRecipeID != nil && ((x.RecipeStepProductRecipeID != nil && *input.RecipeStepProductRecipeID != *x.RecipeStepProductRecipeID) || x.RecipeStepProductRecipeID == nil) {
		x.RecipeStepProductRecipeID = input.RecipeStepProductRecipeID
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.MeasurementUnitID != nil && *input.MeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit = ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	if input.Quantity.Min != nil && *input.Quantity.Min != x.Quantity.Min {
		x.Quantity.Min = *input.Quantity.Min
	}

	if input.Quantity.Max != nil && ((x.Quantity.Max != nil && *input.Quantity.Max != *x.Quantity.Max) || x.Quantity.Max == nil) {
		x.Quantity.Max = input.Quantity.Max
	}

	if input.QuantityNotes != nil && *input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = *input.QuantityNotes
	}

	if input.IngredientNotes != nil && *input.IngredientNotes != x.IngredientNotes {
		x.IngredientNotes = *input.IngredientNotes
	}

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}

	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
	}

	if input.OptionIndex != nil && *input.OptionIndex != x.OptionIndex {
		x.OptionIndex = *input.OptionIndex
	}

	if input.VesselIndex != nil && ((x.VesselIndex != nil && *input.VesselIndex != *x.VesselIndex) || x.VesselIndex == nil) {
		x.VesselIndex = input.VesselIndex
	}

	if input.ToTaste != nil && *input.ToTaste != x.ToTaste {
		x.ToTaste = *input.ToTaste
	}

	if input.ProductPercentageToUse != nil && ((x.ProductPercentageToUse != nil && *input.ProductPercentageToUse != *x.ProductPercentageToUse) || x.ProductPercentageToUse == nil) {
		x.ProductPercentageToUse = input.ProductPercentageToUse
	}

	if input.ScaleFactor != nil && *input.ScaleFactor != x.ScaleFactor {
		x.ScaleFactor = *input.ScaleFactor
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientCreationRequestInput.
func (x *RecipeStepIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	err := &multierror.Error{}

	// When not referencing a recipe step product, bridge table IDs are required
	if isRecipeStepProduct := x.ProductOfRecipeStepIndex != nil; !isRecipeStepProduct {
		if x.ValidIngredientPreparationID == nil || *x.ValidIngredientPreparationID == "" {
			err = multierror.Append(err, errValidIngredientPreparationIDRequired)
		}
		if x.ValidIngredientMeasurementUnitID == nil || *x.ValidIngredientMeasurementUnitID == "" {
			err = multierror.Append(err, errValidIngredientMeasurementUnitIDRequired)
		}
	}

	validationErr := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Quantity, validation.Required),
	)
	if validationErr != nil {
		err = multierror.Append(err, validationErr)
	}

	return err.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientDatabaseCreationInput.
func (x *RecipeStepIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepIngredientUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepIngredientUpdateRequestInput.
func (x *RecipeStepIngredientUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
	)
}
