package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func pasetoTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_setSignatureForRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleBody").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Qual("strings", "Repeat").Call(
						jen.Lit("A"),
						jen.ID("validClientSecretSize"),
					)),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.ID("expected").Op(":=").Lit("_l92KZfsYpDrPeP8CGTgHQiAtpEg3TyECry5Bd0ibdI"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("setSignatureForRequest").Call(
							jen.ID("req"),
							jen.ID("exampleBody"),
							jen.ID("exampleSecretKey"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("signatureHeaderKey")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildAPIClientAuthTokenRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePASETOCreationInput").Call(),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Qual("strings", "Repeat").Call(
						jen.Lit("A"),
						jen.ID("validClientSecretSize"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleSecretKey"),
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
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Qual("strings", "Repeat").Call(
						jen.Lit("A"),
						jen.ID("validClientSecretSize"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("nil"),
						jen.ID("exampleSecretKey"),
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
				jen.Lit("with invalid secret key"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePASETOCreationInput").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
						jen.ID("nil"),
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
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Qual("strings", "Repeat").Call(
						jen.Lit("A"),
						jen.ID("validClientSecretSize"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleSecretKey"),
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
				jen.Lit("with error building request"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("helper").Dot("builder").Op("=").ID("buildTestRequestBuilderWithInvalidURL").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePASETOCreationInput").Call(),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Qual("strings", "Repeat").Call(
						jen.Lit("A"),
						jen.ID("validClientSecretSize"),
					)),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleSecretKey"),
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
				jen.Lit("with error encoding input to buffer"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePASETOCreationInput").Call(),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Qual("strings", "Repeat").Call(
						jen.Lit("A"),
						jen.ID("validClientSecretSize"),
					)),
					jen.ID("clientEncoder").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/mock", "ClientEncoder").Valuesln(),
					jen.ID("clientEncoder").Dot("On").Call(
						jen.Lit("EncodeReader"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("exampleInput"),
					).Dot("Return").Call(
						jen.Qual("io", "Reader").Call(jen.Qual("bytes", "NewReader").Call(jen.Index().ID("byte").Call(jen.Lit("")))),
						jen.ID("nil"),
					),
					jen.ID("clientEncoder").Dot("On").Call(jen.Lit("ContentType")).Dot("Return").Call(jen.Lit("application/fart")),
					jen.ID("clientEncoder").Dot("On").Call(
						jen.Lit("Encode"),
						jen.ID("mock").Dot("Anything"),
						jen.ID("mock").Dot("IsType").Call(jen.Op("&").Qual("bytes", "Buffer").Valuesln()),
						jen.ID("exampleInput"),
					).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
					jen.ID("helper").Dot("builder").Dot("encoder").Op("=").ID("clientEncoder"),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleSecretKey"),
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
				jen.Lit("with error setting signature"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakePASETOCreationInput").Call(),
					jen.ID("exampleSecretKey").Op(":=").Index().ID("byte").Call(jen.Lit("A")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildAPIClientAuthTokenRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
						jen.ID("exampleSecretKey"),
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
