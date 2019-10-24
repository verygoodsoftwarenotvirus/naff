package users

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
	mockauth "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock"
	mockencoding "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock"
	"gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics"
	mockmetrics "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	expectedUserCount := uint64(123)
	mockDB := database.BuildMockDatabase()
	mockDB.UserDataManager.On("GetUserCount", mock.Anything, (*models.QueryFilter)(nil)).Return(expectedUserCount, nil)

	uc := &mockmetrics.UnitCounter{}
	uc.On("IncrementBy", mock.Anything)
	var ucp metrics.UnitCounterProvider = func(
		counterName metrics.CounterName,
		description string,
	) (metrics.UnitCounter, error) {
		return uc, nil
	}

	service, err := ProvideUsersService(
		context.Background(),
		config.AuthSettings{},
		noop.ProvideNoopLogger(),
		mockDB,
		&mockauth.Authenticator{},
		func(req *http.Request) uint64 { return 0 },
		&mockencoding.EncoderDecoder{},
		ucp,
		newsman.NewNewsman(nil, nil),
	)
	require.NoError(t, err)

	return service
}

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		mockUserCount := uint64(0)
		mockDB := database.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserCount", mock.Anything, mock.Anything).Return(mockUserCount, nil)

		uc := &mockmetrics.UnitCounter{}
		uc.On("IncrementBy", mockUserCount).Return()

		var ucp metrics.UnitCounterProvider = func(
			counterName metrics.CounterName,
			description string,
		) (metrics.UnitCounter, error) {
			return uc, nil
		}

		service, err := ProvideUsersService(
			context.Background(),
			config.AuthSettings{},
			noop.ProvideNoopLogger(),
			mockDB,
			&mockauth.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mockencoding.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	T.Run("with nil userIDFetcher", func(t *testing.T) {
		mockUserCount := uint64(0)
		mockDB := database.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserCount", mock.Anything, mock.Anything).Return(mockUserCount, nil)

		uc := &mockmetrics.UnitCounter{}
		uc.On("IncrementBy", mockUserCount).Return()

		var ucp metrics.UnitCounterProvider = func(
			counterName metrics.CounterName,
			description string,
		) (metrics.UnitCounter, error) {
			return uc, nil
		}
		service, err := ProvideUsersService(
			context.Background(),
			config.AuthSettings{},
			noop.ProvideNoopLogger(),
			mockDB,
			&mockauth.Authenticator{},
			nil,
			&mockencoding.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	T.Run("with error initializing counter", func(t *testing.T) {
		mockUserCount := uint64(0)
		mockDB := database.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserCount", mock.Anything, mock.Anything).Return(mockUserCount, nil)

		uc := &mockmetrics.UnitCounter{}
		uc.On("IncrementBy", mockUserCount).Return()

		var ucp metrics.UnitCounterProvider = func(
			counterName metrics.CounterName,
			description string,
		) (metrics.UnitCounter, error) {
			return uc, errors.New("blah")
		}

		service, err := ProvideUsersService(
			context.Background(),
			config.AuthSettings{},
			noop.ProvideNoopLogger(),
			mockDB,
			&mockauth.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mockencoding.EncoderDecoder{},
			ucp,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	T.Run("with error getting user count", func(t *testing.T) {
		mockUserCount := uint64(0)
		mockDB := database.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserCount", mock.Anything, mock.Anything).Return(mockUserCount, errors.New("blah"))

		uc := &mockmetrics.UnitCounter{}
		var ucp metrics.UnitCounterProvider = func(
			counterName metrics.CounterName,
			description string,
		) (metrics.UnitCounter, error) {
			return uc, nil
		}

		service, err := ProvideUsersService(
			context.Background(),
			config.AuthSettings{},
			noop.ProvideNoopLogger(),
			mockDB,
			&mockauth.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mockencoding.EncoderDecoder{},
			ucp,
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, service)
	})
}