package mealplanning

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/primandproper/platform/database/filtering"
	"github.com/primandproper/platform/numbers"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepProductIngredientType represents one of the valid recipe step product type values.
	RecipeStepProductIngredientType = "ingredient"
	// RecipeStepProductInstrumentType represents one of the valid recipe step product type values.
	RecipeStepProductInstrumentType = "instrument"
	// RecipeStepProductVesselType represents one of the valid recipe step product type values.
	RecipeStepProductVesselType = "vessel"

	// RecipeStepProductCreatedServiceEventType indicates a recipe step product was created.
	RecipeStepProductCreatedServiceEventType = "recipe_step_product_created"
	// RecipeStepProductUpdatedServiceEventType indicates a recipe step product was updated.
	RecipeStepProductUpdatedServiceEventType = "recipe_step_product_updated"
	// RecipeStepProductArchivedServiceEventType indicates a recipe step product was archived.
	RecipeStepProductArchivedServiceEventType = "recipe_step_product_archived"
)

func init() {
	gob.Register(new(RecipeStepProduct))
	gob.Register(new(RecipeStepProductCreationRequestInput))
	gob.Register(new(RecipeStepProductUpdateRequestInput))
}

type (
	// RecipeStepProduct represents a recipe step product.
	RecipeStepProduct struct {
		_ struct{} `json:"-"`

		CreatedAt                   time.Time                  `json:"createdAt"`
		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"storageTemperatureInCelsius"`
		StorageDurationInSeconds    numbers.OpenRange[uint32]  `json:"storageDurationInSeconds"`
		MeasurementQuantity         numbers.OpenRange[float32] `json:"measurementQuantity"`
		ItemQuantity                numbers.OpenRange[float32] `json:"itemQuantity"`
		ArchivedAt                  *time.Time                 `json:"archivedAt"`
		LastUpdatedAt               *time.Time                 `json:"lastUpdatedAt"`
		MeasurementUnit             *ValidMeasurementUnit      `json:"measurementUnit"`
		ContainedInVesselIndex      *uint16                    `json:"containedInVesselIndex"`
		Name                        string                     `json:"name"`
		BelongsToRecipeStep         string                     `json:"belongsToRecipeStep"`
		Type                        string                     `json:"type"`
		ID                          string                     `json:"id"`
		StorageInstructions         string                     `json:"storageInstructions"`
		QuantityNotes               string                     `json:"quantityNotes"`
		Index                       uint16                     `json:"index"`
		IsWaste                     bool                       `json:"isWaste"`
		IsLiquid                    bool                       `json:"isLiquid"`
		Compostable                 bool                       `json:"compostable"`
	}

	// RecipeStepProductCreationRequestInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductCreationRequestInput struct {
		_ struct{} `json:"-"`

		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"storageTemperatureInCelsius"`
		StorageDurationInSeconds    numbers.OpenRange[uint32]  `json:"storageDurationInSeconds"`
		MeasurementQuantity         numbers.OpenRange[float32] `json:"measurementQuantity"`
		ItemQuantity                numbers.OpenRange[float32] `json:"itemQuantity"`
		MeasurementUnitID           *string                    `json:"measurementUnitID"`
		ContainedInVesselIndex      *uint16                    `json:"containedInVesselIndex"`
		QuantityNotes               string                     `json:"quantityNotes"`
		Name                        string                     `json:"name"`
		StorageInstructions         string                     `json:"storageInstructions"`
		Type                        string                     `json:"type"`
		Index                       uint16                     `json:"index"`
		Compostable                 bool                       `json:"compostable"`
		IsLiquid                    bool                       `json:"isLiquid"`
		IsWaste                     bool                       `json:"isWaste"`
	}

	// RecipeStepProductDatabaseCreationInput represents what a user could set as input for creating recipe step products.
	RecipeStepProductDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"-"`
		StorageDurationInSeconds    numbers.OpenRange[uint32]  `json:"-"`
		MeasurementQuantity         numbers.OpenRange[float32] `json:"-"`
		ItemQuantity                numbers.OpenRange[float32] `json:"-"`
		MeasurementUnitID           *string                    `json:"-"`
		ContainedInVesselIndex      *uint16                    `json:"-"`
		Name                        string                     `json:"-"`
		BelongsToRecipeStep         string                     `json:"-"`
		StorageInstructions         string                     `json:"-"`
		QuantityNotes               string                     `json:"-"`
		ID                          string                     `json:"-"`
		Type                        string                     `json:"-"`
		Index                       uint16                     `json:"-"`
		Compostable                 bool                       `json:"-"`
		IsLiquid                    bool                       `json:"-"`
		IsWaste                     bool                       `json:"-"`
	}

	// RecipeStepProductUpdateRequestInput represents what a user could set as input for updating recipe step products.
	RecipeStepProductUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Name                        *string                    `json:"name,omitempty"`
		Type                        *string                    `json:"type,omitempty"`
		MeasurementUnitID           *string                    `json:"measurementUnitID,omitempty"`
		QuantityNotes               *string                    `json:"quantityNotes,omitempty"`
		BelongsToRecipeStep         *string                    `json:"belongsToRecipeStep,omitempty"`
		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"storageTemperatureInCelsius"`
		StorageDurationInSeconds    numbers.OpenRange[uint32]  `json:"storageDurationInSeconds"`
		MeasurementQuantity         numbers.OpenRange[float32] `json:"measurementQuantity"`
		ItemQuantity                numbers.OpenRange[float32] `json:"itemQuantity"`
		Compostable                 *bool                      `json:"compostable,omitempty"`
		StorageInstructions         *string                    `json:"storageInstructions,omitempty"`
		IsLiquid                    *bool                      `json:"isLiquid,omitempty"`
		IsWaste                     *bool                      `json:"isWaste,omitempty"`
		Index                       *uint16                    `json:"index,omitempty"`
		ContainedInVesselIndex      *uint16                    `json:"containedInVesselIndex,omitempty"`
	}

	// RecipeStepProductDataManager describes a structure capable of storing recipe step products permanently.
	RecipeStepProductDataManager interface {
		RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (bool, error)
		GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*RecipeStepProduct, error)
		GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[RecipeStepProduct], error)
		CreateRecipeStepProduct(ctx context.Context, input *RecipeStepProductDatabaseCreationInput) (*RecipeStepProduct, error)
		UpdateRecipeStepProduct(ctx context.Context, updated *RecipeStepProduct) error
		ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error
	}

	// RecipeStepProductDataService describes a structure capable of serving traffic related to recipe step products.
	RecipeStepProductDataService interface {
		ListRecipeStepProductsHandler(http.ResponseWriter, *http.Request)
		CreateRecipeStepProductHandler(http.ResponseWriter, *http.Request)
		ReadRecipeStepProductHandler(http.ResponseWriter, *http.Request)
		UpdateRecipeStepProductHandler(http.ResponseWriter, *http.Request)
		ArchiveRecipeStepProductHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepProductUpdateRequestInput with a recipe step product.
func (x *RecipeStepProduct) Update(input *RecipeStepProductUpdateRequestInput) {
	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Type != nil && *input.Type != x.Type {
		x.Type = *input.Type
	}

	if input.MeasurementUnitID != nil && x.MeasurementUnit != nil && *input.MeasurementUnitID != x.MeasurementUnit.ID {
		x.MeasurementUnit = &ValidMeasurementUnit{ID: *input.MeasurementUnitID}
	}

	if input.MeasurementQuantity.Min != nil && input.MeasurementQuantity.Min != x.MeasurementQuantity.Min {
		x.MeasurementQuantity.Min = input.MeasurementQuantity.Min
	}

	if input.MeasurementQuantity.Max != nil && input.MeasurementQuantity.Max != x.MeasurementQuantity.Max {
		x.MeasurementQuantity.Max = input.MeasurementQuantity.Max
	}

	if input.ItemQuantity.Min != nil && input.ItemQuantity.Min != x.ItemQuantity.Min {
		x.ItemQuantity.Min = input.ItemQuantity.Min
	}

	if input.ItemQuantity.Max != nil && input.ItemQuantity.Max != x.ItemQuantity.Max {
		x.ItemQuantity.Max = input.ItemQuantity.Max
	}

	if input.QuantityNotes != nil && *input.QuantityNotes != x.QuantityNotes {
		x.QuantityNotes = *input.QuantityNotes
	}

	if input.Compostable != nil && *input.Compostable != x.Compostable {
		x.Compostable = *input.Compostable
	}

	if input.StorageDurationInSeconds.Max != nil && input.StorageDurationInSeconds.Max != x.StorageDurationInSeconds.Max {
		x.StorageDurationInSeconds.Max = input.StorageDurationInSeconds.Max
	}

	if input.StorageTemperatureInCelsius.Min != nil && input.StorageTemperatureInCelsius.Min != x.StorageTemperatureInCelsius.Min {
		x.StorageTemperatureInCelsius.Min = input.StorageTemperatureInCelsius.Min
	}

	if input.StorageTemperatureInCelsius.Min != nil && input.StorageTemperatureInCelsius.Min != x.StorageTemperatureInCelsius.Min {
		x.StorageTemperatureInCelsius.Min = input.StorageTemperatureInCelsius.Min
	}

	if input.StorageTemperatureInCelsius.Max != nil && input.StorageTemperatureInCelsius.Max != x.StorageTemperatureInCelsius.Max {
		x.StorageTemperatureInCelsius.Max = input.StorageTemperatureInCelsius.Max
	}

	if input.StorageInstructions != nil && *input.StorageInstructions != x.StorageInstructions {
		x.StorageInstructions = *input.StorageInstructions
	}

	if input.IsLiquid != nil && *input.IsLiquid != x.IsLiquid {
		x.IsLiquid = *input.IsLiquid
	}

	if input.IsWaste != nil && *input.IsWaste != x.IsWaste {
		x.IsWaste = *input.IsWaste
	}

	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
	}

	if input.ContainedInVesselIndex != nil && x.ContainedInVesselIndex != nil && *input.ContainedInVesselIndex != *x.ContainedInVesselIndex {
		x.ContainedInVesselIndex = input.ContainedInVesselIndex
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepProductCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductCreationRequestInput.
func (x *RecipeStepProductCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType, RecipeStepProductVesselType)),
		validation.Field(&x.MeasurementQuantity, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepProductDatabaseCreationInput.
func (x *RecipeStepProductDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType, RecipeStepProductVesselType)),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepProductUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepProductUpdateRequestInput.
func (x *RecipeStepProductUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Type, validation.In(RecipeStepProductIngredientType, RecipeStepProductInstrumentType, RecipeStepProductVesselType)),
		validation.Field(&x.MeasurementUnitID, validation.Required),
		validation.Field(&x.MeasurementQuantity, validation.Required),
	)
}
