package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("defaultBrowserWaitTime").Op("=").Qual("time", "Second").Op("/").Lit(2),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestRegistrationFlow").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("helper").Op(":=").ID("setupTestHelper").Call(jen.ID("T")),
			jen.ID("helper").Dot("runForAllBrowsers").Call(
				jen.ID("T"),
				jen.Lit("registration flow"),
				jen.Func().Params(jen.ID("browser").ID("playwright").Dot("Browser")).Params(jen.Params(jen.Op("*").Qual("testing", "T"))).Body(
					jen.Return().Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
						jen.ID("user").Op(":=").ID("fakes").Dot("BuildFakeUserCreationInput").Call(),
						jen.List(jen.ID("page"), jen.ID("err")).Op(":=").ID("browser").Dot("NewPage").Call(),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
							jen.Lit("could not create page"),
						),
						jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("page").Dot("Goto").Call(jen.ID("urlToUse")),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
							jen.Lit("could not navigate to root page"),
						),
						jen.ID("registerLinkClickErr").Op(":=").ID("page").Dot("Click").Call(jen.Lit("#registerLink")),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("registerLinkClickErr"),
							jen.Lit("could not find register link on homepage"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Type").Call(
								jen.Lit("#usernameInput"),
								jen.ID("user").Dot("Username"),
							),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Type").Call(
								jen.Lit("#passwordInput"),
								jen.ID("user").Dot("Password"),
							),
						),
						jen.Qual("time", "Sleep").Call(jen.ID("defaultBrowserWaitTime")),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("urlToUse").Op("+").Lit("/register"),
							jen.ID("page").Dot("URL").Call(),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Click").Call(jen.Lit("#registrationButton")),
						),
						jen.Qual("time", "Sleep").Call(jen.ID("defaultBrowserWaitTime")),
						jen.List(jen.ID("qrCodeElement"), jen.ID("qrCodeElementErr")).Op(":=").ID("page").Dot("QuerySelector").Call(jen.Lit("#twoFactorSecretQRCode")),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("qrCodeElementErr"),
						),
						jen.List(jen.ID("img"), jen.ID("err")).Op(":=").Qual("image/png", "Decode").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("getScreenshotBytes").Call(
							jen.ID("t"),
							jen.ID("qrCodeElement"),
						))),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("err"),
						),
						jen.List(jen.ID("bmp"), jen.ID("bitmapErr")).Op(":=").ID("gozxing").Dot("NewBinaryBitmapFromImage").Call(jen.ID("img")),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("bitmapErr"),
						),
						jen.ID("qrReader").Op(":=").ID("qrcode").Dot("NewQRCodeReader").Call(),
						jen.List(jen.ID("result"), jen.ID("qrCodeDecodeErr")).Op(":=").ID("qrReader").Dot("Decode").Call(
							jen.ID("bmp"),
							jen.ID("nil"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("qrCodeDecodeErr"),
						),
						jen.List(jen.ID("u"), jen.ID("secretParseErr")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("result").Dot("String").Call()),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("secretParseErr"),
						),
						jen.ID("twoFactorSecret").Op(":=").ID("u").Dot("Query").Call().Dot("Get").Call(jen.Lit("secret")),
						jen.ID("require").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("twoFactorSecret"),
						),
						jen.List(jen.ID("code"), jen.ID("firstCodeGenerationErr")).Op(":=").ID("totp").Dot("GenerateCode").Call(
							jen.ID("twoFactorSecret"),
							jen.Qual("time", "Now").Call().Dot("UTC").Call(),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("firstCodeGenerationErr"),
						),
						jen.ID("require").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("code"),
						),
						jen.ID("totpInputFieldFindErr").Op(":=").ID("page").Dot("Type").Call(
							jen.Lit("#totpTokenInput"),
							jen.ID("code"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("totpInputFieldFindErr"),
							jen.Lit("unexpected error finding TOTP token input field: %v"),
							jen.ID("totpInputFieldFindErr"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Click").Call(jen.Lit("#totpTokenSubmitButton")),
						),
						jen.Qual("time", "Sleep").Call(jen.ID("defaultBrowserWaitTime")),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("urlToUse").Op("+").Lit("/login"),
							jen.ID("page").Dot("URL").Call(),
						),
						jen.List(jen.ID("code"), jen.ID("secondCodeGenerationErr")).Op(":=").ID("totp").Dot("GenerateCode").Call(
							jen.ID("twoFactorSecret"),
							jen.Qual("time", "Now").Call().Dot("UTC").Call(),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("secondCodeGenerationErr"),
						),
						jen.ID("require").Dot("NotEmpty").Call(
							jen.ID("t"),
							jen.ID("code"),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Type").Call(
								jen.Lit("#usernameInput"),
								jen.ID("user").Dot("Username"),
							),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Type").Call(
								jen.Lit("#passwordInput"),
								jen.ID("user").Dot("Password"),
							),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Type").Call(
								jen.Lit("#totpTokenInput"),
								jen.ID("code"),
							),
						),
						jen.ID("require").Dot("NoError").Call(
							jen.ID("t"),
							jen.ID("page").Dot("Click").Call(jen.Lit("#loginButton")),
						),
						jen.Qual("time", "Sleep").Call(jen.ID("defaultBrowserWaitTime")),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("urlToUse").Op("+").Lit("/"),
							jen.ID("page").Dot("URL").Call(),
						),
					)),
			),
		),
		jen.Line(),
	)

	return code
}
