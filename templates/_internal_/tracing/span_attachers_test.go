package tracing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_spanAttachersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := spanAttachersDotGo(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	trace "go.opencensus.io/trace"
	"strconv"
)

const (
	itemIDSpanAttachmentKey                 = "item_id"
	userIDSpanAttachmentKey                 = "user_id"
	usernameSpanAttachmentKey               = "username"
	filterPageSpanAttachmentKey             = "filter_page"
	filterLimitSpanAttachmentKey            = "filter_limit"
	oauth2ClientDatabaseIDSpanAttachmentKey = "oauth2client_id"
	oauth2ClientIDSpanAttachmentKey         = "client_id"
	webhookIDSpanAttachmentKey              = "webhook_id"
	requestURISpanAttachmentKey             = "request_uri"
	searchQuerySpanAttachmentKey            = "search_query"
)

func attachUint64ToSpan(span *trace.Span, attachmentKey string, id uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute(attachmentKey, strconv.FormatUint(id, 10)))
	}
}

func attachStringToSpan(span *trace.Span, key, str string) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute(key, str))
	}
}

// AttachFilterToSpan provides a consistent way to attach a filter's info to a span.
func AttachFilterToSpan(span *trace.Span, filter *v1.QueryFilter) {
	if filter != nil && span != nil {
		span.AddAttributes(
			trace.StringAttribute(filterPageSpanAttachmentKey, strconv.FormatUint(filter.QueryPage(), 10)),
			trace.StringAttribute(filterLimitSpanAttachmentKey, strconv.FormatUint(uint64(filter.Limit), 10)),
		)
	}
}

// AttachItemIDToSpan attaches an item ID to a given span.
func AttachItemIDToSpan(span *trace.Span, itemID uint64) {
	attachUint64ToSpan(span, itemIDSpanAttachmentKey, itemID)
}

// AttachUserIDToSpan provides a consistent way to attach a user's ID to a span.
func AttachUserIDToSpan(span *trace.Span, userID uint64) {
	attachUint64ToSpan(span, userIDSpanAttachmentKey, userID)
}

// AttachOAuth2ClientDatabaseIDToSpan is a consistent way to attach an oauth2 client's ID to a span.
func AttachOAuth2ClientDatabaseIDToSpan(span *trace.Span, oauth2ClientID uint64) {
	attachUint64ToSpan(span, oauth2ClientDatabaseIDSpanAttachmentKey, oauth2ClientID)
}

// AttachOAuth2ClientIDToSpan is a consistent way to attach an oauth2 client's Client ID to a span.
func AttachOAuth2ClientIDToSpan(span *trace.Span, clientID string) {
	attachStringToSpan(span, oauth2ClientIDSpanAttachmentKey, clientID)
}

// AttachUsernameToSpan provides a consistent way to attach a user's username to a span.
func AttachUsernameToSpan(span *trace.Span, username string) {
	attachStringToSpan(span, usernameSpanAttachmentKey, username)
}

// AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span.
func AttachWebhookIDToSpan(span *trace.Span, webhookID uint64) {
	attachUint64ToSpan(span, webhookIDSpanAttachmentKey, webhookID)
}

// AttachRequestURIToSpan attaches a given URI to a span.
func AttachRequestURIToSpan(span *trace.Span, uri string) {
	attachStringToSpan(span, requestURISpanAttachmentKey, uri)
}

// AttachSearchQueryToSpan attaches a given search query to a span.
func AttachSearchQueryToSpan(span *trace.Span, query string) {
	attachStringToSpan(span, searchQuerySpanAttachmentKey, query)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConstants(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildConstants(proj)

		expected := `
package example

import ()

const (
	itemIDSpanAttachmentKey                 = "item_id"
	userIDSpanAttachmentKey                 = "user_id"
	usernameSpanAttachmentKey               = "username"
	filterPageSpanAttachmentKey             = "filter_page"
	filterLimitSpanAttachmentKey            = "filter_limit"
	oauth2ClientDatabaseIDSpanAttachmentKey = "oauth2client_id"
	oauth2ClientIDSpanAttachmentKey         = "client_id"
	webhookIDSpanAttachmentKey              = "webhook_id"
	requestURISpanAttachmentKey             = "request_uri"
	searchQuerySpanAttachmentKey            = "search_query"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachUint64ToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachUint64ToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
	"strconv"
)

func attachUint64ToSpan(span *trace.Span, attachmentKey string, id uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute(attachmentKey, strconv.FormatUint(id, 10)))
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachStringToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachStringToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

func attachStringToSpan(span *trace.Span, key, str string) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute(key, str))
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachFilterToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildAttachFilterToSpan(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	trace "go.opencensus.io/trace"
	"strconv"
)

// AttachFilterToSpan provides a consistent way to attach a filter's info to a span.
func AttachFilterToSpan(span *trace.Span, filter *v1.QueryFilter) {
	if filter != nil && span != nil {
		span.AddAttributes(
			trace.StringAttribute(filterPageSpanAttachmentKey, strconv.FormatUint(filter.QueryPage(), 10)),
			trace.StringAttribute(filterLimitSpanAttachmentKey, strconv.FormatUint(uint64(filter.Limit), 10)),
		)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachSomethingIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildAttachSomethingIDToSpan(typ)

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachItemIDToSpan attaches an item ID to a given span.
func AttachItemIDToSpan(span *trace.Span, itemID uint64) {
	attachUint64ToSpan(span, itemIDSpanAttachmentKey, itemID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachUserIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachUserIDToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachUserIDToSpan provides a consistent way to attach a user's ID to a span.
func AttachUserIDToSpan(span *trace.Span, userID uint64) {
	attachUint64ToSpan(span, userIDSpanAttachmentKey, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachOAuth2ClientDatabaseIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachOAuth2ClientDatabaseIDToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachOAuth2ClientDatabaseIDToSpan is a consistent way to attach an oauth2 client's ID to a span.
func AttachOAuth2ClientDatabaseIDToSpan(span *trace.Span, oauth2ClientID uint64) {
	attachUint64ToSpan(span, oauth2ClientDatabaseIDSpanAttachmentKey, oauth2ClientID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachOAuth2ClientIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachOAuth2ClientIDToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachOAuth2ClientIDToSpan is a consistent way to attach an oauth2 client's Client ID to a span.
func AttachOAuth2ClientIDToSpan(span *trace.Span, clientID string) {
	attachStringToSpan(span, oauth2ClientIDSpanAttachmentKey, clientID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachUsernameToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachUsernameToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachUsernameToSpan provides a consistent way to attach a user's username to a span.
func AttachUsernameToSpan(span *trace.Span, username string) {
	attachStringToSpan(span, usernameSpanAttachmentKey, username)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachWebhookIDToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachWebhookIDToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span.
func AttachWebhookIDToSpan(span *trace.Span, webhookID uint64) {
	attachUint64ToSpan(span, webhookIDSpanAttachmentKey, webhookID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachRequestURIToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachRequestURIToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachRequestURIToSpan attaches a given URI to a span.
func AttachRequestURIToSpan(span *trace.Span, uri string) {
	attachStringToSpan(span, requestURISpanAttachmentKey, uri)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAttachSearchQueryToSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAttachSearchQueryToSpan()

		expected := `
package example

import (
	trace "go.opencensus.io/trace"
)

// AttachSearchQueryToSpan attaches a given search query to a span.
func AttachSearchQueryToSpan(span *trace.Span, query string) {
	attachStringToSpan(span, searchQuerySpanAttachmentKey, query)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
