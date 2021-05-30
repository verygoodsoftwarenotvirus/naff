package sqlite

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func genericDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildQueryOnly builds a given query, handles whatever errs and returns just the query and args."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("buildQueryOnly").Params(jen.ID("span").ID("tracing").Dot("Span"), jen.ID("builder").ID("squirrel").Dot("Sqlizer")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("query"), jen.ID("_"), jen.ID("err")).Op(":=").ID("builder").Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Return().ID("query"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildQuery builds a given query, handles whatever errs and returns just the query and args."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("buildQuery").Params(jen.ID("span").ID("tracing").Dot("Span"), jen.ID("builder").ID("squirrel").Dot("Sqlizer")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op(":=").ID("builder").Dot("ToSql").Call(),
			jen.ID("b").Dot("logQueryBuildingError").Call(
				jen.ID("span"),
				jen.ID("err"),
			),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("buildTotalCountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("tableName"), jen.ID("ownershipColumn")).ID("string"), jen.ID("userID").ID("uint64"), jen.List(jen.ID("forAdmin"), jen.ID("includeArchived")).ID("bool")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("where").Op(":=").ID("squirrel").Dot("Eq").Valuesln(),
			jen.ID("totalCountQueryBuilder").Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.ID("columnCountQueryTemplate"),
				jen.ID("tableName"),
			)).Dot("From").Call(jen.ID("tableName")),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.If(jen.ID("userID").Op("!=").Lit(0).Op("&&").ID("ownershipColumn").Op("!=").Lit("")).Body(
					jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("tableName"),
						jen.ID("ownershipColumn"),
					)).Op("=").ID("userID")),
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				)).Op("=").ID("nil"),
			).Else().If(jen.Op("!").ID("includeArchived")).Body(
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				)).Op("=").ID("nil")),
			jen.If(jen.ID("len").Call(jen.ID("where")).Op(">").Lit(0)).Body(
				jen.ID("totalCountQueryBuilder").Op("=").ID("totalCountQueryBuilder").Dot("Where").Call(jen.ID("where"))),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("totalCountQueryBuilder"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).ID("buildFilteredCountQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("tableName"), jen.ID("ownershipColumn")).ID("string"), jen.ID("userID").ID("uint64"), jen.List(jen.ID("forAdmin"), jen.ID("includeArchived")).ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.ID("where").Op(":=").ID("squirrel").Dot("Eq").Valuesln(),
			jen.ID("filteredCountQueryBuilder").Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.ID("columnCountQueryTemplate"),
				jen.ID("tableName"),
			)).Dot("From").Call(jen.ID("tableName")),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.If(jen.ID("userID").Op("!=").Lit(0).Op("&&").ID("ownershipColumn").Op("!=").Lit("")).Body(
					jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("tableName"),
						jen.ID("ownershipColumn"),
					)).Op("=").ID("userID")),
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				)).Op("=").ID("nil"),
			).Else().If(jen.Op("!").ID("includeArchived")).Body(
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				)).Op("=").ID("nil")),
			jen.If(jen.ID("len").Call(jen.ID("where")).Op(">").Lit(0)).Body(
				jen.ID("filteredCountQueryBuilder").Op("=").ID("filteredCountQueryBuilder").Dot("Where").Call(jen.ID("where"))),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("filteredCountQueryBuilder").Op("=").ID("querybuilding").Dot("ApplyFilterToSubCountQueryBuilder").Call(
					jen.ID("filter"),
					jen.ID("tableName"),
					jen.ID("filteredCountQueryBuilder"),
				)),
			jen.Return().ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("filteredCountQueryBuilder"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildListQuery builds a SQL query selecting rows that adhere to a given QueryFilter and belong to a given account,"),
		jen.Line(),
		jen.Func().Comment("and returns both the query and the relevant args to pass to the query executor.").Params(jen.ID("b").Op("*").ID("Sqlite")).ID("buildListQuery").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("tableName"), jen.ID("ownershipColumn")).ID("string"), jen.ID("columns").Index().ID("string"), jen.ID("ownerID").ID("uint64"), jen.ID("forAdmin").ID("bool"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("tracing").Dot("AttachFilterToSpan").Call(
					jen.ID("span"),
					jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				)),
			jen.Var().Defs(
				jen.ID("includeArchived").ID("bool"),
			),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("includeArchived").Op("=").ID("filter").Dot("IncludeArchived")),
			jen.List(jen.ID("filteredCountQuery"), jen.ID("filteredCountQueryArgs")).Op(":=").ID("b").Dot("buildFilteredCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("tableName"),
				jen.ID("ownershipColumn"),
				jen.ID("ownerID"),
				jen.ID("forAdmin"),
				jen.ID("includeArchived"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("totalCountQuery"), jen.ID("totalCountQueryArgs")).Op(":=").ID("b").Dot("buildTotalCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("tableName"),
				jen.ID("ownershipColumn"),
				jen.ID("ownerID"),
				jen.ID("forAdmin"),
				jen.ID("includeArchived"),
			),
			jen.ID("builder").Op(":=").ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.ID("append").Call(
				jen.ID("columns"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s) as total_count"),
					jen.ID("totalCountQuery"),
				),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("(%s) as filtered_count"),
					jen.ID("filteredCountQuery"),
				),
			).Op("...")).Dot("From").Call(jen.ID("tableName")),
			jen.If(jen.Op("!").ID("forAdmin")).Body(
				jen.ID("w").Op(":=").ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("querybuilding").Dot("ArchivedOnColumn"),
				).Op(":").ID("nil")),
				jen.If(jen.ID("ownershipColumn").Op("!=").Lit("").Op("&&").ID("ownerID").Op("!=").Lit(0)).Body(
					jen.ID("w").Index(jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.ID("tableName"),
						jen.ID("ownershipColumn"),
					)).Op("=").ID("ownerID")),
				jen.ID("builder").Op("=").ID("builder").Dot("Where").Call(jen.ID("w")),
			),
			jen.ID("builder").Op("=").ID("builder").Dot("GroupBy").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.ID("tableName"),
				jen.ID("querybuilding").Dot("IDColumn"),
			)),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.ID("builder").Op("=").ID("querybuilding").Dot("ApplyFilterToQueryBuilder").Call(
					jen.ID("filter"),
					jen.ID("tableName"),
					jen.ID("builder"),
				)),
			jen.List(jen.ID("query"), jen.ID("selectArgs")).Op(":=").ID("b").Dot("buildQuery").Call(
				jen.ID("span"),
				jen.ID("builder"),
			),
			jen.Return().List(jen.ID("query"), jen.ID("append").Call(
				jen.ID("append").Call(
					jen.ID("filteredCountQueryArgs"),
					jen.ID("totalCountQueryArgs").Op("..."),
				),
				jen.ID("selectArgs").Op("..."),
			)),
		),
		jen.Line(),
	)

	return code
}
