package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	client "gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http"

	"github.com/icrowley/fake"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
)

var (
	debug     bool
	urlToUse  string
	oa2Client *models.OAuth2Client
)

func init() {
	urlToUse = testutil.DetermineServiceURL()
	logger := zerolog.NewZeroLogger()
	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)
	fake.Seed(time.Now().UnixNano())
	u, err := testutil.CreateObligatoryUser(urlToUse, debug)
	if err != nil {
		logger.Fatal(err)
	}
	oa2Client, err = testutil.CreateObligatoryClient(urlToUse, u)
	if err != nil {
		logger.Fatal(err)
	}
	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}

func buildHTTPClient() *http.Client {
	httpc := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}
	return httpc
}

func initializeClient(oa2Client *models.OAuth2Client) *client.V1Client {
	uri, err := url.Parse(urlToUse)
	if err != nil {
		panic(err)
	}
	c, err := client.NewClient(
		context.Background(),
		oa2Client.ClientID,
		oa2Client.ClientSecret,
		uri,
		zerolog.NewZeroLogger(),
		buildHTTPClient(),
		oa2Client.Scopes,
		debug,
	)
	if err != nil {
		panic(err)
	}
	return c
}
