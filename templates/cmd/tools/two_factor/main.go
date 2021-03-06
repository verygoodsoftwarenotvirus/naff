package two_factor

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"cmd/tools/two_factor/main.go": mainDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("main")

	utils.AddImports(proj, code)

	code.PackageComment(`Command two_factor is a CLI that takes in a secret as a positional argument
and draws the TOTP code for that secret in big ASCII numbers. This command is
helpful when you need to repeatedly test the logic of registering an account
and logging in.`)

	code.Add(buildConstDeclarations()...)
	code.Add(buildVarDeclarations()...)
	code.Add(buildLimitSlice()...)
	code.Add(buildMustnt()...)
	code.Add(buildClearTheScreen()...)
	code.Add(buildBuildTheThing()...)
	code.Add(buildDoTheThing()...)
	code.Add(buildRequestTOTPSecret()...)
	code.Add(buildMain()...)

	return code
}

func buildConstDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("zero").Equals().Lit("  ___   & / _ \\  &| | | | &| |_| | & \\___/  "),
			jen.ID("one").Equals().Lit("    _    &  /_ |   &   | |   &  _| |_  & |_____| "),
			jen.ID("two").Equals().Lit(" ____   &|___ \\  &  __) | & / __/  &|_____| "),
			jen.ID("three").Equals().Lit("_____   &|___ /  &  |_ \\  & ___) | &|____/  "),
			jen.ID("four").Equals().Lit(" _   _   &| | | |  &| |_| |_ &|___   _ &    |_|  "),
			jen.ID("five").Equals().Lit(" ____   &| ___|  &|___ \\  & ___) | &|____/  "),
			jen.ID("six").Equals().Lit("  __    & / /_   &| '_ \\  &| (_) | & \\___/  "),
			jen.ID("seven").Equals().Lit(" _____  &|___  | &   / /  &  / /   & /_/    "),
			jen.ID("eight").Equals().Lit("  ___   & ( o )  & /   \\  &|  O  | & \\___/  "),
			jen.ID("nine").Equals().Lit("  ___   & /   \\  &| (_) | & \\__, | &   /_/  "),
		),
		jen.Line(),
	}

	return lines
}

