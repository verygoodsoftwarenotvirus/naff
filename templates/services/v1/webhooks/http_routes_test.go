package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_httpRoutesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := httpRoutesDotGo(proj)

		expected := `
package example

import (
	"database/sql"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
	"net/url"
	"strings"
)

const (
	// URIParamKey is a standard string that we'll use to refer to webhook IDs with.
	URIParamKey = "webhookID"
)

// validateWebhook does some validation on a WebhookCreationInput and returns an error if anything runs foul.
func validateWebhook(input *v1.WebhookCreationInput) error {
	_, err := url.Parse(input.URL)
	if err != nil {
		return fmt.Errorf("invalid URL provided: %w", err)
	}

	input.Method = strings.ToUpper(input.Method)
	switch input.Method {
	// allowed methods.
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead:
		break
	default:
		return fmt.Errorf("invalid method provided: %q", input.Method)
	}

	return nil
}

// CreateHandler is our webhook creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// figure out who this is all for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// try to pluck the parsed input from the request context.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*v1.WebhookCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	input.BelongsToUser = userID

	// ensure everythings on the up-and-up
	if err := validateWebhook(input); err != nil {
		logger.Info("invalid method provided")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// create the webhook.
	wh, err := s.webhookDataManager.CreateWebhook(ctx, input)
	if err != nil {
		logger.Error(err, "error creating webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify the relevant parties.
	tracing.AttachWebhookIDToSpan(span, wh.ID)
	s.webhookCounter.Increment(ctx)
	s.eventManager.Report(newsman.Event{
		EventType: string(v1.Create),
		Data:      wh,
		Topics:    []string{topicName},
	})

	l := wh.ToListener(logger)
	s.eventManager.TuneIn(l)

	// let everybody know we're good.
	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, wh); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// figure out how specific we need to be.
	filter := v1.ExtractQueryFilter(req)

	// figure out who this is all for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// find the webhooks.
	webhooks, err := s.webhookDataManager.GetWebhooks(ctx, userID, filter)
	if err == sql.ErrNoRows {
		webhooks = &v1.WebhookList{
			Webhooks: []v1.Webhook{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching webhooks")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode the response.
	if err = s.encoderDecoder.EncodeResponse(res, webhooks); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ReadHandler returns a GET handler that returns an webhook.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine relevant user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant webhook ID.
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachWebhookIDToSpan(span, webhookID)
	logger = logger.WithValue("webhook_id", webhookID)

	// fetch the webhook from the database.
	x, err := s.webhookDataManager.GetWebhook(ctx, webhookID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("No rows found in webhook database")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "Error fetching webhook from webhook database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode the response.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}

// UpdateHandler returns a handler that updates an webhook.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine relevant user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant webhook ID.
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachWebhookIDToSpan(span, webhookID)
	logger = logger.WithValue("webhook_id", webhookID)

	// fetch parsed creation input from request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*v1.WebhookUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// fetch the webhook in question.
	wh, err := s.webhookDataManager.GetWebhook(ctx, webhookID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("no rows found for webhook")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update it.
	wh.Update(input)

	// save the update in the database.
	if err = s.webhookDataManager.UpdateWebhook(ctx, wh); err != nil {
		logger.Error(err, "error encountered updating webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify the relevant parties.
	s.eventManager.Report(newsman.Event{
		EventType: string(v1.Update),
		Data:      wh,
		Topics:    []string{topicName},
	})

	// let everybody know we're good.
	if err = s.encoderDecoder.EncodeResponse(res, wh); err != nil {
		logger.Error(err, "encoding response")
	}
}

// ArchiveHandler returns a handler that archives an webhook.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "delete_route")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine relevant user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant webhook ID.
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachWebhookIDToSpan(span, webhookID)
	logger = logger.WithValue("webhook_id", webhookID)

	// do the deed.
	err := s.webhookDataManager.ArchiveWebhook(ctx, webhookID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("no rows found for webhook")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// let the interested parties know.
	s.webhookCounter.Decrement(ctx)
	s.eventManager.Report(newsman.Event{
		EventType: string(v1.Archive),
		Data:      v1.Webhook{ID: webhookID},
		Topics:    []string{topicName},
	})

	// let everybody go home.
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesConstDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildWebhookHTTPRoutesConstDefs()

		expected := `
package example

import ()

const (
	// URIParamKey is a standard string that we'll use to refer to webhook IDs with.
	URIParamKey = "webhookID"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesValidateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWebhookHTTPRoutesValidateWebhook(proj)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"net/url"
	"strings"
)

