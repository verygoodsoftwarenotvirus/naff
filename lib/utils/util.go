package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	a = "assert"
	r = "require"
	t = "t"
	T = "T"
)

// Comments turns a bunch of lines of strings into jen comments
func Comments(input ...string) []jen.Code {
	out := []jen.Code{}
	for i, c := range input {
		if i == len(input)-1 {
			out = append(out, jen.Comment(c))
		} else {
			out = append(out, jen.Comment(c), jen.Line())
		}
	}
	return out
}

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
	return jen.Call(jen.Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter")).Call(jen.Nil())
}

func DefaultQueryFilter(proj *models.Project) jen.Code {
	return jen.Qual(filepath.Join(proj.OutputPath, "models/v1"), "DefaultQueryFilter").Call()
}

const FilterVarName = "filter"

func CreateNilQueryFilter(proj *models.Project) jen.Code {
	return jen.ID(FilterVarName).Op(":=").Call(jen.Op("*").Qual(filepath.Join(proj.OutputPath, "models/v1"), "QueryFilter")).Call(jen.Nil())
}

func CreateDefaultQueryFilter(proj *models.Project) jen.Code {
	return jen.ID(FilterVarName).Op(":=").Qual(filepath.Join(proj.OutputPath, "models/v1"), "DefaultQueryFilter").Call()
}

// FakeSeedFunc builds a consistent fake library seed init function
func FakeSeedFunc() jen.Code {
	return jen.Func().ID("init").Params().Block(
		InlineFakeSeedFunc(),
	)
}

// InlineFakeSeedFunc builds a consistent fake library seed init function
func InlineFakeSeedFunc() jen.Code {
	return jen.Qual(FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call())
}

func FakeFuncForType(typ string) func() jen.Code {
	switch typ {
	case "string":
		return FakeStringFunc
	default:
		panic(fmt.Sprintf("unknown type! %q", typ))
	}
}

func FakeStringFunc() jen.Code {
	return jen.Qual(FakeLibrary, "Word").Call()
}

func FakeContentTypeFunc() jen.Code {
	return jen.Qual(FakeLibrary, "MimeType").Call()
}

func FakeUUIDFunc() jen.Code {
	return jen.Qual(FakeLibrary, "UUID").Call()
}

func FakeURLFunc() jen.Code {
	return jen.Qual(FakeLibrary, "URL").Call()
}

func FakeHTTPMethodFunc() jen.Code {
	return jen.Qual(FakeLibrary, "HTTPMethod").Call()
}

func FakeUint32Func() jen.Code {
	return jen.Qual(FakeLibrary, "Uint32").Call()
}

func FakeUint64Func() jen.Code {
	return jen.Qual(FakeLibrary, "Uint64").Call()
}

func FakeUsernameFunc() jen.Code {
	return jen.Qual(FakeLibrary, "Username").Call()
}

func FakeUnixTimeFunc() jen.Code {
	return jen.Uint64().Call(jen.Qual(FakeLibrary, "Date").Call().Dot("Unix").Call())
}

func FakePasswordFunc() jen.Code {
	return jen.Qual(FakeLibrary, "Password").Call(
		jen.True(),
		jen.True(),
		jen.True(),
		jen.True(),
		jen.True(),
		jen.Lit(32),
	)
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
		insts = append(insts, CreateCtx())
	}
	insts = append(insts, testInstructions...)

	return jen.ID(T).Dot("Run").Call(
		jen.Lit(name), jen.Func().Params(jen.ID(t).Op("*").Qual("testing", T)).Block(insts...),
	)
}

// BuildTestServer builds a test server with an example handlerfunc
func BuildTestServer(name string, handlerLines ...jen.Code) *jen.Statement {
	return jen.ID(name).Op(":=").Qual("net/http/httptest", "NewTLSServer").Callln(
		jen.Qual("net/http", "HandlerFunc").Callln(
			jen.Func().Params(
				jen.ID("res").Qual("net/http", "ResponseWriter"),
				jen.ID("req").Op("*").Qual("net/http", "Request"),
			).Block(handlerLines...),
		),
	)
}

const ContextVarName = "ctx"

// CreateCtx calls context.Background() and assigns it to a variable called ctx
func CreateCtx() jen.Code {
	return CtxVar().Op(":=").Qual("context", "Background").Call()
}

// InlineCtx calls context.Background() and assigns it to a variable called ctx
func InlineCtx() jen.Code {
	return jen.Qual("context", "Background").Call()
}

// CtxParam is a shorthand for a context param
func CtxParam() jen.Code {
	return CtxVar().Qual("context", "Context")
}

// CtxParam is a shorthand for a context param
func CtxVar() *jen.Statement {
	return jen.ID(ContextVarName)
}

// OuterTestFunc does
func OuterTestFunc(subjectName string) *jen.Statement {
	return jen.Func().ID(fmt.Sprintf("Test%s", subjectName)).Params(jen.ID(T).Op("*").Qual("testing", T))
}

// QueryFilterParam does
func QueryFilterParam() jen.Code {
	return jen.ID(FilterVarName).PointerTo().ID("QueryFilter")
}

const SpanVarName = "span"

func StartSpan(proj *models.Project, saveCtx bool, spanName string) jen.Code {
	/*
		ctx, span := trace.StartSpan(ctx, "UpdateItem")
		defer span.End()
	*/
	g := &jen.Group{}

	g.Add(
		jen.List(func() jen.Code {
			if saveCtx {
				return CtxVar()
			}
			return jen.ID("_")
		}(), jen.ID(SpanVarName)).Op(":=").Qual(filepath.Join(proj.OutputPath, "internal", "v1", "tracing"), "StartSpan").Call(CtxVar(), jen.Lit(spanName)),
		jen.Line(),
		jen.Defer().ID(SpanVarName).Dot("End").Call(),
		jen.Line(),
		jen.Line(),
	)

	return g
}

// RunGoimportsForFile runs the `goimports` binary for a given filename
func RunGoimportsForFile(filename string) error {
	hd, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return exec.Command(filepath.Join(hd, "bin/goimports"), "-w", filename).Run()
}

// RunGoFormatForFile does
func RunGoFormatForFile(filename string) error {
	if runtime.GOOS == "linux" {
		return exec.Command("/usr/local/go/bin/gofmt", "-s", "-w", filename).Run()
	} else if runtime.GOOS == "darwin" {
		return exec.Command("/usr/local/bin/gofmt", "-s", "-w", filename).Run()
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
