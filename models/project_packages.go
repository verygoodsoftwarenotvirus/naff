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
	return p.TypesPackage(append([]string{"fakes"}, parts...)...)
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

func (p *Project) ConfigPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"config"}, parts...)...)
}

func (p *Project) EncodingPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"encoding"}, parts...)...)
}

func (p *Project) MetricsPackage(parts ...string) string {
	return p.ObservabilityPackage(append([]string{"metrics"}, parts...)...)
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

func (p *Project) InternalPubSubPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"events"}, parts...)...)
}

func (p *Project) InternalImagesPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"images"}, parts...)...)
}

func (p *Project) ServicePackage(service string) string {
	return p.InternalPackage(append([]string{"services"}, service)...)
}

func (p *Project) ServiceAuthPackage() string {
	return p.ServicePackage("auth")
}

func (p *Project) ServiceFrontendPackage() string {
	return p.ServicePackage("frontend")
}

func (p *Project) ServiceOAuth2ClientsPackage() string {
	return p.ServicePackage("oauth2clients")
}

func (p *Project) ServiceUsersPackage() string {
	return p.ServicePackage("users")
}

func (p *Project) ServiceWebhooksPackage() string {
	return p.ServicePackage("webhooks")
}

func (p *Project) TestUtilPackage(parts ...string) string {
	return p.RelativePath(append([]string{"tests", "utils"}, parts...)...)
}

func (p *Project) ObservabilityPackage(parts ...string) string {
	return p.RelativePath(append([]string{"internal", "observability"}, parts...)...)
}

func (p *Project) ConstantKeysPackage() string {
	return p.ObservabilityPackage("keys")
}
