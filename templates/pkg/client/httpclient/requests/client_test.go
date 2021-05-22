package requests

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_newClientMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := newClientMethod("Example").Params().Body()

		expected := `
package example

import ()

func (c *V1Client) Example() {}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	http2curl "github.com/moul/http2curl"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	ochttp "go.opencensus.io/plugin/ochttp"
	oauth2 "golang.org/x/oauth2"
	clientcredentials "golang.org/x/oauth2/clientcredentials"
	logger "logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	clientName     = "v1_client"
)

var (
	// ErrNotFound is a handy error to return when we receive a 404 response.
	ErrNotFound = fmt.Errorf("%d: not found", http.StatusNotFound)

	// ErrUnauthorized is a handy error to return when we receive a 401 response.
	ErrUnauthorized = fmt.Errorf("%d: not authorized", http.StatusUnauthorized)

	// ErrInvalidTOTPToken is an error for when our TOTP validation request goes awry.
	ErrInvalidTOTPToken = errors.New("invalid TOTP token")
)

// V1Client is a client for interacting with v1 of our HTTP API.
type V1Client struct {
	plainClient  *http.Client
	authedClient *http.Client
	logger       v1.Logger
	Debug        bool
	URL          *url.URL
	Scopes       []string
	tokenSource  oauth2.TokenSource
}

// AuthenticatedClient returns the authenticated *http.Client that we use to make most requests.
func (c *V1Client) AuthenticatedClient() *http.Client {
	return c.authedClient
}

// PlainClient returns the unauthenticated *http.Client that we use to make certain requests.
func (c *V1Client) PlainClient() *http.Client {
	return c.plainClient
}

// TokenSource provides the client's token source.
func (c *V1Client) TokenSource() oauth2.TokenSource {
	return c.tokenSource
}

// tokenEndpoint provides the oauth2 Endpoint for a given host.
func tokenEndpoint(baseURL *url.URL) oauth2.Endpoint {
	tu, au := *baseURL, *baseURL
	tu.Path, au.Path = "oauth2/token", "oauth2/authorize"

	return oauth2.Endpoint{
		TokenURL: tu.String(),
		AuthURL:  au.String(),
	}
}

// NewClient builds a new API client for us.
func NewClient(
	ctx context.Context,
	clientID,
	clientSecret string,
	address *url.URL,
	logger v1.Logger,
	hclient *http.Client,
	scopes []string,
	debug bool,
) (*V1Client, error) {
	var client = hclient
	if client == nil {
		client = &http.Client{
			Timeout: defaultTimeout,
		}
	}
	if client.Timeout == 0 {
		client.Timeout = defaultTimeout
	}

	if debug {
		logger.SetLevel(v1.DebugLevel)
		logger.Debug("log level set to debug!")
	}

	ac, ts := buildOAuthClient(ctx, address, clientID, clientSecret, scopes, client.Timeout)

	c := &V1Client{
		URL:          address,
		plainClient:  client,
		logger:       logger.WithName(clientName),
		Debug:        debug,
		authedClient: ac,
		tokenSource:  ts,
	}

	logger.WithValue("url", address.String()).Debug("returning client")
	return c, nil
}

// NewSimpleClient is a client that is capable of much less than the normal client
// and has noops or empty values for most of its authentication and debug parts.
// Its purpose at the time of this writing is merely so I can make users (which
// is a route that doesn't require authentication.)
func NewSimpleClient(ctx context.Context, address *url.URL, debug bool) (*V1Client, error) {
	return NewClient(
		ctx,
		"",
		"",
		address,
		noop.ProvideNoopLogger(),
		&http.Client{Timeout: 5 * time.Second},
		[]string{"*"},
		debug,
	)
}

// buildOAuthClient takes care of all the OAuth2 noise and returns a nice pretty *http.Client for us to use.
func buildOAuthClient(
	ctx context.Context,
	uri *url.URL,
	clientID,
	clientSecret string,
	scopes []string,
	timeout time.Duration,
) (*http.Client, oauth2.TokenSource) {
	conf := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		EndpointParams: url.Values{
			"client_id":     []string{clientID},
			"client_secret": []string{clientSecret},
		},
		TokenURL: tokenEndpoint(uri).TokenURL,
	}

	ts := oauth2.ReuseTokenSource(nil, conf.TokenSource(ctx))
	client := &http.Client{
		Transport: &oauth2.Transport{
			Base: &ochttp.Transport{
				Base: newDefaultRoundTripper(),
			},
			Source: ts,
		},
		Timeout: timeout,
	}

	return client, ts
}

// closeResponseBody takes a given HTTP response and closes its body, logging if an error occurs.
func (c *V1Client) closeResponseBody(res *http.Response) {
	if res != nil {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(err, "closing response body")
		}
	}
}

// BuildURL builds standard service URLs.
func (c *V1Client) BuildURL(qp url.Values, parts ...string) string {
	var u *url.URL
	if qp != nil {
		u = c.buildURL(qp, parts...)
	} else {
		u = c.buildURL(nil, parts...)
	}

	if u != nil {
		return u.String()
	}
	return ""
}

// buildURL takes a given set of query parameters and URL parts, and returns.
// a parsed URL object from them.
func (c *V1Client) buildURL(queryParams url.Values, parts ...string) *url.URL {
	tu := *c.URL

	parts = append([]string{"api", "v1"}, parts...)
	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		c.logger.Error(err, "building URL")
		return nil
	}

	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	return tu.ResolveReference(u)
}

// buildVersionlessURL builds a URL without the ` + "`" + `/api/v1/` + "`" + ` prefix. It should
// otherwise be identical to buildURL.
func (c *V1Client) buildVersionlessURL(qp url.Values, parts ...string) string {
	tu := *c.URL

	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		c.logger.Error(err, "building URL")
		return ""
	}

	if qp != nil {
		u.RawQuery = qp.Encode()
	}

	return tu.ResolveReference(u).String()
}

// BuildWebsocketURL builds a standard URL and then converts its scheme to the websocket protocol.
func (c *V1Client) BuildWebsocketURL(parts ...string) string {
	u := c.buildURL(nil, parts...)
	u.Scheme = "ws"

	return u.String()
}

// BuildHealthCheckRequest builds a health check HTTP request.
func (c *V1Client) BuildHealthCheckRequest(ctx context.Context) (*http.Request, error) {
	u := *c.URL
	uri := fmt.Sprintf("%s://%s/_meta_/ready", u.Scheme, u.Host)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// IsUp returns whether or not the service's health endpoint is returning 200s.
func (c *V1Client) IsUp(ctx context.Context) bool {
	req, err := c.BuildHealthCheckRequest(ctx)
	if err != nil {
		c.logger.Error(err, "building request")
		return false
	}

	res, err := c.plainClient.Do(req)
	if err != nil {
		c.logger.Error(err, "health check")
		return false
	}
	c.closeResponseBody(res)

	return res.StatusCode == http.StatusOK
}

// buildDataRequest builds an HTTP request for a given method, URL, and body data.
func (c *V1Client) buildDataRequest(ctx context.Context, method, uri string, in interface{}) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "buildDataRequest")
	defer span.End()

	body, err := createBodyFromStruct(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	return req, nil
}

// executeRequest takes a given request and executes it with the auth client. It returns some errors
// upon receiving certain status codes, but otherwise will return nil upon success.
func (c *V1Client) executeRequest(ctx context.Context, req *http.Request, out interface{}) error {
	ctx, span := tracing.StartSpan(ctx, "executeRequest")
	defer span.End()

	res, err := c.executeRawRequest(ctx, c.authedClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	}

	if out != nil {
		if resErr := unmarshalBody(ctx, res, out); resErr != nil {
			return fmt.Errorf("loading response from server: %w", err)
		}
	}

	return nil
}

// executeRawRequest takes a given *http.Request and executes it with the provided.
// client, alongside some debugging logging.
func (c *V1Client) executeRawRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	ctx, span := tracing.StartSpan(ctx, "executeRawRequest")
	defer span.End()

	var logger = c.logger
	if command, err := http2curl.GetCurlCommand(req); err == nil && c.Debug {
		logger = c.logger.WithValue("curl", command.String())
	}

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	if c.Debug {
		bdump, err := httputil.DumpResponse(res, true)
		if err == nil && req.Method != http.MethodGet {
			logger = logger.WithValue("response_body", string(bdump))
		}
		logger.Debug("request executed")
	}

	return res, nil
}

// checkExistence executes an HTTP request and loads the response content into a bool.
func (c *V1Client) checkExistence(ctx context.Context, req *http.Request) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "checkExistence")
	defer span.End()

	res, err := c.executeRawRequest(ctx, c.authedClient, req)
	if err != nil {
		return false, err
	}
	c.closeResponseBody(res)

	return res.StatusCode == http.StatusOK, nil
}

// retrieve executes an HTTP request and loads the response content into a struct. In the event of a 404,
// the provided ErrNotFound is returned.
func (c *V1Client) retrieve(ctx context.Context, req *http.Request, obj interface{}) error {
	ctx, span := tracing.StartSpan(ctx, "retrieve")
	defer span.End()

	if err := argIsNotPointerOrNil(obj); err != nil {
		return fmt.Errorf("struct to load must be a pointer: %w", err)
	}

	res, err := c.executeRawRequest(ctx, c.authedClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	return unmarshalBody(ctx, res, &obj)
}

// executeUnauthenticatedDataRequest takes a given request and loads the response into an interface value.
func (c *V1Client) executeUnauthenticatedDataRequest(ctx context.Context, req *http.Request, out interface{}) error {
	ctx, span := tracing.StartSpan(ctx, "executeUnauthenticatedDataRequest")
	defer span.End()

	res, err := c.executeRawRequest(ctx, c.plainClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	}

	if out != nil {
		if resErr := unmarshalBody(ctx, res, out); resErr != nil {
			return fmt.Errorf("loading response from server: %w", err)
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildClientConstDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildClientConstDecls()

		expected := `
