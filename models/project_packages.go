package models

import "path/filepath"

func (p *Project) RelativePath(parts ...string) string {
	return filepath.Join(append([]string{p.OutputPath}, parts...)...)
}

func (p *Project) HTTPClientV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"pkg", "client", "httpclient"}, parts...)...)
}

func (p *Project) TypesPackage(parts ...string) string {
	return p.RelativePath(append([]string{"pkg", "types"}, parts...)...)
}

func (p *Project) FakeModelsPackage(parts ...string) string {
	return p.TypesPackage(append([]string{"fake"}, parts...)...)
}

func (p *Project) DatabasePackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"database"}, parts...)...)
}

func (p *Project) InternalV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"internal"}, parts...)...)
}

func (p *Project) InternalAuthPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"auth"}, parts...)...)
}

func (p *Project) InternalConfigPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"config"}, parts...)...)
}

func (p *Project) InternalEncodingPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"encoding"}, parts...)...)
}

func (p *Project) InternalMetricsPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"metrics"}, parts...)...)
}

func (p *Project) InternalTracingPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"observability", "tracing"}, parts...)...)
}

func (p *Project) InternalLoggingPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"observability", "logging"}, parts...)...)
}

func (p *Project) InternalSearchPackage(parts ...string) string {
	return p.InternalV1Package(append([]string{"search"}, parts...)...)
}

func (p *Project) ServiceV1Package(parts ...string) string {
	return p.InternalV1Package(append([]string{"services"}, parts...)...)
}

func (p *Project) ServiceAuthPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"auth"}, parts...)...)
}

func (p *Project) ServiceFrontendPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"frontend"}, parts...)...)
}

func (p *Project) ServiceOAuth2ClientsPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"oauth2clients"}, parts...)...)
}

func (p *Project) ServiceUsersPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"users"}, parts...)...)
}

func (p *Project) ServiceWebhooksPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"webhooks"}, parts...)...)
}

func (p *Project) TestUtilPackage(parts ...string) string {
	return p.RelativePath(append([]string{"tests", "utils"}, parts...)...)
}
