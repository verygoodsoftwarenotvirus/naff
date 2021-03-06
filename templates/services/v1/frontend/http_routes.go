package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildFrontendBuildStaticFileServer()...)
	code.Add(buildFrontendVarDeclarations(proj)...)
	code.Add(buildFrontendStaticDir(proj)...)

	return code
}

func buildFrontendBuildStaticFileServer() []jen.Code {
	lines := []jen.Code{
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("buildStaticFileServer").Params(jen.ID("fileDir").String()).Params(jen.PointerTo().Qual("github.com/spf13/afero", "HttpFs"), jen.Error()).Body(jen.Var().ID("afs").Qual("github.com/spf13/afero", "Fs"), jen.If(jen.ID("s").Dot("config").Dot("CacheStaticFiles")).Body(
			jen.ID("afs").Equals().Qual("github.com/spf13/afero", "NewMemMapFs").Call(),
			jen.List(jen.ID("files"), jen.Err()).Assign().Qual("io/ioutil", "ReadDir").Call(jen.ID("fileDir")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("reading directory for frontend files: %w"), jen.Err())),
			),
			jen.Line(),
			jen.For(jen.List(jen.Underscore(), jen.ID("file")).Assign().Range().ID("files")).Body(
				jen.If(jen.ID("file").Dot("IsDir").Call()).Body(
					jen.Continue(),
				),
				jen.Line(),
				jen.ID("fp").Assign().Qual("path/filepath", "Join").Call(jen.ID("fileDir"), jen.ID("file").Dot("Name").Call()),
				jen.List(jen.ID("f"), jen.Err()).Assign().ID("afs").Dot("Create").Call(jen.ID("fp")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("creating static file in memory: %w"), jen.Err())),
				),
				jen.Line(),
				jen.List(jen.ID("bs"), jen.Err()).Assign().Qual("io/ioutil", "ReadFile").Call(jen.ID("fp")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("reading static file from directory: %w"), jen.Err())),
				),
				jen.Line(),
				jen.If(jen.List(jen.Underscore(), jen.Err()).Equals().ID("f").Dot("Write").Call(jen.ID("bs")), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("loading static file into memory: %w"), jen.Err())),
				),
				jen.Line(),
				jen.If(jen.Err().Equals().ID("f").Dot("Close").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("closing file while setting up static dir")),
				),
			),
			jen.ID("afs").Equals().Qual("github.com/spf13/afero", "NewReadOnlyFs").Call(jen.ID("afs")),
		).Else().Body(
			jen.ID("afs").Equals().Qual("github.com/spf13/afero", "NewOsFs").Call(),
		),
			jen.Line(),
			jen.Return().List(jen.Qual("github.com/spf13/afero", "NewHttpFs").Call(jen.ID("afs")), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildFrontendVarDeclarations(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			func() []jen.Code {
				pairs := []jen.Code{
					jen.Comment("Here is where you should put route regexes that need to be ignored by the static file server."),
					jen.Comment("For instance, if you allow someone to see an event in the frontend via a URL that contains dynamic."),
					jen.Comment("information, such as `/event/123`, you would want to put something like this below:"),
					jen.Comment("		eventsFrontendPathRegex = regexp.MustCompile(`/event/\\d+`)"),
					jen.Line(),
				}

				for _, typ := range proj.DataTypes {
					tuvn := typ.Name.PluralUnexportedVarName()
					pairs = append(pairs,
						jen.Commentf("%sFrontendPathRegex matches URLs against our frontend router's specification for specific %s routes.", tuvn, typ.Name.SingularCommonName()),
						jen.IDf("%sFrontendPathRegex", tuvn).Equals().Qual("regexp", "MustCompile").Call(jen.RawStringf(`/%s/\d+`, typ.Name.PluralRouteName())),
					)
				}

				return pairs
			}()...,
		),
		jen.Line(),
	}

	return lines
}

func buildFrontendStaticDir(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("StaticDir builds a static directory handler."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("StaticDir").Params(jen.ID("staticFilesDirectory").String()).Params(jen.Qual("net/http", "HandlerFunc"), jen.Error()).Body(
			jen.List(jen.ID("fileDir"), jen.Err()).Assign().Qual("path/filepath", "Abs").Call(jen.ID("staticFilesDirectory")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("determining absolute path of static files directory: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("httpFs"), jen.Err()).Assign().ID("s").Dot("buildStaticFileServer").Call(jen.ID("fileDir")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("establishing static server filesystem: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("s").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("static_dir"), jen.ID("fileDir")).Dot("Debug").Call(jen.Lit("setting static file server")),
			jen.ID("fs").Assign().Qual("net/http", "StripPrefix").Call(jen.Lit("/"), jen.Qual("net/http", "FileServer").Call(jen.ID("httpFs").Dot("Dir").Call(jen.ID("fileDir")))),
			jen.Line(),
			jen.Return().List(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Body(
				func() []jen.Code {
					lines := []jen.Code{
						jen.ID("rl").Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
						jen.ID("rl").Dot("Debug").Call(jen.Lit("static file requested")),
						jen.Switch(jen.ID(constants.RequestVarName).Dot("URL").Dot("Path")).Body(
							jen.Comment("list your frontend history routes here."),
							jen.Caseln(
								func() []jen.Code {
									lines := []jen.Code{
										jen.Lit("/register"),
										jen.Lit("/login"),
									}

									for _, typ := range proj.DataTypes {
										prn := typ.Name.PluralRouteName()
										lines = append(lines, jen.Litf("/%s", prn), jen.Litf("/%s/new", prn))
									}

									lines = append(lines, jen.Lit("/password/new"))

									return lines
								}()...,
							).Body(jen.ID("rl").Dot("Debug").Call(jen.Lit("rerouting")),
								jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Lit("/"),
							),
						),
					}

					for _, typ := range proj.DataTypes {
						tpuvn := typ.Name.PluralUnexportedVarName()
						lines = append(lines,
							jen.If(jen.IDf("%sFrontendPathRegex", tpuvn).Dot("MatchString").Call(jen.ID(constants.RequestVarName).Dot("URL").Dot("Path"))).Body(
								jen.ID("rl").Dot("Debug").Call(jen.Litf("rerouting %s request", typ.Name.SingularCommonName())),
								jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Lit("/"),
							),
						)
					}

					lines = append(lines,
						jen.Line(),
						jen.ID("fs").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
					)

					return lines
				}()...,
			), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
