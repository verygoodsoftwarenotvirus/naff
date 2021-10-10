package authentication

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/services/authentication"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"doc.go":                              docDotGo(proj),
		"http_routes_test.go":                 httpRoutesTestDotGo(proj),
		"service_test.go":                     serviceTestDotGo(proj),
		"config.go":                           configDotGo(proj),
		"helpers_test.go":                     helpersTestDotGo(proj),
		"http_helpers_test.go":                httpHelpersTestDotGo(proj),
		"service.go":                          serviceDotGo(proj),
		"session_manager_test.go":             sessionManagerTestDotGo(proj),
		"middleware_test.go":                  middlewareTestDotGo(proj),
		"mock_cookie_encoder_decoder_test.go": mockCookieEncoderDecoderTestDotGo(proj),
		"wire.go":                             wireDotGo(proj),
		"config_test.go":                      configTestDotGo(proj),
		"helpers.go":                          helpersDotGo(proj),
		"http_routes.go":                      httpRoutesDotGo(proj),
		"middleware.go":                       middlewareDotGo(proj),
		"session_manager.go":                  sessionManagerDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed http_routes_test.gotpl
var httpRoutesTestTemplate string

func httpRoutesTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpRoutesTestTemplate, nil)
}

//go:embed service_test.gotpl
var serviceTestTemplate string

func serviceTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serviceTestTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed helpers_test.gotpl
var helpersTestTemplate string

func helpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTestTemplate, nil)
}

//go:embed http_helpers_test.gotpl
var httpHelpersTestTemplate string

func httpHelpersTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpHelpersTestTemplate, nil)
}

//go:embed service.gotpl
var serviceTemplate string

func serviceDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, serviceTemplate, nil)
}

//go:embed session_manager_test.gotpl
var sessionManagerTestTemplate string

func sessionManagerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, sessionManagerTestTemplate, nil)
}

//go:embed middleware_test.gotpl
var middlewareTestTemplate string

func middlewareTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, middlewareTestTemplate, nil)
}

//go:embed mock_cookie_encoder_decoder_test.gotpl
var mockCookieEncoderDecoderTestTemplate string

func mockCookieEncoderDecoderTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockCookieEncoderDecoderTestTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed helpers.gotpl
var helpersTemplate string

func helpersDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, helpersTemplate, nil)
}

//go:embed http_routes.gotpl
var httpRoutesTemplate string

func httpRoutesDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, httpRoutesTemplate, nil)
}

//go:embed middleware.gotpl
var middlewareTemplate string

func middlewareDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, middlewareTemplate, nil)
}

//go:embed session_manager.gotpl
var sessionManagerTemplate string

func sessionManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, sessionManagerTemplate, nil)
}
