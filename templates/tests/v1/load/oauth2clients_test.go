package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsDotGo(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"math/rand"
	http1 "net/http"
	"time"
)

// fetchRandomOAuth2Client retrieves a random client from the list of available clients.
func fetchRandomOAuth2Client(c *http.V1Client) *v1.OAuth2Client {
	clientsRes, err := c.GetOAuth2Clients(context.Background(), nil)
	if err != nil || clientsRes == nil || len(clientsRes.Clients) <= 1 {
		return nil
	}

	var selectedClient *v1.OAuth2Client
	for selectedClient == nil {
		ri := rand.Intn(len(clientsRes.Clients))
		c := &clientsRes.Clients[ri]
		if c.ClientID != "FIXME" {
			selectedClient = c
		}
	}

	return selectedClient
}

func mustBuildCode(totpSecret string) string {
	code, err := totp.GenerateCode(totpSecret, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return code
}

func buildOAuth2ClientActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateOAuth2Client": {
			Name: "CreateOAuth2Client",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				ui := fake.BuildFakeUserCreationInput()
				u, err := c.CreateUser(ctx, ui)
				if err != nil {
					return c.BuildHealthCheckRequest(ctx)
				}

				uli := &v1.UserLoginInput{
					Username:  ui.Username,
					Password:  ui.Password,
					TOTPToken: mustBuildCode(u.TwoFactorSecret),
				}

				cookie, err := c.Login(ctx, uli)
				if err != nil {
					return c.BuildHealthCheckRequest(ctx)
				}

				req, err := c.BuildCreateOAuth2ClientRequest(
					ctx,
					cookie,
					&v1.OAuth2ClientCreationInput{
						UserLoginInput: *uli,
					},
				)
				return req, err
			},
			Weight: 100,
		},
		"GetOAuth2Client": {
			Name: "GetOAuth2Client",
			Action: func() (*http1.Request, error) {
				if randomOAuth2Client := fetchRandomOAuth2Client(c); randomOAuth2Client != nil {
					return c.BuildGetOAuth2ClientRequest(context.Background(), randomOAuth2Client.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetOAuth2ClientsForUser": {
			Name: "GetOAuth2ClientsForUser",
			Action: func() (*http1.Request, error) {
				return c.BuildGetOAuth2ClientsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFetchRandomOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildFetchRandomOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"math/rand"
)

// fetchRandomOAuth2Client retrieves a random client from the list of available clients.
func fetchRandomOAuth2Client(c *http.V1Client) *v1.OAuth2Client {
	clientsRes, err := c.GetOAuth2Clients(context.Background(), nil)
	if err != nil || clientsRes == nil || len(clientsRes.Clients) <= 1 {
		return nil
	}

	var selectedClient *v1.OAuth2Client
	for selectedClient == nil {
		ri := rand.Intn(len(clientsRes.Clients))
		c := &clientsRes.Clients[ri]
		if c.ClientID != "FIXME" {
			selectedClient = c
		}
	}

	return selectedClient
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMustBuildCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMustBuildCode()

		expected := `
package example

import (
	totp "github.com/pquerna/otp/totp"
	"time"
)

func mustBuildCode(totpSecret string) string {
	code, err := totp.GenerateCode(totpSecret, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return code
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildOAuth2ClientActions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildOAuth2ClientActions(proj)

		expected := `
package example

import (
	"context"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	http1 "net/http"
)

func buildOAuth2ClientActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateOAuth2Client": {
			Name: "CreateOAuth2Client",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()
				ui := fake.BuildFakeUserCreationInput()
				u, err := c.CreateUser(ctx, ui)
				if err != nil {
					return c.BuildHealthCheckRequest(ctx)
				}

				uli := &v1.UserLoginInput{
					Username:  ui.Username,
					Password:  ui.Password,
					TOTPToken: mustBuildCode(u.TwoFactorSecret),
				}

				cookie, err := c.Login(ctx, uli)
				if err != nil {
					return c.BuildHealthCheckRequest(ctx)
				}

				req, err := c.BuildCreateOAuth2ClientRequest(
					ctx,
					cookie,
					&v1.OAuth2ClientCreationInput{
						UserLoginInput: *uli,
					},
				)
				return req, err
			},
			Weight: 100,
		},
		"GetOAuth2Client": {
			Name: "GetOAuth2Client",
			Action: func() (*http1.Request, error) {
				if randomOAuth2Client := fetchRandomOAuth2Client(c); randomOAuth2Client != nil {
					return c.BuildGetOAuth2ClientRequest(context.Background(), randomOAuth2Client.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetOAuth2ClientsForUser": {
			Name: "GetOAuth2ClientsForUser",
			Action: func() (*http1.Request, error) {
				return c.BuildGetOAuth2ClientsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
