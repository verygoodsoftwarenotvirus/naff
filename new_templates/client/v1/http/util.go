package client

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

func comments(input ...string) []jen.Code {
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

func writeHeader(status string) jen.Code {
	return jen.ID("res").Dot("WriteHeader").Call(
		jen.Qual("net/http", status),
	)
}

func expectMethod(varName, method string) jen.Code {
	return jen.ID(varName).Op(":=").Qual("net/http", method)
}

const (
	a = "assert"
	r = "require"
)

func parallelTest(tee *jen.Statement) jen.Code {
	if tee == nil {
		return jen.ID(T).Dot("Parallel").Call()
	}
	return tee.Dot("Parallel").Call()
}

func requireNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NoError")(value, message, formatArgs...)
}

func requireNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(r, "NotNil")(value, message, formatArgs...)
}

func assertTrue(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "True")(value, message, formatArgs...)
}

func assertFalse(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "False")(value, message, formatArgs...)
}

func assertNotNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NotNil")(value, message, formatArgs...)
}

func assertNil(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Nil")(value, message, formatArgs...)
}

func assertError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "Error")(value, message, formatArgs...)
}

func assertNoError(value, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
	return buildSingleValueTestifyFunc(a, "NoError")(value, message, formatArgs...)
}

func assertEqual(expected, actual, message *jen.Statement, formatArgs ...*jen.Statement) jen.Code {
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

func buildSubTest(name string, testInstructions ...jen.Code) jen.Code {
	return jen.ID(T).Dot("Run").Call(
		jen.Lit(name), jen.Func().Params(jen.ID(t).Op("*").Qual("testing", "T")).Block(testInstructions...),
	)
}

func createCtx() jen.Code {
	return jen.ID("ctx").Op(":=").Qual("context", "Background").Call()
}

func ctxParam() jen.Code {
	return jen.ID("ctx").Qual("context", "Context")
}

func testFunc(subjectName string) *jen.Statement {
	return jen.Func().ID(fmt.Sprintf("Test%s", subjectName)).Params(
		jen.ID(T).Op("*").Qual("testing", "T"),
	)
}
