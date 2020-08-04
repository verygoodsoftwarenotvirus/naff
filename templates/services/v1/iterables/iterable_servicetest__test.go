package iterables

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterableServiceTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterableServiceTestDotGo(proj, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func buildTestService() *Service {
	return &Service{
		logger:          noop.ProvideNoopLogger(),
		itemCounter:     &mock.UnitCounter{},
		itemDataManager: &mock1.ItemDataManager{},
		itemIDFetcher:   func(req *http.Request) uint64 { return 0 },
		userIDFetcher:   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:  &mock2.EncoderDecoder{},
		reporter:        nil,
		search:          &mock3.IndexManager{},
	}
}

func TestProvideItemsService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock.UnitCounter{}, nil
		}

		s, err := ProvideItemsService(
			noop.ProvideNoopLogger(),
			&mock1.ItemDataManager{},
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
			nil,
			&mock3.IndexManager{},
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	T.Run("with error providing unit counter", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, errors.New("blah")
		}

		s, err := ProvideItemsService(
			noop.ProvideNoopLogger(),
			&mock1.ItemDataManager{},
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
			nil,
			&mock3.IndexManager{},
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildbuildTestServiceFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildbuildTestServiceFuncDecl(proj, typ)

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
)

func buildTestService() *Service {
	return &Service{
		logger:          noop.ProvideNoopLogger(),
		itemCounter:     &mock.UnitCounter{},
		itemDataManager: &mock1.ItemDataManager{},
		itemIDFetcher:   func(req *http.Request) uint64 { return 0 },
		userIDFetcher:   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:  &mock2.EncoderDecoder{},
		reporter:        nil,
		search:          &mock3.IndexManager{},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideServiceFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestProvideServiceFuncDecl(proj, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/mock"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestProvideItemsService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock.UnitCounter{}, nil
		}

		s, err := ProvideItemsService(
			noop.ProvideNoopLogger(),
			&mock1.ItemDataManager{},
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
			nil,
			&mock3.IndexManager{},
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	T.Run("with error providing unit counter", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, errors.New("blah")
		}

		s, err := ProvideItemsService(
			noop.ProvideNoopLogger(),
			&mock1.ItemDataManager{},
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
			nil,
			&mock3.IndexManager{},
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
