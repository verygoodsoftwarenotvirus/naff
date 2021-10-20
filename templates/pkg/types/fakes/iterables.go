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
		jen.ID("ID").MapAssign().Qual("github.com/segmentio/ksuid", "New").Call().Dot("String").Call(),
	}

	for _, field := range typ.Fields {
		fakeFields = append(fakeFields, jen.ID(field.Name.Singular()).MapAssign().Add(utils.FakeFuncForType(field.Type, field.IsPointer)()))
	}

	fakeFields = append(fakeFields, jen.ID("CreatedOn").MapAssign().Uint64().Call(jen.ID("uint32").Call(jen.Qual(constants.FakeLibrary, "Date").Call().Dot("Unix").Call())))

	if typ.BelongsToStruct != nil {
		fakeFields = append(fakeFields, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().Qual(constants.FakeLibrary, "UUID").Call())
	}

	fakeFields = append(fakeFields,
		func() jen.Code {
			if typ.BelongsToAccount {
				return jen.ID("BelongsToAccount").MapAssign().Qual(constants.FakeLibrary, "UUID").Call()
			}
			return jen.Null()
		}(),
	)

	code.Add(
		jen.Commentf("BuildFake%s builds a faked %s.", sn, scn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%s", sn).Params().Params(jen.PointerTo().Qual(proj.TypesPackage(), sn)).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), sn).Valuesln(fakeFields...),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sList builds a faked %sList.", sn, sn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sList", sn).Params().Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn))).Body(
			jen.Var().ID("examples").Index().PointerTo().Qual(proj.TypesPackage(), sn),
			jen.For(jen.ID("i").Assign().Zero(),
				jen.ID("i").Op("<").ID("exampleQuantity"),
				jen.ID("i").Op("++")).Body(
				jen.ID("examples").Equals().ID("append").Call(
					jen.ID("examples"),
					jen.IDf("BuildFake%s", sn).Call(),
				)),
			jen.Newline(),
			jen.Return().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.TypesPackage(), fmt.Sprintf("Pagination")).Valuesln(
					jen.ID("Page").MapAssign().One(),
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
		jen.Commentf("BuildFake%sUpdateRequestInput builds a faked %sUpdateRequestInput from %s.", sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sUpdateRequestInput", sn).Params().Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn))).Body(
			jen.ID(uvn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.Return().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn)).Valuesln(updateVals...),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sUpdateRequestInputFrom%s builds a faked %sUpdateRequestInput from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sUpdateRequestInputFrom%s", sn, sn).Params(jen.ID(uvn).PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn))).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sUpdateRequestInput", sn)).Valuesln(updateVals...)),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sCreationRequestInput builds a faked %sCreationRequestInput.", sn, sn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sCreationRequestInput", sn).Params().Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", sn))).Body(
			jen.ID(uvn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.Return().IDf("BuildFake%sCreationRequestInputFrom%s", sn, sn).Call(jen.ID(uvn)),
		),
		jen.Newline(),
	)

	creationVals := []jen.Code{
		jen.ID("ID").MapAssign().ID(uvn).Dot("ID"),
	}
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
		jen.Commentf("BuildFake%sCreationRequestInputFrom%s builds a faked %sCreationRequestInput from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sCreationRequestInputFrom%s", sn, sn).Params(jen.ID(uvn).PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", sn))).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sCreationRequestInput", sn)).Valuesln(creationVals...)),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sDatabaseCreationInput builds a faked %sDatabaseCreationInput.", sn, sn),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sDatabaseCreationInput", sn).Params().Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn))).Body(
			jen.ID(uvn).Assign().IDf("BuildFake%s", sn).Call(),
			jen.Return().IDf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn).Call(jen.ID(uvn)),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Commentf("BuildFake%sDatabaseCreationInputFrom%s builds a faked %sDatabaseCreationInput from %s.", sn, sn, sn, scnwp),
		jen.Newline(),
		jen.Func().IDf("BuildFake%sDatabaseCreationInputFrom%s", sn, sn).Params(jen.ID(uvn).PointerTo().Qual(proj.TypesPackage(), sn)).Params(jen.PointerTo().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn))).Body(
			jen.Return().AddressOf().Qual(proj.TypesPackage(), fmt.Sprintf("%sDatabaseCreationInput", sn)).Valuesln(creationVals...)),
		jen.Newline(),
	)

	return code
}
