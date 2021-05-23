package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(buildSomethingServiceConstantDefs(proj, typ)...)
	code.Add(buildSomethingServiceVarDefs(proj, typ)...)
	code.Add(buildServiceTypeDecls(proj, typ)...)
	code.Add(buildProvideServiceFuncDecl(proj, typ)...)

	if typ.SearchEnabled {
		code.Add(buildProvideServiceSearchIndexFuncDecl(proj, typ)...)
	}

	return code
}

func buildSomethingServiceConstantDefs(proj *models.Project, typ models.DataType) []jen.Code {

	cn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	srn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	lines := []jen.Code{
		jen.Const().Defs(
			jen.Commentf("createMiddlewareCtxKey is a string alias we can use for referring to %s input data in contexts.", cn),
			jen.ID("createMiddlewareCtxKey").Qual(proj.TypesPackage(), "ContextKey").Equals().Lit(fmt.Sprintf("%s_create_input", srn)),
			jen.Commentf("updateMiddlewareCtxKey is a string alias we can use for referring to %s update data in contexts.", cn),
			jen.ID("updateMiddlewareCtxKey").Qual(proj.TypesPackage(), "ContextKey").Equals().Lit(fmt.Sprintf("%s_update_input", srn)),
			jen.Line(),
			jen.ID("counterName").Qual(proj.InternalMetricsPackage(), "CounterName").Equals().Lit(puvn),
			jen.ID("counterDescription").String().Equals().Lit(fmt.Sprintf("the number of %s managed by the %s service", puvn, puvn)),
			jen.ID("topicName").String().Equals().Lit(prn),
			jen.ID("serviceName").String().Equals().Lit(fmt.Sprintf("%s_service", prn)),
		),
		jen.Line(),
	}

	return lines
}

