package models

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/oauth2.v3"
)

var OAuth2ClientKey ContextKey = "oauth2_client"

type (
	OAuth2ClientDataManager interface {
		GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*OAuth2Client, error)
		GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*OAuth2Client, error)
		GetAllOAuth2ClientCount(ctx context.Context) (uint64, error)
		GetOAuth2ClientCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetOAuth2Clients(ctx context.Context, filter *QueryFilter, userID uint64) (*OAuth2ClientList, error)
		GetAllOAuth2Clients(ctx context.Context) ([]*OAuth2Client, error)
		GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*OAuth2Client, error)
		CreateOAuth2Client(ctx context.Context, input *OAuth2ClientCreationInput) (*OAuth2Client, error)
		UpdateOAuth2Client(ctx context.Context, updated *OAuth2Client) error
		ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error
	}
	OAuth2ClientDataServer interface {
		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
		CreationInputMiddleware(next http.Handler) http.Handler
		OAuth2ClientInfoMiddleware(next http.Handler) http.Handler
		ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*OAuth2Client, error)
		HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error
		HandleTokenRequest(res http.ResponseWriter, req *http.Request) error
	}
	OAuth2Client struct {
		ID              uint64
		Name            string
		ClientID        string
		ClientSecret    string
		RedirectURI     string
		Scopes          []string
		ImplicitAllowed bool
		CreatedOn       uint64
		UpdatedOn       *uint64
		ArchivedOn      *uint64
		BelongsTo       uint64
	}
	OAuth2ClientList struct {
		Pagination
		Clients []OAuth2Client
	}
	OAuth2ClientCreationInput struct {
		UserLoginInput
		Name         string
		ClientID     string
		ClientSecret string
		RedirectURI  string
		BelongsTo    uint64
		Scopes       []string
	}
	OAuth2ClientUpdateInput struct {
		RedirectURI string
		Scopes      []string
	}
)

var _ oauth2.ClientInfo = (*OAuth2Client)(nil)

// GetID returns the client ID. NOTE: I believe this is implemented for the above interface spec (oauth2.ClientInfo)
func (c *OAuth2Client) GetID() string {
	return c.ClientID
}

// GetSecret returns the ClientSecret
func (c *OAuth2Client) GetSecret() string {
	return c.ClientSecret
}

// GetDomain returns the client's domain
func (c *OAuth2Client) GetDomain() string {
	return c.RedirectURI
}

// GetUserID returns the client's UserID
func (c *OAuth2Client) GetUserID() string {
	return strconv.FormatUint(c.BelongsTo, 10)
}

// HasScope returns whether or not the provided scope is included in the scope list
func (c *OAuth2Client) HasScope(scope string) (found bool) {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return false
	}
	if c != nil && c.Scopes != nil {
		for _, s := range c.Scopes {
			if strings.TrimSpace(strings.ToLower(s)) == strings.TrimSpace(strings.ToLower(scope)) || strings.TrimSpace(s) == "*" {
				return true
			}
		}
	}
	return false
}
