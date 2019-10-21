package oauth2clients

import (
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

var (
	// Providers are what we provide for dependency injection
	Providers = wire.NewSet(
		ProvideOAuth2ClientsService,
		ProvideOAuth2ClientDataServer,
	)
)

// ProvideOAuth2ClientDataServer is an arbitrary function for dependency injection's sake
func ProvideOAuth2ClientDataServer(s *Service) models.OAuth2ClientDataServer {
	return s
}