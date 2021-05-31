package indexinitializer

import (
	"bytes"
	_ "embed"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

//go:embed main.gotpl
var mainTemplate string

func mainDotGoString(proj *models.Project) string {
	switchCases := buildSwitchCases(proj)

	var b bytes.Buffer
	if err := jen.Null().Add(switchCases...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	return models.RenderCodeFile(proj, mainTemplate, map[string]string{
		"switchCases": b.String(),
	})
}

func buildSwitchCases(proj *models.Project) []jen.Code {
	switchCases := []jen.Code{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		if typ.SearchEnabled {
			switchCases = append(switchCases,
				jen.Case(jen.Lit(typ.Name.RouteName())).Body(
					jen.ID("outputChan").Assign().Make(jen.Chan().Index().PointerTo().Qual(proj.TypesPackage(), typ.Name.Singular())),
					jen.If(
						jen.ID("queryErr").Assign().ID("dbClient").Dotf("GetAll%s", pn).Call(
							constants.CtxVar(),
							jen.ID("outputChan"),
							jen.Lit(1000),
						),
						jen.ID("queryErr").DoesNotEqual().Nil(),
					).Body(
						// this statement is goofy because it renders a format variable.
						jen.Qual("log", "Fatalf").Call(jen.Lit("error fetching "+typ.Name.PluralCommonName()+" from database: %v"), jen.Err()),
					),
					jen.Line(),
					jen.For().Body(
						jen.Select().Body(
							jen.Case(jen.ID(typ.Name.PluralUnexportedVarName()).Assign().ReceiveFromChannel().ID("outputChan")).Body(
								jen.For(jen.List(jen.Underscore(), jen.ID("x").Assign().Range().ID(typ.Name.PluralUnexportedVarName()))).Body(
									jen.If(
										jen.ID("searchIndexErr").Assign().ID("im").Dot("Index").Call(
											constants.CtxVar(),
											jen.ID("x").Dot("ID"),
											jen.ID("x"),
										),
										jen.ID("searchIndexErr").DoesNotEqual().Nil(),
									).Body(
										jen.ID(constants.LoggerVarName).
											Dot("WithValue").Call(jen.Lit("id"), jen.ID("x").Dot("ID")).
											Dot("Error").Call(jen.ID("searchIndexErr"), jen.Lit("error adding to search index")),
									),
								),
							),
							jen.Case(jen.ReceiveFromChannel().Qual("time", "After").Call(jen.ID("deadline"))).Body(
								jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("terminating")),
								jen.Return(),
							),
						),
					),
				),
			)
		}
	}

	return switchCases
}
