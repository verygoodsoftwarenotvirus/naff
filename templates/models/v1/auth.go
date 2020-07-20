package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("models")

	utils.AddImports(proj, code)

	code.Add(
		jen.Const().Defs(
			jen.Comment("SessionInfoKey is the non-string type we use for referencing SessionInfo structs"),
			jen.ID("SessionInfoKey").ID("ContextKey").Equals().Lit("session_info"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Block(
			jen.Qual("encoding/gob", "Register").Call(jen.AddressOf().ID("SessionInfo").Values()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.Comment("SessionInfo represents what we encode in our authentication cookies."),
			jen.Line(),
			jen.ID("SessionInfo").Struct(
				jen.ID(constants.UserIDFieldName).Uint64().Tag(jsonTag("-")),
				jen.ID("UserIsAdmin").Bool().Tag(jsonTag("-")),
			),
			jen.Line(),
			jen.Comment("StatusResponse is what we encode when the frontend wants to check auth status"),
			jen.Line(),
			jen.ID("StatusResponse").Struct(
				jen.ID("Authenticated").Bool().Tag(jsonTag("isAuthenticated")),
				jen.ID("IsAdmin").Bool().Tag(jsonTag("isAdmin")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ToBytes returns the gob encoded session info"),
		jen.Line(),
		jen.Func().Params(jen.ID("i").PointerTo().ID("SessionInfo")).ID("ToBytes").Params().Params(jen.Index().Byte()).Block(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Line(),
			jen.If(
				jen.Err().Assign().Qual("encoding/gob", "NewEncoder").Call(jen.AddressOf().ID("b")).Dot("Encode").Call(jen.ID("i")),
				jen.Err().DoesNotEqual().Nil(),
			).Block(
				jen.Panic(jen.Err()),
			),
			jen.Line(),
			jen.Return(jen.ID("b").Dot("Bytes").Call()),
		),
		jen.Line(),
	)

	return code
}
