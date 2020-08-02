package mocksearch

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildInterfaceImplementationStatement(proj)...)
	code.Add(buildIndexManager()...)

	return code
}

func buildInterfaceImplementationStatement(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual(proj.InternalSearchV1Package(), "IndexManager").Equals().Parens(jen.PointerTo().ID("IndexManager")).Parens(jen.Nil()),
		jen.Line(),
	}

	return lines
}

func buildIndexManager() []jen.Code {
	lines := []jen.Code{
		jen.Comment("IndexManager is a mock IndexManager"),
		jen.Line(),
		jen.Type().ID("IndexManager").Struct(
			jen.Qual(constants.MockPkg, "Mock"),
		),
		jen.Line(),
		jen.Comment("Index implements our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("IndexManager")).ID("Index").Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
			jen.ID("value").Interface(),
		).Error().Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(
				constants.CtxVar(),
				jen.ID("id"),
				jen.ID("value"),
			),
			jen.Return(jen.ID("args").Dot("Error").Call(jen.Zero())),
		),
		jen.Line(),
		jen.Comment("Search implements our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("IndexManager")).ID("Search").Params(
			constants.CtxParam(),
			jen.ID("query").String(),
			constants.UserIDParam(),
		).Params(
			jen.ID("ids").Index().Uint64(),
			jen.Err().Error(),
		).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(
				constants.CtxVar(),
				jen.ID("query"),
				jen.ID(constants.UserIDVarName),
			),
			jen.Return(
				jen.ID("args").Dot("Get").Call(jen.Zero()).Dot("").Call(jen.Index().Uint64()),
				jen.ID("args").Dot("Error").Call(jen.One()),
			),
		),
		jen.Line(),
		jen.Comment("Delete implements our interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("IndexManager")).ID("Delete").Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Error().Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(
				constants.CtxVar(),
				jen.ID("id"),
			),
			jen.Return(jen.ID("args").Dot("Error").Call(jen.Zero())),
		),
		jen.Line(),
	}

	return lines
}
