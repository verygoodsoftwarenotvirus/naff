package mealplanning

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/primandproper/platform/database/filtering"
	"github.com/primandproper/platform/numbers"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	// RecipePrepTaskStorageTypeUncovered is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeUncovered = "uncovered"
	// RecipePrepTaskStorageTypeCovered is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeCovered = "covered"
	// RecipePrepTaskStorageTypeWireRack is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeWireRack = "on a wire rack"
	// RecipePrepTaskStorageTypeAirtightContainer is a valid storage type for a recipe step task.
	RecipePrepTaskStorageTypeAirtightContainer = "in an airtight container"

	// RecipePrepTaskCreatedServiceEventType indicates a recipe prep task was created.
	RecipePrepTaskCreatedServiceEventType = "recipe_created"
	// RecipePrepTaskUpdatedServiceEventType indicates a recipe prep task was updated.
	RecipePrepTaskUpdatedServiceEventType = "recipe_updated"
	// RecipePrepTaskArchivedServiceEventType indicates a recipe prep task was archived.
	RecipePrepTaskArchivedServiceEventType = "recipe_archived"
)

func init() {
	gob.Register(new(RecipePrepTask))
	gob.Register(new(RecipePrepTaskCreationRequestInput))
	gob.Register(new(RecipePrepTaskUpdateRequestInput))
}

