package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("cfg").Assign().AddressOf().ID("Config").Values(),
		jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
		jen.ID("authService").Assign().AddressOf().Qual(proj.TypesPackage("mock"), "AuthService").Values(),
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("s").Assign().ID("ProvideService").Call(
			jen.ID("cfg"),
			jen.ID("logger"),
			jen.ID("authService"),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("authService"),
		),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(jen.ID("t"), jen.ID("s")),
	)

	code.Add(
		jen.Func().ID("TestProvideService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	)

	return code
}
