package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func eventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("dataType").String(),
			jen.Newline(),
			jen.Comment("PreWriteMessage represents an event that asks a worker to write data to the datastore."),
			jen.ID("PreWriteMessage").Struct(
				jen.Underscore().Struct(),
				jen.Newline(),
				jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
			),
			jen.Newline(),
			jen.Comment("PreUpdateMessage represents an event that asks a worker to update data in the datastore."),
			jen.ID("PreUpdateMessage").Struct(
				jen.Underscore().Struct(),
				jen.Newline(),
				jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
			),
			jen.Newline(),
			jen.Comment("PreArchiveMessage represents an event that asks a worker to archive data in the datastore."),
			jen.ID("PreArchiveMessage").Struct(
				jen.Underscore().Struct(),
				jen.Newline(),
				jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
			),
			jen.Newline(),
			jen.Comment("DataChangeMessage represents an event that asks a worker to write data to the datastore."),
			jen.ID("DataChangeMessage").Struct(
				jen.Underscore().Struct(),
				jen.Newline(),
				jen.ID("DataType").ID("dataType").Tag(jsonTag("dataType")),
			),
		),
	)

	return code
}
