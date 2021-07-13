package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableEventsDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	n := typ.Name

	code.Add(
		jen.Const().Defs(
			jen.Commentf("%sAssignmentKey is the key we use to indicate that an audit log entry is associated with an item.", n.Singular()),
			jen.IDf("%sAssignmentKey", n.Singular()).Equals().Litf("%s_id", n.RouteName()),
			jen.Commentf("%sCreationEvent is the event type used to indicate an item was created.", n.Singular()),
			jen.IDf("%sCreationEvent", n.Singular()).Equals().Litf("%s_created", n.RouteName()),
			jen.Commentf("%sUpdateEvent is the event type used to indicate an item was updated.", n.Singular()),
			jen.IDf("%sUpdateEvent", n.Singular()).Equals().Litf("%s_updated", n.RouteName()),
			jen.Commentf("%sArchiveEvent is the event type used to indicate an item was archived.", n.Singular()),
			jen.IDf("%sArchiveEvent", n.Singular()).Equals().Litf("%s_archived", n.RouteName()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Build%sCreationEventEntry builds an entry creation input for when %s is created.", n.Singular(), n.SingularCommonNameWithPrefix()),
		jen.Newline(),
		jen.Func().IDf("Build%sCreationEventEntry", n.Singular()).Params(
			jen.ID(n.UnexportedVarName()).PointerTo().Qual(proj.TypesPackage(), n.Singular()),
			jen.ID("createdByUser").Uint64(),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").IDf("%sCreationEvent", n.Singular()),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("createdByUser"),
					jen.IDf("%sAssignmentKey", n.Singular()).Op(":").ID(n.UnexportedVarName()).Dot("ID"),
					jen.ID("CreationAssignmentKey").Op(":").ID(n.UnexportedVarName()),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID("AccountAssignmentKey").Op(":").ID(n.UnexportedVarName()).Dot("BelongsToAccount")
						}
						return jen.Null()
					}(),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Build%sUpdateEventEntry builds an entry creation input for when %s is updated.", n.Singular(), n.SingularCommonNameWithPrefix()),
		jen.Newline(),
		jen.Func().IDf("Build%sUpdateEventEntry", n.Singular()).Params(
			jen.List(
				jen.ID("changedByUser"),
				jen.IDf("%sID", n.UnexportedVarName()),
				func() jen.Code {
					if typ.BelongsToAccount {
						return jen.ID("accountID")
					}
					return jen.Null()
				}(),
			).Uint64(),
			jen.ID("changes").Index().PointerTo().Qual(proj.TypesPackage(), "FieldChangeSummary"),
		).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").IDf("%sUpdateEvent", n.Singular()),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("changedByUser"),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID("AccountAssignmentKey").Op(":").ID("accountID")
						}
						return jen.Null()
					}(),
					jen.IDf("%sAssignmentKey", n.Singular()).Op(":").IDf("%sID", n.UnexportedVarName()),
					jen.ID("ChangesAssignmentKey").Op(":").ID("changes"),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("Build%sArchiveEventEntry builds an entry creation input for when %s is archived.", n.Singular(), n.SingularCommonNameWithPrefix()),
		jen.Newline(),
		jen.Func().IDf("Build%sArchiveEventEntry", n.Singular()).Params(
			jen.List(
				jen.ID("archivedByUser"),
				func() jen.Code {
					if typ.BelongsToAccount {
						return jen.ID("accountID")
					}
					return jen.Null()
				}(),
				jen.IDf("%sID", n.UnexportedVarName()),
			).ID("uint64")).Params(
			jen.PointerTo().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").IDf("%sArchiveEvent", n.Singular()),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("archivedByUser"),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID("AccountAssignmentKey").Op(":").ID("accountID")
						}
						return jen.Null()
					}(),
					jen.IDf("%sAssignmentKey", n.Singular()).Op(":").IDf("%sID", n.UnexportedVarName()),
				),
			),
		),
		jen.Newline(),
	)

	return code
}