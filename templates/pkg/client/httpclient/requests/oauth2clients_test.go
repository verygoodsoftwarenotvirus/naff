package requests

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
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
	"strconv"
)

const (
	oauth2ClientsBasePath = "oauth2/clients"
)

// BuildGetOAuth2ClientRequest builds an HTTP request for fetching an OAuth2 client.
func (c *V1Client) BuildGetOAuth2ClientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetOAuth2ClientRequest")
	defer span.End()

	uri := c.BuildURL(nil, oauth2ClientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetOAuth2Client gets an OAuth2 client.
func (c *V1Client) GetOAuth2Client(ctx context.Context, id uint64) (oauth2Client *v1.OAuth2Client, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2Client")
	defer span.End()

	req, err := c.BuildGetOAuth2ClientRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &oauth2Client)
	return oauth2Client, err
}

// BuildGetOAuth2ClientsRequest builds an HTTP request for fetching a list of OAuth2 clients.
func (c *V1Client) BuildGetOAuth2ClientsRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetOAuth2ClientsRequest")
	defer span.End()

	uri := c.BuildURL(filter.ToValues(), oauth2ClientsBasePath)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetOAuth2Clients gets a list of OAuth2 clients.
func (c *V1Client) GetOAuth2Clients(ctx context.Context, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2Clients")
	defer span.End()

	req, err := c.BuildGetOAuth2ClientsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	var oauth2Clients *v1.OAuth2ClientList
	err = c.retrieve(ctx, req, &oauth2Clients)
	return oauth2Clients, err
}

// BuildCreateOAuth2ClientRequest builds an HTTP request for creating OAuth2 clients.
func (c *V1Client) BuildCreateOAuth2ClientRequest(
	ctx context.Context,
	cookie *http.Cookie,
	body *v1.OAuth2ClientCreationInput,
) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateOAuth2ClientRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, "oauth2", "client")

	req, err := c.buildDataRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)

	return req, nil
}

// CreateOAuth2Client creates an OAuth2 client. Note that cookie must not be nil
// in order to receive a valid response.
func (c *V1Client) CreateOAuth2Client(
	ctx context.Context,
	cookie *http.Cookie,
	input *v1.OAuth2ClientCreationInput,
) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateOAuth2Client")
	defer span.End()

	var oauth2Client *v1.OAuth2Client
	if cookie == nil {
		return nil, errors.New("cookie required for request")
	}

	req, err := c.BuildCreateOAuth2ClientRequest(ctx, cookie, input)
	if err != nil {
		return nil, err
	}

	if resErr := c.executeUnauthenticatedDataRequest(ctx, req, &oauth2Client); resErr != nil {
		return nil, fmt.Errorf("loading response from server: %w", resErr)
	}

	return oauth2Client, nil
}

// BuildArchiveOAuth2ClientRequest builds an HTTP request for archiving an oauth2 client.
func (c *V1Client) BuildArchiveOAuth2ClientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveOAuth2ClientRequest")
	defer span.End()

	uri := c.BuildURL(nil, oauth2ClientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (c *V1Client) ArchiveOAuth2Client(ctx context.Context, id uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveOAuth2Client")
	defer span.End()

	req, err := c.BuildArchiveOAuth2ClientRequest(ctx, id)
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

func Test_buildBuildGetOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetOAuth2ClientRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildGetOAuth2ClientRequest builds an HTTP request for fetching an OAuth2 client.
func (c *V1Client) BuildGetOAuth2ClientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetOAuth2ClientRequest")
	defer span.End()

	uri := c.BuildURL(nil, oauth2ClientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetOAuth2Client gets an OAuth2 client.
func (c *V1Client) GetOAuth2Client(ctx context.Context, id uint64) (oauth2Client *v1.OAuth2Client, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2Client")
	defer span.End()

	req, err := c.BuildGetOAuth2ClientRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &oauth2Client)
	return oauth2Client, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetOAuth2ClientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetOAuth2ClientsRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
)

// BuildGetOAuth2ClientsRequest builds an HTTP request for fetching a list of OAuth2 clients.
func (c *V1Client) BuildGetOAuth2ClientsRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetOAuth2ClientsRequest")
	defer span.End()

	uri := c.BuildURL(filter.ToValues(), oauth2ClientsBasePath)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Clients(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetOAuth2Clients gets a list of OAuth2 clients.
func (c *V1Client) GetOAuth2Clients(ctx context.Context, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetOAuth2Clients")
	defer span.End()

	req, err := c.BuildGetOAuth2ClientsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	var oauth2Clients *v1.OAuth2ClientList
	err = c.retrieve(ctx, req, &oauth2Clients)
	return oauth2Clients, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateOAuth2ClientRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
)

// BuildCreateOAuth2ClientRequest builds an HTTP request for creating OAuth2 clients.
func (c *V1Client) BuildCreateOAuth2ClientRequest(
	ctx context.Context,
	cookie *http.Cookie,
	body *v1.OAuth2ClientCreationInput,
) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateOAuth2ClientRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, "oauth2", "client")

	req, err := c.buildDataRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)

	return req, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"net/http"
)

// CreateOAuth2Client creates an OAuth2 client. Note that cookie must not be nil
// in order to receive a valid response.
func (c *V1Client) CreateOAuth2Client(
	ctx context.Context,
	cookie *http.Cookie,
	input *v1.OAuth2ClientCreationInput,
) (*v1.OAuth2Client, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateOAuth2Client")
	defer span.End()

	var oauth2Client *v1.OAuth2Client
	if cookie == nil {
		return nil, errors.New("cookie required for request")
	}

	req, err := c.BuildCreateOAuth2ClientRequest(ctx, cookie, input)
	if err != nil {
		return nil, err
	}

	if resErr := c.executeUnauthenticatedDataRequest(ctx, req, &oauth2Client); resErr != nil {
		return nil, fmt.Errorf("loading response from server: %w", resErr)
	}

	return oauth2Client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildArchiveOAuth2ClientRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildArchiveOAuth2ClientRequest builds an HTTP request for archiving an oauth2 client.
func (c *V1Client) BuildArchiveOAuth2ClientRequest(ctx context.Context, id uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveOAuth2ClientRequest")
	defer span.End()

	uri := c.BuildURL(nil, oauth2ClientsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildArchiveOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveOAuth2Client archives an OAuth2 client.
func (c *V1Client) ArchiveOAuth2Client(ctx context.Context, id uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveOAuth2Client")
	defer span.End()

	req, err := c.BuildArchiveOAuth2ClientRequest(ctx, id)
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
