package constants

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

const (
	// UserOwnershipFieldName represents the allowed field name for representing ownership by a user
	UserOwnershipFieldName = "BelongsToUser"

	ContextVarName = "ctx"

	FilterVarName = "filter"

	LoggerVarName = "logger"

	SpanVarName = "span"

	RequestVarName = "req"

	ResponseVarName = "res"
)

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

// UserIDVar is a shorthand for a context param
func UserIDVar() *jen.Statement {
	return jen.ID("userID")
}

// UserIDParam is a shorthand for a context param
func UserIDParam() jen.Code {
	return UserIDVar().Uint64()
}

// CtxParam is a shorthand for a context param
func CtxVar() *jen.Statement {
	return jen.ID(ContextVarName)
}

func err(str string) jen.Code {
	return jen.Qual("errors", "New").Call(jen.Lit(str))
}

func ObligatoryError() jen.Code {
	return err("blah")
}
