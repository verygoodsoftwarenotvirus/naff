package authentication

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "authentication"

	basePackagePath = "internal/services/authentication"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"service_test.go":                     serviceTestDotGo(proj),
		"session_manager_test.go":             sessionManagerTestDotGo(proj),
		"middleware_test.go":                  middlewareTestDotGo(proj),
		"session_manager.go":                  sessionManagerDotGo(proj),
		"wire.go":                             wireDotGo(proj),
		"doc.go":                              docDotGo(proj),
		"helpers.go":                          helpersDotGo(proj),
		"middleware.go":                       middlewareDotGo(proj),
		"service.go":                          serviceDotGo(proj),
		"config.go":                           configDotGo(proj),
		"helpers_test.go":                     helpersTestDotGo(proj),
		"http_routes_test.go":                 httpRoutesTestDotGo(proj),
		"mock_cookie_encoder_decoder_test.go": mockCookieEncoderDecoderTestDotGo(proj),
		"config_test.go":                      configTestDotGo(proj),
		"http_helpers_test.go":                httpHelpersTestDotGo(proj),
		"http_routes.go":                      httpRoutesDotGo(proj),
	}

	//for _, typ := range types {
	//	files[fmt.Sprintf("%s.go", typ.Name.PluralRouteName)] = itemsDotGo(typ)
	//	files[fmt.Sprintf("%s_test.go", typ.Name.PluralRouteName)] = itemsTestDotGo(typ)
	//}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
