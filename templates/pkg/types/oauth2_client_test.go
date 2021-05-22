package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientDotGo(proj)

		expected := `
package example

import (
	"context"
	oauth2v3 "gopkg.in/oauth2.v3"
	"net/http"
	"strconv"
	"strings"
)

const (
	// OAuth2ClientKey is a ContextKey for use with contexts involving OAuth2 clients.
	OAuth2ClientKey ContextKey = "oauth2_client"
)

type (
	// OAuth2Client represents a user-authorized API client
	OAuth2Client struct {
		ID              uint64   ` + "`" + `json:"id"` + "`" + `
		Name            string   ` + "`" + `json:"name"` + "`" + `
		ClientID        string   ` + "`" + `json:"clientID"` + "`" + `
		ClientSecret    string   ` + "`" + `json:"clientSecret"` + "`" + `
		RedirectURI     string   ` + "`" + `json:"redirectURI"` + "`" + `
		Scopes          []string ` + "`" + `json:"scopes"` + "`" + `
		ImplicitAllowed bool     ` + "`" + `json:"implicitAllowed"` + "`" + `
		CreatedOn       uint64   ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn   *uint64  ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn      *uint64  ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToUser   uint64   ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	// OAuth2ClientList is a response struct containing a list of OAuth2Clients.
	OAuth2ClientList struct {
		Pagination
		Clients []OAuth2Client ` + "`" + `json:"clients"` + "`" + `
	}

	// OAuth2ClientCreationInput is a struct for use when creating OAuth2 clients.
	OAuth2ClientCreationInput struct {
		UserLoginInput
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ClientID      string   ` + "`" + `json:"-"` + "`" + `
		ClientSecret  string   ` + "`" + `json:"-"` + "`" + `
		RedirectURI   string   ` + "`" + `json:"redirectURI"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"-"` + "`" + `
		Scopes        []string ` + "`" + `json:"scopes"` + "`" + `
	}

	// OAuth2ClientUpdateInput is a struct for use when updating OAuth2 clients.
	OAuth2ClientUpdateInput struct {
		RedirectURI string   ` + "`" + `json:"redirectURI"` + "`" + `
		Scopes      []string ` + "`" + `json:"scopes"` + "`" + `
	}

	// OAuth2ClientDataManager handles OAuth2 clients.
	OAuth2ClientDataManager interface {
		GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*OAuth2Client, error)
		GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*OAuth2Client, error)
		GetAllOAuth2ClientCount(ctx context.Context) (uint64, error)
		GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *QueryFilter) (*OAuth2ClientList, error)
		CreateOAuth2Client(ctx context.Context, input *OAuth2ClientCreationInput) (*OAuth2Client, error)
		UpdateOAuth2Client(ctx context.Context, updated *OAuth2Client) error
		ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error
	}

	// OAuth2ClientDataServer describes a structure capable of serving traffic related to oauth2 clients.
	OAuth2ClientDataServer interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		// There is deliberately no update function.
		ArchiveHandler(res http.ResponseWriter, req *http.Request)

		CreationInputMiddleware(next http.Handler) http.Handler
		OAuth2ClientInfoMiddleware(next http.Handler) http.Handler
		ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*OAuth2Client, error)

		// wrappers for our implementation library.
		HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error
		HandleTokenRequest(res http.ResponseWriter, req *http.Request) error
	}
)

var _ oauth2v3.ClientInfo = (*OAuth2Client)(nil)

// GetID returns the client ID. NOTE: I believe this is implemented for the above interface spec (oauth2.ClientInfo)
func (c *OAuth2Client) GetID() string {
	return c.ClientID
}

// GetSecret returns the ClientSecret.
func (c *OAuth2Client) GetSecret() string {
	return c.ClientSecret
}

// GetDomain returns the client's domain.
func (c *OAuth2Client) GetDomain() string {
	return c.RedirectURI
}

// GetUserID returns the client's UserID.
func (c *OAuth2Client) GetUserID() string {
	return strconv.FormatUint(c.BelongsToUser, 10)
}

// HasScope returns whether or not the provided scope is included in the scope list.
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
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientsConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientsConstantDefs()

		expected := `
package example

import ()

