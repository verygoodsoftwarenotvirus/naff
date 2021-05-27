package secrets

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("buildExampleKey").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.ID("string")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("rawBytes"), jen.ID("err")).Op(":=").ID("random").Dot("GenerateRawBytes").Call(
				jen.ID("ctx"),
				jen.ID("expectedLocalKeyLength"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().Qual("encoding/base64", "URLEncoding").Dot("EncodeToString").Call(jen.ID("rawBytes")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestProvideSecretKeeper").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_local"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderLocal"), jen.ID("Key").Op(":").ID("buildExampleKey").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_aws"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderAWS"), jen.ID("Key").Op(":").ID("buildExampleKey").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard_vault"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Provider").Op(":").ID("ProviderHashicorpVault"), jen.ID("Key").Op(":").ID("buildExampleKey").Call(
						jen.ID("ctx"),
						jen.ID("t"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("ProvideSecretKeeper").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
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
