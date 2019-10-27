package v1

import (
	"fmt"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildBaseModelStructFields(fields []models.DataField) []jen.Code {
	out := []jen.Code{jen.ID("ID").ID("uint64").Tag(jsonTag("id"))}

	for _, field := range fields {
		out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
	}

	out = append(out,
		jen.ID("CreatedOn").ID("uint64").Tag(jsonTag("created_on")),
		jen.ID("UpdatedOn").Op("*").ID("uint64").Tag(jsonTag("updated_on")),
		jen.ID("ArchivedOn").Op("*").ID("uint64").Tag(jsonTag("archived_on")),
		jen.ID("BelongsTo").ID("uint64").Tag(jsonTag("belongs_to")),
	)

	return out
}

func buildUpdateModelStructFields(fields []models.DataField) []jen.Code {
	var out []jen.Code

	for _, field := range fields {
		if field.ValidForUpdateInput {
			out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
		}
	}

	out = append(out, jen.ID("BelongsTo").ID("uint64").Tag(jsonTag("-")))
	return out
}

func buildCreateModelStructFields(fields []models.DataField) []jen.Code {
	var out []jen.Code

	for _, field := range fields {
		if field.ValidForCreationInput {
			out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
		}
	}

	out = append(out, jen.ID("BelongsTo").ID("uint64").Tag(jsonTag("-")))
	return out
}

func iterableDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	uvn := n.UnexportedVarName()
	cnwp := n.SingularCommonNameWithPrefix()
	pcn := n.PluralCommonName()
	prn := n.PluralRouteName()

	ret.Add(
		jen.Type().Defs(
			jen.Commentf("%s represents %s", sn, cnwp),
			jen.ID(sn).Struct(buildBaseModelStructFields(typ.Fields)...),
			jen.Line(),
			jen.Commentf("%sList represents a list of %s", sn, pcn),
			jen.IDf("%sList", sn).Struct(
				jen.ID("Pagination"),
				jen.ID(pn).Index().ID(sn).Tag(jsonTag(prn)),
			),
			jen.Line(),
			jen.Commentf("%sCreationInput represents what a user could set as input for creating %s", sn, pcn),
			jen.IDf("%sCreationInput", sn).Struct(buildCreateModelStructFields(typ.Fields)...),
			jen.Line(),
			jen.Commentf("%sUpdateInput represents what a user could set as input for updating %s", sn, pcn),
			jen.IDf("%sUpdateInput", sn).Struct(buildUpdateModelStructFields(typ.Fields)...),
			jen.Line(),
			jen.Commentf("%sDataManager describes a structure capable of storing %s permanently", sn, pcn),
			jen.IDf("%sDataManager", sn).Interface(
				jen.IDf("Get%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID(sn), jen.ID("error")),
				jen.IDf("Get%sCount", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.IDf("Get%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").IDf("%sList", sn), jen.ID("error")),
				jen.IDf("GetAll%sForUser", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().ID(sn), jen.ID("error")),
				jen.IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").IDf("%sCreationInput", sn)).Params(jen.Op("*").ID(sn), jen.ID("error")),
				jen.IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID(sn)).Params(jen.ID("error")),
				jen.IDf("Archive%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("id"), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Commentf("%sDataServer describes a structure capable of serving traffic related to %s", sn, pcn),
			jen.IDf("%sDataServer", sn).Interface(
				jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Line(),
				jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ReadHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("UpdateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ArchiveHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Update merges an %sInput with %s", sn, cnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID(sn)).ID("Update").Params(jen.ID("input").Op("*").IDf("%sUpdateInput", sn)).Block(buildUpdateFunctionLogic(typ.Fields)...),
		jen.Line(),
	)
	return ret
}

// var zeroValMap = map[string]interface{}{
// 	"string": "",
// 	"uint64": 0,
// }

func buildUpdateFunctionLogic(fields []models.DataField) []jen.Code {
	var out []jen.Code

	for i, field := range fields {
		fsn := field.Name.Singular()
		switch strings.ToLower(field.Type) {
		case "string":
			out = append(
				out,
				jen.If(jen.ID("input").Dot(fsn).Op("!=").Lit("").Op("||").ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Block(
					jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
				),
			)
		case "uint64":
			out = append(
				out,
				jen.If(jen.ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Block(
					jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
				),
			)
		default:
			panic(fmt.Sprintf("unaccounted for type!: %q", strings.ToLower(field.Type)))
		}
		if i != len(fields)-1 {
			out = append(out, jen.Line())
		}
	}

	return out
}