const (
	// OAuth2ClientKey is a ContextKey for use with contexts involving OAuth2 clients.
	OAuth2ClientKey ContextKey = "oauth2_client"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientsTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientsTypeDefs()

		expected := `
package example

import (
	"context"
	"net/http"
)

type (
	// OAuth2Client represents a user-authorized API client
	OAuth2Client struct {
		ID              uint64   ` + "`" + `json:"id"` + "`" + `
		Name            string   ` + "`" + `json:"name"` + "`" + `
		ClientID        string   ` + "`" + `json:"clientID"` + "`" + `
		ClientSecret    string   ` + "`" + `json:"clientSecret"` + "`" + `
		RedirectURI     string   ` + "`" + `json:"redirectURI"` + "`" + `
		Scopes          []string ` + "`" + `json:"scopes"` + "`" + `
		ImplicitAllowed bool     ` + "`" + `json:"implicitAllowed"` + "`" + `
		CreatedOn       uint64   ` + "`" + `json:"createdOn"` + "`" + `
		LastUpdatedOn   *uint64  ` + "`" + `json:"lastUpdatedOn"` + "`" + `
		ArchivedOn      *uint64  ` + "`" + `json:"archivedOn"` + "`" + `
		BelongsToUser   uint64   ` + "`" + `json:"belongsToUser"` + "`" + `
	}

	// OAuth2ClientList is a response struct containing a list of OAuth2Clients.
	OAuth2ClientList struct {
		Pagination
		Clients []OAuth2Client ` + "`" + `json:"clients"` + "`" + `
	}

	// OAuth2ClientCreationInput is a struct for use when creating OAuth2 clients.
	OAuth2ClientCreationInput struct {
		UserLoginInput
		Name          string   ` + "`" + `json:"name"` + "`" + `
		ClientID      string   ` + "`" + `json:"-"` + "`" + `
		ClientSecret  string   ` + "`" + `json:"-"` + "`" + `
		RedirectURI   string   ` + "`" + `json:"redirectURI"` + "`" + `
		BelongsToUser uint64   ` + "`" + `json:"-"` + "`" + `
		Scopes        []string ` + "`" + `json:"scopes"` + "`" + `
	}

	// OAuth2ClientUpdateInput is a struct for use when updating OAuth2 clients.
	OAuth2ClientUpdateInput struct {
		RedirectURI string   ` + "`" + `json:"redirectURI"` + "`" + `
		Scopes      []string ` + "`" + `json:"scopes"` + "`" + `
	}

	// OAuth2ClientDataManager handles OAuth2 clients.
	OAuth2ClientDataManager interface {
		GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*OAuth2Client, error)
		GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*OAuth2Client, error)
		GetAllOAuth2ClientCount(ctx context.Context) (uint64, error)
		GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *QueryFilter) (*OAuth2ClientList, error)
		CreateOAuth2Client(ctx context.Context, input *OAuth2ClientCreationInput) (*OAuth2Client, error)
		UpdateOAuth2Client(ctx context.Context, updated *OAuth2Client) error
		ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error
	}

	// OAuth2ClientDataServer describes a structure capable of serving traffic related to oauth2 clients.
	OAuth2ClientDataServer interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		// There is deliberately no update function.
		ArchiveHandler(res http.ResponseWriter, req *http.Request)

		CreationInputMiddleware(next http.Handler) http.Handler
		OAuth2ClientInfoMiddleware(next http.Handler) http.Handler
		ExtractOAuth2ClientFromRequest(ctx context.Context, req *http.Request) (*OAuth2Client, error)

		// wrappers for our implementation library.
		HandleAuthorizeRequest(res http.ResponseWriter, req *http.Request) error
		HandleTokenRequest(res http.ResponseWriter, req *http.Request) error
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2PkgInterfaceImplementationStatement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2PkgInterfaceImplementationStatement()

		expected := `
package example

import (
	oauth2v3 "gopkg.in/oauth2.v3"
)

var _ oauth2v3.ClientInfo = (*OAuth2Client)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDotGetID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDotGetID()

		expected := `
package example

import ()

// GetID returns the client ID. NOTE: I believe this is implemented for the above interface spec (oauth2.ClientInfo)
func (c *OAuth2Client) GetID() string {
	return c.ClientID
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDotGetSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDotGetSecret()

		expected := `
package example

import ()

// GetSecret returns the ClientSecret.
func (c *OAuth2Client) GetSecret() string {
	return c.ClientSecret
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDotGetDomain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDotGetDomain()

		expected := `
package example

import ()

// GetDomain returns the client's domain.
func (c *OAuth2Client) GetDomain() string {
	return c.RedirectURI
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDotGetUserID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDotGetUserID()

		expected := `
package example

import (
	"strconv"
)

// GetUserID returns the client's UserID.
func (c *OAuth2Client) GetUserID() string {
	return strconv.FormatUint(c.BelongsToUser, 10)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDotHasScope(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDotHasScope()

		expected := `
package example

import (
	"strings"
)

// HasScope returns whether or not the provided scope is included in the scope list.
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
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
