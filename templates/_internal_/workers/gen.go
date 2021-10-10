package workers

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "workers"

	basePackagePath = "internal/workers"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"doc.go":                      docDotGo(proj),
		"pre_archives_worker.go":      preArchivesWorkerDotGo(proj),
		"pre_archives_worker_test.go": preArchivesWorkerTestDotGo(proj),
		"pre_updates_worker.go":       preUpdatesWorkerDotGo(proj),
		"pre_updates_worker_test.go":  preUpdatesWorkerTestDotGo(proj),
		"pre_writes_worker.go":        preWritesWorkerDotGo(proj),
		"pre_writes_worker_test.go":   preWritesWorkerTestDotGo(proj),
		"data_changes_worker.go":      dataChangesWorkerDotGo(proj),
		"data_changes_worker_test.go": dataChangesWorkerTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}
