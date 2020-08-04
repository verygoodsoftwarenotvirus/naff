package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterableDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
		BelongsToUser uint64  ` + "`" + `json:"belongsToUser"` + "`" + `
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
		BelongsToUser uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// ItemUpdateInput represents what a user could set as input for updating items.
	ItemUpdateInput struct {
		Name          string ` + "`" + `json:"name"` + "`" + `
		Details       string ` + "`" + `json:"details"` + "`" + `
		BelongsToUser uint64 ` + "`" + `json:"-"` + "`" + `
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
}

func Test_buildUpdateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
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
		t.Parallel()

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
		BelongsToUser uint64  ` + "`" + `json:"belongsToUser"` + "`" + `
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
		BelongsToUser uint64 ` + "`" + `json:"-"` + "`" + `
	}

	// ItemUpdateInput represents what a user could set as input for updating items.
	ItemUpdateInput struct {
		Name          string ` + "`" + `json:"name"` + "`" + `
		Details       string ` + "`" + `json:"details"` + "`" + `
		BelongsToUser uint64 ` + "`" + `json:"-"` + "`" + `
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
}

func Test_buildSomethingToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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

func Test_buildBaseModelStructFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	BelongsToUser uint64  ` + "`" + `json:"belongsToUser"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateModelStructFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateModelStructFields(typ)

		expected := `
package example

import ()

type Example struct {
	Name          string ` + "`" + `json:"name"` + "`" + `
	Details       string ` + "`" + `json:"details"` + "`" + `
	BelongsToUser uint64 ` + "`" + `json:"-"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateModelStructFields(typ)

		expected := `
package example

import ()

type Example struct {
	Name          string ` + "`" + `json:"name"` + "`" + `
	Details       string ` + "`" + `json:"details"` + "`" + `
	BelongsToUser uint64 ` + "`" + `json:"-"` + "`" + `
}
`
		actual := testutils.RenderStructFieldsToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfaceMethods(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
		t.Parallel()

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
