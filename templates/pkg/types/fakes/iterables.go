package fakes

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	fakeFields := []jen.Code{
		jen.ID("ID").MapAssign().Uint64().Call(jen.Qual(constants.FakeLibrary, "Uint32").Call()),
		jen.ID("ExternalID").MapAssign().Qual(constants.FakeLibrary, "UUID").Call(),
	}

	for _, field := range typ.Fields {
		fakeFields = append(fakeFields, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type, field.IsPointer)()))
	}

	fakeFields = append(fakeFields, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.ID("uint32").Call(jen.Qual(constants.FakeLibrary, "Date").Call().Dot("Unix").Call())))

	if typ.BelongsToStruct != nil {
		fakeFields = append(fakeFields, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Qual(constants.FakeLibrary, "Uint64").Call())
	}

	fakeFields = append(fakeFields,
		func() jen.Code {
			if typ.BelongsToAccount {
				return jen.ID("BelongsToAccount").MapAssign().Qual(constants.FakeLibrary, "Uint64").Call()
			}
			return jen.Null()
		}(),
	)

	code.Add(
		jen.Commentf("BuildFake%s builds a faked %s.", sn, scn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%s", sn).Params().Params(jen.Op("*").Qual(proj.TypesPackage(), sn)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), sn).Valuesln(fakeFields...),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sList builds a faked %sList.", sn, sn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sList", sn).Params().Params(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Body(
			jen.Var().ID("examples").Index().Op("*").Qual(proj.TypesPackage(), sn),
			jen.For(jen.ID("i").Op(":=").Lit(0),
				jen.ID("i").Op("<").ID("exampleQuantity"),
				jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.IDf("BuildFake%s", sn).Call(),
				)),
			jen.Newline(),
			jen.Return().Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.TypesPackage(), fmt.Sprintf("Pagination")).Valuesln(
					jen.ID("Page").MapAssign().Lit(1),
					jen.ID("Limit").MapAssign().Lit(20),
					jen.ID("FilteredCount").MapAssign().ID("exampleQuantity").Op("/").Lit(2),
					jen.ID("TotalCount").MapAssign().ID("exampleQuantity")),
				jen.ID(pn).MapAssign().ID("examples")),
		),
		jen.Newline(),
	)

	updateVals := []jen.Code{}
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			updateVals = append(updateVals, jen.ID(field.Name.Singular()).MapAssign().ID(uvn).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToStruct != nil {
		updateVals = append(updateVals, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID(uvn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToAccount {
		updateVals = append(updateVals, jen.ID("BelongsToAccount").MapAssign().ID(uvn).Dot("BelongsToAccount"))
	}

	code.Add(
		jen.Commentf("BuildFake%sUpdateInput builds a faked %sUpdateInput from %s.", sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sUpdateInput", sn).Params().Params(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))).Body(
			jen.ID(uvn).Op(":=").IDf("BuildFake%s", sn).Call(),
			jen.Return().Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Valuesln(updateVals...),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sUpdateInputFrom%s builds a faked %sUpdateInput from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sUpdateInputFrom%s", sn, sn).Params(jen.ID(uvn).Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn))).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateInput", sn)).Valuesln(updateVals...)),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sCreationInput builds a faked %sCreationInput.", sn, sn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sCreationInput", sn).Params().Params(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn))).Body(
			jen.ID(uvn).Op(":=").IDf("BuildFake%s", sn).Call(),
			jen.Return().IDf("BuildFake%sCreationInputFrom%s", sn, sn).Call(jen.ID(uvn)),
		),
		jen.Newline(),
	)

	creationVals := []jen.Code{}
	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			creationVals = append(creationVals, jen.ID(field.Name.Singular()).MapAssign().ID(uvn).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToStruct != nil {
		creationVals = append(creationVals, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID(uvn).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToAccount {
		creationVals = append(creationVals, jen.ID("BelongsToAccount").MapAssign().ID(uvn).Dot("BelongsToAccount"))
	}

	code.Add(
		jen.Commentf("BuildFake%sCreationInputFrom%s builds a faked %sCreationInput from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sCreationInputFrom%s", sn, sn).Params(jen.ID(uvn).Op("*").Qual(proj.TypesPackage(), sn)).Params(jen.Op("*").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn))).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationInput", sn)).Valuesln(creationVals...)),
		jen.Newline(),
	)

	return code
}
