package project

import (
	"log"
	"sync"
	"time"

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
	renderFunc func(string, []naffmodels.DataType) error
	activated  bool
}

type otherRenderHelper struct {
	renderFunc func(string, wordsmith.SuperPalabra, []naffmodels.DataType) error
	activated  bool
}

func RenderProject(in *naffmodels.Project) error {
	allActive := true

	packageRenderers := map[string]renderHelper{
		// completed
		"httpclient":       {renderFunc: httpclient.RenderPackage, activated: allActive},
		"configgen":        {renderFunc: configgen.RenderPackage, activated: allActive},
		"servercmd":        {renderFunc: servercmd.RenderPackage, activated: allActive},
		"twofactorcmd":     {renderFunc: twofactorcmd.RenderPackage, activated: allActive},
		"database":         {renderFunc: database.RenderPackage, activated: allActive},
		"internalauth":     {renderFunc: internalauth.RenderPackage, activated: allActive},
		"internalauthmock": {renderFunc: internalauthmock.RenderPackage, activated: allActive},
		"config":           {renderFunc: config.RenderPackage, activated: allActive},
		"encoding":         {renderFunc: encoding.RenderPackage, activated: allActive},
		"encodingmock":     {renderFunc: encodingmock.RenderPackage, activated: allActive},
		"metrics":          {renderFunc: metrics.RenderPackage, activated: allActive},
		"metricsmock":      {renderFunc: metricsmock.RenderPackage, activated: allActive},
		"server":           {renderFunc: server.RenderPackage, activated: allActive},
		"testutil":         {renderFunc: testutil.RenderPackage, activated: allActive},
		"testutilmock":     {renderFunc: testutilmock.RenderPackage, activated: allActive},
		"frontendtests":    {renderFunc: frontendtests.RenderPackage, activated: allActive},
		"webhooks":         {renderFunc: webhooks.RenderPackage, activated: allActive},
		"oauth2clients":    {renderFunc: oauth2clients.RenderPackage, activated: allActive},
		"frontend":         {renderFunc: frontend.RenderPackage, activated: allActive},
		"auth":             {renderFunc: auth.RenderPackage, activated: allActive},
		"users":            {renderFunc: users.RenderPackage, activated: allActive},
		"httpserver":       {renderFunc: httpserver.RenderPackage, activated: allActive},
		"modelsmock":       {renderFunc: modelsmock.RenderPackage, activated: allActive},
		"models":           {renderFunc: models.RenderPackage, activated: allActive},
		"randmodel":        {renderFunc: randmodel.RenderPackage, activated: allActive},
		"iterables":        {renderFunc: iterables.RenderPackage, activated: allActive},
		"dbclient":         {renderFunc: dbclient.RenderPackage, activated: allActive},
		"integrationtests": {renderFunc: integrationtests.RenderPackage, activated: allActive},
		"loadtests":        {renderFunc: loadtests.RenderPackage, activated: allActive},
		"queriers":         {renderFunc: queriers.RenderPackage, activated: allActive},
	}

	specialSnowflakes := map[string]otherRenderHelper{
		"composefiles":  {renderFunc: composefiles.RenderPackage, activated: allActive},
		"deployfiles":   {renderFunc: deploy.RenderPackage, activated: allActive},
		"dockerfiles":   {renderFunc: dockerfiles.RenderPackage, activated: allActive},
		"miscellaneous": {renderFunc: misc.RenderPackage, activated: allActive},
		"frontendmisc":  {renderFunc: frontendmisc.RenderPackage, activated: allActive},
		"frontendsrc":   {renderFunc: frontendsrc.RenderPackage, activated: allActive},
	}

	var wg sync.WaitGroup

	if in != nil {
		for name, x := range packageRenderers {
			if x.activated {
				wg.Add(1)
				go func(taskName string, renderer renderHelper) {
					start := time.Now()
					if err := renderer.renderFunc(in.OutputPath, in.DataTypes); err != nil {
						log.Printf("error rendering %q after %s\n", taskName, time.Since(start))
					}
					log.Printf("rendered %s after %s\n", taskName, time.Since(start))
					wg.Done()
				}(name, x)
			}
		}
		for name, x := range specialSnowflakes {
			if x.activated {
				wg.Add(1)
				go func(taskName string, packageName wordsmith.SuperPalabra, renderer otherRenderHelper) {
					start := time.Now()
					if err := renderer.renderFunc(in.OutputPath, packageName, in.DataTypes); err != nil {
						log.Printf("error rendering %q after %s\n", taskName, time.Since(start))
					}
					log.Printf("rendered %s after %s\n", taskName, time.Since(start))
					wg.Done()
				}(name, in.Name, x)
			}
		}
	}

	wg.Wait()

	return nil
}
