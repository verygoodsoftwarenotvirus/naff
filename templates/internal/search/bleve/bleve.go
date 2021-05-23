package bleve

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const searchPackage = "github.com/blevesearch/bleve"

func bleveDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildConstantDefinitions(proj)...)
	code.Add(buildInterfaceImplementationStatement(proj)...)
	code.Add(buildTypeDefinitions(proj)...)
	code.Add(buildNewBleveIndexManager(proj)...)
	code.Add(buildNewBleveIndexManager_Index(proj)...)
	code.Add(buildNewBleveIndexManager_Search(proj)...)
	code.Add(buildNewBleveIndexManager_Delete(proj)...)

	return code
}

func buildConstantDefinitions(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("base").Equals().Lit(10),
			jen.ID("bitSize").Equals().Lit(64),
			jen.Line(),
			jen.Comment("testingSearchIndexName is an index name that is only valid for testing's sake."),
			jen.ID("testingSearchIndexName").Qual(proj.InternalSearchPackage(), "IndexName").Equals().Lit("testing"),
		),
	}

	return lines
}

func buildInterfaceImplementationStatement(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().Qual(proj.InternalSearchPackage(), "IndexManager").Equals().Parens(jen.PointerTo().ID("bleveIndexManager")).Parens(jen.Nil()),
		jen.Line(),
	}

	return lines
}

func buildTypeDefinitions(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.ID("bleveIndexManager").Struct(
				jen.ID("index").Qual(searchPackage, "Index"),
				proj.LoggerParam(),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildNewBleveIndexManager(proj *models.Project) []jen.Code {
	nameCases := []jen.Code{
		jen.Case(jen.ID("testingSearchIndexName")).Body(
			jen.List(jen.ID("index"), jen.ID("newIndexErr")).Equals().Qual(searchPackage, "New").Call(
				jen.String().Call(jen.ID("path")),
				jen.Qual(searchPackage, "NewIndexMapping").Call(),
			),
			jen.If(jen.ID("newIndexErr").DoesNotEqual().Nil()).Body(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.ID("newIndexErr"),
					jen.Lit("failed to create new index"),
				),
				jen.Return(jen.Nil(), jen.ID("newIndexErr")),
			),
		),
	}

	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			nameCases = append(nameCases,
				jen.Case(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sSearchIndexName", typ.Name.Plural()))).Body(
					jen.List(jen.ID("index"), jen.ID("newIndexErr")).Equals().Qual(searchPackage, "New").Call(
						jen.String().Call(jen.ID("path")),
						jen.IDf("build%sMapping", typ.Name.Singular()).Call(),
					),
					jen.If(jen.ID("newIndexErr").DoesNotEqual().Nil()).Body(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(
							jen.ID("newIndexErr"),
							jen.Lit("failed to create new index"),
						),
						jen.Return(jen.Nil(), jen.ID("newIndexErr")),
					),
				),
			)
		}
	}

	nameCases = append(nameCases,
		jen.Default().Body(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("invalid index name: %q"),
				jen.ID("name"),
			)),
		),
	)

	lines := []jen.Code{
		jen.Comment("NewBleveIndexManager instantiates a bleve index"),
		jen.Line(),
		jen.Func().ID("NewBleveIndexManager").Params(
			jen.ID("path").Qual(proj.InternalSearchPackage(), "IndexPath"),
			jen.ID("name").Qual(proj.InternalSearchPackage(), "IndexName"),
			proj.LoggerParam(),
		).Params(
			jen.Qual(proj.InternalSearchPackage(), "IndexManager"),
			jen.Error(),
		).Body(
			jen.Var().ID("index").Qual(searchPackage, "Index"),
			jen.Line(),
			jen.List(jen.ID("preexistingIndex"), jen.ID("openIndexErr")).Assign().Qual(searchPackage, "Open").Call(
				jen.String().Call(jen.ID("path")),
			),
			jen.Switch(jen.ID("openIndexErr")).Body(
				jen.Case(jen.Nil()).Body(
					jen.ID("index").Equals().ID("preexistingIndex"),
				),
				jen.Case(jen.Qual(searchPackage, "ErrorIndexPathDoesNotExist")).Body(
					jen.ID(constants.LoggerVarName).Dot("WithValue").Call(
						jen.Lit("path"),
						jen.ID("path"),
					).Dot("Debug").Call(
						jen.Lit("tried to open existing index, but didn't find it"),
					),
					jen.Var().ID("newIndexErr").Error(),
					jen.Line(),
					jen.Switch(jen.ID("name")).Body(
						nameCases...,
					),
				),
				jen.Default().Body(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(
						jen.ID("openIndexErr"),
						jen.Lit("failed to open index"),
					),
					jen.Return(jen.Nil(), jen.ID("openIndexErr")),
				),
			),
			jen.Line(),
			jen.ID("im").Assign().AddressOf().ID("bleveIndexManager").Valuesln(
				jen.ID("index").MapAssign().ID("index"),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s_search"), jen.ID("name")),
				),
			),
			jen.Line(),
			jen.Return(jen.ID("im"), jen.Nil()),
		),
	}

	return lines
}

