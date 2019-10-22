package dbclient

import "github.com/google/wire"

var Providers = wire.NewSet(ProvideDatabaseClient)
