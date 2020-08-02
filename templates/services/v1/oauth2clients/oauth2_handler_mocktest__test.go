package oauth2clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2HandlerMockTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := oauth2HandlerMockTestDotGo(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	oauth2v3 "gopkg.in/oauth2.v3"
	server "gopkg.in/oauth2.v3/server"
	"net/http"
)

var _ oauth2Handler = (*mockOAuth2Handler)(nil)

type mockOAuth2Handler struct {
	mock.Mock
}

func (m *mockOAuth2Handler) SetAllowGetAccessRequest(allowed bool) {
	m.Called(allowed)
}

func (m *mockOAuth2Handler) SetClientAuthorizedHandler(handler server.ClientAuthorizedHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) SetClientScopeHandler(handler server.ClientScopeHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) SetClientInfoHandler(handler server.ClientInfoHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) SetUserAuthorizationHandler(handler server.UserAuthorizationHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) SetAuthorizeScopeHandler(handler server.AuthorizeScopeHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) SetResponseErrorHandler(handler server.ResponseErrorHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) SetInternalErrorHandler(handler server.InternalErrorHandler) {
	m.Called(handler)
}

func (m *mockOAuth2Handler) ValidationBearerToken(req *http.Request) (oauth2v3.TokenInfo, error) {
	args := m.Called(req)
	return args.Get(0).(oauth2v3.TokenInfo), args.Error(1)
}

func (m *mockOAuth2Handler) HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error {
	return m.Called(res, req).Error(0)
}

func (m *mockOAuth2Handler) HandleTokenRequest(res http.ResponseWriter, req *http.Request) error {
	return m.Called(res, req).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestMockOAuth2Handler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestMockOAuth2Handler()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

var _ oauth2Handler = (*mockOAuth2Handler)(nil)

type mockOAuth2Handler struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetAllowGetAccessRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetAllowGetAccessRequest()

		expected := `
package example

import ()

func (m *mockOAuth2Handler) SetAllowGetAccessRequest(allowed bool) {
	m.Called(allowed)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetClientAuthorizedHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetClientAuthorizedHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetClientAuthorizedHandler(handler server.ClientAuthorizedHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetClientScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetClientScopeHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetClientScopeHandler(handler server.ClientScopeHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetClientInfoHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetClientInfoHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetClientInfoHandler(handler server.ClientInfoHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetUserAuthorizationHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetUserAuthorizationHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetUserAuthorizationHandler(handler server.UserAuthorizationHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetAuthorizeScopeHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetAuthorizeScopeHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetAuthorizeScopeHandler(handler server.AuthorizeScopeHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetResponseErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetResponseErrorHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetResponseErrorHandler(handler server.ResponseErrorHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestSetInternalErrorHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestSetInternalErrorHandler()

		expected := `
package example

import (
	server "gopkg.in/oauth2.v3/server"
)

func (m *mockOAuth2Handler) SetInternalErrorHandler(handler server.InternalErrorHandler) {
	m.Called(handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestValidationBearerToken(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestValidationBearerToken()

		expected := `
package example

import (
	oauth2v3 "gopkg.in/oauth2.v3"
	"net/http"
)

func (m *mockOAuth2Handler) ValidationBearerToken(req *http.Request) (oauth2v3.TokenInfo, error) {
	args := m.Called(req)
	return args.Get(0).(oauth2v3.TokenInfo), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestHandleAuthorizeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestHandleAuthorizeRequest()

		expected := `
package example

import (
	"net/http"
)

func (m *mockOAuth2Handler) HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error {
	return m.Called(res, req).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2HandlerMockTestHandleTokenRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2HandlerMockTestHandleTokenRequest()

		expected := `
package example

import (
	"net/http"
)

func (m *mockOAuth2Handler) HandleTokenRequest(res http.ResponseWriter, req *http.Request) error {
	return m.Called(res, req).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
