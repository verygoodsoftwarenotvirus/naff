package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("usersBasePath").Op("=").Lit("users"),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetUserRequest builds an HTTP request for fetching a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetUserRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.ID("id").Call(jen.ID("userID")),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetUsersRequest builds an HTTP request for fetching a list of users."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetUsersRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("b").Dot("logger")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("usersBasePath"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildSearchForUsersByUsernameRequest builds an HTTP request that searches for a user by their username."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildSearchForUsersByUsernameRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("username").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrEmptyUsernameProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("username"),
			),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("username"),
			),
			jen.ID("u").Op(":=").ID("b").Dot("buildAPIV1URL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("search"),
			),
			jen.ID("q").Op(":=").ID("u").Dot("Query").Call(),
			jen.ID("q").Dot("Set").Call(
				jen.ID("types").Dot("SearchQueryKey"),
				jen.ID("username"),
			),
			jen.ID("u").Dot("RawQuery").Op("=").ID("q").Dot("Encode").Call(),
			jen.ID("uri").Op(":=").ID("u").Dot("String").Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCreateUserRequest builds an HTTP request for creating a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildCreateUserRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserRegistrationInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("Username"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildArchiveUserRequest builds an HTTP request for archiving a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildArchiveUserRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("buildAPIV1URL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.ID("id").Call(jen.ID("userID")),
			).Dot("String").Call(),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("BuildAvatarUploadRequest builds an HTTP request that sets a user's avatar to the provided content.").Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildAvatarUploadRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("avatar").Index().ID("byte"), jen.ID("extension").ID("string")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("len").Call(jen.ID("avatar")).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.Var().ID("ct").ID("string"),
			jen.Switch(jen.Qual("strings", "ToLower").Call(jen.Qual("strings", "TrimSpace").Call(jen.ID("extension")))).Body(
				jen.Case(jen.Lit("jpeg")).Body(
					jen.ID("ct").Op("=").Lit("image/jpeg")),
				jen.Case(jen.Lit("png")).Body(
					jen.ID("ct").Op("=").Lit("image/png")),
				jen.Case(jen.Lit("gif")).Body(
					jen.ID("ct").Op("=").Lit("image/gif")),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("%s: %w"),
						jen.ID("extension"),
						jen.ID("ErrInvalidPhotoEncodingForUpload"),
					))),
			),
			jen.ID("logger").Op(":=").ID("b").Dot("logger"),
			jen.ID("body").Op(":=").Op("&").Qual("bytes", "Buffer").Valuesln(),
			jen.ID("writer").Op(":=").Qual("mime/multipart", "NewWriter").Call(jen.ID("body")),
			jen.List(jen.ID("part"), jen.ID("err")).Op(":=").ID("writer").Dot("CreateFormFile").Call(
				jen.Lit("avatar"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("avatar.%s"),
					jen.ID("extension"),
				),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating form file"),
				))),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op("=").Qual("io", "Copy").Call(
				jen.ID("part"),
				jen.Qual("bytes", "NewReader").Call(jen.ID("avatar")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("copying file contents to request"),
				))),
			jen.If(jen.ID("err").Op("=").ID("writer").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("closing avatar writer"),
				))),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("avatar"),
				jen.Lit("upload"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building avatar upload request"),
				))),
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.Lit("Content-Type"),
				jen.ID("writer").Dot("FormDataContentType").Call(),
			),
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.Lit("X-Upload-Content-Type"),
				jen.ID("ct"),
			),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildGetAuditLogForUserRequest builds an HTTP request for fetching a list of audit log entries for a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildGetAuditLogForUserRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("BuildURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.ID("id").Call(jen.ID("userID")),
				jen.Lit("audit"),
			),
			jen.ID("tracing").Dot("AttachRequestURIToSpan").Call(
				jen.ID("span"),
				jen.ID("uri"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
