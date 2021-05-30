package images

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func thumbnailsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_newThumbnailer").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.For(jen.List(jen.ID("_"), jen.ID("ct")).Op(":=").Range().Index().ID("string").Valuesln(jen.ID("imagePNG"), jen.ID("imageJPEG"), jen.ID("imageGIF"))).Body(
						jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("newThumbnailer").Call(jen.ID("ct")),
						jen.ID("assert").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("x"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("invalid content type"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("newThumbnailer").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_preprocess").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("imgBytes").Op(":=").ID("buildPNGBytes").Call(jen.ID("t")).Dot("Bytes").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("imgBytes"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("imgBytes"))),
					jen.List(jen.ID("img"), jen.ID("err")).Op(":=").ID("preprocess").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("img"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(1024)),
					jen.List(jen.ID("img"), jen.ID("err")).Op(":=").ID("preprocess").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("img"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_jpegThumbnailer_Thumbnail").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("imgBytes").Op(":=").ID("buildJPEGBytes").Call(jen.ID("t")).Dot("Bytes").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("imgBytes"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("imgBytes"))),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").Parens(jen.Op("&").ID("jpegThumbnailer").Valuesln()).Dot("Thumbnail").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(1024)),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").Parens(jen.Op("&").ID("jpegThumbnailer").Valuesln()).Dot("Thumbnail").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_pngThumbnailer_Thumbnail").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("imgBytes").Op(":=").ID("buildPNGBytes").Call(jen.ID("t")).Dot("Bytes").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("imgBytes"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("imgBytes"))),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").Parens(jen.Op("&").ID("pngThumbnailer").Valuesln()).Dot("Thumbnail").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(1024)),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").Parens(jen.Op("&").ID("pngThumbnailer").Valuesln()).Dot("Thumbnail").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_gifThumbnailer_Thumbnail").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("imgBytes").Op(":=").ID("buildGIFBytes").Call(jen.ID("t")).Dot("Bytes").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("imgBytes"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("imgBytes"))),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").Parens(jen.Op("&").ID("gifThumbnailer").Valuesln()).Dot("Thumbnail").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Lit("whatever.png"), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(1024)),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").Parens(jen.Op("&").ID("gifThumbnailer").Valuesln()).Dot("Thumbnail").Call(
						jen.ID("i"),
						jen.Lit(128),
						jen.Lit(128),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
