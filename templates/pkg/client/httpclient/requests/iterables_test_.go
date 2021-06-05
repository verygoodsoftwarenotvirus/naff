package requests

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
	prn := typ.Name.PluralRouteName()

	code.Add(
		jen.Func().IDf("TestBuilder_Build%sExistsRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("Build%sExistsRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodHead"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("Build%sExistsRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
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
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuilder_BuildGet%sRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildGet%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildGet%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
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
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuilder_BuildGet%sRequest", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Newline(),
			jen.Const().ID("expectedPath").Op("=").Litf("/api/v1/%s", prn),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.ID("filter").Op(":=").Parens(jen.Op("*").ID("types").Dot("QueryFilter")).Call(jen.ID("nil")),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildGet%sRequest", pn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("filter"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
		),
		jen.Newline(),
	)

	if typ.SearchEnabled {
		code.Add(
			jen.Func().IDf("TestBuilder_BuildSearch%sRequest", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Newline(),
				jen.Const().ID("expectedPath").Op("=").Litf("/api/v1/%s/search", prn),
				jen.Newline(),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard"),
					jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
						jen.Newline(),
						jen.ID("limit").Op(":=").ID("types").Dot("DefaultQueryFilter").Call().Dot("Limit"),
						jen.ID("exampleQuery").Op(":=").Lit("whatever"),
						jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
							jen.ID("true"),
							jen.Qual("net/http", "MethodGet"),
							jen.Lit("limit=20&q=whatever"),
							jen.ID("expectedPath"),
						),
						jen.Newline(),
						jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildSearch%sRequest", pn).Call(
							jen.ID("helper").Dot("ctx"),
							jen.ID("exampleQuery"),
							jen.ID("limit"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.Newline(),
						jen.ID("assertRequestQuality").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("spec"),
						),
					),
				),
			),
			jen.Newline(),
		)
	}

	code.Add(
		jen.Func().IDf("TestBuilder_BuildCreate%sRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Const().ID("expectedPath").Op("=").Litf("/api/v1/%s", prn),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dotf("BuildFake%sCreationInput", sn).Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleInput"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit(""),
						jen.ID("expectedPath"),
					),
					jen.Newline(),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
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
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildCreate%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.Op("&").ID("types").Dotf("%sCreationInput", sn).Values(),
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
		jen.Func().IDf("TestBuilder_BuildUpdate%sRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("false"),
						jen.Qual("net/http", "MethodPut"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildUpdate%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.IDf("example%s", sn),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildUpdate%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
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
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuilder_BuildArchive%sRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Const().ID("expectedPathFormat").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d"),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodDelete"),
						jen.Lit(""),
						jen.ID("expectedPathFormat"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildArchive%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildArchive%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
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
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuilder_BuildGetAuditLogFor%sRequest", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Const().ID("expectedPath").Op("=").Lit(fmt.Sprintf("/api/v1/%s/", prn)+"%d/audit"),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.IDf("example%s", sn).Op(":=").ID("fakes").Dotf("BuildFake%s", sn).Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("assertRequestQuality").Call(
						jen.ID("t"),
						jen.ID("actual"),
						jen.ID("spec"),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.Newline(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dotf("BuildGetAuditLogFor%sRequest", sn).Call(
						jen.ID("helper").Dot("ctx"),
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
		),
		jen.Newline(),
	)

	return code
}