func buildVarDeclarations() []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("lastChange").Qual("time", "Time"),
			jen.ID("currentCode").String(),
			jen.Line(),
			jen.ID("numbers").Equals().Index(jen.Lit(10)).Index(jen.Lit(5)).String().Valuesln(
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("zero"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("one"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("two"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("three"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("four"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("five"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("six"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("seven"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("eight"), jen.Lit("&"))),
				jen.ID("limitSlice").Call(jen.Qual("strings", "Split").Call(jen.ID("nine"), jen.Lit("&"))),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildLimitSlice() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("limitSlice").Params(jen.ID("in").Index().String()).Params(jen.ID("out").Index(jen.Lit(5)).String()).Body(
			jen.If(jen.Len(jen.ID("in")).DoesNotEqual().Lit(5)).Body(
				jen.ID("panic").Call(jen.Lit("wut")),
			),
			jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
				jen.ID("out").Index(jen.ID("i")).Equals().ID("in").Index(jen.ID("i")),
			),
			jen.Return(),
		),
		jen.Line(),
	}

	return lines
}

func buildMustnt() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("mustnt").Params(jen.Err().Error()).Body(
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildClearTheScreen() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("clearTheScreen").Params().Body(
			jen.Qual("fmt", "Println").Call(jen.Lit("\033[2J")),
			jen.Qual("fmt", "Printf").Call(jen.Lit("\033[0;0H")),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildTheThing() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTheThing").Params(jen.ID("token").String()).Params(jen.String()).Body(
			jen.Var().ID("out").String(),
			jen.For(jen.ID("i").Assign().Zero(), jen.ID("i").LessThan().Lit(5), jen.ID("i").Op("++")).Body(
				jen.If(jen.ID("i").DoesNotEqual().Zero()).Body(
					jen.ID("out").Op("+=").Lit("\n"),
				),
				jen.For(jen.List(jen.Underscore(), jen.ID("x")).Assign().Range().Qual("strings", "Split").Call(jen.ID("token"), jen.EmptyString())).Body(
					jen.List(jen.ID("y"), jen.Err()).Assign().Qual("strconv", "Atoi").Call(jen.ID("x")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
						jen.ID("panic").Call(jen.Err()),
					),
					jen.ID("out").Op("+=").Lit("  "),
					jen.ID("out").Op("+=").ID("numbers").Index(jen.ID("y")).Index(jen.ID("i")),
				),
			),
			jen.Line(),
			jen.ID("timeLeft").Assign().Parens(jen.Lit(30).Times().Qual("time", "Second").Minus().Qual("time", "Since").Call(jen.ID("lastChange")).Dot("Round").Call(jen.Qual("time", "Second"))).Dot("String").Call(),
			jen.ID("out").Op("+=").Qual("fmt", "Sprintf").Call(jen.Lit("\n\n%s\n"), jen.ID("timeLeft")),
			jen.Line(),
			jen.Return().ID("out"),
		),
		jen.Line(),
	}

	return lines
}

func buildDoTheThing() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("doTheThing").Params(jen.ID("secret").String()).Body(
			jen.ID("t").Assign().Qual("strings", "ToUpper").Call(jen.ID("secret")),
			jen.ID("n").Assign().Qual("time", "Now").Call().Dot("UTC").Call(),
			jen.List(jen.ID("code"), jen.Err()).Assign().Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(jen.ID("t"), jen.ID("n")),
			jen.ID("mustnt").Call(jen.Err()),
			jen.Line(),
			jen.If(jen.ID("code").DoesNotEqual().ID("currentCode")).Body(
				jen.ID("lastChange").Equals().Qual("time", "Now").Call(),
				jen.ID("currentCode").Equals().ID("code"),
			),
			jen.Line(),
			jen.If(jen.Not().Qual("github.com/pquerna/otp/totp", "Validate").Call(jen.ID("code"), jen.ID("t"))).Body(
				jen.ID("panic").Call(jen.Lit("this shouldn't happen")),
			),
			jen.Line(),
			jen.ID("clearTheScreen").Call(),
			jen.Qual("fmt", "Println").Call(jen.ID("buildTheThing").Call(jen.ID("code"))),
		),
		jen.Line(),
	}

	return lines
}

func buildRequestTOTPSecret() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("requestTOTPSecret").Params().Params(jen.String()).Body(
			jen.Var().Defs(
				jen.ID("token").String(),
				jen.Err().Error(),
			),
			jen.Line(),
			jen.If(jen.Len(jen.Qual("os", "Args")).IsEqualTo().One()).Body(
				jen.ID("reader").Assign().Qual("bufio", "NewReader").Call(jen.Qual("os", "Stdin")),
				jen.Qual("fmt", "Print").Call(jen.Lit("token: ")),
				jen.List(jen.ID("token"), jen.Err()).Equals().ID("reader").Dot("ReadString").Call(jen.ID(`'\n'`)),
				jen.ID("mustnt").Call(jen.Err()),
			).Else().Body(
				jen.ID("token").Equals().Qual("os", "Args").Index(jen.One()),
			),
			jen.Line(),
			jen.Return().ID("token"),
		),
		jen.Line(),
	}

	return lines
}

func buildMain() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("main").Params().Body(
			jen.ID("secret").Assign().ID("requestTOTPSecret").Call(),
			jen.ID("clearTheScreen").Call(),
			jen.ID("doTheThing").Call(jen.ID("secret")),
			jen.ID("every").Assign().Qual("time", "Tick").Call(jen.One().Times().Qual("time", "Second")),
			jen.ID("lastChange").Equals().Qual("time", "Now").Call(),
			jen.Line(),
			jen.For().Range().ID("every").Body(
				jen.ID("doTheThing").Call(jen.ID("secret")),
			),
		),
	}

	return lines
}
