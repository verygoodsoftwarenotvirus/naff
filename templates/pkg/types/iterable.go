package types

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildSomethingConstantDefinitions(proj, typ)...)
	code.Add(buildSomethingTypeDefinitions(proj, typ)...)
	code.Add(buildUpdateSomething(typ)...)
	code.Add(buildSomethingToUpdateInput(typ)...)

	if typ.SearchEnabled && len(proj.FindOwnerTypeChain(typ)) > 0 {
		code.Add(buildSomethingToSearchHelper(proj, typ)...)
	}

	return code
}

func buildSomethingConstantDefinitions(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	pn := n.Plural()
	pcn := n.PluralCommonName()

	lines := []jen.Code{}

	if typ.SearchEnabled {
		lines = append(
			lines,
			jen.Const().Defs(
				jen.Commentf("%sSearchIndexName is the name of the index used to search through %s.", pn, pcn),
				jen.IDf("%sSearchIndexName", pn).Qual(proj.InternalSearchV1Package(), "IndexName").Equals().Lit(typ.Name.PluralRouteName()),
			),
			jen.Line(),
		)
	}

	return lines
}

func buildSomethingTypeDefinitions(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	cnwp := n.SingularCommonNameWithPrefix()
	pcn := n.PluralCommonName()
	puvn := n.PluralUnexportedVarName()

	lines := []jen.Code{
		jen.Commentf("%s represents %s.", sn, cnwp),
		jen.ID(sn).Struct(buildBaseModelStructFields(typ)...),
		jen.Line(),
		jen.Commentf("%sList represents a list of %s.", sn, pcn),
		jen.IDf("%sList", sn).Struct(
			jen.ID("Pagination"),
			jen.ID(pn).Index().ID(sn).Tag(jsonTag(puvn)),
		),
	}

	if typ.SearchEnabled && len(proj.FindOwnerTypeChain(typ)) > 0 {
		fields := buildBaseModelStructFields(typ)
		for _, o := range proj.FindOwnerTypeChain(typ) {
			fields = append(fields, jen.IDf("BelongsTo%s", o.Name.Singular()).Uint64().Tag(jsonTag(fmt.Sprintf("belongsTo%s", o.Name.Singular()))))
		}

		lines = append(lines,
			jen.Commentf("%sSearchHelper contains all the owner IDs for search purposes.", sn),
			jen.IDf("%sSearchHelper", sn).Struct(fields...),
			jen.Line(),
		)
	}

	lines = append(lines,
		jen.Line(),
		jen.Commentf("%sCreationInput represents what a user could set as input for creating %s.", sn, pcn),
		jen.IDf("%sCreationInput", sn).Struct(buildCreateModelStructFields(typ)...),
		jen.Line(),
		jen.Commentf("%sUpdateInput represents what a user could set as input for updating %s.", sn, pcn),
		jen.IDf("%sUpdateInput", sn).Struct(buildUpdateModelStructFields(typ)...),
		jen.Line(),
		jen.Commentf("%sDataManager describes a structure capable of storing %s permanently.", sn, pcn),
		jen.IDf("%sDataManager", sn).Interface(buildInterfaceMethods(proj, typ)...),
		jen.Line(),
		jen.Commentf("%sDataServer describes a structure capable of serving traffic related to %s.", sn, pcn),
		jen.IDf("%sDataServer", sn).Interface(
			jen.ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
			jen.ID("UpdateInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
			jen.Line(),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("SearchHandler").Params(
						jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
						jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
					)
				}
				return jen.Null()
			}(),
			jen.ID("ListHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
			jen.ID("CreateHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
			jen.ID("ExistenceHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
			jen.ID("ReadHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
			jen.ID("UpdateHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
			jen.ID("ArchiveHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
		),
	)

	return []jen.Code{
		jen.Type().Defs(
			lines...,
		),
	}
}

func buildUpdateSomething(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	cnwp := n.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("Update merges an %sInput with %s.", sn, cnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("x").PointerTo().ID(sn)).ID("Update").Params(jen.ID("input").PointerTo().IDf("%sUpdateInput", sn)).Body(buildUpdateFunctionLogic(typ.Fields)...),
		jen.Line(),
	}

	return lines
}

func buildSomethingToUpdateInput(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	cnwp := n.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("ToUpdateInput creates a %sUpdateInput struct for %s.", sn, cnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("x").PointerTo().ID(sn)).ID("ToUpdateInput").Params().Params(
			jen.PointerTo().IDf("%sUpdateInput", sn),
		).Body(
			func() jen.Code {
				lines := []jen.Code{}

				for _, typ := range typ.Fields {
					if typ.ValidForUpdateInput {
						fsn := typ.Name.Singular()
						lines = append(lines, jen.ID(fsn).MapAssign().ID("x").Dot(fsn))
					}
				}

				return jen.Return(jen.AddressOf().IDf("%sUpdateInput", sn).Valuesln(lines...))
			}(),
		),
		jen.Line(),
	}

	return lines
}

func buildSomethingToSearchHelper(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	cnwp := n.SingularCommonNameWithPrefix()

	owners := proj.FindOwnerTypeChain(typ)

	params := []jen.Code{}
	for _, o := range owners {
		params = append(params, jen.IDf("%sID", o.Name.UnexportedVarName()).Uint64())
	}

	lines := []jen.Code{
		jen.Commentf("ToSearchHelper creates a %sSearchHelper struct for %s.", sn, cnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("x").PointerTo().ID(sn)).ID("ToSearchHelper").Params(params...).ReturnParams(
			jen.PointerTo().IDf("%sSearchHelper", sn),
		).Body(
			func() jen.Code {
				lines := []jen.Code{}

				for _, typ := range typ.Fields {
					fsn := typ.Name.Singular()
					lines = append(lines, jen.ID(fsn).MapAssign().ID("x").Dot(fsn))
				}

				for _, o := range owners {
					lines = append(lines, jen.IDf("BelongsTo%s", o.Name.Singular()).MapAssign().IDf("%sID", o.Name.UnexportedVarName()))
				}

				return jen.Return(jen.AddressOf().IDf("%sSearchHelper", sn).Valuesln(lines...))
			}(),
		),
		jen.Line(),
	}

	return lines
}

func buildBaseModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{jen.ID("ID").Uint64().Tag(jsonTag("id"))}

	for _, field := range typ.Fields {
		if field.Pointer {
			out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
		} else {
			out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
		}
	}

	out = append(out,
		jen.ID("CreatedOn").Uint64().Tag(jsonTag("createdOn")),
		jen.ID("LastUpdatedOn").PointerTo().Uint64().Tag(jsonTag("lastUpdatedOn")),
		jen.ID("ArchivedOn").PointerTo().Uint64().Tag(jsonTag("archivedOn")),
	)

	if typ.BelongsToUser {
		out = append(out, jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("belongsToUser")))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag(fmt.Sprintf("belongsTo%s", typ.BelongsToStruct.Singular()))))
	}

	return out
}

