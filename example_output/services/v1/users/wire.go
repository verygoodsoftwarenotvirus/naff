package users

import (
	"github.com/google/wire"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

var Providers = wire.NewSet(ProvideUsersService, ProvideUserDataServer, ProvideUserDataManager)

// ProvideUserDataManager is an arbitrary function for dependency injection's sake
func ProvideUserDataManager(db database.Database) models.UserDataManager {
	return db
}

// ProvideUserDataServer is an arbitrary function for dependency injection's sake
func ProvideUserDataServer(s *Service) models.UserDataServer {
	return s
}
