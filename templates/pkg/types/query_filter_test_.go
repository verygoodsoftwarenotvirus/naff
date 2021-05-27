package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestQueryFilter_AttachToLogger").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(100), jen.ID("Limit").Op(":").ID("MaxLimit"), jen.ID("CreatedAfter").Op(":").Lit(123456789), jen.ID("CreatedBefore").Op(":").Lit(123456789), jen.ID("UpdatedAfter").Op(":").Lit(123456789), jen.ID("UpdatedBefore").Op(":").Lit(123456789), jen.ID("SortBy").Op(":").ID("SortDescending"), jen.ID("IncludeArchived").Op(":").ID("true")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("qf").Dot("AttachToLogger").Call(jen.ID("logger")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.Parens(jen.Op("*").ID("QueryFilter")).Call(jen.ID("nil")).Dot("AttachToLogger").Call(jen.ID("logger")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQueryFilter_FromParams").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("actual").Op(":=").Op("&").ID("QueryFilter").Values(),
					jen.ID("expected").Op(":=").Op("&").ID("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(100), jen.ID("Limit").Op(":").ID("MaxLimit"), jen.ID("CreatedAfter").Op(":").Lit(123456789), jen.ID("CreatedBefore").Op(":").Lit(123456789), jen.ID("UpdatedAfter").Op(":").Lit(123456789), jen.ID("UpdatedBefore").Op(":").Lit(123456789), jen.ID("SortBy").Op(":").ID("SortDescending"), jen.ID("IncludeArchived").Op(":").ID("true")),
					jen.ID("exampleInput").Op(":=").Qual("net/url", "Values").Valuesln(jen.ID("pageQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Page")))), jen.ID("LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Limit")))), jen.ID("createdBeforeQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedAfter")))), jen.ID("createdAfterQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedBefore")))), jen.ID("updatedBeforeQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedAfter")))), jen.ID("updatedAfterQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedBefore")))), jen.ID("sortByQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("string").Call(jen.ID("expected").Dot("SortBy"))), jen.ID("includeArchivedQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "FormatBool").Call(jen.ID("true")))),
					jen.ID("actual").Dot("FromParams").Call(jen.ID("exampleInput")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
					jen.ID("exampleInput").Index(jen.ID("sortByQueryKey")).Op("=").Index().ID("string").Valuesln(jen.ID("string").Call(jen.ID("SortAscending"))),
					jen.ID("actual").Dot("FromParams").Call(jen.ID("exampleInput")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("SortAscending"),
						jen.ID("actual").Dot("SortBy"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQueryFilter_SetPage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("QueryFilter").Values(),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
					jen.ID("qf").Dot("SetPage").Call(jen.ID("expected")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("qf").Dot("Page"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQueryFilter_QueryPage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("QueryFilter").Valuesln(jen.ID("Limit").Op(":").Lit(10), jen.ID("Page").Op(":").Lit(11)),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(100)),
					jen.ID("actual").Op(":=").ID("qf").Dot("QueryPage").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQueryFilter_ToValues").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(100), jen.ID("Limit").Op(":").Lit(50), jen.ID("CreatedAfter").Op(":").Lit(123456789), jen.ID("CreatedBefore").Op(":").Lit(123456789), jen.ID("UpdatedAfter").Op(":").Lit(123456789), jen.ID("UpdatedBefore").Op(":").Lit(123456789), jen.ID("IncludeArchived").Op(":").ID("true"), jen.ID("SortBy").Op(":").ID("SortDescending")),
					jen.ID("expected").Op(":=").Qual("net/url", "Values").Valuesln(jen.ID("pageQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("Page")))), jen.ID("LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("Limit")))), jen.ID("createdBeforeQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("CreatedAfter")))), jen.ID("createdAfterQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("CreatedBefore")))), jen.ID("updatedBeforeQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("UpdatedAfter")))), jen.ID("updatedAfterQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("UpdatedBefore")))), jen.ID("includeArchivedQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "FormatBool").Call(jen.ID("qf").Dot("IncludeArchived"))), jen.ID("sortByQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("string").Call(jen.ID("qf").Dot("SortBy")))),
					jen.ID("actual").Op(":=").ID("qf").Dot("ToValues").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Parens(jen.Op("*").ID("QueryFilter")).Call(jen.ID("nil")),
					jen.ID("expected").Op(":=").ID("DefaultQueryFilter").Call().Dot("ToValues").Call(),
					jen.ID("actual").Op(":=").ID("qf").Dot("ToValues").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestExtractQueryFilter").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expected").Op(":=").Op("&").ID("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(100), jen.ID("Limit").Op(":").ID("MaxLimit"), jen.ID("CreatedAfter").Op(":").Lit(123456789), jen.ID("CreatedBefore").Op(":").Lit(123456789), jen.ID("UpdatedAfter").Op(":").Lit(123456789), jen.ID("UpdatedBefore").Op(":").Lit(123456789), jen.ID("SortBy").Op(":").ID("SortDescending")),
					jen.ID("exampleInput").Op(":=").Qual("net/url", "Values").Valuesln(jen.ID("pageQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Page")))), jen.ID("LimitQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Limit")))), jen.ID("createdBeforeQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedAfter")))), jen.ID("createdAfterQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedBefore")))), jen.ID("updatedBeforeQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedAfter")))), jen.ID("updatedAfterQueryKey").Op(":").Index().ID("string").Valuesln(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedBefore")))), jen.ID("sortByQueryKey").Op(":").Index().ID("string").Valuesln(jen.ID("string").Call(jen.ID("expected").Dot("SortBy")))),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("https://verygoodsoftwarenotvirus.ru"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("require").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("req"),
					),
					jen.ID("req").Dot("URL").Dot("RawQuery").Op("=").ID("exampleInput").Dot("Encode").Call(),
					jen.ID("actual").Op(":=").ID("ExtractQueryFilter").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
