package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1/client"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"strings"
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

func TestService_validateCredentialChangeRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			examplePassword,
			exampleUser.TwoFactorSecret,
			exampleTOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Equal(t, exampleUser, actual)
		assert.Equal(t, http.StatusOK, sc)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with no rows found in database", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), sql.ErrNoRows)
		s.userDataManager = mockDB

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusNotFound, sc)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), errors.New("blah"))
		s.userDataManager = mockDB

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusInternalServerError, sc)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error validating login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			examplePassword,
			exampleUser.TwoFactorSecret,
			exampleTOTPToken,
			exampleUser.Salt,
		).Return(false, errors.New("blah"))
		s.authenticator = auth

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusInternalServerError, sc)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with invalid login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			examplePassword,
			exampleUser.TwoFactorSecret,
			exampleTOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = auth

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusUnauthorized, sc)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})
}

func TestService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUserList := fake.BuildFakeUserList()

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, mock.Anything).Return(exampleUserList, nil)
		s.userDataManager = mockDB

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.UserList")).Return(nil)
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, mock.Anything).Return((*v11.UserList)(nil), errors.New("blah"))
		s.userDataManager = mockDB

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ListHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUserList := fake.BuildFakeUserList()

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, mock.Anything).Return(exampleUserList, nil)
		s.userDataManager = mockDB

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.UserList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})
}

func TestService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		mc := &mock3.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.userCounter = mc

		r := &mock4.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.UserCreationResponse")).Return(nil)
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, auth, db, mc, r, ed)
	})

	T.Run("with user creation disabled", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.userCreationEnabled = false
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	T.Run("with missing input", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error hashing password", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, errors.New("blah"))
		s.authenticator = auth

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, auth)
	})

	T.Run("with error generating two factor secret", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		sg := &mockSecretGenerator{}
		sg.On("GenerateTwoFactorSecret").Return("", errors.New("blah"))
		s.secretGenerator = sg

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, auth, db, sg)
	})

	T.Run("with error generating salt", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		sg := &mockSecretGenerator{}
		sg.On("GenerateTwoFactorSecret").Return("PRETENDTHISISASECRET", nil)
		sg.On("GenerateSalt").Return([]byte{}, errors.New("blah"))
		s.secretGenerator = sg

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, auth, db, sg)
	})

	T.Run("with error creating entry in database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, errors.New("blah"))
		s.userDataManager = db

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, auth, db)
	})

	T.Run("with pre-existing entry in database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, client.ErrUserExists)
		s.userDataManager = db

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, auth, db)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock1.Authenticator{}
		auth.On("HashPassword", mock.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock.Anything, mock.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		mc := &mock3.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.userCounter = mc

		r := &mock4.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.UserCreationResponse")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, auth, db, mc, r, ed)
	})
}

func TestService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with no rows found", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, sql.ErrNoRows)
		s.userDataManager = mockDB

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, errors.New("blah"))
		s.userDataManager = mockDB

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})
}

func TestService_NewTOTPSecretHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.TOTPSecretRefreshResponse")).Return(nil)
		s.encoderDecoder = ed

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth, ed)
	})

	T.Run("without input attached to request", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with input attached but without user information", func(t *testing.T) {
		s := buildTestService(t)

		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	T.Run("with error validating login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(false, errors.New("blah"))
		s.authenticator = auth

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error generating secret", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		sg := &mockSecretGenerator{}
		sg.On("GenerateTwoFactorSecret").Return("", errors.New("blah"))
		s.secretGenerator = sg

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth, sg)
	})

	T.Run("with error updating user in database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("blah"))
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.TOTPSecretRefreshResponse")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth, ed)
	})
}

func TestService_TOTPSecretValidationHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("VerifyUserTwoFactorSecret", mock.Anything, exampleUser.ID).Return(nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without valid input attached", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil

		res, req := httptest.NewRecorder(), buildRequest(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), errors.New("blah"))
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with secret already validated", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		og := exampleUser.TwoFactorSecretVerifiedOn
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		exampleUser.TwoFactorSecretVerifiedOn = og

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusAlreadyReported, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid code", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)
		exampleInput.TOTPToken = "INVALID"

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error verifying two factor secret", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("VerifyUserTwoFactorSecret", mock.Anything, exampleUser.ID).Return(errors.New("blah"))
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

}

func TestService_UpdatePasswordHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		auth.On("HashPassword", mock.Anything, exampleInput.NewPassword).Return("blah", nil)
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("without input attached to request", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with input but without user info", func(t *testing.T) {
		s := buildTestService(t)

		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	T.Run("with error validating login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(false, errors.New("blah"))
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error hashing password", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		auth.On("HashPassword", mock.Anything, exampleInput.NewPassword).Return("blah", errors.New("blah"))
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error updating user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v11.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(errors.New("blah"))
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		auth.On("HashPassword", mock.Anything, exampleInput.NewPassword).Return("blah", nil)
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})
}

