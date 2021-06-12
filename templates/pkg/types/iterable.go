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

	utils.AddImports(proj, code, false)

	code.Add(buildSomethingConstantDefinitions(proj, typ)...)
	code.Add(buildSomethingTypeDefinitions(proj, typ)...)
	code.Add(buildSomethingToUpdateInput(typ)...)
	code.Add(buildSomethingCreationInputValidateWithContext(typ)...)
	code.Add(buildSomethingUpdateInputValidateWithContext(typ)...)

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
				jen.IDf("%sSearchIndexName", pn).Qual(proj.InternalSearchPackage(), "IndexName").Equals().Lit(typ.Name.PluralRouteName()),
			),
			jen.Newline(),
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
		jen.Newline(),
		jen.Commentf("%sList represents a list of %s.", sn, pcn),
		jen.IDf("%sList", sn).Struct(
			jen.ID("Pagination"),
			jen.ID(pn).Index().PointerTo().ID(sn).Tag(jsonTag(puvn)),
		),
	}

	lines = append(lines,
		jen.Newline(),
		jen.Commentf("%sCreationInput represents what a user could set as input for creating %s.", sn, pcn),
		jen.IDf("%sCreationInput", sn).Struct(buildCreateModelStructFields(typ)...),
		jen.Newline(),
		jen.Commentf("%sUpdateInput represents what a user could set as input for updating %s.", sn, pcn),
		jen.IDf("%sUpdateInput", sn).Struct(buildUpdateModelStructFields(typ)...),
		jen.Newline(),
		jen.Commentf("%sDataManager describes a structure capable of storing %s permanently.", sn, pcn),
		jen.IDf("%sDataManager", sn).Interface(buildInterfaceMethods(proj, typ)...),
		jen.Newline(),
		jen.Commentf("%sDataService describes a structure capable of serving traffic related to %s.", sn, pcn),
		jen.IDf("%sDataService", sn).Interface(
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("SearchHandler").Params(
						jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
						jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
					)
				}
				return jen.Null()
			}(),
			jen.ID("AuditEntryHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
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

func buildSomethingToUpdateInput(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	updateLines := []jen.Code{
		jen.Var().ID("out").Index().Op("*").ID("FieldChangeSummary"),
		jen.Newline(),
	}
	for _, field := range typ.Fields {
		fsn := field.Name.Singular()

		updateLines = append(updateLines,
			jen.If(jen.ID("input").Dot(fsn).Op("!=").Lit("").Op("&&").ID("input").Dot(fsn).Op("!=").ID("x").Dot(fsn)).Body(
				jen.ID("out").Op("=").ID("append").Call(
					jen.ID("out"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit(fsn), jen.ID("OldValue").Op(":").ID("x").Dot(fsn), jen.ID("NewValue").Op(":").ID("input").Dot(fsn)),
				),
				jen.Newline(),
				jen.ID("x").Dot(fsn).Op("=").ID("input").Dot(fsn),
			),
			jen.Newline(),
		)
	}
	updateLines = append(updateLines,
		jen.Newline(),
		jen.Return().ID("out"),
	)

	lines := []jen.Code{
		jen.Commentf("Update merges an %sUpdateInput with %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").Op("*").ID(sn)).ID("Update").Params(jen.ID("input").Op("*").IDf("%sUpdateInput", sn)).Params(jen.Index().Op("*").ID("FieldChangeSummary")).Body(
			updateLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildBaseModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{
		jen.ID("ID").Uint64().Tag(jsonTag("id")),
		jen.ID("ExternalID").String().Tag(jsonTag("externalID")),
	}

	for _, field := range typ.Fields {
		if field.IsPointer {
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

	if typ.BelongsToAccount {
		out = append(out, jen.ID(constants.UserOwnershipFieldName).Uint64().Tag(jsonTag("belongsToAccount")))
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
			if field.IsPointer {
				out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			}
		}
	}

	if typ.BelongsToAccount {
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
			if field.IsPointer {
				out = append(out, jen.ID(field.Name.Singular()).PointerTo().ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			} else {
				out = append(out, jen.ID(field.Name.Singular()).ID(field.Type).Tag(jsonTag(field.Name.UnexportedVarName())))
			}
		}
	}

	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).Uint64().Tag(jsonTag("-")))
	}
	if typ.BelongsToAccount {
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
			jen.ID("resultChannel").Chan().Index().PointerTo().ID(sn),
			jen.ID("bucketSize").Uint16(),
		).Params(jen.Error()),
		jen.IDf("Get%s", pn).Params(typ.BuildInterfaceDefinitionListRetrievalMethodParams(proj)...).Params(jen.PointerTo().IDf("%sList", sn), jen.Error()),
		jen.IDf("Get%sWithIDs", pn).Params(getWithIDsParams...).Params(jen.Index().PointerTo().ID(sn), jen.Error()),
	}

	interfaceMethods = append(interfaceMethods,
		jen.IDf("Create%s", sn).Params(typ.BuildInterfaceDefinitionCreationMethodParams(proj)...).Params(jen.PointerTo().ID(sn), jen.Error()),
		jen.IDf("Update%s", sn).Params(typ.BuildInterfaceDefinitionUpdateMethodParams(proj, "updated")...).Params(jen.Error()),
		jen.IDf("Archive%s", sn).Params(typ.BuildInterfaceDefinitionArchiveMethodParams()...).Params(jen.Error()),
		jen.IDf("GetAuditLogEntriesFor%s", sn).Params(typ.BuildInterfaceDefinitionAuditLogEntryRetrievalMethodParams()...).Params(jen.Index().PointerTo().ID("AuditLogEntry"), jen.Error()),
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
			if field.IsPointer {
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
			if field.IsPointer {
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
			out = append(out, jen.Newline())
		}
	}

	return out
}

