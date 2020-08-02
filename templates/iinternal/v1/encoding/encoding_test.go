package encoding

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
	"encoding/json"
	"encoding/xml"
	wire "github.com/google/wire"
	"net/http"
	"strings"
)

const (
	// ContentTypeHeader is the HTTP standard header name for content type.
	ContentTypeHeader = "Content-type"
	// XMLContentType represents the XML content type.
	XMLContentType = "application/xml"
	// JSONContentType represents the JSON content type.
	JSONContentType = "application/json"
	// DefaultContentType is what the library defaults to.
	DefaultContentType = JSONContentType
)

var (
	// Providers provides ResponseEncoders for dependency injection.
	Providers = wire.NewSet(
		ProvideResponseEncoder,
	)
)

type (
	// EncoderDecoder is an interface that allows for multiple implementations of HTTP response formats.
	EncoderDecoder interface {
		EncodeResponse(http.ResponseWriter, interface{}) error
		DecodeRequest(*http.Request, interface{}) error
	}

	// ServerEncoderDecoder is our concrete implementation of EncoderDecoder.
	ServerEncoderDecoder struct{}

	encoder interface {
		Encode(v interface{}) error
	}

	decoder interface {
		Decode(v interface{}) error
	}
)

// EncodeResponse encodes responses.
func (ed *ServerEncoderDecoder) EncodeResponse(res http.ResponseWriter, v interface{}) error {
	var ct = strings.ToLower(res.Header().Get(ContentTypeHeader))
	if ct == "" {
		ct = DefaultContentType
	}

	var e encoder
	switch ct {
	case XMLContentType:
		e = xml.NewEncoder(res)
	default:
		e = json.NewEncoder(res)
	}

	res.Header().Set(ContentTypeHeader, ct)
	return e.Encode(v)
}

// DecodeRequest decodes responses.
func (ed *ServerEncoderDecoder) DecodeRequest(req *http.Request, v interface{}) error {
	var ct = strings.ToLower(req.Header.Get(ContentTypeHeader))
	if ct == "" {
		ct = DefaultContentType
	}

	var d decoder
	switch ct {
	case XMLContentType:
		d = xml.NewDecoder(req.Body)
	default:
		d = json.NewDecoder(req.Body)
	}

	return d.Decode(v)
}

// ProvideResponseEncoder provides a jsonResponseEncoder.
func ProvideResponseEncoder() EncoderDecoder {
	return &ServerEncoderDecoder{}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncodingConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildEncodingConstDeclarations()

		expected := `
package example

import ()

const (
	// ContentTypeHeader is the HTTP standard header name for content type.
	ContentTypeHeader = "Content-type"
	// XMLContentType represents the XML content type.
	XMLContentType = "application/xml"
	// JSONContentType represents the JSON content type.
	JSONContentType = "application/json"
	// DefaultContentType is what the library defaults to.
	DefaultContentType = JSONContentType
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncodingVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildEncodingVarDeclarations()

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers provides ResponseEncoders for dependency injection.
	Providers = wire.NewSet(
		ProvideResponseEncoder,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncodingTypeDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildEncodingTypeDeclarations()

		expected := `
package example

import (
	"net/http"
)

type (
	// EncoderDecoder is an interface that allows for multiple implementations of HTTP response formats.
	EncoderDecoder interface {
		EncodeResponse(http.ResponseWriter, interface{}) error
		DecodeRequest(*http.Request, interface{}) error
	}

	// ServerEncoderDecoder is our concrete implementation of EncoderDecoder.
	ServerEncoderDecoder struct{}

	encoder interface {
		Encode(v interface{}) error
	}

	decoder interface {
		Decode(v interface{}) error
	}
)
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
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

// EncodeResponse encodes responses.
func (ed *ServerEncoderDecoder) EncodeResponse(res http.ResponseWriter, v interface{}) error {
	var ct = strings.ToLower(res.Header().Get(ContentTypeHeader))
	if ct == "" {
		ct = DefaultContentType
	}

	var e encoder
	switch ct {
	case XMLContentType:
		e = xml.NewEncoder(res)
	default:
		e = json.NewEncoder(res)
	}

	res.Header().Set(ContentTypeHeader, ct)
	return e.Encode(v)
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
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

// DecodeRequest decodes responses.
func (ed *ServerEncoderDecoder) DecodeRequest(req *http.Request, v interface{}) error {
	var ct = strings.ToLower(req.Header.Get(ContentTypeHeader))
	if ct == "" {
		ct = DefaultContentType
	}

	var d decoder
	switch ct {
	case XMLContentType:
		d = xml.NewDecoder(req.Body)
	default:
		d = json.NewDecoder(req.Body)
	}

	return d.Decode(v)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideResponseEncoder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildProvideResponseEncoder()

		expected := `
package example

import ()

// ProvideResponseEncoder provides a jsonResponseEncoder.
func ProvideResponseEncoder() EncoderDecoder {
	return &ServerEncoderDecoder{}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
