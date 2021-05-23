package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockOauth2ClientDataServerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mockOauth2ClientDataServerDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
)

var _ v1.OAuth2ClientDataServer = (*OAuth2ClientDataServer)(nil)

// OAuth2ClientDataServer is a mocked models.OAuth2ClientDataServer for testing
type OAuth2ClientDataServer struct {
	mock.Mock
}

// ListHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreationInputMiddleware is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// OAuth2ClientInfoMiddleware is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) OAuth2ClientInfoMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ExtractOAuth2ClientFromRequest is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}

// HandleAuthorizeRequest is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error {
	args := m.Called(res, req)
	return args.Error(0)
}

// HandleTokenRequest is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) HandleTokenRequest(res http.ResponseWriter, req *http.Request) error {
	args := m.Called(res, req)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDataServer()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// OAuth2ClientDataServer is a mocked models.OAuth2ClientDataServer for testing
type OAuth2ClientDataServer struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientListHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientListHandler()

		expected := `
package example

import (
	"net/http"
)

// ListHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientCreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientCreateHandler()

		expected := `
package example

import (
	"net/http"
)

// CreateHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientReadHandler()

		expected := `
package example

import (
	"net/http"
)

// ReadHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientArchiveHandler()

		expected := `
package example

import (
	"net/http"
)

// ArchiveHandler is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientCreationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientCreationInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// CreationInputMiddleware is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientInfoMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientInfoMiddleware()

		expected := `
package example

import (
	"net/http"
)

// OAuth2ClientInfoMiddleware is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) OAuth2ClientInfoMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExtractOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildExtractOAuth2ClientFromRequest(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
)

// ExtractOAuth2ClientFromRequest is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHandleAuthorizeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientHandleAuthorizeRequest()

		expected := `
package example

import (
	"net/http"
)

// HandleAuthorizeRequest is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error {
	args := m.Called(res, req)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientHandleTokenRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientHandleTokenRequest()

		expected := `
package example

import (
	"net/http"
)

// HandleTokenRequest is the obligatory implementation for our interface.
func (m *OAuth2ClientDataServer) HandleTokenRequest(res http.ResponseWriter, req *http.Request) error {
	args := m.Called(res, req)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
