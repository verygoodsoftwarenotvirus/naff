package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	newsmanImp              = "gitlab.com/verygoodsoftwarenotvirus/newsman"
	loggingImp              = "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	internalConfigImp       = "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config"
	internalMetricsImp      = "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics"
	internalEncodingImp     = "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding"
	databaseClientImp       = "gitlab.com/verygoodsoftwarenotvirus/todo/database/v1"
	internalAuthImp         = "gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth"
	authServiceImp          = "gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth"
	usersServiceImp         = "gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users"
	itemsServiceImp         = "gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items"
	frontendServiceImp      = "gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend"
	webhooksServiceImp      = "gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks"
	oauth2ClientsServiceImp = "gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients"
	httpServerImp           = "gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http"
	serverImp               = "gitlab.com/verygoodsoftwarenotvirus/todo/server/v1"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	files := map[string]*jen.File{
		"cmd/server/v1/coverage_test.go": coverageTestDotGo(),
		"cmd/server/v1/doc.go":           docDotGo(),
		"cmd/server/v1/main.go":          mainDotGo(),
		"cmd/server/v1/wire.go":          wireDotGo(),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