func TestService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}
		res, req := httptest.NewRecorder(), buildRequest(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("ArchiveUser", mock.Anything, exampleUser.ID).Return(nil)
		s.userDataManager = mockDB

		r := &mock4.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mock3.UnitCounter{}
		mc.On("Decrement", mock.Anything)
		s.userCounter = mc

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, r, mc)
	})

	T.Run("with error updating database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}
		res, req := httptest.NewRecorder(), buildRequest(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("ArchiveUser", mock.Anything, exampleUser.ID).Return(errors.New("blah"))
		s.userDataManager = mockDB

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildQRCode(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		actual := s.buildQRCode(ctx, exampleUser.Username, exampleUser.TwoFactorSecret)

		assert.NotEmpty(t, actual)
		assert.True(t, strings.HasPrefix(actual, base64ImagePrefix))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

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

func Test_buildTestService_validateCredentialChangeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_validateCredentialChangeRequest(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_validateCredentialChangeRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			examplePassword,
			exampleUser.TwoFactorSecret,
			exampleTOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Equal(t, exampleUser, actual)
		assert.Equal(t, http.StatusOK, sc)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with no rows found in database", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), sql.ErrNoRows)
		s.userDataManager = mockDB

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusNotFound, sc)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), errors.New("blah"))
		s.userDataManager = mockDB

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusInternalServerError, sc)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error validating login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			examplePassword,
			exampleUser.TwoFactorSecret,
			exampleTOTPToken,
			exampleUser.Salt,
		).Return(false, errors.New("blah"))
		s.authenticator = auth

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusInternalServerError, sc)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with invalid login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleTOTPToken := "123456"
		examplePassword := "password"

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			examplePassword,
			exampleUser.TwoFactorSecret,
			exampleTOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = auth

		actual, sc := s.validateCredentialChangeRequest(
			ctx,
			exampleUser.ID,
			examplePassword,
			exampleTOTPToken,
		)

		assert.Nil(t, actual)
		assert.Equal(t, http.StatusUnauthorized, sc)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
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
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_ListHandler(proj)

		expected := `
package example

import (
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

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUserList := fake.BuildFakeUserList()

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, mock.Anything).Return(exampleUserList, nil)
		s.userDataManager = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.UserList")).Return(nil)
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, mock.Anything).Return((*v11.UserList)(nil), errors.New("blah"))
		s.userDataManager = mockDB

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ListHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUserList := fake.BuildFakeUserList()

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, mock.Anything).Return(exampleUserList, nil)
		s.userDataManager = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.UserList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
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
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_CreateHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1/client"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"testing"
)

func TestService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock1.Anything, mock1.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.userCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.UserCreationResponse")).Return(nil)
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock1.AssertExpectationsForObjects(t, auth, db, mc, r, ed)
	})

	T.Run("with user creation disabled", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.userCreationEnabled = false
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	T.Run("with missing input", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error hashing password", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, errors.New("blah"))
		s.authenticator = auth

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, auth)
	})

	T.Run("with error generating two factor secret", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock1.Anything, mock1.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		sg := &mockSecretGenerator{}
		sg.On("GenerateTwoFactorSecret").Return("", errors.New("blah"))
		s.secretGenerator = sg

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, auth, db, sg)
	})

	T.Run("with error generating salt", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock1.Anything, mock1.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		sg := &mockSecretGenerator{}
		sg.On("GenerateTwoFactorSecret").Return("PRETENDTHISISASECRET", nil)
		sg.On("GenerateSalt").Return([]byte{}, errors.New("blah"))
		s.secretGenerator = sg

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, auth, db, sg)
	})

	T.Run("with error creating entry in database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock1.Anything, mock1.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, errors.New("blah"))
		s.userDataManager = db

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, auth, db)
	})

	T.Run("with pre-existing entry in database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock1.Anything, mock1.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, client.ErrUserExists)
		s.userDataManager = db

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock1.AssertExpectationsForObjects(t, auth, db)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		auth := &mock.Authenticator{}
		auth.On("HashPassword", mock1.Anything, exampleInput.Password).Return(exampleUser.HashedPassword, nil)
		s.authenticator = auth

		db := v1.BuildMockDatabase()
		db.UserDataManager.On("CreateUser", mock1.Anything, mock1.AnythingOfType("models.UserDatabaseCreationInput")).Return(exampleUser, nil)
		s.userDataManager = db

		mc := &mock2.UnitCounter{}
		mc.On("Increment", mock1.Anything)
		s.userCounter = mc

		r := &mock3.Reporter{}
		r.On("Report", mock1.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mock4.EncoderDecoder{}
		ed.On("EncodeResponse", mock1.Anything, mock1.AnythingOfType("*models.UserCreationResponse")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				userCreationMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.userCreationEnabled = true
		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock1.AssertExpectationsForObjects(t, auth, db, mc, r, ed)
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
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_ReadHandler(proj)

		expected := `
