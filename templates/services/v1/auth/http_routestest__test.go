package auth

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
	securecookie "github.com/gorilla/securecookie"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func attachCookieToRequestForTest(t *testing.T, s *Service, req *http.Request, user *v1.User) (context.Context, *http.Request) {
	t.Helper()

	ctx, sessionErr := s.sessionManager.Load(req.Context(), "")
	require.NoError(t, sessionErr)
	require.NoError(t, s.sessionManager.RenewToken(ctx))

	// Then make the privilege-level change.
	s.sessionManager.Put(ctx, sessionInfoKey, user.ToSessionInfo())

	token, _, err := s.sessionManager.Commit(ctx)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	c, err := s.buildCookie(token, time.Now().Add(s.config.CookieLifetime))
	require.NoError(t, err)
	req.AddCookie(c)

	return ctx, req.WithContext(ctx)
}

func TestService_DecodeCookieFromRequest(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/api/v1/something", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		ctx, req := attachCookieToRequestForTest(t, s, req, exampleUser)

		cookie, err := s.DecodeCookieFromRequest(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, cookie)
	})

	T.Run("with invalid cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/api/v1/something", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		// begin building bad cookie.
		// NOTE: any code here is duplicated from service.buildAuthCookie
		// any changes made there might need to be reflected here.
		c := &http.Cookie{
			Name:     CookieName,
			Value:    "blah blah blah this is not a real cookie",
			Path:     "/",
			HttpOnly: true,
		}
		// end building bad cookie.
		req.AddCookie(c)

		cookie, err := s.DecodeCookieFromRequest(req.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, cookie)
	})

	T.Run("without cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/api/v1/something", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		cookie, err := s.DecodeCookieFromRequest(req.Context(), req)
		assert.Error(t, err)
		assert.Equal(t, err, http.ErrNoCookie)
		assert.Nil(t, cookie)
	})
}

func TestService_WebsocketAuthFunction(T *testing.T) {
	T.Run("with valid oauth2 client", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		oacv := &mockOAuth2ClientValidator{}
		oacv.On(
			"ExtractOAuth2ClientFromRequest",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = oacv

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		actual := s.WebsocketAuthFunction(req)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, oacv)
	})

	T.Run("with valid cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		oacv := &mockOAuth2ClientValidator{}
		oacv.On(
			"ExtractOAuth2ClientFromRequest",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return(exampleOAuth2Client, errors.New("blah"))
		s.oauth2ClientsService = oacv

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		actual := s.WebsocketAuthFunction(req)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, oacv)
	})

	T.Run("with nothing", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		oacv := &mockOAuth2ClientValidator{}
		oacv.On(
			"ExtractOAuth2ClientFromRequest",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return(exampleOAuth2Client, errors.New("blah"))
		s.oauth2ClientsService = oacv

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		actual := s.WebsocketAuthFunction(req)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, oacv)
	})
}

func TestService_fetchUserFromCookie(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		ctx, req := attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUser",
			mock.Anything,
			exampleUser.ID,
		).Return(exampleUser, nil)
		s.userDB = udb

		actualUser, err := s.fetchUserFromCookie(ctx, req)
		assert.Equal(t, exampleUser, actualUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, udb)
	})

	T.Run("without cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		actualUser, err := s.fetchUserFromCookie(req.Context(), req)
		assert.Nil(t, actualUser)
		assert.Error(t, err)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		expectedError := errors.New("blah")
		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUser",
			mock.Anything,
			exampleUser.ID,
		).Return((*v1.User)(nil), expectedError)
		s.userDB = udb

		actualUser, err := s.fetchUserFromCookie(req.Context(), req)
		assert.Nil(t, actualUser)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, udb)
	})
}

func TestService_LoginHandler(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)
		assert.NotEmpty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, udb, authr)
	})

	T.Run("with error fetching login data from request", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, errors.New("arbitrary"))
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error encoding error fetching login data", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ed := &mock3.EncoderDecoder{}
		ed.On(
			"EncodeResponse",
			mock.AnythingOfType("*httptest.ResponseRecorder"),
			mock.AnythingOfType("*models.ErrorResponse"),
		).Return(errors.New("blah"))
		s.encoderDecoder = ed

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, errors.New("arbitrary"))
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, ed, udb)
	})

	T.Run("with invalid login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, udb, authr)
	})

	T.Run("with error validating login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, errors.New("blah"))
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, udb, authr)
	})

	T.Run("with error building cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		cb := &mockCookieEncoderDecoder{}
		cb.On(
			"Encode",
			CookieName,
			mock.AnythingOfType("string"),
		).Return("", errors.New("blah"))
		s.cookieManager = cb

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, cb, udb, authr)
	})

	T.Run("with error building cookie and error encoding cookie response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		cb := &mockCookieEncoderDecoder{}
		cb.On(
			"Encode",
			CookieName,
			mock.AnythingOfType("string"),
		).Return("", errors.New("blah"))
		s.cookieManager = cb

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock.AssertExpectationsForObjects(t, cb, udb, authr)
	})
}

