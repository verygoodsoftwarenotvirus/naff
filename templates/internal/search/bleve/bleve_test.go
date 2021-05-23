package bleve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_bleveDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := bleveDotGo(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	bleve "github.com/blevesearch/bleve"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"strconv"
)

const (
	base    = 10
	bitSize = 64

	// testingSearchIndexName is an index name that is only valid for testing's sake.
	testingSearchIndexName search.IndexName = "testing"
)

var _ search.IndexManager = (*bleveIndexManager)(nil)

type (
	bleveIndexManager struct {
		index  bleve.Index
		logger v1.Logger
	}
)

// NewBleveIndexManager instantiates a bleve index
func NewBleveIndexManager(path search.IndexPath, name search.IndexName, logger v1.Logger) (search.IndexManager, error) {
	var index bleve.Index

	preexistingIndex, openIndexErr := bleve.Open(string(path))
	switch openIndexErr {
	case nil:
		index = preexistingIndex
	case bleve.ErrorIndexPathDoesNotExist:
		logger.WithValue("path", path).Debug("tried to open existing index, but didn't find it")
		var newIndexErr error

		switch name {
		case testingSearchIndexName:
			index, newIndexErr = bleve.New(string(path), bleve.NewIndexMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		case v11.ItemsSearchIndexName:
			index, newIndexErr = bleve.New(string(path), buildItemMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		default:
			return nil, fmt.Errorf("invalid index name: %q", name)
		}
	default:
		logger.Error(openIndexErr, "failed to open index")
		return nil, openIndexErr
	}

	im := &bleveIndexManager{
		index:  index,
		logger: logger.WithName(fmt.Sprintf("%s_search", name)),
	}

	return im, nil
}

// Index implements our IndexManager interface
func (sm *bleveIndexManager) Index(ctx context.Context, id uint64, value interface{}) error {
	_, span := tracing.StartSpan(ctx, "Index")
	defer span.End()

	sm.logger.WithValue("id", id).Debug("adding to index")
	return sm.index.Index(strconv.FormatUint(id, base), value)
}

// Search implements our IndexManager interface
func (sm *bleveIndexManager) Search(ctx context.Context, query string, userID uint64) (ids []uint64, err error) {
	_, span := tracing.StartSpan(ctx, "Search")
	defer span.End()

	query = ensureQueryIsRestrictedToUser(query, userID)
	tracing.AttachSearchQueryToSpan(span, query)
	sm.logger.WithValues(map[string]interface{}{
		"search_query": query,
		"user_id":      userID,
	}).Debug("performing search")

	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResults, err := sm.index.SearchInContext(ctx, searchRequest)
	if err != nil {
		sm.logger.Error(err, "performing search query")
		return nil, err
	}

	out := []uint64{}
	for _, result := range searchResults.Hits {
		x, err := strconv.ParseUint(result.ID, base, bitSize)
		if err != nil {
			// this should literally never happen
			return nil, err
		}
		out = append(out, x)
	}

	return out, nil
}

// Delete implements our IndexManager interface
func (sm *bleveIndexManager) Delete(ctx context.Context, id uint64) error {
	_, span := tracing.StartSpan(ctx, "Delete")
	defer span.End()

	sm.logger.WithValue("id", id).Debug("removing from index")
	return sm.index.Delete(strconv.FormatUint(id, base))
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConstantDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildConstantDefinitions(proj)

		expected := `
package example

import (
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
)

const (
	base    = 10
	bitSize = 64

	// testingSearchIndexName is an index name that is only valid for testing's sake.
	testingSearchIndexName search.IndexName = "testing"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfaceImplementationStatement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildInterfaceImplementationStatement(proj)

		expected := `
package example

import (
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
)

var _ search.IndexManager = (*bleveIndexManager)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTypeDefinitions()

		expected := `
package example

import (
	bleve "github.com/blevesearch/bleve"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

type (
	bleveIndexManager struct {
		index  bleve.Index
		logger v1.Logger
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewBleveIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildNewBleveIndexManager(proj)

		expected := `
package example

import (
	"fmt"
	bleve "github.com/blevesearch/bleve"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// NewBleveIndexManager instantiates a bleve index
func NewBleveIndexManager(path search.IndexPath, name search.IndexName, logger v1.Logger) (search.IndexManager, error) {
	var index bleve.Index

	preexistingIndex, openIndexErr := bleve.Open(string(path))
	switch openIndexErr {
	case nil:
		index = preexistingIndex
	case bleve.ErrorIndexPathDoesNotExist:
		logger.WithValue("path", path).Debug("tried to open existing index, but didn't find it")
		var newIndexErr error

		switch name {
		case testingSearchIndexName:
			index, newIndexErr = bleve.New(string(path), bleve.NewIndexMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		case v11.ItemsSearchIndexName:
			index, newIndexErr = bleve.New(string(path), buildItemMapping())
			if newIndexErr != nil {
				logger.Error(newIndexErr, "failed to create new index")
				return nil, newIndexErr
			}
		default:
			return nil, fmt.Errorf("invalid index name: %q", name)
		}
	default:
		logger.Error(openIndexErr, "failed to open index")
		return nil, openIndexErr
	}

	im := &bleveIndexManager{
		index:  index,
		logger: logger.WithName(fmt.Sprintf("%s_search", name)),
	}

	return im, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewBleveIndexManager_Index(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildNewBleveIndexManager_Index(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"strconv"
)

// Index implements our IndexManager interface
func (sm *bleveIndexManager) Index(ctx context.Context, id uint64, value interface{}) error {
	_, span := tracing.StartSpan(ctx, "Index")
	defer span.End()

	sm.logger.WithValue("id", id).Debug("adding to index")
	return sm.index.Index(strconv.FormatUint(id, base), value)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewBleveIndexManager_Search(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildNewBleveIndexManager_Search(proj)

		expected := `
package example

import (
	"context"
	bleve "github.com/blevesearch/bleve"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"strconv"
)

// Search implements our IndexManager interface
func (sm *bleveIndexManager) Search(ctx context.Context, query string, userID uint64) (ids []uint64, err error) {
	_, span := tracing.StartSpan(ctx, "Search")
	defer span.End()

	query = ensureQueryIsRestrictedToUser(query, userID)
	tracing.AttachSearchQueryToSpan(span, query)
	sm.logger.WithValues(map[string]interface{}{
		"search_query": query,
		"user_id":      userID,
	}).Debug("performing search")

	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResults, err := sm.index.SearchInContext(ctx, searchRequest)
	if err != nil {
		sm.logger.Error(err, "performing search query")
		return nil, err
	}

	out := []uint64{}
	for _, result := range searchResults.Hits {
		x, err := strconv.ParseUint(result.ID, base, bitSize)
		if err != nil {
			// this should literally never happen
			return nil, err
		}
		out = append(out, x)
	}

	return out, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewBleveIndexManager_Delete(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildNewBleveIndexManager_Delete(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"strconv"
)

// Delete implements our IndexManager interface
func (sm *bleveIndexManager) Delete(ctx context.Context, id uint64) error {
	_, span := tracing.StartSpan(ctx, "Delete")
	defer span.End()

	sm.logger.WithValue("id", id).Debug("removing from index")
	return sm.index.Delete(strconv.FormatUint(id, base))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
