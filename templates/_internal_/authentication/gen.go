package authentication

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/authentication"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"authenticator.go":      authenticatorDotGoString(proj),
		"authenticator_test.go": authenticatorTestDotGoString(proj),
		"doc.go":                docDotGoString(proj),
		"wire.go":               wireDotGoString(proj),
		"mock_authenticator.go": mockAuthenticatorDotGoString(proj),
		"config.go":             configDotGoString(proj),
		"config_test.go":        configTestDotGoString(proj),
		"argon2.go":             argon2DotGoString(proj),
		"argon2_test.go":        argon2TestDotGoString(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

//go:embed authenticator.gotpl
var authenticatorTemplate string

func authenticatorDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, authenticatorTemplate, nil)
}

//go:embed authenticator_test.gotpl
var authenticatorTestTemplate string

func authenticatorTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, authenticatorTestTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed wire.gotpl
var wireTemplate string

func wireDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, wireTemplate, nil)
}

//go:embed mock_authenticator.gotpl
var mockAuthenticatorTemplate string

func mockAuthenticatorDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, mockAuthenticatorTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed argon2.gotpl
var argon2Template string

func argon2DotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, argon2Template, nil)
}

//go:embed argon2_test.gotpl
var argon2TestTemplate string

func argon2TestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, argon2TestTemplate, nil)
}
