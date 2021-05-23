package load

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
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"math/rand"
	http1 "net/http"
)

// fetchRandomWebhook retrieves a random webhook from the list of available webhooks.
func fetchRandomWebhook(c *http.V1Client) *v1.Webhook {
	webhooks, err := c.GetWebhooks(context.Background(), nil)
	if err != nil || webhooks == nil || len(webhooks.Webhooks) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(webhooks.Webhooks))
	return &webhooks.Webhooks[randIndex]
}

func buildWebhookActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"GetWebhooks": {
			Name: "GetWebhooks",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				return c.BuildGetWebhooksRequest(ctx, nil)
			},
			Weight: 100,
		},
		"GetWebhook": {
			Name: "GetWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				if randomWebhook := fetchRandomWebhook(c); randomWebhook != nil {
					return c.BuildGetWebhookRequest(ctx, randomWebhook.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"CreateWebhook": {
			Name: "CreateWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				exampleInput := fake.BuildFakeWebhookCreationInput()
				return c.BuildCreateWebhookRequest(ctx, exampleInput)
			},
			Weight: 1,
		},
		"UpdateWebhook": {
			Name: "UpdateWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				if randomWebhook := fetchRandomWebhook(c); randomWebhook != nil {
					randomWebhook.Name = fake.BuildFakeWebhook().Name
					return c.BuildUpdateWebhookRequest(ctx, randomWebhook)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 50,
		},
		"ArchiveWebhook": {
			Name: "ArchiveWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				if randomWebhook := fetchRandomWebhook(c); randomWebhook != nil {
					return c.BuildArchiveWebhookRequest(ctx, randomWebhook.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 50,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFetchRandomWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildFetchRandomWebhook(proj)

		expected := `
package example

import (
	"context"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"math/rand"
)

// fetchRandomWebhook retrieves a random webhook from the list of available webhooks.
func fetchRandomWebhook(c *http.V1Client) *v1.Webhook {
	webhooks, err := c.GetWebhooks(context.Background(), nil)
	if err != nil || webhooks == nil || len(webhooks.Webhooks) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(webhooks.Webhooks))
	return &webhooks.Webhooks[randIndex]
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildWebhookActions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildWebhookActions(proj)

		expected := `
package example

import (
	"context"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	http1 "net/http"
)

func buildWebhookActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"GetWebhooks": {
			Name: "GetWebhooks",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				return c.BuildGetWebhooksRequest(ctx, nil)
			},
			Weight: 100,
		},
		"GetWebhook": {
			Name: "GetWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				if randomWebhook := fetchRandomWebhook(c); randomWebhook != nil {
					return c.BuildGetWebhookRequest(ctx, randomWebhook.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"CreateWebhook": {
			Name: "CreateWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				exampleInput := fake.BuildFakeWebhookCreationInput()
				return c.BuildCreateWebhookRequest(ctx, exampleInput)
			},
			Weight: 1,
		},
		"UpdateWebhook": {
			Name: "UpdateWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				if randomWebhook := fetchRandomWebhook(c); randomWebhook != nil {
					randomWebhook.Name = fake.BuildFakeWebhook().Name
					return c.BuildUpdateWebhookRequest(ctx, randomWebhook)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 50,
		},
		"ArchiveWebhook": {
			Name: "ArchiveWebhook",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				if randomWebhook := fetchRandomWebhook(c); randomWebhook != nil {
					return c.BuildArchiveWebhookRequest(ctx, randomWebhook.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 50,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
