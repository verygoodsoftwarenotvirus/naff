package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("TestFromParams").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("actual").Assign().VarPointer().ID("QueryFilter").Values(),
				jen.ID("expected").Assign().VarPointer().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().ID("MaxLimit"),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.Line(),
				jen.ID("exampleInput").Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("pageKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Page")))),
					jen.ID("limitKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Limit")))),
					jen.ID("createdBeforeKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedAfter")))),
					jen.ID("createdAfterKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedBefore")))),
					jen.ID("updatedBeforeKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedAfter")))),
					jen.ID("updatedAfterKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedBefore")))),
					jen.ID("sortByKey").MapAssign().Index().ID("string").Values(jen.ID("string").Call(jen.ID("expected").Dot("SortBy"))),
				),
				jen.Line(),
				jen.ID("actual").Dot("FromParams").Call(jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("exampleInput").Index(jen.ID("sortByKey")).Equals().Index().ID("string").Values(jen.ID("string").Call(jen.ID("SortAscending"))),
				jen.Line(),
				jen.ID("actual").Dot("FromParams").Call(jen.ID("exampleInput")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("SortAscending"), jen.ID("actual").Dot("SortBy")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestQueryFilter_SetPage").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("qf").Assign().VarPointer().ID("QueryFilter").Values(),
				jen.ID("expected").Assign().ID("uint64").Call(utils.FakeUint32Func()),
				jen.ID("qf").Dot("SetPage").Call(jen.ID("expected")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("qf").Dot("Page")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestQueryFilter_QueryPage").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("qf").Assign().VarPointer().ID("QueryFilter").Values(jen.ID("Limit").MapAssign().Lit(10), jen.ID("Page").MapAssign().Lit(11)),
				jen.ID("expected").Assign().ID("uint64").Call(jen.Lit(100)),
				jen.ID("actual").Assign().ID("qf").Dot("QueryPage").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestQueryFilter_ToValues").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("qf").Assign().VarPointer().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().Lit(50),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.ID("expected").Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("pageKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("Page")))),
					jen.ID("limitKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("Limit")))),
					jen.ID("createdBeforeKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("CreatedAfter")))),
					jen.ID("createdAfterKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("CreatedBefore")))),
					jen.ID("updatedBeforeKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("UpdatedAfter")))),
					jen.ID("updatedAfterKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("qf").Dot("UpdatedBefore")))),
					jen.ID("sortByKey").MapAssign().Index().ID("string").Values(jen.ID("string").Call(jen.ID("qf").Dot("SortBy"))),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("qf").Dot("ToValues").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("qf").Assign().Parens(jen.Op("*").ID("QueryFilter")).Call(jen.Nil()),
				jen.ID("expected").Assign().ID("DefaultQueryFilter").Call().Dot("ToValues").Call(),
				jen.ID("actual").Assign().ID("qf").Dot("ToValues").Call(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestQueryFilter_ApplyToQueryBuilder").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("baseQueryBuilder").Assign().Qual("github.com/Masterminds/squirrel", "StatementBuilder").Dot("PlaceholderFormat").Call(jen.Qual("github.com/Masterminds/squirrel", "Dollar")).
				Dotln("Select").Call(jen.Lit("things")).
				Dotln("From").Call(jen.Lit("stuff")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("condition").MapAssign().ID("true"))),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("qf").Assign().VarPointer().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().Lit(50),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.Line(),
				jen.ID("sb").Assign().Qual("github.com/Masterminds/squirrel", "StatementBuilder").Dot("Select").Call(jen.Lit("*")).Dot("From").Call(jen.Lit("testing")),
				jen.ID("qf").Dot("ApplyToQueryBuilder").Call(jen.ID("sb")),
				jen.ID("expected").Assign().Lit("SELECT * FROM testing"),
				jen.List(jen.ID("actual"), jen.ID("_"), jen.Err()).Assign().ID("sb").Dot("ToSql").Call(),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("basic usecase"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleQF").Assign().VarPointer().ID("QueryFilter").Values(jen.ID("Limit").MapAssign().Lit(15), jen.ID("Page").MapAssign().Lit(2)),
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE condition = $1 LIMIT 15 OFFSET 15`),
				jen.ID("x").Assign().ID("exampleQF").Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder")),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("args")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("returns query builder if query filter is nil"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE condition = $1`),
				jen.ID("x").Assign().Parens(jen.Op("*").ID("QueryFilter")).Call(jen.Nil()).Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder")),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("args")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("whole kit and kaboodle"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleQF").Assign().VarPointer().ID("QueryFilter").Valuesln(
					jen.ID("Limit").MapAssign().Lit(20), jen.ID("Page").MapAssign().Lit(6),
					jen.ID("CreatedAfter").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("CreatedBefore").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("UpdatedAfter").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
					jen.ID("UpdatedBefore").MapAssign().ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()),
				),
				jen.Line(),
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE condition = $1 AND created_on > $2 AND created_on < $3 AND updated_on > $4 AND updated_on < $5 LIMIT 20 OFFSET 100`),
				jen.ID("x").Assign().ID("exampleQF").Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder")),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("args")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with zero limit"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("exampleQF").Assign().VarPointer().ID("QueryFilter").Values(jen.ID("Limit").MapAssign().Lit(0), jen.ID("Page").MapAssign().Lit(1)),
				jen.ID("expected").Assign().Lit(`SELECT things FROM stuff WHERE condition = $1 LIMIT 250`),
				jen.ID("x").Assign().ID("exampleQF").Dot("ApplyToQueryBuilder").Call(jen.ID("baseQueryBuilder")),
				jen.List(jen.ID("actual"), jen.ID("args"), jen.Err()).Assign().ID("x").Dot("ToSql").Call(),
				jen.Line(),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual"), jen.Lit("expected and actual queries don't match")),
				jen.Qual("github.com/stretchr/testify/assert", "Nil").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/assert", "NotEmpty").Call(jen.ID("t"), jen.ID("args")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestExtractQueryFilter").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.ID("expected").Assign().VarPointer().ID("QueryFilter").Valuesln(
					jen.ID("Page").MapAssign().Lit(100),
					jen.ID("Limit").MapAssign().ID("MaxLimit"),
					jen.ID("CreatedAfter").MapAssign().Lit(123456789),
					jen.ID("CreatedBefore").MapAssign().Lit(123456789),
					jen.ID("UpdatedAfter").MapAssign().Lit(123456789),
					jen.ID("UpdatedBefore").MapAssign().Lit(123456789),
					jen.ID("SortBy").MapAssign().ID("SortDescending"),
				),
				jen.ID("exampleInput").Assign().Qual("net/url", "Values").Valuesln(
					jen.ID("pageKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Page")))),
					jen.ID("limitKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("Limit")))),
					jen.ID("createdBeforeKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedAfter")))),
					jen.ID("createdAfterKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("CreatedBefore")))),
					jen.ID("updatedBeforeKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedAfter")))),
					jen.ID("updatedAfterKey").MapAssign().Index().ID("string").Values(jen.Qual("strconv", "Itoa").Call(jen.ID("int").Call(jen.ID("expected").Dot("UpdatedBefore")))),
					jen.ID("sortByKey").MapAssign().Index().ID("string").Values(jen.ID("string").Call(jen.ID("expected").Dot("SortBy"))),
				),
				jen.Line(),
				jen.List(jen.ID("req"), jen.Err()).Assign().Qual("net/http", "NewRequest").Call(jen.Qual("net/http", "MethodGet"), jen.Lit("https://verygoodsoftwarenotvirus.ru"), jen.Nil()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				jen.Qual("github.com/stretchr/testify/require", "NotNil").Call(jen.ID("t"), jen.ID("req")),
				jen.Line(),
				jen.ID("req").Dot("URL").Dot("RawQuery").Equals().ID("exampleInput").Dot("Encode").Call(),
				jen.ID("actual").Assign().ID("ExtractQueryFilter").Call(jen.ID("req")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)
	return ret
}
