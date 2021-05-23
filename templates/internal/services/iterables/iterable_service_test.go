package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterableServiceDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterableServiceDotGo(proj, typ)

		expected := `
package example

import (
	"fmt"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to item input data in contexts.
	createMiddlewareCtxKey v1.ContextKey = "item_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to item update data in contexts.
	updateMiddlewareCtxKey v1.ContextKey = "item_update_input"

	counterName        metrics.CounterName = "items"
	counterDescription string              = "the number of items managed by the items service"
	topicName          string              = "items"
	serviceName        string              = "items_service"
)

var (
	_ v1.ItemDataServer = (*Service)(nil)
)

type (
	// SearchIndex is a type alias for dependency injection's sake
	SearchIndex search.IndexManager

	// Service handles to-do list items
	Service struct {
		logger          v11.Logger
		itemDataManager v1.ItemDataManager
		itemIDFetcher   ItemIDFetcher
		userIDFetcher   UserIDFetcher
		itemCounter     metrics.UnitCounter
		encoderDecoder  encoding.EncoderDecoder
		reporter        newsman.Reporter
		search          SearchIndex
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// ItemIDFetcher is a function that fetches item IDs.
	ItemIDFetcher func(*http.Request) uint64
)

// ProvideItemsService builds a new ItemsService.
func ProvideItemsService(
	logger v11.Logger,
	itemDataManager v1.ItemDataManager,
	itemIDFetcher ItemIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	itemCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
	searchIndexManager SearchIndex,
) (*Service, error) {
	itemCounter, err := itemCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:          logger.WithName(serviceName),
		itemIDFetcher:   itemIDFetcher,
		userIDFetcher:   userIDFetcher,
		itemDataManager: itemDataManager,
		encoderDecoder:  encoder,
		itemCounter:     itemCounter,
		reporter:        reporter,
		search:          searchIndexManager,
	}

	return svc, nil
}

// ProvideItemsServiceSearchIndex provides a search index for the service
func ProvideItemsServiceSearchIndex(
	searchSettings config.SearchSettings,
	indexProvider search.IndexManagerProvider,
	logger v11.Logger,
) (SearchIndex, error) {
	logger.WithValue("index_path", searchSettings.ItemsIndexPath).Debug("setting up items search index")

	searchIndex, indexInitErr := indexProvider(searchSettings.ItemsIndexPath, v1.ItemsSearchIndexName, logger)
	if indexInitErr != nil {
		logger.Error(indexInitErr, "setting up items search index")
		return nil, indexInitErr
	}

	return searchIndex, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingServiceConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingServiceConstantDefs(proj, typ)

		expected := `
package example

import (
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to item input data in contexts.
	createMiddlewareCtxKey v1.ContextKey = "item_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to item update data in contexts.
	updateMiddlewareCtxKey v1.ContextKey = "item_update_input"

	counterName        metrics.CounterName = "items"
	counterDescription string              = "the number of items managed by the items service"
	topicName          string              = "items"
	serviceName        string              = "items_service"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingServiceVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingServiceVarDefs(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

var (
	_ v1.ItemDataServer = (*Service)(nil)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServiceTypeDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildServiceTypeDecls(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

type (
	// SearchIndex is a type alias for dependency injection's sake
	SearchIndex search.IndexManager

	// Service handles to-do list items
	Service struct {
		logger          v1.Logger
		itemDataManager v11.ItemDataManager
		itemIDFetcher   ItemIDFetcher
		userIDFetcher   UserIDFetcher
		itemCounter     metrics.UnitCounter
		encoderDecoder  encoding.EncoderDecoder
		reporter        newsman.Reporter
		search          SearchIndex
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// ItemIDFetcher is a function that fetches item IDs.
	ItemIDFetcher func(*http.Request) uint64
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildServiceTypeDecls(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

type (
	// Service handles to-do list yet another things
	Service struct {
		logger                     v1.Logger
		thingDataManager           v11.ThingDataManager
		anotherThingDataManager    v11.AnotherThingDataManager
		yetAnotherThingDataManager v11.YetAnotherThingDataManager
		thingIDFetcher             ThingIDFetcher
		anotherThingIDFetcher      AnotherThingIDFetcher
		yetAnotherThingIDFetcher   YetAnotherThingIDFetcher
		yetAnotherThingCounter     metrics.UnitCounter
		encoderDecoder             encoding.EncoderDecoder
		reporter                   newsman.Reporter
	}

	// ThingIDFetcher is a function that fetches thing IDs.
	ThingIDFetcher func(*http.Request) uint64

	// AnotherThingIDFetcher is a function that fetches another thing IDs.
	AnotherThingIDFetcher func(*http.Request) uint64

	// YetAnotherThingIDFetcher is a function that fetches yet another thing IDs.
	YetAnotherThingIDFetcher func(*http.Request) uint64
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideServiceFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildProvideServiceFuncDecl(proj, typ)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideItemsService builds a new ItemsService.
func ProvideItemsService(
	logger v1.Logger,
	itemDataManager v11.ItemDataManager,
	itemIDFetcher ItemIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	itemCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
	searchIndexManager SearchIndex,
) (*Service, error) {
	itemCounter, err := itemCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:          logger.WithName(serviceName),
		itemIDFetcher:   itemIDFetcher,
		userIDFetcher:   userIDFetcher,
		itemDataManager: itemDataManager,
		encoderDecoder:  encoder,
		itemCounter:     itemCounter,
		reporter:        reporter,
		search:          searchIndexManager,
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildProvideServiceFuncDecl(proj, typ)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideYetAnotherThingsService builds a new YetAnotherThingsService.
func ProvideYetAnotherThingsService(
	logger v1.Logger,
	thingDataManager v11.ThingDataManager,
	anotherThingDataManager v11.AnotherThingDataManager,
	yetAnotherThingDataManager v11.YetAnotherThingDataManager,
	thingIDFetcher ThingIDFetcher,
	anotherThingIDFetcher AnotherThingIDFetcher,
	yetAnotherThingIDFetcher YetAnotherThingIDFetcher,
	encoder encoding.EncoderDecoder,
	yetAnotherThingCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	yetAnotherThingCounter, err := yetAnotherThingCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:                     logger.WithName(serviceName),
		thingIDFetcher:             thingIDFetcher,
		anotherThingIDFetcher:      anotherThingIDFetcher,
		yetAnotherThingIDFetcher:   yetAnotherThingIDFetcher,
		thingDataManager:           thingDataManager,
		anotherThingDataManager:    anotherThingDataManager,
		yetAnotherThingDataManager: yetAnotherThingDataManager,
		encoderDecoder:             encoder,
		yetAnotherThingCounter:     yetAnotherThingCounter,
		reporter:                   reporter,
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideServiceSearchIndexFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildProvideServiceSearchIndexFuncDecl(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// ProvideItemsServiceSearchIndex provides a search index for the service
func ProvideItemsServiceSearchIndex(
	searchSettings config.SearchSettings,
	indexProvider search.IndexManagerProvider,
	logger v1.Logger,
) (SearchIndex, error) {
	logger.WithValue("index_path", searchSettings.ItemsIndexPath).Debug("setting up items search index")

	searchIndex, indexInitErr := indexProvider(searchSettings.ItemsIndexPath, v11.ItemsSearchIndexName, logger)
	if indexInitErr != nil {
		logger.Error(indexInitErr, "setting up items search index")
		return nil, indexInitErr
	}

	return searchIndex, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
