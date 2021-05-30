package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func utilsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("belongsToAccountWithMandatedRestrictionRegexp").Op("=").Qual("regexp", "MustCompile").Call(jen.Lit(`\+belongsToAccount:\d+`)),
			jen.ID("belongsToAccountWithoutMandatedRestrictionRegexp").Op("=").Qual("regexp", "MustCompile").Call(jen.Lit(`belongsToAccount:\d+`)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ensureQueryIsRestrictedToUser takes a query and userID and ensures that query"),
		jen.Line(),
		jen.Func().Comment("asks that results be restricted to a given user.").ID("ensureQueryIsRestrictedToUser").Params(jen.ID("query").ID("string"), jen.ID("userID").ID("uint64")).Params(jen.ID("string")).Body(
			jen.Switch().Body(
				jen.Case(jen.ID("belongsToAccountWithMandatedRestrictionRegexp").Dot("MatchString").Call(jen.ID("query"))).Body(
					jen.Return().ID("query")),
				jen.Case(jen.ID("belongsToAccountWithoutMandatedRestrictionRegexp").Dot("MatchString").Call(jen.ID("query"))).Body(
					jen.ID("query").Op("=").ID("belongsToAccountWithoutMandatedRestrictionRegexp").Dot("ReplaceAllString").Call(
						jen.ID("query"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("+belongsToAccount:%d"),
							jen.ID("userID"),
						),
					)),
				jen.Case(jen.Op("!").ID("belongsToAccountWithMandatedRestrictionRegexp").Dot("MatchString").Call(jen.ID("query"))).Body(
					jen.ID("query").Op("=").Qual("fmt", "Sprintf").Call(
						jen.Lit("%s +belongsToAccount:%d"),
						jen.ID("query"),
						jen.ID("userID"),
					)),
			),
			jen.Return().ID("query"),
		),
		jen.Line(),
	)

	return code
}
