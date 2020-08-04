package encoding

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_encodingTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := encodingTestDotGo(proj)

		expected := `
package example

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type example struct {
	Name string ` + "`" + `json:"name" xml:"name"` + "`" + `
}

func TestServerEncoderDecoder_EncodeResponse(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectation := "name"
		ex := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		res := httptest.NewRecorder()
		err := ed.EncodeResponse(res, ex)

		assert.NoError(t, err)
		assert.Equal(t, res.Body.String(), fmt.Sprintf("{%q:%q}\n", "name", ex.Name))
	})

	T.Run("as XML", func(t *testing.T) {
		expectation := "name"
		ex := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeader, "application/xml")

		err := ed.EncodeResponse(res, ex)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("<example><name>%s</name></example>", expectation), res.Body.String())
	})
}

func TestServerEncoderDecoder_DecodeRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectation := "name"
		e := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		bs, err := json.Marshal(e)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(bs))
		require.NoError(t, err)

		var x example
		assert.NoError(t, ed.DecodeRequest(req, &x))
		assert.Equal(t, x.Name, e.Name)
	})

	T.Run("as XML", func(t *testing.T) {
		expectation := "name"
		e := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		bs, err := xml.Marshal(e)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(bs))
		require.NoError(t, err)
		req.Header.Set(ContentTypeHeader, XMLContentType)

		var x example
		assert.NoError(t, ed.DecodeRequest(req, &x))
		assert.Equal(t, x.Name, e.Name)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncodingTestTypeDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildEncodingTestTypeDeclarations()

		expected := `
package example

import ()

type example struct {
	Name string ` + "`" + `json:"name" xml:"name"` + "`" + `
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServerEncoderDecoder_EncodeResponse(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestServerEncoderDecoder_EncodeResponse()

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestServerEncoderDecoder_EncodeResponse(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectation := "name"
		ex := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		res := httptest.NewRecorder()
		err := ed.EncodeResponse(res, ex)

		assert.NoError(t, err)
		assert.Equal(t, res.Body.String(), fmt.Sprintf("{%q:%q}\n", "name", ex.Name))
	})

	T.Run("as XML", func(t *testing.T) {
		expectation := "name"
		ex := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		res := httptest.NewRecorder()
		res.Header().Set(ContentTypeHeader, "application/xml")

		err := ed.EncodeResponse(res, ex)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("<example><name>%s</name></example>", expectation), res.Body.String())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestServerEncoderDecoder_DecodeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestServerEncoderDecoder_DecodeRequest()

		expected := `
package example

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestServerEncoderDecoder_DecodeRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectation := "name"
		e := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		bs, err := json.Marshal(e)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(bs))
		require.NoError(t, err)

		var x example
		assert.NoError(t, ed.DecodeRequest(req, &x))
		assert.Equal(t, x.Name, e.Name)
	})

	T.Run("as XML", func(t *testing.T) {
		expectation := "name"
		e := &example{Name: expectation}
		ed := ProvideResponseEncoder()

		bs, err := xml.Marshal(e)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "http://todo.verygoodsoftwarenotvirus.ru", bytes.NewReader(bs))
		require.NoError(t, err)
		req.Header.Set(ContentTypeHeader, XMLContentType)

		var x example
		assert.NoError(t, ed.DecodeRequest(req, &x))
		assert.Equal(t, x.Name, e.Name)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
