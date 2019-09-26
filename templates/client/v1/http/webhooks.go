package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

var webhooksBasePath = "webhooks"

//
func (c *V1Client) BuildGetWebhookRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(id, 10))
	return http.NewRequest(http.MethodGet, uri, nil)
}

//
func (c *V1Client) GetWebhook(ctx context.Context, id uint64) (webhook *models.Webhook, err error) {
	req, err := c.BuildGetWebhookRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.retrieve(ctx, req, &webhook)
	return webhook, err
}

//
func (c *V1Client) BuildGetWebhooksRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), webhooksBasePath)
	return http.NewRequest(http.MethodGet, uri, nil)
}

//
func (c *V1Client) GetWebhooks(ctx context.Context, filter *models.QueryFilter) (webhooks *models.WebhookList, err error) {
	req, err := c.BuildGetWebhooksRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.retrieve(ctx, req, &webhooks)
	return webhooks, err
}

//
func (c *V1Client) BuildCreateWebhookRequest(ctx context.Context, body *models.WebhookCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, webhooksBasePath)
	return c.buildDataRequest(http.MethodPost, uri, body)
}

//
func (c *V1Client) CreateWebhook(ctx context.Context, input *models.WebhookCreationInput) (webhook *models.Webhook, err error) {
	req, err := c.BuildCreateWebhookRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.executeRequest(ctx, req, &webhook)
	return webhook, err
}

//
func (c *V1Client) BuildUpdateWebhookRequest(ctx context.Context, updated *models.Webhook) (*http.Request, error) {
	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(updated.ID, 10))
	return c.buildDataRequest(http.MethodPut, uri, updated)
}

//
func (c *V1Client) UpdateWebhook(ctx context.Context, updated *models.Webhook) error {
	req, err := c.BuildUpdateWebhookRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	return c.executeRequest(ctx, req, &updated)
}

//
func (c *V1Client) BuildArchiveWebhookRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, webhooksBasePath, strconv.FormatUint(id, 10))
	return http.NewRequest(http.MethodDelete, uri, nil)
}

//
func (c *V1Client) ArchiveWebhook(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveWebhookRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	return c.executeRequest(ctx, req, nil)
}