// validateWebhook does some validation on a WebhookCreationInput and returns an error if anything runs foul.
func validateWebhook(input *v1.WebhookCreationInput) error {
	_, err := url.Parse(input.URL)
	if err != nil {
		return fmt.Errorf("invalid URL provided: %w", err)
	}

	input.Method = strings.ToUpper(input.Method)
	switch input.Method {
	// allowed methods.
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead:
		break
	default:
		return fmt.Errorf("invalid method provided: %q", input.Method)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesCreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWebhookHTTPRoutesCreateHandler(proj)

		expected := `
package example

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

// CreateHandler is our webhook creation route.
func (s *Service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "CreateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// figure out who this is all for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// try to pluck the parsed input from the request context.
	input, ok := ctx.Value(createMiddlewareCtxKey).(*v1.WebhookCreationInput)
	if !ok {
		logger.Info("valid input not attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	input.BelongsToUser = userID

	// ensure everythings on the up-and-up
	if err := validateWebhook(input); err != nil {
		logger.Info("invalid method provided")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// create the webhook.
	wh, err := s.webhookDataManager.CreateWebhook(ctx, input)
	if err != nil {
		logger.Error(err, "error creating webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify the relevant parties.
	tracing.AttachWebhookIDToSpan(span, wh.ID)
	s.webhookCounter.Increment(ctx)
	s.eventManager.Report(newsman.Event{
		EventType: string(v1.Create),
		Data:      wh,
		Topics:    []string{topicName},
	})

	l := wh.ToListener(logger)
	s.eventManager.TuneIn(l)

	// let everybody know we're good.
	res.WriteHeader(http.StatusCreated)
	if err = s.encoderDecoder.EncodeResponse(res, wh); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesListHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWebhookHTTPRoutesListHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// ListHandler is our list route.
func (s *Service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ListHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// figure out how specific we need to be.
	filter := v1.ExtractQueryFilter(req)

	// figure out who this is all for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue("user_id", userID)
	tracing.AttachUserIDToSpan(span, userID)

	// find the webhooks.
	webhooks, err := s.webhookDataManager.GetWebhooks(ctx, userID, filter)
	if err == sql.ErrNoRows {
		webhooks = &v1.WebhookList{
			Webhooks: []v1.Webhook{},
		}
	} else if err != nil {
		logger.Error(err, "error encountered fetching webhooks")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode the response.
	if err = s.encoderDecoder.EncodeResponse(res, webhooks); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWebhookHTTPRoutesReadHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// ReadHandler returns a GET handler that returns an webhook.
func (s *Service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "ReadHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine relevant user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant webhook ID.
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachWebhookIDToSpan(span, webhookID)
	logger = logger.WithValue("webhook_id", webhookID)

	// fetch the webhook from the database.
	x, err := s.webhookDataManager.GetWebhook(ctx, webhookID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("No rows found in webhook database")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "Error fetching webhook from webhook database")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode the response.
	if err = s.encoderDecoder.EncodeResponse(res, x); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesUpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWebhookHTTPRoutesUpdateHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

// UpdateHandler returns a handler that updates an webhook.
func (s *Service) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "UpdateHandler")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine relevant user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant webhook ID.
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachWebhookIDToSpan(span, webhookID)
	logger = logger.WithValue("webhook_id", webhookID)

	// fetch parsed creation input from request context.
	input, ok := ctx.Value(updateMiddlewareCtxKey).(*v1.WebhookUpdateInput)
	if !ok {
		logger.Info("no input attached to request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// fetch the webhook in question.
	wh, err := s.webhookDataManager.GetWebhook(ctx, webhookID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("no rows found for webhook")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered getting webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update it.
	wh.Update(input)

	// save the update in the database.
	if err = s.webhookDataManager.UpdateWebhook(ctx, wh); err != nil {
		logger.Error(err, "error encountered updating webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// notify the relevant parties.
	s.eventManager.Report(newsman.Event{
		EventType: string(v1.Update),
		Data:      wh,
		Topics:    []string{topicName},
	})

	// let everybody know we're good.
	if err = s.encoderDecoder.EncodeResponse(res, wh); err != nil {
		logger.Error(err, "encoding response")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhookHTTPRoutesArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWebhookHTTPRoutesArchiveHandler(proj)

		expected := `
package example

import (
	"database/sql"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
	"net/http"
)

// ArchiveHandler returns a handler that archives an webhook.
func (s *Service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := tracing.StartSpan(req.Context(), "delete_route")
	defer span.End()

	logger := s.logger.WithRequest(req)

	// determine relevant user ID.
	userID := s.userIDFetcher(req)
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue("user_id", userID)

	// determine relevant webhook ID.
	webhookID := s.webhookIDFetcher(req)
	tracing.AttachWebhookIDToSpan(span, webhookID)
	logger = logger.WithValue("webhook_id", webhookID)

	// do the deed.
	err := s.webhookDataManager.ArchiveWebhook(ctx, webhookID, userID)
	if err == sql.ErrNoRows {
		logger.Debug("no rows found for webhook")
		res.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error(err, "error encountered deleting webhook")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// let the interested parties know.
	s.webhookCounter.Decrement(ctx)
	s.eventManager.Report(newsman.Event{
		EventType: string(v1.Archive),
		Data:      v1.Webhook{ID: webhookID},
		Topics:    []string{topicName},
	})

	// let everybody go home.
	res.WriteHeader(http.StatusNoContent)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
