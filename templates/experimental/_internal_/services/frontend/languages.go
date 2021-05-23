package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func languagesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("displayLanguage").Op("*").ID("string").Type().ID("languageDetail").Struct(
			jen.ID("Tag").ID("language").Dot("Tag"),
			jen.ID("Name").ID("string"),
			jen.ID("Abbreviation").ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("english").Op("=").ID("new").Call(jen.ID("displayLanguage")),
			jen.ID("englishDetails").Op("=").Op("&").ID("languageDetail").Valuesln(
				jen.ID("Name").Op(":").Lit("English"), jen.ID("Abbreviation").Op(":").Lit("en-US"), jen.ID("Tag").Op(":").ID("language").Dot("AmericanEnglish")),
			jen.ID("spanish").Op("=").ID("new").Call(jen.ID("displayLanguage")),
			jen.ID("spanishDetails").Op("=").Op("&").ID("languageDetail").Valuesln(
				jen.ID("Name").Op(":").Lit("Spanish"), jen.ID("Abbreviation").Op(":").Lit("es-419"), jen.ID("Tag").Op(":").ID("language").Dot("LatinAmericanSpanish")),
			jen.ID("languageDetails").Op("=").Map(jen.Op("*").ID("displayLanguage")).Op("*").ID("languageDetail").Valuesln(
				jen.ID("english").Op(":").ID("englishDetails"), jen.ID("spanish").Op(":").ID("spanishDetails")),
			jen.ID("supportedLanguages").Op("=").Index().Op("*").ID("displayLanguage").Valuesln(
				jen.ID("english"), jen.ID("spanish")),
			jen.ID("defaultLanguage").Op("=").ID("english"),
			jen.ID("defaultLanguageDetails").Op("=").ID("languageDetails").Index(jen.ID("defaultLanguage")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("detailsForLanguage").Params(jen.ID("l").Op("*").ID("displayLanguage")).Params(jen.Op("*").ID("languageDetail")).Body(
			jen.Switch(jen.ID("l")).Body(
				jen.Case(jen.ID("spanish")).Body(
					jen.Return().ID("spanishDetails")),
				jen.Case(jen.ID("english")).Body(
					jen.Return().ID("englishDetails")),
				jen.Default().Body(
					jen.Return().ID("defaultLanguageDetails")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("determineLanguage").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("displayLanguage")).Body(
			jen.If(jen.ID("req").Op("==").ID("nil")).Body(
				jen.Return().ID("defaultLanguage")),
			jen.List(jen.ID("langs"), jen.ID("_"), jen.ID("err")).Op(":=").ID("language").Dot("ParseAcceptLanguage").Call(jen.ID("req").Dot("Header").Dot("Get").Call(jen.Lit("Accept-Language"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("defaultLanguage")),
			jen.If(jen.ID("len").Call(jen.ID("langs")).Op("!=").Lit(1)).Body(
				jen.Return().ID("defaultLanguage")),
			jen.Switch(jen.ID("langs").Index(jen.Lit(0))).Body(
				jen.Case(jen.ID("language").Dot("LatinAmericanSpanish"), jen.ID("language").Dot("EuropeanSpanish"), jen.ID("language").Dot("Spanish")).Body(
					jen.Return().ID("spanish")),
				jen.Case(jen.ID("language").Dot("AmericanEnglish"), jen.ID("language").Dot("BritishEnglish"), jen.ID("language").Dot("English")).Body(
					jen.Return().ID("english")),
				jen.Default().Body(
					jen.Return().ID("defaultLanguage")),
			),
		),
		jen.Line(),
	)

	return code
}
