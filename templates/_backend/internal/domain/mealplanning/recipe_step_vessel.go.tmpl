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
	errValidPreparationVesselIDRequired = errors.New("validPreparationVesselID is required when not referencing a recipe step product")
)

const (
	// RecipeStepVesselCreatedServiceEventType indicates a recipe step instrument was created.
	RecipeStepVesselCreatedServiceEventType = "recipe_step_vessel_created"
	// RecipeStepVesselUpdatedServiceEventType indicates a recipe step instrument was updated.
	RecipeStepVesselUpdatedServiceEventType = "recipe_step_vessel_updated"
	// RecipeStepVesselArchivedServiceEventType indicates a recipe step instrument was archived.
	RecipeStepVesselArchivedServiceEventType = "recipe_step_vessel_archived"
)

func init() {
	gob.Register(new(RecipeStepVessel))
	gob.Register(new(RecipeStepVesselCreationRequestInput))
	gob.Register(new(RecipeStepVesselUpdateRequestInput))
}

type (
	// RecipeStepVessel represents a recipe step instrument.
	RecipeStepVessel struct {
		_ struct{} `json:"-"`

		CreatedAt            time.Time                `json:"createdAt"`
		Quantity             numbers.MinRange[uint16] `json:"quantity"`
		LastUpdatedAt        *time.Time               `json:"lastUpdatedAt"`
		ArchivedAt           *time.Time               `json:"archivedAt"`
		RecipeStepProductID  *string                  `json:"recipeStepProductID"`
		Vessel               *ValidVessel             `json:"vessel"`
		ID                   string                   `json:"id"`
		Notes                string                   `json:"notes"`
		BelongsToRecipeStep  string                   `json:"belongsToRecipeStep"`
		VesselPreposition    string                   `json:"vesselPreposition"`
		Name                 string                   `json:"name"`
		Index                uint16                   `json:"index"`
		OptionIndex          uint16                   `json:"optionIndex"`
		UnavailableAfterStep bool                     `json:"unavailableAfterStep"`
		ScaleFactor          float32                  `json:"scaleFactor"`
	}

	// RecipeStepVesselCreationRequestInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselCreationRequestInput struct {
		_                               struct{}                 `json:"-"`
		Quantity                        numbers.MinRange[uint16] `json:"quantity"`
		RecipeStepProductID             *string                  `json:"recipeStepProductID"`
		ProductOfRecipeStepIndex        *uint64                  `json:"productOfRecipeStepIndex"`
		ProductOfRecipeStepProductIndex *uint64                  `json:"productOfRecipeStepProductIndex"`
		ValidPreparationVesselID        *string                  `json:"validPreparationVesselID"`
		Index                           *uint16                  `json:"index,omitempty"`
		Name                            string                   `json:"name"`
		Notes                           string                   `json:"notes"`
		VesselPreposition               string                   `json:"vesselPreposition"`
		OptionIndex                     uint16                   `json:"optionIndex"`
		UnavailableAfterStep            bool                     `json:"unavailableAfterStep"`
		ScaleFactor                     float32                  `json:"scaleFactor"`
	}

	// RecipeStepVesselDatabaseCreationInput represents what a user could set as input for creating recipe step instruments.
	RecipeStepVesselDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		VesselID                        *string                  `json:"-"`
		ValidPreparationVesselID        *string                  `json:"-"`
		RecipeStepProductID             *string                  `json:"-"`
		ProductOfRecipeStepIndex        *uint64                  `json:"-"`
		ProductOfRecipeStepProductIndex *uint64                  `json:"-"`
		Quantity                        numbers.MinRange[uint16] `json:"-"`
		ID                              string                   `json:"-"`
		Notes                           string                   `json:"-"`
		BelongsToRecipeStep             string                   `json:"-"`
		VesselPreposition               string                   `json:"-"`
		Name                            string                   `json:"-"`
		Index                           uint16                   `json:"-"`
		OptionIndex                     uint16                   `json:"-"`
		UnavailableAfterStep            bool                     `json:"-"`
		ScaleFactor                     float32                  `json:"-"`
	}

	// RecipeStepVesselUpdateRequestInput represents what a user could set as input for updating recipe step instruments.
	RecipeStepVesselUpdateRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepProductID  *string                                     `json:"recipeStepProductID,omitempty"`
		Name                 *string                                     `json:"name,omitempty"`
		Notes                *string                                     `json:"notes,omitempty"`
		BelongsToRecipeStep  *string                                     `json:"belongsToRecipeStep,omitempty"`
		VesselID             *string                                     `json:"vesselID,omitempty"`
		Quantity             numbers.OpenRangeUpdateRequestInput[uint16] `json:"quantity"`
		Index                *uint16                                     `json:"index,omitempty"`
		OptionIndex          *uint16                                     `json:"optionIndex,omitempty"`
		VesselPreposition    *string                                     `json:"vesselPreposition,omitempty"`
		UnavailableAfterStep *bool                                       `json:"unavailableAfterStep,omitempty"`
		ScaleFactor          *float32                                    `json:"scaleFactor,omitempty"`
	}

	// RecipeStepVesselDataManager describes a structure capable of storing recipe step instruments permanently.
	RecipeStepVesselDataManager interface {
		RecipeStepVesselExists(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (bool, error)
		GetRecipeStepVessel(ctx context.Context, recipeID, recipeStepID, recipeStepInstrumentID string) (*RecipeStepVessel, error)
		GetRecipeStepVessels(ctx context.Context, recipeID, recipeStepID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[RecipeStepVessel], error)
		CreateRecipeStepVessel(ctx context.Context, input *RecipeStepVesselDatabaseCreationInput) (*RecipeStepVessel, error)
		UpdateRecipeStepVessel(ctx context.Context, updated *RecipeStepVessel) error
		ArchiveRecipeStepVessel(ctx context.Context, recipeStepID, recipeStepInstrumentID string) error
	}

	// RecipeStepVesselDataService describes a structure capable of serving traffic related to recipe step instruments.
	RecipeStepVesselDataService interface {
		ListRecipeStepVesselsHandler(http.ResponseWriter, *http.Request)
		CreateRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
		ReadRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
		UpdateRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
		ArchiveRecipeStepVesselHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepVesselUpdateRequestInput with a recipe step instrument.
func (x *RecipeStepVessel) Update(input *RecipeStepVesselUpdateRequestInput) {
	if input.RecipeStepProductID != nil && input.RecipeStepProductID != x.RecipeStepProductID {
		x.RecipeStepProductID = input.RecipeStepProductID
	}

	if input.Name != nil && *input.Name != x.Name {
		x.Name = *input.Name
	}

	if input.BelongsToRecipeStep != nil && *input.BelongsToRecipeStep != x.BelongsToRecipeStep {
		x.BelongsToRecipeStep = *input.BelongsToRecipeStep
	}

	if input.Quantity.Min != nil && *input.Quantity.Min != x.Quantity.Min {
		x.Quantity.Min = *input.Quantity.Min
	}

	if input.Quantity.Max != nil && x.Quantity.Max != nil && *input.Quantity.Max != *x.Quantity.Max {
		x.Quantity.Max = input.Quantity.Max
	}

	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
	}

	if input.OptionIndex != nil && *input.OptionIndex != x.OptionIndex {
		x.OptionIndex = *input.OptionIndex
	}

	if input.VesselPreposition != nil && *input.VesselPreposition != x.VesselPreposition {
		x.VesselPreposition = *input.VesselPreposition
	}

	if input.UnavailableAfterStep != nil && *input.UnavailableAfterStep != x.UnavailableAfterStep {
		x.UnavailableAfterStep = *input.UnavailableAfterStep
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.VesselID != nil && x.Vessel != nil && *input.VesselID != x.Vessel.ID {
		x.Vessel = &ValidVessel{ID: *input.VesselID}
	}

	if input.ScaleFactor != nil && *input.ScaleFactor != x.ScaleFactor {
		x.ScaleFactor = *input.ScaleFactor
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselCreationRequestInput.
func (x *RecipeStepVesselCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	err := &multierror.Error{}

	// When not referencing a recipe step product, bridge table ID is required
	if isRecipeStepProduct := x.ProductOfRecipeStepIndex != nil; !isRecipeStepProduct {
		if x.ValidPreparationVesselID == nil || *x.ValidPreparationVesselID == "" {
			err = multierror.Append(err, errValidPreparationVesselIDRequired)
		}
	}

	validationErr := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Name, validation.Required),
		validation.Field(&x.Quantity, validation.Required),
	)
	if validationErr != nil {
		err = multierror.Append(err, validationErr)
	}

	return err.ErrorOrNil()
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepVesselDatabaseCreationInput.
func (x *RecipeStepVesselDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepVesselUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepVesselUpdateRequestInput.
func (x *RecipeStepVesselUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}
