package models

import "path/filepath"

func (p *Project) RelativePath(parts ...string) string {
	return filepath.Join(append([]string{p.OutputPath}, parts...)...)
}

func (p *Project) HTTPServerPackage() string {
	return p.InternalPackage("server")
}

func (p *Project) HTTPClientPackage(parts ...string) string {
	return p.RelativePath(append([]string{"pkg", "client", "httpclient"}, parts...)...)
}

func (p *Project) TypesPackage(parts ...string) string {
	return p.RelativePath(append([]string{"pkg", "types"}, parts...)...)
}

func (p *Project) FakeTypesPackage(parts ...string) string {
	return p.TypesPackage(append([]string{"fakes"}, parts...)...)
}

func (p *Project) DatabasePackage(parts ...string) string {
	return p.InternalPackage(append([]string{"database"}, parts...)...)
}

func (p *Project) QuerybuildingPackage(parts ...string) string {
	return p.DatabasePackage(append([]string{"querybuilding"}, parts...)...)
}

func (p *Project) InternalPackage(parts ...string) string {
	return p.RelativePath(append([]string{"internal"}, parts...)...)
}

func (p *Project) InternalAuthenticationPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"authentication"}, parts...)...)
}

func (p *Project) InternalAuthorizationPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"authorization"}, parts...)...)
}

func (p *Project) InternalAuditPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"audit"}, parts...)...)
}

func (p *Project) ConfigPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"config"}, parts...)...)
}

func (p *Project) UploadsPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"uploads"}, parts...)...)
}

func (p *Project) StoragePackage() string {
	return p.InternalPackage("storage")
}

func (p *Project) CapitalismPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"capitalism"}, parts...)...)
}

func (p *Project) RoutingPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"routing"}, parts...)...)
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

func (p *Project) InternalRoutingPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"routing"}, parts...)...)
}

func (p *Project) InternalSearchPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"search"}, parts...)...)
}

func (p *Project) InternalSecretsPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"secrets"}, parts...)...)
}

func (p *Project) InternalEventsPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"events"}, parts...)...)
}

func (p *Project) InternalPubSubPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"events"}, parts...)...)
}

func (p *Project) InternalImagesPackage() string {
	return p.UploadsPackage("images")
}

func (p *Project) ServicePackage(service string) string {
	return p.InternalPackage(append([]string{"services"}, service)...)
}

func (p *Project) AuditServicePackage() string {
	return p.ServicePackage("audit")
}

func (p *Project) AuthServicePackage() string {
	return p.ServicePackage("authentication")
}

func (p *Project) AdminServicePackage() string {
	return p.ServicePackage("admin")
}

func (p *Project) APIClientsServicePackage() string {
	return p.ServicePackage("apiclients")
}

func (p *Project) FrontendServicePackage() string {
	return p.ServicePackage("frontend")
}

func (p *Project) UsersServicePackage() string {
	return p.ServicePackage("users")
}

func (p *Project) AccountsServicePackage() string {
	return p.ServicePackage("accounts")
}

func (p *Project) WebhooksServicePackage() string {
	return p.ServicePackage("webhooks")
}

func (p *Project) TestsPackage(parts ...string) string {
	return p.RelativePath(append([]string{"tests"}, parts...)...)
}

func (p *Project) TestUtilsPackage() string {
	return p.TestsPackage("utils")
}

func (p *Project) ObservabilityPackage(parts ...string) string {
	return p.InternalPackage(append([]string{"observability"}, parts...)...)
}

func (p *Project) ConstantKeysPackage() string {
	return p.ObservabilityPackage("keys")
}
