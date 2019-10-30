package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	client "gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http"
	models "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
	"gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil"

	"github.com/icrowley/fake"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
)

const (
	debug         = true
	nonexistentID = 999999999
)

var (
	urlToUse   string
	todoClient *client.V1Client
)

func init() {
	urlToUse = testutil.DetermineServiceURL()
	logger := zerolog.NewZeroLogger()

	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)

	fake.Seed(time.Now().UnixNano())

	ogUser, err := testutil.CreateObligatoryUser(urlToUse, debug)
	if err != nil {
		logger.Fatal(err)
	}

	oa2Client, err := testutil.CreateObligatoryClient(urlToUse, ogUser)
	if err != nil {
		logger.Fatal(err)
	}

	todoClient = initializeClient(oa2Client)

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
