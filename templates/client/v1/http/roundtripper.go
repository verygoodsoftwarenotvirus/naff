package client

import (
	"net"
	"net/http"
	"time"
)

const (
	userAgentHeader = "User-Agent"
	userAgent       = "TODO Service Client"
)

type defaultRoundTripper struct {
	baseTransport *http.Transport
}

//
func newDefaultRoundTripper() *defaultRoundTripper {
	return &defaultRoundTripper{baseTransport: buildDefaultTransport()}
}

//
func (t *defaultRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(userAgentHeader, userAgent)

	return t.baseTransport.RoundTrip(req)
}

//
func buildDefaultTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			DualStack: true,
			KeepAlive: 30 * time.Second,
			Timeout:   30 * time.Second,
		}).DialContext,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   10 * time.Second,
	}
}
