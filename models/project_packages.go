package models

import "path/filepath"

func (p *Project) RelativePath(parts ...string) string {
	return filepath.Join(append([]string{p.OutputPath}, parts...)...)
}

func (p *Project) HTTPClientV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"client", "v1", "http"}, parts...)...)
}

func (p *Project) ModelsV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"models", "v1"}, parts...)...)
}

func (p *Project) FakeModelsPackage(parts ...string) string {
	return p.ModelsV1Package(append([]string{"fake"}, parts...)...)
}

func (p *Project) DatabaseV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"database", "v1"}, parts...)...)
}

func (p *Project) InternalV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"internal", "v1"}, parts...)...)
}

func (p *Project) InternalAuthV1Package(parts ...string) string {
	return p.InternalV1Package(append([]string{"auth"}, parts...)...)
}

func (p *Project) InternalConfigV1Package(parts ...string) string {
	return p.InternalV1Package(append([]string{"config"}, parts...)...)
}

func (p *Project) InternalEncodingV1Package(parts ...string) string {
	return p.InternalV1Package(append([]string{"encoding"}, parts...)...)
}

func (p *Project) InternalMetricsV1Package(parts ...string) string {
	return p.InternalV1Package(append([]string{"metrics"}, parts...)...)
}

func (p *Project) InternalTracingV1Package(parts ...string) string {
	return p.InternalV1Package(append([]string{"tracing"}, parts...)...)
}

func (p *Project) ServiceV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"services", "v1"}, parts...)...)
}

func (p *Project) ServiceV1AuthPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"auth"}, parts...)...)
}

func (p *Project) ServiceV1FrontendPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"frontend"}, parts...)...)
}

func (p *Project) ServiceV1OAuth2ClientsPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"oauth2clients"}, parts...)...)
}

func (p *Project) ServiceV1UsersPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"users"}, parts...)...)
}

func (p *Project) ServiceV1WebhooksPackage(parts ...string) string {
	return p.ServiceV1Package(append([]string{"webhooks"}, parts...)...)
}

func (p *Project) TestutilV1Package(parts ...string) string {
	return p.RelativePath(append([]string{"tests", "v1", "testutil"}, parts...)...)
}
