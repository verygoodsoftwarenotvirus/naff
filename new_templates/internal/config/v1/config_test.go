package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func configTestDotGo() *jen.File {
	ret := jen.NewFile("config")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("Test_randString").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
		jen.ID("t").Dot(
			"Parallel",
		).Call(),
		jen.ID("actual").Op(":=").ID("randString").Call(),
		jen.ID("assert").Dot(
			"NotEmpty",
		).Call(jen.ID("t"), jen.ID("actual")),
		jen.ID("assert").Dot(
			"Len",
		).Call(jen.ID("t"), jen.ID("actual"), jen.Lit(52)),
	),
	)
	ret.Add(jen.Func().ID("Test_buildConfig").Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
		jen.ID("t").Dot(
			"Parallel",
		).Call(),
		jen.ID("actual").Op(":=").ID("BuildConfig").Call(),
		jen.ID("assert").Dot(
			"NotNil",
		).Call(jen.ID("t"), jen.ID("actual")),
	),
	)
	ret.Add(jen.Func().ID("TestParseConfigFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("tf"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(jen.Qual("os", "TempDir").Call(), jen.Lit("*.toml")),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("expected").Op(":=").Lit("thisisatest"),
			jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("tf").Dot(
				"Write",
			).Call(jen.Index().ID("byte").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit(`
[server]
http_port = 1234
debug = false

[database]
provider = "postgres"
debug = true
connection_details = "%s"
`), jen.ID("expected")))),
			jen.ID("require").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("expectedConfig").Op(":=").Op("&").ID("ServerConfig").Valuesln(jen.ID("Server").Op(":").ID("ServerSettings").Valuesln(jen.ID("HTTPPort").Op(":").Lit(1234), jen.ID("Debug").Op(":").ID("false")), jen.ID("Database").Op(":").ID("DatabaseSettings").Valuesln(jen.ID("Provider").Op(":").Lit("postgres"), jen.ID("Debug").Op(":").ID("true"), jen.ID("ConnectionDetails").Op(":").ID("database").Dot(
				"ConnectionDetails",
			).Call(jen.ID("expected")))),
			jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(jen.ID("tf").Dot(
				"Name",
			).Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedConfig").Dot(
				"Server",
			).Dot(
				"HTTPPort",
			), jen.ID("cfg").Dot(
				"Server",
			).Dot(
				"HTTPPort",
			)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedConfig").Dot(
				"Server",
			).Dot(
				"Debug",
			), jen.ID("cfg").Dot(
				"Server",
			).Dot(
				"Debug",
			)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedConfig").Dot(
				"Database",
			).Dot(
				"Provider",
			), jen.ID("cfg").Dot(
				"Database",
			).Dot(
				"Provider",
			)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedConfig").Dot(
				"Database",
			).Dot(
				"Debug",
			), jen.ID("cfg").Dot(
				"Database",
			).Dot(
				"Debug",
			)),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("expectedConfig").Dot(
				"Database",
			).Dot(
				"ConnectionDetails",
			), jen.ID("cfg").Dot(
				"Database",
			).Dot(
				"ConnectionDetails",
			)),
			jen.Qual("os", "Remove").Call(jen.ID("tf").Dot(
				"Name",
			).Call()),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with nonexistent file"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(jen.Lit("/this/doesn't/even/exist/lol")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("cfg")),
		)),
	),
	)
	return ret
}
