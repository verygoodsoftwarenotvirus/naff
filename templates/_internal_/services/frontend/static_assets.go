package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func staticAssetsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Anon("embed")

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildStaticFileServer").Params().Params(jen.Op("*").ID("afero").Dot("HttpFs")).Body(
			jen.Return().ID("afero").Dot("NewHttpFs").Call(jen.ID("afero").Dot("NewOsFs").Call())),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.Comment("Here is where you should put route regexes that need to be ignored by the static file server."),
			jen.Comment("For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic."),
			jen.Comment("information, such as `/event/123`, you would want to put something like this below:"),
			jen.Comment(`		eventsFrontendPathRegex = regexp.MustCompile(`+"`"+`/event/\d+`+"`"+`)`),
		),
		jen.Newline(),
	)

	routes := []jen.Code{
		jen.Lit("/register"),
		jen.Lit("/login"),
		jen.Lit("/home"),
		jen.Lit("/account"),
		jen.Lit("/admin"),
		jen.Lit("/admin/dashboard"),
		jen.Lit("/admin/users"),
		jen.Lit("/"),
	}

	code.Add(
		jen.Comment("StaticDir builds a static directory handler."),
		jen.Newline(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("StaticDir").Params(jen.ID("staticFilesDirectory").ID("string")).Params(jen.Qual("net/http", "HandlerFunc"), jen.ID("error")).Body(
			jen.List(jen.ID("fileDir"), jen.ID("err")).Op(":=").Qual("path/filepath", "Abs").Call(jen.ID("staticFilesDirectory")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("determining absolute path of static files directory: %w"),
					jen.ID("err"),
				))),
			jen.Newline(),
			jen.ID("httpFs").Op(":=").ID("s").Dot("buildStaticFileServer").Call(),
			jen.Newline(),
			jen.ID("s").Dot("logger").Dot("WithValue").Call(
				jen.Lit("static_dir"),
				jen.ID("fileDir"),
			).Dot("Debug").Call(jen.Lit("setting static file server")),
			jen.ID("fs").Op(":=").Qual("net/http", "StripPrefix").Call(
				jen.Lit("/"),
				jen.Qual("net/http", "FileServer").Call(jen.ID("httpFs").Dot("Dir").Call(jen.ID("fileDir"))),
			),
			jen.Newline(),
			jen.Return().List(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.ID("logger").Dot("Debug").Call(jen.Lit("static file requested")),
				jen.Newline(),
				jen.List(jen.ID("sessCtxData"), jen.ID("sessCtxErr")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
				jen.If(jen.ID("sessCtxErr").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("Error").Call(
						jen.ID("sessCtxErr"),
						jen.Lit("fetching session context data"),
					)),
				jen.Newline(),
				jen.If(jen.Qual("strings", "HasPrefix").Call(
					jen.ID("req").Dot("URL").Dot("Path"),
					jen.Lit("/admin"),
				).Op("&&").ID("sessCtxData").Op("!=").ID("nil").Op("&&").Op("!").ID("sessCtxData").Dot("ServiceRolePermissionChecker").Call().Dot("IsServiceAdmin").Call()).Body(
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Qual("net/http", "Redirect").Call(
						jen.ID("res"),
						jen.ID("req"),
						jen.Lit("/login"),
						jen.Qual("net/http", "StatusUnauthorized"),
					),
					jen.Return(),
				),
				jen.Newline(),
				jen.Switch(jen.ID("req").Dot("URL").Dot("Path")).Body(
					jen.Comment("list your frontend history routes here."),
					jen.Caseln(routes...).Body(
						jen.ID("logger").Dot("Debug").Call(jen.Lit("rerouting")), jen.ID("req").Dot("URL").Dot("Path").Op("=").Lit("/"))),
				jen.Newline(),
				jen.Comment("if eventsFrontendPathRegex.MatchString(req.URL.Path) {"),
				jen.Comment(`	logger.Debug("rerouting request")`),
				jen.Comment(`	req.URL.Path = "/"`),
				jen.Comment("}"),
				jen.Newline(),
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("destination"),
					jen.ID("req").Dot("URL").Dot("Path"),
				).Dot("Debug").Call(jen.Lit("heading to frontend path")),
				jen.Newline(),
				jen.ID("fs").Dot("ServeHTTP").Call(
					jen.ID("res"),
					jen.ID("req"),
				),
			), jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
