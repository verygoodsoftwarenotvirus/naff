package server

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
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	require "github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var _ http.Handler = (*mockHTTPHandler)(nil)

type mockHTTPHandler struct {
	mock.Mock
}

func (m *mockHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
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

func Test_formatSpanNameForRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req := buildRequest(t)
		req.Method = http.MethodPatch
		req.URL.Path = "/blah"

		expected := "PATCH /blah"
		actual := formatSpanNameForRequest(req)

		assert.Equal(t, expected, actual)
	})
}

func TestServer_loggingMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestServer()

		mh := &mockHTTPHandler{}
		mh.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.loggingMiddleware(mh).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		mock.AssertExpectationsForObjects(t, mh)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMiddlewareTestTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMiddlewareTestTypeDefinitions()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

type mockHTTPHandler struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockHTTPHandlerServeHTTP(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockHTTPHandlerServeHTTP()

		expected := `
package example

import (
	"net/http"
)

func (m *mockHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
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

func Test_buildTest_formatSpanNameForRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTest_formatSpanNameForRequest()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_formatSpanNameForRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		req := buildRequest(t)
		req.Method = http.MethodPatch
		req.URL.Path = "/blah"

		expected := "PATCH /blah"
		actual := formatSpanNameForRequest(req)

		assert.Equal(t, expected, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServer_loggingMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestServer_loggingMiddleware()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestServer_loggingMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestServer()

		mh := &mockHTTPHandler{}
		mh.On("ServeHTTP", mock.Anything, mock.Anything).Return()

		res, req := httptest.NewRecorder(), buildRequest(t)
		s.loggingMiddleware(mh).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		mock.AssertExpectationsForObjects(t, mh)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