package example

import (
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, ed)
	})

	T.Run("with no rows found", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, sql.ErrNoRows)
		s.userDataManager = mockDB

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, errors.New("blah"))
		s.userDataManager = mockDB

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(_ *http.Request) uint64 {
			return exampleUser.ID
		}

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		ed := &mock1.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res, req := httptest.NewRecorder(), buildRequest(t)
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

func Test_buildTestService_NewTOTPSecretHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_NewTOTPSecretHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_NewTOTPSecretHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.TOTPSecretRefreshResponse")).Return(nil)
		s.encoderDecoder = ed

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth, ed)
	})

	T.Run("without input attached to request", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with input attached but without user information", func(t *testing.T) {
		s := buildTestService(t)

		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	T.Run("with error validating login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(false, errors.New("blah"))
		s.authenticator = auth

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error generating secret", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		sg := &mockSecretGenerator{}
		sg.On("GenerateTwoFactorSecret").Return("", errors.New("blah"))
		s.secretGenerator = sg

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth, sg)
	})

	T.Run("with error updating user in database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("blah"))
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretRefreshInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretRefreshMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = auth

		ed := &mock2.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.TOTPSecretRefreshResponse")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		s.NewTOTPSecretHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_TOTPSecretValidationHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_TOTPSecretValidationHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_TOTPSecretValidationHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("VerifyUserTwoFactorSecret", mock.Anything, exampleUser.ID).Return(nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without valid input attached", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil

		res, req := httptest.NewRecorder(), buildRequest(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), errors.New("blah"))
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with secret already validated", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		og := exampleUser.TwoFactorSecretVerifiedOn
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		exampleUser.TwoFactorSecretVerifiedOn = og

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusAlreadyReported, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid code", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)
		exampleInput.TOTPToken = "INVALID"

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error verifying two factor secret", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedOn = nil
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				totpSecretVerificationMiddlewareCtxKey,
				exampleInput,
			),
		)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("VerifyUserTwoFactorSecret", mock.Anything, exampleUser.ID).Return(errors.New("blah"))
		s.userDataManager = mockDB

		s.TOTPSecretVerificationHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_UpdatePasswordHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_UpdatePasswordHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_UpdatePasswordHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		auth.On("HashPassword", mock.Anything, exampleInput.NewPassword).Return("blah", nil)
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("without input attached to request", func(t *testing.T) {
		s := buildTestService(t)

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with input but without user info", func(t *testing.T) {
		s := buildTestService(t)

		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	T.Run("with error validating login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(false, errors.New("blah"))
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error hashing password", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(nil)
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		auth.On("HashPassword", mock.Anything, exampleInput.NewPassword).Return("blah", errors.New("blah"))
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})

	T.Run("with error updating user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakePasswordUpdateInput()

		res, req := httptest.NewRecorder(), buildRequest(t)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				passwordChangeMiddlewareCtxKey,
				exampleInput,
			),
		)
		req = req.WithContext(
			context.WithValue(req.Context(), v1.SessionInfoKey, exampleUser.ToSessionInfo()),
		)

		mockDB := v11.BuildMockDatabase()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, mock.AnythingOfType("string")).Return(errors.New("blah"))
		s.userDataManager = mockDB

		auth := &mock1.Authenticator{}
		auth.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleInput.CurrentPassword,
			exampleUser.TwoFactorSecret,
			exampleInput.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		auth.On("HashPassword", mock.Anything, exampleInput.NewPassword).Return("blah", nil)
		s.authenticator = auth

		s.UpdatePasswordHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, auth)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_Archive(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
	"net/http"
	"testing"
)

func TestService_Archive(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}
		res, req := httptest.NewRecorder(), buildRequest(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("ArchiveUser", mock.Anything, exampleUser.ID).Return(nil)
		s.userDataManager = mockDB

		r := &mock1.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mock2.UnitCounter{}
		mc.On("Decrement", mock.Anything)
		s.userCounter = mc

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, r, mc)
	})

	T.Run("with error updating database", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		s.userIDFetcher = func(req *http.Request) uint64 {
			return exampleUser.ID
		}
		res, req := httptest.NewRecorder(), buildRequest(t)

		mockDB := v1.BuildMockDatabase()
		mockDB.UserDataManager.On("ArchiveUser", mock.Anything, exampleUser.ID).Return(errors.New("blah"))
		s.userDataManager = mockDB

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

func Test_buildTestService_buildQRCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_buildQRCode(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestService_buildQRCode(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		actual := s.buildQRCode(ctx, exampleUser.Username, exampleUser.TwoFactorSecret)

		assert.NotEmpty(t, actual)
		assert.True(t, strings.HasPrefix(actual, base64ImagePrefix))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