func buildSomethingServiceVarDefs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.Var().Defs(
			jen.Underscore().Qual(proj.TypesPackage(), fmt.Sprintf("%sDataServer", sn)).Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceTypeDecls(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	cn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()

	structFields := []jen.Code{
		proj.LoggerParam(),
	}

	// data managers
	for _, ot := range proj.FindOwnerTypeChain(typ) {
		structFields = append(structFields,
			jen.ID(fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", ot.Name.Singular())),
		)
	}
	structFields = append(structFields, jen.ID(fmt.Sprintf("%sDataManager", uvn)).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)))

	// id fetchers
	for _, ot := range proj.FindOwnerTypeChain(typ) {
		structFields = append(structFields,
			jen.IDf("%sIDFetcher", ot.Name.UnexportedVarName()).IDf("%sIDFetcher", ot.Name.Singular()),
		)
	}

	structFields = append(structFields, jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).ID(fmt.Sprintf("%sIDFetcher", sn)))
	if typ.OwnedByAUserAtSomeLevel(proj) {
		structFields = append(structFields,
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
		)
	}

	structFields = append(structFields,
		jen.ID(fmt.Sprintf("%sCounter", uvn)).Qual(proj.InternalMetricsPackage(), "UnitCounter"),
	)

	structFields = append(structFields,
		jen.ID("encoderDecoder").Qual(proj.InternalEncodingPackage(), "EncoderDecoder"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)

	typeDefs := []jen.Code{}
	if typ.SearchEnabled {
		typeDefs = append(typeDefs,
			jen.Comment("SearchIndex is a type alias for dependency injection's sake"),
			jen.ID("SearchIndex").Qual(proj.InternalSearchPackage(), "IndexManager"),
			jen.Line(),
		)
		structFields = append(structFields,
			jen.ID("search").ID("SearchIndex"),
		)
	}

	typeDefs = append(typeDefs,
		jen.Commentf("Service handles to-do list %s", pcn),
		jen.ID("Service").Struct(structFields...),
		jen.Line(),
	)

	if typ.OwnedByAUserAtSomeLevel(proj) {
		typeDefs = append(typeDefs,
			jen.Comment("UserIDFetcher is a function that fetches user IDs."),
			jen.ID("UserIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
			jen.Line(),
		)
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		typeDefs = append(typeDefs,
			jen.Commentf("%sIDFetcher is a function that fetches %s IDs.", ot.Name.Singular(), ot.Name.SingularCommonName()),
			jen.IDf("%sIDFetcher", ot.Name.Singular()).Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
			jen.Line(),
		)
	}

	typeDefs = append(typeDefs,
		jen.Line(),
		jen.Commentf("%sIDFetcher is a function that fetches %s IDs.", sn, cn),
		jen.ID(fmt.Sprintf("%sIDFetcher", sn)).Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
	)

	lines := []jen.Code{
		jen.Type().Defs(typeDefs...),
		jen.Line(),
	}

	return lines
}

func buildProvideServiceFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()

	params := []jen.Code{
		proj.LoggerParam(),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		params = append(params, jen.IDf("%sDataManager", ot.Name.UnexportedVarName()).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", ot.Name.Singular())))
	}
	params = append(params, jen.IDf("%sDataManager", uvn).Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)))

	serviceValues := []jen.Code{
		jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		params = append(params, jen.IDf("%sIDFetcher", ot.Name.UnexportedVarName()).IDf("%sIDFetcher", ot.Name.Singular()))
		serviceValues = append(serviceValues, jen.IDf("%sIDFetcher", ot.Name.UnexportedVarName()).MapAssign().IDf("%sIDFetcher", ot.Name.UnexportedVarName()))
	}
	params = append(params, jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).ID(fmt.Sprintf("%sIDFetcher", sn)))
	serviceValues = append(serviceValues, jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).MapAssign().ID(fmt.Sprintf("%sIDFetcher", uvn)))

	if typ.OwnedByAUserAtSomeLevel(proj) {
		params = append(params, jen.ID("userIDFetcher").ID("UserIDFetcher"))
		serviceValues = append(serviceValues, jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"))
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		serviceValues = append(serviceValues, jen.IDf("%sDataManager", ot.Name.UnexportedVarName()).MapAssign().IDf("%sDataManager", ot.Name.UnexportedVarName()))
	}

	serviceValues = append(serviceValues,
		jen.IDf("%sDataManager", uvn).MapAssign().IDf("%sDataManager", uvn),
		jen.ID("encoderDecoder").MapAssign().ID("encoder"),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).MapAssign().ID(fmt.Sprintf("%sCounter", uvn)),
	)

	params = append(params,
		jen.ID("encoder").Qual(proj.InternalEncodingPackage(), "EncoderDecoder"),
		jen.ID(fmt.Sprintf("%sCounterProvider", uvn)).Qual(proj.InternalMetricsPackage(), "UnitCounterProvider"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)
	serviceValues = append(serviceValues,
		jen.ID("reporter").MapAssign().ID("reporter"),
	)

	if typ.SearchEnabled {
		params = append(params,
			jen.ID("searchIndexManager").ID("SearchIndex"),
		)
		serviceValues = append(serviceValues,
			jen.ID("search").MapAssign().ID("searchIndexManager"),
		)
	}

	lines := []jen.Code{
		jen.Commentf("Provide%sService builds a new %sService.", pn, pn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sService", pn)).Paramsln(params...).Params(jen.PointerTo().ID("Service"), jen.Error()).Body(
			jen.List(jen.ID(fmt.Sprintf("%sCounter", uvn)), jen.Err()).Assign().ID(fmt.Sprintf("%sCounterProvider", uvn)).Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(serviceValues...),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideServiceSearchIndexFuncDecl(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Commentf("Provide%sServiceSearchIndex provides a search index for the service", pn),
		jen.Line(),
		jen.Func().IDf("Provide%sServiceSearchIndex", pn).Paramsln(
			jen.ID("searchSettings").Qual(proj.InternalConfigPackage(), "SearchSettings"),
			jen.ID("indexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider"),
			proj.LoggerParam(),
		).Params(jen.ID("SearchIndex"), jen.Error()).Body(
			jen.ID(constants.LoggerVarName).Dot("WithValue").Call(
				jen.Lit("index_path"),
				jen.ID("searchSettings").Dotf("%sIndexPath", pn),
			).Dot("Debug").Call(jen.Litf("setting up %s search index", pcn)),
			jen.Line(),
			jen.List(jen.ID("searchIndex"), jen.ID("indexInitErr")).Assign().ID("indexProvider").Call(
				jen.ID("searchSettings").Dotf("%sIndexPath", pn),
				jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sSearchIndexName", pn)),
				jen.ID(constants.LoggerVarName),
			),
			jen.If(jen.ID("indexInitErr").DoesNotEqual().Nil()).Body(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.ID("indexInitErr"), jen.Litf("setting up %s search index", pcn)),
				jen.Return(jen.Nil(), jen.ID("indexInitErr")),
			),
			jen.Line(),
			jen.Return(jen.ID("searchIndex"), jen.Nil()),
		),
	}

	return lines
}
