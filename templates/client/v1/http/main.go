package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

var (
	defaultTimeout  = 5 * time.Second
	ErrNotFound     = errors.New("404: not found")
	ErrUnauthorized = errors.New("401: not authorized")
)

type V1Client struct {
	plainClient  *http.Client
	authedClient *http.Client
	logger       logging.Logger
	Debug        bool
	URL          *url.URL
	Scopes       []string
	tokenSource  oauth2.TokenSource
}

func NewClient(ctx context.Context, clientID, clientSecret string, address *url.URL, logger logging.Logger, hclient *http.Client, scopes []string, debug bool) (*V1Client, error) {
	var client = hclient
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	if client.Timeout == 0 {
		client.Timeout = defaultTimeout
	}
	if debug {
		logger.SetLevel(logging.DebugLevel)
		logger.Debug("log level set to debug!")
	}
	ac, ts := buildOAuthClient(ctx, address, clientID, clientSecret, scopes)
	c := &V1Client{URL: address, plainClient: client, logger: logger.WithName("v1_client"), Debug: debug, authedClient: ac, tokenSource: ts}
	logger.WithValue("url", address.String()).Debug("returning client")
	return c, nil
}

func buildOAuthClient(ctx context.Context, uri *url.URL, clientID, clientSecret string, scopes []string) (*http.Client, oauth2.TokenSource) {
	conf := clientcredentials.Config{ClientID: clientID, ClientSecret: clientSecret, Scopes: scopes, EndpointParams: url.Values{"client_id": []string{clientID}, "client_secret": []string{clientSecret}}, TokenURL: tokenEndpoint(uri).TokenURL}
	ts := oauth2.ReuseTokenSource(nil, conf.TokenSource(ctx))
	client := &http.Client{Transport: &oauth2.Transport{Base: &ochttp.Transport{Base: newDefaultRoundTripper()}, Source: ts}, Timeout: 5 * time.Second}
	return client, ts
}

func tokenEndpoint(baseURL *url.URL) oauth2.Endpoint {
	tu, au := *baseURL, *baseURL
	tu.Path, au.Path = "oauth2/token", "oauth2/authorize"
	return oauth2.Endpoint{TokenURL: tu.String(), AuthURL: au.String()}
}

func NewSimpleClient(ctx context.Context, address *url.URL, debug bool) (*V1Client, error) {
	l := noop.ProvideNoopLogger()
	h := &http.Client{Timeout: 5 * time.Second}
	c, err := NewClient(ctx, "", "", address, l, h, []string{"*"}, debug)
	return c, err
}

func (c *V1Client) AuthenticatedClient() *http.Client {
	return c.authedClient
}

func (c *V1Client) PlainClient() *http.Client {
	return c.plainClient
}

func (c *V1Client) TokenSource() oauth2.TokenSource {
	return c.tokenSource
}

func (c *V1Client) executeRawRequest(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
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

func (c *V1Client) BuildURL(qp url.Values, parts ...string) string {
	if qp != nil {
		return c.buildURL(qp, parts...).String()
	}
	return c.buildURL(nil, parts...).String()
}

func (c *V1Client) buildURL(queryParams url.Values, parts ...string) *url.URL {
	tu := *c.URL
	parts = append([]string{"api", "v1"}, parts...)
	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		panic(fmt.Sprintf("was asked to build an invalid URL: %v", err))
	}
	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}
	return tu.ResolveReference(u)
}

func (c *V1Client) buildVersionlessURL(qp url.Values, parts ...string) string {
	tu := *c.URL
	u, err := url.Parse(path.Join(parts...))
	if err != nil {
		panic(fmt.Sprintf("user tried to build an invalid URL: %v", err))
	}
	if qp != nil {
		u.RawQuery = qp.Encode()
	}
	return tu.ResolveReference(u).String()
}

func (c *V1Client) BuildWebsocketURL(parts ...string) string {
	u := c.buildURL(nil, parts...)
	u.Scheme = "ws"
	return u.String()
}

func (c *V1Client) BuildHealthCheckRequest() (*http.Request, error) {
	u := *c.URL
	uri := fmt.Sprintf("%s://%s/_meta_/ready", u.Scheme, u.Host)
	return http.NewRequest(http.MethodGet, uri, nil)
}

func (c *V1Client) IsUp() bool {
	req, err := c.BuildHealthCheckRequest()
	if err != nil {
		c.logger.Error(err, "building request")
		return false
	}
	res, err := c.plainClient.Do(req)
	if err != nil {
		c.logger.Error(err, "health check")
		return false
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(err, "closing response body")
		}
	}()
	return res.StatusCode == http.StatusOK
}

func (c *V1Client) buildDataRequest(method, uri string, in interface{}) (*http.Request, error) {
	body, err := createBodyFromStruct(in)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	return req, nil
}

func (c *V1Client) retrieve(ctx context.Context, req *http.Request, obj interface{}) error {
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
	return unmarshalBody(res, &obj)
}

func (c *V1Client) executeRequest(ctx context.Context, req *http.Request, out interface{}) error {
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
		resErr := unmarshalBody(res, &out)
		if resErr != nil {
			return fmt.Errorf("loading response from server: %w", err)
		}
	}
	return nil
}

func (c *V1Client) executeUnathenticatedDataRequest(ctx context.Context, req *http.Request, out interface{}) error {
	if out != nil {
		if np, err := argIsNotPointer(out); np || err != nil {
			return fmt.Errorf("struct to load must be a pointer: %w", err)
		}
	}
	res, err := c.executeRawRequest(ctx, c.plainClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	if res.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if out != nil {
		resErr := unmarshalBody(res, &out)
		if resErr != nil {
			return fmt.Errorf("loading response from server: %w", err)
		}
	}
	return nil
}
