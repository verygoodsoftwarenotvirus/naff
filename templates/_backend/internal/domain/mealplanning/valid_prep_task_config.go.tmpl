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
	// ValidPrepTaskConfigCreatedServiceEventType indicates a valid ingredient preparation storage config was created.
	ValidPrepTaskConfigCreatedServiceEventType = "valid_prep_task_config_created"
	// ValidPrepTaskConfigUpdatedServiceEventType indicates a valid ingredient preparation storage config was updated.
	ValidPrepTaskConfigUpdatedServiceEventType = "valid_prep_task_config_updated"
	// ValidPrepTaskConfigArchivedServiceEventType indicates a valid ingredient preparation storage config was archived.
	ValidPrepTaskConfigArchivedServiceEventType = "valid_prep_task_config_archived"
)

func init() {
	gob.Register(new(ValidPrepTaskConfig))
	gob.Register(new(ValidPrepTaskConfigCreationRequestInput))
	gob.Register(new(ValidPrepTaskConfigUpdateRequestInput))
}

type (
	// ValidPrepTaskConfig represents reusable knowledge about how long
	// a prepared ingredient can be stored under specific conditions.
	// Example: "diced onion" can be stored for 72 hours in an airtight container at 0-4°C.
	ValidPrepTaskConfig struct {
		_ struct{} `json:"-"`

		CreatedAt                   time.Time                  `json:"createdAt"`
		LastUpdatedAt               *time.Time                 `json:"lastUpdatedAt"`
		ArchivedAt                  *time.Time                 `json:"archivedAt"`
		StorageDurationInSeconds    numbers.MinRange[uint32]   `json:"storageDurationInSeconds"`
		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"storageTemperatureInCelsius"`
		ID                          string                     `json:"id"`
		StorageType                 string                     `json:"storageType"`
		StorageInstructions         string                     `json:"storageInstructions"`
		Notes                       string                     `json:"notes"`
		Source                      string                     `json:"source"`
		Preparation                 ValidPreparation           `json:"preparation"`
		Ingredient                  ValidIngredient            `json:"ingredient"`
	}

	// ValidPrepTaskConfigCreationRequestInput represents what a user could set as input for creating valid ingredient preparation storage configs.
	ValidPrepTaskConfigCreationRequestInput struct {
		_ struct{} `json:"-"`

		StorageDurationInSeconds    numbers.MinRange[uint32]   `json:"storageDurationInSeconds"`
		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"storageTemperatureInCelsius"`
		StorageType                 string                     `json:"storageType"`
		StorageInstructions         string                     `json:"storageInstructions"`
		Notes                       string                     `json:"notes"`
		Source                      string                     `json:"source"`
		ValidPreparationID          string                     `json:"validPreparationID"`
		ValidIngredientID           string                     `json:"validIngredientID"`
	}

	// ValidPrepTaskConfigDatabaseCreationInput represents what a user could set as input for creating valid ingredient preparation storage configs.
	ValidPrepTaskConfigDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		StorageDurationInSeconds    numbers.MinRange[uint32]   `json:"-"`
		StorageTemperatureInCelsius numbers.OpenRange[float32] `json:"-"`
		ID                          string                     `json:"-"`
		StorageType                 string                     `json:"-"`
		StorageInstructions         string                     `json:"-"`
		Notes                       string                     `json:"-"`
		Source                      string                     `json:"-"`
		ValidPreparationID          string                     `json:"-"`
		ValidIngredientID           string                     `json:"-"`
	}

	// ValidPrepTaskConfigUpdateRequestInput represents what a user could set as input for updating valid ingredient preparation storage configs.
	ValidPrepTaskConfigUpdateRequestInput struct {
		_ struct{} `json:"-"`

		StorageDurationInSeconds    numbers.OpenRangeUpdateRequestInput[uint32] `json:"storageDurationInSeconds"`
		StorageTemperatureInCelsius numbers.OpenRange[float32]                  `json:"storageTemperatureInCelsius"`
		StorageType                 *string                                     `json:"storageType,omitempty"`
		StorageInstructions         *string                                     `json:"storageInstructions,omitempty"`
		Notes                       *string                                     `json:"notes,omitempty"`
		Source                      *string                                     `json:"source,omitempty"`
		ValidPreparationID          *string                                     `json:"validPreparationID,omitempty"`
		ValidIngredientID           *string                                     `json:"validIngredientID,omitempty"`
	}

	// ValidPrepTaskConfigDataManager describes a structure capable of storing valid ingredient preparation storage configs permanently.
	ValidPrepTaskConfigDataManager interface {
		ValidPrepTaskConfigExists(ctx context.Context, validIngredientPreparationStorageConfigID string) (bool, error)
		GetValidPrepTaskConfig(ctx context.Context, validIngredientPreparationStorageConfigID string) (*ValidPrepTaskConfig, error)
		GetValidPrepTaskConfigs(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ValidPrepTaskConfig], error)
		GetValidPrepTaskConfigsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ValidPrepTaskConfig], error)
		GetValidPrepTaskConfigsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ValidPrepTaskConfig], error)
		GetValidPrepTaskConfigsForIngredientAndPreparation(ctx context.Context, ingredientID, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[ValidPrepTaskConfig], error)
		CreateValidPrepTaskConfig(ctx context.Context, input *ValidPrepTaskConfigDatabaseCreationInput) (*ValidPrepTaskConfig, error)
		UpdateValidPrepTaskConfig(ctx context.Context, updated *ValidPrepTaskConfig) error
		ArchiveValidPrepTaskConfig(ctx context.Context, validIngredientPreparationStorageConfigID string) error
	}

	// ValidPrepTaskConfigDataService describes a structure capable of serving traffic related to valid ingredient preparation storage configs.
	ValidPrepTaskConfigDataService interface {
		ListValidPrepTaskConfigsHandler(http.ResponseWriter, *http.Request)
		CreateValidPrepTaskConfigHandler(http.ResponseWriter, *http.Request)
		ReadValidPrepTaskConfigHandler(http.ResponseWriter, *http.Request)
		UpdateValidPrepTaskConfigHandler(http.ResponseWriter, *http.Request)
		ArchiveValidPrepTaskConfigHandler(http.ResponseWriter, *http.Request)
		ListValidPrepTaskConfigsForIngredientHandler(http.ResponseWriter, *http.Request)
		ListValidPrepTaskConfigsForPreparationHandler(http.ResponseWriter, *http.Request)
		ListValidPrepTaskConfigsForIngredientAndPreparationHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges a ValidPrepTaskConfigUpdateRequestInput with a valid ingredient preparation storage config.
func (x *ValidPrepTaskConfig) Update(input *ValidPrepTaskConfigUpdateRequestInput) {
	if input.StorageType != nil && *input.StorageType != x.StorageType {
		x.StorageType = *input.StorageType
	}

	if input.StorageInstructions != nil && *input.StorageInstructions != x.StorageInstructions {
		x.StorageInstructions = *input.StorageInstructions
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.Source != nil && *input.Source != x.Source {
		x.Source = *input.Source
	}

	if input.ValidPreparationID != nil && *input.ValidPreparationID != x.Preparation.ID {
		x.Preparation.ID = *input.ValidPreparationID
	}

	if input.ValidIngredientID != nil && *input.ValidIngredientID != x.Ingredient.ID {
		x.Ingredient.ID = *input.ValidIngredientID
	}

	if input.StorageDurationInSeconds.Min != nil && *input.StorageDurationInSeconds.Min != x.StorageDurationInSeconds.Min {
		x.StorageDurationInSeconds.Min = *input.StorageDurationInSeconds.Min
	}

	if input.StorageDurationInSeconds.Max != nil && input.StorageDurationInSeconds.Max != x.StorageDurationInSeconds.Max {
		x.StorageDurationInSeconds.Max = input.StorageDurationInSeconds.Max
	}

	if input.StorageTemperatureInCelsius.Min != nil && input.StorageTemperatureInCelsius.Min != x.StorageTemperatureInCelsius.Min {
		x.StorageTemperatureInCelsius.Min = input.StorageTemperatureInCelsius.Min
	}

	if input.StorageTemperatureInCelsius.Max != nil && input.StorageTemperatureInCelsius.Max != x.StorageTemperatureInCelsius.Max {
		x.StorageTemperatureInCelsius.Max = input.StorageTemperatureInCelsius.Max
	}
}

var _ validation.ValidatableWithContext = (*ValidPrepTaskConfigCreationRequestInput)(nil)

// ValidateWithContext validates a ValidPrepTaskConfigCreationRequestInput.
func (x *ValidPrepTaskConfigCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.StorageType, validation.Required),
		validation.Field(&x.StorageDurationInSeconds, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPrepTaskConfigDatabaseCreationInput)(nil)

// ValidateWithContext validates a ValidPrepTaskConfigDatabaseCreationInput.
func (x *ValidPrepTaskConfigDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
		validation.Field(&x.StorageType, validation.Required),
		validation.Field(&x.StorageDurationInSeconds, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ValidPrepTaskConfigUpdateRequestInput)(nil)

// ValidateWithContext validates a ValidPrepTaskConfigUpdateRequestInput.
func (x *ValidPrepTaskConfigUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ValidPreparationID, validation.Required),
		validation.Field(&x.ValidIngredientID, validation.Required),
	)
}
