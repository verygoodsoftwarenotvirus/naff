package oauth2clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := httpRoutesTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func Test_randString(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		actual := randString()
		assert.NotEmpty(t, actual)
	})
}

func buildRequest(t *testing.T) *http.Request {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodGet,
		"https://verygoodsoftwarenotvirus.ru",
		nil,
	)

	require.NotNil(t, req)
	assert.NoError(t, err)
	return req
}

func Test_fetchUserID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req := buildRequest(t)
		exampleUser := fake.BuildFakeUser()

		// for the service.fetchUserID() call
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		s := buildTestService(t)

		actual := s.fetchUserID(req)
		assert.Equal(t, exampleUser.ID, actual)
	})

	T.Run("without context value present", func(t *testing.T) {
		req := buildRequest(t)

		expected := uint64(0)
		s := buildTestService(t)

		actual := s.fetchUserID(req)
		assert.Equal(t, expected, actual)
	})
}

func TestService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return(exampleOAuth2ClientList, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2ClientList")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		// for the service.fetchUserID() call
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return((*v1.OAuth2ClientList)(nil), sql.ErrNoRows)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2ClientList")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return((*v1.OAuth2ClientList)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return(exampleOAuth2ClientList, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2ClientList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})
}

func TestService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock2.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = a

		uc := &mock3.UnitCounter{}
		uc.On("Increment", mock.Anything).Return()
		s.oauth2ClientCounter = uc

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, uc, ed)
	})

	T.Run("with missing input", func(t *testing.T) {
		s := buildTestService(t)

		req := buildRequest(t)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error getting user", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return((*v1.User)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid credentials", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock2.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = a

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error validating password", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock2.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, errors.New("blah"))
		s.authenticator = a

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error creating oauth2 client", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		a := &mock2.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = a

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock2.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = a

		uc := &mock3.UnitCounter{}
		uc.On("Increment", mock.Anything).Return()
		s.oauth2ClientCounter = uc

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, uc, ed)
	})
}

func TestService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with no rows found", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(exampleOAuth2Client, sql.ErrNoRows)
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching client from database", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})
}

func TestService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"ArchiveOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(nil)
		s.database = mockDB

		uc := &mock3.UnitCounter{}
		uc.On("Decrement", mock.Anything).Return()
		s.oauth2ClientCounter = uc

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, uc)
	})

	T.Run("with no rows found", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"ArchiveOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(sql.ErrNoRows)
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error deleting record", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"ArchiveOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_randString(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTest_randString()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_randString(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		actual := randString()
		assert.NotEmpty(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildRequest()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func buildRequest(t *testing.T) *http.Request {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodGet,
		"https://verygoodsoftwarenotvirus.ru",
		nil,
	)

	require.NotNil(t, req)
	assert.NoError(t, err)
	return req
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_fetchUserID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_fetchUserID(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func Test_fetchUserID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req := buildRequest(t)
		exampleUser := fake.BuildFakeUser()

		// for the service.fetchUserID() call
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		s := buildTestService(t)

		actual := s.fetchUserID(req)
		assert.Equal(t, exampleUser.ID, actual)
	})

	T.Run("without context value present", func(t *testing.T) {
		req := buildRequest(t)

		expected := uint64(0)
		s := buildTestService(t)

		actual := s.fetchUserID(req)
		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_ListHandler(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return(exampleOAuth2ClientList, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2ClientList")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		// for the service.fetchUserID() call
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return((*v11.OAuth2ClientList)(nil), sql.ErrNoRows)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2ClientList")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return((*v11.OAuth2ClientList)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientsForUser",
			mock.Anything,
			exampleUser.ID,
			mock.AnythingOfType("*models.QueryFilter"),
		).Return(exampleOAuth2ClientList, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2ClientList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ListHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_CreateHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock1.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = a

		uc := &mock2.UnitCounter{}
		uc.On("Increment", mock.Anything).Return()
		s.oauth2ClientCounter = uc

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, uc, ed)
	})

	T.Run("with missing input", func(t *testing.T) {
		s := buildTestService(t)

		req := buildRequest(t)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error getting user", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return((*v11.User)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid credentials", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock1.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = a

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error validating password", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock1.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, errors.New("blah"))
		s.authenticator = a

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error creating oauth2 client", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		a := &mock1.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = a

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On(
			"GetUserByUsername",
			mock.Anything,
			exampleInput.Username,
		).Return(exampleUser, nil)
		mockDB.OAuth2ClientDataManager.On(
			"CreateOAuth2Client",
			mock.Anything,
			exampleInput,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		a := &mock1.Authenticator{}
		a.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.Password,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = a

		uc := &mock2.UnitCounter{}
		uc.On("Increment", mock.Anything).Return()
		s.oauth2ClientCounter = uc

		ed := &mock3.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), creationMiddlewareCtxKey, exampleInput),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.CreateHandler(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, uc, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_ReadHandler(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(nil)
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with no rows found", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(exampleOAuth2Client, sql.ErrNoRows)
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching client from database", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.OAuth2Client")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ReadHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_ArchiveHandler(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"ArchiveOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(nil)
		s.database = mockDB

		uc := &mock1.UnitCounter{}
		uc.On("Decrement", mock.Anything).Return()
		s.oauth2ClientCounter = uc

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, uc)
	})

	T.Run("with no rows found", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"ArchiveOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(sql.ErrNoRows)
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error deleting record", func(t *testing.T) {
		s := buildTestService(t)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleOAuth2Client.BelongsToUser = exampleUser.ID

		s.urlClientIDExtractor = func(req *http.Request) uint64 {
			return exampleOAuth2Client.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"ArchiveOAuth2Client",
			mock.Anything,
			exampleOAuth2Client.ID,
			exampleOAuth2Client.BelongsToUser,
		).Return(errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)
		res := httptest.NewRecorder()

		s.ArchiveHandler(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
