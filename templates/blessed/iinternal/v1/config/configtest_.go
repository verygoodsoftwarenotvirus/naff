package config

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("config")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("Test_randString").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Comment("obligatory"),
			jen.Line(),
			jen.ID("actual").Op(":=").ID("randString").Call(),
			jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("actual")),
			jen.Qual("github.com/stretchr/testify/assert", "Len").Call(jen.ID("t"), jen.ID("actual"), jen.Lit(52)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_buildConfig").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Comment("obligatory"),
			jen.Line(),
			jen.ID("actual").Op(":=").ID("BuildConfig").Call(),
			jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestParseConfigFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("tf"), jen.Err()).Op(":=").Qual("io/ioutil", "TempFile").Call(jen.Qual("os", "TempDir").Call(), jen.Lit("*.toml")),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.ID("expected").Op(":=").Lit("thisisatest"),
				jen.Line(),
				jen.List(jen.ID("_"), jen.Err()).Op("=").ID("tf").Dot("Write").Call(
					jen.Index().ID("byte").Call(
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(`
[server]
http_port = 1234
debug = false

[database]
provider = "postgres"
debug = true
connection_details = "%s"
`),
							jen.ID("expected"),
						),
					),
				),
				jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.ID("expectedConfig").Op(":=").Op("&").ID("ServerConfig").Valuesln(
					jen.ID("Server").Op(":").ID("ServerSettings").Valuesln(
						jen.ID("HTTPPort").Op(":").Lit(1234),
						jen.ID("Debug").Op(":").ID("false"),
					),
					jen.ID("Database").Op(":").ID("DatabaseSettings").Valuesln(
						jen.ID("Provider").Op(":").Lit("postgres"),
						jen.ID("Debug").Op(":").ID("true"),
						jen.ID("ConnectionDetails").Op(":").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "ConnectionDetails").Call(jen.ID("expected")),
					),
				),
				jen.Line(),
				jen.List(jen.ID("cfg"), jen.Err()).Op(":=").ID("ParseConfigFile").Call(jen.ID("tf").Dot("Name").Call()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedConfig").Dot("Server").Dot("HTTPPort"), jen.ID("cfg").Dot("Server").Dot("HTTPPort")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedConfig").Dot("Server").Dot("Debug"), jen.ID("cfg").Dot("Server").Dot("Debug")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedConfig").Dot("Database").Dot("Provider"), jen.ID("cfg").Dot("Database").Dot("Provider")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedConfig").Dot("Database").Dot("Debug"), jen.ID("cfg").Dot("Database").Dot("Debug")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expectedConfig").Dot("Database").Dot("ConnectionDetails"), jen.ID("cfg").Dot("Database").Dot("ConnectionDetails")),
				jen.Line(),
				jen.Qual("os", "Remove").Call(jen.ID("tf").Dot("Name").Call()),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nonexistent file"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.List(jen.ID("cfg"), jen.Err()).Op(":=").ID("ParseConfigFile").Call(jen.Lit("/this/doesn't/even/exist/lol")),
				jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.ID("cfg")),
			)),
		),
		jen.Line(),
	)
	return ret
}
