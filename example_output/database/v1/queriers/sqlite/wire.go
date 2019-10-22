package sqlite

import "github.com/google/wire"

var Providers = wire.NewSet(ProvideSqliteDB, ProvideSqlite)
