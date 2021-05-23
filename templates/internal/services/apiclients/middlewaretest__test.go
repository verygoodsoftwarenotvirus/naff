package apiclients

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
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	models "gopkg.in/oauth2.v3/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var _ http.Handler = (*mockHTTPHandler)(nil)

type mockHTTPHandler struct {
	mock.Mock
}

func (m *mockHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

func TestService_CreationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		ed := &mock1.EncoderDecoder{}
		ed.On(
			"DecodeRequest",
			mock.AnythingOfType("*http.Request"),
			mock.Anything,
		).Return(nil)
		s.encoderDecoder = ed

		mh := &mockHTTPHandler{}
		mh.On(
			"ServeHTTP",
			mock.Anything,
			mock.Anything,
		)

		h := s.CreationInputMiddleware(mh)
		req := buildRequest(t)
		res := httptest.NewRecorder()

		expected := v1.OAuth2ClientCreationInput{
			RedirectURI: "https://blah.com",
		}
		bs, err := json.Marshal(expected)
		require.NoError(t, err)
		req.Body = ioutil.NopCloser(bytes.NewReader(bs))

		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, ed, mh)
	})

	T.Run("with error decoding request", func(t *testing.T) {
		s := buildTestService(t)

		ed := &mock1.EncoderDecoder{}
		ed.On(
			"DecodeRequest",
			mock.AnythingOfType("*http.Request"),
			mock.Anything,
		).Return(errors.New("blah"))
		s.encoderDecoder = ed

		mh := &mockHTTPHandler{}
		h := s.CreationInputMiddleware(mh)
		req := buildRequest(t)
		res := httptest.NewRecorder()

		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, ed, mh)
	})
}

func TestService_RequestIsAuthenticated(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = fmt.Sprintf("/api/v1/%s", exampleOAuth2Client.Scopes[0])
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mh, mockDB)
	})

	T.Run("with error validating token", func(t *testing.T) {
		s := buildTestService(t)

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return((*models.Token)(nil), errors.New("blah"))
		s.oauth2Handler = mh

		req := buildRequest(t)
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mh)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mh, mockDB)
	})

	T.Run("with invalid scope", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = "/api/v1/stuff"
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mh, mockDB)
	})
}

func TestService_OAuth2TokenAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	// These tests have a lot of overlap to those of ExtractOAuth2ClientFromRequest, which is deliberate.

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = fmt.Sprintf("/api/v1/%s", exampleOAuth2Client.Scopes[0])
		res := httptest.NewRecorder()

		mhh := &mockHTTPHandler{}
		mhh.On(
			"ServeHTTP",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return()

		s.OAuth2TokenAuthenticationMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mh, mhh, mockDB)
	})

	T.Run("with error authenticating request", func(t *testing.T) {
		s := buildTestService(t)

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return((*models.Token)(nil), errors.New("blah"))
		s.oauth2Handler = mh

		res := httptest.NewRecorder()
		req := buildRequest(t)

		mhh := &mockHTTPHandler{}
		s.OAuth2TokenAuthenticationMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, mh, mhh)
	})
}

func TestService_OAuth2ClientInfoMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		expected := "blah"

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mhh := &mockHTTPHandler{}
		mhh.On(
			"ServeHTTP",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return()

		res, req := httptest.NewRecorder(), buildRequest(t)
		q := url.Values{}
		q.Set(oauth2ClientIDURIParamKey, expected)
		req.URL.RawQuery = q.Encode()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			expected,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		s.OAuth2ClientInfoMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mhh, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)
		expected := "blah"

		res, req := httptest.NewRecorder(), buildRequest(t)
		q := url.Values{}
		q.Set(oauth2ClientIDURIParamKey, expected)
		req.URL.RawQuery = q.Encode()

		mockDB := v11.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			expected,
		).Return((*v1.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		mhh := &mockHTTPHandler{}
		s.OAuth2ClientInfoMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, mhh, mockDB)
	})
}

