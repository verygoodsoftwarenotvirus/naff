package client

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

	ts := typ.Name.Singular()
	vn := typ.Name.UnexportedVarName()
	pvn := typ.Name.PluralUnexportedVarName()
	prn := typ.Name.PluralRouteName()
	tp := typ.Name.Plural() // title plural

	commonName := strings.Join(strings.Split(typ.Name.RouteName(), "_"), " ")
	commonNameWithPrefix := fmt.Sprintf("%s %s", wordsmith.AOrAn(ts), commonName)
	basePath := fmt.Sprintf("%sBasePath", pvn)

	utils.AddImports(ret)
	ret.Add(jen.Const().Defs(
		jen.ID(basePath).Op("=").Lit(prn)),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", ts)).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Get%s retrieves %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", ts)).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID(vn).Op("*").Qual(utils.ModelsPkg, ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					)),
			),
			jen.Line(),
			jen.If(jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(vn),
			),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("retrieveErr"),
				),
			),
			jen.Line(),
			jen.Return().List(
				jen.ID(vn),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", tp, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", tp)).Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Get%s retrieves a list of %s", tp, typ.Name.Plural())),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", tp)).Params(
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
		).Params(
			jen.ID(pvn).Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sList", ts)),
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
				jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID(pvn),
				),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(jen.ID("nil"),
					jen.ID("retrieveErr"),
				),
			),
			jen.Line(),
			jen.Return().List(
				jen.ID(pvn),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildCreate%sRequest builds an HTTP request for creating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildCreate%sRequest", ts)).Params(
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)),
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Create%s creates %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Create%s", ts)).Params(
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)),
		).Params(
			jen.ID(vn).Op("*").Qual(utils.ModelsPkg, ts),
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildUpdate%sRequest builds an HTTP request for updating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildUpdate%sRequest", ts)).Params(
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(utils.ModelsPkg, ts),
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Update%s updates %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Update%s", ts)).Params(utils.CtxParam(), jen.ID("updated").Op("*").Qual(utils.ModelsPkg, ts)).Params(jen.ID("error")).Block(
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildArchive%sRequest builds an HTTP request for updating %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildArchive%sRequest", ts)).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
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
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Archive%s archives %s", ts, commonNameWithPrefix)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Archive%s", ts)).Params(
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.ID("id"),
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
	)
	return ret
}
