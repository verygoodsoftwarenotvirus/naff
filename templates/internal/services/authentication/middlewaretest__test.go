package authentication

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_middlewareTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := middlewareTestDotGo(proj)

		expected := `
package example

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestService_CookieAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		exampleUser := fake.BuildFakeUser()

		md := &mock.UserDataManager{}
		md.On("GetUser", mock1.Anything, mock1.Anything).Return(exampleUser, nil)
		s.userDB = md

		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		res := httptest.NewRecorder()

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		h := s.CookieAuthenticationMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, md, ms)
	})

	T.Run("with nil user", func(t *testing.T) {
		s := buildTestService(t)
		exampleUser := fake.BuildFakeUser()

		md := &mock.UserDataManager{}
		md.On("GetUser", mock1.Anything, mock1.Anything).Return((*v1.User)(nil), nil)
		s.userDB = md

		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		res := httptest.NewRecorder()

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		ms := &MockHTTPHandler{}
		h := s.CookieAuthenticationMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, md, ms)
	})

	T.Run("without user attached", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		res := httptest.NewRecorder()

		ms := &MockHTTPHandler{}
		h := s.CookieAuthenticationMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, ms)
	})
}

func TestService_AuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock1.Anything, mock1.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v11.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock1.Anything, exampleOAuth2Client.BelongsToUser).Return(exampleUser, nil)
		s.userDB = mockDB

		h := &MockHTTPHandler{}
		h.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})

	T.Run("happy path without allowing cookies", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock1.Anything, mock1.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v11.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock1.Anything, exampleOAuth2Client.BelongsToUser).Return(exampleUser, nil)
		s.userDB = mockDB

		h := &MockHTTPHandler{}
		h.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})

	T.Run("with error fetching client but able to use cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		mockDB := v11.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock1.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDB = mockDB

		h := &MockHTTPHandler{}
		h.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, mockDB, h)
	})

	T.Run("able to use cookies but error fetching user info", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		mockDB := v11.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock1.Anything, exampleUser.ID).Return((*v1.User)(nil), errors.New("blah"))
		s.userDB = mockDB

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, mockDB, h)
	})

	T.Run("no cookies allowed, with error fetching user info", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock1.Anything, mock1.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v11.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock1.Anything, exampleOAuth2Client.BelongsToUser).Return((*v1.User)(nil), errors.New("blah"))
		s.userDB = mockDB

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock1.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})

	T.Run("with error fetching client but able to use cookie but unable to decode cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock1.Anything, mock1.Anything).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.oauth2ClientsService = ocv

		cb := &mockCookieEncoderDecoder{}
		cb.On("Decode", CookieName, mock1.Anything, mock1.Anything).Return(errors.New("blah"))
		cb.On("Encode", CookieName, mock1.Anything).Return("", nil)
		s.cookieManager = cb

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, ocv, cb, h)
	})

	T.Run("with invalid authentication", func(t *testing.T) {
		s := buildTestService(t)

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock1.Anything, mock1.Anything).Return((*v1.OAuth2Client)(nil), nil)
		s.oauth2ClientsService = ocv

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, ocv, h)
	})

	T.Run("nightmare path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock1.Anything, mock1.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v11.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock1.Anything, exampleOAuth2Client.BelongsToUser).Return((*v1.User)(nil), nil)
		s.userDB = mockDB

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})
}

func Test_parseLoginInputFromForm(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		exampleUser := fake.BuildFakeUser()
		expected := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		req.Form = map[string][]string{
			usernameFormKey:  {expected.Username},
			passwordFormKey:  {expected.Password},
			totpTokenFormKey: {expected.TOTPToken},
		}

		actual := parseLoginInputFromForm(req)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
	})

	T.Run("returns nil with error parsing form", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		req.URL.RawQuery = "%gh&%ij"
		req.Form = nil

		actual := parseLoginInputFromForm(req)
		assert.Nil(t, actual)
	})
}

func TestService_UserLoginInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		var b bytes.Buffer
		require.NoError(t, json.NewEncoder(&b).Encode(exampleInput))

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", &b)
		require.NoError(t, err)
		require.NotNil(t, req)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		h := s.UserLoginInputMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, ms)
	})

	T.Run("with error decoding request", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		var b bytes.Buffer
		require.NoError(t, json.NewEncoder(&b).Encode(exampleInput))

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", &b)
		require.NoError(t, err)
		require.NotNil(t, req)

		s := buildTestService(t)
		ed := &mock2.EncoderDecoder{}
		ed.On("DecodeRequest", mock1.Anything, mock1.Anything).Return(errors.New("blah"))
		s.encoderDecoder = ed

		ms := &MockHTTPHandler{}
		h := s.UserLoginInputMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, ed, ms)
	})

	T.Run("with error decoding request but valid value attached to form", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		form := url.Values{
			usernameFormKey:  {exampleInput.Username},
			passwordFormKey:  {exampleInput.Password},
			totpTokenFormKey: {exampleInput.TOTPToken},
		}

		req, err := http.NewRequest(
			http.MethodPost,
			"http://todo.verygoodsoftwarenotvirus.ru",
			strings.NewReader(form.Encode()),
		)
		require.NoError(t, err)
		require.NotNil(t, req)

		res := httptest.NewRecorder()
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")

		s := buildTestService(t)
		ed := &mock2.EncoderDecoder{}
		ed.On("DecodeRequest", mock1.Anything, mock1.Anything).Return(errors.New("blah"))
		s.encoderDecoder = ed

		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		h := s.UserLoginInputMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, ed, ms)
	})
}

func TestService_AdminMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		exampleUser := fake.BuildFakeUser()
		exampleUser.IsAdmin = true

		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				v1.SessionInfoKey,
				exampleUser.ToSessionInfo(),
			),
		)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		h := s.AdminMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, ms)
	})

	T.Run("without user attached", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}

		h := s.AdminMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, ms)
	})

	T.Run("with non-admin user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		exampleUser := fake.BuildFakeUser()
		exampleUser.IsAdmin = false

		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				v1.SessionInfoKey,
				exampleUser.ToSessionInfo(),
			),
		)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}

		h := s.AdminMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, ms)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_CookieAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_CookieAuthenticationMiddleware(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestService_CookieAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		exampleUser := fake.BuildFakeUser()

		md := &mock.UserDataManager{}
		md.On("GetUser", mock1.Anything, mock1.Anything).Return(exampleUser, nil)
		s.userDB = md

		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock1.Anything, mock1.Anything).Return()

		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		res := httptest.NewRecorder()

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		h := s.CookieAuthenticationMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, md, ms)
	})

	T.Run("with nil user", func(t *testing.T) {
		s := buildTestService(t)
		exampleUser := fake.BuildFakeUser()

		md := &mock.UserDataManager{}
		md.On("GetUser", mock1.Anything, mock1.Anything).Return((*v1.User)(nil), nil)
		s.userDB = md

		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		res := httptest.NewRecorder()

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		ms := &MockHTTPHandler{}
		h := s.CookieAuthenticationMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock1.AssertExpectationsForObjects(t, md, ms)
	})

	T.Run("without user attached", func(t *testing.T) {
		s := buildTestService(t)

		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)
		res := httptest.NewRecorder()

		ms := &MockHTTPHandler{}
		h := s.CookieAuthenticationMiddleware(ms)
		h.ServeHTTP(res, req)

		mock1.AssertExpectationsForObjects(t, ms)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_AuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_AuthenticationMiddleware(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_AuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock.Anything, mock.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v1.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock.Anything, exampleOAuth2Client.BelongsToUser).Return(exampleUser, nil)
		s.userDB = mockDB

		h := &MockHTTPHandler{}
		h.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})

	T.Run("happy path without allowing cookies", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock.Anything, mock.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v1.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock.Anything, exampleOAuth2Client.BelongsToUser).Return(exampleUser, nil)
		s.userDB = mockDB

		h := &MockHTTPHandler{}
		h.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})

	T.Run("with error fetching client but able to use cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		mockDB := v1.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)
		s.userDB = mockDB

		h := &MockHTTPHandler{}
		h.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		mock.AssertExpectationsForObjects(t, mockDB, h)
	})

	T.Run("able to use cookies but error fetching user info", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		mockDB := v1.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock.Anything, exampleUser.ID).Return((*v11.User)(nil), errors.New("blah"))
		s.userDB = mockDB

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, h)
	})

	T.Run("no cookies allowed, with error fetching user info", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock.Anything, mock.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v1.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock.Anything, exampleOAuth2Client.BelongsToUser).Return((*v11.User)(nil), errors.New("blah"))
		s.userDB = mockDB

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})

	T.Run("with error fetching client but able to use cookie but unable to decode cookie", func(t *testing.T) {
		s := buildTestService(t)

		exampleUser := fake.BuildFakeUser()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock.Anything, mock.Anything).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.oauth2ClientsService = ocv

		cb := &mockCookieEncoderDecoder{}
		cb.On("Decode", CookieName, mock.Anything, mock.Anything).Return(errors.New("blah"))
		cb.On("Encode", CookieName, mock.Anything).Return("", nil)
		s.cookieManager = cb

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		_, req = attachCookieToRequestForTest(t, s, req, exampleUser)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(true)(h).ServeHTTP(res, req)

		mock.AssertExpectationsForObjects(t, ocv, cb, h)
	})

	T.Run("with invalid authentication", func(t *testing.T) {
		s := buildTestService(t)

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock.Anything, mock.Anything).Return((*v11.OAuth2Client)(nil), nil)
		s.oauth2ClientsService = ocv

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, ocv, h)
	})

	T.Run("nightmare path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		ocv := &mockOAuth2ClientValidator{}
		ocv.On("ExtractOAuth2ClientFromRequest", mock.Anything, mock.Anything).Return(exampleOAuth2Client, nil)
		s.oauth2ClientsService = ocv

		mockDB := v1.BuildMockDatabase().UserDataManager
		mockDB.On("GetUser", mock.Anything, exampleOAuth2Client.BelongsToUser).Return((*v11.User)(nil), nil)
		s.userDB = mockDB

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		h := &MockHTTPHandler{}
		s.AuthenticationMiddleware(false)(h).ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, ocv, mockDB, h)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_parseLoginInputFromForm(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTest_parseLoginInputFromForm(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func Test_parseLoginInputFromForm(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		exampleUser := fake.BuildFakeUser()
		expected := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		req.Form = map[string][]string{
			usernameFormKey:  {expected.Username},
			passwordFormKey:  {expected.Password},
			totpTokenFormKey: {expected.TOTPToken},
		}

		actual := parseLoginInputFromForm(req)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
	})

	T.Run("returns nil with error parsing form", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		req.URL.RawQuery = "%gh&%ij"
		req.Form = nil

		actual := parseLoginInputFromForm(req)
		assert.Nil(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_UserLoginInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_UserLoginInputMiddleware(proj)

		expected := `
