package iterables

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, ret)

	sn := typ.Name.Singular()
	cn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	srn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	ret.Add(
		jen.Const().Defs(
			jen.Commentf("CreateMiddlewareCtxKey is a string alias we can use for referring to %s input data in contexts.", cn),
			jen.ID("CreateMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit(fmt.Sprintf("%s_create_input", srn)),
			jen.Commentf("UpdateMiddlewareCtxKey is a string alias we can use for referring to %s update data in contexts.", cn),
			jen.ID("UpdateMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit(fmt.Sprintf("%s_update_input", srn)),
			jen.Line(),
			jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName").Equals().Lit(puvn),
			jen.ID("counterDescription").String().Equals().Lit(fmt.Sprintf("the number of %s managed by the %s service", puvn, puvn)),
			jen.ID("topicName").String().Equals().Lit(prn),
			jen.ID("serviceName").String().Equals().Lit(fmt.Sprintf("%s_service", prn)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataServer", sn)).Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(buildServiceTypeDecl(proj, typ)...)
	ret.Add(buildProvideServiceFuncDecl(proj, typ)...)

	return ret
}

func buildServiceTypeDecl(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	cn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()

	structFields := []jen.Code{
		jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
	}

	// data managers
	for _, ot := range proj.FindOwnerTypeChain(typ) {
		structFields = append(structFields,
			jen.ID(fmt.Sprintf("%sDataManager", ot.Name.UnexportedVarName())).Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", ot.Name.Singular())),
		)
	}
	structFields = append(structFields, jen.ID(fmt.Sprintf("%sDataManager", uvn)).Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn)))

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
		jen.ID(fmt.Sprintf("%sCounter", uvn)).Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
	)

	structFields = append(structFields,
		jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)

	typeDefs := []jen.Code{
		jen.Commentf("Service handles to-do list %s", pcn),
		jen.ID("Service").Struct(structFields...),
		jen.Line(),
	}

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
		jen.ID(constants.LoggerVarName).Qual(utils.LoggingPkg, "Logger"),
	}

	for _, ot := range proj.FindOwnerTypeChain(typ) {
		params = append(params, jen.IDf("%sDataManager", ot.Name.UnexportedVarName()).Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", ot.Name.Singular())))
	}
	params = append(params, jen.IDf("%sDataManager", uvn).Qual(proj.ModelsV1Package(), fmt.Sprintf("%sDataManager", sn)))

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
		jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
		jen.ID(fmt.Sprintf("%sCounterProvider", uvn)).Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)
	serviceValues = append(serviceValues,
		jen.ID("reporter").MapAssign().ID("reporter"),
	)

	lines := []jen.Code{
		jen.Commentf("Provide%sService builds a new %sService.", pn, pn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sService", pn)).Paramsln(params...).Params(jen.PointerTo().ID("Service"), jen.Error()).Block(
			jen.List(jen.ID(fmt.Sprintf("%sCounter", uvn)), jen.Err()).Assign().ID(fmt.Sprintf("%sCounterProvider", uvn)).Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
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
