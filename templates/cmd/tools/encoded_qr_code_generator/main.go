package encoded_qr_code_generator

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("base64ImagePrefix").Op("=").Lit("data:image/png;base64,"),
			jen.ID("otpAuthURL").Op("=").Litf(`otpauth://totp/%s:username?secret=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=&issuer=%s`, proj.Name.RouteName(), proj.Name.RouteName()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("main").Params().Body(
			jen.List(jen.ID("bmp"), jen.ID("err")).Op(":=").ID("qrcode").Dot("NewQRCodeWriter").Call().Dot("EncodeWithoutHint").Call(
				jen.ID("otpAuthURL"),
				jen.ID("gozxing").Dot("BarcodeFormat_QR_CODE"),
				jen.Lit(128),
				jen.Lit(128),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual("log", "Fatal").Call(jen.ID("err"))),
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.If(jen.ID("err").Op("=").Qual("image/png", "Encode").Call(
				jen.Op("&").ID("b"),
				jen.ID("bmp"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Qual("log", "Fatal").Call(jen.ID("err"))),
			jen.Qual("log", "Println").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s%s"),
				jen.ID("base64ImagePrefix"),
				jen.Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b").Dot("Bytes").Call()),
			)),
		),
		jen.Newline(),
	)

	return code
}
