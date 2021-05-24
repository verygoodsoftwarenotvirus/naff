package tracing

import (
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "tracing"

	basePackagePath = "internal/observability/tracing"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"config.go":                          configDotGo(proj),
		"config_test_.go":                    configTestDotGo(proj),
		"instrumented_span_wrapper.go":       instrumentedSpanWrapperDotGo(proj),
		"span_attachers.go":                  spanAttachersDotGo(proj),
		"span_manager.go":                    spanManagerDotGo(proj),
		"span_manager_test_.go":              spanManagerTestDotGo(proj),
		"spans.go":                           spansDotGo(proj),
		"tracer_test_.go":                    tracerTestDotGo(proj),
		"caller.go":                          callerDotGo(proj),
		"doc.go":                             docDotGo(proj),
		"instrumentedsql_test_.go":           instrumentedsqlTestDotGo(proj),
		"span_attachers_test_.go":            spanAttachersTestDotGo(proj),
		"caller_test_.go":                    callerTestDotGo(proj),
		"instrumentedsql.go":                 instrumentedsqlDotGo(proj),
		"tracer.go":                          tracerDotGo(proj),
		"instrumented_span_wrapper_test_.go": instrumentedSpanWrapperTestDotGo(proj),
		"spans_test_.go":                     spansTestDotGo(proj),
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
