package utils

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	a = "assert"
	r = "require"
	t = "t"
	T = "T"
)

// WriteHeader calls res.WriteHeader for a given status
func WriteHeader(status string) jen.Code {
	return WriteXHeader("res", status)
}

// WriteXHeader calls WriteHeader for a given variable name
func WriteXHeader(varName, status string) jen.Code {
	return jen.ID(varName).Dot("WriteHeader").Call(
		jen.Qual("net/http", status),
	)
}

// ExpectMethod creates a test expectation for a gievn method
func ExpectMethod(varName, method string) jen.Code {
	return jen.ID(varName).Op(":=").Qual("net/http", method)
}

// ParallelTest creates a new t.Parallel call
func ParallelTest(tee *jen.Statement) jen.Code {
	if tee == nil {
		return jen.ID(T).Dot("Parallel").Call()
	}
	return tee.Dot("Parallel").Call()
}

func NilQueryFilter(proj *models.Project) jen.Code {
	return jen.Call(jen.PointerTo().Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter")).Call(jen.Nil())
}

func DefaultQueryFilter(proj *models.Project) jen.Code {
	return jen.Qual(filepath.Join(proj.OutputPath, "models/v1"), "DefaultQueryFilter").Call()
}

func CreateNilQueryFilter(proj *models.Project) jen.Code {
	return jen.ID(constants.FilterVarName).Op(":=").Call(jen.PointerTo().Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter")).Call(jen.Nil())
}

func CreateDefaultQueryFilter(proj *models.Project) jen.Code {
	return jen.ID(constants.FilterVarName).Op(":=").Qual(filepath.Join(proj.OutputPath, "models/v1"), "DefaultQueryFilter").Call()
}

func AppendItemsToList(list jen.Code, items ...jen.Code) jen.Code {
	return jen.Add(list).Equals().Append(append([]jen.Code{list}, items...)...)
}

func FakeError() jen.Code {
	return jen.Qual("errors", "New").Call(jen.Lit("blah"))
}

func buildSingleValueTestifyFunc(pkg, method string) func(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return func(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
		args := []jen.Code{
			jen.ID(t),
			value,
		}

		if message != nil {
			args = append(args, message)
		}
		for _, arg := range formatArgs {
			args = append(args, arg)
		}

		return jen.Qual(fmt.Sprintf("github.com/stretchr/testify/%s", pkg), method).Call(args...)
	}
}

func buildDoubleValueTestifyFunc(pkg, method string) func(expected, actual, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return func(first, second, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
		args := []jen.Code{
			jen.ID(t),
			first,
			second,
		}

		if message != nil {
			args = append(args, message)
		}
		for _, arg := range formatArgs {
			args = append(args, arg)
		}

		return jen.Qual(fmt.Sprintf("github.com/stretchr/testify/%s", pkg), method).Call(args...)
	}
}

// BuildTemplatePath builds a template path
func BuildTemplatePath(pkgRoot, tail string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", pkgRoot, tail)
}

// BuildSubTest builds a subtest
func BuildSubTest(name string, testInstructions ...jen.Code) jen.Code {
	return _buildSubtest(name, true, testInstructions...)
}

// BuildSubTestWithoutContext builds a subtest without context
func BuildSubTestWithoutContext(name string, testInstructions ...jen.Code) jen.Code {
	return _buildSubtest(name, false, testInstructions...)
}

func _buildSubtest(name string, includeContext bool, testInstructions ...jen.Code) jen.Code {
	insts := []jen.Code{}
	if includeContext {
		insts = append(insts, constants.CreateCtx(), jen.Line())
	}
	insts = append(insts, testInstructions...)

	return jen.ID(T).Dot("Run").Call(
		jen.Lit(name), jen.Func().Params(jen.ID(t).PointerTo().Qual("testing", T)).Block(insts...),
	)
}

// BuildTestServer builds a test server with an example handlerfunc
func BuildTestServer(name string, handlerLines ...jen.Code) *jen.Statement {
	return jen.ID(name).Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
		jen.Qual("net/http", "HandlerFunc").Callln(
			jen.Func().Params(
				jen.ID("res").Qual("net/http", "ResponseWriter"),
				jen.ID("req").PointerTo().Qual("net/http", "Request"),
			).Block(handlerLines...),
		),
	)
}

// OuterTestFunc does
func OuterTestFunc(subjectName string) *jen.Statement {
	return jen.Func().ID(fmt.Sprintf("Test%s", subjectName)).Params(jen.ID(T).PointerTo().Qual("testing", T))
}

// QueryFilterParam does
func QueryFilterParam(proj *models.Project) jen.Code {
	if proj != nil {
		return jen.ID(constants.FilterVarName).PointerTo().Qual(proj.ModelsV1Package(), "QueryFilter")
	}
	return jen.ID(constants.FilterVarName).PointerTo().ID("QueryFilter")
}

func FormatString(str string, args ...jen.Code) *jen.Statement {
	return jen.Qual("fmt", "Sprintf").Call(append([]jen.Code{jen.Lit(str)}, args...)...)
}

func FormatStringWithArg(arg jen.Code, args ...jen.Code) *jen.Statement {
	return jen.Qual("fmt", "Sprintf").Call(append([]jen.Code{arg}, args...)...)
}

func BuildError(args ...jen.Code) jen.Code {
	return jen.Qual("errors", "New").Call(args...)
}

func Error(str string) jen.Code {
	return jen.Qual("errors", "New").Call(jen.Lit(str))
}

func Errorf(str string, args ...interface{}) jen.Code {
	return jen.Qual("errors", "New").Call(jen.Litf(str, args...))
}

const SpanVarName = "span"

func StartSpan(proj *models.Project, saveCtx bool, spanName string) jen.Code {
	return StartSpanWithVar(proj, saveCtx, jen.Lit(spanName))
}

func StartSpanWithVar(proj *models.Project, saveCtx bool, spanName jen.Code) jen.Code {
	/*
		ctx, span := trace.StartSpan(ctx, "UpdateItem")
		defer span.End()
	*/
	g := &jen.Group{}

	g.Add(
		jen.List(
			func() jen.Code {
				if saveCtx {
					return constants.CtxVar()
				}
				return jen.Underscore()
			}(),
			jen.ID(SpanVarName),
		).Op(":=").Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(
			constants.CtxVar(),
			spanName,
		),
		jen.Line(),
		jen.Defer().ID(SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.Line(),
	)

	return g
}

func StartSpanWithInlineCtx(proj *models.Project, saveCtx bool, spanName jen.Code) jen.Code {
	/*
		ctx, span := trace.StartSpan(context.Background(), "UpdateItem")
		defer span.End()
	*/
	g := &jen.Group{}

	g.Add(
		jen.List(
			func() jen.Code {
				if saveCtx {
					return constants.CtxVar()
				}
				return jen.ID("_")
			}(),
			jen.ID(SpanVarName),
		).Op(":=").Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(
			constants.InlineCtx(),
			spanName,
		),
		jen.Line(),
		jen.Defer().ID(SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.Line(),
	)

	return g
}

func AssertExpectationsFor(varNames ...string) jen.Code {
	callArgs := []jen.Code{
		jen.ID("t"),
	}

	for _, name := range varNames {
		callArgs = append(callArgs, jen.ID(name))
	}

	return jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(callArgs...)
}

// RunGoimportsForFile runs the `goimports` binary for a given filename
func RunGoimportsForFile(filename string) error {
	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return exec.Command(filepath.Join(hd, "bin/goimports"), "-w", filename).Run()
}

func determineGofmtPath() string {
	gofmtLocation, err := exec.Command("which", "gofmt").Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(gofmtLocation)
}

// RunGoFormatForFile does something
func RunGoFormatForFile(filename string) error {
	if runtime.GOOS == "linux" {
		return exec.Command("gofmt", "-s", "-w", filename).Run()
	} else if runtime.GOOS == "darwin" {
		return exec.Command("gofmt", "-s", "-w", filename).Run()
	} else {
		return errors.New("invalid platform")
	}
}

// RenderGoFile does
func RenderGoFile(proj *models.Project, path string, file *jen.File) error {
	fp := BuildTemplatePath(proj.OutputPath, path)

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if mkdirErr := os.MkdirAll(filepath.Dir(fp), os.ModePerm); mkdirErr != nil {
			log.Printf("error making directory: %v\n", mkdirErr)
			return err
		}

		var b bytes.Buffer
		if err := file.Render(&b); err != nil {
			return fmt.Errorf("error rendering file %q: %w", path, err)
		}

		if err := ioutil.WriteFile(fp, b.Bytes(), 0644); err != nil {
			return fmt.Errorf("error rendering file %q: %w", path, err)
		}

		if gie := RunGoimportsForFile(fp); gie != nil {
			return fmt.Errorf("error rendering file %q: %w", path, gie)
		}

		if ferr := FindAndFixImportBlock(proj.OutputPath, fp); ferr != nil {
			return fmt.Errorf("error sorting imports for file %q: %w", path, ferr)
		}

		if gfe := RunGoFormatForFile(fp); gfe != nil {
			return fmt.Errorf("error rendering file %q: %w", path, gfe)
		}
	} else {
		return err
	}

	return nil
}

func BuildFakeVarName(typName string) string {
	return fmt.Sprintf("example%s", typName)
}

func BuildFakeVar(proj *models.Project, typName string, args ...jen.Code) jen.Code {
	return BuildFakeVarWithCustomName(proj, BuildFakeVarName(typName), typName, args...)
}

func BuildFakeVarWithCustomName(proj *models.Project, varName, funcName string, args ...jen.Code) jen.Code {
	if !strings.HasPrefix(funcName, "BuildFake") {
		funcName = fmt.Sprintf("BuildFake%s", funcName)
	}
	return jen.ID(varName).Assign().Qual(proj.FakeModelsPackage(), funcName).Call(args...)
}
