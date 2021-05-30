package secrets

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func secretsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("example").Struct(jen.ID("Name").ID("string")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestSecretKeeper").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("gocloud.dev/secrets", "Keeper")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("b"), jen.ID("err")).Op(":=").ID("random").Dot("GenerateRawBytes").Call(
				jen.ID("ctx"),
				jen.ID("expectedLocalKeyLength"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("b"),
			),
			jen.ID("key").Op(":=").Qual("encoding/base64", "URLEncoding").Dot("EncodeToString").Call(jen.ID("b")),
			jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderLocal"), jen.ID("Key").Op(":").ID("key")),
			jen.List(jen.ID("k"), jen.ID("err")).Op(":=").ID("ProvideSecretKeeper").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("k"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("k"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestSecretManager").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.ID("SecretManager")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("k").Op(":=").ID("buildTestSecretKeeper").Call(
				jen.ID("ctx"),
				jen.ID("t"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("k"),
			),
			jen.List(jen.ID("sm"), jen.ID("err")).Op(":=").ID("ProvideSecretManager").Call(
				jen.ID("logger"),
				jen.ID("k"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("sm"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("sm"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideSecretManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("k").Op(":=").ID("buildTestSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.List(jen.ID("sm"), jen.ID("err")).Op(":=").ID("ProvideSecretManager").Call(
						jen.ID("nil"),
						jen.ID("k"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil keeper"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("k"), jen.ID("err")).Op(":=").ID("ProvideSecretManager").Call(
						jen.ID("nil"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("k"),
					),
					jen.ID("require").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("broken").Struct(jen.ID("Thing").Qual("encoding/json", "Number")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_secretManager_Encrypt").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.List(jen.ID("rawKey"), jen.ID("err")).Op(":=").ID("localsecrets").Dot("NewRandomKey").Call(),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("key").Op(":=").Index().ID("byte").Valuesln(),
					jen.For(jen.List(jen.ID("i"), jen).Op(":=").Range().ID("rawKey")).Body(
						jen.ID("key").Op("=").ID("append").Call(
							jen.ID("key"),
							jen.ID("rawKey").Index(jen.ID("i")),
						)),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderLocal"), jen.ID("Key").Op(":").Qual("encoding/base64", "URLEncoding").Dot("EncodeToString").Call(jen.ID("key"))),
					jen.List(jen.ID("k"), jen.ID("err")).Op(":=").ID("ProvideSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("k"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("sm"), jen.ID("err")).Op(":=").ID("ProvideSecretManager").Call(
						jen.ID("logger"),
						jen.ID("k"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("sm").Dot("Encrypt").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid value"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.List(jen.ID("rawKey"), jen.ID("err")).Op(":=").ID("localsecrets").Dot("NewRandomKey").Call(),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("key").Op(":=").Index().ID("byte").Valuesln(),
					jen.For(jen.List(jen.ID("i"), jen).Op(":=").Range().ID("rawKey")).Body(
						jen.ID("key").Op("=").ID("append").Call(
							jen.ID("key"),
							jen.ID("rawKey").Index(jen.ID("i")),
						)),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderLocal"), jen.ID("Key").Op(":").Qual("encoding/base64", "URLEncoding").Dot("EncodeToString").Call(jen.ID("key"))),
					jen.List(jen.ID("k"), jen.ID("err")).Op(":=").ID("ProvideSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("k"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("sm"), jen.ID("err")).Op(":=").ID("ProvideSecretManager").Call(
						jen.ID("logger"),
						jen.ID("k"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("sm"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("exampleInput").Op(":=").Op("&").ID("broken").Valuesln(jen.ID("Thing").Op(":").Qual("encoding/json", "Number").Call(jen.ID("t").Dot("Name").Call())),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("sm").Dot("Encrypt").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_secretManager_Decrypt").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("sm").Op(":=").ID("buildTestSecretManager").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("encrypted"), jen.ID("err")).Op(":=").ID("sm").Dot("Encrypt").Call(
						jen.ID("ctx"),
						jen.ID("expected"),
					),
					jen.ID("require").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("encrypted"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Var().Defs(
						jen.ID("actual").Op("*").ID("example"),
					),
					jen.ID("err").Op("=").ID("sm").Dot("Decrypt").Call(
						jen.ID("ctx"),
						jen.ID("encrypted"),
						jen.Op("&").ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid value"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("sm").Op(":=").ID("buildTestSecretManager").Call(jen.ID("t")),
					jen.Var().Defs(
						jen.ID("actual").Op("*").ID("example"),
					),
					jen.ID("err").Op(":=").ID("sm").Dot("Decrypt").Call(
						jen.ID("ctx"),
						jen.Lit(" this isn't a real string lol "),
						jen.Op("&").ID("actual"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with inability to decrypt"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("sm").Op(":=").ID("buildTestSecretManager").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("encrypted"), jen.ID("err")).Op(":=").ID("sm").Dot("Encrypt").Call(
						jen.ID("ctx"),
						jen.ID("expected"),
					),
					jen.ID("require").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("encrypted"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("sm").Assert(jen.Op("*").ID("secretManager")).Dot("keeper").Op("=").ID("buildTestSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					),
					jen.Var().Defs(
						jen.ID("actual").Op("*").ID("example"),
					),
					jen.ID("err").Op("=").ID("sm").Dot("Decrypt").Call(
						jen.ID("ctx"),
						jen.ID("encrypted"),
						jen.Op("&").ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid JSON value"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("sm").Op(":=").ID("buildTestSecretManager").Call(jen.ID("t")),
					jen.List(jen.ID("encrypted"), jen.ID("err")).Op(":=").ID("sm").Assert(jen.Op("*").ID("secretManager")).Dot("keeper").Dot("Encrypt").Call(
						jen.ID("ctx"),
						jen.Index().ID("byte").Call(jen.Lit(` this isn't a real JSON string lol `)),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("encoded").Op(":=").Qual("encoding/base64", "URLEncoding").Dot("EncodeToString").Call(jen.ID("encrypted")),
					jen.Var().Defs(
						jen.ID("actual").Op("*").ID("example"),
					),
					jen.ID("err").Op("=").ID("sm").Dot("Decrypt").Call(
						jen.ID("ctx"),
						jen.ID("encoded"),
						jen.Op("&").ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
