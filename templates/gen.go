package project

import (
	"log"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"

	httpclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/client/v1/http"
	configgen "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/config_gen/v1"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/server/v1"
	indexinitializercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/index_initializer"
	twofactorcmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/two_factor"
	database "gitlab.com/verygoodsoftwarenotvirus/naff/templates/database/v1"
	dbclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/database/v1/client"
	queriers "gitlab.com/verygoodsoftwarenotvirus/naff/templates/database/v1/queriers"
	composefiles "gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/composefiles"
	deploy "gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/deploy"
	dockerfiles "gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/dockerfiles"
	frontendmisc "gitlab.com/verygoodsoftwarenotvirus/naff/templates/frontend/v1"
	frontendsrc "gitlab.com/verygoodsoftwarenotvirus/naff/templates/frontend/v1/src"
	internalauth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/auth"
	internalauthmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/auth/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/encoding"
	encodingmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/metrics"
	metricsmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/metrics/mock"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/search"
	bleve "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/search/bleve"
	searchmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/search/mock"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/templates/iinternal/v1/tracing"
	misc "gitlab.com/verygoodsoftwarenotvirus/naff/templates/misc"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/templates/models/v1"
	fakemodels "gitlab.com/verygoodsoftwarenotvirus/naff/templates/models/v1/fake"
	modelsmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/models/v1/mock"
	server "gitlab.com/verygoodsoftwarenotvirus/naff/templates/server/v1"
	httpserver "gitlab.com/verygoodsoftwarenotvirus/naff/templates/server/v1/http"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/templates/services/v1/frontend"
	iterables "gitlab.com/verygoodsoftwarenotvirus/naff/templates/services/v1/iterables"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/templates/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/templates/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/templates/services/v1/webhooks"
	frontendtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/v1/frontend"
	integrationtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/v1/integration"
	loadtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/v1/load"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/v1/testutil"
)

type renderHelper struct {
	name       string
	renderFunc func(*naffmodels.Project) error
	activated  bool
}

const async = true

// RenderProject renders a project
func RenderProject(proj *naffmodels.Project) error {
	allActive := true
	searchActive := false

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled && !searchActive {
			searchActive = true
			break
		}
	}

	packageRenderers := []renderHelper{
		{name: "httpclient", renderFunc: httpclient.RenderPackage, activated: allActive},
		{name: "configgen", renderFunc: configgen.RenderPackage, activated: allActive},
		{name: "servercmd", renderFunc: servercmd.RenderPackage, activated: allActive},
		{name: "twofactorcmd", renderFunc: twofactorcmd.RenderPackage, activated: allActive},
		{name: "indexinitializercmd", renderFunc: indexinitializercmd.RenderPackage, activated: searchActive},
		{name: "database", renderFunc: database.RenderPackage, activated: allActive},
		{name: "internalauth", renderFunc: internalauth.RenderPackage, activated: allActive},
		{name: "internalauthmock", renderFunc: internalauthmock.RenderPackage, activated: allActive},
		{name: "config", renderFunc: config.RenderPackage, activated: allActive},
		{name: "encoding", renderFunc: encoding.RenderPackage, activated: allActive},
		{name: "encodingmock", renderFunc: encodingmock.RenderPackage, activated: allActive},
		{name: "metrics", renderFunc: metrics.RenderPackage, activated: allActive},
		{name: "tracing", renderFunc: tracing.RenderPackage, activated: allActive},
		{name: "search", renderFunc: search.RenderPackage, activated: searchActive},
		{name: "searchmock", renderFunc: searchmock.RenderPackage, activated: searchActive},
		{name: "bleve", renderFunc: bleve.RenderPackage, activated: searchActive},
		{name: "metricsmock", renderFunc: metricsmock.RenderPackage, activated: allActive},
		{name: "server", renderFunc: server.RenderPackage, activated: allActive},
		{name: "testutil", renderFunc: testutil.RenderPackage, activated: allActive},
		{name: "frontendtests", renderFunc: frontendtests.RenderPackage, activated: allActive},
		{name: "webhooks", renderFunc: webhooks.RenderPackage, activated: allActive},
		{name: "oauth2clients", renderFunc: oauth2clients.RenderPackage, activated: allActive},
		{name: "frontend", renderFunc: frontend.RenderPackage, activated: allActive},
		{name: "auth", renderFunc: auth.RenderPackage, activated: allActive},
		{name: "users", renderFunc: users.RenderPackage, activated: allActive},
		{name: "httpserver", renderFunc: httpserver.RenderPackage, activated: allActive},
		{name: "modelsmock", renderFunc: modelsmock.RenderPackage, activated: allActive},
		{name: "models", renderFunc: models.RenderPackage, activated: allActive},
		{name: "fakemodels", renderFunc: fakemodels.RenderPackage, activated: allActive},
		{name: "iterables", renderFunc: iterables.RenderPackage, activated: allActive},
		{name: "dbclient", renderFunc: dbclient.RenderPackage, activated: allActive},
		{name: "integrationtests", renderFunc: integrationtests.RenderPackage, activated: allActive},
		{name: "loadtests", renderFunc: loadtests.RenderPackage, activated: allActive},
		{name: "queriers", renderFunc: queriers.RenderPackage, activated: allActive},
		{name: "composefiles", renderFunc: composefiles.RenderPackage, activated: allActive},
		{name: "deployfiles", renderFunc: deploy.RenderPackage, activated: allActive},
		{name: "dockerfiles", renderFunc: dockerfiles.RenderPackage, activated: allActive},
		{name: "miscellaneous", renderFunc: misc.RenderPackage, activated: allActive},
		{name: "frontendmisc", renderFunc: frontendmisc.RenderPackage, activated: allActive},
		{name: "frontendsrc", renderFunc: frontendsrc.RenderPackage, activated: allActive},
	}

	var wg sync.WaitGroup

	uiprogress.Start()
	progressBar := uiprogress.AddBar(len(packageRenderers)).PrependElapsed().AppendCompleted()

	wg.Add(1)
	if proj != nil {
		for _, x := range packageRenderers {
			if x.activated {
				if async {
					go renderTask(proj, &wg, x, progressBar)
				} else {
					renderTask(proj, &wg, x, progressBar)
				}
			}
		}
	}

	// probably unnecessary?
	time.Sleep(2 * time.Second)
	wg.Done()
	wg.Wait()

	return nil
}

func renderTask(proj *naffmodels.Project, wg *sync.WaitGroup, renderer renderHelper, progressBar *uiprogress.Bar) {
	wg.Add(1)
	start := time.Now()
	if err := renderer.renderFunc(proj); err != nil {
		log.Panicf("error rendering %q after %s: %v\n", renderer.name, time.Since(start), err)
	}
	progressBar.Incr()
	wg.Done()
}
