package client

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// func buildPathForType(pkg *models.Project, typ models.DataType) string {
// 	if typ.BelongsToStruct != nil {
// 		parentType := pkg.FindType(typ.BelongsToStruct.Singular())
// 		if parentType != nil {
// 			return filepath.Join(buildPathForType(pkg, *parentType), )
// 		}
// 	}
// }

func iterablesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

	basePath := fmt.Sprintf("%sBasePath", typ.Name.PluralUnexportedVarName())

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)
	ret.Add(jen.Const().Defs(
		jen.ID(basePath).Op("=").Lit(typ.Name.PluralRouteName())),
	)

	ret.Add(buildBuildGetSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildGetSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildBuildGetSomethingsRequestFuncDecl(pkg, typ)...)
	ret.Add(buildGetSomethingsFuncDecl(pkg, typ)...)
	ret.Add(buildBuildCreateSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildCreateSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildBuildUpdateSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildUpdateSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildBuildArchiveSomethingRequestFuncDecl(pkg, typ)...)
	ret.Add(buildArchiveSomethingFuncDecl(pkg, typ)...)

	return ret
}

func buildBuildGetSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	pvn := typ.Name.PluralUnexportedVarName()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)
	basePath := fmt.Sprintf("%sBasePath", pvn)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", ts)).Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(pkg, typ)...).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildParamsForMethodThatHandlesAnInstanceOfADataType(pkg *models.Project, typ models.DataType) []jen.Code {
	parents := pkg.FindOwnerTypeChain(typ)
	listParams := []jen.Code{}
	params := []jen.Code{
		utils.CtxParam(),
	}

	if len(parents) > 0 {
		for _, pt := range parents {
			listParams = append(listParams, jen.IDf("%sID", pt.Name.UnexportedVarName()))
		}
		listParams = append(listParams, jen.IDf("%sID", typ.Name.UnexportedVarName()))
		params = append(params, jen.List(listParams...).ID("uint64"))
	} else {
		params = append(params, jen.IDf("%sID", typ.Name.UnexportedVarName()).ID("uint64"))
	}

	return params
}

func buildGetSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()

	commonNameWithPrefix := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Get%s retrieves %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", ts)).Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(pkg, typ)...).Params(
			jen.ID(uvn).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(jen.ID("ctx"), jen.IDf("%sID", uvn)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(jen.Lit("building request: %w"), jen.ID("err")),
				),
			),
			jen.Line(),
			jen.If(jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(jen.ID("ctx"), jen.ID("req"), jen.Op("&").ID(uvn)), jen.ID("retrieveErr").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("retrieveErr")),
			),
			jen.Line(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildGetSomethingsRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	tp := typ.Name.Plural() // title plural
	pvn := typ.Name.PluralUnexportedVarName()

	basePath := fmt.Sprintf("%sBasePath", pvn)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", tp, typ.Name.PluralCommonName())),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", tp)).Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID(basePath),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildGetSomethingsFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	tp := typ.Name.Plural() // title plural
	pvn := typ.Name.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Get%s retrieves a list of %s", tp, typ.Name.PluralCommonName())),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", tp)).Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
		).Params(
			jen.ID(pvn).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", ts)),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.If(
				jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(jen.ID("ctx"), jen.ID("req"), jen.Op("&").ID(pvn)),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("retrieveErr")),
			),
			jen.Line(),
			jen.Return().List(jen.ID(pvn), jen.ID("nil")),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildCreateSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	pvn := typ.Name.PluralUnexportedVarName()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)
	basePath := fmt.Sprintf("%sBasePath", pvn)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildCreate%sRequest builds an HTTP request for creating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildCreate%sRequest", ts)).Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", ts)),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildCreateSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	vn := typ.Name.UnexportedVarName()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Create%s creates %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Create%s", ts)).Params(
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", ts)),
		).Params(
			jen.ID(vn).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.Line(),
			jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(vn),
			),
			jen.Return().List(
				jen.ID(vn),
				jen.ID("err"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildUpdateSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	pvn := typ.Name.PluralUnexportedVarName()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)
	basePath := fmt.Sprintf("%sBasePath", pvn)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildUpdate%sRequest builds an HTTP request for updating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildUpdate%sRequest", ts)).Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("updated").Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("updated"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildUpdateSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Update%s updates %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Update%s", ts)).Params(utils.CtxParam(), jen.ID("updated").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), ts)).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("updated"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildArchiveSomethingRequestFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()
	pvn := typ.Name.PluralUnexportedVarName()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)
	basePath := fmt.Sprintf("%sBasePath", pvn)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("BuildArchive%sRequest builds an HTTP request for updating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildArchive%sRequest", ts)).Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(pkg, typ)...).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.IDf("%sID", typ.Name.UnexportedVarName()),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildArchiveSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	ts := typ.Name.Singular()

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)

	lines := []jen.Code{
		jen.Comment(fmt.Sprintf("Archive%s archives %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Archive%s", ts)).Params(buildParamsForMethodThatHandlesAnInstanceOfADataType(pkg, typ)...).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", typ.Name.UnexportedVarName()),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
	}

	return lines
}
