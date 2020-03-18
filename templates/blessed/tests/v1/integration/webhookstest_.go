package integration

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("integration")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("checkWebhookEquality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("expected"), jen.ID("actual")).Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ID")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Name"), jen.ID("actual").Dot("Name")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("ContentType"), jen.ID("actual").Dot("ContentType")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("URL"), jen.ID("actual").Dot("URL")),
			jen.Qual("github.com/stretchr/testify/assert", "Equal").Call(jen.ID("t"), jen.ID("expected").Dot("Method"), jen.ID("actual").Dot("Method")),
			jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("CreatedOn")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyWebhookInput").Params().Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput")).Block(
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
				jen.ID("Name").Op(":").Qual(utils.FakeLibrary, "Word").Call(),
				jen.ID("URL").Op(":").Qual(utils.FakeLibrary, "DomainName").Call(),
				jen.ID("ContentType").Op(":").Lit("application/json"),
				jen.ID("Method").Op(":").Qual("net/http", "MethodPost"),
			),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("buildDummyWebhook").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.List(jen.ID("y"), jen.Err()).Op(":=").ID("todoClient").Dot("CreateWebhook").Call(jen.Qual("context", "Background").Call(), jen.ID("buildDummyWebhookInput").Call()),
			jen.Qual("github.com/stretchr/testify/require", "NoError").Call(jen.ID("t"), jen.Err()),
			jen.Return().ID("y"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("reverse").Params(jen.ID("s").ID("string")).Params(jen.ID("string")).Block(
			jen.ID("runes").Op(":=").Index().ID("rune").Call(jen.ID("s")),
			jen.For(jen.List(jen.ID("i"), jen.ID("j")).Op(":=").List(jen.Lit(0), jen.ID("len").Call(jen.ID("runes")).Op("-").Add(utils.FakeUint64Func())), jen.ID("i").Op("<").ID("j"), jen.List(jen.ID("i"), jen.ID("j")).Op("=").List(jen.ID("i").Op("+").Add(utils.FakeUint64Func()), jen.ID("j").Op("-").Add(utils.FakeUint64Func()))).Block(
				jen.List(jen.ID("runes").Index(jen.ID("i")), jen.ID("runes").Index(jen.ID("j"))).Op("=").List(jen.ID("runes").Index(jen.ID("j")), jen.ID("runes").Index(jen.ID("i"))),
			),
			jen.Return().ID("string").Call(jen.ID("runes")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestWebhooks").Params(jen.ID("test").Op("*").Qual("testing", "T")).Block(
			jen.ID("test").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Creating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be createable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.ID("Name").Op(":").ID("input").Dot("Name"),
						jen.ID("URL").Op(":").ID("input").Dot("URL"),
						jen.ID("ContentType").Op(":").ID("input").Dot("ContentType"),
						jen.ID("Method").Op(":").ID("input").Dot("Method")),
					jen.List(jen.ID("premade"), jen.Err()).Op(":=").ID("todoClient").Dot("CreateWebhook").Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
						jen.ID("Name").Op(":").ID("expected").Dot("Name"),
						jen.ID("ContentType").Op(":").ID("expected").Dot("ContentType"),
						jen.ID("URL").Op(":").ID("expected").Dot("URL"),
						jen.ID("Method").Op(":").ID("expected").Dot("Method"))),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("premade")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Op("=").ID("todoClient").Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("premade").Dot("ID")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
					jen.Line(), // REPEATME NOTICEME
					jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("todoClient").Dot("GetWebhook").Call(jen.ID("tctx"), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Qual("github.com/stretchr/testify/assert", "NotZero").Call(jen.ID("t"), jen.ID("actual").Dot("ArchivedOn")),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Listing"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be read in a list"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create webhooks"),
					jen.Var().ID("expected").Index().Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook"),
					jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").Lit(5), jen.ID("i").Op("++")).Block(
						jen.ID("expected").Op("=").ID("append").Call(jen.ID("expected"), jen.ID("buildDummyWebhook").Call(jen.ID("t"))),
					),
					jen.Line(),
					jen.Comment("Assert webhook list equality"),
					jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("todoClient").Dot("GetWebhooks").Call(jen.ID("tctx"), jen.ID("nil")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Qual("github.com/stretchr/testify/assert", "True").Call(jen.ID("t"), jen.ID("len").Call(jen.ID("expected")).Op("<=").ID("len").Call(jen.ID("actual").Dot("Webhooks"))),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.For(jen.List(jen.ID("_"), jen.ID("webhook")).Op(":=").Range().ID("actual").Dot("Webhooks")).Block(
						jen.Err().Op("=").ID("todoClient").Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("webhook").Dot("ID")),
						jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
					),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Reading"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to read something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("_"), jen.Err()).Op(":=").ID("todoClient").Dot("GetWebhook").Call(jen.ID("tctx"), jen.ID("nonexistentID")),
					jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be readable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.ID("Name").Op(":").ID("input").Dot("Name"),
						jen.ID("URL").Op(":").ID("input").Dot("URL"),
						jen.ID("ContentType").Op(":").ID("input").Dot("ContentType"),
						jen.ID("Method").Op(":").ID("input").Dot("Method")),
					jen.List(jen.ID("premade"), jen.Err()).Op(":=").ID("todoClient").Dot("CreateWebhook").Call(
						jen.ID("tctx"),
						jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
							jen.ID("Name").Op(":").ID("expected").Dot("Name"),
							jen.ID("ContentType").Op(":").ID("expected").Dot("ContentType"),
							jen.ID("URL").Op(":").ID("expected").Dot("URL"),
							jen.ID("Method").Op(":").ID("expected").Dot("Method"),
						),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("todoClient").Dot("GetWebhook").Call(jen.ID("tctx"), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Op("=").ID("todoClient").Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("actual").Dot("ID")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Updating"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("it should return an error when trying to update something that doesn't exist"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Err().Op(":=").ID("todoClient").Dot("UpdateWebhook").Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Values(jen.ID("ID").Op(":").ID("nonexistentID"))),
					jen.Qual("github.com/stretchr/testify/assert", "Error").Call(jen.ID("t"), jen.Err()),
				)),
				jen.Line(),
				jen.ID("T").Dot("Run").Call(jen.Lit("it should be updatable"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.ID("Name").Op(":").ID("input").Dot("Name"),
						jen.ID("URL").Op(":").ID("input").Dot("URL"),
						jen.ID("ContentType").Op(":").ID("input").Dot("ContentType"),
						jen.ID("Method").Op(":").ID("input").Dot("Method")),
					jen.List(jen.ID("premade"), jen.Err()).Op(":=").ID("todoClient").Dot("CreateWebhook").Call(
						jen.ID("tctx"),
						jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
							jen.ID("Name").Op(":").ID("expected").Dot("Name"),
							jen.ID("ContentType").Op(":").ID("expected").Dot("ContentType"),
							jen.ID("URL").Op(":").ID("expected").Dot("URL"),
							jen.ID("Method").Op(":").ID("expected").Dot("Method"),
						),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Change webhook"),
					jen.ID("premade").Dot("Name").Op("=").ID("reverse").Call(jen.ID("premade").Dot("Name")),
					jen.ID("expected").Dot("Name").Op("=").ID("premade").Dot("Name"),
					jen.Err().Op("=").ID("todoClient").Dot("UpdateWebhook").Call(jen.ID("tctx"), jen.ID("premade")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
					jen.Line(),
					jen.Comment("Fetch webhook"),
					jen.List(jen.ID("actual"), jen.Err()).Op(":=").ID("todoClient").Dot("GetWebhook").Call(jen.ID("tctx"), jen.ID("premade").Dot("ID")),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("actual"), jen.Err()),
					jen.Line(),
					jen.Comment("Assert webhook equality"),
					jen.ID("checkWebhookEquality").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
					jen.Qual("github.com/stretchr/testify/assert", "NotNil").Call(jen.ID("t"), jen.ID("actual").Dot("UpdatedOn")),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Op("=").ID("todoClient").Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("actual").Dot("ID")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				)),
			)),
			jen.Line(),
			jen.ID("test").Dot("Run").Call(jen.Lit("Deleting"), jen.Func().Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
				jen.ID("T").Dot("Run").Call(jen.Lit("should be able to be deleted"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
					jen.ID("tctx").Op(":=").Qual("context", "Background").Call(),
					jen.Line(),
					jen.Comment("Create webhook"),
					jen.ID("input").Op(":=").ID("buildDummyWebhookInput").Call(),
					jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Webhook").Valuesln(
						jen.ID("Name").Op(":").ID("input").Dot("Name"),
						jen.ID("URL").Op(":").ID("input").Dot("URL"),
						jen.ID("ContentType").Op(":").ID("input").Dot("ContentType"),
						jen.ID("Method").Op(":").ID("input").Dot("Method"),
					),
					jen.List(jen.ID("premade"), jen.Err()).Op(":=").ID("todoClient").Dot("CreateWebhook").Call(jen.ID("tctx"), jen.Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "WebhookCreationInput").Valuesln(
						jen.ID("Name").Op(":").ID("expected").Dot("Name"),
						jen.ID("ContentType").Op(":").ID("expected").Dot("ContentType"),
						jen.ID("URL").Op(":").ID("expected").Dot("URL"),
						jen.ID("Method").Op(":").ID("expected").Dot("Method")),
					),
					jen.ID("checkValueAndError").Call(jen.ID("t"), jen.ID("premade"), jen.Err()),
					jen.Line(),
					jen.Comment("Clean up"),
					jen.Err().Op("=").ID("todoClient").Dot("ArchiveWebhook").Call(jen.ID("tctx"), jen.ID("premade").Dot("ID")),
					jen.Qual("github.com/stretchr/testify/assert", "NoError").Call(jen.ID("t"), jen.Err()),
				)),
			)),
		),
		jen.Line(),
	)
	return ret
}
