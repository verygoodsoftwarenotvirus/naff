package querier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhooksDotGo(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var _ v1.WebhookDataManager = (*Client)(nil)

// GetWebhook fetches a webhook from the database.
func (c *Client) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhook")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	c.logger.WithValues(map[string]interface{}{
		"webhook_id": webhookID,
		"user_id":    userID,
	}).Debug("GetWebhook called")

	return c.querier.GetWebhook(ctx, webhookID, userID)
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (c *Client) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhooks")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetWebhookCount called")

	return c.querier.GetWebhooks(ctx, userID, filter)
}

// GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (c *Client) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllWebhooks")
	defer span.End()

	c.logger.Debug("GetWebhookCount called")

	return c.querier.GetAllWebhooks(ctx)
}

// GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter.
func (c *Client) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllWebhooksCount")
	defer span.End()

	c.logger.Debug("GetAllWebhooksCount called")

	return c.querier.GetAllWebhooksCount(ctx)
}

// CreateWebhook creates a webhook in a database.
func (c *Client) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateWebhook")
	defer span.End()

	tracing.AttachUserIDToSpan(span, input.BelongsToUser)
	c.logger.WithValue("user_id", input.BelongsToUser).Debug("CreateWebhook called")

	return c.querier.CreateWebhook(ctx, input)
}

// UpdateWebhook updates a particular webhook.
// NOTE: this function expects the provided input to have a non-zero ID.
func (c *Client) UpdateWebhook(ctx context.Context, input *v1.Webhook) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateWebhook")
	defer span.End()

	tracing.AttachWebhookIDToSpan(span, input.ID)
	tracing.AttachUserIDToSpan(span, input.BelongsToUser)

	c.logger.WithValue("webhook_id", input.ID).Debug("UpdateWebhook called")

	return c.querier.UpdateWebhook(ctx, input)
}

// ArchiveWebhook archives a webhook from the database.
func (c *Client) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveWebhook")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	c.logger.WithValues(map[string]interface{}{
		"webhook_id": webhookID,
		"user_id":    userID,
	}).Debug("ArchiveWebhook called")

	return c.querier.ArchiveWebhook(ctx, webhookID, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetWebhook(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhook fetches a webhook from the database.
func (c *Client) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhook")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	c.logger.WithValues(map[string]interface{}{
		"webhook_id": webhookID,
		"user_id":    userID,
	}).Debug("GetWebhook called")

	return c.querier.GetWebhook(ctx, webhookID, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetWebhooks(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (c *Client) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhooks")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetWebhookCount called")

	return c.querier.GetWebhooks(ctx, userID, filter)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetAllWebhooks(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (c *Client) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllWebhooks")
	defer span.End()

	c.logger.Debug("GetWebhookCount called")

	return c.querier.GetAllWebhooks(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllWebhooksCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetAllWebhooksCount(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter.
func (c *Client) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllWebhooksCount")
	defer span.End()

	c.logger.Debug("GetAllWebhooksCount called")

	return c.querier.GetAllWebhooksCount(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateWebhook(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateWebhook creates a webhook in a database.
func (c *Client) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateWebhook")
	defer span.End()

	tracing.AttachUserIDToSpan(span, input.BelongsToUser)
	c.logger.WithValue("user_id", input.BelongsToUser).Debug("CreateWebhook called")

	return c.querier.CreateWebhook(ctx, input)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUpdateWebhook(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateWebhook updates a particular webhook.
// NOTE: this function expects the provided input to have a non-zero ID.
func (c *Client) UpdateWebhook(ctx context.Context, input *v1.Webhook) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateWebhook")
	defer span.End()

	tracing.AttachWebhookIDToSpan(span, input.ID)
	tracing.AttachUserIDToSpan(span, input.BelongsToUser)

	c.logger.WithValue("webhook_id", input.ID).Debug("UpdateWebhook called")

	return c.querier.UpdateWebhook(ctx, input)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildArchiveWebhook(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveWebhook archives a webhook from the database.
func (c *Client) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveWebhook")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachWebhookIDToSpan(span, webhookID)

	c.logger.WithValues(map[string]interface{}{
		"webhook_id": webhookID,
		"user_id":    userID,
	}).Debug("ArchiveWebhook called")

	return c.querier.ArchiveWebhook(ctx, webhookID, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
