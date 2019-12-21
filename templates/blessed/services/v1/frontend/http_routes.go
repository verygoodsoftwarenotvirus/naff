package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(pkg.OutputPath, pkg.DataTypes, ret)

	ret.Add(
		jen.Comment("Routes returns a map of route to HandlerFunc for the parent router to set"),
		jen.Line(),
		jen.Comment("this keeps routing logic in the frontend service and not in the server itself."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("Routes").Params().Params(jen.Map(jen.ID("string")).Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Map(jen.ID("string")).Qual("net/http", "HandlerFunc").Valuesln(
				jen.Comment(`"/login":    s.LoginPage`),
				jen.Comment(`"/register": s.RegistrationPage`),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("buildStaticFileServer").Params(jen.ID("fileDir").ID("string")).Params(jen.Op("*").Qual("github.com/spf13/afero", "HttpFs"), jen.ID("error")).Block(jen.Var().ID("afs").Qual("github.com/spf13/afero", "Fs"), jen.If(jen.ID("s").Dot("config").Dot("CacheStaticFiles")).Block(
			jen.ID("afs").Op("=").Qual("github.com/spf13/afero", "NewMemMapFs").Call(),
			jen.List(jen.ID("files"), jen.ID("err")).Op(":=").Qual("io/ioutil", "ReadDir").Call(jen.ID("fileDir")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("reading directory for frontend files: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("_"), jen.ID("file")).Op(":=").Range().ID("files")).Block(
				jen.If(jen.ID("file").Dot("IsDir").Call()).Block(
					jen.Continue(),
				),
				jen.Line(),
				jen.ID("fp").Op(":=").Qual("path/filepath", "Join").Call(jen.ID("fileDir"), jen.ID("file").Dot("Name").Call()),
				jen.List(jen.ID("f"), jen.ID("err")).Op(":=").ID("afs").Dot("Create").Call(jen.ID("fp")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("creating static file in memory: %w"), jen.ID("err"))),
				),
				jen.Line(),
				jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("io/ioutil", "ReadFile").Call(jen.ID("fp")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("reading static file from directory: %w"), jen.ID("err"))),
				),
				jen.Line(),
				jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("f").Dot("Write").Call(jen.ID("bs")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading static file into memory: %w"), jen.ID("err"))),
				),
				jen.Line(),
				jen.If(jen.ID("err").Op("=").ID("f").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("closing file while setting up static dir")),
				),
			),
			jen.ID("afs").Op("=").Qual("github.com/spf13/afero", "NewReadOnlyFs").Call(jen.ID("afs")),
		).Else().Block(
			jen.ID("afs").Op("=").Qual("github.com/spf13/afero", "NewOsFs").Call(),
		),
			jen.Line(),
			jen.Return().List(jen.Qual("github.com/spf13/afero", "NewHttpFs").Call(jen.ID("afs")), jen.ID("nil")),
		),
		jen.Line(),
	)

	buildRegexPairs := func() []jen.Code {
		pairs := []jen.Code{
			jen.Comment("Here is where you should put route regexes that need to be ignored by the static file server."),
			jen.Comment("For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic"),
			jen.Comment("information, such as `/event/123`, you would want to put something like this below:"),
			jen.Comment("		eventsFrontendPathRegex = regexp.MustCompile(`/event/\\d+`)"),
			jen.Line(),
		}

		for _, typ := range pkg.DataTypes {
			tuvn := typ.Name.PluralUnexportedVarName()
			pairs = append(pairs,
				jen.Commentf("%sFrontendPathRegex matches URLs against our frontend router's specification for specific %s routes", tuvn, typ.Name.SingularCommonName()),
				jen.IDf("%sFrontendPathRegex", tuvn).Op("=").Qual("regexp", "MustCompile").Call(jen.RawStringf(`/%s/\d+`, typ.Name.PluralRouteName())),
			)
		}

		return pairs
	}

	ret.Add(
		jen.Var().Defs(
			buildRegexPairs()...,
		),
		jen.Line(),
	)

	buildHistoryRoutes := func() []jen.Code {
		lines := []jen.Code{
			jen.Lit("/register"),
			jen.Lit("/login"),
		}

		for _, typ := range pkg.DataTypes {
			prn := typ.Name.PluralRouteName()
			lines = append(lines, jen.Litf("/%s", prn), jen.Litf("/%s/new", prn))
		}

		lines = append(lines, jen.Lit("/password/new"))

		return lines
	}

	buildIndexRoute := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("rl").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("rl").Dot("Debug").Call(jen.Lit("static file requested")),
			jen.Switch(jen.ID("req").Dot("URL").Dot("Path")).Block(
				jen.Comment("list your frontend history routes here"),
				jen.Caseln(
					buildHistoryRoutes()...,
				).Block(jen.ID("rl").Dot("Debug").Call(jen.Lit("rerouting")),
					jen.ID("req").Dot("URL").Dot("Path").Op("=").Lit("/"),
				),
			),
		}

		for _, typ := range pkg.DataTypes {
			tpuvn := typ.Name.PluralUnexportedVarName()
			lines = append(lines,
				jen.If(jen.IDf("%sFrontendPathRegex", tpuvn).Dot("MatchString").Call(jen.ID("req").Dot("URL").Dot("Path"))).Block(
					jen.ID("rl").Dot("Debug").Call(jen.Lit("rerouting request")),
					jen.ID("req").Dot("URL").Dot("Path").Op("=").Lit("/"),
				),
			)
		}

		lines = append(lines,
			jen.Line(),
			jen.ID("fs").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
		)

		return lines
	}

	ret.Add(
		jen.Comment("StaticDir builds a static directory handler"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("StaticDir").Params(jen.ID("staticFilesDirectory").ID("string")).Params(jen.Qual("net/http", "HandlerFunc"), jen.ID("error")).Block(
			jen.List(jen.ID("fileDir"), jen.ID("err")).Op(":=").Qual("path/filepath", "Abs").Call(jen.ID("staticFilesDirectory")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("determining absolute path of static files directory: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("httpFs"), jen.ID("err")).Op(":=").ID("s").Dot("buildStaticFileServer").Call(jen.ID("fileDir")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("establishing static server filesystem: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("static_dir"), jen.ID("fileDir")).Dot("Debug").Call(jen.Lit("setting static file server")),
			jen.ID("fs").Op(":=").Qual("net/http", "StripPrefix").Call(jen.Lit("/"), jen.Qual("net/http", "FileServer").Call(jen.ID("httpFs").Dot("Dir").Call(jen.ID("fileDir")))),
			jen.Line(),
			jen.Return().List(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				buildIndexRoute()...,
			), jen.ID("nil")),
		),
		jen.Line(),
	)
	return ret
}
