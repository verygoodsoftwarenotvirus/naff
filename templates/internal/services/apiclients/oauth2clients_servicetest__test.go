package apiclients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsServiceTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsServiceTestDotGo(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock3 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	manage "gopkg.in/oauth2.v3/manage"
	server "gopkg.in/oauth2.v3/server"
	store "gopkg.in/oauth2.v3/store"
	"net/http"
	"testing"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	manager := manage.NewDefaultManager()
	tokenStore, err := store.NewMemoryTokenStore()
	require.NoError(t, err)
	manager.MustTokenStorage(tokenStore, err)
	server := server.NewDefaultServer(manager)

	service := &Service{
		database:             v1.BuildMockDatabase(),
		logger:               noop.ProvideNoopLogger(),
		encoderDecoder:       &mock.EncoderDecoder{},
		authenticator:        &mock1.Authenticator{},
		urlClientIDExtractor: func(req *http.Request) uint64 { return 0 },
		oauth2ClientCounter:  &mock2.UnitCounter{},
		oauth2Handler:        server,
	}

	return service
}

func TestProvideOAuth2ClientsService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On("GetAllOAuth2Clients", mock3.Anything).Return([]*v11.OAuth2Client{}, nil)

		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, nil
		}

		service, err := ProvideOAuth2ClientsService(
			noop.ProvideNoopLogger(),
			mockDB,
			&mock1.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock.EncoderDecoder{},
			ucp,
		)
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mock3.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error providing counter", func(t *testing.T) {
		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On("GetAllOAuth2Clients", mock3.Anything).Return([]*v11.OAuth2Client{}, nil)

		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, errors.New("blah")
		}

		service, err := ProvideOAuth2ClientsService(
			noop.ProvideNoopLogger(),
			mockDB,
			&mock1.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock.EncoderDecoder{},
			ucp,
		)
		assert.Error(t, err)
		assert.Nil(t, service)

		mock3.AssertExpectationsForObjects(t, mockDB)
	})
}

func Test_clientStore_GetByID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock3.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)

		c := &clientStore{database: mockDB}
		actual, err := c.GetByID(exampleOAuth2Client.ClientID)

		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client.ClientID, actual.GetID())

		mock3.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with no rows", func(t *testing.T) {
		exampleID := "blah"

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock3.Anything,
			exampleID,
		).Return((*v11.OAuth2Client)(nil), sql.ErrNoRows)

		c := &clientStore{database: mockDB}
		_, err := c.GetByID(exampleID)

		assert.Error(t, err)

		mock3.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		exampleID := "blah"

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock3.Anything,
			exampleID,
		).Return((*v11.OAuth2Client)(nil), errors.New(exampleID))

		c := &clientStore{database: mockDB}
		_, err := c.GetByID(exampleID)

		assert.Error(t, err)

		mock3.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_HandleAuthorizeRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		moah := &mockOAuth2Handler{}
		moah.On(
			"HandleAuthorizeRequest",
			mock3.Anything,
			mock3.Anything,
		).Return(nil)
		s.oauth2Handler = moah
		req, res := buildRequest(t), httptest.NewRecorder()

		assert.NoError(t, s.HandleAuthorizeRequest(res, req))

		mock3.AssertExpectationsForObjects(t, moah)
	})
}

func TestService_HandleTokenRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		moah := &mockOAuth2Handler{}
		moah.On(
			"HandleTokenRequest",
			mock3.Anything,
			mock3.Anything,
		).Return(nil)
		s.oauth2Handler = moah
		req, res := buildRequest(t), httptest.NewRecorder()

		assert.NoError(t, s.HandleTokenRequest(res, req))

		mock3.AssertExpectationsForObjects(t, moah)
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
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	manage "gopkg.in/oauth2.v3/manage"
	server "gopkg.in/oauth2.v3/server"
	store "gopkg.in/oauth2.v3/store"
	"net/http"
	"testing"
)

func buildTestService(t *testing.T) *Service {
	t.Helper()

	manager := manage.NewDefaultManager()
	tokenStore, err := store.NewMemoryTokenStore()
	require.NoError(t, err)
	manager.MustTokenStorage(tokenStore, err)
	server := server.NewDefaultServer(manager)

	service := &Service{
		database:             v1.BuildMockDatabase(),
		logger:               noop.ProvideNoopLogger(),
		encoderDecoder:       &mock.EncoderDecoder{},
		authenticator:        &mock1.Authenticator{},
		urlClientIDExtractor: func(req *http.Request) uint64 { return 0 },
		oauth2ClientCounter:  &mock2.UnitCounter{},
		oauth2Handler:        server,
	}

	return service
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideOAuth2ClientsService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestProvideOAuth2ClientsService(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
	"testing"
)

func TestProvideOAuth2ClientsService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On("GetAllOAuth2Clients", mock.Anything).Return([]*v11.OAuth2Client{}, nil)

		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, nil
		}

		service, err := ProvideOAuth2ClientsService(
			noop.ProvideNoopLogger(),
			mockDB,
			&mock1.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
		)
		assert.NoError(t, err)
		assert.NotNil(t, service)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error providing counter", func(t *testing.T) {
		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On("GetAllOAuth2Clients", mock.Anything).Return([]*v11.OAuth2Client{}, nil)

		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, errors.New("blah")
		}

		service, err := ProvideOAuth2ClientsService(
			noop.ProvideNoopLogger(),
			mockDB,
			&mock1.Authenticator{},
			func(req *http.Request) uint64 { return 0 },
			&mock2.EncoderDecoder{},
			ucp,
		)
		assert.Error(t, err)
		assert.Nil(t, service)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_clientStore_GetByID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_clientStore_GetByID(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func Test_clientStore_GetByID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)

		c := &clientStore{database: mockDB}
		actual, err := c.GetByID(exampleOAuth2Client.ClientID)

		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client.ClientID, actual.GetID())

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with no rows", func(t *testing.T) {
		exampleID := "blah"

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleID,
		).Return((*v11.OAuth2Client)(nil), sql.ErrNoRows)

		c := &clientStore{database: mockDB}
		_, err := c.GetByID(exampleID)

		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		exampleID := "blah"

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleID,
		).Return((*v11.OAuth2Client)(nil), errors.New(exampleID))

		c := &clientStore{database: mockDB}
		_, err := c.GetByID(exampleID)

		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_HandleAuthorizeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestService_HandleAuthorizeRequest()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

func TestService_HandleAuthorizeRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		moah := &mockOAuth2Handler{}
		moah.On(
			"HandleAuthorizeRequest",
			mock.Anything,
			mock.Anything,
		).Return(nil)
		s.oauth2Handler = moah
		req, res := buildRequest(t), httptest.NewRecorder()

		assert.NoError(t, s.HandleAuthorizeRequest(res, req))

		mock.AssertExpectationsForObjects(t, moah)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_HandleTokenRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestService_HandleTokenRequest()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

func TestService_HandleTokenRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		moah := &mockOAuth2Handler{}
		moah.On(
			"HandleTokenRequest",
			mock.Anything,
			mock.Anything,
		).Return(nil)
		s.oauth2Handler = moah
		req, res := buildRequest(t), httptest.NewRecorder()

		assert.NoError(t, s.HandleTokenRequest(res, req))

		mock.AssertExpectationsForObjects(t, moah)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
