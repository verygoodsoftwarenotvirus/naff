package integration

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	sn := typ.Name.Singular()
	pn := typ.Name.Plural()

	ret.Add(
		jen.Func().IDf("check%sEquality", sn).Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Block(
			buildEqualityCheckLines(typ)...,
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("buildDummy%s", sn).Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
				buildFakeCallForCreationInput(pkg, typ)...,
			),
			jen.List(jen.ID("y"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("x")),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.Return().ID("y"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("Test%s", pn).Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be createable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestCreating(pkg, typ)...,
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestListing(pkg, typ)...,
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(typ)...,
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestReadingShouldBeReadable(pkg, typ)...,
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to update something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(pkg, typ)...,
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be updatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestUpdatingShouldBeUpdatable(pkg, typ)...,
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					buildTestDeletingShouldBeAbleToBeDeleted(pkg, typ)...,
				)),
			)),
		),
		jen.Line(),
	)

	return ret
}

func buildFakeCallForCreationInput(pkg *models.Project, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, field := range typ.Fields {
		lines = append(lines, jen.ID(field.Name.Singular()).Op(":").Add(utils.FakeCallForField(pkg.OutputPath, field)))
	}

	return lines
}

func fieldToExpectedDotField(varName string, typ models.DataType) []jen.Code {
	lines := []jen.Code{}

	for _, field := range typ.Fields {
		sn := field.Name.Singular()
		lines = append(lines, jen.ID(sn).Op(":").ID(varName).Dot(sn))
	}

	return lines
}

func buildEqualityCheckLines(typ models.DataType) []jen.Code {
	lines := []jen.Code{
		jen.ID("t").Dot("Helper").Call(),
		jen.Line(),
		jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ID")),
	}

	for _, field := range typ.Fields {
		sn := field.Name.Singular()
		if !field.Pointer {
			lines = append(lines, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(
				jen.ID("t"),
				jen.ID("expected").Dot(sn),
				jen.ID("actual").Dot(sn),
				jen.Lit("expected "+sn+" for ID %d to be %v, but it was %v "), jen.ID("expected").Dot("ID"), jen.ID("expected").Dot(sn), jen.ID("actual").Dot(sn),
			))
		} else {
			lines = append(lines, jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(
				jen.ID("t"),
				jen.Op("*").ID("expected").Dot(sn),
				jen.Op("*").ID("actual").Dot(sn),
				jen.Lit("expected "+sn+" to be %v, but it was %v "), jen.ID("expected").Dot(sn), jen.ID("actual").Dot(sn),
			))
		}
	}
	lines = append(lines, jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("CreatedOn")))

	return lines
}

func buildCreationArguments(pkg *models.Project, varPrefix string, typ models.DataType) []jen.Code {
	creationArgs := []jen.Code{}

	if typ.BelongsToStruct != nil {
		parentTyp := pkg.FindType(typ.BelongsToStruct.Singular())
		if parentTyp != nil {
			nca := buildCreationArguments(pkg, varPrefix, *parentTyp)
			creationArgs = append(creationArgs, nca...)
		}
	}

	creationArgs = append(creationArgs, jen.IDf("%s%s", varPrefix, typ.Name.Singular()).Dot("ID"))

	return creationArgs
}

