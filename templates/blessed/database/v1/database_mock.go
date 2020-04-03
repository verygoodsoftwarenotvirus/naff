package v1

import (
	"fmt"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	mockImp = "github.com/stretchr/testify/mock"
)

func databaseMockDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("database")

	mockModelsImp := proj.ModelsV1Package("mock")
	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().ID("_").ID("Database").Equals().Parens(jen.PointerTo().ID("MockDatabase")).Call(jen.Nil()),
		jen.Line(),
	)

	buildMockDatabaseLines := func() []jen.Code {
		var lines []jen.Code

		for _, typ := range proj.DataTypes {
			lines = append(lines, jen.IDf("%sDataManager", typ.Name.Singular()).MapAssign().VarPointer().Qual(mockModelsImp, fmt.Sprintf("%sDataManager", typ.Name.Singular())).Values())
		}

		lines = append(lines,
			jen.ID("UserDataManager").MapAssign().VarPointer().Qual(mockModelsImp, "UserDataManager").Values(),
			jen.ID("OAuth2ClientDataManager").MapAssign().VarPointer().Qual(mockModelsImp, "OAuth2ClientDataManager").Values(),
			jen.ID("WebhookDataManager").MapAssign().VarPointer().Qual(mockModelsImp, "WebhookDataManager").Values(),
		)

		return lines
	}

	ret.Add(
		jen.Comment("BuildMockDatabase builds a mock database"),
		jen.Line(),
		jen.Func().ID("BuildMockDatabase").Params().Params(jen.PointerTo().ID("MockDatabase")).Block(
			jen.Return().VarPointer().ID("MockDatabase").Valuesln(
				buildMockDatabaseLines()...,
			),
		),
		jen.Line(),
	)

	buildMockDBLines := func() []jen.Code {
		lines := []jen.Code{
			jen.Qual(mockImp, "Mock"),
			jen.Line(),
		}

		for _, typ := range proj.DataTypes {
			lines = append(lines, jen.PointerTo().Qual(mockModelsImp, fmt.Sprintf("%sDataManager", typ.Name.Singular())))
		}

		lines = append(lines,
			jen.PointerTo().Qual(mockModelsImp, "UserDataManager"),
			jen.PointerTo().Qual(mockModelsImp, "OAuth2ClientDataManager"),
			jen.PointerTo().Qual(mockModelsImp, "WebhookDataManager"),
		)

		return lines
	}

	ret.Add(
		jen.Comment("MockDatabase is our mock database structure"),
		jen.Line(),
		jen.Type().ID("MockDatabase").Struct(
			buildMockDBLines()...,
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("Migrate satisfies the database.Database interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockDatabase")).ID("Migrate").Params(utils.CtxParam()).Params(jen.Error()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().ID("args").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("IsReady satisfies the database.Database interface"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").PointerTo().ID("MockDatabase")).ID("IsReady").Params(utils.CtxParam()).Params(jen.ID("ready").Bool()).Block(
			jen.ID("args").Assign().ID("m").Dot("Called").Call(utils.CtxVar()),
			jen.Return().ID("args").Dot("Bool").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