func TestService_LogoutHandler(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		res := httptest.NewRecorder()

		s.LogoutHandler(res, req)

		actualCookie := res.Header().Get("Set-Cookie")
		assert.Contains(t, actualCookie, "Max-Age=0")
	})

	T.Run("without cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		s.LogoutHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error building cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)
		s.cookieManager = securecookie.New(
			securecookie.GenerateRandomKey(0),
			[]byte(""),
		)

		res := httptest.NewRecorder()

		s.LogoutHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestService_fetchLoginDataFromRequest(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))
		loginData, err := s.fetchLoginDataFromRequest(req)

		require.NotNil(t, loginData)
		assert.Equal(t, loginData.user, exampleUser)
		assert.Nil(t, err)

		mock.AssertExpectationsForObjects(t, udb)
	})

	T.Run("without login data attached to request", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, err = s.fetchLoginDataFromRequest(req)
		assert.Error(t, err)
	})

	T.Run("with DB error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return((*v1.User)(nil), sql.ErrNoRows)
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))
		_, err = s.fetchLoginDataFromRequest(req)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock.Anything,
			exampleUser.Username,
		).Return((*v1.User)(nil), errors.New("blah"))
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))
		_, err = s.fetchLoginDataFromRequest(req)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, udb)
	})
}

func TestService_validateLogin(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, authr)
	})

	T.Run("with too weak a password hash", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, auth.ErrCostTooLow)
		s.authenticator = authr

		authr.On(
			"HashPassword",
			mock.Anything,
			exampleLoginData.Password,
		).Return("blah", nil)

		udb := &mock1.UserDataManager{}
		udb.On(
			"UpdateUser",
			mock.Anything,
			mock.AnythingOfType("*models.User"),
		).Return(nil)
		s.userDB = udb

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, authr, udb)
	})

	T.Run("with too weak a password hash and error hashing the password", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		expectedErr := errors.New("arbitrary")

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, auth.ErrCostTooLow)

		authr.On(
			"HashPassword",
			mock.Anything,
			exampleLoginData.Password,
		).Return("", expectedErr)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, authr)
	})

	T.Run("with too weak a password hash and error updating user", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		expectedErr := errors.New("arbitrary")
		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, auth.ErrCostTooLow)

		authr.On(
			"HashPassword",
			mock.Anything,
			exampleLoginData.Password,
		).Return("blah", nil)
		s.authenticator = authr

		udb := &mock1.UserDataManager{}
		udb.On(
			"UpdateUser",
			mock.Anything,
			mock.AnythingOfType("*models.User"),
		).Return(expectedErr)
		s.userDB = udb

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, authr, udb)
	})

	T.Run("with error validating login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		expectedErr := errors.New("arbitrary")
		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(false, expectedErr)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, authr)
	})

	T.Run("with invalid login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, authr)
	})
}

func TestService_StatusHandler(T *testing.T) {
	T.Run("normal operation", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUser",
			mock.Anything,
			exampleUser.ID,
		).Return(exampleUser, nil)
		s.userDB = udb

		s.StatusHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUser",
			mock.Anything,
			exampleUser.ID,
		).Return((*v1.User)(nil), errors.New("blah"))
		s.userDB = udb

		s.StatusHandler(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock1.UserDataManager{}
		udb.On(
			"GetUser",
			mock.Anything,
			exampleUser.ID,
		).Return(exampleUser, nil)
		s.userDB = udb

		ed := &mock3.EncoderDecoder{}
		ed.On(
			"EncodeResponse",
			mock.AnythingOfType("*httptest.ResponseRecorder"),
			mock.AnythingOfType("*models.StatusResponse"),
		).Return(errors.New("blah"))
		s.encoderDecoder = ed

		s.StatusHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, udb, ed)
	})
}

