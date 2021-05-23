package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildAuthConstantDefinitions()...)
	code.Add(buildAuthInit()...)
	code.Add(buildAuthTypeDefinitions()...)
	code.Add(buildAuthSessionInfoToBytes()...)

	return code
}

func buildAuthConstantDefinitions() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("SessionInfoKey is the non-string type we use for referencing SessionInfo structs"),
			jen.ID("SessionInfoKey").ID("ContextKey").Equals().Lit("session_info"),
		),
		jen.Line(),
	}

	return lines
}

func buildAuthInit() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.Qual("encoding/gob", "Register").Call(jen.AddressOf().ID("SessionInfo").Values()),
		),
		jen.Line(),
	}

	return lines
}

func buildAuthTypeDefinitions() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("SessionInfo represents what we encode in our authentication cookies."),
			jen.ID("SessionInfo").Struct(
				jen.ID(constants.UserIDFieldName).Uint64().Tag(jsonTag("-")),
				jen.ID("UserIsAdmin").Bool().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("StatusResponse is what we encode when the frontend wants to check auth status"),
			jen.ID("StatusResponse").Struct(
				jen.ID("Authenticated").Bool().Tag(jsonTag("isAuthenticated")),
				jen.ID("IsAdmin").Bool().Tag(jsonTag("isAdmin")),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildAuthSessionInfoToBytes() []jen.Code {
	lines := []jen.Code{
		jen.Comment("ToBytes returns the gob encoded session info"),
		jen.Line(),
		jen.Func().Params(jen.ID("i").PointerTo().ID("SessionInfo")).ID("ToBytes").Params().Params(jen.Index().Byte()).Body(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Line(),
			jen.If(
				jen.Err().Assign().Qual("encoding/gob", "NewEncoder").Call(jen.AddressOf().ID("b")).Dot("Encode").Call(jen.ID("i")),
				jen.Err().DoesNotEqual().Nil(),
			).Body(
				jen.Panic(jen.Err()),
			),
			jen.Line(),
			jen.Return(jen.ID("b").Dot("Bytes").Call()),
		),
		jen.Line(),
	}

	return lines
}
