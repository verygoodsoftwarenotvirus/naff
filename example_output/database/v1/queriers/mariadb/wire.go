package mariadb

import "github.com/google/wire"

var Providers = wire.NewSet(ProvideMariaDBConnection, ProvideMariaDB)
