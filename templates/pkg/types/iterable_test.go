package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterableDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterableDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	"net/http"
)

const (
	// ItemsSearchIndexName is the name of the index used to search through items.
	ItemsSearchIndexName search.IndexName = "items"
)

type (
	// Item represents an item.
	Item struct {
		ID            uint64  ` + "`" + `json:"id"` + "`" + `
		Name          string  ` + "`" + `json:"name"` + "`" + `
		Details       string  ` + "`" + `json:"details"` + "`" + `
		CreatedOn     uint64  ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn    *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToAccount uint64  ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	// ItemList represents a list of items.
	ItemList struct {
		Pagination
		Items []Item ` + "`" + `json:"items"` + "`" + `
	}

	// ItemCreationInput represents what a user could set as input for creating items.
	ItemCreationInput struct {
		Name          string ` + "`" + `json:"name"` + "`" + `
		Details       string ` + "`" + `json:"details"` + "`" + `
		BelongsToAccount uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// ItemUpdateInput represents what a user could set as input for updating items.
	ItemUpdateInput struct {
		Name          string ` + "`" + `json:"name"` + "`" + `
		Details       string ` + "`" + `json:"details"` + "`" + `
		BelongsToAccount uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// ItemDataManager describes a structure capable of storing items permanently.
	ItemDataManager interface {
		ItemExists(ctx context.Context, itemID, userID uint64) (bool, error)
		GetItem(ctx context.Context, itemID, userID uint64) (*Item, error)
		GetAllItemsCount(ctx context.Context) (uint64, error)
		GetAllItems(ctx context.Context, resultChannel chan []Item) error
		GetItems(ctx context.Context, userID uint64, filter *QueryFilter) (*ItemList, error)
		GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]Item, error)
		CreateItem(ctx context.Context, input *ItemCreationInput) (*Item, error)
		UpdateItem(ctx context.Context, updated *Item) error
		ArchiveItem(ctx context.Context, itemID, userID uint64) error
	}

	// ItemDataServer describes a structure capable of serving traffic related to items.
	ItemDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ItemInput with an item.
func (x *Item) Update(input *ItemUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Details != "" && input.Details != x.Details {
		x.Details = input.Details
	}
}

// ToUpdateInput creates a ItemUpdateInput struct for an item.
func (x *Item) ToUpdateInput() *ItemUpdateInput {
	return &ItemUpdateInput{
		Name:    x.Name,
		Details: x.Details,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with search enabled and owner types", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		typ.Fields = []models.DataField{
			{
				Name:                wordsmith.FromSingularPascalCase("FieldOne"),
				ValidForUpdateInput: true,
				Type:                "string",
			},
			{
				Name:                wordsmith.FromSingularPascalCase("FieldTwo"),
				ValidForUpdateInput: true,
				Type:                "string",
			},
			{
				Name:                wordsmith.FromSingularPascalCase("FieldThree"),
				ValidForUpdateInput: true,
				Type:                "string",
			},
		}
		typ.SearchEnabled = true
		x := iterableDotGo(proj, typ)

		//+ "`" + `

		expected := `
package example

import (
	"context"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	"net/http"
)

const (
	// YetAnotherThingsSearchIndexName is the name of the index used to search through yet another things.
	YetAnotherThingsSearchIndexName search.IndexName = "yet_another_things"
)

type (
	// YetAnotherThing represents a yet another thing.
	YetAnotherThing struct {
		ID                    uint64  ` + "`" + `json:"id"` + "`" + `
		FieldOne              string  ` + "`" + `json:"fieldOne"` + "`" + `
		FieldTwo              string  ` + "`" + `json:"fieldTwo"` + "`" + `
		FieldThree            string  ` + "`" + `json:"fieldThree"` + "`" + `
		CreatedOn             uint64  ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn         *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn            *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToAnotherThing uint64  ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
	}

	// YetAnotherThingList represents a list of yet another things.
	YetAnotherThingList struct {
		Pagination
		YetAnotherThings []YetAnotherThing ` + "`" + `json:"yetAnotherThings"` + "`" + `
	}
	// YetAnotherThingSearchHelper contains all the owner IDs for search purposes.
	YetAnotherThingSearchHelper struct {
		ID                    uint64  ` + "`" + `json:"id"` + "`" + `
		FieldOne              string  ` + "`" + `json:"fieldOne"` + "`" + `
		FieldTwo              string  ` + "`" + `json:"fieldTwo"` + "`" + `
		FieldThree            string  ` + "`" + `json:"fieldThree"` + "`" + `
		CreatedOn             uint64  ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn         *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn            *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToAnotherThing uint64  ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
		BelongsToThing        uint64  ` + "`" + `json:"belongsToThing"` + "`" + `
		BelongsToAnotherThing uint64  ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
	}

	// YetAnotherThingCreationInput represents what a user could set as input for creating yet another things.
	YetAnotherThingCreationInput struct {
		BelongsToAnotherThing uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// YetAnotherThingUpdateInput represents what a user could set as input for updating yet another things.
	YetAnotherThingUpdateInput struct {
		FieldOne              string ` + "`" + `json:"fieldOne"` + "`" + `
		FieldTwo              string ` + "`" + `json:"fieldTwo"` + "`" + `
		FieldThree            string ` + "`" + `json:"fieldThree"` + "`" + `
		BelongsToAnotherThing uint64 ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
	}

	// YetAnotherThingDataManager describes a structure capable of storing yet another things permanently.
	YetAnotherThingDataManager interface {
		YetAnotherThingExists(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) (bool, error)
		GetYetAnotherThing(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) (*YetAnotherThing, error)
		GetAllYetAnotherThingsCount(ctx context.Context) (uint64, error)
		GetAllYetAnotherThings(ctx context.Context, resultChannel chan []YetAnotherThing) error
		GetYetAnotherThings(ctx context.Context, thingID, anotherThingID uint64, filter *QueryFilter) (*YetAnotherThingList, error)
		GetYetAnotherThingsWithIDs(ctx context.Context, thingID, anotherThingID uint64, limit uint8, ids []uint64) ([]YetAnotherThing, error)
		CreateYetAnotherThing(ctx context.Context, input *YetAnotherThingCreationInput) (*YetAnotherThing, error)
		UpdateYetAnotherThing(ctx context.Context, updated *YetAnotherThing) error
		ArchiveYetAnotherThing(ctx context.Context, anotherThingID, yetAnotherThingID uint64) error
	}

	// YetAnotherThingDataServer describes a structure capable of serving traffic related to yet another things.
	YetAnotherThingDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an YetAnotherThingInput with a yet another thing.
func (x *YetAnotherThing) Update(input *YetAnotherThingUpdateInput) {
	if input.FieldOne != "" && input.FieldOne != x.FieldOne {
		x.FieldOne = input.FieldOne
	}

	if input.FieldTwo != "" && input.FieldTwo != x.FieldTwo {
		x.FieldTwo = input.FieldTwo
	}

	if input.FieldThree != "" && input.FieldThree != x.FieldThree {
		x.FieldThree = input.FieldThree
	}
}

// ToUpdateInput creates a YetAnotherThingUpdateInput struct for a yet another thing.
func (x *YetAnotherThing) ToUpdateInput() *YetAnotherThingUpdateInput {
	return &YetAnotherThingUpdateInput{
		FieldOne:   x.FieldOne,
		FieldTwo:   x.FieldTwo,
		FieldThree: x.FieldThree,
	}
}

// ToSearchHelper creates a YetAnotherThingSearchHelper struct for a yet another thing.
func (x *YetAnotherThing) ToSearchHelper(thingID uint64, anotherThingID uint64) *YetAnotherThingSearchHelper {
	return &YetAnotherThingSearchHelper{
		FieldOne:              x.FieldOne,
		FieldTwo:              x.FieldTwo,
		FieldThree:            x.FieldThree,
		BelongsToThing:        thingID,
		BelongsToAnotherThing: anotherThingID,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomething(typ)

		expected := `
package example

import ()

// Update merges an ItemInput with an item.
func (x *Item) Update(input *ItemUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Details != "" && input.Details != x.Details {
		x.Details = input.Details
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingConstantDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.SearchEnabled = true

		x := buildSomethingConstantDefinitions(proj, typ)

		expected := `
package example

import (
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
)

const (
	// ItemsSearchIndexName is the name of the index used to search through items.
	ItemsSearchIndexName search.IndexName = "items"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingTypeDefinitions(proj, typ)

		expected := `
package example

import (
	"context"
	"net/http"
)

type (
	// Item represents an item.
	Item struct {
		ID            uint64  ` + "`" + `json:"id"` + "`" + `
		Name          string  ` + "`" + `json:"name"` + "`" + `
		Details       string  ` + "`" + `json:"details"` + "`" + `
		CreatedOn     uint64  ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn    *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToAccount uint64  ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	// ItemList represents a list of items.
	ItemList struct {
		Pagination
		Items []Item ` + "`" + `json:"items"` + "`" + `
	}

	// ItemCreationInput represents what a user could set as input for creating items.
	ItemCreationInput struct {
		Name          string ` + "`" + `json:"name"` + "`" + `
		Details       string ` + "`" + `json:"details"` + "`" + `
		BelongsToAccount uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// ItemUpdateInput represents what a user could set as input for updating items.
	ItemUpdateInput struct {
		Name          string ` + "`" + `json:"name"` + "`" + `
		Details       string ` + "`" + `json:"details"` + "`" + `
		BelongsToAccount uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// ItemDataManager describes a structure capable of storing items permanently.
	ItemDataManager interface {
		ItemExists(ctx context.Context, itemID, userID uint64) (bool, error)
		GetItem(ctx context.Context, itemID, userID uint64) (*Item, error)
		GetAllItemsCount(ctx context.Context) (uint64, error)
		GetAllItems(ctx context.Context, resultChannel chan []Item) error
		GetItems(ctx context.Context, userID uint64, filter *QueryFilter) (*ItemList, error)
		GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]Item, error)
		CreateItem(ctx context.Context, input *ItemCreationInput) (*Item, error)
		UpdateItem(ctx context.Context, updated *Item) error
		ArchiveItem(ctx context.Context, itemID, userID uint64) error
	}

	// ItemDataServer describes a structure capable of serving traffic related to items.
	ItemDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with search enabled and ownership", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		typ.SearchEnabled = true
		x := buildSomethingTypeDefinitions(proj, typ)

		expected := `
package example

import (
	"context"
	"net/http"
)

type (
	// YetAnotherThing represents a yet another thing.
	YetAnotherThing struct {
		ID                    uint64  ` + "`" + `json:"id"` + "`" + `
		CreatedOn             uint64  ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn         *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn            *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToAnotherThing uint64  ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
	}

	// YetAnotherThingList represents a list of yet another things.
	YetAnotherThingList struct {
		Pagination
		YetAnotherThings []YetAnotherThing ` + "`" + `json:"yetAnotherThings"` + "`" + `
	}
	// YetAnotherThingSearchHelper contains all the owner IDs for search purposes.
	YetAnotherThingSearchHelper struct {
		ID                    uint64  ` + "`" + `json:"id"` + "`" + `
		CreatedOn             uint64  ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn         *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn            *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToAnotherThing uint64  ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
		BelongsToThing        uint64  ` + "`" + `json:"belongsToThing"` + "`" + `
		BelongsToAnotherThing uint64  ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
	}

	// YetAnotherThingCreationInput represents what a user could set as input for creating yet another things.
	YetAnotherThingCreationInput struct {
		BelongsToAnotherThing uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// YetAnotherThingUpdateInput represents what a user could set as input for updating yet another things.
	YetAnotherThingUpdateInput struct {
		BelongsToAnotherThing uint64 ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
	}

	// YetAnotherThingDataManager describes a structure capable of storing yet another things permanently.
	YetAnotherThingDataManager interface {
		YetAnotherThingExists(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) (bool, error)
		GetYetAnotherThing(ctx context.Context, thingID, anotherThingID, yetAnotherThingID uint64) (*YetAnotherThing, error)
		GetAllYetAnotherThingsCount(ctx context.Context) (uint64, error)
		GetAllYetAnotherThings(ctx context.Context, resultChannel chan []YetAnotherThing) error
		GetYetAnotherThings(ctx context.Context, thingID, anotherThingID uint64, filter *QueryFilter) (*YetAnotherThingList, error)
		GetYetAnotherThingsWithIDs(ctx context.Context, thingID, anotherThingID uint64, limit uint8, ids []uint64) ([]YetAnotherThing, error)
		CreateYetAnotherThing(ctx context.Context, input *YetAnotherThingCreationInput) (*YetAnotherThing, error)
		UpdateYetAnotherThing(ctx context.Context, updated *YetAnotherThing) error
		ArchiveYetAnotherThing(ctx context.Context, anotherThingID, yetAnotherThingID uint64) error
	}

	// YetAnotherThingDataServer describes a structure capable of serving traffic related to yet another things.
	YetAnotherThingDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingToUpdateInput(typ)

		expected := `
package example

import ()

// ToUpdateInput creates a ItemUpdateInput struct for an item.
func (x *Item) ToUpdateInput() *ItemUpdateInput {
	return &ItemUpdateInput{
		Name:    x.Name,
		Details: x.Details,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingToSearchHelper(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		typ.Fields = []models.DataField{
			{
				Name: wordsmith.FromSingularPascalCase("FieldOne"),
				Type: "string",
			},
			{
				Name: wordsmith.FromSingularPascalCase("FieldTwo"),
				Type: "string",
			},
			{
				Name: wordsmith.FromSingularPascalCase("FieldThree"),
				Type: "string",
			},
		}
		typ.SearchEnabled = true
		x := buildSomethingToSearchHelper(proj, typ)

		expected := `
package example

import ()

// ToSearchHelper creates a YetAnotherThingSearchHelper struct for a yet another thing.
func (x *YetAnotherThing) ToSearchHelper(thingID uint64, anotherThingID uint64) *YetAnotherThingSearchHelper {
	return &YetAnotherThingSearchHelper{
		FieldOne:              x.FieldOne,
		FieldTwo:              x.FieldTwo,
		FieldThree:            x.FieldThree,
		BelongsToThing:        thingID,
		BelongsToAnotherThing: anotherThingID,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBaseModelStructFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBaseModelStructFields(typ)

		expected := `
package example

import ()

type Example struct {
	ID            uint64  ` + "`" + `json:"id"` + "`" + `
	Name          string  ` + "`" + `json:"name"` + "`" + `
	Details       string  ` + "`" + `json:"details"` + "`" + `
	CreatedOn     uint64  ` + "`" + `json:"createdOn"` + "`" + `
	LastUpdatedOn *uint64 ` + "`" + `json:"lastUpdatedOn"` + "`" + `
	ArchivedOn    *uint64 ` + "`" + `json:"archivedOn"` + "`" + `
	BelongsToAccount uint64  ` + "`" + `json:"belongsToUser"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateModelStructFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateModelStructFields(typ)

		expected := `
package example

import ()

type Example struct {
	Name          string ` + "`" + `json:"name"` + "`" + `
	Details       string ` + "`" + `json:"details"` + "`" + `
	BelongsToAccount uint64 ` + "`" + `json:"-"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildUpdateModelStructFields(proj.LastDataType())

		expected := `
package example

import ()

type Example struct {
	BelongsToAnotherThing uint64 ` + "`" + `json:"belongsToAnotherThing"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateModelStructFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateModelStructFields(typ)

		expected := `
package example

import ()

type Example struct {
	Name          string ` + "`" + `json:"name"` + "`" + `
	Details       string ` + "`" + `json:"details"` + "`" + `
	BelongsToAccount uint64 ` + "`" + `json:"-"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfaceMethods(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildInterfaceMethods(proj, typ)

		expected := `
package example

import (
	"context"
)

type Example interface {
	ItemExists(ctx context.Context, itemID, userID uint64) (bool, error)
	GetItem(ctx context.Context, itemID, userID uint64) (*Item, error)
	GetAllItemsCount(ctx context.Context) (uint64, error)
	GetAllItems(ctx context.Context, resultChannel chan []Item) error
	GetItems(ctx context.Context, userID uint64, filter *QueryFilter) (*ItemList, error)
	GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]Item, error)
	CreateItem(ctx context.Context, input *ItemCreationInput) (*Item, error)
	UpdateItem(ctx context.Context, updated *Item) error
	ArchiveItem(ctx context.Context, itemID, userID uint64) error
}
`
		actual := testutils.RenderInterfaceMethodsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateFunctionLogic(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUpdateFunctionLogic(proj.DataTypes[0].Fields)

		expected := `
package main

import ()

func main() {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Details != "" && input.Details != x.Details {
		x.Details = input.Details
	}
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with every type", func(t *testing.T) {
		proj := testprojects.BuildEveryTypeApp()
		x := buildUpdateFunctionLogic(proj.DataTypes[0].Fields)

		expected := `
package main

import ()

func main() {
	if input.String != "" && input.String != x.String {
		x.String = input.String
	}

	if input.PointerToString != nil && *input.PointerToString != "" && input.PointerToString != x.PointerToString {
		x.PointerToString = input.PointerToString
	}

	if input.Bool != x.Bool {
		x.Bool = input.Bool
	}

	if input.PointerToBool != nil && input.PointerToBool != x.PointerToBool {
		x.PointerToBool = input.PointerToBool
	}

	if input.Int != x.Int {
		x.Int = input.Int
	}

	if input.PointerToInt != nil && input.PointerToInt != x.PointerToInt {
		x.PointerToInt = input.PointerToInt
	}

	if input.Int8 != x.Int8 {
		x.Int8 = input.Int8
	}

	if input.PointerToInt8 != nil && input.PointerToInt8 != x.PointerToInt8 {
		x.PointerToInt8 = input.PointerToInt8
	}

	if input.Int16 != x.Int16 {
		x.Int16 = input.Int16
	}

	if input.PointerToInt16 != nil && input.PointerToInt16 != x.PointerToInt16 {
		x.PointerToInt16 = input.PointerToInt16
	}

	if input.Int32 != x.Int32 {
		x.Int32 = input.Int32
	}

	if input.PointerToInt32 != nil && input.PointerToInt32 != x.PointerToInt32 {
		x.PointerToInt32 = input.PointerToInt32
	}

	if input.Int64 != x.Int64 {
		x.Int64 = input.Int64
	}

	if input.PointerToInt64 != nil && input.PointerToInt64 != x.PointerToInt64 {
		x.PointerToInt64 = input.PointerToInt64
	}

	if input.Uint != x.Uint {
		x.Uint = input.Uint
	}

	if input.PointerToUint != nil && input.PointerToUint != x.PointerToUint {
		x.PointerToUint = input.PointerToUint
	}

	if input.Uint8 != x.Uint8 {
		x.Uint8 = input.Uint8
	}

	if input.PointerToUint8 != nil && input.PointerToUint8 != x.PointerToUint8 {
		x.PointerToUint8 = input.PointerToUint8
	}

	if input.Uint16 != x.Uint16 {
		x.Uint16 = input.Uint16
	}

	if input.PointerToUint16 != nil && input.PointerToUint16 != x.PointerToUint16 {
		x.PointerToUint16 = input.PointerToUint16
	}

	if input.Uint32 != x.Uint32 {
		x.Uint32 = input.Uint32
	}

	if input.PointerToUint32 != nil && input.PointerToUint32 != x.PointerToUint32 {
		x.PointerToUint32 = input.PointerToUint32
	}

	if input.Uint64 != x.Uint64 {
		x.Uint64 = input.Uint64
	}

	if input.PointerToUint64 != nil && input.PointerToUint64 != x.PointerToUint64 {
		x.PointerToUint64 = input.PointerToUint64
	}

	if input.Float32 != x.Float32 {
		x.Float32 = input.Float32
	}

	if input.PointerToFloat32 != nil && input.PointerToFloat32 != x.PointerToFloat32 {
		x.PointerToFloat32 = input.PointerToFloat32
	}

	if input.Float64 != x.Float64 {
		x.Float64 = input.Float64
	}

	if input.PointerToFloat64 != nil && input.PointerToFloat64 != x.PointerToFloat64 {
		x.PointerToFloat64 = input.PointerToFloat64
	}
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
