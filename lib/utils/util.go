package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
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

func BuildTemplatePath(pkgRoot, tail string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", pkgRoot, tail)
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
)

func RunGoimportsForFile(filename string) error {
	return exec.Command("/home/jeffrey/bin/goimports", "-w", filename).Run()
}

func RunGoFormatForFile(filename string) error {
	return exec.Command("/usr/local/go/bin/gofmt", "-s", "-w", filename).Run()
}

func RenderGoFile(pkgRoot, path string, file *jen.File) error {
	// start := time.Now()
	fp := BuildTemplatePath(pkgRoot, path)

	_ = os.Remove(fp)
	if mkdirErr := os.MkdirAll(filepath.Dir(fp), os.ModePerm); mkdirErr != nil {
		log.Printf("error making directory: %v\n", mkdirErr)
	}

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

func ExampleValueForField(field models.DataField) jen.Code {
	switch field.Type {
	case "string":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("example"))
		}
		return jen.Lit("example")
	case "float32":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1.23))
		}
		return jen.Lit(1.23)
	case "float64":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1.23))
		}
		return jen.Lit(1.23)
	case "uint8":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "uint16":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "uint32":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "uint64":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "int8":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "int16":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "int32":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	case "int64":
		if field.Pointer {
			jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("1"))
		}
		return jen.Lit("1")
	default:
		return nil
	}
}

const FakeLibrary = "github.com/brianvoe/gofakeit"

func FakeCallForField(field models.DataField) jen.Code {
	switch field.Type {
	case "string":
		x := jen.Qual(FakeLibrary, "Word").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "float32":
		x := jen.Qual(FakeLibrary, "Float32").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "float64":
		x := jen.Qual(FakeLibrary, "Float64").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint8":
		x := jen.Qual(FakeLibrary, "Uint8").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint16":
		x := jen.Qual(FakeLibrary, "Uint16").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint32":
		x := jen.Qual(FakeLibrary, "Uint32").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint64":
		x := jen.Qual(FakeLibrary, "Uint64").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint":
		x := jen.Qual(FakeLibrary, "Uint64").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int8":
		x := jen.Qual(FakeLibrary, "Int8").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int16":
		x := jen.Qual(FakeLibrary, "Int16").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int32":
		x := jen.Qual(FakeLibrary, "Int32").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int64":
		x := jen.Qual(FakeLibrary, "Int64").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int":
		x := jen.Qual(FakeLibrary, "Int").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	default:
		return jen.Null()
	}
}