type (
	// RecipePrepTask represents a recipe prep task.
	RecipePrepTask struct {
		_ struct{} `json:"-"`

		CreatedAt                       time.Time                  `json:"createdAt"`
		StorageTemperatureInCelsius     numbers.OpenRange[float32] `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds numbers.MinRange[uint32]   `json:"timeBufferBeforeRecipeInSeconds"`
		ArchivedAt                      *time.Time                 `json:"archivedAt"`
		LastUpdatedAt                   *time.Time                 `json:"lastUpdatedAt"`
		ID                              string                     `json:"id"`
		StorageType                     string                     `json:"storageType"`
		BelongsToRecipe                 string                     `json:"belongsToRecipe"`
		ExplicitStorageInstructions     string                     `json:"explicitStorageInstructions"`
		Notes                           string                     `json:"notes"`
		Name                            string                     `json:"name"`
		Description                     string                     `json:"description"`
		TaskSteps                       []*RecipePrepTaskStep      `json:"recipeSteps"`
		Optional                        bool                       `json:"optional"`
	}

	// RecipePrepTaskCreationRequestInput represents what a user could set as input for creating recipes.
	RecipePrepTaskCreationRequestInput struct {
		_ struct{} `json:"-"`

		StorageTemperatureInCelsius     numbers.OpenRange[float32]                `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds numbers.MinRange[uint32]                  `json:"timeBufferBeforeRecipeInSeconds"`
		StorageType                     string                                    `json:"storageType"`
		ExplicitStorageInstructions     string                                    `json:"explicitStorageInstructions"`
		Notes                           string                                    `json:"notes"`
		Name                            string                                    `json:"name"`
		Description                     string                                    `json:"description"`
		BelongsToRecipe                 string                                    `json:"belongsToRecipe"`
		RecipeSteps                     []*RecipePrepTaskStepCreationRequestInput `json:"recipeSteps"`
		Optional                        bool                                      `json:"optional"`
	}

	// RecipePrepTaskWithinRecipeCreationRequestInput represents what a user could set as input for creating recipes.
	RecipePrepTaskWithinRecipeCreationRequestInput struct {
		_ struct{} `json:"-"`

		StorageTemperatureInCelsius     numbers.OpenRange[float32]                            `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds numbers.MinRange[uint32]                              `json:"timeBufferBeforeRecipeInSeconds"`
		StorageType                     string                                                `json:"storageType"`
		Name                            string                                                `json:"name"`
		Description                     string                                                `json:"description"`
		ExplicitStorageInstructions     string                                                `json:"explicitStorageInstructions"`
		Notes                           string                                                `json:"notes"`
		BelongsToRecipe                 string                                                `json:"belongsToRecipe"`
		RecipeSteps                     []*RecipePrepTaskStepWithinRecipeCreationRequestInput `json:"recipeSteps"`
		Optional                        bool                                                  `json:"optional"`
	}

	// RecipePrepTaskDatabaseCreationInput represents what a user could set as input for creating recipes.
	RecipePrepTaskDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		StorageTemperatureInCelsius     numbers.OpenRange[float32]                 `json:"-"`
		TimeBufferBeforeRecipeInSeconds numbers.MinRange[uint32]                   `json:"-"`
		ExplicitStorageInstructions     string                                     `json:"-"`
		Notes                           string                                     `json:"-"`
		ID                              string                                     `json:"-"`
		Name                            string                                     `json:"-"`
		Description                     string                                     `json:"-"`
		StorageType                     string                                     `json:"-"`
		BelongsToRecipe                 string                                     `json:"-"`
		TaskSteps                       []*RecipePrepTaskStepDatabaseCreationInput `json:"-"`
		Optional                        bool                                       `json:"-"`
	}

	// RecipePrepTaskUpdateRequestInput represents what a user could set as input for updating recipes.
	RecipePrepTaskUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes                           *string                                     `json:"notes,omitempty"`
		ExplicitStorageInstructions     *string                                     `json:"explicitStorageInstructions,omitempty"`
		StorageType                     *string                                     `json:"storageType,omitempty"`
		Name                            *string                                     `json:"name"`
		Optional                        *bool                                       `json:"optional"`
		Description                     *string                                     `json:"description"`
		StorageTemperatureInCelsius     numbers.OpenRange[float32]                  `json:"storageTemperatureInCelsius"`
		TimeBufferBeforeRecipeInSeconds numbers.OpenRangeUpdateRequestInput[uint32] `json:"timeBufferBeforeRecipeInSeconds"`
		BelongsToRecipe                 *string                                     `json:"belongsToRecipe,omitempty"`
		TaskSteps                       []*RecipePrepTaskStepUpdateRequestInput     `json:"recipeSteps,omitempty"`
	}

	// RecipePrepTaskDataManager describes a structure capable of storing recipes permanently.
	RecipePrepTaskDataManager interface {
		RecipePrepTaskExists(ctx context.Context, recipeID, recipePrepTaskID string) (bool, error)
		GetRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) (*RecipePrepTask, error)
		GetRecipePrepTasks(ctx context.Context, recipeID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[RecipePrepTask], error)
		GetRecipePrepTasksForRecipe(ctx context.Context, recipeID string) ([]*RecipePrepTask, error)
		CreateRecipePrepTask(ctx context.Context, input *RecipePrepTaskDatabaseCreationInput) (*RecipePrepTask, error)
		UpdateRecipePrepTask(ctx context.Context, updated *RecipePrepTask) error
		ArchiveRecipePrepTask(ctx context.Context, recipeID, recipePrepTaskID string) error
	}

	// RecipePrepTaskDataService describes a structure capable of serving traffic related to recipes.
	RecipePrepTaskDataService interface {
		ListRecipePrepTaskHandler(http.ResponseWriter, *http.Request)
		CreateRecipePrepTaskHandler(http.ResponseWriter, *http.Request)
		ReadRecipePrepTaskHandler(http.ResponseWriter, *http.Request)
		UpdateRecipePrepTaskHandler(http.ResponseWriter, *http.Request)
		ArchiveRecipePrepTaskHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipePrepTaskUpdateRequestInput with a recipe prep task.
func (x *RecipePrepTask) Update(input *RecipePrepTaskUpdateRequestInput) {
	if input.BelongsToRecipe != nil && *input.BelongsToRecipe != x.BelongsToRecipe {
		x.BelongsToRecipe = *input.BelongsToRecipe
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.ExplicitStorageInstructions != nil && *input.ExplicitStorageInstructions != x.ExplicitStorageInstructions {
		x.ExplicitStorageInstructions = *input.ExplicitStorageInstructions
	}

	if input.TimeBufferBeforeRecipeInSeconds.Min != nil && *input.TimeBufferBeforeRecipeInSeconds.Min != x.TimeBufferBeforeRecipeInSeconds.Min {
		x.TimeBufferBeforeRecipeInSeconds.Min = *input.TimeBufferBeforeRecipeInSeconds.Min
	}

	if input.TimeBufferBeforeRecipeInSeconds.Max != nil && input.TimeBufferBeforeRecipeInSeconds.Max != x.TimeBufferBeforeRecipeInSeconds.Max {
		x.TimeBufferBeforeRecipeInSeconds.Max = input.TimeBufferBeforeRecipeInSeconds.Max
	}

	if input.StorageType != nil && *input.StorageType != x.StorageType {
		x.StorageType = *input.StorageType
	}

	if input.StorageTemperatureInCelsius.Min != nil && input.StorageTemperatureInCelsius.Min != x.StorageTemperatureInCelsius.Min {
		x.StorageTemperatureInCelsius.Min = input.StorageTemperatureInCelsius.Min
	}

	if input.StorageTemperatureInCelsius.Max != nil && input.StorageTemperatureInCelsius.Max != x.StorageTemperatureInCelsius.Max {
		x.StorageTemperatureInCelsius.Max = input.StorageTemperatureInCelsius.Max
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.Description != nil && *input.Description != x.Description {
		x.Description = *input.Description
	}

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskCreationRequestInput)(nil)

// ValidateWithContext validates a RecipePrepTaskCreationRequestInput.
func (x *RecipePrepTaskCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	result := &multierror.Error{}

	if err := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.StorageType, validation.Required),
		validation.Field(&x.TimeBufferBeforeRecipeInSeconds, validation.Required),
		validation.Field(&x.StorageTemperatureInCelsius, validation.Required),
	); err != nil {
		result = multierror.Append(err, result)
	}

	// TODO: uncomment me
	// if x.StorageTemperatureInCelsius.Min != nil && x.StorageTemperatureInCelsius.Max != nil && *x.StorageTemperatureInCelsius.Min > *x.StorageTemperatureInCelsius.Max {
	//	result = multierror.Append(fmt.Errorf("minimum storage temperature (%d) is greater than maximum storage temperature (%d)", x.StorageTemperatureInCelsius.Min, x.StorageTemperatureInCelsius.Max))
	// }

	return result.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipePrepTaskDatabaseCreationInput.
func (x *RecipePrepTaskDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipe, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipePrepTaskUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipePrepTaskUpdateRequestInput.
func (x *RecipePrepTaskUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipe, validation.Required),
	)
}
