package models

import "path/filepath"

func (p *Project) RelativePath(parts ...string) string {
	return filepath.Join(append([]string{p.OutputPath}, parts...)...)
}

func (p *Project) HTTPClientPackage(parts ...string) string {
	return p.RelativePath(append([]string{"pkg", "client", "httpclient"}, parts...)...)
}

func (p *Project) TypesPackage(parts ...string) string {
	return p.RelativePath(append([]string{"pkg", "types"}, parts...)...)
}

func (p *Project) FakeModelsPackage(parts ...string) string {
	return p.TypesPackage(append([]string{"fake"}, parts...)...)
}

func (p *Project) DatabasePackage(parts ...string) string {
	return p.InternalPackage(append([]string{"database"}, parts...)...)
}

func (p *Project) QuerybuildersPackage(parts ...string) string {
	return p.DatabasePackage(append([]string{"querybuilders"}, parts...)...)
}

func (p *Project) InternalPackage(parts ...string) string {
	return p.RelativePath(append([]string{"internal"}, parts...)...)
}

func (p *Project) InternalAuthPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"authentication"}, parts...)...)
}

func (p *Project) InternalAuditPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"audit"}, parts...)...)
}

func (p *Project) InternalConfigPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"config"}, parts...)...)
}

func (p *Project) InternalEncodingPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"encoding"}, parts...)...)
}

func (p *Project) InternalMetricsPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"metrics"}, parts...)...)
}

func (p *Project) InternalTracingPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"observability", "tracing"}, parts...)...)
}

func (p *Project) InternalLoggingPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"observability", "logging"}, parts...)...)
}

func (p *Project) InternalSearchPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"search"}, parts...)...)
}

func (p *Project) InternalSecretsPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"secrets"}, parts...)...)
}

func (p *Project) InternalEventzPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"events"}, parts...)...)
}

func (p *Project) ServicePackage(parts ...string) string {
	return p.InternalPackage(append([]string{"services"}, parts...)...)
}

func (p *Project) ServiceAuthPackage(parts ...string) string {
	return p.ServicePackage(append([]string{"auth"}, parts...)...)
}

func (p *Project) ServiceFrontendPackage(parts ...string) string {
	return p.ServicePackage(append([]string{"frontend"}, parts...)...)
}

func (p *Project) ServiceOAuth2ClientsPackage(parts ...string) string {
	return p.ServicePackage(append([]string{"oauth2clients"}, parts...)...)
}

func (p *Project) ServiceUsersPackage(parts ...string) string {
	return p.ServicePackage(append([]string{"users"}, parts...)...)
}

func (p *Project) ServiceWebhooksPackage(parts ...string) string {
	return p.ServicePackage(append([]string{"webhooks"}, parts...)...)
}

func (p *Project) TestUtilPackage(parts ...string) string {
	return p.RelativePath(append([]string{"tests", "utils"}, parts...)...)
}