package example

import (
	"bytes"
	"encoding/json"
	"errors"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestService_UserLoginInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		var b bytes.Buffer
		require.NoError(t, json.NewEncoder(&b).Encode(exampleInput))

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", &b)
		require.NoError(t, err)
		require.NotNil(t, req)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		h := s.UserLoginInputMiddleware(ms)
		h.ServeHTTP(res, req)

		mock.AssertExpectationsForObjects(t, ms)
	})

	T.Run("with error decoding request", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		var b bytes.Buffer
		require.NoError(t, json.NewEncoder(&b).Encode(exampleInput))

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", &b)
		require.NoError(t, err)
		require.NotNil(t, req)

		s := buildTestService(t)
		ed := &mock1.EncoderDecoder{}
		ed.On("DecodeRequest", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.encoderDecoder = ed

		ms := &MockHTTPHandler{}
		h := s.UserLoginInputMiddleware(ms)
		h.ServeHTTP(res, req)

		mock.AssertExpectationsForObjects(t, ed, ms)
	})

	T.Run("with error decoding request but valid value attached to form", func(t *testing.T) {
		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		form := url.Values{
			usernameFormKey:  {exampleInput.Username},
			passwordFormKey:  {exampleInput.Password},
			totpTokenFormKey: {exampleInput.TOTPToken},
		}

		req, err := http.NewRequest(
			http.MethodPost,
			"http://todo.verygoodsoftwarenotvirus.ru",
			strings.NewReader(form.Encode()),
		)
		require.NoError(t, err)
		require.NotNil(t, req)

		res := httptest.NewRecorder()
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")

		s := buildTestService(t)
		ed := &mock1.EncoderDecoder{}
		ed.On("DecodeRequest", mock.Anything, mock.Anything).Return(errors.New("blah"))
		s.encoderDecoder = ed

		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		h := s.UserLoginInputMiddleware(ms)
		h.ServeHTTP(res, req)

		mock.AssertExpectationsForObjects(t, ed, ms)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_AdminMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_AdminMiddleware(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
)

func TestService_AdminMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		exampleUser := fake.BuildFakeUser()
		exampleUser.IsAdmin = true

		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				v1.SessionInfoKey,
				exampleUser.ToSessionInfo(),
			),
		)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}
		ms.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		h := s.AdminMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, ms)
	})

	T.Run("without user attached", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}

		h := s.AdminMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, ms)
	})

	T.Run("with non-admin user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://todo.verygoodsoftwarenotvirus.ru", nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		exampleUser := fake.BuildFakeUser()
		exampleUser.IsAdmin = false

		res := httptest.NewRecorder()
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				v1.SessionInfoKey,
				exampleUser.ToSessionInfo(),
			),
		)

		s := buildTestService(t)
		ms := &MockHTTPHandler{}

		h := s.AdminMiddleware(ms)
		h.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, ms)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
