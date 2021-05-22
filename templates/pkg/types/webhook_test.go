package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhookDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhookDotGo(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

type (
	// Webhook represents a webhook listener, an endpoint to send an HTTP request to upon an event.
	Webhook struct {
		ID            uint64   ` + "`" + `json:"id"` + "`" + `
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ContentType   string   ` + "`" + `json:"contentType"` + "`" + `
		URL           string   ` + "`" + `json:"url"` + "`" + `
		Method        string   ` + "`" + `json:"method"` + "`" + `
		Events        []string ` + "`" + `json:"events"` + "`" + `
		DataTypes     []string ` + "`" + `json:"dataTypes"` + "`" + `
		Topics        []string ` + "`" + `json:"topics"` + "`" + `
		CreatedOn     uint64   ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn *uint64  ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn    *uint64  ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	// WebhookCreationInput represents what a user could set as input for creating a webhook.
	WebhookCreationInput struct {
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ContentType   string   ` + "`" + `json:"contentType"` + "`" + `
		URL           string   ` + "`" + `json:"url"` + "`" + `
		Method        string   ` + "`" + `json:"method"` + "`" + `
		Events        []string ` + "`" + `json:"events"` + "`" + `
		DataTypes     []string ` + "`" + `json:"dataTypes"` + "`" + `
		Topics        []string ` + "`" + `json:"topics"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"-"` + "`" + `
	}

	// WebhookUpdateInput represents what a user could set as input for updating a webhook.
	WebhookUpdateInput struct {
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ContentType   string   ` + "`" + `json:"contentType"` + "`" + `
		URL           string   ` + "`" + `json:"url"` + "`" + `
		Method        string   ` + "`" + `json:"method"` + "`" + `
		Events        []string ` + "`" + `json:"events"` + "`" + `
		DataTypes     []string ` + "`" + `json:"dataTypes"` + "`" + `
		Topics        []string ` + "`" + `json:"topics"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"-"` + "`" + `
	}

	// WebhookList represents a list of webhooks.
	WebhookList struct {
		Pagination
		Webhooks []Webhook ` + "`" + `json:"webhooks"` + "`" + `
	}

	// WebhookDataManager describes a structure capable of storing webhooks.
	WebhookDataManager interface {
		GetWebhook(ctx context.Context, webhookID, userID uint64) (*Webhook, error)
		GetAllWebhooksCount(ctx context.Context) (uint64, error)
		GetWebhooks(ctx context.Context, userID uint64, filter *QueryFilter) (*WebhookList, error)
		GetAllWebhooks(ctx context.Context) (*WebhookList, error)
		CreateWebhook(ctx context.Context, input *WebhookCreationInput) (*Webhook, error)
		UpdateWebhook(ctx context.Context, updated *Webhook) error
		ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error
	}

	// WebhookDataServer describes a structure capable of serving traffic related to webhooks.
	WebhookDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an WebhookCreationInput with an Webhook.
func (w *Webhook) Update(input *WebhookUpdateInput) {
	if input.Name != "" {
		w.Name = input.Name
	}
	if input.ContentType != "" {
		w.ContentType = input.ContentType
	}
	if input.URL != "" {
		w.URL = input.URL
	}
	if input.Method != "" {
		w.Method = input.Method
	}

	if input.Events != nil && len(input.Events) > 0 {
		w.Events = input.Events
	}
	if input.DataTypes != nil && len(input.DataTypes) > 0 {
		w.DataTypes = input.DataTypes
	}
	if input.Topics != nil && len(input.Topics) > 0 {
		w.Topics = input.Topics
	}
}

func buildErrorLogFunc(w *Webhook, logger v1.Logger) func(error) {
	return func(err error) {
		logger.WithValues(map[string]interface{}{
			"url":          w.URL,
			"method":       w.Method,
			"content_type": w.ContentType,
		}).Error(err, "error executing webhook")
	}
}

// ToListener creates a newsman Listener from a Webhook.
func (w *Webhook) ToListener(logger v1.Logger) newsman.Listener {
	return newsman.NewWebhookListener(
		buildErrorLogFunc(w, logger),
		&newsman.WebhookConfig{
			Method:      w.Method,
			URL:         w.URL,
			ContentType: w.ContentType,
		},
		&newsman.ListenerConfig{
			Events:    w.Events,
			DataTypes: w.DataTypes,
			Topics:    w.Topics,
		},
	)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildWebhookTypeDefinitions()

		expected := `
package example

import (
	"context"
	"net/http"
)

type (
	// Webhook represents a webhook listener, an endpoint to send an HTTP request to upon an event.
	Webhook struct {
		ID            uint64   ` + "`" + `json:"id"` + "`" + `
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ContentType   string   ` + "`" + `json:"contentType"` + "`" + `
		URL           string   ` + "`" + `json:"url"` + "`" + `
		Method        string   ` + "`" + `json:"method"` + "`" + `
		Events        []string ` + "`" + `json:"events"` + "`" + `
		DataTypes     []string ` + "`" + `json:"dataTypes"` + "`" + `
		Topics        []string ` + "`" + `json:"topics"` + "`" + `
		CreatedOn     uint64   ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn *uint64  ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn    *uint64  ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	// WebhookCreationInput represents what a user could set as input for creating a webhook.
	WebhookCreationInput struct {
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ContentType   string   ` + "`" + `json:"contentType"` + "`" + `
		URL           string   ` + "`" + `json:"url"` + "`" + `
		Method        string   ` + "`" + `json:"method"` + "`" + `
		Events        []string ` + "`" + `json:"events"` + "`" + `
		DataTypes     []string ` + "`" + `json:"dataTypes"` + "`" + `
		Topics        []string ` + "`" + `json:"topics"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"-"` + "`" + `
	}

	// WebhookUpdateInput represents what a user could set as input for updating a webhook.
	WebhookUpdateInput struct {
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ContentType   string   ` + "`" + `json:"contentType"` + "`" + `
		URL           string   ` + "`" + `json:"url"` + "`" + `
		Method        string   ` + "`" + `json:"method"` + "`" + `
		Events        []string ` + "`" + `json:"events"` + "`" + `
		DataTypes     []string ` + "`" + `json:"dataTypes"` + "`" + `
		Topics        []string ` + "`" + `json:"topics"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"-"` + "`" + `
	}

	// WebhookList represents a list of webhooks.
	WebhookList struct {
		Pagination
		Webhooks []Webhook ` + "`" + `json:"webhooks"` + "`" + `
	}

	// WebhookDataManager describes a structure capable of storing webhooks.
	WebhookDataManager interface {
		GetWebhook(ctx context.Context, webhookID, userID uint64) (*Webhook, error)
		GetAllWebhooksCount(ctx context.Context) (uint64, error)
		GetWebhooks(ctx context.Context, userID uint64, filter *QueryFilter) (*WebhookList, error)
		GetAllWebhooks(ctx context.Context) (*WebhookList, error)
		CreateWebhook(ctx context.Context, input *WebhookCreationInput) (*Webhook, error)
		UpdateWebhook(ctx context.Context, updated *Webhook) error
		ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error
	}

	// WebhookDataServer describes a structure capable of serving traffic related to webhooks.
	WebhookDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookUpdate(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildWebhookUpdate()

		expected := `
package example

import ()

// Update merges an WebhookCreationInput with an Webhook.
func (w *Webhook) Update(input *WebhookUpdateInput) {
	if input.Name != "" {
		w.Name = input.Name
	}
	if input.ContentType != "" {
		w.ContentType = input.ContentType
	}
	if input.URL != "" {
		w.URL = input.URL
	}
	if input.Method != "" {
		w.Method = input.Method
	}

	if input.Events != nil && len(input.Events) > 0 {
		w.Events = input.Events
	}
	if input.DataTypes != nil && len(input.DataTypes) > 0 {
		w.DataTypes = input.DataTypes
	}
	if input.Topics != nil && len(input.Topics) > 0 {
		w.Topics = input.Topics
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookbuildErrorLogFunc(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildWebhookbuildErrorLogFunc()

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

func buildErrorLogFunc(w *Webhook, logger v1.Logger) func(error) {
	return func(err error) {
		logger.WithValues(map[string]interface{}{
			"url":          w.URL,
			"method":       w.Method,
			"content_type": w.ContentType,
		}).Error(err, "error executing webhook")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookToListener(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildWebhookToListener()

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ToListener creates a newsman Listener from a Webhook.
func (w *Webhook) ToListener(logger v1.Logger) newsman.Listener {
	return newsman.NewWebhookListener(
		buildErrorLogFunc(w, logger),
		&newsman.WebhookConfig{
			Method:      w.Method,
			URL:         w.URL,
			ContentType: w.ContentType,
		},
		&newsman.ListenerConfig{
			Events:    w.Events,
			DataTypes: w.DataTypes,
			Topics:    w.Topics,
		},
	)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
