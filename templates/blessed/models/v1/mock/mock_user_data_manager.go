package mock

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockUserDataManagerDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("mock")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Var().ID("_").Qual(filepath.Join(pkgRoot, "models/v1"), "UserDataManager").Op("=").Parens(jen.Op("*").ID("UserDataManager")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UserDataManager is a mocked models.UserDataManager for testing"),
		jen.Line(),
		jen.Type().ID("UserDataManager").Struct(jen.Qual("github.com/stretchr/testify/mock", "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("userID")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserByUsername is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("username")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUserCount is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("GetUserCount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter")).Params(jen.ID("uint64"), jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetUsers is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "QueryFilter")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "UserList"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("filter")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "UserList")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreateUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "UserInput")).Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User"),
			jen.ID("error")).Block(
			jen.ID("args").Op(":=").ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("input")),
			jen.Return().List(jen.ID("args").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")), jen.ID("args").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("UpdateUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), "User")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("updated")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ArchiveUser is a mock function"),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UserDataManager")).ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Block(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("ctx"), jen.ID("userID")).Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)
	return ret
}
