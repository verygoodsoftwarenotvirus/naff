package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("AdminService").Interface(jen.ID("UserReputationChangeHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request"))),
			jen.ID("AdminAuditManager").Interface(
				jen.ID("LogUserBanEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("banGiver"), jen.ID("banReceiver")).ID("uint64"), jen.ID("reason").ID("string")),
				jen.ID("LogAccountTerminationEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("terminator"), jen.ID("terminee")).ID("uint64"), jen.ID("reason").ID("string")),
			),
			jen.ID("UserReputationUpdateInput").Struct(
				jen.ID("NewReputation").ID("accountStatus"),
				jen.ID("Reason").ID("string"),
				jen.ID("TargetUserID").ID("uint64"),
			),
			jen.ID("FrontendService").Interface(jen.ID("StaticDir").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("staticFilesDirectory").ID("string")).Params(jen.Qual("net/http", "HandlerFunc"), jen.ID("error"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("UserReputationUpdateInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our struct is validatable."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("UserReputationUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("NewReputation"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("Reason"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("TargetUserID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Min").Call(jen.ID("uint64").Call(jen.Lit(1))),
				),
			)),
		jen.Line(),
	)

	return code
}
