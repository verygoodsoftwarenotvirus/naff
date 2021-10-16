package types

import (
	"fmt"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildSomethingConstantDefinitions(proj, typ)...)
	code.Add(buildInitFunc(typ)...)
	code.Add(buildSomethingTypeDefinitions(proj, typ)...)
	code.Add(buildSomethingUpdate(typ)...)
	code.Add(buildSomethingCreationInputValidateWithContext(typ)...)
	code.Add(buildSomethingDatabaseCreationInputValidateWithContext(typ)...)
	code.Add(buildSomethingDatabaseCreationInputFromSomething(typ)...)
	code.Add(buildSomethingUpdateInputValidateWithContext(typ)...)

	return code
}

func buildSomethingConstantDefinitions(_ *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	lines := []jen.Code{}

	if typ.SearchEnabled {
		lines = append(
			lines,
			jen.Const().Defs(
				jen.Commentf("%sDataType indicates an event is related to %s.", sn, scnwp),
				jen.IDf("%sDataType", sn).ID("dataType").Equals().Lit(typ.Name.RouteName()),
			),
			jen.Newline(),
		)
	}

	return lines
}

func buildInitFunc(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.Qual("encoding/gob", "Register").Call(jen.New(jen.ID(sn))),
			jen.Qual("encoding/gob", "Register").Call(jen.New(jen.IDf("%sList", sn))),
			jen.Qual("encoding/gob", "Register").Call(jen.New(jen.IDf("%sCreationRequestInput", sn))),
			jen.Qual("encoding/gob", "Register").Call(jen.New(jen.IDf("%sUpdateRequestInput", sn))),
		),
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
			jen.Underscore().Struct(),
			jen.Newline(),
			jen.ID("Pagination"),
			jen.ID(pn).Index().PointerTo().ID(sn).Tag(jsonTag(puvn)),
		),
	}

	lines = append(lines,
		jen.Newline(),
		jen.Commentf("%sCreationRequestInput represents what a user could set as input for creating %s.", sn, pcn),
		jen.IDf("%sCreationRequestInput", sn).Struct(buildCreateModelStructFields(typ)...),
		jen.Newline(),
		jen.Commentf("%sDatabaseCreationInput represents what a user could set as input for creating %s.", sn, pcn),
		jen.IDf("%sDatabaseCreationInput", sn).Struct(buildDatabaseCreationStructFields(typ)...),
		jen.Newline(),
		jen.Commentf("%sUpdateRequestInput represents what a user could set as input for updating %s.", sn, pcn),
		jen.IDf("%sUpdateRequestInput", sn).Struct(buildUpdateModelStructFields(typ)...),
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
			jen.ID("ListHandler").Params(
				jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"),
				jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
			),
			jen.ID("CreateHandler").Params(
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
func fakeEmptyValue(field models.DataField) jen.Code {
	switch field.Type {
	case "string":
		return jen.EmptyString()
	case "bool":
		return jen.False()
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		return jen.Zero()
	default:
		panic(fmt.Sprintf("unknown type! %q", field.Type))
	}
}

func buildSomethingUpdate(typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()

	updateLines := []jen.Code{}
	for i, field := range typ.Fields {
		fsn := field.Name.Singular()

		conditionBodyLines := []jen.Code{
			jen.ID("x").Dot(fsn).Equals().ID("input").Dot(fsn),
		}

		switch field.Type {
		case "string":
			if field.IsPointer {
				updateLines = append(updateLines,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Nil().And().Parens(jen.ID("x").Dot(fsn).IsEqualTo().Nil().Or().Parens(jen.PointerTo().ID("input").Dot(fsn).DoesNotEqual().EmptyString().And().PointerTo().ID("input").Dot(fsn).DoesNotEqual().PointerTo().ID("x").Dot(fsn)))).Body(
						conditionBodyLines...,
					),
					utils.ConditionalCode(i != len(typ.Fields)-1, jen.Newline()),
				)
			} else {
				updateLines = append(updateLines,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Lit("").Op("&&").ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
						conditionBodyLines...,
					),
					utils.ConditionalCode(i != len(typ.Fields)-1, jen.Newline()),
				)
			}
		case "bool":
			if field.IsPointer {
				updateLines = append(updateLines,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Nil().And().Parens(jen.ID("x").Dot(fsn).IsEqualTo().Nil().Or().Parens(jen.PointerTo().ID("input").Dot(fsn).DoesNotEqual().PointerTo().ID("x").Dot(fsn)))).Body(
						conditionBodyLines...,
					),
					utils.ConditionalCode(i != len(typ.Fields)-1, jen.Newline()),
				)
			} else {
				updateLines = append(updateLines,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
						conditionBodyLines...,
					),
					utils.ConditionalCode(i != len(typ.Fields)-1, jen.Newline()),
				)
			}
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			if field.IsPointer {
				updateLines = append(
					updateLines,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Nil().And().Parens(jen.ID("x").Dot(fsn).IsEqualTo().Nil().Or().Parens(jen.PointerTo().ID("input").Dot(fsn).DoesNotEqual().Zero().And().PointerTo().ID("input").Dot(fsn).DoesNotEqual().PointerTo().ID("x").Dot(fsn)))).Body(
						conditionBodyLines...,
					),
					utils.ConditionalCode(i != len(typ.Fields)-1, jen.Newline()),
				)
			} else {
				updateLines = append(
					updateLines,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Zero().Op("&&").ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
						conditionBodyLines...,
					),
					utils.ConditionalCode(i != len(typ.Fields)-1, jen.Newline()),
				)
			}
		default:
			panic(fmt.Sprintf("unknown type! %q", field.Type))
		}
	}

	lines := []jen.Code{
		jen.Commentf("Update merges an %sUpdateRequestInput with %s.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").PointerTo().ID(sn)).ID("Update").Params(jen.ID("input").PointerTo().IDf("%sUpdateRequestInput", sn)).Params().Body(
			updateLines...,
		),
	}

	return lines
}

func buildBaseModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("ID").String().Tag(jsonTag("id")),
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
		out = append(out, jen.ID(constants.UserOwnershipFieldName).String().Tag(jsonTag("belongsToAccount")))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).String().Tag(jsonTag(fmt.Sprintf("belongsTo%s", typ.BelongsToStruct.Singular()))))
	}

	return out
}

func buildUpdateModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
	}

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
		out = append(out, jen.ID(constants.UserOwnershipFieldName).String().Tag(jsonTag("-")))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).String().Tag(jsonTag(fmt.Sprintf("belongsTo%s", typ.BelongsToStruct.Singular()))))
	}

	return out
}

func buildCreateModelStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("ID").String().Tag(jsonTag("")),
	}

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
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).String().Tag(jsonTag("-")))
	}
	if typ.BelongsToAccount {
		out = append(out, jen.ID(constants.UserOwnershipFieldName).String().Tag(jsonTag("-")))
	}

	return out
}

func buildDatabaseCreationStructFields(typ models.DataType) []jen.Code {
	out := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("ID").String().Tag(jsonTag("id")),
	}

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
		out = append(out, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).String().Tag(jsonTag(fmt.Sprintf("belongsTo%s", typ.BelongsToStruct.Singular()))))
	}
	if typ.BelongsToAccount {
		out = append(out, jen.ID(constants.UserOwnershipFieldName).String().Tag(jsonTag("belongsToAccount")))
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
		jen.IDf("GetTotal%sCount", sn).Params(constants.CtxParam()).Params(jen.Uint64(), jen.Error()),
		jen.IDf("Get%s", pn).Params(typ.BuildInterfaceDefinitionListRetrievalMethodParams(proj)...).Params(jen.PointerTo().IDf("%sList", sn), jen.Error()),
		jen.IDf("Get%sWithIDs", pn).Params(getWithIDsParams...).Params(jen.Index().PointerTo().ID(sn), jen.Error()),
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
			if field.IsPointer {
				out = append(
					out,
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Nil().And().PointerTo().ID("input").Dot(fsn).DoesNotEqual().EmptyString().And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
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
					jen.If(jen.ID("input").Dot(fsn).DoesNotEqual().Nil().And().ID("input").Dot(fsn).DoesNotEqual().ID("x").Dot(fsn)).Body(
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
		if field.Type == "bool" {
			continue
		}
		validationFields = append(validationFields,
			jen.Qual(constants.ValidationLibrary, "Field").Call(
				jen.AddressOf().ID("x").Dot(field.Name.Singular()),
				jen.Qual(constants.ValidationLibrary, "Required"),
			),
		)
	}

	lines := []jen.Code{
		jen.Var().Underscore().Qual(constants.ValidationLibrary, "ValidatableWithContext").Equals().Parens(jen.PointerTo().IDf("%sCreationRequestInput", sn)).Call(jen.Nil()),
		jen.Newline(),
		jen.Newline(),
		jen.Commentf("ValidateWithContext validates a %sCreationRequestInput.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").PointerTo().IDf("%sCreationRequestInput", sn)).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual(constants.ValidationLibrary, "ValidateStructWithContext").Callln(
				validationFields...,
			)),
		jen.Newline(),
	}

	return lines
}

