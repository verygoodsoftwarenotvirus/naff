package items

import (
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var Providers = wire.NewSet(ProvideItemsService, ProvideItemDataManager, ProvideItemDataServer)

// ProvideItemDataManager turns a database into an ItemDataManager
func ProvideItemDataManager(db database.Database) models.ItemDataManager {
	return db
}

// ProvideItemDataServer is an arbitrary function for dependency injection's sake
func ProvideItemDataServer(s *Service) models.ItemDataServer {
	return s
}
