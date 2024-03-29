package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func dataChangesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("DataChangesWorker observes data changes in the database."),
		jen.Newline(),
		jen.Type().ID("DataChangesWorker").Struct(
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ClientEncoder"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideDataChangesWorker provides a DataChangesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvideDataChangesWorker").Params(jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger")).Params(jen.Op("*").ID("DataChangesWorker")).Body(
			jen.ID("name").Assign().Lit("post_writes"),
			jen.Newline(),
			jen.Return().Op("&").ID("DataChangesWorker").Valuesln(
				jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")), jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().Qual(proj.EncodingPackage(), "ProvideClientEncoder").Call(
					jen.ID("logger"),
					jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
				)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HandleMessage handles a pending write."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("DataChangesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Var().ID("msg").Op("*").Qual(proj.TypesPackage(), "DataChangeMessage"),
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
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("msg").Dot("AttributableToUserID"),
			),
			jen.ID("w").Dot("logger").Dot("WithValue").Call(
				jen.Lit("message"),
				jen.ID("message"),
			).Dot("Info").Call(jen.Lit("message received")),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	return code
}
