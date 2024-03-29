package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	firstSupportedDatabase := ""
	for _, db := range proj.EnabledDatabases() {
		firstSupportedDatabase = db
	}

	code.Add(
		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("cfg").Assign().Op("&").ID("Config").Valuesln(
						jen.ID("Provider").MapAssign().IDf("%sProvider", firstSupportedDatabase), jen.ID("ConnectionDetails").MapAssign().Lit("example_connection_string")),
					jen.Newline(),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestProvideSessionManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("cookieConfig").Assign().Qual(proj.AuthServicePackage(), "CookieConfig").Values(),
					jen.ID("store").Assign().Qual("github.com/alexedwards/scs/v2/memstore", "New").Call(),
					jen.Newline(),
					jen.ID("mdm").Assign().AddressOf().Qual(proj.DatabasePackage(), "MockDatabase").Values(),
					jen.ID("mdm").Dot("On").Call(jen.Lit("ProvideSessionStore")).Dot("Return").Call(jen.ID("store")),
					jen.Newline(),
					jen.List(jen.ID("sessionManager"), jen.ID("err")).Assign().ID("ProvideSessionManager").Call(
						jen.ID("cookieConfig"),
						jen.ID("mdm"),
					),
					jen.Qual(constants.AssertionLibrary, "NotNil").Call(
						jen.ID("t"),
						jen.ID("sessionManager"),
					),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