func TestService_fetchOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		req := buildRequest(t).WithContext(
			context.WithValue(
				ctx,
				v1.OAuth2ClientKey,
				exampleOAuth2Client,
			),
		)

		actual := s.fetchOAuth2ClientFromRequest(req)
		assert.Equal(t, exampleOAuth2Client, actual)
	})

	T.Run("without value present", func(t *testing.T) {
		s := buildTestService(t)
		assert.Nil(t, s.fetchOAuth2ClientFromRequest(buildRequest(t)))
	})
}

func TestService_fetchOAuth2ClientIDFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		req := buildRequest(t).WithContext(
			context.WithValue(
				ctx,
				clientIDKey,
				exampleOAuth2Client.ClientID,
			),
		)

		actual := s.fetchOAuth2ClientIDFromRequest(req)
		assert.Equal(t, exampleOAuth2Client.ClientID, actual)
	})

	T.Run("without value present", func(t *testing.T) {
		s := buildTestService(t)

		assert.Empty(t, s.fetchOAuth2ClientIDFromRequest(buildRequest(t)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockHTTPHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockHTTPHandler()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	"net/http"
)

var _ http.Handler = (*mockHTTPHandler)(nil)

type mockHTTPHandler struct {
	mock.Mock
}

func (m *mockHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_CreationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_CreationInputMiddleware(proj)

		expected := `
package example

import (
	"bytes"
	"encoding/json"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock1 "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestService_CreationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		ed := &mock.EncoderDecoder{}
		ed.On(
			"DecodeRequest",
			mock1.AnythingOfType("*http.Request"),
			mock1.Anything,
		).Return(nil)
		s.encoderDecoder = ed

		mh := &mockHTTPHandler{}
		mh.On(
			"ServeHTTP",
			mock1.Anything,
			mock1.Anything,
		)

		h := s.CreationInputMiddleware(mh)
		req := buildRequest(t)
		res := httptest.NewRecorder()

		expected := v1.OAuth2ClientCreationInput{
			RedirectURI: "https://blah.com",
		}
		bs, err := json.Marshal(expected)
		require.NoError(t, err)
		req.Body = ioutil.NopCloser(bytes.NewReader(bs))

		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock1.AssertExpectationsForObjects(t, ed, mh)
	})

	T.Run("with error decoding request", func(t *testing.T) {
		s := buildTestService(t)

		ed := &mock.EncoderDecoder{}
		ed.On(
			"DecodeRequest",
			mock1.AnythingOfType("*http.Request"),
			mock1.Anything,
		).Return(errors.New("blah"))
		s.encoderDecoder = ed

		mh := &mockHTTPHandler{}
		h := s.CreationInputMiddleware(mh)
		req := buildRequest(t)
		res := httptest.NewRecorder()

		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock1.AssertExpectationsForObjects(t, ed, mh)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_RequestIsAuthenticated(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_RequestIsAuthenticated(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	models "gopkg.in/oauth2.v3/models"
	"testing"
)

func TestService_RequestIsAuthenticated(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = fmt.Sprintf("/api/v1/%s", exampleOAuth2Client.Scopes[0])
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mh, mockDB)
	})

	T.Run("with error validating token", func(t *testing.T) {
		s := buildTestService(t)

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return((*models.Token)(nil), errors.New("blah"))
		s.oauth2Handler = mh

		req := buildRequest(t)
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mh)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		req := buildRequest(t)
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mh, mockDB)
	})

	T.Run("with invalid scope", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = "/api/v1/stuff"
		actual, err := s.ExtractOAuth2ClientFromRequest(req.Context(), req)

		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mh, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_OAuth2TokenAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_OAuth2TokenAuthenticationMiddleware(proj)

		expected := `
package example

import (
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	models "gopkg.in/oauth2.v3/models"
	"net/http"
	"testing"
)

func TestService_OAuth2TokenAuthenticationMiddleware(T *testing.T) {
	T.Parallel()

	// These tests have a lot of overlap to those of ExtractOAuth2ClientFromRequest, which is deliberate.

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return(&models.Token{ClientID: exampleOAuth2Client.ClientID}, nil)
		s.oauth2Handler = mh

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			exampleOAuth2Client.ClientID,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		req := buildRequest(t)
		req.URL.Path = fmt.Sprintf("/api/v1/%s", exampleOAuth2Client.Scopes[0])
		res := httptest.NewRecorder()

		mhh := &mockHTTPHandler{}
		mhh.On(
			"ServeHTTP",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return()

		s.OAuth2TokenAuthenticationMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mh, mhh, mockDB)
	})

	T.Run("with error authenticating request", func(t *testing.T) {
		s := buildTestService(t)

		mh := &mockOAuth2Handler{}
		mh.On(
			"ValidationBearerToken",
			mock.AnythingOfType("*http.Request"),
		).Return((*models.Token)(nil), errors.New("blah"))
		s.oauth2Handler = mh

		res := httptest.NewRecorder()
		req := buildRequest(t)

		mhh := &mockHTTPHandler{}
		s.OAuth2TokenAuthenticationMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, mh, mhh)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_OAuth2ClientInfoMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_OAuth2ClientInfoMiddleware(proj)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/url"
	"testing"
)

func TestService_OAuth2ClientInfoMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService(t)
		expected := "blah"

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		mhh := &mockHTTPHandler{}
		mhh.On(
			"ServeHTTP",
			mock.Anything,
			mock.AnythingOfType("*http.Request"),
		).Return()

		res, req := httptest.NewRecorder(), buildRequest(t)
		q := url.Values{}
		q.Set(oauth2ClientIDURIParamKey, expected)
		req.URL.RawQuery = q.Encode()

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			expected,
		).Return(exampleOAuth2Client, nil)
		s.database = mockDB

		s.OAuth2ClientInfoMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mhh, mockDB)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		s := buildTestService(t)
		expected := "blah"

		res, req := httptest.NewRecorder(), buildRequest(t)
		q := url.Values{}
		q.Set(oauth2ClientIDURIParamKey, expected)
		req.URL.RawQuery = q.Encode()

		mockDB := v1.BuildMockDatabase()
		mockDB.OAuth2ClientDataManager.On(
			"GetOAuth2ClientByClientID",
			mock.Anything,
			expected,
		).Return((*v11.OAuth2Client)(nil), errors.New("blah"))
		s.database = mockDB

		mhh := &mockHTTPHandler{}
		s.OAuth2ClientInfoMiddleware(mhh).ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		mock.AssertExpectationsForObjects(t, mhh, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_fetchOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_fetchOAuth2ClientFromRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestService_fetchOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		req := buildRequest(t).WithContext(
			context.WithValue(
				ctx,
				v1.OAuth2ClientKey,
				exampleOAuth2Client,
			),
		)

		actual := s.fetchOAuth2ClientFromRequest(req)
		assert.Equal(t, exampleOAuth2Client, actual)
	})

	T.Run("without value present", func(t *testing.T) {
		s := buildTestService(t)
		assert.Nil(t, s.fetchOAuth2ClientFromRequest(buildRequest(t)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_fetchOAuth2ClientIDFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_fetchOAuth2ClientIDFromRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestService_fetchOAuth2ClientIDFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		req := buildRequest(t).WithContext(
			context.WithValue(
				ctx,
				clientIDKey,
				exampleOAuth2Client.ClientID,
			),
		)

		actual := s.fetchOAuth2ClientIDFromRequest(req)
		assert.Equal(t, exampleOAuth2Client.ClientID, actual)
	})

	T.Run("without value present", func(t *testing.T) {
		s := buildTestService(t)

		assert.Empty(t, s.fetchOAuth2ClientIDFromRequest(buildRequest(t)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
