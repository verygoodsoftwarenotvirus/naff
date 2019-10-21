package frontend

import "github.com/google/wire"

var Providers = wire.NewSet(ProvideFrontendService)
