package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func httpRoutesDotGo() *jen.File {
	ret := jen.NewFile("frontend")
	utils.AddImports(ret)

	ret.Add(jen.Func())
	ret.Add(jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("buildStaticFileServer").Params(jen.ID("fileDir").ID("string")).Params(jen.Op("*").ID("afero").Dot(
		"HttpFs",
	), jen.ID("error")).Block(
		jen.Null().Var().ID("afs").ID("afero").Dot(
			"Fs",
		),
		jen.If(
			jen.ID("s").Dot(
				"config",
			).Dot(
				"CacheStaticFiles",
			),
		).Block(
			jen.ID("afs").Op("=").ID("afero").Dot(
				"NewMemMapFs",
			).Call(),
			jen.List(jen.ID("files"), jen.ID("err")).Op(":=").Qual("io/ioutil", "ReadDir").Call(jen.ID("fileDir")),
			jen.If(
				jen.ID("err").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("reading directory for frontend files: %w"), jen.ID("err"))),
			),
			jen.For(jen.List(jen.ID("_"), jen.ID("file")).Op(":=").Range().ID("files")).Block(
				jen.If(
					jen.ID("file").Dot(
						"IsDir",
					).Call(),
				).Block(
					jen.Continue(),
				),
				jen.ID("fp").Op(":=").Qual("path/filepath", "Join").Call(jen.ID("fileDir"), jen.ID("file").Dot(
					"Name",
				).Call()),
				jen.List(jen.ID("f"), jen.ID("err")).Op(":=").ID("afs").Dot(
					"Create",
				).Call(jen.ID("fp")),
				jen.If(
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("creating static file in memory: %w"), jen.ID("err"))),
				),
				jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("io/ioutil", "ReadFile").Call(jen.ID("fp")),
				jen.If(
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("reading static file from directory: %w"), jen.ID("err"))),
				),
				jen.If(
					jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("f").Dot(
						"Write",
					).Call(jen.ID("bs")),
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading static file into memory: %w"), jen.ID("err"))),
				),
				jen.If(
					jen.ID("err").Op("=").ID("f").Dot(
						"Close",
					).Call(),
					jen.ID("err").Op("!=").ID("nil"),
				).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Error",
					).Call(jen.ID("err"), jen.Lit("closing file while setting up static dir")),
				),
			),
			jen.ID("afs").Op("=").ID("afero").Dot(
				"NewReadOnlyFs",
			).Call(jen.ID("afs")),
		).Else().Block(
			jen.ID("afs").Op("=").ID("afero").Dot(
				"NewOsFs",
			).Call(),
		),
		jen.Return().List(jen.ID("afero").Dot(
			"NewHttpFs",
		).Call(jen.ID("afs")), jen.ID("nil")),
	),
	)
	ret.Add(jen.Null().Var().ID("itemsFrontendPathRegex").Op("=").Qual("regexp", "MustCompile").Call(jen.Lit(`/items/\d+`)))
	ret.Add(jen.Func())
	return ret
}