func buildSomethingCreationInputValidateWithContext(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	validationFields := []jen.Code{
		jen.ID("ctx"),
		jen.ID("x"),
	}
	for _, field := range typ.Fields {
		validationFields = append(validationFields,
			jen.Qual(constants.ValidationLibrary, "Field").Call(
				jen.Op("&").ID("x").Dot(field.Name.Singular()),
				jen.Qual(constants.ValidationLibrary, "Required"),
			),
		)
	}

	lines := []jen.Code{
		jen.Var().ID("_").Qual(constants.ValidationLibrary, "ValidatableWithContext").Op("=").Parens(jen.Op("*").IDf("%sCreationInput", sn)).Call(jen.ID("nil")),
		jen.Newline(),
		jen.Newline(),
		jen.Commentf("ValidateWithContext validates a %sCreationInput.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").Op("*").IDf("%sCreationInput", sn)).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual(constants.ValidationLibrary, "ValidateStructWithContext").Callln(
				validationFields...,
			)),
		jen.Newline(),
	}

	return lines
}

func buildSomethingUpdateInputValidateWithContext(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	validationFields := []jen.Code{
		jen.ID("ctx"),
		jen.ID("x"),
	}
	for _, field := range typ.Fields {
		validationFields = append(validationFields,
			jen.Qual(constants.ValidationLibrary, "Field").Call(
				jen.Op("&").ID("x").Dot(field.Name.Singular()),
				jen.Qual(constants.ValidationLibrary, "Required"),
			),
		)
	}

	lines := []jen.Code{
		jen.Var().ID("_").Qual(constants.ValidationLibrary, "ValidatableWithContext").Op("=").Parens(jen.Op("*").IDf("%sUpdateInput", sn)).Call(jen.ID("nil")),
		jen.Newline(),
		jen.Newline(),
		jen.Commentf("ValidateWithContext validates a %sUpdateInput.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").Op("*").IDf("%sUpdateInput", sn)).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual(constants.ValidationLibrary, "ValidateStructWithContext").Callln(
				validationFields...,
			)),
		jen.Newline(),
	}

	return lines
}
