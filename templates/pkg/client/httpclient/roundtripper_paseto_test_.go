package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func roundtripperPasetoTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_newPASETORoundTripper").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleClientID").Op(":=").Lit("example_client_id"),
					jen.ID("exampleSecret").Op(":=").ID("make").Call(
						jen.Index().ID("byte"),
						jen.ID("validClientSecretSize"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("newPASETORoundTripper").Call(
							jen.ID("c"),
							jen.ID("exampleClientID"),
							jen.ID("exampleSecret"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
