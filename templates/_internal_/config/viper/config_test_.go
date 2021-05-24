package viper

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuildViperConfig").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("actual").Op(":=").ID("BuildViperConfig").Call(),
			jen.ID("assert").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("actual"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestFromConfig").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("ServerConfig").Valuesln(
						jen.ID("Server").Op(":").ID("server").Dot("Config").Valuesln(
							jen.ID("HTTPPort").Op(":").Lit(1234), jen.ID("Debug").Op(":").ID("false"), jen.ID("StartupDeadline").Op(":").Qual("time", "Minute")), jen.ID("AuditLog").Op(":").ID("audit").Dot("Config").Valuesln(
							jen.ID("Enabled").Op(":").ID("true")), jen.ID("Meta").Op(":").ID("config").Dot("MetaSettings").Valuesln(
							jen.ID("RunMode").Op(":").ID("config").Dot("DevelopmentRunMode")), jen.ID("Encoding").Op(":").ID("encoding").Dot("Config").Valuesln(
							jen.ID("ContentType").Op(":").Lit("application/json")), jen.ID("Capitalism").Op(":").ID("capitalism").Dot("Config").Valuesln(
							jen.ID("Enabled").Op(":").ID("false"), jen.ID("Provider").Op(":").ID("capitalism").Dot("StripeProvider"), jen.ID("Stripe").Op(":").Op("&").ID("capitalism").Dot("StripeConfig").Valuesln(
								jen.ID("APIKey").Op(":").Lit("whatever"), jen.ID("SuccessURL").Op(":").Lit("whatever"), jen.ID("CancelURL").Op(":").Lit("whatever"), jen.ID("WebhookSecret").Op(":").Lit("whatever"))), jen.ID("Auth").Op(":").ID("authentication").Dot("Config").Valuesln(
							jen.ID("Cookies").Op(":").ID("authentication").Dot("CookieConfig").Valuesln(
								jen.ID("Name").Op(":").Lit("todocookie"), jen.ID("Domain").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Lifetime").Op(":").Qual("time", "Second")), jen.ID("MinimumUsernameLength").Op(":").Lit(4), jen.ID("MinimumPasswordLength").Op(":").Lit(8), jen.ID("EnableUserSignup").Op(":").ID("true")), jen.ID("Observability").Op(":").ID("observability").Dot("Config").Valuesln(
							jen.ID("Metrics").Op(":").ID("metrics").Dot("Config").Valuesln(
								jen.ID("Provider").Op(":").Lit(""), jen.ID("RouteToken").Op(":").Lit(""), jen.ID("RuntimeMetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second")), jen.ID("Tracing").Op(":").ID("tracing").Dot("Config").Valuesln(
								jen.ID("Jaeger").Op(":").Op("&").ID("tracing").Dot("JaegerConfig").Valuesln(
									jen.ID("CollectorEndpoint").Op(":").Lit("things"), jen.ID("ServiceName").Op(":").Lit("stuff")), jen.ID("Provider").Op(":").Lit("blah"), jen.ID("SpanCollectionProbability").Op(":").Lit(0))), jen.ID("Uploads").Op(":").ID("uploads").Dot("Config").Valuesln(
							jen.ID("Storage").Op(":").ID("storage").Dot("Config").Valuesln(
								jen.ID("FilesystemConfig").Op(":").Op("&").ID("storage").Dot("FilesystemConfig").Valuesln(
									jen.ID("RootDirectory").Op(":").Lit("/blah")), jen.ID("AzureConfig").Op(":").Op("&").ID("storage").Dot("AzureConfig").Valuesln(
									jen.ID("BucketName").Op(":").Lit("blahs"), jen.ID("Retrying").Op(":").Op("&").ID("storage").Dot("AzureRetryConfig").Values()), jen.ID("GCSConfig").Op(":").Op("&").ID("storage").Dot("GCSConfig").Valuesln(
									jen.ID("ServiceAccountKeyFilepath").Op(":").Lit("/blah/blah"), jen.ID("BucketName").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").ID("nil")), jen.ID("S3Config").Op(":").Op("&").ID("storage").Dot("S3Config").Valuesln(
									jen.ID("BucketName").Op(":").Lit("blahs")), jen.ID("BucketName").Op(":").Lit("blahs"), jen.ID("UploadFilenameKey").Op(":").Lit("blahs"), jen.ID("Provider").Op(":").Lit("blahs")), jen.ID("Debug").Op(":").ID("false")), jen.ID("Frontend").Op(":").ID("frontend").Dot("Config").Values(), jen.ID("Search").Op(":").ID("search").Dot("Config").Valuesln(
							jen.ID("ItemsIndexPath").Op(":").Lit("/items_index_path")), jen.ID("Database").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "Config").Valuesln(
							jen.ID("Provider").Op(":").Lit("postgres"), jen.ID("MetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"), jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("ConnectionDetails").Op(":").ID("database").Dot("ConnectionDetails").Call(jen.Lit("postgres://username:passwords@host/table")), jen.ID("CreateTestUser").Op(":").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(
								jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("HashedPassword").Op(":").Lit("hashashashashash"), jen.ID("IsServiceAdmin").Op(":").ID("false")))),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("FromConfig").Call(jen.ID("exampleConfig")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("FromConfig").Call(jen.ID("nil")),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("ServerConfig").Values(),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("FromConfig").Call(jen.ID("exampleConfig")),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestParseConfigFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("tf"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Qual("os", "TempDir").Call(),
						jen.Lit("*.json"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("filename").Op(":=").ID("tf").Dot("Name").Call(),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("ServerConfig").Valuesln(
						jen.ID("Server").Op(":").ID("server").Dot("Config").Valuesln(
							jen.ID("HTTPPort").Op(":").Lit(1234), jen.ID("Debug").Op(":").ID("false"), jen.ID("StartupDeadline").Op(":").Qual("time", "Minute")), jen.ID("AuditLog").Op(":").ID("audit").Dot("Config").Valuesln(
							jen.ID("Enabled").Op(":").ID("true")), jen.ID("Meta").Op(":").ID("config").Dot("MetaSettings").Valuesln(
							jen.ID("RunMode").Op(":").ID("config").Dot("DevelopmentRunMode")), jen.ID("Encoding").Op(":").ID("encoding").Dot("Config").Valuesln(
							jen.ID("ContentType").Op(":").Lit("application/json")), jen.ID("Auth").Op(":").ID("authentication").Dot("Config").Valuesln(
							jen.ID("Cookies").Op(":").ID("authentication").Dot("CookieConfig").Valuesln(
								jen.ID("Name").Op(":").Lit("todocookie"), jen.ID("Domain").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Lifetime").Op(":").Qual("time", "Second")), jen.ID("MinimumUsernameLength").Op(":").Lit(4), jen.ID("MinimumPasswordLength").Op(":").Lit(8), jen.ID("EnableUserSignup").Op(":").ID("true")), jen.ID("Observability").Op(":").ID("observability").Dot("Config").Valuesln(
							jen.ID("Metrics").Op(":").ID("metrics").Dot("Config").Valuesln(
								jen.ID("Provider").Op(":").Lit(""), jen.ID("RouteToken").Op(":").Lit(""), jen.ID("RuntimeMetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"))), jen.ID("Frontend").Op(":").ID("frontend").Dot("Config").Values(), jen.ID("Search").Op(":").ID("search").Dot("Config").Valuesln(
							jen.ID("ItemsIndexPath").Op(":").Lit("/items_index_path")), jen.ID("Database").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "Config").Valuesln(
							jen.ID("Provider").Op(":").Lit("postgres"), jen.ID("MetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"), jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("ConnectionDetails").Op(":").ID("database").Dot("ConnectionDetails").Call(jen.Lit("postgres://username:passwords@host/table")))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("exampleConfig").Dot("EncodeToFile").Call(
							jen.ID("filename"),
							jen.Qual("encoding/json", "Marshal"),
						),
					),
					jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("filename"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleConfig"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tf").Dot("Name").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("unparseable garbage"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("tf"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Qual("os", "TempDir").Call(),
						jen.Lit("*.toml"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("tf").Dot("Write").Call(jen.Index().ID("byte").Call(jen.Lit(`
[server]
http_port = "blah"
debug = ":banana:"
`))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("tf").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tf").Dot("Name").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nonexistent file"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.Lit("/this/doesn't/even/exist/lol"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cfg"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with test user creation on production error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("tf"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Qual("os", "TempDir").Call(),
						jen.Lit("*.json"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("filename").Op(":=").ID("tf").Dot("Name").Call(),
					jen.ID("exampleConfig").Op(":=").Op("&").ID("config").Dot("ServerConfig").Valuesln(
						jen.ID("Server").Op(":").ID("server").Dot("Config").Valuesln(
							jen.ID("HTTPPort").Op(":").Lit(1234), jen.ID("Debug").Op(":").ID("false"), jen.ID("StartupDeadline").Op(":").Qual("time", "Minute")), jen.ID("AuditLog").Op(":").ID("audit").Dot("Config").Valuesln(
							jen.ID("Enabled").Op(":").ID("true")), jen.ID("Meta").Op(":").ID("config").Dot("MetaSettings").Valuesln(
							jen.ID("RunMode").Op(":").ID("config").Dot("ProductionRunMode")), jen.ID("Encoding").Op(":").ID("encoding").Dot("Config").Valuesln(
							jen.ID("ContentType").Op(":").Lit("application/json")), jen.ID("Auth").Op(":").ID("authentication").Dot("Config").Valuesln(
							jen.ID("Cookies").Op(":").ID("authentication").Dot("CookieConfig").Valuesln(
								jen.ID("Name").Op(":").Lit("todocookie"), jen.ID("Domain").Op(":").Lit("https://verygoodsoftwarenotvirus.ru"), jen.ID("Lifetime").Op(":").Qual("time", "Second")), jen.ID("MinimumUsernameLength").Op(":").Lit(4), jen.ID("MinimumPasswordLength").Op(":").Lit(8), jen.ID("EnableUserSignup").Op(":").ID("true")), jen.ID("Observability").Op(":").ID("observability").Dot("Config").Valuesln(
							jen.ID("Metrics").Op(":").ID("metrics").Dot("Config").Valuesln(
								jen.ID("Provider").Op(":").Lit(""), jen.ID("RouteToken").Op(":").Lit(""), jen.ID("RuntimeMetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"))), jen.ID("Frontend").Op(":").ID("frontend").Dot("Config").Values(), jen.ID("Search").Op(":").ID("search").Dot("Config").Valuesln(
							jen.ID("ItemsIndexPath").Op(":").Lit("/items_index_path")), jen.ID("Database").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "Config").Valuesln(
							jen.ID("Provider").Op(":").Lit("postgres"), jen.ID("MetricsCollectionInterval").Op(":").Lit(2).Op("*").Qual("time", "Second"), jen.ID("Debug").Op(":").ID("true"), jen.ID("RunMigrations").Op(":").ID("true"), jen.ID("ConnectionDetails").Op(":").ID("database").Dot("ConnectionDetails").Call(jen.Lit("postgres://username:passwords@host/table")), jen.ID("CreateTestUser").Op(":").Op("&").ID("types").Dot("TestUserCreationConfig").Valuesln(
								jen.ID("Username").Op(":").Lit("username"), jen.ID("Password").Op(":").Lit("password"), jen.ID("HashedPassword").Op(":").Lit("blahblahblah"), jen.ID("IsServiceAdmin").Op(":").ID("false")))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("exampleConfig").Dot("EncodeToFile").Call(
							jen.ID("filename"),
							jen.Qual("encoding/json", "Marshal"),
						),
					),
					jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("filename"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cfg"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with validation error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("tf"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Qual("os", "TempDir").Call(),
						jen.Lit("*.toml"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("_"), jen.ID("err")).Op("=").ID("tf").Dot("Write").Call(jen.Index().ID("byte").Call(jen.Lit(`
[server]
http_port = 8888
`))),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("cfg"), jen.ID("err")).Op(":=").ID("ParseConfigFile").Call(
						jen.ID("ctx"),
						jen.ID("logger"),
						jen.ID("tf").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("cfg"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tf").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
