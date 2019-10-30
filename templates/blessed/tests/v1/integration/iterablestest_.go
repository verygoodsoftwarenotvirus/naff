package integration

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(rootPkg string, typ models.DataType) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(ret)

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	ret.Add(
		jen.Func().IDf("check%sEquality", sn).Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").Qual(filepath.Join(rootPkg, "models/v1"), sn)).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("assert").Dot("NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ID")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Name"), jen.ID("actual").Dot("Name")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Details"), jen.ID("actual").Dot("Details")),
			jen.ID("assert").Dot("NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("CreatedOn")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().IDf("buildDummy%s", sn).Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual(filepath.Join(rootPkg, "models/v1"), sn)).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(rootPkg, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
				jen.ID("Name").Op(":").ID("fake").Dot("Word").Call(),
				jen.ID("Details").Op(":").ID("fake").Dot("Sentence").Call(),
			),
			jen.List(jen.ID("y"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.Qual("context", "Background").Call(), jen.ID("x")),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
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
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Commentf("Create %s", scn),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(rootPkg, "models/v1"), sn).Valuesln(
						jen.ID("Name").Op(":").Lit("name"),
						jen.ID("Details").Op(":").Lit("details"),
					),
					jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(
						jen.ID("ctx"),
						jen.Op("&").Qual(filepath.Join(rootPkg, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
							jen.ID("Name").Op(":").ID("expected").Dot("Name"),
							jen.ID("Details").Op(":").ID("expected").Dot("Details"),
						),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
					jen.Line(),
					jen.Commentf("Assert %s equality", scn),
					jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID("expected"), jen.ID("premade")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
					jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
					jen.Line(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
					jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.ID("assert").Dot("NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ArchivedOn")),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Commentf("Create %s", pcn),
					jen.Var().ID("expected").Index().Op("*").Qual(filepath.Join(rootPkg, "models/v1"), sn),
					jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
						jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.IDf("buildDummy%s", sn).Call(jen.ID("t"))),
					),
					jen.Line(),
					jen.Commentf("Assert %s list equality", scn),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", pn).Call(jen.ID("ctx"), jen.ID("nil")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
					jen.ID("assert").Dot("True").Callln(
						jen.ID("t"),
						jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot(pn)),
						jen.Lit("expected %d to be <= %d"), jen.ID("len").Call(jen.ID("expected")),
						jen.ID("len").Call(jen.ID("actual").Dot(pn)),
					),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("actual").Dot(pn)).Block(
						jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("x").Dot("ID")),
						jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
					),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Commentf("Fetch %s", scn),
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("nonexistentID")),
					jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Commentf("Create %s", scn),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(rootPkg, "models/v1"), sn).Valuesln(
						jen.ID("Name").Op(":").Lit("name"),
						jen.ID("Details").Op(":").Lit("details"),
					),
					jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.ID("ctx"), jen.Op("&").Qual(filepath.Join(rootPkg, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
						jen.ID("Name").Op(":").ID("expected").Dot("Name"),
						jen.ID("Details").Op(":").ID("expected").Dot("Details")),
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
					jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to update something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.ID("err").Op(":=").ID("todoClient").Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.Op("&").Qual(filepath.Join(rootPkg, "models/v1"), sn).Values(jen.ID("ID").Op(":").ID("nonexistentID"))),
					jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be updatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Commentf("Create %s", scn),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(rootPkg, "models/v1"), sn).Valuesln(
						jen.ID("Name").Op(":").Lit("new name"),
						jen.ID("Details").Op(":").Lit("new details"),
					),
					jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(rootPkg, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
						jen.ID("Name").Op(":").Lit("old name"),
						jen.ID("Details").Op(":").Lit("old details"),
					)),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
					jen.Line(),
					jen.Commentf("Change %s", scn),
					jen.List(jen.ID("premade").Dot("Name"), jen.ID("premade").Dot("Details")).Op("=").List(jen.ID("expected").Dot("Name"), jen.ID("expected").Dot("Details")),
					jen.ID("err").Op("=").ID("todoClient").Dotf("Update%s", sn).Call(jen.ID("ctx"), jen.ID("premade")),
					jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
					jen.Line(),
					jen.Commentf("Fetch %s", scn),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Get%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.ID("err")),
					jen.Line(),
					jen.Commentf("Assert %s equality", scn),
					jen.IDf("check%sEquality", sn).Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.ID("actual").Dot("UpdatedOn")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("actual").Dot("ID")),
					jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("tctx"), jen.ID("t").Dot("Name").Call()),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Comment("CreateHandler item"),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(rootPkg, "models/v1"), sn).Valuesln(
						jen.ID("Name").Op(":").Lit("name"),
						jen.ID("Details").Op(":").Lit("details"),
					),
					jen.List(jen.ID("premade"), jen.ID("err")).Op(":=").ID("todoClient").Dotf("Create%s", sn).Call(jen.ID("ctx"), jen.Op("&").Qual(filepath.Join(rootPkg, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(
						jen.ID("Name").Op(":").ID("expected").Dot("Name"),
						jen.ID("Details").Op(":").ID("expected").Dot("Details")),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.ID("err")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.ID("err").Op("=").ID("todoClient").Dotf("Archive%s", sn).Call(jen.ID("ctx"), jen.ID("premade").Dot("ID")),
					jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				)),
			)),
		),
		jen.Line(),
	)
	return ret
}
