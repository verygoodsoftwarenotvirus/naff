package webhooks

import (
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

var Providers = wire.NewSet(ProvideWebhooksService, ProvideWebhookDataManager, ProvideWebhookDataServer)

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake
func ProvideWebhookDataManager(db database.Database) models.WebhookDataManager {
	return db
}

// ProvideWebhookDataServer is an arbitrary function for dependency injection's sake
func ProvideWebhookDataServer(s *Service) models.WebhookDataServer {
	return s
}
