package images

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func thumbnailsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("allSupportedColors").Op("=").Lit(2).Op("<<").Lit(7),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("thumbnailer").Interface(jen.ID("Thumbnail").Params(jen.ID("i").Op("*").ID("Image"), jen.List(jen.ID("width"), jen.ID("height")).ID("uint"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("newThumbnailer provides a thumbnailer given a particular content type."),
		jen.Line(),
		jen.Func().ID("newThumbnailer").Params(jen.ID("contentType").ID("string")).Params(jen.ID("thumbnailer"), jen.ID("error")).Body(
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("contentType")))).Body(
				jen.Case(jen.ID("imagePNG")).Body(
					jen.Return().List(jen.Op("&").ID("pngThumbnailer").Valuesln(), jen.ID("nil"))),
				jen.Case(jen.ID("imageJPEG")).Body(
					jen.Return().List(jen.Op("&").ID("jpegThumbnailer").Valuesln(), jen.ID("nil"))),
				jen.Case(jen.ID("imageGIF")).Body(
					jen.Return().List(jen.Op("&").ID("gifThumbnailer").Valuesln(), jen.ID("nil"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %s"),
						jen.ID("ErrInvalidImageContentType"),
						jen.ID("contentType"),
					))),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("preprocess").Params(jen.ID("i").Op("*").ID("Image"), jen.List(jen.ID("width"), jen.ID("height")).ID("uint")).Params(jen.Qual("image", "Image"), jen.ID("error")).Body(
			jen.List(jen.ID("img"), jen.ID("_"), jen.ID("err")).Op(":=").Qual("image", "Decode").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("i").Dot("Data"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("decoding image: %w"),
					jen.ID("err"),
				))),
			jen.ID("thumbnail").Op(":=").ID("resize").Dot("Thumbnail").Call(
				jen.ID("width"),
				jen.ID("height"),
				jen.ID("img"),
				jen.ID("resize").Dot("Lanczos3"),
			),
			jen.Return().List(jen.ID("thumbnail"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("jpegThumbnailer").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Thumbnail creates a GIF thumbnail."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("gifThumbnailer")).ID("Thumbnail").Params(jen.ID("img").Op("*").ID("Image"), jen.List(jen.ID("width"), jen.ID("height")).ID("uint"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error")).Body(
			jen.List(jen.ID("thumbnail"), jen.ID("err")).Op(":=").ID("preprocess").Call(
				jen.ID("img"),
				jen.ID("width"),
				jen.ID("height"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.Var().Defs(
				jen.ID("b").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op("=").Qual("image/png", "Encode").Call(
				jen.Op("&").ID("b"),
				jen.ID("thumbnail"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encoding PNG: %w"),
					jen.ID("err"),
				))),
			jen.ID("bs").Op(":=").ID("b").Dot("Bytes").Call(),
			jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.png"),
				jen.ID("filename"),
			), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("bs"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("bs"))),
			jen.Return().List(jen.ID("i"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("pngThumbnailer").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Thumbnail creates a GIF thumbnail."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("gifThumbnailer")).ID("Thumbnail").Params(jen.ID("img").Op("*").ID("Image"), jen.List(jen.ID("width"), jen.ID("height")).ID("uint"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error")).Body(
			jen.List(jen.ID("thumbnail"), jen.ID("err")).Op(":=").ID("preprocess").Call(
				jen.ID("img"),
				jen.ID("width"),
				jen.ID("height"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.Var().Defs(
				jen.ID("b").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op("=").Qual("image/png", "Encode").Call(
				jen.Op("&").ID("b"),
				jen.ID("thumbnail"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encoding PNG: %w"),
					jen.ID("err"),
				))),
			jen.ID("bs").Op(":=").ID("b").Dot("Bytes").Call(),
			jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.png"),
				jen.ID("filename"),
			), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("bs"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("bs"))),
			jen.Return().List(jen.ID("i"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("gifThumbnailer").Struct(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Thumbnail creates a GIF thumbnail."),
		jen.Line(),
		jen.Func().Params(jen.ID("t").Op("*").ID("gifThumbnailer")).ID("Thumbnail").Params(jen.ID("img").Op("*").ID("Image"), jen.List(jen.ID("width"), jen.ID("height")).ID("uint"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error")).Body(
			jen.List(jen.ID("thumbnail"), jen.ID("err")).Op(":=").ID("preprocess").Call(
				jen.ID("img"),
				jen.ID("width"),
				jen.ID("height"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.Var().Defs(
				jen.ID("b").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op("=").Qual("image/png", "Encode").Call(
				jen.Op("&").ID("b"),
				jen.ID("thumbnail"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encoding PNG: %w"),
					jen.ID("err"),
				))),
			jen.ID("bs").Op(":=").ID("b").Dot("Bytes").Call(),
			jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.png"),
				jen.ID("filename"),
			), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("bs"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("bs"))),
			jen.Return().List(jen.ID("i"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
