package types

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildPreWriteFields(proj *models.Project) []jen.Code {
	fields := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
	}

	for _, typ := range proj.DataTypes {
		if len(proj.FindDependentsOfType(typ)) > 0 {
			fields = append(fields, jen.IDf("%sID", typ.Name.Singular()).String().Tag(jsonTag(typ.Name.UnexportedVarName())))
		}

		fields = append(fields, jen.ID(typ.Name.Singular()).PointerTo().IDf("%sDatabaseCreationInput", typ.Name.Singular()).Tag(jsonTagOmittingEmpty(typ.Name.UnexportedVarName())))
	}

	fields = append(fields,
		jen.ID("Webhook").PointerTo().ID("WebhookDatabaseCreationInput").Tag(jsonTagOmittingEmpty("webhook")),
		jen.ID("UserMembership").PointerTo().ID("AddUserToAccountInput").Tag(jsonTagOmittingEmpty("userMembership")),
		jen.ID("AttributableToUserID").String().Tag(jsonTag("attributableToUserID")),
		jen.ID("AttributableToAccountID").String().Tag(jsonTag("attributableToAccountID")),
	)

	return fields
}

func buildPreUpdateFields(proj *models.Project) []jen.Code {
	fields := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
	}

	for _, typ := range proj.DataTypes {
		if len(proj.FindDependentsOfType(typ)) > 0 {
			fields = append(fields, jen.IDf("%sID", typ.Name.Singular()).String().Tag(jsonTag(typ.Name.UnexportedVarName())))
		}

		fields = append(fields, jen.ID(typ.Name.Singular()).PointerTo().ID(typ.Name.Singular()).Tag(jsonTagOmittingEmpty(typ.Name.UnexportedVarName())))
	}

	fields = append(fields,
		jen.ID("AttributableToUserID").String().Tag(jsonTag("attributableToUserID")),
		jen.ID("AttributableToAccountID").String().Tag(jsonTag("attributableToAccountID")),
	)

	return fields
}

func buildPreArchiveFields(proj *models.Project) []jen.Code {
	fields := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
	}

	for _, typ := range proj.DataTypes {
		fields = append(fields, jen.IDf("%sID", typ.Name.Singular()).String().Tag(jsonTag(fmt.Sprintf("%sID", typ.Name.UnexportedVarName()))))
	}

	fields = append(fields,
		jen.ID("WebhookID").String().Tag(jsonTag("webhookID")),
		jen.ID("AttributableToUserID").String().Tag(jsonTag("attributableToUserID")),
		jen.ID("AttributableToAccountID").String().Tag(jsonTag("attributableToAccountID")),
	)

	return fields
}

func buildDataChangeMessageFields(proj *models.Project) []jen.Code {
	fields := []jen.Code{
		jen.Underscore().Struct(),
		jen.Newline(),
		jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
		jen.ID("MessageType").String().Tag(jsonTag("messageType")),
	}

	for _, typ := range proj.DataTypes {
		if len(proj.FindDependentsOfType(typ)) > 0 {
			fields = append(fields, jen.IDf("%sID", typ.Name.Singular()).String().Tag(jsonTag(typ.Name.UnexportedVarName())))
		}

		fields = append(fields, jen.ID(typ.Name.Singular()).PointerTo().ID(typ.Name.Singular()).Tag(jsonTagOmittingEmpty(typ.Name.UnexportedVarName())))
	}

	fields = append(fields,
		jen.ID("Webhook").PointerTo().ID("Webhook").Tag(jsonTagOmittingEmpty("webhook")),
		jen.ID("UserMembership").PointerTo().ID("AccountUserMembership").Tag(jsonTagOmittingEmpty("userMembership")),
		jen.ID("Context").Map(jen.String()).String().Tag(jsonTagOmittingEmpty("context")),
		jen.ID("AttributableToUserID").String().Tag(jsonTag("attributableToUserID")),
		jen.ID("AttributableToAccountID").String().Tag(jsonTag("attributableToAccountID")),
	)

	return fields
}

func eventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("dataType").String(),
			jen.Newline(),
			jen.Comment("PreWriteMessage represents an event that asks a worker to write data to the datastore."),
			jen.ID("PreWriteMessage").Struct(
				buildPreWriteFields(proj)...,
			),
			jen.Newline(),
			jen.Comment("PreUpdateMessage represents an event that asks a worker to update data in the datastore."),
			jen.ID("PreUpdateMessage").Struct(
				buildPreUpdateFields(proj)...,
			),
			jen.Newline(),
			jen.Comment("PreArchiveMessage represents an event that asks a worker to archive data in the datastore."),
			jen.ID("PreArchiveMessage").Struct(
				buildPreArchiveFields(proj)...,
			),
			jen.Newline(),
			jen.Comment("DataChangeMessage represents an event that asks a worker to write data to the datastore."),
			jen.ID("DataChangeMessage").Struct(
				buildDataChangeMessageFields(proj)...,
			),
		),
	)

	return code
}
