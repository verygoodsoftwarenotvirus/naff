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

	ret.Add()

	ret.Add(buildTestClientGetSomething(pkg, typ)...)
	ret.Add(buildTestClientGetSomethingCount(pkg, typ)...)
	ret.Add(buildTestClientGetAllOfSomethingCount(pkg, typ)...)
	ret.Add(buildTestClientGetListOfSomething(pkg, typ)...)

	if typ.BelongsToUser {
		ret.Add(buildTestClientGetAllOfSomethingForUser(pkg, typ)...)
	}

	ret.Add(buildTestClientCreateSomething(pkg, typ)...)
	ret.Add(buildTestClientUpdateSomething(pkg, typ)...)
	ret.Add(buildTestClientArchiveSomething(pkg, typ)...)

	return ret
}

func buildTestClientGetSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	block := []jen.Code{
		jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
	}
	mockCallArgs := []jen.Code{
		jen.Litf("Get%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
		jen.IDf("example%sID", sn),
	}
	callArgs := []jen.Code{
		jen.Qual("context", "Background").Call(),
		jen.IDf("example%sID", sn),
	}

	if typ.BelongsToUser {
		block = append(block, jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
		mockCallArgs = append(mockCallArgs, jen.ID("exampleUserID"))
		callArgs = append(callArgs, jen.ID("exampleUserID"))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
		mockCallArgs = append(mockCallArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
		callArgs = append(callArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	}

	block = append(block,
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", sn).Call(callArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		jen.Line(),
		jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
	)

	return []jen.Code{jen.Func().IDf("TestClient_Get%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
	),
		jen.Line(),
	}
}

func buildTestClientGetSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	buildSubtest := func(typ models.DataType, nilFilter bool) []jen.Code {
		lines := []jen.Code{
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
		}

		mockCalls := []jen.Code{
			jen.Litf("Get%sCount", sn),
			jen.Qual("github.com/stretchr/testify/mock", "Anything"),
		}

		callArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
		}

		if !nilFilter {
			mockCalls = append(mockCalls, jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call())
			callArgs = append(callArgs, jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call())
		} else {
			mockCalls = append(mockCalls, jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.ID("nil")))
			callArgs = append(callArgs, jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.ID("nil")))
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockCalls = append(mockCalls, jen.ID("exampleUserID"))
			callArgs = append(callArgs, jen.ID("exampleUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockCalls = append(mockCalls, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
			callArgs = append(callArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
		}

		lines = append(lines,
			jen.Line(),
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCalls...).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%sCount", sn).Call(callArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%sCount", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSubtest(typ, false)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSubtest(typ, true)...)),
		),
		jen.Line(),
	}
}

func buildTestClientGetAllOfSomethingForUser(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	return []jen.Code{
		jen.Func().IDf("TestClient_GetAll%sForUser", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("expected").Assign().Slice().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
				jen.Line(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sForUser", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("exampleUserID")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAll%sForUser", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("exampleUserID")),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	}
}

func buildTestClientGetAllOfSomethingCount(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	return []jen.Code{
		jen.Func().IDf("TestClient_GetAll%sCount", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(321)),
				jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
				jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(jen.Litf("GetAll%sCount", pn), jen.Qual("github.com/stretchr/testify/mock", "Anything")).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("GetAll%sCount", pn).Call(jen.Qual("context", "Background").Call()),
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.Line(),
				jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
			)),
		),
		jen.Line(),
	}
}

func buildTestClientGetListOfSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()

	buildSubtest := func(nilFilter bool) []jen.Code {
		mockCalls := []jen.Code{
			jen.Litf("Get%s", pn),
			jen.Qual("github.com/stretchr/testify/mock", "Anything"),
		}
		callArgs := []jen.Code{
			jen.Qual("context", "Background").Call(),
		}
		lines := []jen.Code{}

		if !nilFilter {
			mockCalls = append(mockCalls, jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call())
			callArgs = append(callArgs, jen.Qual(filepath.Join(pkg.OutputPath, "models/v1"), "DefaultQueryFilter").Call())
		} else {
			mockCalls = append(mockCalls, jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.ID("nil")))
			callArgs = append(callArgs, jen.Parens(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter")).Call(jen.ID("nil")))
		}

		if typ.BelongsToUser {
			lines = append(lines, jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockCalls = append(mockCalls, jen.ID("exampleUserID"))
			callArgs = append(callArgs, jen.ID("exampleUserID"))
		} else if typ.BelongsToStruct != nil {
			lines = append(lines, jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
			mockCalls = append(mockCalls, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
			callArgs = append(callArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
		}

		lines = append(lines,
			jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
			jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Values(),
			jen.Line(),
			jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCalls...).Dot("Return").Call(jen.ID("expected"), jen.ID("nil")),
			jen.Line(),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(callArgs...),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		)

		return lines
	}

	return []jen.Code{
		jen.Func().IDf("TestClient_Get%s", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSubtest(false)...)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with nil filter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(buildSubtest(true)...)),
		),
		jen.Line(),
	}
}

func buildTestClientCreateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	return []jen.Code{jen.Func().IDf("TestClient_Create%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
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
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			jen.Line(),
			jen.ID("mockDB").Dot("AssertExpectations").Call(jen.ID("t")),
		)),
	),
		jen.Line(),
	}
}

func buildTestClientUpdateSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	return []jen.Code{
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
				jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	}
}

func buildTestClientArchiveSomething(pkg *models.Project, typ models.DataType) []jen.Code {
	n := typ.Name
	sn := n.Singular()

	block := []jen.Code{}
	callArgs := []jen.Code{
		jen.Qual("context", "Background").Call(),
	}
	mockCallArgs := []jen.Code{
		jen.Litf("Archive%s", sn),
		jen.Qual("github.com/stretchr/testify/mock", "Anything"),
		jen.IDf("example%sID", sn),
	}

	if typ.BelongsToUser {
		block = append(block, jen.ID("exampleUserID").Op(":=").ID("uint64").Call(jen.Lit(123)))
		callArgs = append(callArgs, jen.ID("exampleUserID"))
		mockCallArgs = append(mockCallArgs, jen.ID("exampleUserID"))
	} else if typ.BelongsToStruct != nil {
		block = append(block, jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").ID("uint64").Call(jen.Lit(123)))
		callArgs = append(callArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
		mockCallArgs = append(mockCallArgs, jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	}

	callArgs = append(callArgs, jen.IDf("example%sID", sn))

	block = append(block,
		jen.IDf("example%sID", sn).Op(":=").ID("uint64").Call(jen.Lit(123)),
		jen.Var().ID("expected").ID("error"),
		jen.Line(),
		jen.List(jen.ID("c"), jen.ID("mockDB")).Op(":=").ID("buildTestClient").Call(),
		jen.ID("mockDB").Dotf("%sDataManager", sn).Dot("On").Call(mockCallArgs...).Dot("Return").Call(jen.ID("expected")),
		jen.Line(),
		jen.ID("err").Op(":=").ID("c").Dotf("Archive%s", sn).Call(callArgs...),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
	)

	return []jen.Code{
		jen.Func().IDf("TestClient_Archive%s", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(block...)),
		),
		jen.Line(),
	}
}
