package workers

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preArchivesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildPreArchivesWorker(proj)...)
	code.Add(buildProvidePreArchivesWorker(proj)...)
	code.Add(buildHandlePreArchivesWorkerMessage(proj)...)

	return code
}

func buildPreArchivesWorker(proj *models.Project) []jen.Code {
	indexManagers := []jen.Code{}
	for _, typ := range proj.DataTypes {
		puvn := typ.Name.PluralUnexportedVarName()
		indexManagers = append(indexManagers, jen.IDf("%sIndexManager", puvn).Qual(proj.InternalSearchPackage(), "IndexManager"))
	}

	lines := []jen.Code{
		jen.Comment("PreArchivesWorker archives data from the pending archives topic to the database."),
		jen.Newline(),
		jen.Type().ID("PreArchivesWorker").Struct(
			append(
				[]jen.Code{
					jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
					jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
					jen.ID("encoder").Qual(proj.EncodingPackage(), "ClientEncoder"),
					jen.ID("postArchivesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
					jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
				},
				indexManagers...,
			)...,
		),
		jen.Newline(),
	}

	return lines
}
func buildProvidePreArchivesWorker(proj *models.Project) []jen.Code {
	bodyLines := []jen.Code{
		jen.Const().ID("name").Equals().Lit("pre_archives"),
		jen.Newline(),
	}
	valuesLines := []jen.Code{
		jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")).Dot("WithValue").Call(jen.Lit("topic"), jen.ID("name")),
		jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().Qual(proj.EncodingPackage(), "ProvideClientEncoder").Call(jen.ID("logger"), jen.Qual(proj.EncodingPackage(), "ContentTypeJSON")),
		jen.ID("postArchivesPublisher").MapAssign().ID("postArchivesPublisher"),
		jen.ID("dataManager").MapAssign().ID("dataManager"),
	}

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()
		prn := typ.Name.PluralRouteName()
		puvn := typ.Name.PluralUnexportedVarName()

		providerLines := []jen.Code{
			jen.ID("ctx"),
			jen.ID("logger"),
			jen.ID("client"),
			jen.ID("searchIndexLocation"),
			jen.Lit(prn),
		}

		for _, field := range typ.Fields {
			if field.Type == "string" {
				providerLines = append(providerLines, jen.Lit(field.Name.UnexportedVarName()))
			}
		}

		bodyLines = append(bodyLines,
			jen.List(jen.IDf("%sIndexManager", puvn), jen.ID("err")).Assign().ID("searchIndexProvider").Call(providerLines...),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
					jen.Lit(fmt.Sprintf("setting up %s search index manager", pcn)+": %w"),
					jen.ID("err"),
				))),
			jen.Newline(),
		)

		valuesLines = append(valuesLines, jen.IDf("%sIndexManager", puvn).MapAssign().IDf("%sIndexManager", puvn))
	}

	bodyLines = append(bodyLines,
		jen.ID("w").Assign().Op("&").ID("PreArchivesWorker").Valuesln(
			valuesLines...,
		),
		jen.Newline(),
		jen.Return().List(jen.ID("w"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Comment("ProvidePreArchivesWorker provides a PreArchivesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvidePreArchivesWorker").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("client").Op("*").Qual("net/http", "Client"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("postArchivesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("searchIndexLocation").Qual(proj.InternalSearchPackage(), "IndexPath"),
			jen.ID("searchIndexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider"),
		).Params(jen.Op("*").ID("PreArchivesWorker"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildHandlePreArchivesWorkerMessage(proj *models.Project) []jen.Code {
	switchCases := []jen.Code{}
	for _, typ := range proj.DataTypes {
		puvn := typ.Name.PluralUnexportedVarName()
		sn := typ.Name.Singular()
		scn := typ.Name.SingularCommonName()

		switchCases = append(switchCases, jen.Case(jen.Qualf(proj.TypesPackage(), "%sDataType", sn)).Body(
			jen.If(jen.ID("err").Assign().ID("w").Dot("dataManager").Dotf("Archive%s", sn).Call(
				jen.ID("ctx"),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.ID("msg").Dotf("%sID", typ.BelongsToStruct.Singular())
					}
					return jen.Null()
				}(),
				jen.ID("msg").Dotf("%sID", sn),
				utils.ConditionalCode(typ.BelongsToAccount, jen.ID("msg").Dot("AttributableToAccountID")),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("w").Dot("logger"),
					jen.ID("span"),
					jen.Litf("archiving %s", scn),
				)),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("w").Dotf("%sIndexManager", puvn).Dot("Delete").Call(
				jen.ID("ctx"),
				jen.ID("msg").Dotf("%sID", sn),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("w").Dot("logger"),
					jen.ID("span"),
					jen.Litf("removing %s from index", scn),
				)),
			jen.Newline(),
			jen.If(jen.ID("w").Dot("postArchivesPublisher").DoesNotEqual().Nil()).Body(
				jen.ID("dcm").Assign().Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
					jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
					jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
					jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
				),
				jen.Newline(),
				jen.If(jen.ID("err").Assign().ID("w").Dot("postArchivesPublisher").Dot("Publish").Call(
					jen.ID("ctx"),
					jen.ID("dcm"),
				), jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("publishing data change message"),
					)),
			)),
		)
	}

	switchCases = append(switchCases, jen.Case(jen.Qual(proj.TypesPackage(), "WebhookDataType")).Body(
		jen.If(jen.ID("err").Assign().ID("w").Dot("dataManager").Dot("ArchiveWebhook").Call(
			jen.ID("ctx"),
			jen.ID("msg").Dot("WebhookID"),
			jen.ID("msg").Dot("AttributableToAccountID"),
		), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("w").Dot("logger"),
				jen.ID("span"),
				jen.Lit("creating webhook"),
			)),
		jen.Newline(),
		jen.If(jen.ID("w").Dot("postArchivesPublisher").DoesNotEqual().Nil()).Body(
			jen.ID("dcm").Assign().Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
				jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
				jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
				jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
			),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("w").Dot("postArchivesPublisher").Dot("Publish").Call(
				jen.ID("ctx"),
				jen.ID("dcm"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("publishing data change message"),
				)),
		)),
		jen.Case(jen.Qual(proj.TypesPackage(), "UserMembershipDataType")).Body(
			jen.Break()),
	)

	lines := []jen.Code{
		jen.Comment("HandleMessage handles a pending archive."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("PreArchivesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Var().ID("msg").Op("*").Qual(proj.TypesPackage(), "PreArchiveMessage"),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("w").Dot("encoder").Dot("Unmarshal").Call(
				jen.ID("ctx"),
				jen.ID("message"),
				jen.Op("&").ID("msg"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("w").Dot("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling message"),
				)),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("msg").Dot("AttributableToUserID"),
			),
			jen.ID("logger").Assign().ID("w").Dot("logger").Dot("WithValue").Call(
				jen.Lit("data_type"),
				jen.ID("msg").Dot("DataType"),
			),
			jen.Newline(),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("message read")),
			jen.Newline(),
			jen.Switch(jen.ID("msg").Dot("DataType")).Body(
				switchCases...,
			),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	}

	return lines
}