func buildNewBleveIndexManager_Index(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("Index implements our IndexManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("Index").Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
			jen.ID("value").Interface(),
		).Error().Body(
			utils.StartSpan(proj, false, "Index"),
			jen.ID("sm").Dot(constants.LoggerVarName).Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			).Dot("Debug").Call(jen.Lit("adding to index")),
			jen.Return(jen.ID("sm").Dot("index").Dot("Index").Call(
				jen.Qual("strconv", "FormatUint").Call(jen.ID("id"), jen.ID("base")),
				jen.ID("value"),
			)),
		),
	}

	return lines
}

func buildNewBleveIndexManager_Search(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("Search implements our IndexManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("Search").Params(
			constants.CtxParam(),
			jen.ID("query").String(),
			constants.UserIDParam(),
		).Params(
			jen.ID("ids").Index().Uint64(),
			jen.Err().Error(),
		).Body(
			utils.StartSpan(proj, false, "Search"),
			jen.ID("query").Equals().ID("ensureQueryIsRestrictedToUser").Call(jen.ID("query"), jen.ID(constants.UserIDVarName)),
			jen.Qual(proj.InternalTracingPackage(), "AttachSearchQueryToSpan").Call(jen.ID("span"), jen.ID("query")),
			jen.ID("sm").Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
				jen.Lit("search_query").MapAssign().ID("query"),
				jen.Lit("user_id").MapAssign().ID(constants.UserIDVarName),
			)).Dot("Debug").Call(jen.Lit("performing search")),
			jen.Line(),
			jen.ID("searchRequest").Assign().Qual(searchPackage, "NewSearchRequest").Call(
				jen.Qual(searchPackage, "NewQueryStringQuery").Call(jen.ID("query")),
			),
			jen.List(jen.ID("searchResults"), jen.Err()).Assign().ID("sm").Dot("index").Dot("SearchInContext").Call(
				constants.CtxVar(),
				jen.ID("searchRequest"),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Body(
				jen.ID("sm").Dot(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("performing search query"),
				),
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.ID("out").Assign().Index().Uint64().Values(),
			jen.For(jen.List(jen.Underscore(), jen.ID("result")).Assign().Range().ID("searchResults").Dot("Hits")).Body(
				jen.List(jen.ID("x"), jen.Err()).Assign().Qual("strconv", "ParseUint").Call(
					jen.ID("result").Dot("ID"),
					jen.ID("base"),
					jen.ID("bitSize"),
				),
				jen.If(jen.Err().DoesNotEqual().Nil()).Body(
					jen.Comment("this should literally never happen"),
					jen.Return(jen.Nil(), jen.Err()),
				),
				jen.ID("out").Equals().Append(jen.ID("out"), jen.ID("x")),
			),
			jen.Line(),
			jen.Return(jen.ID("out"), jen.Nil()),
		),
	}

	return lines
}

func buildNewBleveIndexManager_Delete(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("Delete implements our IndexManager interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("Delete").Params(
			constants.CtxParam(),
			jen.ID("id").Uint64(),
		).Error().Body(
			utils.StartSpan(proj, false, "Delete"),
			jen.ID("sm").Dot(constants.LoggerVarName).Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			).Dot("Debug").Call(jen.Lit("removing from index")),
			jen.Return(jen.ID("sm").Dot("index").Dot("Delete").Call(
				jen.Qual("strconv", "FormatUint").Call(jen.ID("id"), jen.ID("base")),
			)),
		),
	}

	return lines
}
