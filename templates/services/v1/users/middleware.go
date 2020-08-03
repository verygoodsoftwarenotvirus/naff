package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildMiddlewareConstantDefs(proj)...)
	code.Add(buildUserInputMiddleware(proj)...)
	code.Add(buildPasswordUpdateInputMiddleware(proj)...)
	code.Add(buildTOTPSecretVerificationInputMiddleware(proj)...)
	code.Add(buildTOTPSecretRefreshInputMiddleware(proj)...)

	return code
}

func buildMiddlewareConstantDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("userCreationMiddlewareCtxKey is the context key for creation input."),
			jen.ID("userCreationMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("user_creation_input"),
			jen.Line(),
			jen.Comment("passwordChangeMiddlewareCtxKey is the context key for password changes."),
			jen.ID("passwordChangeMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("user_password_change"),
			jen.Line(),
			jen.Comment("totpSecretVerificationMiddlewareCtxKey is the context key for TOTP token refreshes."),
			jen.ID("totpSecretVerificationMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("totp_verify"),
			jen.Line(),
			jen.Comment("totpSecretRefreshMiddlewareCtxKey is the context key for TOTP token refreshes."),
			jen.ID("totpSecretRefreshMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("totp_refresh"),
		),
		jen.Line(),
	}

	return lines
}

func buildUserInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UserInputMiddleware fetches user input from requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UserInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "UserCreationInput")),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UserInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.Comment("decode the request."),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("attach parsed value to request context."),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("userCreationMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildPasswordUpdateInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("PasswordUpdateInputMiddleware fetches password update input from requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("PasswordUpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "PasswordUpdateInput")),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("PasswordUpdateInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.Comment("decode the request."),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("attach parsed value to request context."),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("passwordChangeMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTOTPSecretVerificationInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("TOTPSecretVerificationInputMiddleware fetches 2FA update input from requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("TOTPSecretVerificationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "TOTPSecretVerificationInput")),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("TOTPSecretVerificationInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.Comment("decode the request."),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("attach parsed value to request context."),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("totpSecretVerificationMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildTOTPSecretRefreshInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("TOTPSecretRefreshInputMiddleware fetches 2FA update input from requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("TOTPSecretRefreshInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "TOTPSecretRefreshInput")),
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("TOTPSecretRefreshInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.Comment("decode the request."),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("attach parsed value to request context."),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("totpSecretRefreshMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}
