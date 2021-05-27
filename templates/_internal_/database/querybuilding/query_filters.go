package querybuilding

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFiltersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("ApplyFilterToQueryBuilder applies the query filter to a query builder."),
		jen.Line(),
		jen.Func().ID("ApplyFilterToQueryBuilder").Params(jen.ID("qf").Op("*").ID("types").Dot("QueryFilter"), jen.ID("tableName").ID("string"), jen.ID("queryBuilder").ID("squirrel").Dot("SelectBuilder")).Params(jen.ID("squirrel").Dot("SelectBuilder")).Body(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Body(
				jen.Return().ID("queryBuilder")),
			jen.ID("qf").Dot("SetPage").Call(jen.ID("qf").Dot("Page")),
			jen.If(jen.ID("qp").Op(":=").ID("qf").Dot("QueryPage").Call(), jen.ID("qp").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Offset").Call(jen.ID("qp"))),
			jen.If(jen.ID("qf").Dot("Limit").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Limit").Call(jen.ID("uint64").Call(jen.ID("qf").Dot("Limit")))).Else().Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Limit").Call(jen.ID("types").Dot("MaxLimit"))),
			jen.If(jen.ID("qf").Dot("CreatedAfter").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("CreatedOnColumn"),
				).Op(":").ID("qf").Dot("CreatedAfter")))),
			jen.If(jen.ID("qf").Dot("CreatedBefore").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("CreatedOnColumn"),
				).Op(":").ID("qf").Dot("CreatedBefore")))),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("LastUpdatedOnColumn"),
				).Op(":").ID("qf").Dot("UpdatedAfter")))),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("LastUpdatedOnColumn"),
				).Op(":").ID("qf").Dot("UpdatedBefore")))),
			jen.Return().ID("queryBuilder"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ApplyFilterToSubCountQueryBuilder applies the query filter to a query builder."),
		jen.Line(),
		jen.Func().ID("ApplyFilterToSubCountQueryBuilder").Params(jen.ID("qf").Op("*").ID("types").Dot("QueryFilter"), jen.ID("tableName").ID("string"), jen.ID("queryBuilder").ID("squirrel").Dot("SelectBuilder")).Params(jen.ID("squirrel").Dot("SelectBuilder")).Body(
			jen.If(jen.ID("qf").Op("==").ID("nil")).Body(
				jen.Return().ID("queryBuilder")),
			jen.If(jen.ID("qf").Dot("CreatedAfter").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("CreatedOnColumn"),
				).Op(":").ID("qf").Dot("CreatedAfter")))),
			jen.If(jen.ID("qf").Dot("CreatedBefore").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("CreatedOnColumn"),
				).Op(":").ID("qf").Dot("CreatedBefore")))),
			jen.If(jen.ID("qf").Dot("UpdatedAfter").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("LastUpdatedOnColumn"),
				).Op(":").ID("qf").Dot("UpdatedAfter")))),
			jen.If(jen.ID("qf").Dot("UpdatedBefore").Op(">").Lit(0)).Body(
				jen.ID("queryBuilder").Op("=").ID("queryBuilder").Dot("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.ID("tableName"),
					jen.ID("LastUpdatedOnColumn"),
				).Op(":").ID("qf").Dot("UpdatedBefore")))),
			jen.Return().ID("queryBuilder"),
		),
		jen.Line(),
	)

	return code
}
