package bleve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_bleveTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := bleveTestDotGo(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"os"
	"testing"
)

type (
	exampleType struct {
		ID            uint64 ` + "`" + `json:"id"` + "`" + `
		Name          string ` + "`" + `json:"name"` + "`" + `
		BelongsToUser uint64 ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	exampleTypeWithStringID struct {
		ID            string ` + "`" + `json:"id"` + "`" + `
		Name          string ` + "`" + `json:"name"` + "`" + `
		BelongsToUser uint64 ` + "`" + `json:"belongsToUser"` + "`" + `
	}
)

func TestNewBleveIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("constructor_test_happy_path.bleve")

		_, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("invalid path", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("")

		_, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.Error(t, err)
	})

	T.Run("invalid name", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("constructor_test_invalid_name.bleve")

		_, err := NewBleveIndexManager(exampleIndexPath, "invalid", noop.ProvideNoopLogger())
		assert.Error(t, err)
	})
}

func TestBleveIndexManager_Index(T *testing.T) {
	T.Parallel()

	exampleUserID := fake.BuildFakeUser().ID

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "index_test"
		exampleIndexPath := search.IndexPath("index_test_obligatory.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}

func TestBleveIndexManager_Search(T *testing.T) {
	T.Parallel()

	exampleUserID := fake.BuildFakeUser().ID

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_obligatory.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		results, err := im.Search(ctx, x.Name, exampleUserID)
		assert.NotEmpty(t, results)
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with empty index and search", func(t *testing.T) {
		ctx := context.Background()

		exampleIndexPath := search.IndexPath("search_test_empty_index.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		results, err := im.Search(ctx, "", exampleUserID)
		assert.Empty(t, results)
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with closed index", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_closed_index.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, im.(*bleveIndexManager).index.Close())

		results, err := im.Search(ctx, x.Name, exampleUserID)
		assert.Empty(t, results)
		assert.Error(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with invalid ID", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_invalid_id.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleTypeWithStringID{
			ID:            "whatever",
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.(*bleveIndexManager).index.Index(x.ID, x))

		results, err := im.Search(ctx, x.Name, exampleUserID)
		assert.Empty(t, results)
		assert.Error(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}

func TestBleveIndexManager_Delete(T *testing.T) {
	T.Parallel()

	exampleUserID := fake.BuildFakeUser().ID

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "delete_test"
		exampleIndexPath := search.IndexPath("delete_test.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, im.Delete(ctx, x.ID))

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBleveTestTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBleveTestTypeDefinitions()

		expected := `
package example

import ()

type (
	exampleType struct {
		ID            uint64 ` + "`" + `json:"id"` + "`" + `
		Name          string ` + "`" + `json:"name"` + "`" + `
		BelongsToUser uint64 ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	exampleTypeWithStringID struct {
		ID            string ` + "`" + `json:"id"` + "`" + `
		Name          string ` + "`" + `json:"name"` + "`" + `
		BelongsToUser uint64 ` + "`" + `json:"belongsToUser"` + "`" + `
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestNewBleveIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestNewBleveIndexManager(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	"os"
	"testing"
)

func TestNewBleveIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("constructor_test_happy_path.bleve")

		_, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("invalid path", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("")

		_, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.Error(t, err)
	})

	T.Run("invalid name", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("constructor_test_invalid_name.bleve")

		_, err := NewBleveIndexManager(exampleIndexPath, "invalid", noop.ProvideNoopLogger())
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBleveIndexManager_Index(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestBleveIndexManager_Index(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"os"
	"testing"
)

func TestBleveIndexManager_Index(T *testing.T) {
	T.Parallel()

	exampleUserID := fake.BuildFakeUser().ID

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "index_test"
		exampleIndexPath := search.IndexPath("index_test_obligatory.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBleveIndexManager_Search(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestBleveIndexManager_Search(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"os"
	"testing"
)

func TestBleveIndexManager_Search(T *testing.T) {
	T.Parallel()

	exampleUserID := fake.BuildFakeUser().ID

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_obligatory.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		results, err := im.Search(ctx, x.Name, exampleUserID)
		assert.NotEmpty(t, results)
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with empty index and search", func(t *testing.T) {
		ctx := context.Background()

		exampleIndexPath := search.IndexPath("search_test_empty_index.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		results, err := im.Search(ctx, "", exampleUserID)
		assert.Empty(t, results)
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with closed index", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_closed_index.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, im.(*bleveIndexManager).index.Close())

		results, err := im.Search(ctx, x.Name, exampleUserID)
		assert.Empty(t, results)
		assert.Error(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with invalid ID", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_invalid_id.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleTypeWithStringID{
			ID:            "whatever",
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.(*bleveIndexManager).index.Index(x.ID, x))

		results, err := im.Search(ctx, x.Name, exampleUserID)
		assert.Empty(t, results)
		assert.Error(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBleveIndexManager_Delete(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestBleveIndexManager_Delete(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"os"
	"testing"
)

func TestBleveIndexManager_Delete(T *testing.T) {
	T.Parallel()

	exampleUserID := fake.BuildFakeUser().ID

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "delete_test"
		exampleIndexPath := search.IndexPath("delete_test.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:            123,
			Name:          exampleQuery,
			BelongsToUser: exampleUserID,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, im.Delete(ctx, x.ID))

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
