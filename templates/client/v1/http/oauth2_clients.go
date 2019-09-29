package client

import (
	"context"
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
	"net/http"
	"strconv"
)

var oauth2ClientsBasePath = "oauth2/clients"

//
func (c *V1Client) BuildGetOAuth2ClientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, oauth2ClientsBasePath, strconv.FormatUint(id, 10))
	return http.NewRequest(http.MethodGet, uri, nil)
}

//
func (c *V1Client) GetOAuth2Client(ctx context.Context, id uint64) (oauth2Client *models.OAuth2Client, err error) {
	req, err := c.BuildGetOAuth2ClientRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	err = c.retrieve(ctx, req, &oauth2Client)
	return oauth2Client, err
}

//
func (c *V1Client) BuildGetOAuth2ClientsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), oauth2ClientsBasePath)
	return http.NewRequest(http.MethodGet, uri, nil)
}

//
func (c *V1Client) GetOAuth2Clients(ctx context.Context, filter *models.QueryFilter) (*models.OAuth2ClientList, error) {
	req, err := c.BuildGetOAuth2ClientsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	var oauth2Clients *models.OAuth2ClientList
	err = c.retrieve(ctx, req, &oauth2Clients)
	return oauth2Clients, err
}

//
func (c *V1Client) BuildCreateOAuth2ClientRequest(ctx context.Context, cookie *http.Cookie, body *models.OAuth2ClientCreationInput) (*http.Request, error) {
	uri := c.buildVersionlessURL(nil, "oauth2", "client")
	req, err := c.buildDataRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)
	return req, nil
}

//
func (c *V1Client) CreateOAuth2Client(ctx context.Context, cookie *http.Cookie, input *models.OAuth2ClientCreationInput) (*models.OAuth2Client, error) {
	var oauth2Client *models.OAuth2Client
	if cookie == nil {
		return nil, errors.New("cookie required for request")
	}
	req, err := c.BuildCreateOAuth2ClientRequest(ctx, cookie, input)
	if err != nil {
		return nil, err
	}
	res, err := c.executeRawRequest(ctx, c.plainClient, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resErr := unmarshalBody(res, &oauth2Client); resErr != nil {
		return nil, errors.Wrap(resErr, "loading response from server")
	}
	return oauth2Client, nil
}

//
func (c *V1Client) BuildArchiveOAuth2ClientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, oauth2ClientsBasePath, strconv.FormatUint(id, 10))
	return http.NewRequest(http.MethodDelete, uri, nil)
}

//
func (c *V1Client) ArchiveOAuth2Client(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveOAuth2ClientRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	return c.executeRequest(ctx, req, nil)
}
