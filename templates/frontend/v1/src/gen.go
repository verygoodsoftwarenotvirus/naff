package frontendsrc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]func() string{
		"frontend/v1/src/models/fakes.ts": fakesDotTS,
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("frontend/v1/src/models/%s.ts", typ.Name.UnexportedVarName())] = buildSomethingFrontendModels(typ)
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(proj.OutputPath, filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return err
		}

		if _, err := f.WriteString(file()); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}
