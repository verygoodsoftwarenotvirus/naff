package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

var (
	// files are all the available files to generate
	files = map[string]*jen.File{
		"internal/config/v1/doc.go":          docDotGo(),
		"internal/config/v1/config.go":       configDotGo(),
		"internal/config/v1/config_test.go":  configTestDotGo(),
		"internal/config/v1/database.go":     databaseDotGo(),
		"internal/config/v1/metrics.go":      metricsDotGo(),
		"internal/config/v1/metrics_test.go": metricsTestDotGo(),
		"internal/config/v1/wire.go":         wireDotGo(),
	}
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) {
	for path, file := range files {
		renderFile(path, file)
	}
}

func renderFile(path string, file *jen.File) {
	fp := utils.BuildTemplatePath(path)
	_ = os.Remove(fp)

	var b bytes.Buffer
	if err := file.Render(&b); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fp, b.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}

if err := utils.RunGoimportsForFile(fp); err != nil {
	log.Println(err)
	}
}
