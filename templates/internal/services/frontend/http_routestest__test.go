package frontend

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
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"net/http"
	"os"
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

func TestService_StaticDir(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}

		cwd, err := os.Getwd()
		require.NoError(t, err)

		hf, err := s.StaticDir(cwd)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/http_routes_test.go"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/login"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend items routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/items/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestService_buildStaticFileServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s := &Service{
			config: config.FrontendSettings{
				CacheStaticFiles: true,
			},
		}
		cwd, err := os.Getwd()
		require.NoError(t, err)

		actual, err := s.buildStaticFileServer(cwd)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
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

func Test_buildTestService_StaticDir(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_StaticDir(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"net/http"
	"os"
	"testing"
)

func TestService_StaticDir(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}

		cwd, err := os.Getwd()
		require.NoError(t, err)

		hf, err := s.StaticDir(cwd)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/http_routes_test.go"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/login"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with frontend items routing path", func(t *testing.T) {
		s := &Service{logger: noop.ProvideNoopLogger()}
		exampleDir := "."

		hf, err := s.StaticDir(exampleDir)
		assert.NoError(t, err)
		assert.NotNil(t, hf)

		req, res := buildRequest(t), httptest.NewRecorder()
		req.URL.Path = "/items/123"
		hf(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestService_buildStaticFileServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestService_buildStaticFileServer(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"os"
	"testing"
)

func TestService_buildStaticFileServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		s := &Service{
			config: config.FrontendSettings{
				CacheStaticFiles: true,
			},
		}
		cwd, err := os.Getwd()
		require.NoError(t, err)

		actual, err := s.buildStaticFileServer(cwd)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
