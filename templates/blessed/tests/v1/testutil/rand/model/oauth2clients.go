package model

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Comment("RandomOAuth2ClientInput creates a random OAuth2ClientCreationInput"),
		jen.Line(),
		jen.Func().ID("RandomOAuth2ClientInput").Params(
			jen.List(jen.ID("username"), jen.ID("password"), jen.ID("totpToken")).ID("string")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientCreationInput")).Block(
			jen.ID("x").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientCreationInput").Valuesln(
				jen.ID("UserLoginInput").MapAssign().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID("username"),
					jen.ID("Password").MapAssign().ID("password"),
					jen.ID("TOTPToken").MapAssign().ID("mustBuildCode").Call(jen.ID("totpToken")),
				),
			),
			jen.Line(),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)
	return ret
}
