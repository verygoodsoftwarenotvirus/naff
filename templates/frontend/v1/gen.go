package frontendv1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"log"
	"os"
	"path/filepath"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]func() string{
		"frontend/v1/package.json": packageDotJSON,
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

	return os.MkdirAll(utils.BuildTemplatePath(proj.OutputPath, "frontend/v1/public"), os.ModePerm)
}

const pdjson = `{
  "name": "frontend",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "build": "echo \"this is where I would build if I could!\" && exit 0"
  },
  "author": "",
  "license": "ISC"
}
`

func packageDotJSON() string {
	return pdjson
}
