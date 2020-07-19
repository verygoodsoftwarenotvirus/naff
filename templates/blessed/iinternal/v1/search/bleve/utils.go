package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func utilsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(
		jen.Var().Defs(
			jen.ID("belongsToUserWithMandatedRestrictionRegexp").Equals().Qual("regexp", "MustCompile").Call(
				jen.RawString(`\+belongsToUser:\d+`),
			),
			jen.ID("belongsToUserWithoutMandatedRestrictionRegexp").Equals().Qual("regexp", "MustCompile").Call(
				jen.RawString(`belongsToUser:\d+`),
			),
		),
		jen.Line(),
		jen.Comment("ensureQueryIsRestrictedToUser takes a query and userID and ensures that query"),
		jen.Line(),
		jen.Comment("asks that results be restricted to a given user."),
		jen.Line(),
		jen.Func().ID("ensureQueryIsRestrictedToUser").Params(
			jen.ID("query").String(),
			constants.UserIDParam(),
		).String().Block(
			jen.Switch().Block(
				jen.Case(jen.ID("belongsToUserWithMandatedRestrictionRegexp").Dot("MatchString").Call(jen.ID("query"))).Block(
					jen.Return(jen.ID("query")),
				),
				jen.Case(jen.ID("belongsToUserWithoutMandatedRestrictionRegexp").Dot("MatchString").Call(jen.ID("query"))).Block(
					jen.ID("query").Equals().ID("belongsToUserWithoutMandatedRestrictionRegexp").Dot("ReplaceAllString").Call(
						jen.ID("query"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("+belongsToUser:%d"),
							constants.UserIDVar(),
						),
					),
				),
				jen.Case(jen.Not().ID("belongsToUserWithMandatedRestrictionRegexp").Dot("MatchString").Call(jen.ID("query"))).Block(
					jen.ID("query").Equals().Qual("fmt", "Sprintf").Call(
						jen.Lit("%s +belongsToUser:%d"),
						jen.ID("query"),
						constants.UserIDVar(),
					),
				),
			),
			jen.Line(),
			jen.Return(jen.ID("query")),
		),
	)

	return code
}
