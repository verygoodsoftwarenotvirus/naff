package fake

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)
	utils.AddImports(proj, code)

	code.Add(buildBuildFakeSomething(proj, typ)...)
	code.Add(buildBuildFakeSomethingList(proj, typ)...)
	code.Add(buildBuildFakeSomethingUpdateInputFromSomething(proj, typ)...)
	code.Add(buildBuildFakeSomethingCreationInput(proj, typ)...)
	code.Add(buildBuildFakeSomethingCreationInputFromSomething(proj, typ)...)

	return code
}

func buildBuildFakeSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	funcName := fmt.Sprintf("BuildFake%s", sn)

	block := []jen.Code{
		jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		if field.Type == "uint64" {
			if field.Pointer {
				block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.ValuePointerWrapper(field.Type, utils.FakeUint64WhichIsActuallyAUint32Func())()))
			} else {
				block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeUint64WhichIsActuallyAUint32Func()))
			}
		} else if field.Type == "int64" {
			if field.Pointer {
				block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.ValuePointerWrapper(field.Type, utils.FakeInt64WhichIsReallyAnInt32Func())()))
			} else {
				block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeInt64WhichIsReallyAnInt32Func()))
			}
		} else if field.Type == "float64" {
			if field.Pointer {
				block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.ValuePointerWrapper(field.Type, utils.FakeFloat64WhichIsActuallyAFloat32Func())()))
			} else {
				block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFloat64WhichIsActuallyAFloat32Func()))
			}
		} else {
			block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type, field.Pointer)()))
		}
	}

	block = append(block,
		jen.ID("CreatedOn").MapAssign().Add(utils.FakeUnixTimeFunc()),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(constants.UserOwnershipFieldName).MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s.", funcName, scn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn)).Block(
			jen.Return(jen.AddressOf().Qual(proj.ModelsV1Package(), sn).Valuesln(block...)),
		),
	}

	return lines
}

func buildBuildFakeSomethingList(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	funcName := fmt.Sprintf("BuildFake%sList", sn)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %sList.", funcName, sn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn))).Block(
			jen.IDf("example%s1", sn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.IDf("example%s2", sn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.IDf("example%s3", sn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.Line(),
			jen.Return(
				jen.AddressOf().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)).Valuesln(
					jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().One(),
						jen.ID("Limit").MapAssign().Lit(20),
					),
					jen.ID(pn).MapAssign().Index().Qual(proj.ModelsV1Package(), sn).Valuesln(
						jen.PointerTo().IDf("example%s1", sn),
						jen.PointerTo().IDf("example%s2", sn),
						jen.PointerTo().IDf("example%s3", sn),
					),
				),
			),
		),
	}

	return lines
}

func buildBuildFakeSomethingUpdateInputFromSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	uvn := typ.Name.UnexportedVarName()
	funcName := fmt.Sprintf("BuildFake%sUpdateInputFrom%s", sn, sn)

	var block []jen.Code
	for _, field := range typ.Fields {
		fns := field.Name.Singular()
		if field.ValidForUpdateInput {
			block = append(block, jen.ID(fns).MapAssign().ID(uvn).Dot(fns))
		}
	}

	block = append(block,
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				tns := typ.BelongsToStruct.Singular()
				return jen.IDf("BelongsTo%s", tns).MapAssign().ID(uvn).Dotf("BelongsTo%s", tns)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(constants.UserOwnershipFieldName).MapAssign().ID(uvn).Dot(constants.UserOwnershipFieldName)
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %sUpdateInput from %s.", funcName, sn, scnwp),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(uvn).PointerTo().Qual(proj.ModelsV1Package(), sn),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sUpdateInput", sn)),
		).Block(
			jen.Return(jen.AddressOf().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sUpdateInput", sn)).Valuesln(block...)),
		),
	}

	return lines
}

func buildBuildFakeSomethingCreationInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	funcName := fmt.Sprintf("BuildFake%sCreationInput", sn)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %sCreationInput.", funcName, sn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)),
		).Block(
			jen.ID(uvn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.Return(jen.IDf("BuildFake%sCreationInputFrom%s", sn, sn).Call(jen.ID(uvn))),
		),
	}

	return lines
}

func buildBuildFakeSomethingCreationInputFromSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	funcName := fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)

	var block []jen.Code
	for _, field := range typ.Fields {
		fns := field.Name.Singular()
		if field.ValidForCreationInput {
			block = append(block, jen.ID(fns).MapAssign().ID(uvn).Dot(fns))
		}
	}

	block = append(block,
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				tns := typ.BelongsToStruct.Singular()
				return jen.IDf("BelongsTo%s", tns).MapAssign().ID(uvn).Dotf("BelongsTo%s", tns)
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID(constants.UserOwnershipFieldName).MapAssign().ID(uvn).Dot(constants.UserOwnershipFieldName)
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %sCreationInput from %s.", funcName, sn, scnwp),
		jen.Line(),
		jen.Func().ID(funcName).Params(
			jen.ID(uvn).PointerTo().Qual(proj.ModelsV1Package(), sn),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)),
		).Block(
			jen.Return(jen.AddressOf().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sCreationInput", sn)).Valuesln(block...)),
		),
	}

	return lines
}
