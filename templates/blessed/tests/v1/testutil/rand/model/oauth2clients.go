package model

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsDotGo() *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(ret)

	ret.Add(
		jen.Comment("RandomOAuth2ClientInput creates a random OAuth2ClientCreationInput"),
		jen.Line(),
		jen.Func().ID("RandomOAuth2ClientInput").Params(
			jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpToken")).ID("string")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2ClientCreationInput")).Block(
			jen.ID("x").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "OAuth2ClientCreationInput").Valuesln(
				jen.ID("UserLoginInput").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("username"),
					jen.ID("Password").Op(":").ID("password"),
					jen.ID("TOTPToken").Op(":").ID("mustBuildCode").Call(jen.ID("totpToken")),
				),
			),
			jen.Line(),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)
	return ret
}