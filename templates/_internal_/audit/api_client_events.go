package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientEventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("APIClientAssignmentKey").Op("=").Lit("api_client_id"),
			jen.ID("APIClientCreationEvent").Op("=").Lit("api_client_created"),
			jen.ID("APIClientUpdateEvent").Op("=").Lit("api_client_created"),
			jen.ID("APIClientArchiveEvent").Op("=").Lit("api_client_archived"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAPIClientCreationEventEntry builds an entry creation input for when an API client is created."),
		jen.Line(),
		jen.Func().ID("BuildAPIClientCreationEventEntry").Params(
			jen.ID("client").Op("*").Qual(proj.TypesPackage(), "APIClient"),
			jen.ID("createdBy").ID("uint64"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("APIClientCreationEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("APIClientAssignmentKey").Op(":").ID("client").Dot("ID"),
					jen.ID("CreationAssignmentKey").Op(":").ID("client"),
					jen.ID("ActorAssignmentKey").Op(":").ID("createdBy"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAPIClientArchiveEventEntry builds an entry creation input for when an API client is archived."),
		jen.Line(),
		jen.Func().ID("BuildAPIClientArchiveEventEntry").Params(
			jen.List(jen.ID("accountID"), jen.ID("clientID"), jen.ID("archivedBy")).ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("APIClientArchiveEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("archivedBy"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.ID("APIClientAssignmentKey").Op(":").ID("clientID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
