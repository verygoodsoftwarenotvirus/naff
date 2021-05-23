package project

import (
	"log"
	"sync"
	"time"

	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/server"
	configgencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/config_gen"
	indexinitializercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/index_initializer"
	templategencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/template_gen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/composefiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/deploy"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/dockerfiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/authentication"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/database"
	dbclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/database/querier"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/database/querybuilding"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/encoding"
	mockencoding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/encoding/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/metrics"
	mockmetrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/metrics/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/search"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/search/bleve"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/search/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/server"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/services/apiclients"
	authnservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/services/authentication"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/services/frontend"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/services/iterables"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/services/users"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/services/webhooks"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/internal/tracing"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/misc"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient/requests"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/fake"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/mock"
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
		"encodingmock":        mockencoding.RenderPackage,
		"metrics":             metrics.RenderPackage,
		"tracing":             tracing.RenderPackage,
		"search":              search.RenderPackage,
		"searchmock":          mocksearch.RenderPackage,
		"bleve":               bleve.RenderPackage,
		"metricsmock":         mockmetrics.RenderPackage,
		"server":              server.RenderPackage,
		"testutil":            testutil.RenderPackage,
		"frontendtests":       frontendtests.RenderPackage,
		"webhooks":            webhooks.RenderPackage,
		"oauth2clients":       apiclients.RenderPackage,
		"frontend":            frontend.RenderPackage,
		"auth":                authnservice.RenderPackage,
		"users":               users.RenderPackage,
		"httpserver":          server.RenderPackage,
		"modelsmock":          mock.RenderPackage,
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
