package v1

import (
	"fmt"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildBaseModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{jen.ID("ID").ID("uint64").Tag(jsonTag("id"))}

	for _, field := range typ.Fields {
		if field.Pointer {
			out = append(out, jen.ID(field.Name.Singular()).Op("*").ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
		} else {
			out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
		}
	}

	out = append(out,
		jen.ID("CreatedOn").ID("uint64").Tag(jsonTag("created_on")),
		jen.ID("UpdatedOn").Op("*").ID("uint64").Tag(jsonTag("updated_on")),
		jen.ID("ArchivedOn").Op("*").ID("uint64").Tag(jsonTag("archived_on")),
	)

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("belongs_to_user")))
	} else if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).ID("uint64").Tag(jsonTag(fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))))
	}

	return out
}

func buildUpdateModelStructFields(typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			if field.Pointer {
				out = append(out, jen.ID(field.Name.Singular()).Op("*").ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			}
		}
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("-")))
	} else if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).ID("uint64").Tag(jsonTag(fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))))
	}

	return out
}

func buildCreateModelStructFields(typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			if field.Pointer {
				out = append(out, jen.ID(field.Name.Singular()).Op("*").ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.RouteName())))
			}
		}
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID("BelongsToUser").ID("uint64").Tag(jsonTag("-")))
	} else if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).ID("uint64").Tag(jsonTag("-")))
	}

	return out
}

func iterableDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	uvn := n.UnexportedVarName()
	cnwp := n.SingularCommonNameWithPrefix()
	pcn := n.PluralCommonName()
	prn := n.PluralRouteName()

	interfaceMethods := []jen.Code{}

	getSomethingParams := []jen.Code{jen.ID("ctx").Qual("context", "Context")}
	secondGetSomethingParams := []jen.Code{jen.IDf("%sID", uvn)}

	getSomethingCountParams := []jen.Code{
		utils.CtxParam(),
		jen.ID("filter").Op("*").ID("QueryFilter"),
	}
	getAllSomethingCount := []jen.Code{jen.ID("ctx").Qual("context", "Context")}
	getListOfSomethingParams := []jen.Code{
		utils.CtxParam(),
		jen.ID("filter").Op("*").ID("QueryFilter"),
	}
	archiveSomethingParams := []jen.Code{jen.ID("ctx").Qual("context", "Context")}
	secondArchiveSomethingParams := []jen.Code{jen.IDf("%sID", uvn)}

	if typ.BelongsToUser {
		secondGetSomethingParams = append(secondGetSomethingParams, jen.ID("userID"))
		secondArchiveSomethingParams = append(secondArchiveSomethingParams, jen.ID("userID"))

		getSomethingCountParams = append(getSomethingCountParams, jen.ID("userID").ID("uint64"))
		getListOfSomethingParams = append(getListOfSomethingParams, jen.ID("userID").ID("uint64"))
	} else if typ.BelongsToStruct != nil {
		secondGetSomethingParams = append(secondGetSomethingParams, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		secondArchiveSomethingParams = append(secondArchiveSomethingParams, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))

		getSomethingCountParams = append(getSomethingCountParams, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"))
		getListOfSomethingParams = append(getListOfSomethingParams, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"))
	}

	getSomethingParams = append(getSomethingParams, jen.List(secondGetSomethingParams...).ID("uint64"))
	archiveSomethingParams = append(archiveSomethingParams, jen.List(secondArchiveSomethingParams...).ID("uint64"))

	interfaceMethods = append(interfaceMethods,
		jen.IDf("Get%s", sn).Params(getSomethingParams...).Params(jen.Op("*").ID(sn), jen.ID("error")),
		jen.IDf("Get%sCount", sn).Params(getSomethingCountParams...).Params(jen.ID("uint64"), jen.ID("error")),
		jen.IDf("GetAll%sCount", pn).Params(getAllSomethingCount...).Params(jen.ID("uint64"), jen.ID("error")),
		jen.IDf("Get%s", pn).Params(getListOfSomethingParams...).Params(jen.Op("*").IDf("%sList", sn), jen.ID("error")),
	)

	if typ.BelongsToUser {
		interfaceMethods = append(interfaceMethods, jen.IDf("GetAll%sForUser", pn).Params(
			utils.CtxParam(),
			jen.ID("userID").ID("uint64"),
		).Params(jen.Index().ID(sn), jen.ID("error")))
	} else if typ.BelongsToStruct != nil {
		interfaceMethods = append(interfaceMethods, jen.IDf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular()).Params(
			utils.CtxParam(),
			jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).ID("uint64"),
		).Params(jen.Index().ID(sn), jen.ID("error")))
	}

	interfaceMethods = append(interfaceMethods,
		jen.IDf("Create%s", sn).Params(utils.CtxParam(), jen.ID("input").Op("*").IDf("%sCreationInput", sn)).Params(jen.Op("*").ID(sn), jen.ID("error")),
		jen.IDf("Update%s", sn).Params(utils.CtxParam(), jen.ID("updated").Op("*").ID(sn)).Params(jen.ID("error")),
		jen.IDf("Archive%s", sn).Params(archiveSomethingParams...).Params(jen.ID("error")),
	)

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
			jen.IDf("%sDataManager", sn).Interface(interfaceMethods...),
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

	buildToUpdateInput := func() jen.Code {
		lines := []jen.Code{}

		for _, typ := range typ.Fields {
			if typ.ValidForUpdateInput {
				fsn := typ.Name.Singular()
				lines = append(lines, jen.ID(fsn).Op(":").ID("x").Dot(fsn))
			}
		}

		return jen.Return(jen.Op("&").IDf("%sUpdateInput", sn).Valuesln(lines...))
	}

	ret.Add(
		jen.Commentf("ToInput creates a %sUpdateInput struct for %s", sn, cnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID(sn)).ID("ToInput").Params().Params(jen.Op("*").IDf("%sUpdateInput", sn)).Block(buildToUpdateInput()),
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
					jen.If(jen.ID("input").Dot(fsn).Op("!=").ID("nil").Op("&&").Op("*").ID("input").Dot(fsn).Op("!=").Lit("").Op("&&").ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
					),
				)
			} else {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).Op("!=").Lit("").Op("&&").ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
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
					jen.If(jen.ID("input").Dot(fsn).Op("!=").ID("nil").Op("&&").ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
					),
				)
			} else {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Block(
						jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
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
