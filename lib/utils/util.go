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

// RequireNoError creates a require call
func RequireNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NoError")(value, message, formatArgs...)
}

// RequireNotNil creates a require call
func RequireNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NotNil")(value, message, formatArgs...)
}

// RequireNil creates a require call
func RequireNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "Nil")(value, message, formatArgs...)
}

// AssertTrue calls assert.True
func AssertTrue(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "True")(value, message, formatArgs...)
}

// AssertFalse calls assert.False
func AssertFalse(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "False")(value, message, formatArgs...)
}

// AssertNotNil calls assert.NotNil
func AssertNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotNil")(value, message, formatArgs...)
}

// AssertNil calls assert.Nil
func AssertNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Nil")(value, message, formatArgs...)
}

// AssertError calls assert.Error
func AssertError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Error")(value, message, formatArgs...)
}

// AssertNoError calls assert.NoError
func AssertNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NoError")(value, message, formatArgs...)
}

// AssertNotEmpty calls assert.NotEmpty
func AssertNotEmpty(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotEmpty")(value, message, formatArgs...)
}

// AssertEqual calls assert.Equal
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

// CreateCtx calls context.Background() and assigns it to a variable called ctx
func CreateCtx() jen.Code {
	return jen.ID("ctx").Op(":=").Qual("context", "Background").Call()
}

// CtxParam is a shorthand for a context param
func CtxParam() jen.Code {
	return jen.ID("ctx").Qual("context", "Context")
}

// OuterTestFunc does
func OuterTestFunc(subjectName string) *jen.Statement {
	return jen.Func().ID(fmt.Sprintf("Test%s", subjectName)).Params(
		jen.ID(T).Op("*").Qual("testing", T),
	)
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
func RenderGoFile(pkgRoot, path string, file *jen.File) error {
	fp := BuildTemplatePath(pkgRoot, path)
	log.Printf("rendering %q\n", fp)

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

		if ferr := FindAndFixImportBlock(pkgRoot, fp); ferr != nil {
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