func buildUpdateModelStructFields(typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			if field.Pointer {
				out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			}
		}
	}

	if typ.BelongsToUser {
		out = append(out, jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("-")))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag(fmt.Sprintf("belongsTo%s", typ.BelongsToStruct.Singular()))))
	}

	return out
}

func buildCreateModelStructFields(typ models.DataType) []jen.Code {
	var out []jen.Code

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			if field.Pointer {
				out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			}
		}
	}

	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag("-")))
	}
	if typ.BelongsToUser {
		out = append(out, jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("-")))
	}

	return out
}

func buildInterfaceMethods(proj *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	getWithIDsParams := typ.BuildGetListOfSomethingFromIDsParams(proj)

	interfaceMethods := []jen.Code{
		jen.IDf("%sExists", sn).Params(typ.BuildInterfaceDefinitionExistenceMethodParams(proj)...).Params(jen.Bool(), jen.Error()),
		jen.IDf("Get%s", sn).Params(typ.BuildInterfaceDefinitionRetrievalMethodParams(proj)...).Params(jen.PointerTo().ID(sn), jen.Error()),
		jen.IDf("GetAll%sCount", pn).Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()),
		jen.IDf("GetAll%s", pn).Params(
			constants.CtxParam(),
			jen.ID("resultChannel").Chan().Index().ID(sn),
		).Params(jen.Error()),
		jen.IDf("Get%s", pn).Params(typ.BuildInterfaceDefinitionListRetrievalMethodParams(proj)...).Params(jen.PointerTo().IDf("%sList", sn), jen.Error()),
		jen.IDf("Get%sWithIDs", pn).Params(getWithIDsParams...).Params(jen.Index().ID(sn), jen.Error()),
	}

	interfaceMethods = append(interfaceMethods,
		jen.IDf("Create%s", sn).Params(typ.BuildInterfaceDefinitionCreationMethodParams(proj)...).Params(jen.PointerTo().ID(sn), jen.Error()),
		jen.IDf("Update%s", sn).Params(typ.BuildInterfaceDefinitionUpdateMethodParams(proj, "updated")...).Params(jen.Error()),
		jen.IDf("Archive%s", sn).Params(typ.BuildInterfaceDefinitionArchiveMethodParams()...).Params(jen.Error()),
	)

	return interfaceMethods
}

func buildUpdateFunctionLogic(fields []models.DataField) []jen.Code {
	var out []jen.Code

	for i, field := range fields {
		fsn := field.Name.Singular()
		tt := strings.TrimPrefix(field.Type, "*")

		switch strings.ToLower(tt) {
		case "string":
			if field.Pointer {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("nil").And().PointerTo().ID("input").Dot(fsn).DoesNotEqual().EmptyString().And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
						jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
					),
				)
			} else {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().EmptyString().And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
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
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("nil").And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
						jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
					),
				)
			} else {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
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
