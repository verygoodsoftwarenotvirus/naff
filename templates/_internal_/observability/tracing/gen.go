package tracing

import (
	_ "embed"
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
		"span_attachers.go":      spanAttachersDotGo(proj),
		"span_attachers_test.go": spanAttachersTestDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	stringFiles := map[string]string{
		"caller.go":                         callerDotGo(proj),
		"config_test.go":                    configTestDotGo(proj),
		"span_manager.go":                   spanManagerDotGo(proj),
		"spans_test.go":                     spansTestDotGo(proj),
		"caller_test.go":                    callerTestDotGo(proj),
		"instrumented_span_wrapper.go":      instrumentedSpanWrapperDotGo(proj),
		"instrumented_span_wrapper_test.go": instrumentedSpanWrapperTestDotGo(proj),
		"instrumentedsql.go":                instrumentedsqlDotGo(proj),
		"config.go":                         configDotGo(proj),
		"doc.go":                            docDotGo(proj),
		"spans.go":                          spansDotGo(proj),
		"instrumentedsql_test.go":           instrumentedsqlTestDotGo(proj),
		"span_manager_test.go":              spanManagerTestDotGo(proj),
		"tracer.go":                         tracerDotGo(proj),
		"tracer_test.go":                    tracerTestDotGo(proj),
		"transport.go":                      transportDotGo(proj),
		"transport_test.go":                 transportTestDotGo(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed caller.gotpl
var callerTemplate string

func callerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, callerTemplate, nil)
}

//go:embed config_test.gotpl
var configTestTemplate string

func configTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTestTemplate, nil)
}

//go:embed span_manager.gotpl
var spanManagerTemplate string

func spanManagerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, spanManagerTemplate, nil)
}

//go:embed spans_test.gotpl
var spansTestTemplate string

func spansTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, spansTestTemplate, nil)
}

//go:embed caller_test.gotpl
var callerTestTemplate string

func callerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, callerTestTemplate, nil)
}

//go:embed instrumented_span_wrapper.gotpl
var instrumentedSpanWrapperTemplate string

func instrumentedSpanWrapperDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, instrumentedSpanWrapperTemplate, nil)
}

//go:embed instrumented_span_wrapper_test.gotpl
var instrumentedSpanWrapperTestTemplate string

func instrumentedSpanWrapperTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, instrumentedSpanWrapperTestTemplate, nil)
}

//go:embed instrumentedsql.gotpl
var instrumentedsqlTemplate string

func instrumentedsqlDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, instrumentedsqlTemplate, nil)
}

//go:embed config.gotpl
var configTemplate string

func configDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, configTemplate, nil)
}

//go:embed doc.gotpl
var docTemplate string

func docDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, docTemplate, nil)
}

//go:embed spans.gotpl
var spansTemplate string

func spansDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, spansTemplate, nil)
}

//go:embed instrumentedsql_test.gotpl
var instrumentedsqlTestTemplate string

func instrumentedsqlTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, instrumentedsqlTestTemplate, nil)
}

//go:embed span_manager_test.gotpl
var spanManagerTestTemplate string

func spanManagerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, spanManagerTestTemplate, nil)
}

//go:embed tracer.gotpl
var tracerTemplate string

func tracerDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, tracerTemplate, nil)
}

//go:embed tracer_test.gotpl
var tracerTestTemplate string

func tracerTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, tracerTestTemplate, nil)
}

//go:embed transport.gotpl
var transportTemplate string

func transportDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, transportTemplate, nil)
}

//go:embed transport_test.gotpl
var transportTestTemplate string

func transportTestDotGo(proj *models.Project) string {
	return models.RenderCodeFile(proj, transportTestTemplate, nil)
}
