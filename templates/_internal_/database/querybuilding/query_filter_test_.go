package querybuilding

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func queryFilterTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestQueryFilter_ApplyFilterToQueryBuilder").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleTableName").Op(":=").Lit("stuff"),
			jen.ID("baseQueryBuilder").Op(":=").ID("squirrel").Dot("StatementBuilder").Dot("PlaceholderFormat").Call(jen.ID("squirrel").Dot("Dollar")).Dot("Select").Call(jen.Lit("things")).Dot("From").Call(jen.ID("exampleTableName")).Dot("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.condition"),
				jen.ID("exampleTableName"),
			).Op(":").ID("true"))),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("types").Dot("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(100), jen.ID("Limit").Op(":").Lit(50), jen.ID("CreatedAfter").Op(":").Lit(123456789), jen.ID("CreatedBefore").Op(":").Lit(123456789), jen.ID("UpdatedAfter").Op(":").Lit(123456789), jen.ID("UpdatedBefore").Op(":").Lit(123456789), jen.ID("SortBy").Op(":").ID("types").Dot("SortDescending")),
					jen.ID("sb").Op(":=").ID("squirrel").Dot("StatementBuilder").Dot("Select").Call(jen.Lit("*")).Dot("From").Call(jen.Lit("testing")),
					jen.ID("sb").Op("=").ID("ApplyFilterToQueryBuilder").Call(
						jen.ID("qf"),
						jen.ID("exampleTableName"),
						jen.ID("sb"),
					),
					jen.ID("expected").Op(":=").Lit("SELECT * FROM testing WHERE stuff.created_on > ? AND stuff.created_on < ? AND stuff.last_updated_on > ? AND stuff.last_updated_on < ? LIMIT 50 OFFSET 4950"),
					jen.List(jen.ID("actual"), jen.ID("_"), jen.ID("err")).Op(":=").ID("sb").Dot("ToSql").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("sb").Op(":=").ID("squirrel").Dot("StatementBuilder").Dot("Select").Call(jen.Lit("*")).Dot("From").Call(jen.Lit("testing")),
					jen.ID("sb").Op("=").ID("ApplyFilterToQueryBuilder").Call(
						jen.ID("nil"),
						jen.ID("exampleTableName"),
						jen.ID("sb"),
					),
					jen.ID("expected").Op(":=").Lit("SELECT * FROM testing"),
					jen.List(jen.ID("actual"), jen.ID("_"), jen.ID("err")).Op(":=").ID("sb").Dot("ToSql").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("basic usage"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("types").Dot("QueryFilter").Valuesln(jen.ID("Limit").Op(":").Lit(15), jen.ID("Page").Op(":").Lit(2)),
					jen.ID("expected").Op(":=").Lit("SELECT things FROM stuff WHERE stuff.condition = $1 LIMIT 15 OFFSET 15"),
					jen.ID("x").Op(":=").ID("ApplyFilterToQueryBuilder").Call(
						jen.ID("qf"),
						jen.ID("exampleTableName"),
						jen.ID("baseQueryBuilder"),
					),
					jen.List(jen.ID("actual"), jen.ID("args"), jen.ID("err")).Op(":=").ID("x").Dot("ToSql").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
						jen.Lit("expected and actual queries don't match"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("args"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("whole kit and kaboodle"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("types").Dot("QueryFilter").Valuesln(jen.ID("Limit").Op(":").Lit(20), jen.ID("Page").Op(":").Lit(6), jen.ID("CreatedAfter").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()), jen.ID("CreatedBefore").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()), jen.ID("UpdatedAfter").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call()), jen.ID("UpdatedBefore").Op(":").ID("uint64").Call(jen.Qual("time", "Now").Call().Dot("Unix").Call())),
					jen.ID("expected").Op(":=").Lit("SELECT things FROM stuff WHERE stuff.condition = $1 AND stuff.created_on > $2 AND stuff.created_on < $3 AND stuff.last_updated_on > $4 AND stuff.last_updated_on < $5 LIMIT 20 OFFSET 100"),
					jen.ID("x").Op(":=").ID("ApplyFilterToQueryBuilder").Call(
						jen.ID("qf"),
						jen.ID("exampleTableName"),
						jen.ID("baseQueryBuilder"),
					),
					jen.List(jen.ID("actual"), jen.ID("args"), jen.ID("err")).Op(":=").ID("x").Dot("ToSql").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
						jen.Lit("expected and actual queries don't match"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("args"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with zero limit"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("types").Dot("QueryFilter").Valuesln(jen.ID("Limit").Op(":").Lit(0), jen.ID("Page").Op(":").Lit(1)),
					jen.ID("expected").Op(":=").Lit("SELECT things FROM stuff WHERE stuff.condition = $1 LIMIT 250"),
					jen.ID("x").Op(":=").ID("ApplyFilterToQueryBuilder").Call(
						jen.ID("qf"),
						jen.ID("exampleTableName"),
						jen.ID("baseQueryBuilder"),
					),
					jen.List(jen.ID("actual"), jen.ID("args"), jen.ID("err")).Op(":=").ID("x").Dot("ToSql").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
						jen.Lit("expected and actual queries don't match"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.ID("args"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestQueryFilter_ApplyFilterToSubCountQueryBuilder").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("exampleTableName").Op(":=").Lit("stuff"),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("qf").Op(":=").Op("&").ID("types").Dot("QueryFilter").Valuesln(jen.ID("Page").Op(":").Lit(100), jen.ID("Limit").Op(":").Lit(50), jen.ID("CreatedAfter").Op(":").Lit(123456789), jen.ID("CreatedBefore").Op(":").Lit(123456789), jen.ID("UpdatedAfter").Op(":").Lit(123456789), jen.ID("UpdatedBefore").Op(":").Lit(123456789), jen.ID("SortBy").Op(":").ID("types").Dot("SortDescending")),
					jen.ID("sb").Op(":=").ID("squirrel").Dot("StatementBuilder").Dot("Select").Call(jen.Lit("*")).Dot("From").Call(jen.Lit("testing")),
					jen.ID("sb").Op("=").ID("ApplyFilterToSubCountQueryBuilder").Call(
						jen.ID("qf"),
						jen.ID("exampleTableName"),
						jen.ID("sb"),
					),
					jen.ID("expected").Op(":=").Lit("SELECT * FROM testing WHERE stuff.created_on > ? AND stuff.created_on < ? AND stuff.last_updated_on > ? AND stuff.last_updated_on < ?"),
					jen.List(jen.ID("actual"), jen.ID("_"), jen.ID("err")).Op(":=").ID("sb").Dot("ToSql").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil filter"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("sb").Op(":=").ID("squirrel").Dot("StatementBuilder").Dot("Select").Call(jen.Lit("*")).Dot("From").Call(jen.Lit("testing")),
					jen.ID("sb").Op("=").ID("ApplyFilterToSubCountQueryBuilder").Call(
						jen.ID("nil"),
						jen.ID("exampleTableName"),
						jen.ID("sb"),
					),
					jen.ID("expected").Op(":=").Lit("SELECT * FROM testing"),
					jen.List(jen.ID("actual"), jen.ID("_"), jen.ID("err")).Op(":=").ID("sb").Dot("ToSql").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
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
