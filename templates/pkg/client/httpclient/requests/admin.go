package requests

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("adminBasePath").Op("=").Lit("admin"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserReputationUpdateInputRequest builds a request to change a user's reputation."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildUserReputationUpdateInputRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserReputationUpdateInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("input").Dot("TargetUserID"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("adminBasePath"),
				jen.ID("usersBasePath"),
				jen.Lit("status"),
			),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
		),
		jen.Line(),
	)

	return code
}
