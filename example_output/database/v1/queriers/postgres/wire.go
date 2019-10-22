package postgres

import "github.com/google/wire"

var Providers = wire.NewSet(ProvidePostgresDB, ProvidePostgres)