package example

import (
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	clientName     = "v1_client"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildClientVarDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildClientVarDecls()

		expected := `
package example

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNotFound is a handy error to return when we receive a 404 response.
	ErrNotFound = fmt.Errorf("%d: not found", http.StatusNotFound)

	// ErrUnauthorized is a handy error to return when we receive a 401 response.
	ErrUnauthorized = fmt.Errorf("%d: not authorized", http.StatusUnauthorized)

	// ErrInvalidTOTPToken is an error for when our TOTP validation request goes awry.
	ErrInvalidTOTPToken = errors.New("invalid TOTP token")
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildClientTypeDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildClientTypeDecls()

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	oauth2 "golang.org/x/oauth2"
	"net/http"
	"net/url"
)

// V1Client is a client for interacting with v1 of our HTTP API.
type V1Client struct {
	plainClient  *http.Client
	authedClient *http.Client
	logger       v1.Logger
	Debug        bool
	URL          *url.URL
	Scopes       []string
	tokenSource  oauth2.TokenSource
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthenticatedClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildAuthenticatedClient()

		expected := `
package example

import (
	"net/http"
)

// AuthenticatedClient returns the authenticated *http.Client that we use to make most requests.
func (c *V1Client) AuthenticatedClient() *http.Client {
	return c.authedClient
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildPlainClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildPlainClient()

		expected := `
package example

import (
	"net/http"
)

// PlainClient returns the unauthenticated *http.Client that we use to make certain requests.
func (c *V1Client) PlainClient() *http.Client {
	return c.plainClient
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTokenSource(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTokenSource()

		expected := `
package example

import ()

// TokenSource provides the client's token source.
func (c *V1Client) TokenSource() oauth2.TokenSource {
	return c.tokenSource
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildNewClient()

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	logger "logger"
	"net/http"
	"net/url"
)

// NewClient builds a new API client for us.
func NewClient(
	ctx context.Context,
	clientID,
	clientSecret string,
	address *url.URL,
	logger v1.Logger,
	hclient *http.Client,
	scopes []string,
	debug bool,
) (*V1Client, error) {
	var client = hclient
	if client == nil {
		client = &http.Client{
			Timeout: defaultTimeout,
		}
	}
	if client.Timeout == 0 {
		client.Timeout = defaultTimeout
	}

	if debug {
		logger.SetLevel(v1.DebugLevel)
		logger.Debug("log level set to debug!")
	}

	ac, ts := buildOAuthClient(ctx, address, clientID, clientSecret, scopes, client.Timeout)

	c := &V1Client{
		URL:          address,
		plainClient:  client,
		logger:       logger.WithName(clientName),
		Debug:        debug,
		authedClient: ac,
		tokenSource:  ts,
	}

	logger.WithValue("url", address.String()).Debug("returning client")
	return c, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildOAuthClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildOAuthClient()

		expected := `
package example

import (
	"context"
	ochttp "go.opencensus.io/plugin/ochttp"
	clientcredentials "golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"time"
)

// buildOAuthClient takes care of all the OAuth2 noise and returns a nice pretty *http.Client for us to use.
func buildOAuthClient(
	ctx context.Context,
	uri *url.URL,
	clientID,
	clientSecret string,
	scopes []string,
	timeout time.Duration,
) (*http.Client, oauth2.TokenSource) {
	conf := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		EndpointParams: url.Values{
			"client_id":     []string{clientID},
			"client_secret": []string{clientSecret},
		},
		TokenURL: tokenEndpoint(uri).TokenURL,
	}

	ts := oauth2.ReuseTokenSource(nil, conf.TokenSource(ctx))
	client := &http.Client{
		Transport: &oauth2.Transport{
			Base: &ochttp.Transport{
				Base: newDefaultRoundTripper(),
			},
			Source: ts,
		},
		Timeout: timeout,
	}

	return client, ts
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTokenEndpoint(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTokenEndpoint()

		expected := `
package example

import (
	"net/url"
)

// tokenEndpoint provides the oauth2 Endpoint for a given host.
func tokenEndpoint(baseURL *url.URL) oauth2.Endpoint {
	tu, au := *baseURL, *baseURL
	tu.Path, au.Path = "oauth2/token", "oauth2/authorize"

	return oauth2.Endpoint{
		TokenURL: tu.String(),
		AuthURL:  au.String(),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewSimpleClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildNewSimpleClient()

		expected := `
package example

import (
	"context"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"net/http"
	"net/url"
	"time"
)

// NewSimpleClient is a client that is capable of much less than the normal client
// and has noops or empty values for most of its authentication and debug parts.
// Its purpose at the time of this writing is merely so I can make users (which
// is a route that doesn't require authentication.)
func NewSimpleClient(ctx context.Context, address *url.URL, debug bool) (*V1Client, error) {
	return NewClient(
		ctx,
		"",
		"",
		address,
		noop.ProvideNoopLogger(),
		&http.Client{Timeout: 5 * time.Second},
		[]string{"*"},
		debug,
	)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCloseResponseBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildCloseResponseBody()

		expected := `
package example

import (
	"net/http"
)

// closeResponseBody takes a given HTTP response and closes its body, logging if an error occurs.
func (c *V1Client) closeResponseBody(res *http.Response) {
	if res != nil {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(err, "closing response body")
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExecuteRawRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildExecuteRawRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	http2curl "github.com/moul/http2curl"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"net/http/httputil"
)

// executeRawRequest takes a given *http.Request and executes it with the provided.
// client, alongside some debugging logging.
func (c *V1Client) executeRawRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	ctx, span := tracing.StartSpan(ctx, "executeRawRequest")
	defer span.End()

	var logger = c.logger
	if command, err := http2curl.GetCurlCommand(req); err == nil && c.Debug {
		logger = c.logger.WithValue("curl", command.String())
	}

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	if c.Debug {
		bdump, err := httputil.DumpResponse(res, true)
		if err == nil && req.Method != http.MethodGet {
			logger = logger.WithValue("response_body", string(bdump))
		}
		logger.Debug("request executed")
	}

	return res, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExportedBuildURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildExportedBuildURL()

		expected := `
package example

import (
	"net/url"
)

// BuildURL builds standard service URLs.
func (c *V1Client) BuildURL(qp url.Values, parts ...string) string {
	var u *url.URL
	if qp != nil {
		u = c.buildURL(qp, parts...)
	} else {
		u = c.buildURL(nil, parts...)
	}

	if u != nil {
		return u.String()
	}
	return ""
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUnexportedBuildURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildUnexportedBuildURL()

		expected := `
package example

import (
	"net/url"
	"strings"
)

// buildURL takes a given set of query parameters and URL parts, and returns.
// a parsed URL object from them.
func (c *V1Client) buildURL(queryParams url.Values, parts ...string) *url.URL {
	tu := *c.URL

	parts = append([]string{"api", "v1"}, parts...)
	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		c.logger.Error(err, "building URL")
		return nil
	}

	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	return tu.ResolveReference(u)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildVersionlessURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildVersionlessURL()

		expected := `
package example

import (
	"net/url"
	"path"
)

// buildVersionlessURL builds a URL without the ` + "`" + `/api/v1/` + "`" + ` prefix. It should
// otherwise be identical to buildURL.
func (c *V1Client) buildVersionlessURL(qp url.Values, parts ...string) string {
	tu := *c.URL

	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		c.logger.Error(err, "building URL")
		return ""
	}

	if qp != nil {
		u.RawQuery = qp.Encode()
	}

	return tu.ResolveReference(u).String()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildWebsocketURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildWebsocketURL()

		expected := `
package example

import ()

// BuildWebsocketURL builds a standard URL and then converts its scheme to the websocket protocol.
func (c *V1Client) BuildWebsocketURL(parts ...string) string {
	u := c.buildURL(nil, parts...)
	u.Scheme = "ws"

	return u.String()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildHealthCheckRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildHealthCheckRequest()

		expected := `
package example

import (
	"context"
	"fmt"
	"net/http"
)

// BuildHealthCheckRequest builds a health check HTTP request.
func (c *V1Client) BuildHealthCheckRequest(ctx context.Context) (*http.Request, error) {
	u := *c.URL
	uri := fmt.Sprintf("%s://%s/_meta_/ready", u.Scheme, u.Host)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIsUp(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildIsUp()

		expected := `
package example

import (
	"context"
	"net/http"
)

// IsUp returns whether or not the service's health endpoint is returning 200s.
func (c *V1Client) IsUp(ctx context.Context) bool {
	req, err := c.BuildHealthCheckRequest(ctx)
	if err != nil {
		c.logger.Error(err, "building request")
		return false
	}

	res, err := c.plainClient.Do(req)
	if err != nil {
		c.logger.Error(err, "health check")
		return false
	}
	c.closeResponseBody(res)

	return res.StatusCode == http.StatusOK
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildDataRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildDataRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// buildDataRequest builds an HTTP request for a given method, URL, and body data.
func (c *V1Client) buildDataRequest(ctx context.Context, method, uri string, in interface{}) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "buildDataRequest")
	defer span.End()

	body, err := createBodyFromStruct(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	return req, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCheckExistence(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCheckExistence(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// checkExistence executes an HTTP request and loads the response content into a bool.
func (c *V1Client) checkExistence(ctx context.Context, req *http.Request) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "checkExistence")
	defer span.End()

	res, err := c.executeRawRequest(ctx, c.authedClient, req)
	if err != nil {
		return false, err
	}
	c.closeResponseBody(res)

	return res.StatusCode == http.StatusOK, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRetrieve(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildRetrieve(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// retrieve executes an HTTP request and loads the response content into a struct. In the event of a 404,
// the provided ErrNotFound is returned.
func (c *V1Client) retrieve(ctx context.Context, req *http.Request, obj interface{}) error {
	ctx, span := tracing.StartSpan(ctx, "retrieve")
	defer span.End()

	if err := argIsNotPointerOrNil(obj); err != nil {
		return fmt.Errorf("struct to load must be a pointer: %w", err)
	}

	res, err := c.executeRawRequest(ctx, c.authedClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	return unmarshalBody(ctx, res, &obj)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExecuteRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildExecuteRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// executeRequest takes a given request and executes it with the auth client. It returns some errors
// upon receiving certain status codes, but otherwise will return nil upon success.
func (c *V1Client) executeRequest(ctx context.Context, req *http.Request, out interface{}) error {
	ctx, span := tracing.StartSpan(ctx, "executeRequest")
	defer span.End()

	res, err := c.executeRawRequest(ctx, c.authedClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	}

	if out != nil {
		if resErr := unmarshalBody(ctx, res, out); resErr != nil {
			return fmt.Errorf("loading response from server: %w", err)
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildExecuteUnauthenticatedDataRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildExecuteUnauthenticatedDataRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// executeUnauthenticatedDataRequest takes a given request and loads the response into an interface value.
func (c *V1Client) executeUnauthenticatedDataRequest(ctx context.Context, req *http.Request, out interface{}) error {
	ctx, span := tracing.StartSpan(ctx, "executeUnauthenticatedDataRequest")
	defer span.End()

	res, err := c.executeRawRequest(ctx, c.plainClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	}

	if out != nil {
		if resErr := unmarshalBody(ctx, res, out); resErr != nil {
			return fmt.Errorf("loading response from server: %w", err)
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
