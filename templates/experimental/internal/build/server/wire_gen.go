package server

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireGenDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("Build builds a server."),
		jen.Line(),
		jen.Func().ID("Build").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("config").Dot("ServerConfig"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("dbm").ID("database").Dot("DataManager"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("authenticator").ID("authentication").Dot("Authenticator")).Params(jen.Op("*").ID("server").Dot("HTTPServer"), jen.ID("error")).Body(
			jen.ID("serverConfig").Op(":=").ID("cfg").Dot("Server"),
			jen.ID("observabilityConfig").Op(":=").Op("&").ID("cfg").Dot("Observability"),
			jen.ID("metricsConfig").Op(":=").Op("&").ID("observabilityConfig").Dot("Metrics"),
			jen.List(jen.ID("instrumentationHandler"), jen.ID("err")).Op(":=").ID("metrics").Dot("ProvideMetricsInstrumentationHandlerForServer").Call(
				jen.ID("metricsConfig"),
				jen.ID("logger"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("authenticationConfig").Op(":=").Op("&").ID("cfg").Dot("Auth"),
			jen.ID("userDataManager").Op(":=").ID("database").Dot("ProvideUserDataManager").Call(jen.ID("dbm")),
			jen.ID("authAuditManager").Op(":=").ID("database").Dot("ProvideAuthAuditManager").Call(jen.ID("dbm")),
			jen.ID("apiClientDataManager").Op(":=").ID("database").Dot("ProvideAPIClientDataManager").Call(jen.ID("dbm")),
			jen.ID("accountUserMembershipDataManager").Op(":=").ID("database").Dot("ProvideAccountUserMembershipDataManager").Call(jen.ID("dbm")),
			jen.ID("cookieConfig").Op(":=").ID("authenticationConfig").Dot("Cookies"),
			jen.ID("databaseConfig").Op(":=").ID("cfg").Dot("Database"),
			jen.List(jen.ID("sessionManager"), jen.ID("err")).Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "ProvideSessionManager").Call(
				jen.ID("cookieConfig"),
				jen.ID("databaseConfig"),
				jen.ID("db"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("encodingConfig").Op(":=").ID("cfg").Dot("Encoding"),
			jen.ID("contentType").Op(":=").ID("encoding").Dot("ProvideContentType").Call(jen.ID("encodingConfig")),
			jen.ID("serverEncoderDecoder").Op(":=").ID("encoding").Dot("ProvideServerEncoderDecoder").Call(
				jen.ID("logger"),
				jen.ID("contentType"),
			),
			jen.ID("routeParamManager").Op(":=").ID("chi").Dot("NewRouteParamManager").Call(),
			jen.List(jen.ID("authService"), jen.ID("err")).Op(":=").ID("authentication").Dot("ProvideService").Call(
				jen.ID("logger"),
				jen.ID("authenticationConfig"),
				jen.ID("authenticator"),
				jen.ID("userDataManager"),
				jen.ID("authAuditManager"),
				jen.ID("apiClientDataManager"),
				jen.ID("accountUserMembershipDataManager"),
				jen.ID("sessionManager"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("routeParamManager"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("auditLogEntryDataManager").Op(":=").ID("database").Dot("ProvideAuditLogEntryDataManager").Call(jen.ID("dbm")),
			jen.ID("auditLogEntryDataService").Op(":=").ID("audit").Dot("ProvideService").Call(
				jen.ID("logger"),
				jen.ID("auditLogEntryDataManager"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("routeParamManager"),
			),
			jen.ID("accountDataManager").Op(":=").ID("database").Dot("ProvideAccountDataManager").Call(jen.ID("dbm")),
			jen.List(jen.ID("unitCounterProvider"), jen.ID("err")).Op(":=").ID("metrics").Dot("ProvideUnitCounterProvider").Call(
				jen.ID("metricsConfig"),
				jen.ID("logger"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("imageUploadProcessor").Op(":=").ID("images").Dot("NewImageUploadProcessor").Call(jen.ID("logger")),
			jen.ID("uploadsConfig").Op(":=").Op("&").ID("cfg").Dot("Uploads"),
			jen.ID("storageConfig").Op(":=").Op("&").ID("uploadsConfig").Dot("Storage"),
			jen.List(jen.ID("uploader"), jen.ID("err")).Op(":=").ID("storage").Dot("NewUploadManager").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("storageConfig"),
				jen.ID("routeParamManager"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("uploadManager").Op(":=").ID("uploads").Dot("ProvideUploadManager").Call(jen.ID("uploader")),
			jen.ID("userDataService").Op(":=").ID("users").Dot("ProvideUsersService").Call(
				jen.ID("authenticationConfig"),
				jen.ID("logger"),
				jen.ID("userDataManager"),
				jen.ID("accountDataManager"),
				jen.ID("authenticator"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("unitCounterProvider"),
				jen.ID("imageUploadProcessor"),
				jen.ID("uploadManager"),
				jen.ID("routeParamManager"),
			),
			jen.ID("accountDataService").Op(":=").ID("accounts").Dot("ProvideService").Call(
				jen.ID("logger"),
				jen.ID("accountDataManager"),
				jen.ID("accountUserMembershipDataManager"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("unitCounterProvider"),
				jen.ID("routeParamManager"),
			),
			jen.ID("apiclientsConfig").Op(":=").ID("apiclients").Dot("ProvideConfig").Call(jen.ID("authenticationConfig")),
			jen.ID("apiClientDataService").Op(":=").ID("apiclients").Dot("ProvideAPIClientsService").Call(
				jen.ID("logger"),
				jen.ID("apiClientDataManager"),
				jen.ID("userDataManager"),
				jen.ID("authenticator"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("unitCounterProvider"),
				jen.ID("routeParamManager"),
				jen.ID("apiclientsConfig"),
			),
			jen.ID("itemDataManager").Op(":=").ID("database").Dot("ProvideItemDataManager").Call(jen.ID("dbm")),
			jen.ID("searchConfig").Op(":=").ID("cfg").Dot("Search"),
			jen.ID("indexManagerProvider").Op(":=").ID("bleve").Dot("ProvideBleveIndexManagerProvider").Call(),
			jen.List(jen.ID("itemDataService"), jen.ID("err")).Op(":=").ID("items").Dot("ProvideService").Call(
				jen.ID("logger"),
				jen.ID("itemDataManager"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("unitCounterProvider"),
				jen.ID("searchConfig"),
				jen.ID("indexManagerProvider"),
				jen.ID("routeParamManager"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("webhookDataManager").Op(":=").ID("database").Dot("ProvideWebhookDataManager").Call(jen.ID("dbm")),
			jen.ID("webhookDataService").Op(":=").ID("webhooks").Dot("ProvideWebhooksService").Call(
				jen.ID("logger"),
				jen.ID("webhookDataManager"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("unitCounterProvider"),
				jen.ID("routeParamManager"),
			),
			jen.ID("adminUserDataManager").Op(":=").ID("database").Dot("ProvideAdminUserDataManager").Call(jen.ID("dbm")),
			jen.ID("adminAuditManager").Op(":=").ID("database").Dot("ProvideAdminAuditManager").Call(jen.ID("dbm")),
			jen.ID("adminService").Op(":=").ID("admin").Dot("ProvideService").Call(
				jen.ID("logger"),
				jen.ID("authenticationConfig"),
				jen.ID("authenticator"),
				jen.ID("adminUserDataManager"),
				jen.ID("adminAuditManager"),
				jen.ID("sessionManager"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("routeParamManager"),
			),
			jen.ID("frontendConfig").Op(":=").Op("&").ID("cfg").Dot("Frontend"),
			jen.ID("frontendAuthService").Op(":=").ID("frontend").Dot("ProvideAuthService").Call(jen.ID("authService")),
			jen.ID("usersService").Op(":=").ID("frontend").Dot("ProvideUsersService").Call(jen.ID("userDataService")),
			jen.ID("capitalismConfig").Op(":=").Op("&").ID("cfg").Dot("Capitalism"),
			jen.ID("stripeConfig").Op(":=").ID("capitalismConfig").Dot("Stripe"),
			jen.ID("paymentManager").Op(":=").ID("stripe").Dot("NewStripePaymentManager").Call(
				jen.ID("logger"),
				jen.ID("stripeConfig"),
			),
			jen.ID("service").Op(":=").ID("frontend").Dot("ProvideService").Call(
				jen.ID("frontendConfig"),
				jen.ID("logger"),
				jen.ID("frontendAuthService"),
				jen.ID("usersService"),
				jen.ID("dbm"),
				jen.ID("routeParamManager"),
				jen.ID("paymentManager"),
			),
			jen.ID("router").Op(":=").ID("chi").Dot("NewRouter").Call(jen.ID("logger")),
			jen.List(jen.ID("httpServer"), jen.ID("err")).Op(":=").ID("server").Dot("ProvideHTTPServer").Call(
				jen.ID("ctx"),
				jen.ID("serverConfig"),
				jen.ID("instrumentationHandler"),
				jen.ID("authService"),
				jen.ID("auditLogEntryDataService"),
				jen.ID("userDataService"),
				jen.ID("accountDataService"),
				jen.ID("apiClientDataService"),
				jen.ID("itemDataService"),
				jen.ID("webhookDataService"),
				jen.ID("adminService"),
				jen.ID("service"),
				jen.ID("logger"),
				jen.ID("serverEncoderDecoder"),
				jen.ID("router"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.Return().List(jen.ID("httpServer"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
