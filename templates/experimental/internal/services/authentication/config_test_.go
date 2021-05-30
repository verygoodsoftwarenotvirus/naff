package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestConfig_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("localKey"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/random", "GenerateRawBytes").Call(
						jen.ID("ctx"),
						jen.ID("pasetoKeyRequiredLength"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("PASETO").Op(":").ID("PASETOConfig").Valuesln(jen.ID("Issuer").Op(":").Lit("issuer"), jen.ID("LocalModeKey").Op(":").ID("localKey")), jen.ID("Cookies").Op(":").ID("CookieConfig").Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("Domain").Op(":").Lit("domain"), jen.ID("Lifetime").Op(":").Qual("time", "Second")), jen.ID("Debug").Op(":").ID("false"), jen.ID("EnableUserSignup").Op(":").ID("false"), jen.ID("MinimumUsernameLength").Op(":").Lit(123), jen.ID("MinimumPasswordLength").Op(":").Lit(123)),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestCookieConfig_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("CookieConfig").Valuesln(jen.ID("Name").Op(":").Lit("name"), jen.ID("Domain").Op(":").Lit("domain"), jen.ID("Lifetime").Op(":").Qual("time", "Second")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestPASETOConfig_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("localKey"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/random", "GenerateRawBytes").Call(
						jen.ID("ctx"),
						jen.ID("pasetoKeyRequiredLength"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("cfg").Op(":=").Op("&").ID("PASETOConfig").Valuesln(jen.ID("Issuer").Op(":").Lit("issuer"), jen.ID("LocalModeKey").Op(":").ID("localKey")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
