package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func auditLogEntriesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuilder_BuildGetAuditLogEntryRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/admin/audit_log/%d"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("exampleAuditLogEntry").Op(":=").ID("fakes").Dot("BuildFakeAuditLogEntry").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit(""),
						jen.ID("expectedPath"),
						jen.ID("exampleAuditLogEntry").Dot("ID"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetAuditLogEntryRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("exampleAuditLogEntry").Dot("ID"),
					),
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
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuilder_BuildGetAuditLogEntriesRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Var().Defs(
				jen.ID("expectedPath").Op("=").Lit("/api/v1/admin/audit_log"),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("helper").Op(":=").ID("buildTestHelper").Call(),
					jen.ID("filter").Op(":=").ID("types").Dot("DefaultQueryFilter").Call(),
					jen.ID("spec").Op(":=").ID("newRequestSpec").Call(
						jen.ID("true"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("includeArchived=false&limit=20&page=1&sortBy=asc"),
						jen.ID("expectedPath"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("helper").Dot("builder").Dot("BuildGetAuditLogEntriesRequest").Call(
						jen.ID("helper").Dot("ctx"),
						jen.ID("filter"),
					),
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
		),
		jen.Line(),
	)

	return code
}
