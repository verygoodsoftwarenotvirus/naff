package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockWebhookDataServerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mockWebhookDataServerDotGo(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

var _ v1.WebhookDataServer = (*WebhookDataServer)(nil)

// WebhookDataServer is a mocked models.WebhookDataServer for testing
type WebhookDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *WebhookDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *WebhookDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements.
func (m *WebhookDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *WebhookDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *WebhookDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *WebhookDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *WebhookDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookDataServer()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// WebhookDataServer is a mocked models.WebhookDataServer for testing
type WebhookDataServer struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookCreationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookCreationInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// CreationInputMiddleware implements our interface requirements.
func (m *WebhookDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookUpdateInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookUpdateInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// UpdateInputMiddleware implements our interface requirements.
func (m *WebhookDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookListHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookListHandler()

		expected := `
package example

import (
	"net/http"
)

// ListHandler implements our interface requirements.
func (m *WebhookDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookCreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookCreateHandler()

		expected := `
package example

import (
	"net/http"
)

// CreateHandler implements our interface requirements.
func (m *WebhookDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookReadHandler()

		expected := `
package example

import (
	"net/http"
)

// ReadHandler implements our interface requirements.
func (m *WebhookDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookUpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookUpdateHandler()

		expected := `
package example

import (
	"net/http"
)

// UpdateHandler implements our interface requirements.
func (m *WebhookDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockWebhookArchiveHandler()

		expected := `
package example

import (
	"net/http"
)

// ArchiveHandler implements our interface requirements.
func (m *WebhookDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
