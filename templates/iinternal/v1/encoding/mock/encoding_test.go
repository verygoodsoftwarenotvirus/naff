package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_encodingDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := encodingDotGo(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	"net/http"
)

var _ encoding.EncoderDecoder = (*EncoderDecoder)(nil)

// EncoderDecoder is a mock EncoderDecoder.
type EncoderDecoder struct {
	mock.Mock
}

// EncodeResponse satisfies our EncoderDecoder interface.
func (m *EncoderDecoder) EncodeResponse(res http.ResponseWriter, v interface{}) error {
	return m.Called(res, v).Error(0)
}

// DecodeRequest satisfies our EncoderDecoder interface.
func (m *EncoderDecoder) DecodeRequest(req *http.Request, v interface{}) error {
	return m.Called(req, v).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfaceImplementationDeclaration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildInterfaceImplementationDeclaration(proj)

		expected := `
package example

import (
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
)

var _ encoding.EncoderDecoder = (*EncoderDecoder)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncoderDecoder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildEncoderDecoder()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// EncoderDecoder is a mock EncoderDecoder.
type EncoderDecoder struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncodeResponse(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildEncodeResponse()

		expected := `
package example

import (
	"net/http"
)

// EncodeResponse satisfies our EncoderDecoder interface.
func (m *EncoderDecoder) EncodeResponse(res http.ResponseWriter, v interface{}) error {
	return m.Called(res, v).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDecodeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildDecodeRequest()

		expected := `
package example

import (
	"net/http"
)

// DecodeRequest satisfies our EncoderDecoder interface.
func (m *EncoderDecoder) DecodeRequest(req *http.Request, v interface{}) error {
	return m.Called(req, v).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
