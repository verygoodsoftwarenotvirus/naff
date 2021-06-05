package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestItems").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("suite").Dot("Run").Call(
				jen.ID("t"),
				jen.ID("new").Call(jen.ID("itemsTestSuite")),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().ID("itemsBaseSuite").Struct(
			jen.ID("suite").Dot("Suite"),
			jen.Newline(),
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("exampleItem").Op("*").ID("types").Dot("Item"),
			jen.ID("exampleItemList").Op("*").ID("types").Dot("ItemList"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().ID("_").ID("suite").Dot("SetupTestSuite").Op("=").Parens(jen.Op("*").ID("itemsBaseSuite")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsBaseSuite")).ID("SetupTest").Params().Body(
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleItem").Op("=").ID("fakes").Dot("BuildFakeItem").Call(),
			jen.ID("s").Dot("exampleItemList").Op("=").ID("fakes").Dot("BuildFakeItemList").Call(),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().ID("itemsTestSuite").Struct(
			jen.ID("suite").Dot("Suite"),
			jen.Newline(),
			jen.ID("itemsBaseSuite"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_ItemExists").Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit("/api/v1/items/%d"),
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
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
				jen.Lit("with invalid item ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodHead"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("ItemExists").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_GetItem").Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit("/api/v1/items/%d"),
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
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleItem"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
						jen.ID("s").Dot("exampleItem"),
						jen.ID("actual"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleItem"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_GetItems").Params().Body(
			jen.Const().ID("expectedPath").Op("=").Lit("/api/v1/items"),
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
						jen.ID("s").Dot("exampleItemList"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
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
						jen.ID("s").Dot("exampleItemList"),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetItems").Call(
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
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_SearchItems").Params().Body(
			jen.Const().ID("expectedPath").Op("=").Lit("/api/v1/items/search"),
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
						jen.ID("s").Dot("exampleItemList").Dot("Items"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchItems").Call(
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
						jen.ID("s").Dot("exampleItemList").Dot("Items"),
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
						jen.ID("s").Dot("exampleItemList").Dot("Items"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchItems").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchItems").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("SearchItems").Call(
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
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_CreateItem").Params().Body(
			jen.Const().ID("expectedPath").Op("=").Lit("/api/v1/items"),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInput").Call(),
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
						jen.Op("&").ID("types").Dot("ItemCreationInput").Values(),
						jen.ID("exampleInput"),
						jen.ID("s").Dot("exampleItem"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
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
						jen.ID("s").Dot("exampleItem"),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
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
					jen.ID("exampleInput").Op(":=").Op("&").ID("types").Dot("ItemCreationInput").Values(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
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
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("s").Dot("exampleItem")),
					jen.Newline(),
					jen.ID("c").Op(":=").ID("buildTestClientWithInvalidURL").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
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
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("s").Dot("exampleItem")),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientThatWaitsTooLong").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("CreateItem").Call(
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
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_UpdateItem").Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit("/api/v1/items/%d"),
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
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("s").Dot("exampleItem"),
					),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dot("UpdateItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem"),
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
					jen.ID("err").Op(":=").ID("c").Dot("UpdateItem").Call(
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
					jen.ID("err").Op(":=").ID("c").Dot("UpdateItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem"),
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
					jen.ID("err").Op(":=").ID("c").Dot("UpdateItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_ArchiveItem").Params().Body(
			jen.Const().ID("expectedPathFormat").Op("=").Lit("/api/v1/items/%d"),
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
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithStatusCodeResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.Newline(),
			jen.ID("s").Dot("Run").Call(
				jen.Lit("with invalid item ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveItem").Call(
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
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
					jen.ID("err").Op(":=").ID("c").Dot("ArchiveItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("itemsTestSuite")).ID("TestClient_GetAuditLogForItem").Params().Body(
			jen.Const().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/items/%d/audit"),
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
						jen.ID("s").Dot("exampleItem").Dot("ID"),
					),
					jen.ID("exampleAuditLogEntryList").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntryList").Call().Dot("Entries"),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildTestClientWithJSONResponse").Call(
						jen.ID("t"),
						jen.ID("spec"),
						jen.ID("exampleAuditLogEntryList"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
				jen.Lit("with invalid item ID"),
				jen.Func().Params().Body(
					jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
					jen.Newline(),
					jen.List(jen.ID("c"), jen.ID("_")).Op(":=").ID("buildSimpleTestClient").Call(jen.ID("t")),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForItem").Call(
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot("GetAuditLogForItem").Call(
						jen.ID("s").Dot("ctx"),
						jen.ID("s").Dot("exampleItem").Dot("ID"),
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
	)

	return code
}
