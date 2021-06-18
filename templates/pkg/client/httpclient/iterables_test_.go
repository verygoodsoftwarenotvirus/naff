package httpclient

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	code.Add(
		jen.Func().IDf("Test%s", pn).Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.IDf("%sTestSuite", puvn)),
			),
		),
		jen.Newline(),
		jen.Type().IDf("%sBaseSuite", puvn).Struct(
			jen.ID("suite").Dot("Suite"),
			jen.Newline(),
			jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("example%s", sn).Op("*").ID("types").Dot(sn),
			jen.IDf("example%sList", sn).Op("*").ID("types").Dotf("%sList", sn),
		),
		jen.Newline(),
		jen.Var().ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").IDf("%sBaseSuite", puvn)).Call(jen.ID("nil")),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sBaseSuite", puvn)).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dotf("example%s", sn).Op("=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
			jen.ID("s").Dotf("example%sList", sn).Op("=").ID("fakes").Dotf("BuildFake%sList", sn).Call(),
		),
		jen.Newline(),
		jen.Type().IDf("%sTestSuite", puvn).Struct(
			jen.ID("suite").Dot("Suite"),
			jen.Newline(),
			jen.IDf("%sBaseSuite", puvn),
		),
		jen.Newline(),
	)

	code.Add(buildTestClientSomethingExists(proj, typ)...)
	code.Add(buildTestClientGetSomething(proj, typ)...)
	code.Add(buildTestClientGetListOfSomething(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildTestClientSearchSomething(proj, typ)...)
	}

	code.Add(buildTestClientCreateSomething(proj, typ)...)
	code.Add(buildTestClientUpdateSomething(proj, typ)...)
	code.Add(buildTestClientArchiveSomething(proj, typ)...)
	code.Add(buildTestClientGetAuditLogForSomething(proj, typ)...)

	return code
}

func buildTestClientSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_%sExists", sn).Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodHead"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodHead"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("%sExists", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_Get%s", sn).Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dotf("example%s", sn),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.Newline(),
					jen.ID("require").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
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
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_Get%s", pn).Params().Body(
			jen.Const().ID("expectedPath").Op("=").Litf("/api/v1/%s", prn),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dotf("example%sList", sn),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dotf("example%sList", sn),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.Newline(),
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
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientSearchSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_Search%s", pn).Params().Body(
			jen.Const().ID("expectedPath").Op("=").Litf("/api/v1/%s/search", prn),
			jen.Newline(),
			jen.ID("exampleQuery").Op(":=").Lit("whatever"),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("limit=20&q=whatever"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dotf("example%sList", sn).Dot(pn),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Search%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleQuery"),
						jen.Lit(0),
					),
					jen.Newline(),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dotf("example%sList", sn).Dot(pn),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with empty query"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("limit=20&q=whatever"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dotf("example%sList", sn).Dot(pn),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Search%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(""),
						jen.Lit(0),
					),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("limit").Op(":=").ID("types").Dot("DefaultQueryFilter").Call().Dot("Limit"),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Search%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleQuery"),
						jen.ID("limit"),
					),
					jen.Newline(),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("limit=20&q=whatever"),
						jen.ID("expectedPath"),
					),
					jen.ID("limit").Op(":=").ID("types").Dot("DefaultQueryFilter").Call().Dot("Limit"),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Search%s", pn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleQuery"),
						jen.ID("limit"),
					),
					jen.Newline(),
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
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_Create%s", sn).Params().Body(
			jen.Const().ID("expectedPath").Op("=").Litf("/api/v1/%s", prn),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dotf("BuildFake%sCreationInput", sn).Call(),
					jen.ID("exampleInput").Dot("BelongsToAccount").Op("=").Lit(0),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithRequestBodyValidation").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Op("&").ID("types").Dotf("%sCreationInput", sn).Values(),
						jen.ID("exampleInput"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("s").Dotf("example%s", sn),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dotf("%sCreationInput", sn).Values(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("exampleInput"),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dotf("BuildFake%sCreationInputFrom%s", sn, sn).Call(jen.ID("s").Dotf("example%s", sn)),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dotf("BuildFake%sCreationInputFrom%s", sn, sn).Call(jen.ID("s").Dotf("example%s", sn)),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
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
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_Update%s", sn).Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Update%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Update%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Update%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Update%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_Archive%s", sn).Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Archive%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Archive%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Archive%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dotf("Archive%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildTestClientGetAuditLogForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	prn := typ.Name.PluralRouteName()
	puvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").Op("*").IDf("%sTestSuite", puvn)).IDf("TestClient_GetAuditLogFor%s", sn).Params().Body(
			jen.Const().Defs(
				jen.ID("expectedPath").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d/audit"),
				jen.ID("expectedMethod").Op("=").Qual("net/http", "MethodGet"),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.ID("expectedMethod"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.ID("exampleAuditLogEntryList").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleAuditLogEntryList"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAuditLogFor%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleAuditLogEntryList"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Litf("with invalid %s ID", scn),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAuditLogFor%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.Lit(0),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error building request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAuditLogFor%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
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
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with error executing request"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAuditLogFor%s", sn).Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dotf("example%s", sn).Dot("ID"),
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
		),
		jen.Newline(),
	}

	return lines
}
