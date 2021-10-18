package project

import (
	"log"
	"sync"
	"time"

	naffmodels "gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authentication"
	internalmockauthentication "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authentication/mock"
	internalauth "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/authorization"
	buildserver "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/build/server"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database"
	dbconfig "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/config"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/queriers/mysql"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/database/queriers/postgres"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/encoding"
	mockencoding "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/encoding/mock"
	msgqconfig "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/messagequeue/config"
	msgqconsumers "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/messagequeue/consumers"
	msgqconsumersmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/messagequeue/consumers/mock"
	msgqpublishers "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/messagequeue/publishers"
	msgqpublishersmock "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/messagequeue/publishers/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/keys"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/logging"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/metrics"
	mockmetrics "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/metrics/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/observability/tracing"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/panicking"
	mockpanicking "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/panicking/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/random"
	mockrandom "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/random/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing/chi"
	mockrouting "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/routing/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search/elasticsearch"
	mocksearch "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/search/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/secrets"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/server"
	accountsservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/accounts"
	adminservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/admin"
	apiclientsservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/apiclients"
	authenticationservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/authentication"
	frontendservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/frontend"
	iterablesservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/iterables"
	usersservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/users"
	webhooksservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/webhooks"
	websocketsservice "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/services/websockets"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/storage"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads/images"
	mockuploads "gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/uploads/mock"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/_internal_/workers"
	servercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/server"
	configgencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/config_gen"
	datascaffoldercmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/data_scaffolder"
	encodedqrcodegeneratorcmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/encoded_qr_code_generator"
	templategencmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/tools/template_gen"
	workerscmd "gitlab.com/verygoodsoftwarenotvirus/naff/templates/cmd/workers"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/composefiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/dockerfiles"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/environments/providerconfigs"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/misc"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/client/httpclient/requests"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/fakes"
	mocktypes "gitlab.com/verygoodsoftwarenotvirus/naff/templates/pkg/types/mock"
	frontendtests "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/frontend"
	"gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/integration"
	testutils "gitlab.com/verygoodsoftwarenotvirus/naff/templates/tests/utils"

	"github.com/gosuri/uiprogress"
)

const async = true

// RenderProject renders a project
func RenderProject(proj *naffmodels.Project) {
	packageRenderers := map[string]func(*naffmodels.Project) error{
		"servercmd":                 servercmd.RenderPackage,
		"workerscmd":                workerscmd.RenderPackage,
		"templategencmd":            templategencmd.RenderPackage,
		"configgen":                 configgencmd.RenderPackage,
		"datascaffoldercmd":         datascaffoldercmd.RenderPackage,
		"encodedqrcodegeneratorcmd": encodedqrcodegeneratorcmd.RenderPackage,
		"composefiles":              composefiles.RenderPackage,
		"deployfiles":               providerconfigs.RenderPackage,
		"dockerfiles":               dockerfiles.RenderPackage,
		"authentication":            authentication.RenderPackage,
		"authenticationmock":        internalmockauthentication.RenderPackage,
		"authorization":             internalauth.RenderPackage,
		"buildserver":               buildserver.RenderPackage,
		"config":                    config.RenderPackage,
		"database":                  database.RenderPackage,
		"dbconfig":                  dbconfig.RenderPackage,
		"encoding":                  encoding.RenderPackage,
		"mockencoding":              mockencoding.RenderPackage,
		"observability":             observability.RenderPackage,
		"keys":                      keys.RenderPackage,
		"logging":                   logging.RenderPackage,
		"metrics":                   metrics.RenderPackage,
		"mockmetrics":               mockmetrics.RenderPackage,
		"tracing":                   tracing.RenderPackage,
		"panicking":                 panicking.RenderPackage,
		"mockpanicking":             mockpanicking.RenderPackage,
		"random":                    random.RenderPackage,
		"mockrandom":                mockrandom.RenderPackage,
		"routing":                   routing.RenderPackage,
		"chi":                       chi.RenderPackage,
		"mockrouting":               mockrouting.RenderPackage,
		"mysql":                     mysql.RenderPackage,
		"postgres":                  postgres.RenderPackage,
		"elasticsearch":             elasticsearch.RenderPackage,
		"workers":                   workers.RenderPackage,
		"websocketsservice":         websocketsservice.RenderPackage,
		"search":                    search.RenderPackage,
		"searchmock":                mocksearch.RenderPackage,
		"secrets":                   secrets.RenderPackage,
		"server":                    server.RenderPackage,
		"accountsservice":           accountsservice.RenderPackage,
		"adminservice":              adminservice.RenderPackage,
		"apiclientsservice":         apiclientsservice.RenderPackage,
		"authenticationservice":     authenticationservice.RenderPackage,
		"msgqconfig":                msgqconfig.RenderPackage,
		"msgqconsumers":             msgqconsumers.RenderPackage,
		"msgqconsumersmock":         msgqconsumersmock.RenderPackage,
		"msgqpublishers":            msgqpublishers.RenderPackage,
		"msgqpublishersmock":        msgqpublishersmock.RenderPackage,
		"frontendservice":           frontendservice.RenderPackage,
		"iterablesservice":          iterablesservice.RenderPackage,
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
		"miscellaneous":             misc.RenderPackage,
		"frontendtests":             frontendtests.RenderPackage,
		"integrationtests":          integration.RenderPackage,
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
