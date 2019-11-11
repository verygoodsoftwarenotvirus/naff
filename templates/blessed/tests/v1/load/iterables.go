package load

import (
	"fmt"
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)

	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	ret.Add(
		jen.Commentf("fetchRandom%s retrieves a random %s from the list of available %s", sn, scn, pcn),
		jen.Line(),
		jen.Func().IDf("fetchRandom%s", sn).Params(jen.ID("c").Op("*").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "V1Client")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Block(
			jen.List(jen.IDf("%sRes", puvn), jen.ID("err")).Op(":=").ID("c").Dotf("Get%s", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
			jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").IDf("%sRes", puvn).Op("==").ID("nil").Op("||").ID("len").Call(jen.IDf("%sRes", puvn).Dot(pn)).Op("==").Lit(0)).Block(
				jen.Return().ID("nil"),
			),
			jen.Line(),
			jen.ID("randIndex").Op(":=").Qual("math/rand", "Intn").Call(jen.ID("len").Call(jen.IDf("%sRes", puvn).Dot(pn))),
			jen.Return().Op("&").IDf("%sRes", puvn).Dot(pn).Index(jen.ID("randIndex")),
		),
		jen.Line(),
	)

	buildRandomLines := func() []jen.Code {
		var lines []jen.Code

		for _, field := range typ.Fields {
			fsn := field.Name.Singular()
			lines = append(lines, jen.IDf("random%s", sn).Dot(fsn).Op("=").Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/rand/model"), fmt.Sprintf("Random%sCreationInput", sn)).Call().Dot(fsn))
		}

		lines = append(lines,
			jen.Return().ID("c").Dotf("BuildUpdate%sRequest", sn).Call(jen.Qual("context", "Background").Call(), jen.IDf("random%s", sn)),
		)

		return lines

	}

	ret.Add(
		jen.Func().IDf("build%sActions", sn).Params(jen.ID("c").Op("*").Qual(filepath.Join(pkg.OutputPath, "client/v1/http"), "V1Client")).Params(jen.Map(jen.ID("string")).Op("*").ID("Action")).Block(
			jen.Return().Map(jen.ID("string")).Op("*").ID("Action").Valuesln(
				jen.Litf("Create%s", sn).Op(":").Valuesln(
					jen.ID("Name").Op(":").Litf("Create%s", sn), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dotf("BuildCreate%sRequest", sn).Call(jen.Qual("context", "Background").Call(), jen.Qual(filepath.Join(pkg.OutputPath, "tests/v1/testutil/rand/model"), fmt.Sprintf("Random%sCreationInput", sn)).Call()),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Litf("Get%s", sn).Op(":").Valuesln(
					jen.ID("Name").Op(":").Litf("Get%s", sn), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.IDf("random%s", sn).Op(":=").IDf("fetchRandom%s", sn).Call(jen.ID("c")), jen.IDf("random%s", sn).Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dotf("BuildGet%sRequest", sn).Call(jen.Qual("context", "Background").Call(), jen.IDf("random%s", sn).Dot("ID")),
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Litf("Get%s", pn).Op(":").Valuesln(
					jen.ID("Name").Op(":").Litf("Get%s", pn), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.Return().ID("c").Dotf("BuildGet%sRequest", pn).Call(jen.Qual("context", "Background").Call(), jen.ID("nil")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Litf("Update%s", sn).Op(":").Valuesln(
					jen.ID("Name").Op(":").Litf("Update%s", sn), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.IDf("random%s", sn).Op(":=").IDf("fetchRandom%s", sn).Call(jen.ID("c")), jen.IDf("random%s", sn).Op("!=").ID("nil")).Block(
							buildRandomLines()...,
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(100)), jen.Litf("Archive%s", sn).Op(":").Valuesln(
					jen.ID("Name").Op(":").Litf("Archive%s", sn), jen.ID("Action").Op(":").Func().Params().Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Block(
						jen.If(jen.IDf("random%s", sn).Op(":=").IDf("fetchRandom%s", sn).Call(jen.ID("c")), jen.IDf("random%s", sn).Op("!=").ID("nil")).Block(
							jen.Return().ID("c").Dotf("BuildArchive%sRequest", sn).Call(jen.Qual("context", "Background").Call(), jen.IDf("random%s", sn).Dot("ID")),
						),
						jen.Return().List(jen.ID("nil"), jen.ID("ErrUnavailableYet")),
					),
					jen.ID("Weight").Op(":").Lit(85))),
		),
		jen.Line(),
	)
	return ret
}
