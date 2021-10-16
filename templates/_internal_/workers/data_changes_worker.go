package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func dataChangesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("DataChangesWorker").Struct(
				jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("encoder").ID("encoding").Dot("ClientEncoder"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideDataChangesWorker provides a DataChangesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvideDataChangesWorker").Params(jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger")).Params(jen.Op("*").ID("DataChangesWorker")).Body(
			jen.ID("name").Op(":=").Lit("post_writes"),
			jen.Return().Op("&").ID("DataChangesWorker").Valuesln(
				jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")), jen.ID("tracer").MapAssign().ID("tracing").Dot("NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().ID("encoding").Dot("ProvideClientEncoder").Call(
					jen.ID("logger"),
					jen.ID("encoding").Dot("ContentTypeJSON"),
				)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HandleMessage handles a pending write."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("DataChangesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().Defs(
				jen.ID("msg").Op("*").ID("types").Dot("DataChangeMessage"),
			),
			jen.If(jen.ID("err").Op(":=").ID("w").Dot("encoder").Dot("Unmarshal").Call(
				jen.ID("ctx"),
				jen.ID("message"),
				jen.Op("&").ID("msg"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("w").Dot("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling message"),
				)),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("msg").Dot("AttributableToUserID"),
			),
			jen.ID("w").Dot("logger").Dot("WithValue").Call(
				jen.Lit("message"),
				jen.ID("message"),
			).Dot("Info").Call(jen.Lit("message received")),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	return code
}