func TestService_CycleSecretHandler(T *testing.T) {
	T.Run("normal operation", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)
		c := req.Cookies()[0]

		var token string
		assert.NoError(t, s.cookieManager.Decode(CookieName, c.Value, &token))
		s.CycleSecretHandler(res, req)

		assert.Error(t, s.cookieManager.Decode(CookieName, c.Value, &token))
	})
}

func TestService_buildCookie(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		cookie, err := s.buildCookie("example", time.Now().Add(s.config.CookieLifetime))
		assert.NotNil(t, cookie)
		assert.NoError(t, err)
	})

	T.Run("with erroneous cookie building setup", func(t *testing.T) {
		s := buildTestService(t)
		s.cookieManager = securecookie.New(
			securecookie.GenerateRandomKey(0),
			[]byte(""),
		)

		cookie, err := s.buildCookie("example", time.Now().Add(s.config.CookieLifetime))
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachCookieToRequestForTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildAttachCookieToRequestForTest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"testing"
	"time"
)

func attachCookieToRequestForTest(t *testing.T, s *Service, req *http.Request, user *v1.User) (context.Context, *http.Request) {
	t.Helper()

	ctx, sessionErr := s.sessionManager.Load(req.Context(), "")
	require.NoError(t, sessionErr)
	require.NoError(t, s.sessionManager.RenewToken(ctx))

	// Then make the privilege-level change.
	s.sessionManager.Put(ctx, sessionInfoKey, user.ToSessionInfo())

	token, _, err := s.sessionManager.Commit(ctx)
	assert.NotEmpty(t, token)
	assert.NoError(t, err)

	c, err := s.buildCookie(token, time.Now().Add(s.config.CookieLifetime))
	require.NoError(t, err)
	req.AddCookie(c)

	return ctx, req.WithContext(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_DecodeCookieFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_DecodeCookieFromRequest(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_DecodeCookieFromRequest(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/api/v1/something", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		ctx, req := attachCookieToRequestForTest(t, s, req, exampleUser)

		cookie, err := s.DecodeCookieFromRequest(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, cookie)
	})

	T.Run("with invalid cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/api/v1/something", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		// begin building bad cookie.
		// NOTE: any code here is duplicated from service.buildAuthCookie
		// any changes made there might need to be reflected here.
		c := &http.Cookie{
			Name:     CookieName,
			Value:    "blah blah blah this is not a real cookie",
			Path:     "/",
			HttpOnly: true,
		}
		// end building bad cookie.
		req.AddCookie(c)

		cookie, err := s.DecodeCookieFromRequest(req.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, cookie)
	})

	T.Run("without cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/api/v1/something", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		cookie, err := s.DecodeCookieFromRequest(req.Context(), req)
		assert.Error(t, err)
		assert.Equal(t, err, http.ErrNoCookie)
		assert.Nil(t, cookie)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_WebsocketAuthFunction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_WebsocketAuthFunction(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_WebsocketAuthFunction(T *testing.T) {
	T.Run("with valid oauth2 client", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		oacv := &mockOAuth2ClientValidator{}
		oacv.On(
			"ExtractOAuth2ClientFromRequest",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = oacv

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		actual := s.WebsocketAuthFunction(req)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, oacv)
	})

	T.Run("with valid cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		oacv := &mockOAuth2ClientValidator{}
		oacv.On(
			"ExtractOAuth2ClientFromRequest",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return(exampleOAuth2Client, errors.New("blah"))
		s.oauth2ClientsService = oacv

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		actual := s.WebsocketAuthFunction(req)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, oacv)
	})

	T.Run("with nothing", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		oacv := &mockOAuth2ClientValidator{}
		oacv.On(
			"ExtractOAuth2ClientFromRequest",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return(exampleOAuth2Client, errors.New("blah"))
		s.oauth2ClientsService = oacv

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		actual := s.WebsocketAuthFunction(req)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, oacv)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_fetchUserFromCookie(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_fetchUserFromCookie(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestService_fetchUserFromCookie(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		ctx, req := attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUser",
			mock1.Anything,
			exampleUser.ID,
		).Return(exampleUser, nil)
		s.userDB = udb

		actualUser, err := s.fetchUserFromCookie(ctx, req)
		assert.Equal(t, exampleUser, actualUser)
		assert.NoError(t, err)

		mock1.AssertExpectationsForObjects(t, udb)
	})

	T.Run("without cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		actualUser, err := s.fetchUserFromCookie(req.Context(), req)
		assert.Nil(t, actualUser)
		assert.Error(t, err)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		expectedError := errors.New("blah")
		udb := &mock.UserDataManager{}
		udb.On(
			"GetUser",
			mock1.Anything,
			exampleUser.ID,
		).Return((*v1.User)(nil), expectedError)
		s.userDB = udb

		actualUser, err := s.fetchUserFromCookie(req.Context(), req)
		assert.Nil(t, actualUser)
		assert.Error(t, err)

		mock1.AssertExpectationsForObjects(t, udb)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_LoginHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_LoginHandler(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestService_LoginHandler(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)
		assert.NotEmpty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, udb, authr)
	})

	T.Run("with error fetching login data from request", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, errors.New("arbitrary"))
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error encoding error fetching login data", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ed := &mock3.EncoderDecoder{}
		ed.On(
			"EncodeResponse",
			mock1.AnythingOfType("*httptest.ResponseRecorder"),
			mock1.AnythingOfType("*models.ErrorResponse"),
		).Return(errors.New("blah"))
		s.encoderDecoder = ed

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, errors.New("arbitrary"))
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, ed, udb)
	})

	T.Run("with invalid login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, udb, authr)
	})

	T.Run("with error validating login", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, errors.New("blah"))
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, udb, authr)
	})

	T.Run("with error building cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		cb := &mockCookieEncoderDecoder{}
		cb.On(
			"Encode",
			CookieName,
			mock1.AnythingOfType("string"),
		).Return("", errors.New("blah"))
		s.cookieManager = cb

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, cb, udb, authr)
	})

	T.Run("with error building cookie and error encoding cookie response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		cb := &mockCookieEncoderDecoder{}
		cb.On(
			"Encode",
			CookieName,
			mock1.AnythingOfType("string"),
		).Return("", errors.New("blah"))
		s.cookieManager = cb

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		authr := &mock2.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))

		s.LoginHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Empty(t, res.Header().Get("Set-Cookie"))

		mock1.AssertExpectationsForObjects(t, cb, udb, authr)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_LogoutHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_LogoutHandler(proj)

		expected := `
package example

import (
	securecookie "github.com/gorilla/securecookie"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_LogoutHandler(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		res := httptest.NewRecorder()

		s.LogoutHandler(res, req)

		actualCookie := res.Header().Get("Set-Cookie")
		assert.Contains(t, actualCookie, "Max-Age=0")
	})

	T.Run("without cookie", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		res := httptest.NewRecorder()
		s.LogoutHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error building cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)
		s.cookieManager = securecookie.New(
			securecookie.GenerateRandomKey(0),
			[]byte(""),
		)

		res := httptest.NewRecorder()

		s.LogoutHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_fetchLoginDataFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_fetchLoginDataFromRequest(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"testing"
)

func TestService_fetchLoginDataFromRequest(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return(exampleUser, nil)
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))
		loginData, err := s.fetchLoginDataFromRequest(req)

		require.NotNil(t, loginData)
		assert.Equal(t, loginData.user, exampleUser)
		assert.Nil(t, err)

		mock1.AssertExpectationsForObjects(t, udb)
	})

	T.Run("without login data attached to request", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, err = s.fetchLoginDataFromRequest(req)
		assert.Error(t, err)
	})

	T.Run("with DB error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return((*v1.User)(nil), sql.ErrNoRows)
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))
		_, err = s.fetchLoginDataFromRequest(req)
		assert.Error(t, err)

		mock1.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUserByUsername",
			mock1.Anything,
			exampleUser.Username,
		).Return((*v1.User)(nil), errors.New("blah"))
		s.userDB = udb

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru/testing", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), userLoginInputMiddlewareCtxKey, exampleLoginData))
		_, err = s.fetchLoginDataFromRequest(req)
		assert.Error(t, err)

		mock1.AssertExpectationsForObjects(t, udb)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_validateLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_validateLogin(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"testing"
)

func TestService_validateLogin(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, nil)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock1.AssertExpectationsForObjects(t, authr)
	})

	T.Run("with too weak a password hash", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, auth.ErrCostTooLow)
		s.authenticator = authr

		authr.On(
			"HashPassword",
			mock1.Anything,
			exampleLoginData.Password,
		).Return("blah", nil)

		udb := &mock2.UserDataManager{}
		udb.On(
			"UpdateUser",
			mock1.Anything,
			mock1.AnythingOfType("*models.User"),
		).Return(nil)
		s.userDB = udb

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock1.AssertExpectationsForObjects(t, authr, udb)
	})

	T.Run("with too weak a password hash and error hashing the password", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		expectedErr := errors.New("arbitrary")

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, auth.ErrCostTooLow)

		authr.On(
			"HashPassword",
			mock1.Anything,
			exampleLoginData.Password,
		).Return("", expectedErr)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock1.AssertExpectationsForObjects(t, authr)
	})

	T.Run("with too weak a password hash and error updating user", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		expectedErr := errors.New("arbitrary")
		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(true, auth.ErrCostTooLow)

		authr.On(
			"HashPassword",
			mock1.Anything,
			exampleLoginData.Password,
		).Return("blah", nil)
		s.authenticator = authr

		udb := &mock2.UserDataManager{}
		udb.On(
			"UpdateUser",
			mock1.Anything,
			mock1.AnythingOfType("*models.User"),
		).Return(expectedErr)
		s.userDB = udb

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock1.AssertExpectationsForObjects(t, authr, udb)
	})

	T.Run("with error validating login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		expectedErr := errors.New("arbitrary")
		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(false, expectedErr)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.Error(t, err)

		mock1.AssertExpectationsForObjects(t, authr)
	})

	T.Run("with invalid login", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleLoginData := fake.BuildFakeUserLoginInputFromUser(exampleUser)
		exampleInput := loginData{
			loginInput: exampleLoginData,
			user:       exampleUser,
		}

		authr := &mock.Authenticator{}
		authr.On(
			"ValidateLogin",
			mock1.Anything,
			exampleUser.HashedPassword,
			exampleLoginData.Password,
			exampleUser.TwoFactorSecret,
			exampleLoginData.TOTPToken,
			exampleUser.Salt,
		).Return(false, nil)
		s.authenticator = authr

		actual, err := s.validateLogin(ctx, exampleInput)
		assert.False(t, actual)
		assert.NoError(t, err)

		mock1.AssertExpectationsForObjects(t, authr)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_StatusHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_StatusHandler(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestService_StatusHandler(T *testing.T) {
	T.Run("normal operation", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUser",
			mock1.Anything,
			exampleUser.ID,
		).Return(exampleUser, nil)
		s.userDB = udb

		s.StatusHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error fetching user", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUser",
			mock1.Anything,
			exampleUser.ID,
		).Return((*v1.User)(nil), errors.New("blah"))
		s.userDB = udb

		s.StatusHandler(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, udb)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		udb := &mock.UserDataManager{}
		udb.On(
			"GetUser",
			mock1.Anything,
			exampleUser.ID,
		).Return(exampleUser, nil)
		s.userDB = udb

		ed := &mock2.EncoderDecoder{}
		ed.On(
			"EncodeResponse",
			mock1.AnythingOfType("*httptest.ResponseRecorder"),
			mock1.AnythingOfType("*models.StatusResponse"),
		).Return(errors.New("blah"))
		s.encoderDecoder = ed

		s.StatusHandler(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, udb, ed)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_CycleSecretHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestService_CycleSecretHandler(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_CycleSecretHandler(T *testing.T) {
	T.Run("normal operation", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "https://blah.com", nil)
		require.NotNil(t, req)
		require.NoError(t, err)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)
		c := req.Cookies()[0]

		var token string
		assert.NoError(t, s.cookieManager.Decode(CookieName, c.Value, &token))
		s.CycleSecretHandler(res, req)

		assert.Error(t, s.cookieManager.Decode(CookieName, c.Value, &token))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_buildCookie(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestService_buildCookie()

		expected := `
package example

import (
	securecookie "github.com/gorilla/securecookie"
	assert "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_buildCookie(T *testing.T) {
	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		cookie, err := s.buildCookie("example", time.Now().Add(s.config.CookieLifetime))
		assert.NotNil(t, cookie)
		assert.NoError(t, err)
	})

	T.Run("with erroneous cookie building setup", func(t *testing.T) {
		s := buildTestService(t)
		s.cookieManager = securecookie.New(
			securecookie.GenerateRandomKey(0),
			[]byte(""),
		)

		cookie, err := s.buildCookie("example", time.Now().Add(s.config.CookieLifetime))
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
