package metrics

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
		"internal/metrics/v1/doc.go":          docDotGo(),
		"internal/metrics/v1/meta.go":         metaDotGo(),
		"internal/metrics/v1/meta_test.go":    metaTestDotGo(),
		"internal/metrics/v1/counter.go":      counterDotGo(),
		"internal/metrics/v1/counter_test.go": counterTestDotGo(),
		// "internal/metrics/v1/runtime.go":      runtimeTestDotGo(),
		"internal/metrics/v1/runtime_test.go": runtimeTestDotGo(),
		"internal/metrics/v1/types.go":        typesDotGo(),
		"internal/metrics/v1/wire.go":         wireDotGo(),
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
