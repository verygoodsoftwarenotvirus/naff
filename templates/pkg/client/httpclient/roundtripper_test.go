package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_roundtripperDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := roundtripperDotGo(proj)

		expected := `
package example

import (
	"net"
	"net/http"
	"time"
)

const (
	userAgentHeader = "User-Agent"
	userAgent       = "Todo Service Client"
)

type defaultRoundTripper struct {
	baseTransport *http.Transport
}

// newDefaultRoundTripper constructs a new http.RoundTripper.
func newDefaultRoundTripper() *defaultRoundTripper {
	return &defaultRoundTripper{
		baseTransport: buildDefaultTransport(),
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (t *defaultRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(userAgentHeader, userAgent)
	return t.baseTransport.RoundTrip(req)
}

// buildDefaultTransport constructs a new http.Transport.
func buildDefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 2 * defaultTimeout,
		IdleConnTimeout:       3 * defaultTimeout,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRoundtripperConstDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildRoundtripperConstDecls(proj)

		expected := `
package example

import ()

const (
	userAgentHeader = "User-Agent"
	userAgent       = "Todo Service Client"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDefaultRoundTripper(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildDefaultRoundTripper()

		expected := `
package example

import (
	"net/http"
)

type defaultRoundTripper struct {
	baseTransport *http.Transport
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildNewDefaultRoundTripper(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildNewDefaultRoundTripper()

		expected := `
package example

import ()

// newDefaultRoundTripper constructs a new http.RoundTripper.
func newDefaultRoundTripper() *defaultRoundTripper {
	return &defaultRoundTripper{
		baseTransport: buildDefaultTransport(),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRoundTrip(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRoundTrip()

		expected := `
package example

import (
	"net/http"
)

// RoundTrip implements the http.RoundTripper interface.
func (t *defaultRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(userAgentHeader, userAgent)
	return t.baseTransport.RoundTrip(req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildDefaultTransport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildDefaultTransport()

		expected := `
package example

import (
	"net"
	"net/http"
	"time"
)

// buildDefaultTransport constructs a new http.Transport.
func buildDefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 2 * defaultTimeout,
		IdleConnTimeout:       3 * defaultTimeout,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
