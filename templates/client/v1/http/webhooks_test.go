package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := webhooksDotGo(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strconv"
)

const (
	webhooksBasePath = "webhooks"
)

// BuildGetWebhookRequest builds an HTTP request for fetching a webhook.
func (c *V1Client) BuildGetWebhookRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetWebhook retrieves a webhook.
func (c *V1Client) GetWebhook(ctx context.Context, id uint64) (webhook *v1.Webhook, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhook")
	defer span.End()

	req, err := c.BuildGetWebhookRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &webhook)
	return webhook, err
}

// BuildGetWebhooksRequest builds an HTTP request for fetching webhooks.
func (c *V1Client) BuildGetWebhooksRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetWebhooksRequest")
	defer span.End()

	uri := c.BuildURL(filter.ToValues(), webhooksBasePath)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetWebhooks gets a list of webhooks.
func (c *V1Client) GetWebhooks(ctx context.Context, filter *v1.QueryFilter) (webhooks *v1.WebhookList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhooks")
	defer span.End()

	req, err := c.BuildGetWebhooksRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &webhooks)
	return webhooks, err
}

// BuildCreateWebhookRequest builds an HTTP request for creating a webhook.
func (c *V1Client) BuildCreateWebhookRequest(ctx context.Context, body *v1.WebhookCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath)

	return c.buildDataRequest(ctx, http.MethodPost, uri, body)
}

// CreateWebhook creates a webhook.
func (c *V1Client) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (webhook *v1.Webhook, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateWebhook")
	defer span.End()

	req, err := c.BuildCreateWebhookRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &webhook)
	return webhook, err
}

// BuildUpdateWebhookRequest builds an HTTP request for updating a webhook.
func (c *V1Client) BuildUpdateWebhookRequest(ctx context.Context, updated *v1.Webhook) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(ctx, http.MethodPut, uri, updated)
}

// UpdateWebhook updates a webhook.
func (c *V1Client) UpdateWebhook(ctx context.Context, updated *v1.Webhook) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateWebhook")
	defer span.End()

	req, err := c.BuildUpdateWebhookRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveWebhookRequest builds an HTTP request for updating a webhook.
func (c *V1Client) BuildArchiveWebhookRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveWebhook archives a webhook.
func (c *V1Client) ArchiveWebhook(ctx context.Context, id uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveWebhook")
	defer span.End()

	req, err := c.BuildArchiveWebhookRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildGetWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildGetWebhookRequest builds an HTTP request for fetching a webhook.
func (c *V1Client) BuildGetWebhookRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildGetWebhook(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhook retrieves a webhook.
func (c *V1Client) GetWebhook(ctx context.Context, id uint64) (webhook *v1.Webhook, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhook")
	defer span.End()

	req, err := c.BuildGetWebhookRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &webhook)
	return webhook, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetWebhooksRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildGetWebhooksRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildGetWebhooksRequest builds an HTTP request for fetching webhooks.
func (c *V1Client) BuildGetWebhooksRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetWebhooksRequest")
	defer span.End()

	uri := c.BuildURL(filter.ToValues(), webhooksBasePath)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildGetWebhooks(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhooks gets a list of webhooks.
func (c *V1Client) GetWebhooks(ctx context.Context, filter *v1.QueryFilter) (webhooks *v1.WebhookList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetWebhooks")
	defer span.End()

	req, err := c.BuildGetWebhooksRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &webhooks)
	return webhooks, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildCreateWebhookRequest builds an HTTP request for creating a webhook.
func (c *V1Client) BuildCreateWebhookRequest(ctx context.Context, body *v1.WebhookCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath)

	return c.buildDataRequest(ctx, http.MethodPost, uri, body)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildCreateWebhook(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateWebhook creates a webhook.
func (c *V1Client) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (webhook *v1.Webhook, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateWebhook")
	defer span.End()

	req, err := c.BuildCreateWebhookRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &webhook)
	return webhook, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildUpdateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildUpdateWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strconv"
)

// BuildUpdateWebhookRequest builds an HTTP request for updating a webhook.
func (c *V1Client) BuildUpdateWebhookRequest(ctx context.Context, updated *v1.Webhook) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(ctx, http.MethodPut, uri, updated)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildUpdateWebhook(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateWebhook updates a webhook.
func (c *V1Client) UpdateWebhook(ctx context.Context, updated *v1.Webhook) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateWebhook")
	defer span.End()

	req, err := c.BuildUpdateWebhookRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildArchiveWebhookRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildArchiveWebhookRequest builds an HTTP request for updating a webhook.
func (c *V1Client) BuildArchiveWebhookRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveWebhookRequest")
	defer span.End()

	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildArchiveWebhook(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveWebhook archives a webhook.
func (c *V1Client) ArchiveWebhook(ctx context.Context, id uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveWebhook")
	defer span.End()

	req, err := c.BuildArchiveWebhookRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