func buildSomethingDatabaseCreationInputValidateWithContext(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	validationFields := []jen.Code{
		jen.ID("ctx"),
		jen.ID("x"),
		jen.Qual(constants.ValidationLibrary, "Field").Call(
			jen.AddressOf().ID("x").Dot("ID"),
			jen.Qual(constants.ValidationLibrary, "Required"),
		),
	}
	for _, field := range typ.Fields {
		if field.Type == "bool" {
			continue
		}
		validationFields = append(validationFields,
			jen.Qual(constants.ValidationLibrary, "Field").Call(
				jen.AddressOf().ID("x").Dot(field.Name.Singular()),
				jen.Qual(constants.ValidationLibrary, "Required"),
			),
		)
	}

	if typ.BelongsToAccount {
		validationFields = append(validationFields, jen.Qual(constants.ValidationLibrary, "Field").Call(jen.AddressOf().ID("x").Dot(constants.AccountOwnershipFieldName), jen.Qual(constants.ValidationLibrary, "Required")))
	}

	lines := []jen.Code{
		jen.Var().Underscore().Qual(constants.ValidationLibrary, "ValidatableWithContext").Equals().Parens(jen.PointerTo().IDf("%sDatabaseCreationInput", sn)).Call(jen.Nil()),
		jen.Newline(),
		jen.Newline(),
		jen.Commentf("ValidateWithContext validates a %sDatabaseCreationInput.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").PointerTo().IDf("%sDatabaseCreationInput", sn)).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual(constants.ValidationLibrary, "ValidateStructWithContext").Callln(
				validationFields...,
			)),
		jen.Newline(),
	}

	return lines
}

func buildSomethingDatabaseCreationInputFromSomething(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.ID("x").Assign().AddressOf().IDf("%sDatabaseCreationInput", sn).Values(),
		jen.Newline(),
	}

	for _, field := range typ.Fields {
		lines = append(lines, jen.ID("x").Dot(field.Name.Singular()).Equals().ID("input").Dot(field.Name.Singular()))
	}

	lines = append(lines,
		jen.Newline(),
		jen.Return(jen.ID("x")),
	)

	return []jen.Code{
		jen.Commentf("%sDatabaseCreationInputFrom%sCreationInput creates a DatabaseCreationInput from a CreationInput.", sn, sn),
		jen.Newline(),
		jen.Func().IDf("%sDatabaseCreationInputFrom%sCreationInput", sn, sn).Params(jen.ID("input").PointerTo().IDf("%sCreationRequestInput", sn)).Params(jen.PointerTo().IDf("%sDatabaseCreationInput", sn)).Body(
			lines...,
		),
	}
}

func buildSomethingUpdateInputValidateWithContext(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	validationFields := []jen.Code{
		jen.ID("ctx"),
		jen.ID("x"),
	}
	for _, field := range typ.Fields {
		if field.Type == "bool" {
			continue
		}
		validationFields = append(validationFields,
			jen.Qual(constants.ValidationLibrary, "Field").Call(
				jen.AddressOf().ID("x").Dot(field.Name.Singular()),
				jen.Qual(constants.ValidationLibrary, "Required"),
			),
		)
	}

	lines := []jen.Code{
		jen.Var().Underscore().Qual(constants.ValidationLibrary, "ValidatableWithContext").Equals().Parens(jen.PointerTo().IDf("%sUpdateRequestInput", sn)).Call(jen.Nil()),
		jen.Newline(),
		jen.Newline(),
		jen.Commentf("ValidateWithContext validates a %sUpdateRequestInput.", sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("x").PointerTo().IDf("%sUpdateRequestInput", sn)).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual(constants.ValidationLibrary, "ValidateStructWithContext").Callln(
				validationFields...,
			)),
		jen.Newline(),
	}

	return lines
}
