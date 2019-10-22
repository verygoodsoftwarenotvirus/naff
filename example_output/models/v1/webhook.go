package models

import (
	"context"
	"net/http"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

type (
	WebhookDataManager interface {
		GetWebhook(ctx context.Context, webhookID, userID uint64) (*Webhook, error)
		GetWebhookCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllWebhooksCount(ctx context.Context) (uint64, error)
		GetWebhooks(ctx context.Context, filter *QueryFilter, userID uint64) (*WebhookList, error)
		GetAllWebhooks(ctx context.Context) (*WebhookList, error)
		GetAllWebhooksForUser(ctx context.Context, userID uint64) ([]Webhook, error)
		CreateWebhook(ctx context.Context, input *WebhookCreationInput) (*Webhook, error)
		UpdateWebhook(ctx context.Context, updated *Webhook) error
		ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error
	}
	WebhookDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler
		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
	Webhook struct {
		ID          uint64
		Name        string
		ContentType string
		URL         string
		Method      string
		Events      []string
		DataTypes   []string
		Topics      []string
		CreatedOn   uint64
		UpdatedOn   *uint64
		ArchivedOn  *uint64
		BelongsTo   uint64
	}
	WebhookCreationInput struct {
		Name        string
		ContentType string
		URL         string
		Method      string
		Events      []string
		DataTypes   []string
		Topics      []string
		BelongsTo   uint64
	}
	WebhookUpdateInput struct {
		Name        string
		ContentType string
		URL         string
		Method      string
		Events      []string
		DataTypes   []string
		Topics      []string
		BelongsTo   uint64
	}
	WebhookList struct {
		Pagination
		Webhooks []Webhook
	}
)

// Update merges an WebhookCreationInput with an Webhook
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

func buildErrorLogFunc(w *Webhook, logger logging.Logger) error {
	return func(err error) {
		logger.WithValues(map[string]interface{}{
			"url":          w.URL,
			"method":       w.Method,
			"content_type": w.ContentType,
		}).Error(err, "error executing webhook")
	}
}

// ToListener creates a newsman Listener from a Webhook
func (w *Webhook) ToListener(logger logging.Logger) newsman.Listener {
	return newsman.NewWebhookListener(buildErrorLogFunc(w, logger), &newsman.WebhookConfig{
		Method:      w.Method,
		URL:         w.URL,
		ContentType: w.ContentType,
	}, &newsman.ListenerConfig{
		Events:    w.Events,
		DataTypes: w.DataTypes,
		Topics:    w.Topics,
	})
}
