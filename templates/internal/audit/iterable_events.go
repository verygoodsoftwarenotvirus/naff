package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableEventsDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	n := typ.Name

	code.Add(
		jen.Const().Defs(
			jen.IDf("%sAssignmentKey", n.Singular()).Op("=").Litf("%s_id", n.RouteName()),
			jen.IDf("%sCreationEvent", n.Singular()).Op("=").Litf("%s_created", n.RouteName()),
			jen.IDf("%sUpdateEvent", n.Singular()).Op("=").Litf("%s_updated", n.RouteName()),
			jen.IDf("%sArchiveEvent", n.Singular()).Op("=").Litf("%s_archived", n.RouteName()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Commentf("Build%sCreationEventEntry builds an entry creation input for when %s is created.", n.Singular(), n.SingularCommonNameWithPrefix()),
		jen.Line(),
		jen.Func().IDf("Build%sCreationEventEntry", n.Singular()).Params(
			jen.ID(n.UnexportedVarName()).Op("*").Qual(proj.TypesPackage(), n.Singular()),
			jen.ID("createdByUser").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").IDf("%sCreationEvent", n.Singular()),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("createdByUser"),
					jen.IDf("%sAssignmentKey", n.Singular()).Op(":").ID(n.UnexportedVarName()).Dot("ID"),
					jen.ID("CreationAssignmentKey").Op(":").ID(n.UnexportedVarName()),
					jen.ID("AccountAssignmentKey").Op(":").ID(n.UnexportedVarName()).Dot("BelongsToAccount"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Commentf("Build%sUpdateEventEntry builds an entry creation input for when %s is updated.", n.Singular(), n.SingularCommonNameWithPrefix()),
		jen.Line(),
		jen.Func().IDf("Build%sUpdateEventEntry", n.Singular()).Params(
			jen.List(jen.ID("changedByUser"), jen.IDf("%sID", n.UnexportedVarName()), jen.ID("accountID")).ID("uint64"),
			jen.ID("changes").Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").IDf("%sUpdateEvent", n.Singular()),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("changedByUser"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.IDf("%sAssignmentKey", n.Singular()).Op(":").IDf("%sID", n.UnexportedVarName()),
					jen.ID("ChangesAssignmentKey").Op(":").ID("changes"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Commentf("Build%sArchiveEventEntry builds an entry creation input for when %s is archived.", n.Singular(), n.SingularCommonNameWithPrefix()),
		jen.Line(),
		jen.Func().IDf("Build%sArchiveEventEntry", n.Singular()).Params(
			jen.List(jen.ID("archivedByUser"), jen.ID("accountID"), jen.IDf("%sID", n.UnexportedVarName())).ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").IDf("%sArchiveEvent", n.Singular()),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("archivedByUser"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.IDf("%sAssignmentKey", n.Singular()).Op(":").IDf("%sID", n.UnexportedVarName()),
				),
			),
		),
		jen.Line(),
	)

	return code
}
