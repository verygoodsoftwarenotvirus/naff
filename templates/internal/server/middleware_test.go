package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_middlewareDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := middlewareDotGo(proj)

		expected := `
package example

import (
	"fmt"
	middleware "github.com/go-chi/chi/middleware"
	"net/http"
	"regexp"
	"time"
)

var (
	idReplacementRegex = regexp.MustCompile(` + "`" + `[^(v|oauth)]\\d+` + "`" + `)
)

func formatSpanNameForRequest(req *http.Request) string {
	return fmt.Sprintf(
		"%s %s",
		req.Method,
		idReplacementRegex.ReplaceAllString(req.URL.Path, "/{id}"),
	)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ww := middleware.NewWrapResponseWriter(res, req.ProtoMajor)

		start := time.Now()
		next.ServeHTTP(ww, req)

		s.logger.WithRequest(req).WithValues(map[string]interface{}{
			"status":        ww.Status(),
			"bytes_written": ww.BytesWritten(),
			"elapsed":       time.Since(start),
		}).Debug("responded to request")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerMiddlewareVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServerMiddlewareVarDeclarations()

		expected := `
package example

import (
	"regexp"
)

var (
	idReplacementRegex = regexp.MustCompile(` + "`" + `[^(v|oauth)]\\d+` + "`" + `)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerFormatSpanNameForRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServerFormatSpanNameForRequest()

		expected := `
package example

import (
	"fmt"
	"net/http"
)

func formatSpanNameForRequest(req *http.Request) string {
	return fmt.Sprintf(
		"%s %s",
		req.Method,
		idReplacementRegex.ReplaceAllString(req.URL.Path, "/{id}"),
	)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServerServerLoggingMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServerServerLoggingMiddleware()

		expected := `
package example

import (
	middleware "github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ww := middleware.NewWrapResponseWriter(res, req.ProtoMajor)

		start := time.Now()
		next.ServeHTTP(ww, req)

		s.logger.WithRequest(req).WithValues(map[string]interface{}{
			"status":        ww.Status(),
			"bytes_written": ww.BytesWritten(),
			"elapsed":       time.Since(start),
		}).Debug("responded to request")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
