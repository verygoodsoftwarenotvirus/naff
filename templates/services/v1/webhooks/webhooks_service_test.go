package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksServiceDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhooksServiceDotGo(proj)

		expected := `
package example

import (
	"fmt"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts.
	createMiddlewareCtxKey v1.ContextKey = "webhook_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts.
	updateMiddlewareCtxKey v1.ContextKey = "webhook_update_input"

	counterName        metrics.CounterName = "webhooks"
	counterDescription string              = "the number of webhooks managed by the webhooks service"
	topicName          string              = "webhooks"
	serviceName        string              = "webhooks_service"
)

var (
	_ v1.WebhookDataServer = (*Service)(nil)
)

type (
	eventManager interface {
		newsman.Reporter

		TuneIn(newsman.Listener)
	}

	// Service handles TODO ListHandler webhooks.
	Service struct {
		logger             v11.Logger
		webhookCounter     metrics.UnitCounter
		webhookDataManager v1.WebhookDataManager
		userIDFetcher      UserIDFetcher
		webhookIDFetcher   WebhookIDFetcher
		encoderDecoder     encoding.EncoderDecoder
		eventManager       eventManager
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// WebhookIDFetcher is a function that fetches webhook IDs.
	WebhookIDFetcher func(*http.Request) uint64
)

// ProvideWebhooksService builds a new WebhooksService.
func ProvideWebhooksService(
	logger v11.Logger,
	webhookDataManager v1.WebhookDataManager,
	userIDFetcher UserIDFetcher,
	webhookIDFetcher WebhookIDFetcher,
	encoder encoding.EncoderDecoder,
	webhookCounterProvider metrics.UnitCounterProvider,
	em *newsman.Newsman,
) (*Service, error) {
	webhookCounter, err := webhookCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:             logger.WithName(serviceName),
		webhookDataManager: webhookDataManager,
		encoderDecoder:     encoder,
		webhookCounter:     webhookCounter,
		userIDFetcher:      userIDFetcher,
		webhookIDFetcher:   webhookIDFetcher,
		eventManager:       em,
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhooksServiceConstDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWebhooksServiceConstDefs(proj)

		expected := `
package example

import (
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts.
	createMiddlewareCtxKey v1.ContextKey = "webhook_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to webhook input data in contexts.
	updateMiddlewareCtxKey v1.ContextKey = "webhook_update_input"

	counterName        metrics.CounterName = "webhooks"
	counterDescription string              = "the number of webhooks managed by the webhooks service"
	topicName          string              = "webhooks"
	serviceName        string              = "webhooks_service"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhooksServiceVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWebhooksServiceVarDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var (
	_ v1.WebhookDataServer = (*Service)(nil)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhooksServiceTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildWebhooksServiceTypeDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

type (
	eventManager interface {
		newsman.Reporter

		TuneIn(newsman.Listener)
	}

	// Service handles TODO ListHandler webhooks.
	Service struct {
		logger             v1.Logger
		webhookCounter     metrics.UnitCounter
		webhookDataManager v11.WebhookDataManager
		userIDFetcher      UserIDFetcher
		webhookIDFetcher   WebhookIDFetcher
		encoderDecoder     encoding.EncoderDecoder
		eventManager       eventManager
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// WebhookIDFetcher is a function that fetches webhook IDs.
	WebhookIDFetcher func(*http.Request) uint64
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideWebhooksService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideWebhooksService(proj)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideWebhooksService builds a new WebhooksService.
func ProvideWebhooksService(
	logger v1.Logger,
	webhookDataManager v11.WebhookDataManager,
	userIDFetcher UserIDFetcher,
	webhookIDFetcher WebhookIDFetcher,
	encoder encoding.EncoderDecoder,
	webhookCounterProvider metrics.UnitCounterProvider,
	em *newsman.Newsman,
) (*Service, error) {
	webhookCounter, err := webhookCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:             logger.WithName(serviceName),
		webhookDataManager: webhookDataManager,
		encoderDecoder:     encoder,
		webhookCounter:     webhookCounter,
		userIDFetcher:      userIDFetcher,
		webhookIDFetcher:   webhookIDFetcher,
		eventManager:       em,
	}

	return svc, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
