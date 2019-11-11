package client

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("dbclient")

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	ret.Add(
		jen.Func().IDf("TestClient_Get%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Get%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.IDf("example%sID", sn), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(jen.Qual("context", "Background").Call(), jen.IDf("example%sID", sn), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestClient_Get%sCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Get%sCount", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sCount", sn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Get%sCount", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sCount", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestClient_GetAll%sCount", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sCount", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAll%sCount", pn).Call(jen.Qual("context", "Background").Call()),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestClient_Get%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Values(),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Get%s", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call(), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Values(),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Get%s", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.ID("nil")), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("nil"), jen.ID("exampleUserID")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestClient_Create%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Create%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestClient_Update%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Update%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleInput")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("err").Op(":=").ID("c").Dotf("Update%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleInput")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("TestClient_Archive%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.Var().ID("expected").ID("error"),
				jen.Line(),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("Archive%s", sn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.IDf("example%sID", sn), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected")),
				jen.Line(),
				jen.ID("err").Op(":=").ID("c").Dotf("Archive%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUserID"), jen.IDf("example%sID", sn)),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)
	return ret
}
