package project

import (
	"log"
	"sync"
	"time"

	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/audit"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authentication"
	internalauth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authorization"
	buildserver "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/build/server"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/capitalism"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/capitalism/stripe"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/config/viper"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database"
	dbconfig "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querier"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding/mariadb"
	mockquerybuilding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding/postgres"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding/sqlite"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/encoding"
	mockencoding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/encoding/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/events"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/keys"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/logging"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/metrics"
	mockmetrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/metrics/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/tracing"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/panicking"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/random"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing/chi"
	mockrouting "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search/bleve"
	mocksearch "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/secrets"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/server"
	accountsservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/accounts"
	adminservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/admin"
	apiclientsservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/apiclients"
	auditservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/audit"
	authenticationservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/authentication"
	frontendservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/frontend"
	usersservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/users"
	webhooksservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/webhooks"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/storage"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads/images"
	mockuploads "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads/mock"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/server"
	configgencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/config_gen"
	datascaffoldercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/data_scaffolder"
	encodedqrcodegeneratorcmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/encoded_qr_code_generator"
	indexinitializercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/index_initializer"
	templategencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/template_gen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/composefiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/dockerfiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/providerconfigs"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/misc"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient/requests"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/converters"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/fakes"
	mocktypes "gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/mock"
	frontendtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/frontend"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/integration"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/load"
	testutils "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/utils"

	"github.com/gosuri/uiprogress"
)

const async = true

// RenderProject renders a project
func RenderProject(proj *naffmodels.Project) {
	packageRenderers := map[string]func(*naffmodels.Project) error{
		"servercmd":                 servercmd.RenderPackage,
		"templategencmd":            templategencmd.RenderPackage,
		"configgen":                 configgencmd.RenderPackage,
		"datascaffoldercmd":         datascaffoldercmd.RenderPackage,
		"encodedqrcodegeneratorcmd": encodedqrcodegeneratorcmd.RenderPackage,
		"indexinitializercmd":       indexinitializercmd.RenderPackage,
		"composefiles":              composefiles.RenderPackage,
		"deployfiles":               providerconfigs.RenderPackage,
		"dockerfiles":               dockerfiles.RenderPackage,
		"audit":                     audit.RenderPackage,
		"authentication":            authentication.RenderPackage,
		"authorization":             internalauth.RenderPackage,
		"buildserver":               buildserver.RenderPackage,
		"config":                    config.RenderPackage,
		"viper":                     viper.RenderPackage,
		"database":                  database.RenderPackage,
		"dbconfig":                  dbconfig.RenderPackage,
		"querier":                   querier.RenderPackage,
		"querybuilding":             querybuilding.RenderPackage,
		"mariadb":                   mariadb.RenderPackage,
		"postgres":                  postgres.RenderPackage,
		"sqlite":                    sqlite.RenderPackage,
		"dbmock":                    mockquerybuilding.RenderPackage,
		"capitalism":                capitalism.RenderPackage,
		"stripe":                    stripe.RenderPackage,
		"encoding":                  encoding.RenderPackage,
		"mockencoding":              mockencoding.RenderPackage,
		"events":                    events.RenderPackage,
		"observability":             observability.RenderPackage,
		"keys":                      keys.RenderPackage,
		"logging":                   logging.RenderPackage,
		"metrics":                   metrics.RenderPackage,
		"mockmetrics":               mockmetrics.RenderPackage,
		"tracing":                   tracing.RenderPackage,
		"panicking":                 panicking.RenderPackage,
		"random":                    random.RenderPackage,
		"routing":                   routing.RenderPackage,
		"chi":                       chi.RenderPackage,
		"mockrouting":               mockrouting.RenderPackage,
		"search":                    search.RenderPackage,
		"searchmock":                mocksearch.RenderPackage,
		"bleve":                     bleve.RenderPackage,
		"secrets":                   secrets.RenderPackage,
		"server":                    server.RenderPackage,
		"accountsservice":           accountsservice.RenderPackage,
		"adminservice":              adminservice.RenderPackage,
		"apiclientsservice":         apiclientsservice.RenderPackage,
		"auditservice":              auditservice.RenderPackage,
		"authenticationservice":     authenticationservice.RenderPackage,
		"frontendservice":           frontendservice.RenderPackage, // iterable service would go on next line
		"usersservice":              usersservice.RenderPackage,
		"webhooksservice":           webhooksservice.RenderPackage,
		"storage":                   storage.RenderPackage,
		"images":                    images.RenderPackage,
		"uploads":                   uploads.RenderPackage,
		"mockuploads":               mockuploads.RenderPackage,
		"httpclient":                httpclient.RenderPackage,
		"requests":                  requests.RenderPackage,
		"mocktypes":                 mocktypes.RenderPackage,
		"types":                     types.RenderPackage,
		"fakes":                     fakes.RenderPackage,
		"converters":                converters.RenderPackage,
		"miscellaneous":             misc.RenderPackage,
		"frontendtests":             frontendtests.RenderPackage,
		"integrationtests":          integration.RenderPackage,
		"loadtests":                 load.RenderPackage,
		"testutils":                 testutils.RenderPackage,
	}

	var wg sync.WaitGroup

	uiprogress.Start()
	progressBar := uiprogress.AddBar(len(packageRenderers)).PrependElapsed().AppendCompleted()

	wg.Add(1)
	if proj != nil {
		for name, x := range packageRenderers {
			if async {
				go renderTask(proj, &wg, name, x, progressBar)
			} else {
				renderTask(proj, &wg, name, x, progressBar)
			}
		}
	}

	// probably unnecessary?
	time.Sleep(2 * time.Second)
	wg.Done()
	wg.Wait()
}

func renderTask(proj *naffmodels.Project, wg *sync.WaitGroup, name string, renderer func(*naffmodels.Project) error, progressBar *uiprogress.Bar) {
	wg.Add(1)
	start := time.Now()

	if err := renderer(proj); err != nil {
		log.Panicf("error rendering %q after %s: %v\n", name, time.Since(start), err)
	}

	progressBar.Incr()
	wg.Done()
}
