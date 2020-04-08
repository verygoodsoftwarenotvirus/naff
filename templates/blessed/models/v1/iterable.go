package v1

import (
	"fmt"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildBaseModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{jen.ID("ID").Uint64().Tag(jsonTag("id"))}

	for _, field := range typ.Fields {
		if field.Pointer {
			out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
		} else {
			out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
		}
	}

	out = append(out,
		jen.ID("CreatedOn").Uint64().Tag(jsonTag("created_on")),
		jen.ID("UpdatedOn").PointerTo().Uint64().Tag(jsonTag("updated_on")),
		jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archived_on")),
	)

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").Uint64().Tag(jsonTag("belongs_to_user")))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag(fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))))
	}

	return out
}

func buildUpdateModelStructFields(typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			if field.Pointer {
				out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			}
		}
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").Uint64().Tag(jsonTag("-")))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag(fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))))
	}

	return out
}

func buildCreateModelStructFields(typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			if field.Pointer {
				out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			}
		}
	}

	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag("-")))
	}
	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").Uint64().Tag(jsonTag("-")))
	}

	return out
}

func buildInterfaceMethods(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	interfaceMethods := []jen.Code{
		jen.IDf("%sExists", sn).Params(typ.BuildGetSomethingParams(proj)...).Params(jen.Bool(), jen.Error()),
		jen.IDf("Get%s", sn).Params(typ.BuildGetSomethingParams(proj)...).Params(jen.PointerTo().ID(sn), jen.Error()),
		jen.IDf("GetAll%sCount", pn).Params(utils.CtxParam()).Params(jen.Uint64(), jen.Error()),
		jen.IDf("Get%s", pn).Params(typ.BuildGetListOfSomethingParams(proj, true)...).Params(jen.PointerTo().IDf("%sList", sn), jen.Error()),
	}

	//if typ.BelongsToUser {
	//	interfaceMethods = append(interfaceMethods,
	//		jen.IDf("GetAll%sForUser", pn).Params(
	//			utils.CtxParam(), jen.ID("userID").Uint64(),
	//		).Params(jen.Index().ID(sn), jen.Error()))
	//}
	if typ.BelongsToStruct != nil {
		interfaceMethods = append(interfaceMethods,
			jen.IDf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular()).Params(typ.BuildGetSomethingForSomethingElseParamsForModelsPackage(proj)...).Params(jen.Index().ID(sn), jen.Error()))
	}

	interfaceMethods = append(interfaceMethods,
		jen.IDf("Create%s", sn).Params(typ.BuildCreateSomethingParams(proj, true)...).Params(jen.PointerTo().ID(sn), jen.Error()),
		jen.IDf("Update%s", sn).Params(typ.BuildUpdateSomethingParams(proj, "updated", true)...).Params(jen.Error()),
		jen.IDf("Archive%s", sn).Params(typ.BuildGetSomethingParams(proj)...).Params(jen.Error()),
	)

	return interfaceMethods
}

func iterableDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	cnwp := n.SingularCommonNameWithPrefix()
	pcn := n.PluralCommonName()
	prn := n.PluralRouteName()

	ret.Add(
		jen.Type().Defs(
			jen.Commentf("%s represents %s", sn, cnwp),
			jen.ID(sn).Struct(buildBaseModelStructFields(typ)...),
			jen.Line(),
			jen.Commentf("%sList represents a list of %s", sn, pcn),
			jen.IDf("%sList", sn).Struct(
				jen.ID("Pagination"),
				jen.ID(pn).Index().ID(sn).Tag(jsonTag(prn)),
			),
			jen.Line(),
			jen.Commentf("%sCreationInput represents what a user could set as input for creating %s", sn, pcn),
			jen.IDf("%sCreationInput", sn).Struct(buildCreateModelStructFields(typ)...),
			jen.Line(),
			jen.Commentf("%sUpdateInput represents what a user could set as input for updating %s", sn, pcn),
			jen.IDf("%sUpdateInput", sn).Struct(buildUpdateModelStructFields(typ)...),
			jen.Line(),
			jen.Commentf("%sDataManager describes a structure capable of storing %s permanently", sn, pcn),
			jen.IDf("%sDataManager", sn).Interface(buildInterfaceMethods(proj, typ)...),
			jen.Line(),
			jen.Commentf("%sDataServer describes a structure capable of serving traffic related to %s", sn, pcn),
			jen.IDf("%sDataServer", sn).Interface(
				jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Line(),
				jen.ID("ListHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("CreateHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
				jen.ID("ExistenceHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")),
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
		jen.Func().Params(jen.ID("x").PointerTo().ID(sn)).ID("Update").Params(jen.ID("input").PointerTo().IDf("%sUpdateInput", sn)).Block(buildUpdateFunctionLogic(typ.Fields)...),
		jen.Line(),
	)

	buildToUpdateInput := func() jen.Code {
		lines := []jen.Code{}

		for _, typ := range typ.Fields {
			if typ.ValidForUpdateInput {
				fsn := typ.Name.Singular()
				lines = append(lines, jen.ID(fsn).MapAssign().ID("x").Dot(fsn))
			}
		}

		return jen.Return(jen.AddressOf().IDf("%sUpdateInput", sn).Valuesln(lines...))
	}

	ret.Add(
		jen.Commentf("ToUpdateInput creates a %sUpdateInput struct for %s", sn, cnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("x").PointerTo().ID(sn)).ID("ToUpdateInput").Params().Params(jen.PointerTo().IDf("%sUpdateInput", sn)).Block(buildToUpdateInput()),
		jen.Line(),
	)

	return ret
}

func buildUpdateFunctionLogic(fields []models.DataField) []jen.Code {
	var out []jen.Code

	for i, field := range fields {
		fsn := field.Name.Singular()
		switch strings.ToLower(field.Type) {
		case "string":
			if field.Pointer {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("nil").And().PointerTo().ID("input").Dot(fsn).DoesNotEqual().EmptyString().And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
					),
				)
			} else {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().EmptyString().And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
					),
				)
			}
		case "bool",
			"float32",
			"float64",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"int",
			"int8",
			"int16",
			"int32",
			"int64":
			if field.Pointer {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("nil").And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
					),
				)
			} else {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
					),
				)
			}
		default:
			panic(fmt.Sprintf("unaccounted for type!: %q", strings.ToLower(field.Type)))
		}
		if i != len(fields)-1 {
			out = append(out, jen.Line())
		}
	}

	return out
}
