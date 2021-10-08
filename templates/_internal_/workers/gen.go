package workers

import (
	_ "embed"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	basePackagePath = "internal/workers"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]string{
		"data_changes_worker.go": dataChangesWorkerDotGo(proj),
		"data_changes_worker_test.go": dataChangesWorkerTestDotGo(proj),
		"pre_archives_worker.go": preArchivesWorkerDotGo(proj),
		"pre_archives_worker_test.go": preArchivesWorkerTestDotGo(proj),
		"pre_updates_worker.go": preUpdatesWorkerDotGo(proj),
		"pre_updates_worker_test.go": preUpdatesWorkerTestDotGo(proj),
		"pre_writes_worker.go": preWritesWorkerDotGo(proj),
		"pre_writes_worker_test.go": preWritesWorkerTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}



//go:embed data_changes_worker.gotpl
var dataChangesWorkerTemplate string
func dataChangesWorkerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, dataChangesWorkerTemplate, nil)
}



//go:embed data_changes_worker_test.gotpl
var dataChangesWorkerTestTemplate string
func dataChangesWorkerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, dataChangesWorkerTestTemplate, nil)
}



//go:embed pre_archives_worker.gotpl
var preArchivesWorkerTemplate string
func preArchivesWorkerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, preArchivesWorkerTemplate, nil)
}




//go:embed pre_archives_worker_test.gotpl
var preArchivesWorkerTestTemplate string
func preArchivesWorkerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, preArchivesWorkerTestTemplate, nil)
}



//go:embed pre_updates_worker.gotpl
var preUpdatesWorkerTemplate string
func preUpdatesWorkerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, preUpdatesWorkerTemplate, nil)
}



//go:embed pre_updates_worker_test.gotpl
var preUpdatesWorkerTestTemplate string
func preUpdatesWorkerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, preUpdatesWorkerTestTemplate, nil)
}



//go:embed pre_writes_worker.gotpl
var preWritesWorkerTemplate string
func preWritesWorkerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, preWritesWorkerTemplate, nil)
}



//go:embed pre_writes_worker_test.gotpl
var preWritesWorkerTestTemplate string
func preWritesWorkerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, preWritesWorkerTestTemplate, nil)
}