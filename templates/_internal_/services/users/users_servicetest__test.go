package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersServiceTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := usersServiceTestDotGo(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/authentication/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
	"testing"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	expectedUserCount := uint64(123)

	mockDB := v1.BuildMockDatabase()
	mockDB.UserDataManager.On("GetAllUsersCount", mock.Anything).Return(expectedUserCount, nil)

	uc := &mock1.UnitCounter{}
	var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
		return uc, nil
	}

	service, err := ProvideUsersService(
		config.AuthSettings{},
		logging.NewNonOperationalLogger(),
		v1.BuildMockDatabase(),
		&mock2.Authenticator{},
		func(req *http.Request) uint64 { return 0 },
		&mock3.EncoderDecoder{},
		ucp,
		newsman.NewNewsman(nil, nil),
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mockDB, uc)

	return service
}

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock1.UnitCounter{}, nil
		}

		service, err := ProvideUsersService(
			config.AuthSettings{},
			logging.NewNonOperationalLogger(),
			v1.BuildMockDatabase(),
			&mock2.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock3.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	T.Run("with nil userIDFetcher", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock1.UnitCounter{}, nil
		}

		service, err := ProvideUsersService(
			config.AuthSettings{},
			logging.NewNonOperationalLogger(),
			v1.BuildMockDatabase(),
			&mock2.Authenticator{},
			nil,
			&mock3.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	T.Run("with error initializing counter", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock1.UnitCounter{}, errors.New("blah")
		}

		service, err := ProvideUsersService(
			config.AuthSettings{},
			logging.NewNonOperationalLogger(),
			v1.BuildMockDatabase(),
			&mock2.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock3.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, service)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildTestService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildTestService(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/authentication/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
	"testing"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	expectedUserCount := uint64(123)

	mockDB := v1.BuildMockDatabase()
	mockDB.UserDataManager.On("GetAllUsersCount", mock.Anything).Return(expectedUserCount, nil)

	uc := &mock1.UnitCounter{}
	var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
		return uc, nil
	}

	service, err := ProvideUsersService(
		config.AuthSettings{},
		logging.NewNonOperationalLogger(),
		v1.BuildMockDatabase(),
		&mock2.Authenticator{},
		func(req *http.Request) uint64 { return 0 },
		&mock3.EncoderDecoder{},
		ucp,
		newsman.NewNewsman(nil, nil),
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, mockDB, uc)

	return service
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestProvideUsersService(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/authentication/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	"net/http"
	"testing"
)

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock.UnitCounter{}, nil
		}

		service, err := ProvideUsersService(
			config.AuthSettings{},
			logging.NewNonOperationalLogger(),
			v1.BuildMockDatabase(),
			&mock1.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	T.Run("with nil userIDFetcher", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock.UnitCounter{}, nil
		}

		service, err := ProvideUsersService(
			config.AuthSettings{},
			logging.NewNonOperationalLogger(),
			v1.BuildMockDatabase(),
			&mock1.Authenticator{},
			nil,
			&mock2.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	T.Run("with error initializing counter", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mock.UnitCounter{}, errors.New("blah")
		}

		service, err := ProvideUsersService(
			config.AuthSettings{},
			logging.NewNonOperationalLogger(),
			v1.BuildMockDatabase(),
			&mock1.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, service)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
