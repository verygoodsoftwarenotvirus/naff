package project

import (
	authorization2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authorization"
	capitalism2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/capitalism"
	stripe2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/capitalism/stripe"
	viper2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/config/viper"
	config2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/config"
	mock3 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding/mock"
	mock4 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/encoding/mock"
	observability2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability"
	keys2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/keys"
	logging2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/logging"
	zerolog2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/logging/zerolog"
	panicking2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/panicking"
	random2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/random"
	routing2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing"
	chi2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing/chi"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/accounts"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/admin"
	audit2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/audit"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/frontend"
	storage2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/storage"
	uploads2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads"
	images2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads/images"
	mock2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/data_scaffolder"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/encoded_qr_code_generator"
	converters2 "gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/converters"
	"log"
	"sync"
	"time"

	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/audit"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authentication"
	buildserver "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/build/server"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database"
	dbclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querier"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/querybuilding"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/encoding"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/metrics"
	mockmetrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/metrics/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search/bleve"
	mocksearch "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/server"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/apiclients"
	authnservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/authentication"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/iterables"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/users"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/webhooks"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/tracing"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/server"
	configgencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/config_gen"
	indexinitializercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/index_initializer"
	templategencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/template_gen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/composefiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/deploy"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/dockerfiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/misc"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient/requests"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/fake"
	mocktypes "gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/mock"
	frontendtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/frontend"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/integration"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/load"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/testutil"

	"github.com/gosuri/uiprogress"
)

const async = true

// RenderProject renders a project
func RenderProject(proj *naffmodels.Project) {
	packageRenderers := map[string]func(*naffmodels.Project) error{
		"httpclient":          client.RenderPackage,
		"requests":            requests.RenderPackage,
		"configgen":           configgencmd.RenderPackage,
		"servercmd":           servercmd.RenderPackage,
		"templategencmd":      templategencmd.RenderPackage,
		"indexinitializercmd": indexinitializercmd.RenderPackage,
		"database":            database.RenderPackage,
		"internalauth":        authentication.RenderPackage,
		"config":              config.RenderPackage,
		"encoding":            encoding.RenderPackage,
		"metrics":             metrics.RenderPackage,
		"tracing":             tracing.RenderPackage,
		"search":              search.RenderPackage,
		"audit":               audit.RenderPackage,
		"searchmock":          mocksearch.RenderPackage,
		"bleve":               bleve.RenderPackage,
		"metricsmock":         mockmetrics.RenderPackage,
		"server":              server.RenderPackage,
		"testutil":            testutil.RenderPackage,
		"frontendtests":       frontendtests.RenderPackage,
		"webhooks":            webhooks.RenderPackage,
		"oauth2clients":       apiclients.RenderPackage,
		"auth":                authnservice.RenderPackage,
		"users":               users.RenderPackage,
		"httpserver":          server.RenderPackage,
		"mocktypes":           mocktypes.RenderPackage,
		"models":              types.RenderPackage,
		"fakemodels":          fake.RenderPackage,
		"iterables":           iterables.RenderPackage,
		"dbclient":            dbclient.RenderPackage,
		"integrationtests":    integration.RenderPackage,
		"loadtests":           load.RenderPackage,
		"queriers":            querybuilding.RenderPackage,
		"composefiles":        composefiles.RenderPackage,
		"deployfiles":         deploy.RenderPackage,
		"dockerfiles":         dockerfiles.RenderPackage,
		"miscellaneous":       misc.RenderPackage,
		//
		// experimental
		"buildserver":               buildserver.RenderPackage,
		"authorization":             authorization2.RenderPackage,
		"capitalism":                capitalism2.RenderPackage,
		"stripe":                    stripe2.RenderPackage,
		"viper":                     viper2.RenderPackage,
		"dbconfig":                  config2.RenderPackage,
		"dbmock":                    mock3.RenderPackage,
		"mockencoding":              mock4.RenderPackage,
		"observability":             observability2.RenderPackage,
		"keys":                      keys2.RenderPackage,
		"logging":                   logging2.RenderPackage,
		"zerolog":                   zerolog2.RenderPackage,
		"panicking":                 panicking2.RenderPackage,
		"random":                    random2.RenderPackage,
		"routing":                   routing2.RenderPackage,
		"chi":                       chi2.RenderPackage,
		"mockrouting":               mock.RenderPackage,
		"accountsservice":           accounts.RenderPackage,
		"adminservice":              admin.RenderPackage,
		"auditservice":              audit2.RenderPackage,
		"frontendservice":           frontend.RenderPackage,
		"storage":                   storage2.RenderPackage,
		"uploads":                   uploads2.RenderPackage,
		"images":                    images2.RenderPackage,
		"mockuploads":               mock2.RenderPackage,
		"datascaffoldercmd":         data_scaffolder.RenderPackage,
		"encodedqrcodegeneratorcmd": encoded_qr_code_generator.RenderPackage,
		"converters":                converters2.RenderPackage,
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
