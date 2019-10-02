package twofactor

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

var (
	// files are all the available files to generate
	files = map[string]*jen.File{
		"cmd/tools/two_factor/main.go": mainDotGo(),
	}
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) {
	for path, file := range files {
		renderFile(path, file)
	}
}

func renderFile(path string, file *jen.File) {
	fp := utils.BuildTemplatePath(path)
	_ = os.Remove(fp)

	var b bytes.Buffer
	if err := file.Render(&b); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fp, b.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func mainDotGo() *jen.File {
	ret := jen.NewFile("main")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("zero").Op("=").Lit("  ___   & / _ \\  &| | | | &| |_| | & \\___/  "),
			jen.ID("one").Op("=").Lit("    _    &  /_ |   &   | |   &  _| |_  & |_____| "),
			jen.ID("two").Op("=").Lit(" ____   &|___ \\  &  __) | & / __/  &|_____| "),
			jen.ID("three").Op("=").Lit("_____   &|___ /  &  |_ \\  & ___) | &|____/  "),
			jen.ID("four").Op("=").Lit(" _   _   &| | | |  &| |_| |_ &|___   _ &    |_|  "),
			jen.ID("five").Op("=").Lit(" ____   &| ___|  &|___ \\  & ___) | &|____/  "),
			jen.ID("six").Op("=").Lit("  __    & / /_   &| '_ \\  &| (_) | & \\___/  "),
			jen.ID("seven").Op("=").Lit(" _____  &|___  | &   / /  &  / /   & /_/    "),
			jen.ID("eight").Op("=").Lit("  ___   & ( o )  & /   \\  &|  O  | & \\___/  "),
			jen.ID("nine").Op("=").Lit("  ___   & /   \\  &| (_) | & \\__, | &   /_/  "),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("lastChange").Qual("time", "Time"),
			jen.ID("currentCode").ID("string"),
			jen.Line(),
			jen.Comment("feel free to link to this variable and the related  non-stdlib"),
			jen.Comment("functions as a demonstration of useless over-engineering"),
			jen.ID("numbers").Op("=").Index(jen.Lit(10)).Index(jen.Lit(5)).ID("string").Valuesln(
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
	)

	ret.Add(
		jen.Func().ID("limitSlice").Params(jen.ID("in").Index().ID("string")).Params(jen.ID("out").Index(jen.Lit(5)).ID("string")).Block(
			jen.If(jen.ID("len").Call(jen.ID("in")).Op("!=").Lit(5)).Block(
				jen.ID("panic").Call(jen.Lit("wut")),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
				jen.ID("out").Index(jen.ID("i")).Op("=").ID("in").Index(jen.ID("i")),
			),
			jen.Return(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("mustnt").Params(
			jen.ID("err").ID("error")).Block(
			jen.If(
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.ID("panic").Call(jen.ID("err")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("clearTheScreen").Params().Block(
			jen.Qual("fmt", "Println").Call(jen.Lit(`033[2J`)),
			jen.Qual("fmt", "Printf").Call(jen.Lit(`033[0;0H`)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildTheThing").Params(
			jen.ID("token").ID("string")).Params(
			jen.ID("string")).Block(
			jen.Var().ID("out").ID("string"),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
				jen.If(jen.ID("i").Op("!=").Lit(0)).Block(jen.ID("out").Op("+=").Lit("\n")),
				jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().Qual("strings", "Split").Call(jen.ID("token"), jen.Lit(""))).Block(
					jen.List(jen.ID("y"), jen.ID("err")).Op(":=").Qual("strconv", "Atoi").Call(jen.ID("x")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(jen.ID("panic").Call(jen.ID("err"))),
					jen.ID("out").Op("+=").Lit("  "),
					jen.ID("out").Op("+=").ID("numbers").Index(jen.ID("y")).Index(jen.ID("i")),
				),
			),
			jen.Line(),
			jen.ID("out").Op("+=").Lit("\n\n").Op("+").Parens(
				jen.Lit(30).Op("*").Qual("time", "Second").
					Op("-").Qual("time", "Since").
					Call(jen.ID("lastChange")).Dot("Round").
					Call(jen.Qual("time", "Second")),
			).Dot("String").Call().Op("+").Lit("\n"),
			jen.Line(),
			jen.Return().ID("out"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("doTheThing").Params(
			jen.ID("secret").ID("string")).Block(
			jen.ID("t").Op(":=").Qual("strings", "ToUpper").Call(jen.ID("secret")),
			jen.ID("n").Op(":=").Qual("time", "Now").Call().Dot(
				"UTC",
			).Call(),
			jen.List(jen.ID("code"), jen.ID("err")).Op(":=").Qual("github.com/pquerna/otp/totp", "GenerateCode").Call(
				jen.ID("t"), jen.ID("n")),
			jen.ID("mustnt").Call(jen.ID("err")),
			jen.Line(),
			jen.If(jen.ID("code").Op("!=").ID("currentCode")).Block(
				jen.ID("lastChange").Op("=").Qual("time", "Now").Call(),
				jen.ID("currentCode").Op("=").ID("code"),
			),
			jen.Line(),
			jen.If(
				jen.Op("!").ID("totp").Dot(
					"Validate",
				).Call(jen.ID("code"), jen.ID("t")),
			).Block(
				jen.ID("panic").Call(jen.Lit("this shouldn't happen")),
			),
			jen.Line(),
			jen.ID("clearTheScreen").Call(),
			jen.Qual("fmt", "Println").Call(jen.ID("buildTheThing").Call(jen.ID("code"))),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("requestTOTPSecret").Params().Params(
			jen.ID("string")).Block(
			jen.Var().Defs(
				jen.ID("token").ID("string"),
				jen.ID("err").ID("error"),
			),
			jen.Line(),
			jen.If(
				jen.ID("len").Call(jen.Qual("os", "Args")).Op("==").Lit(1),
			).Block(
				jen.ID("reader").Op(":=").Qual("bufio", "NewReader").Call(jen.Qual("os", "Stdin")),
				jen.Qual("fmt", "Print").Call(jen.Lit("token: ")),
				jen.List(jen.ID("token"), jen.ID("err")).Op("=").ID("reader").Dot(
					"ReadString",
				).Call(jen.ID(`'\n'`)),
				jen.ID("mustnt").Call(jen.ID("err")),
			).Else().Block(
				jen.ID("token").Op("=").Qual("os", "Args").Index(jen.Lit(1)),
			),
			jen.Line(),
			jen.Return().ID("token"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("main").Params().Block(
			jen.ID("secret").Op(":=").ID("requestTOTPSecret").Call(),
			jen.ID("clearTheScreen").Call(),
			jen.ID("doTheThing").Call(jen.ID("secret")),
			jen.ID("every").Op(":=").Qual("time", "Tick").Call(jen.Lit(1).Op("*").Qual("time", "Second")),
			jen.ID("lastChange").Op("=").Qual("time", "Now").Call(),
			jen.Line(),
			jen.For().Range().ID("every").Block(
				jen.ID("doTheThing").Call(jen.ID("secret")),
			),
		),
	)
	return ret
}
