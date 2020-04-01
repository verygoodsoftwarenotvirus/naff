package iterables

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableServiceDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(pkg, ret)

	sn := typ.Name.Singular()
	cn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	srn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	ret.Add(
		jen.Const().Defs(
			jen.Commentf("CreateMiddlewareCtxKey is a string alias we can use for referring to %s input data in contexts", cn),
			jen.ID("CreateMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Equals().Lit(fmt.Sprintf("%s_create_input", srn)),
			jen.Commentf("UpdateMiddlewareCtxKey is a string alias we can use for referring to %s update data in contexts", cn),
			jen.ID("UpdateMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Equals().Lit(fmt.Sprintf("%s_update_input", srn)),
			jen.Line(),
			jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName").Equals().Lit(puvn),
			jen.ID("counterDescription").Equals().Lit(fmt.Sprintf("the number of %s managed by the %s service", puvn, puvn)),
			jen.ID("topicName").ID("string").Equals().Lit(prn),
			jen.ID("serviceName").ID("string").Equals().Lit(fmt.Sprintf("%s_service", prn)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataServer", sn)).Equals().Parens(jen.Op("*").ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(buildServiceTypeDecl(pkg, typ)...)
	ret.Add(buildProvideServiceFuncDecl(pkg, typ)...)

	return ret
}

func buildServiceTypeDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	cn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	uvn := typ.Name.UnexportedVarName()

	serviceLines := []jen.Code{
		jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
		jen.ID(fmt.Sprintf("%sDatabase", uvn)).Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataManager", sn)),
	}

	if typ.BelongsToUser {
		serviceLines = append(serviceLines,
			jen.ID("userIDFetcher").ID("UserIDFetcher"),
		)
	}
	if typ.BelongsToStruct != nil {
		serviceLines = append(serviceLines,
			jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).IDf("%sIDFetcher", typ.BelongsToStruct.Singular()),
		)
	}

	serviceLines = append(serviceLines,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).ID(fmt.Sprintf("%sIDFetcher", sn)),
		jen.ID("encoderDecoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)

	typeDefs := []jen.Code{
		jen.Commentf("Service handles to-do list %s", pcn),
		jen.ID("Service").Struct(serviceLines...),
		jen.Line(),
	}

	if typ.BelongsToUser {
		typeDefs = append(typeDefs,
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")),
		)
	}
	if typ.BelongsToStruct != nil {
		typeDefs = append(typeDefs,
			jen.Commentf("%sIDFetcher is a function that fetches %s IDs", typ.BelongsToStruct.Singular(), typ.BelongsToStruct.SingularCommonName()),
			jen.IDf("%sIDFetcher", typ.BelongsToStruct.Singular()).Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")),
		)
	}

	typeDefs = append(typeDefs,
		jen.Line(),
		jen.Commentf("%sIDFetcher is a function that fetches %s IDs", sn, cn),
		jen.ID(fmt.Sprintf("%sIDFetcher", sn)).Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.ID("uint64")),
	)

	lines := []jen.Code{
		jen.Type().Defs(typeDefs...),
		jen.Line(),
	}

	return lines
}

func buildProvideServiceFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	cn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()

	params := []jen.Code{
		utils.CtxParam(),
		jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		jen.ID("db").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataManager", sn)),
	}
	serviceValues := []jen.Code{
		jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
		jen.ID(fmt.Sprintf("%sDatabase", uvn)).MapAssign().ID("db"),
		jen.ID("encoderDecoder").MapAssign().ID("encoder"),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).MapAssign().ID(fmt.Sprintf("%sCounter", uvn)),
	}

	if typ.BelongsToUser {
		params = append(params, jen.ID("userIDFetcher").ID("UserIDFetcher"))
		serviceValues = append(serviceValues, jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"))
	}
	if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).IDf("%sIDFetcher", typ.BelongsToStruct.Singular()))
		serviceValues = append(serviceValues, jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).MapAssign().IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()))
	}

	params = append(params,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).ID(fmt.Sprintf("%sIDFetcher", sn)),
		jen.ID("encoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
		jen.ID(fmt.Sprintf("%sCounterProvider", uvn)).Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)
	serviceValues = append(serviceValues,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).MapAssign().ID(fmt.Sprintf("%sIDFetcher", uvn)),
		jen.ID("reporter").MapAssign().ID("reporter"),
	)

	lines := []jen.Code{
		jen.Commentf("Provide%sService builds a new %sService", pn, pn),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sService", pn)).Paramsln(params...).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID(fmt.Sprintf("%sCounter", uvn)), jen.Err()).Assign().ID(fmt.Sprintf("%sCounterProvider", uvn)).Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("svc").Assign().VarPointer().ID("Service").Valuesln(serviceValues...),
			jen.Line(),
			jen.List(jen.ID(fmt.Sprintf("%sCount", uvn)), jen.Err()).Assign().ID("svc").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("GetAll%sCount", pn)).Call(utils.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("setting current %s count: ", cn)+"%w"), jen.Err())),
			),
			jen.ID("svc").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("IncrementBy").Call(utils.CtxVar(), jen.ID(fmt.Sprintf("%sCount", uvn))),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
