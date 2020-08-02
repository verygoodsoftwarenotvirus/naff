package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_initDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := initDotGo(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	http1 "net/http"
	"net/url"
	"strings"
	"time"
)

const (
	debug         = true
	nonexistentID = 999999999
)

var (
	urlToUse   string
	todoClient *http.V1Client
)

func init() {
	urlToUse = testutil.DetermineServiceURL()
	logger := zerolog.NewZeroLogger()

	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)

	ogUser, err := testutil.CreateObligatoryUser(urlToUse, debug)
	if err != nil {
		logger.Fatal(err)
	}

	oa2Client, err := testutil.CreateObligatoryClient(urlToUse, ogUser)
	if err != nil {
		logger.Fatal(err)
	}

	todoClient = initializeClient(oa2Client)
	todoClient.Debug = urlToUse == "" // change this for debug logs

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}

func buildHTTPClient() *http1.Client {
	return &http1.Client{
		Transport: http1.DefaultTransport,
		Timeout:   5 * time.Second,
	}
}

func initializeClient(oa2Client *v1.OAuth2Client) *http.V1Client {
	uri, err := url.Parse(urlToUse)
	if err != nil {
		panic(err)
	}

	c, err := http.NewClient(
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
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitConstDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInitConstDefs()

		expected := `
package example

import ()

const (
	debug         = true
	nonexistentID = 999999999
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildInitVarDefs(proj)

		expected := `
package example

import (
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
)

var (
	urlToUse   string
	todoClient *http.V1Client
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildInitInit(proj)

		expected := `
package example

import (
	"fmt"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"strings"
)

func init() {
	urlToUse = testutil.DetermineServiceURL()
	logger := zerolog.NewZeroLogger()

	logger.WithValue("url", urlToUse).Info("checking server")
	testutil.EnsureServerIsUp(urlToUse)

	ogUser, err := testutil.CreateObligatoryUser(urlToUse, debug)
	if err != nil {
		logger.Fatal(err)
	}

	oa2Client, err := testutil.CreateObligatoryClient(urlToUse, ogUser)
	if err != nil {
		logger.Fatal(err)
	}

	todoClient = initializeClient(oa2Client)
	todoClient.Debug = urlToUse == "" // change this for debug logs

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitBuildHTTPClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInitBuildHTTPClient()

		expected := `
package example

import (
	"net/http"
	"time"
)

func buildHTTPClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInitInitializeClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildInitInitializeClient(proj)

		expected := `
package example

import (
	"context"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/url"
)

func initializeClient(oa2Client *v1.OAuth2Client) *http.V1Client {
	uri, err := url.Parse(urlToUse)
	if err != nil {
		panic(err)
	}

	c, err := http.NewClient(
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
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