func buildRequisiteCreationCode(pkg *models.Project, typ models.DataType) []jen.Code {
	var lines []jen.Code
	sn := typ.Name.Singular()

	const (
		sourceVarPrefix  = "example"
		createdVarPrefix = "created"
	)

	creationArgs := []jen.Code{
		jen.ID("ctx"),
	}
	ca := buildCreationArguments(pkg, createdVarPrefix, typ)
	creationArgs = append(creationArgs, ca[:len(ca)-1]...)
	creationArgs = append(creationArgs,
		jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
			fieldToExpectedDotField(fmt.Sprintf("%s%s", sourceVarPrefix, typ.Name.Singular()), typ)...,
		),
	)

	if typ.BelongsToStruct != nil {
		if parentTyp := pkg.FindType(typ.BelongsToStruct.Singular()); parentTyp != nil {
			newLines := buildRequisiteCreationCode(pkg, *parentTyp)
			lines = append(lines, newLines...)
		}
	}

	lines = append(lines,
		jen.Commentf("Create %s", typ.Name.SingularCommonName()),
		jen.IDf("%s%s", sourceVarPrefix, typ.Name.Singular()).Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), typ.Name.Singular()).Valuesln(
			buildFakeCallForCreationInput(pkg, typ)...,
		),
		jen.Line(),
		jen.List(jen.IDf("%s%s", createdVarPrefix, typ.Name.Singular()), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(creationArgs...),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.IDf("%s%s", createdVarPrefix, typ.Name.Singular()), jen.ID("err")),
		jen.Line(),
	)

	return lines
}

func buildTestCreating(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
	}

	lines = append(lines, buildRequisiteCreationCode(pkg, typ)...)

	lines = append(lines,
		jen.Commentf("Assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.IDf("created%s", typ.Name.Singular()), jen.ID("premade")),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Line(),
		jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.IDf("created%s", typ.Name.Singular()), jen.ID("actual")),
		jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ArchivedOn")),
	)

	return lines
}

func buildTestListing(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Commentf("Create %s", pcn),
		jen.Var().ID("expected").Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
		jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
			jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.IDf("buildDummy%s", sn).Call(jen.ID("t"))),
		),
		jen.Line(),
		jen.Commentf("Assert %s list equality", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", pn).Call(jen.ID("ctx"), jen.ID("nil")),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
		jen.Qual("github.com/stretchr/testify/assert", "True").Callln(
			jen.ID("t"),
			jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(pn)),
			jen.Lit("expected %d to be <= %d"), jen.ID("len").Call(jen.ID("expected")),
			jen.ID("len").Call(jen.ID("actual").Dot(pn)),
		),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("actual").Dot(pn)).Block(
			jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("x").Dot("ID")),
			jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
		),
	}

	return lines
}

func buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Commentf("Fetch %s", scn),
		jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("nonexistentID")),
		jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
	}

	return lines
}

func buildTestReadingShouldBeReadable(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Commentf("Create %s", scn),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
			buildFakeCallForCreationInput(pkg, typ)...,
		),
		jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.ID("ctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
			fieldToExpectedDotField("expected", typ)...,
		),
		),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
		jen.Line(),
		jen.Commentf("Fetch %s", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
		jen.Line(),
		jen.Commentf("Assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("actual").Dot("ID")),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
	}

	return lines
}

func buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.ID("err").Op(":=").ID("todoClient").Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(jen.ID("ID").Op(":").ID("nonexistentID"))),
		jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.ID("err")),
	}

	return lines
}

func buildTestUpdatingShouldBeUpdatable(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Commentf("Create %s", scn),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
			buildFakeCallForCreationInput(pkg, typ)...,
		),
		jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
			buildFakeCallForCreationInput(pkg, typ)...,
		)),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
		jen.Line(),
		jen.Commentf("Change %s", scn),
		jen.List(jen.ID("premade").Dot("Update").Call(jen.ID("expected").Dot("ToInput").Call())),
		jen.ID("err").Op("=").ID("todoClient").Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.ID("premade")),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
		jen.Line(),
		jen.Commentf("Fetch %s", scn),
		jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
		jen.Line(),
		jen.Commentf("Assert %s equality", scn),
		jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual").Dot("UpdatedOn")),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("actual").Dot("ID")),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
	}

	return lines
}

func buildTestDeletingShouldBeAbleToBeDeleted(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Line(),
		jen.Commentf("Create %s", scn),
		jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(
			buildFakeCallForCreationInput(pkg, typ)...,
		),
		jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.ID("ctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
			fieldToExpectedDotField("expected", typ)...,
		)),
		jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
		jen.Line(),
		jen.Comment("Clean up"),
		jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
		jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.ID("err")),
	}

	return lines
}
