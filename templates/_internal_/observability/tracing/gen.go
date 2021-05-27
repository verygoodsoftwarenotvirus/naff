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
		"tracer_test.go":                    tracerTestDotGo(proj),
		"caller.go":                         callerDotGo(proj),
		"config_test.go":                    configTestDotGo(proj),
		"span_manager.go":                   spanManagerDotGo(proj),
		"spans_test.go":                     spansTestDotGo(proj),
		"tracer.go":                         tracerDotGo(proj),
		"caller_test.go":                    callerTestDotGo(proj),
		"instrumented_span_wrapper.go":      instrumentedSpanWrapperDotGo(proj),
		"instrumented_span_wrapper_test.go": instrumentedSpanWrapperTestDotGo(proj),
		"instrumentedsql.go":                instrumentedsqlDotGo(proj),
		"config.go":                         configDotGo(proj),
		"doc.go":                            docDotGo(proj),
		"span_attachers.go":                 spanAttachersDotGo(proj),
		"spans.go":                          spansDotGo(proj),
		"instrumentedsql_test.go":           instrumentedsqlTestDotGo(proj),
		"span_attachers_test.go":            spanAttachersTestDotGo(proj),
		"span_manager_test.go":              spanManagerTestDotGo(proj),
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
