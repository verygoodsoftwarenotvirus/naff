package workers

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preWritesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildPreWritesWorker(proj)...)
	code.Add(buildProvidePreWritesWorker(proj)...)
	code.Add(buildHandleMessage(proj)...)

	return code
}

func buildPreWritesWorker(proj *models.Project) []jen.Code {
	indexManagers := []jen.Code{}
	for _, typ := range proj.DataTypes {
		puvn := typ.Name.PluralUnexportedVarName()
		indexManagers = append(indexManagers, jen.IDf("%sIndexManager", puvn).Qual(proj.InternalSearchPackage(), "IndexManager"))
	}

	lines := []jen.Code{
		jen.Comment("PreWritesWorker writes data from the pending writes topic to the database."),
		jen.Newline(),
		jen.Type().ID("PreWritesWorker").Struct(
			append(
				[]jen.Code{
					jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
					jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
					jen.ID("encoder").Qual(proj.EncodingPackage(), "ClientEncoder"),
					jen.ID("postWritesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
					jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
				},
				indexManagers...,
			)...,
		),
		jen.Newline(),
	}

	return lines
}

func buildProvidePreWritesWorker(proj *models.Project) []jen.Code {
	bodyLines := []jen.Code{jen.Const().ID("name").Equals().Lit("pre_writes"),
		jen.Newline(),
	}
	valuesLines := []jen.Code{
		jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")).Dot("WithValue").Call(jen.Lit("topic"), jen.ID("name")),
		jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().Qual(proj.EncodingPackage(), "ProvideClientEncoder").Call(jen.ID("logger"), jen.Qual(proj.EncodingPackage(), "ContentTypeJSON")),
		jen.ID("postWritesPublisher").MapAssign().ID("postWritesPublisher"),
		jen.ID("dataManager").MapAssign().ID("dataManager"),
	}

	for _, typ := range proj.DataTypes {
		pcn := typ.Name.PluralCommonName()
		prn := typ.Name.PluralRouteName()
		puvn := typ.Name.PluralUnexportedVarName()

		stringFields := []jen.Code{}
		for _, field := range typ.Fields {
			if field.Type == "string" {
				stringFields = append(stringFields, jen.Lit(field.Name.UnexportedVarName()))
			}
		}

		bodyLines = append(bodyLines,
			jen.List(jen.IDf("%sIndexManager", puvn), jen.ID("err")).Assign().ID("searchIndexProvider").Call(
				append([]jen.Code{
					jen.ID("ctx"),
					jen.ID("logger"),
					jen.ID("client"),
					jen.ID("searchIndexLocation"),
					jen.Lit(prn),
				},
					stringFields...,
				)...,
			),
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
		jen.ID("w").Assign().Op("&").ID("PreWritesWorker").Valuesln(
			valuesLines...,
		),
		jen.Newline(),
		jen.Return().List(jen.ID("w"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Comment("ProvidePreWritesWorker provides a PreWritesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvidePreWritesWorker").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("client").Op("*").Qual("net/http", "Client"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("postWritesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("searchIndexLocation").Qual(proj.InternalSearchPackage(), "IndexPath"),
			jen.ID("searchIndexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider"),
		).Params(jen.Op("*").ID("PreWritesWorker"), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildHandleMessage(proj *models.Project) []jen.Code {
	switchCases := []jen.Code{}
	for _, typ := range proj.DataTypes {
		uvn := typ.Name.UnexportedVarName()
		puvn := typ.Name.PluralUnexportedVarName()
		sn := typ.Name.Singular()
		scn := typ.Name.SingularCommonName()

		switchCases = append(switchCases,
			jen.Case(jen.Qualf(proj.TypesPackage(), "%sDataType", sn)).Body(
				jen.List(jen.ID(uvn), jen.ID("err")).Assign().ID("w").Dot("dataManager").Dotf("Create%s", sn).Call(
					jen.ID("ctx"),
					jen.ID("msg").Dot(sn),
				), jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("creating %s", scn),
					)),
				jen.Newline(),
				jen.If(jen.ID("err").Equals().ID("w").Dotf("%sIndexManager", puvn).Dot("Index").Call(
					jen.ID("ctx"),
					jen.ID(uvn).Dot("ID"),
					jen.ID(uvn),
				), jen.ID("err").DoesNotEqual().Nil()).Body(
					jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Litf("indexing the %s", scn),
					)),
				jen.Newline(),
				jen.If(jen.ID("w").Dot("postWritesPublisher").DoesNotEqual().Nil()).Body(
					jen.ID("dcm").Assign().Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
						jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
						jen.ID(sn).MapAssign().ID(uvn),
						jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
						jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
					),
					jen.Newline(),
					jen.If(jen.ID("err").Equals().ID("w").Dot("postWritesPublisher").Dot("Publish").Call(jen.ID("ctx"), jen.ID("dcm")), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("publishing to post-writes topic"),
						),
					),
				),
			),
		)
	}

	switchCases = append(switchCases,
		jen.Case(jen.Qual(proj.TypesPackage(), "WebhookDataType")).Body(
			jen.List(jen.ID("webhook"), jen.ID("err")).Assign().ID("w").Dot("dataManager").Dot("CreateWebhook").Call(
				jen.ID("ctx"),
				jen.ID("msg").Dot("Webhook"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating webhook"),
				)),
			jen.Newline(),
			jen.If(jen.ID("w").Dot("postWritesPublisher").DoesNotEqual().Nil()).Body(
				jen.ID("dcm").Assign().Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
					jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
					jen.ID("Webhook").MapAssign().ID("webhook"),
					jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
					jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
				),
				jen.Newline(),
				jen.If(jen.ID("err").Equals().ID("w").Dot("postWritesPublisher").Dot("Publish").Call(
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
			jen.If(jen.ID("err").Assign().ID("w").Dot("dataManager").Dot("AddUserToAccount").Call(
				jen.ID("ctx"),
				jen.ID("msg").Dot("UserMembership"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating webhook"),
				)),
			jen.Newline(),
			jen.If(jen.ID("w").Dot("postWritesPublisher").DoesNotEqual().Nil()).Body(
				jen.ID("dcm").Assign().Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
					jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"), jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"), jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID")),
				jen.If(jen.ID("err").Assign().ID("w").Dot("postWritesPublisher").Dot("Publish").Call(
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

	lines := []jen.Code{
		jen.Comment("HandleMessage handles a pending write."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("PreWritesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Var().ID("msg").Op("*").Qual(proj.TypesPackage(), "PreWriteMessage"),
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
