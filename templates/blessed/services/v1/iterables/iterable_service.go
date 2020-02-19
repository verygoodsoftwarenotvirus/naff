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

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	sn := typ.Name.Singular()
	cn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	srn := typ.Name.RouteName()
	prn := typ.Name.PluralRouteName()

	ret.Add(
		jen.Const().Defs(
			jen.Comment(fmt.Sprintf("CreateMiddlewareCtxKey is a string alias we can use for referring to %s input data in contexts", cn)),
			jen.ID("CreateMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit(fmt.Sprintf("%s_create_input", srn)),
			jen.Comment(fmt.Sprintf("UpdateMiddlewareCtxKey is a string alias we can use for referring to %s update data in contexts", cn)),
			jen.ID("UpdateMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit(fmt.Sprintf("%s_update_input", srn)),
			jen.Line(),
			jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName").Op("=").Lit(puvn),
			jen.ID("counterDescription").Op("=").Lit(fmt.Sprintf("the number of %s managed by the %s service", puvn, puvn)),
			jen.ID("topicName").ID("string").Op("=").Lit(prn),
			jen.ID("serviceName").ID("string").Op("=").Lit(fmt.Sprintf("%s_service", prn)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataServer", sn)).Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
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
	} else if typ.BelongsToStruct != nil {
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
		jen.Comment(fmt.Sprintf("Service handles to-do list %s", pcn)),
		jen.ID("Service").Struct(serviceLines...),
		jen.Line(),
	}

	if typ.BelongsToUser {
		typeDefs = append(typeDefs,
			jen.Comment("UserIDFetcher is a function that fetches user IDs"),
			jen.ID("UserIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		)
	} else if typ.BelongsToStruct != nil {
		typeDefs = append(typeDefs,
			jen.Commentf("%sIDFetcher is a function that fetches %s IDs", typ.BelongsToStruct.Singular(), typ.BelongsToStruct.SingularCommonName()),
			jen.IDf("%sIDFetcher", typ.BelongsToStruct.Singular()).Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		)
	}

	typeDefs = append(typeDefs,
		jen.Line(),
		jen.Comment(fmt.Sprintf("%sIDFetcher is a function that fetches %s IDs", sn, cn)),
		jen.ID(fmt.Sprintf("%sIDFetcher", sn)).Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
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
		jen.ID("ctx").Qual("context", "Context"),
		jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
		jen.ID("db").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sDataManager", sn)),
	}
	serviceValues := []jen.Code{
		jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
		jen.ID(fmt.Sprintf("%sDatabase", uvn)).Op(":").ID("db"),
		jen.ID("encoderDecoder").Op(":").ID("encoder"),
		jen.ID(fmt.Sprintf("%sCounter", uvn)).Op(":").ID(fmt.Sprintf("%sCounter", uvn)),
	}

	if typ.BelongsToUser {
		params = append(params, jen.ID("userIDFetcher").ID("UserIDFetcher"))
		serviceValues = append(serviceValues, jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"))
	} else if typ.BelongsToStruct != nil {
		params = append(params, jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).IDf("%sIDFetcher", typ.BelongsToStruct.Singular()))
		serviceValues = append(serviceValues, jen.IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()).Op(":").IDf("%sIDFetcher", typ.BelongsToStruct.UnexportedVarName()))
	}

	params = append(params,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).ID(fmt.Sprintf("%sIDFetcher", sn)),
		jen.ID("encoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
		jen.ID(fmt.Sprintf("%sCounterProvider", uvn)).Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider"),
		jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
	)
	serviceValues = append(serviceValues,
		jen.ID(fmt.Sprintf("%sIDFetcher", uvn)).Op(":").ID(fmt.Sprintf("%sIDFetcher", uvn)),
		jen.ID("reporter").Op(":").ID("reporter"),
	)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Provide%sService builds a new %sService", pn, pn)),
		jen.Line(),
		jen.Func().ID(fmt.Sprintf("Provide%sService", pn)).Paramsln(params...).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID(fmt.Sprintf("%sCounter", uvn)), jen.ID("err")).Op(":=").ID(fmt.Sprintf("%sCounterProvider", uvn)).Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(serviceValues...),
			jen.Line(),
			jen.List(jen.ID(fmt.Sprintf("%sCount", uvn)), jen.ID("err")).Op(":=").ID("svc").Dot(fmt.Sprintf("%sDatabase", uvn)).Dot(fmt.Sprintf("GetAll%sCount", pn)).Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("setting current %s count: ", cn)+"%w"), jen.ID("err"))),
			),
			jen.ID("svc").Dot(fmt.Sprintf("%sCounter", uvn)).Dot("IncrementBy").Call(jen.ID("ctx"), jen.ID(fmt.Sprintf("%sCount", uvn))),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.ID("nil")),
		),
		jen.Line(),
	}

	return lines
}
