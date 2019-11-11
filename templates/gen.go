package project

import (
	"log"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"

	// completed
	httpclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/client/v1/http"
	configgen "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/cmd/config_gen/v1"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/cmd/server/v1"
	twofactorcmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/cmd/tools/two_factor"
	composefiles "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/composefiles"
	database "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/database/v1"
	dbclient "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/database/v1/client"
	queriers "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/database/v1/queriers"
	deploy "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/deploy"
	dockerfiles "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/dockerfiles"
	frontendmisc "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/frontend/v1"
	frontendsrc "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/frontend/v1/src"
	internalauth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/auth"
	internalauthmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/auth/mock"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/encoding"
	encodingmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/encoding/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/metrics"
	metricsmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/iinternal/v1/metrics/mock"
	misc "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/misc"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/models/v1"
	modelsmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/models/v1/mock"
	server "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/server/v1"
	httpserver "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/server/v1/http"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/services/v1/frontend"
	iterables "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/services/v1/iterables"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/services/v1/webhooks"
	frontendtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/tests/v1/frontend"
	integrationtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/tests/v1/integration"
	loadtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/tests/v1/load"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/tests/v1/testutil"
	testutilmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/tests/v1/testutil/mock"
	randmodel "gitlab.com/verygoodsoftwarenotvirus/naff/templates/blessed/tests/v1/testutil/rand/model"
)

type renderHelper struct {
	name       string
	renderFunc func(string, []naffmodels.DataType) error
	activated  bool
}

type otherRenderHelper struct {
	name       string
	renderFunc func(string, wordsmith.SuperPalabra, []naffmodels.DataType) error
	activated  bool
}

func RenderProject(in *naffmodels.Project) error {
	allActive := true

	packageRenderers := []renderHelper{
		// completed
		{name: "httpclient", renderFunc: httpclient.RenderPackage, activated: allActive},
		{name: "configgen", renderFunc: configgen.RenderPackage, activated: allActive},
		{name: "servercmd", renderFunc: servercmd.RenderPackage, activated: allActive},
		{name: "twofactorcmd", renderFunc: twofactorcmd.RenderPackage, activated: allActive},
		{name: "database", renderFunc: database.RenderPackage, activated: allActive},
		{name: "internalauth", renderFunc: internalauth.RenderPackage, activated: allActive},
		{name: "internalauthmock", renderFunc: internalauthmock.RenderPackage, activated: allActive},
		{name: "config", renderFunc: config.RenderPackage, activated: allActive},
		{name: "encoding", renderFunc: encoding.RenderPackage, activated: allActive},
		{name: "encodingmock", renderFunc: encodingmock.RenderPackage, activated: allActive},
		{name: "metrics", renderFunc: metrics.RenderPackage, activated: allActive},
		{name: "metricsmock", renderFunc: metricsmock.RenderPackage, activated: allActive},
		{name: "server", renderFunc: server.RenderPackage, activated: allActive},
		{name: "testutil", renderFunc: testutil.RenderPackage, activated: allActive},
		{name: "testutilmock", renderFunc: testutilmock.RenderPackage, activated: allActive},
		{name: "frontendtests", renderFunc: frontendtests.RenderPackage, activated: allActive},
		{name: "webhooks", renderFunc: webhooks.RenderPackage, activated: allActive},
		{name: "oauth2clients", renderFunc: oauth2clients.RenderPackage, activated: allActive},
		{name: "frontend", renderFunc: frontend.RenderPackage, activated: allActive},
		{name: "auth", renderFunc: auth.RenderPackage, activated: allActive},
		{name: "users", renderFunc: users.RenderPackage, activated: allActive},
		{name: "httpserver", renderFunc: httpserver.RenderPackage, activated: allActive},
		{name: "modelsmock", renderFunc: modelsmock.RenderPackage, activated: allActive},
		{name: "models", renderFunc: models.RenderPackage, activated: allActive},
		{name: "randmodel", renderFunc: randmodel.RenderPackage, activated: allActive},
		{name: "iterables", renderFunc: iterables.RenderPackage, activated: allActive},
		{name: "dbclient", renderFunc: dbclient.RenderPackage, activated: allActive},
		{name: "integrationtests", renderFunc: integrationtests.RenderPackage, activated: allActive},
		{name: "loadtests", renderFunc: loadtests.RenderPackage, activated: allActive},
		{name: "queriers", renderFunc: queriers.RenderPackage, activated: allActive},
	}

	specialSnowflakes := []otherRenderHelper{
		{name: "composefiles", renderFunc: composefiles.RenderPackage, activated: allActive},
		{name: "deployfiles", renderFunc: deploy.RenderPackage, activated: allActive},
		{name: "dockerfiles", renderFunc: dockerfiles.RenderPackage, activated: allActive},
		{name: "miscellaneous", renderFunc: misc.RenderPackage, activated: allActive},
		{name: "frontendmisc", renderFunc: frontendmisc.RenderPackage, activated: allActive},
		{name: "frontendsrc", renderFunc: frontendsrc.RenderPackage, activated: allActive},
	}

	var wg sync.WaitGroup

	uiprogress.Start()
	progressBar := uiprogress.AddBar(len(packageRenderers) + len(specialSnowflakes)).PrependElapsed().AppendCompleted()

	if in != nil {
		for _, x := range packageRenderers {
			if x.activated {
				wg.Add(1)
				go func(taskName string, renderer renderHelper) {
					start := time.Now()
					if err := renderer.renderFunc(in.OutputPath, in.DataTypes); err != nil {
						log.Printf("error rendering %q after %s: %v\n", taskName, time.Since(start), err)
					}
					// fmt.Printf("rendered %s after %s\n", taskName, time.Since(start))
					progressBar.Incr()
					wg.Done()
				}(x.name, x)
			}
		}
		for _, x := range specialSnowflakes {
			if x.activated {
				wg.Add(1)
				go func(taskName string, packageName wordsmith.SuperPalabra, renderer otherRenderHelper) {
					start := time.Now()
					if err := renderer.renderFunc(in.OutputPath, packageName, in.DataTypes); err != nil {
						log.Printf("error rendering %q after %s\n", taskName, time.Since(start))
					}
					// fmt.Printf("rendered %s after %s\n", taskName, time.Since(start))
					progressBar.Incr()
					wg.Done()
				}(x.name, in.Name, x)
			}
		}
	}

	wg.Wait()

	return nil
}
