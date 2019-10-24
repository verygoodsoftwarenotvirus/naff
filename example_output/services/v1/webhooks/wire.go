package webhooks

import (
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	database "gitlab.com/verygoodsoftwarenotvirus/todo/database/v1"
)

var (
	// Providers is our collection of what we provide to other services
	Providers = wire.NewSet(
		ProvideWebhooksService,
		ProvideWebhookDataManager,
		ProvideWebhookDataServer,
	)
)

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake
func ProvideWebhookDataManager(db database.Database) models.WebhookDataManager {
	return db
}

// ProvideWebhookDataServer is an arbitrary function for dependency injection's sake
func ProvideWebhookDataServer(s *Service) models.WebhookDataServer {
	return s
}
