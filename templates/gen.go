package project

import (
	"log"

	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"

	// completed
	httpclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/client/v1/http"
	configgen "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/cmd/config_gen/v1"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/cmd/server/v1"
	twofactorcmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/cmd/tools/two_factor"
	database "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/database/v1"

	// to do
	dbclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/database/v1/client"
	mariaDB "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/database/v1/queriers/mariadb"
	postgresql "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/database/v1/queriers/postgres"
	sqlite3 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/database/v1/queriers/sqlite"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/models/v1"
	modelsmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/models/v1/mock"
	internalauth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/auth"
	internalauthmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/auth/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/encoding"
	encodingmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/metrics"
	metricsmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/notreallyinternal/v1/metrics/mock"
	server "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/server/v1"
	httpserver "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/server/v1/http"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/services/v1/webhooks"
	frontendtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/tests/v1/frontend"
	integrationtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/tests/v1/integration"
	loadtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/tests/v1/load"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/tests/v1/testutil"
	testutilmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/tests/v1/testutil/mock"
	randmodel "gitlab.com/verygoodsoftwarenotvirus/naff/templates/experimental/tests/v1/testutil/rand/model"
)

func RenderProject(in *naffmodels.Project) error {

	packageRenderers := map[string]struct {
		renderFunc func([]naffmodels.DataType) error
		activated  bool
	}{
		// completed
		"httpclient": {
			renderFunc: httpclient.RenderPackage,
			activated:  false,
		},
		"configgen": {
			renderFunc: configgen.RenderPackage,
			activated:  false,
		},
		"servercmd": {
			renderFunc: servercmd.RenderPackage,
			activated:  false,
		},
		"twofactorcmd": {
			renderFunc: twofactorcmd.RenderPackage,
			activated:  false,
		},
		"database": {
			renderFunc: database.RenderPackage,
			activated:  false,
		},
		// to do
		"internalauth": {
			renderFunc: internalauth.RenderPackage,
			activated:  true,
		},
		"internalauthmock": {
			renderFunc: internalauthmock.RenderPackage,
			activated:  false,
		},
		"config": {
			renderFunc: config.RenderPackage,
			activated:  false,
		},
		"encoding": {
			renderFunc: encoding.RenderPackage,
			activated:  false,
		},
		"encodingmock": {
			renderFunc: encodingmock.RenderPackage,
			activated:  false,
		},
		"metrics": {
			renderFunc: metrics.RenderPackage,
			activated:  false,
		},
		"metricsmock": {
			renderFunc: metricsmock.RenderPackage,
			activated:  false,
		},
		"models": {
			renderFunc: models.RenderPackage,
			activated:  false,
		},
		"modelsmock": {
			renderFunc: modelsmock.RenderPackage,
			activated:  false,
		},
		"dbclient": {
			renderFunc: dbclient.RenderPackage,
			activated:  false,
		},
		"mariaDB": {
			renderFunc: mariaDB.RenderPackage,
			activated:  false,
		},
		"postgresql": {
			renderFunc: postgresql.RenderPackage,
			activated:  false,
		},
		"sqlite3": {
			renderFunc: sqlite3.RenderPackage,
			activated:  false,
		},
		"server": {
			renderFunc: server.RenderPackage,
			activated:  false,
		},
		"httpserver": {
			renderFunc: httpserver.RenderPackage,
			activated:  false,
		},
		"auth": {
			renderFunc: auth.RenderPackage,
			activated:  false,
		},
		"frontend": {
			renderFunc: frontend.RenderPackage,
			activated:  false,
		},
		"items": {
			renderFunc: items.RenderPackage,
			activated:  false,
		},
		"oauth2clients": {
			renderFunc: oauth2clients.RenderPackage,
			activated:  false,
		},
		"users": {
			renderFunc: users.RenderPackage,
			activated:  false,
		},
		"webhooks": {
			renderFunc: webhooks.RenderPackage,
			activated:  false,
		},
		"frontendtests": {
			renderFunc: frontendtests.RenderPackage,
			activated:  false,
		},
		"integrationtests": {
			renderFunc: integrationtests.RenderPackage,
			activated:  false,
		},
		"loadtests": {
			renderFunc: loadtests.RenderPackage,
			activated:  false,
		},
		"testutil": {
			renderFunc: testutil.RenderPackage,
			activated:  false,
		},
		"testutilmock": {
			renderFunc: testutilmock.RenderPackage,
			activated:  false,
		},
		"randmodel": {
			renderFunc: randmodel.RenderPackage,
			activated:  false,
		},
	}

	if in != nil {
		for name, x := range packageRenderers {
			if x.activated {
				if err := x.renderFunc(in.DataTypes); err != nil {
					log.Printf("error rendering %q", name)
					return err
				}
			}
		}
	}

	return nil
}
