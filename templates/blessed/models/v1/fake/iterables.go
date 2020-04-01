package fake

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"path/filepath"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile(packageName)
	utils.AddImports(proj, ret)

	ret.Add(buildBuildFakeSomething(proj, typ)...)
	ret.Add(buildBuildFakeItemList(proj, typ)...)
	ret.Add(buildBuildFakeItemUpdateInputFromItem(proj, typ)...)
	ret.Add(buildBuildFakeItemCreationInput(proj, typ)...)
	ret.Add(buildBuildFakeItemCreationInputFromItem(proj, typ)...)

	return ret
}

//
//// BuildFakeItem builds a faked item
//func BuildFakeItem() *models.Item {
//	return &models.Item{
//		ID:            fake.Uint64(),
//		Name:          fake.Word(),
//		Details:       fake.Word(),
//		CreatedOn:     uint64(uint32(fake.Date().Unix())),
//		BelongsToUser: fake.Uint64(),
//	}
//}
//

func buildBuildFakeSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	funcName := fmt.Sprintf("BuildFake%s", sn)

	block := []jen.Code{
		jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
	}

	block = append(block,
		jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Uint32().Call(jen.Qual(utils.FakeLibrary, "Date").Call().Dot("Unix").Call())),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, scn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(jen.ParamPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn)).Block(
			jen.Return(jen.VarPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn).Valuesln(block...)),
		),
	}

	return lines
}

//
//// BuildFakeItemList builds a faked ItemList
//func BuildFakeItemList() *models.ItemList {
//	exampleItem1 := BuildFakeItem()
//	exampleItem2 := BuildFakeItem()
//	exampleItem3 := BuildFakeItem()
//
//	return &models.ItemList{
//		Pagination: models.Pagination{
//			Page:       1,
//			Limit:      20,
//			TotalCount: 3,
//		},
//		Items: []models.Item{
//			*exampleItem1,
//			*exampleItem2,
//			*exampleItem3,
//		},
//	}
//}
//

func buildBuildFakeItemList(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	funcName := fmt.Sprintf("BuildFake%sList", sn)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %sList", funcName, sn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Params(jen.ParamPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), fmt.Sprintf("%sList", sn))).Block(
			jen.IDf("example%s1", sn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.IDf("example%s2", sn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.IDf("example%s3", sn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.Line(),
			jen.Return(
				jen.VarPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), fmt.Sprintf("%sList", sn)).Valuesln(
					jen.ID("Pagination").MapAssign().Qual(filepath.Join(proj.OutputPath, "models", "v1"), "Pagination").Valuesln(
						jen.ID("Page").MapAssign().Lit(1),
						jen.ID("Limit").MapAssign().Lit(20),
						jen.ID("TotalCount").MapAssign().Lit(3),
					),
					jen.ID(pn).MapAssign().Index().Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn).Valuesln(
						jen.PointerTo().ID("exampleItem1"),
						jen.PointerTo().ID("exampleItem2"),
						jen.PointerTo().ID("exampleItem3"),
					),
				),
			),
		),
	}

	return lines
}

//
//// BuildFakeItemUpdateInputFromItem builds a faked ItemUpdateInput from an item
//func BuildFakeItemUpdateInputFromItem(item *models.Item) *models.ItemUpdateInput {
//	return &models.ItemUpdateInput{
//		Name:          item.Name,
//		Details:       item.Details,
//		BelongsToUser: item.BelongsToUser,
//	}
//}
//

func buildBuildFakeItemUpdateInputFromItem(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.Singular()
	funcName := fmt.Sprintf("BuildFakeItemUpdateInputFrom%s", sn)

	block := []jen.Code{
		jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
	}

	block = append(block,
		jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Uint32().Call(jen.Qual(utils.FakeLibrary, "Date").Call().Dot("Unix").Call())),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, scn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Block(
			jen.Return(jen.VarPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn).Valuesln(block...)),
		),
	}

	return lines
}

//
//// BuildFakeItemCreationInput builds a faked ItemCreationInput
//func BuildFakeItemCreationInput() *models.ItemCreationInput {
//	item := BuildFakeItem()
//	return BuildFakeItemCreationInputFromItem(item)
//}
//

func buildBuildFakeItemCreationInput(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.Singular()
	funcName := fmt.Sprintf("BuildFake%sCreationInput", sn)

	block := []jen.Code{
		jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
	}

	block = append(block,
		jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Uint32().Call(jen.Qual(utils.FakeLibrary, "Date").Call().Dot("Unix").Call())),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, scn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Block(
			jen.Return(jen.VarPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn).Valuesln(block...)),
		),
	}

	return lines
}

//
//// BuildFakeItemCreationInputFromItem builds a faked ItemCreationInput from an item
//func BuildFakeItemCreationInputFromItem(item *models.Item) *models.ItemCreationInput {
//	return &models.ItemCreationInput{
//		Name:          item.Name,
//		Details:       item.Details,
//		BelongsToUser: item.BelongsToUser,
//	}
//}
//

func buildBuildFakeItemCreationInputFromItem(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.Singular()
	funcName := fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)

	block := []jen.Code{
		jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
	}

	for _, field := range typ.Fields {
		block = append(block, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type)()))
	}

	block = append(block,
		jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.Uint32().Call(jen.Qual(utils.FakeLibrary, "Date").Call().Dot("Unix").Call())),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
		func() jen.Code {
			if typ.BelongsToUser {
				return jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func())
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Commentf("%s builds a faked %s", funcName, scn),
		jen.Line(),
		jen.Func().ID(funcName).Params().Block(
			jen.Return(jen.VarPointer().Qual(filepath.Join(proj.OutputPath, "models", "v1"), sn).Valuesln(block...)),
		),
	}

	return lines
}
