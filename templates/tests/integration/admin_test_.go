package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAdmin_Returns404WhenModifyingUserReputation").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should not be possible to ban a user that does not exist"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.ID("input").Op(":=").ID("fakes").Dot("BuildFakeUserReputationUpdateInput").Call(),
						jen.ID("input").Dot("TargetUserID").Op("=").ID("nonexistentID"),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("UpdateUserReputation").Call(
								jen.ID("ctx"),
								jen.ID("input"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("TestSuite")).ID("TestAdmin_BanningUsers").Params().Body(
			jen.ID("s").Dot("runForEachClientExcept").Call(
				jen.Lit("should be possible to ban users"),
				jen.Func().Params(jen.ID("testClients").Op("*").ID("testClientWrapper")).Params(jen.Func().Params()).Body(
					jen.Return().Func().Params().Body(
						jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
						jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartCustomSpan").Call(
							jen.ID("s").Dot("ctx"),
							jen.ID("t").Dot("Name").Call(),
						),
						jen.Defer().ID("span").Dot("End").Call(),
						jen.Var().Defs(
							jen.ID("user").Op("*").ID("types").Dot("User"),
							jen.ID("userClient").Op("*").ID("httpclient").Dot("Client"),
						),
						jen.Switch(jen.ID("testClients").Dot("authType")).Body(
							jen.Case(jen.ID("cookieAuthType")).Body(
								jen.List(jen.ID("user"), jen.ID("_"), jen.ID("userClient"), jen.ID("_")).Op("=").ID("createUserAndClientForTest").Call(
									jen.ID("ctx"),
									jen.ID("t"),
								)),
							jen.Case(jen.ID("pasetoAuthType")).Body(
								jen.List(jen.ID("user"), jen.ID("_"), jen.ID("_"), jen.ID("userClient")).Op("=").ID("createUserAndClientForTest").Call(
									jen.ID("ctx"),
									jen.ID("t"),
								)),
							jen.Default().Body(
								jen.Qual("log", "Panicf").Call(
									jen.Lit("invalid auth type: %q"),
									jen.ID("testClients").Dot("authType"),
								)),
						),
						jen.List(jen.ID("_"), jen.ID("initialCheckErr")).Op(":=").ID("userClient").Dot("GetAPIClients").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("initialCheckErr"),
						),
						jen.ID("input").Op(":=").Op("&").ID("types").Dot("UserReputationUpdateInput").Valuesln(jen.ID("TargetUserID").Op(":").ID("user").Dot("ID"), jen.ID("NewReputation").Op(":").ID("types").Dot("BannedUserAccountStatus"), jen.ID("Reason").Op(":").Lit("testing")),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("UpdateUserReputation").Call(
								jen.ID("ctx"),
								jen.ID("input"),
							),
						),
						jen.List(jen.ID("_"), jen.ID("subsequentCheckErr")).Op(":=").ID("userClient").Dot("GetAPIClients").Call(
							jen.ID("ctx"),
							jen.ID("nil"),
						),
						jen.ID("assert").Dot("Error").Call(
							jen.ID("t"),
							jen.ID("subsequentCheckErr"),
						),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("testClients").Dot("admin").Dot("ArchiveUser").Call(
								jen.ID("ctx"),
								jen.ID("user").Dot("ID"),
							),
						),
					)),
			)),
		jen.Line(),
	)

	return code
}
