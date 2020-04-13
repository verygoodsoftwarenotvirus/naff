package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("Test_randString").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Comment("obligatory"),
			jen.Line(),
			jen.ID("actual").Assign().ID("randString").Call(),
			utils.AssertNotEmpty(jen.ID("actual"), nil),
			utils.AssertLength(jen.ID("actual"), jen.Lit(52), nil),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildConfig").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Comment("obligatory"),
			jen.Line(),
			jen.ID("actual").Assign().ID("BuildConfig").Call(),
			utils.AssertNotNil(jen.ID("actual"), nil),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestParseConfigFile").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID("tf"), jen.Err()).Assign().Qual("io/ioutil", "TempFile").Call(jen.Qual("os", "TempDir").Call(), jen.Lit("*.toml")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("expected").Assign().Lit("thisisatest"),
				jen.Line(),
				jen.List(jen.Underscore(), jen.Err()).Equals().ID("tf").Dot("Write").Call(
					jen.Index().Byte().Call(
						utils.FormatString(`
[server]
http_port = 1234
debug = false

[database]
provider = "postgres"
debug = true
connection_details = "%s"
`,
							jen.ID("expected"),
						),
					),
				),
				utils.RequireNoError(jen.Err(), nil),
				jen.Line(),
				jen.ID("expectedConfig").Assign().AddressOf().ID("ServerConfig").Valuesln(
					jen.ID("Server").MapAssign().ID("ServerSettings").Valuesln(
						jen.ID("HTTPPort").MapAssign().Lit(1234),
						jen.ID("Debug").MapAssign().False(),
					),
					jen.ID("Database").MapAssign().ID("DatabaseSettings").Valuesln(
						jen.ID("Provider").MapAssign().Lit("postgres"),
						jen.ID("Debug").MapAssign().True(),
						jen.ID("ConnectionDetails").MapAssign().Qual(proj.DatabaseV1Package(), "ConnectionDetails").Call(jen.ID("expected")),
					),
				),
				jen.Line(),
				jen.List(jen.ID("cfg"), jen.Err()).Assign().ID("ParseConfigFile").Call(jen.ID("tf").Dot("Name").Call()),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertEqual(jen.ID("expectedConfig").Dot("Server").Dot("HTTPPort"), jen.ID("cfg").Dot("Server").Dot("HTTPPort"), nil),
				utils.AssertEqual(jen.ID("expectedConfig").Dot("Server").Dot("Debug"), jen.ID("cfg").Dot("Server").Dot("Debug"), nil),
				utils.AssertEqual(jen.ID("expectedConfig").Dot("Database").Dot("Provider"), jen.ID("cfg").Dot("Database").Dot("Provider"), nil),
				utils.AssertEqual(jen.ID("expectedConfig").Dot("Database").Dot("Debug"), jen.ID("cfg").Dot("Database").Dot("Debug"), nil),
				utils.AssertEqual(jen.ID("expectedConfig").Dot("Database").Dot("ConnectionDetails"), jen.ID("cfg").Dot("Database").Dot("ConnectionDetails"), nil),
				jen.Line(),
				utils.AssertNoError(jen.Qual("os", "Remove").Call(jen.ID("tf").Dot("Name").Call()), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nonexistent file",
				jen.List(jen.ID("cfg"), jen.Err()).Assign().ID("ParseConfigFile").Call(jen.Lit("/this/doesn't/even/exist/lol")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("cfg"), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
