package v1

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	mockImp = "github.com/stretchr/testify/mock"
)

func databaseMockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("database")

	utils.AddImports(proj, code)

	code.Add(
		utils.BuildInterfaceCheck("DataManager", "MockDatabase"),
		jen.Line(),
	)

	code.Add(buildBuildMockDatabase(proj)...)
	code.Add(buildMockDatabase(proj)...)
	code.Add(buildMigrate()...)
	code.Add(buildIsReady()...)
	code.Add(buildResultIterator()...)

	return code
}

func buildBuildMockDatabase(proj *models.Project) []jen.Code {
	mockModelsImp := proj.ModelsV1Package("mock")

	var mockDatabaseLines []jen.Code

	for _, typ := range proj.DataTypes {
		mockDatabaseLines = append(mockDatabaseLines, jen.IDf("%sDataManager", typ.Name.Singular()).MapAssign().AddressOf().Qual(mockModelsImp, fmt.Sprintf("%sDataManager", typ.Name.Singular())).Values())
	}

	mockDatabaseLines = append(mockDatabaseLines,
		jen.ID("UserDataManager").MapAssign().AddressOf().Qual(mockModelsImp, "UserDataManager").Values(),
		jen.ID("OAuth2ClientDataManager").MapAssign().AddressOf().Qual(mockModelsImp, "OAuth2ClientDataManager").Values(),
		jen.ID("WebhookDataManager").MapAssign().AddressOf().Qual(mockModelsImp, "WebhookDataManager").Values(),
	)

	lines := []jen.Code{
		jen.Comment("BuildMockDatabase builds a mock database."),
		jen.Line(),
		jen.Func().ID("BuildMockDatabase").Params().Params(jen.PointerTo().ID("MockDatabase")).Block(
			jen.Return().AddressOf().ID("MockDatabase").Valuesln(
				mockDatabaseLines...,
			),
		),
		jen.Line(),
	}

	return lines
}

func buildMockDatabase(proj *models.Project) []jen.Code {
	mockModelsImp := proj.ModelsV1Package("mock")

	mockDBLines := []jen.Code{
		jen.Qual(mockImp, "Mock"),
		jen.Line(),
	}

	for _, typ := range proj.DataTypes {
		mockDBLines = append(mockDBLines, jen.PointerTo().Qual(mockModelsImp, fmt.Sprintf("%sDataManager", typ.Name.Singular())))
	}

	mockDBLines = append(mockDBLines,
		jen.PointerTo().Qual(mockModelsImp, "UserDataManager"),
		jen.PointerTo().Qual(mockModelsImp, "OAuth2ClientDataManager"),
		jen.PointerTo().Qual(mockModelsImp, "WebhookDataManager"),
	)

	lines := []jen.Code{
		jen.Comment("MockDatabase is our mock database structure."),
		jen.Line(),
		jen.Type().ID("MockDatabase").Struct(
			mockDBLines...,
		),
		jen.Line(),
	}

	return lines
}

func buildMigrate() []jen.Code {
	lines := []jen.Code{
		jen.Comment("Migrate satisfies the DataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockDatabase")).ID("Migrate").Params(constants.CtxParam()).Params(jen.Error()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar()).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildIsReady() []jen.Code {
	lines := []jen.Code{
		jen.Comment("IsReady satisfies the DataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockDatabase")).ID("IsReady").Params(constants.CtxParam()).Params(jen.ID("ready").Bool()).Block(
			jen.Return().ID("m").Dot("Called").Call(constants.CtxVar()).Dot("Bool").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}

func buildResultIterator() []jen.Code {
	lines := []jen.Code{
		jen.Var().Underscore().ID("ResultIterator").Equals().Parens(jen.PointerTo().ID("MockResultIterator")).Call(jen.Nil()),
		jen.Line(),
		jen.Line(),
		jen.Comment("MockResultIterator is our mock sql.Rows structure."), jen.Line(),
		jen.Type().ID("MockResultIterator").Struct(
			jen.Qual(constants.MockPkg, "Mock"),
		),
		jen.Line(),
		jen.Comment("Scan satisfies the ResultIterator interface."), jen.Line(),
		jen.Func().Parens(jen.ID("m").PointerTo().ID("MockResultIterator")).ID("Scan").Params(jen.ID("dest").Spread().Interface()).Params(jen.Error()).Block(
			jen.Return(jen.ID("m").Dot("Called").Call(jen.ID("dest").Spread())).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
		jen.Comment("Next satisfies the ResultIterator interface."), jen.Line(),
		jen.Func().Parens(jen.ID("m").PointerTo().ID("MockResultIterator")).ID("Next").Params().Params(jen.Bool()).Block(
			jen.Return(jen.ID("m").Dot("Called").Call()).Dot("Bool").Call(jen.Zero()),
		),
		jen.Line(),
		jen.Comment("Err satisfies the ResultIterator interface."), jen.Line(),
		jen.Func().Parens(jen.ID("m").PointerTo().ID("MockResultIterator")).ID("Err").Params().Params(jen.Error()).Block(
			jen.Return(jen.ID("m").Dot("Called").Call()).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
		jen.Comment("Close satisfies the ResultIterator interface."), jen.Line(),
		jen.Func().Parens(jen.ID("m").PointerTo().ID("MockResultIterator")).ID("Close").Params().Params(jen.Error()).Block(
			jen.Return(jen.ID("m").Dot("Called").Call()).Dot("Error").Call(jen.Zero()),
		),
		jen.Line(),
	}

	return lines
}
