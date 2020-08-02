package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mockTestDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

var _ OAuth2ClientValidator = (*mockOAuth2ClientValidator)(nil)

type mockOAuth2ClientValidator struct {
	mock.Mock
}

func (m *mockOAuth2ClientValidator) ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}

var _ cookieEncoderDecoder = (*mockCookieEncoderDecoder)(nil)

type mockCookieEncoderDecoder struct {
	mock.Mock
}

func (m *mockCookieEncoderDecoder) Encode(name string, value interface{}) (string, error) {
	args := m.Called(name, value)
	return args.String(0), args.Error(1)
}

func (m *mockCookieEncoderDecoder) Decode(name, value string, dst interface{}) error {
	args := m.Called(name, value, dst)
	return args.Error(0)
}

var _ http.Handler = (*MockHTTPHandler)(nil)

type MockHTTPHandler struct {
	mock.Mock
}

func (m *MockHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockOAuth2ClientValidator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockOAuth2ClientValidator()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

type mockOAuth2ClientValidator struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockOAuth2ClientValidatorExtractOAuth2ClientFromRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMockOAuth2ClientValidatorExtractOAuth2ClientFromRequest(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

func (m *mockOAuth2ClientValidator) ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockCookieEncoderDecoder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockCookieEncoderDecoder()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

type mockCookieEncoderDecoder struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockCookieEncoderDecoderEncode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockCookieEncoderDecoderEncode()

		expected := `
package example

import ()

func (m *mockCookieEncoderDecoder) Encode(name string, value interface{}) (string, error) {
	args := m.Called(name, value)
	return args.String(0), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockCookieEncoderDecoderDecode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockCookieEncoderDecoderDecode()

		expected := `
package example

import ()

func (m *mockCookieEncoderDecoder) Decode(name, value string, dst interface{}) error {
	args := m.Called(name, value, dst)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockHTTPHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockHTTPHandler()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

type MockHTTPHandler struct {
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
		t.Parallel()

		x := buildMockHTTPHandlerServeHTTP()

		expected := `
package example

import (
	"net/http"
)

func (m *MockHTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
