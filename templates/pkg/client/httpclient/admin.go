package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("UpdateUserReputation updates a user's reputation."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateUserReputation").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserReputationUpdateInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("input").Dot("TargetUserID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("TargetUserID"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				)),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildUserReputationUpdateInputRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user reputation update request"),
				)),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("authedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating user reputation"),
				)),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.If(jen.ID("err").Op("=").ID("errorFromResponse").Call(jen.ID("res")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("invalid response status"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
