package constants

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

const (
	CoreOAuth2Pkg          = "golang.org/x/oauth2"
	LoggingPkg             = "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	NoopLoggingPkg         = "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	AssertPkg              = "github.com/stretchr/testify/assert"
	MustAssertPkg          = "github.com/stretchr/testify/require"
	MockPkg                = "github.com/stretchr/testify/mock"
	DependencyInjectionPkg = "github.com/google/wire"
	FakeLibrary            = "github.com/brianvoe/gofakeit/v5"
	TracingLibrary         = "go.opencensus.io/trace"
	FlagParsingLibrary     = "github.com/spf13/pflag"
	SessionManagerLibrary  = "github.com/alexedwards/scs/v2"

	// UserIDVarName is what we normally call a user ID
	UserIDVarName = "userID"

	// UserIDVarName is what we normally call a user ID in a struct definition
	UserIDFieldName = "UserID"

	// UserOwnershipFieldName represents the allowed field name for representing ownership by a user
	UserOwnershipFieldName = "BelongsToUser"

	// ContextVarName is what we normally call a context.Context
	ContextVarName = "ctx"

	// FilterVarName is what we normally call a models.QueryFilter
	FilterVarName = "filter"

	// LoggerVarName is what we normally call a logging.Logger
	LoggerVarName = "logger"

	// SpanVarName is what we normally call a tracing span
	SpanVarName = "span"

	// RequestVarName is what we normally call an HTTP request
	RequestVarName = "req"

	// ResponseVarName is what we normally call an HTTP response
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

// LoggerParam is a shorthand for a context param
func LoggerParam() jen.Code {
	return jen.ID(LoggerVarName).Qual(LoggingPkg, "Logger")
}

// UserIDVar is a shorthand for a context param
func UserIDVar() *jen.Statement {
	return jen.ID(UserIDVarName)
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
