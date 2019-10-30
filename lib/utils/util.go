package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

const (
	a = "assert"
	r = "require"
	T = "T"
	t = "t"
)

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

func WriteHeader(status string) jen.Code {
	return jen.ID("res").Dot("WriteHeader").Call(
		jen.Qual("net/http", status),
	)
}

func ExpectMethod(varName, method string) jen.Code {
	return jen.ID(varName).Op(":=").Qual("net/http", method)
}

func ParallelTest(tee *jen.Statement) jen.Code {
	if tee == nil {
		return jen.ID(T).Dot("Parallel").Call()
	}
	return tee.Dot("Parallel").Call()
}

func RequireNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NoError")(value, message, formatArgs...)
}

func RequireNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NotNil")(value, message, formatArgs...)
}

func RequireNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "Nil")(value, message, formatArgs...)
}

func AssertTrue(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "True")(value, message, formatArgs...)
}

func AssertFalse(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "False")(value, message, formatArgs...)
}

func AssertNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotNil")(value, message, formatArgs...)
}

func AssertNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Nil")(value, message, formatArgs...)
}

func AssertError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Error")(value, message, formatArgs...)
}

func AssertNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NoError")(value, message, formatArgs...)
}

func AssertNotEmpty(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotEmpty")(value, message, formatArgs...)
}

func AssertEqual(expected, actual, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildDoubleValueTestifyFunc(a, "Equal")(expected, actual, message, formatArgs...)
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

func BuildTemplatePath(tail string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/naff/example_output", tail)
}

func BuildSubTest(name string, testInstructions ...jen.Code) jen.Code {
	return _buildSubtest(name, true, testInstructions...)
}

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

func CreateCtx() jen.Code {
	return jen.ID("ctx").Op(":=").Qual("context", "Background").Call()
}

func CtxParam() jen.Code {
	return jen.ID("ctx").Qual("context", "Context")
}

func OuterTestFunc(subjectName string) *jen.Statement {
	return jen.Func().ID(fmt.Sprintf("Test%s", subjectName)).Params(
		jen.ID(T).Op("*").Qual("testing", T),
	)
}

const (
	CoreOAuth2Pkg  = "golang.org/x/oauth2"
	LoggingPkg     = "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	NoopLoggingPkg = "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	AssertPkg      = "github.com/stretchr/testify/assert"
	MustAssertPkg  = "github.com/stretchr/testify/require"
	MockPkg        = "github.com/stretchr/testify/mock"
	ModelsPkg      = "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

func RunGoimportsForFile(filename string) error {
	return exec.Command("/home/jeffrey/bin/goimports", "-w", filename).Run()
}

func RunGoFormatForFile(filename string) error {
	return exec.Command("/usr/local/go/bin/gofmt", "-s", "-w", filename).Run()
}

func RenderFile(pkgRoot, path string, file *jen.File) error {
	// start := time.Now()
	fp := BuildTemplatePath(path)
	_ = os.Remove(fp)

	var b bytes.Buffer
	if err := file.Render(&b); err != nil {
		return fmt.Errorf("error rendering file %q: %w", path, err)
	}

	if err := ioutil.WriteFile(fp, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("error rendering file %q: %w", path, err)
	}

	gie := RunGoimportsForFile(fp)
	if gie != nil {
		return fmt.Errorf("error rendering file %q: %w", path, gie)
	}

	if ferr := FindAndFixImportBlock(pkgRoot, fp); ferr != nil {
		return fmt.Errorf("error sorting imports for file %q: %w", path, ferr)
	}

	if gfe := RunGoFormatForFile(fp); gfe != nil {
		return fmt.Errorf("error rendering file %q: %w", path, gfe)
	}
	// log.Printf("took %s to render %q", time.Since(start), path)

	return nil
}
