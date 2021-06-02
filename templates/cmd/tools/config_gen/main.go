package config_gen

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	searchIndices := []jen.Code{
		jen.Comment("search index paths."),
	}
	defaultSearchIndexCount := len(searchIndices)
	for _, typ := range proj.DataTypes {
		searchIndices = append(searchIndices, jen.IDf("default%sSearchIndexPath", typ.Name.Plural()).Equals().Litf("%s.bleve", typ.Name.PluralRouteName()))
	}

	code.Add(
		jen.Const().Defs(
			jen.ID("defaultPort").Equals().Lit(8888),
			jen.ID("defaultCookieDomain").Equals().Lit("localhost"),
			jen.ID("debugCookieSecret").Equals().Lit("HEREISA32CHARSECRETWHICHISMADEUP"),
			jen.ID("devPostgresDBConnDetails").Equals().Lit("postgres://dbuser:hunter2@database:5432/todo?sslmode=disable"),
			jen.ID("devSqliteConnDetails").Equals().Lit("/tmp/db"),
			jen.ID("devMariaDBConnDetails").Equals().Lit("dbuser:hunter2@tcp(database:3306)/todo"),
			jen.ID("defaultCookieName").Equals().Qual(proj.AuthServicePackage(), "DefaultCookieName"),
			jen.Line(),
			jen.Comment("run modes."),
			jen.ID("developmentEnv").Equals().Lit("development"),
			jen.ID("testingEnv").Equals().Lit("testing"),
			jen.Line(),
			jen.Comment("database providers."),
			jen.ID("postgres").Equals().Lit("postgres"),
			jen.ID("sqlite").Equals().Lit("sqlite"),
			jen.ID("mariadb").Equals().Lit("mariadb"),
			jen.Line(),
			jen.Comment("test user stuff."),
			jen.ID("defaultPassword").Equals().Lit("password"),
			jen.Line(),
			func() jen.Code {
				if len(searchIndices) > defaultSearchIndexCount {
					return jen.Null().Add(utils.IntersperseWithNewlines(searchIndices)...)
				}
				return jen.Null()
			}(),
			jen.Line(),
			jen.ID("pasetoSecretSize").Equals().Lit(32),
			jen.ID("maxAttempts").Equals().Lit(50),
			jen.ID("defaultPASETOLifetime").Equals().Lit(1).Op("*").Qual("time", "Minute"),
			jen.Line(),
			jen.ID("contentTypeJSON").Equals().Lit("application/json"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("examplePASETOKey").Equals().ID("generatePASETOKey").Call(),
			jen.Line(),
			jen.ID("noopTracingConfig").Equals().Qual(proj.InternalTracingPackage(), "Config").Valuesln(
				jen.ID("Provider").MapAssign().Lit(""),
				jen.ID("SpanCollectionProbability").MapAssign().Lit(1),
			),
			jen.Line(),
			jen.ID("localServer").Equals().Qual(proj.HTTPServerPackage(), "Config").Valuesln(
				jen.ID("Debug").MapAssign().ID("true"),
				jen.ID("HTTPPort").MapAssign().ID("defaultPort"),
				jen.ID("StartupDeadline").MapAssign().Qual("time", "Minute"),
			),
			jen.Line(),
			jen.ID("localCookies").Equals().Qual(proj.AuthServicePackage(), "CookieConfig").Valuesln(
				jen.ID("Name").MapAssign().ID("defaultCookieName"),
				jen.ID("Domain").MapAssign().ID("defaultCookieDomain"),
				jen.ID("HashKey").MapAssign().ID("debugCookieSecret"),
				jen.ID("SigningKey").MapAssign().ID("debugCookieSecret"),
				jen.ID("Lifetime").MapAssign().Qual(proj.AuthServicePackage(), "DefaultCookieLifetime"),
				jen.ID("SecureOnly").MapAssign().ID("false"),
			),
			jen.Line(),
			jen.ID("localTracingConfig").Equals().Qual(proj.InternalTracingPackage(), "Config").Valuesln(
				jen.ID("Provider").MapAssign().Lit("jaeger"),
				jen.ID("SpanCollectionProbability").MapAssign().Lit(1),
				jen.ID("Jaeger").MapAssign().Op("&").Qual(proj.InternalTracingPackage(), "JaegerConfig").Valuesln(
					jen.ID("CollectorEndpoint").MapAssign().Lit("http://tracing-server:14268/api/traces"),
					jen.ID("ServiceName").MapAssign().Lit("todo_service"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("initializeLocalSecretManager").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Qual(proj.InternalSecretsPackage(), "SecretManager")).Body(
			jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
			jen.Line(),
			jen.ID("cfg").Op(":=").Op("&").Qual(proj.InternalSecretsPackage(), "Config").Valuesln(
				jen.ID("Provider").MapAssign().Qual(proj.InternalSecretsPackage(), "ProviderLocal"),
				jen.ID("Key").MapAssign().Lit("SUFNQVdBUkVUSEFUVEhJU1NFQ1JFVElTVU5TRUNVUkU="),
			),
			jen.Line(),
			jen.List(jen.ID("k"), jen.ID("err")).Op(":=").Qual(proj.InternalSecretsPackage(), "ProvideSecretKeeper").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.List(jen.ID("sm"), jen.ID("err")).Op(":=").Qual(proj.InternalSecretsPackage(), "ProvideSecretManager").Call(
				jen.ID("logger"),
				jen.ID("k"),
			),
			jen.If(jen.ID("err").DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("sm"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("encryptAndSaveConfig").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("outputPath").ID("string"),
			jen.ID("cfg").Op("*").Qual(proj.ConfigPackage(), "InstanceConfig")).Params(jen.ID("error")).Body(
			jen.ID("sm").Op(":=").ID("initializeLocalSecretManager").Call(jen.ID("ctx")),
			jen.List(jen.ID("output"), jen.ID("err")).Op(":=").ID("sm").Dot("Encrypt").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			),
			jen.If(jen.ID("err").DoesNotEqual().ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("encrypting config: %v"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().Qual("os", "WriteFile").Call(
				jen.ID("outputPath"),
				jen.Index().ID("byte").Call(jen.ID("output")),
				jen.Octal(644),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("configFunc").Func().Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filePath").ID("string")).Params(jen.ID("error")),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("files").Equals().Map(jen.ID("string")).ID("configFunc").Valuesln(
			jen.Lit("environments/local/service.config").MapAssign().ID("localDevelopmentConfig"),
			jen.Lit("environments/testing/config_files/frontend-tests.config").MapAssign().ID("frontendTestsConfig"),
			jen.Lit("environments/testing/config_files/integration-tests-postgres.config").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(
				jen.ID("postgres"),
				jen.ID("devPostgresDBConnDetails"),
			),
			jen.Lit("environments/testing/config_files/integration-tests-sqlite.config").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(
				jen.ID("sqlite"),
				jen.ID("devSqliteConnDetails"),
			),
			jen.Lit("environments/testing/config_files/integration-tests-mariadb.config").MapAssign().ID("buildIntegrationTestForDBImplementation").Call(
				jen.ID("mariadb"),
				jen.ID("devMariaDBConnDetails"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildLocalFrontendServiceConfig").Params().Params(jen.Qual(proj.FrontendServicePackage(), "Config")).Body(
			jen.Return().Qual(proj.FrontendServicePackage(), "Config").Valuesln(
				jen.ID("UseFakeData").MapAssign().ID("false"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("mustHashPass").Params(jen.ID("password").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("hashed"), jen.ID("err")).Op(":=").Qual(proj.InternalAuthenticationPackage(), "ProvideArgon2Authenticator").Call(jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call()).
				Dotln("HashPassword").Call(jen.Qual("context", "Background").Call(), jen.ID("password")),
			jen.Line(),
			jen.If(jen.ID("err").DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("hashed"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("generatePASETOKey").Params().Params(jen.Index().ID("byte")).Body(
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("pasetoSecretSize"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")),
				jen.ID("err").DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("b"),
		),
		jen.Line(),
	)

	serviceConfigs := []jen.Code{}
	for _, typ := range proj.DataTypes {
		serviceConfigs = append(serviceConfigs, jen.ID(typ.Name.Plural()).MapAssign().Qual(proj.ServicePackage(typ.Name.PackageName()), "Config").Valuesln(
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("SearchIndexPath").MapAssign().Qual("fmt", "Sprintf").Call(jen.Lit("/search_indices/%s"), jen.IDf("default%sSearchIndexPath", typ.Name.Plural()))
				}
				return jen.Null()
			}(),
			jen.ID("Logger").MapAssign().Qual(proj.InternalLoggingPackage(), "Config").Valuesln(
				jen.ID("Name").MapAssign().Lit(typ.Name.PluralRouteName()),
				jen.ID("Level").MapAssign().Qual(proj.InternalLoggingPackage(), "InfoLevel"),
				jen.ID("Provider").MapAssign().Qual(proj.InternalLoggingPackage(), "ProviderZerolog"),
			),
		))
	}

	code.Add(
		jen.Func().ID("localDevelopmentConfig").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("filePath").ID("string")).Params(jen.ID("error")).Body(
			jen.ID("cfg").Op(":=").Op("&").Qual(proj.ConfigPackage(), "InstanceConfig").Valuesln(
				jen.ID("Meta").MapAssign().Qual(proj.ConfigPackage(), "MetaSettings").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("RunMode").MapAssign().ID("developmentEnv"),
				),
				jen.ID("Encoding").MapAssign().Qual(proj.EncodingPackage(), "Config").Valuesln(
					jen.ID("ContentType").MapAssign().ID("contentTypeJSON"),
				),
				jen.ID("Server").MapAssign().ID("localServer"),
				jen.ID("Database").MapAssign().Qual(proj.DatabasePackage("config"), "Config").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("RunMigrations").MapAssign().ID("true"),
					jen.ID("MaxPingAttempts").MapAssign().ID("maxAttempts"),
					jen.ID("Provider").MapAssign().ID("postgres"),
					jen.ID("ConnectionDetails").MapAssign().ID("devPostgresDBConnDetails"),
					jen.ID("MetricsCollectionInterval").MapAssign().Qual("time", "Second"),
					jen.ID("CreateTestUser").MapAssign().Op("&").Qual(proj.TypesPackage(), "TestUserCreationConfig").Valuesln(
						jen.ID("Username").MapAssign().Lit("username"),
						jen.ID("Password").MapAssign().ID("defaultPassword"),
						jen.ID("HashedPassword").MapAssign().ID("mustHashPass").Call(jen.ID("defaultPassword")),
						jen.ID("IsServiceAdmin").MapAssign().ID("true"),
					),
				),
				jen.ID("Observability").MapAssign().Qual(proj.ObservabilityPackage(), "Config").Valuesln(
					jen.ID("Metrics").MapAssign().Qual(proj.MetricsPackage(), "Config").Valuesln(
						jen.ID("Provider").MapAssign().Lit("prometheus"),
						jen.ID("RouteToken").MapAssign().Lit(""),
						jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Qual("time", "Second"),
					),
					jen.ID("Tracing").MapAssign().ID("localTracingConfig"),
				),
				jen.ID("Uploads").MapAssign().Qual(proj.UploadsPackage(), "Config").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("Storage").MapAssign().Qual(proj.StoragePackage(), "Config").Valuesln(
						jen.ID("UploadFilenameKey").MapAssign().Lit("avatar"),
						jen.ID("Provider").MapAssign().Lit("filesystem"),
						jen.ID("BucketName").MapAssign().Lit("avatars"),
						jen.ID("AzureConfig").MapAssign().ID("nil"),
						jen.ID("GCSConfig").MapAssign().ID("nil"),
						jen.ID("S3Config").MapAssign().ID("nil"),
						jen.ID("FilesystemConfig").MapAssign().Op("&").Qual(proj.StoragePackage(), "FilesystemConfig").Valuesln(
							jen.ID("RootDirectory").MapAssign().Lit("/avatars"),
						),
					),
				),
				jen.ID("Search").MapAssign().Qual(proj.InternalSearchPackage(), "Config").Valuesln(
					jen.ID("Provider").MapAssign().Lit("bleve"),
				),
				jen.ID("Services").MapAssign().Qual(proj.ConfigPackage(), "ServicesConfigurations").Valuesln(
					append([]jen.Code{
						jen.ID("Auth").MapAssign().Qual(proj.AuthServicePackage(), "Config").Valuesln(
							jen.ID("PASETO").MapAssign().Qual(proj.AuthServicePackage(), "PASETOConfig").Valuesln(
								jen.ID("Issuer").MapAssign().Lit("todo_service"),
								jen.ID("Lifetime").MapAssign().ID("defaultPASETOLifetime"),
								jen.ID("LocalModeKey").MapAssign().ID("examplePASETOKey"),
							),
							jen.ID("Cookies").MapAssign().ID("localCookies"),
							jen.ID("Debug").MapAssign().ID("true"),
							jen.ID("EnableUserSignup").MapAssign().ID("true"),
							jen.ID("MinimumUsernameLength").MapAssign().Lit(4),
							jen.ID("MinimumPasswordLength").MapAssign().Lit(8),
						),
						jen.ID("Frontend").MapAssign().ID("buildLocalFrontendServiceConfig").Call(),
						jen.ID("Webhooks").MapAssign().Qual(proj.WebhooksServicePackage(), "Config").Valuesln(
							jen.ID("Debug").MapAssign().ID("true"),
							jen.ID("Enabled").MapAssign().ID("false"),
						),
					},
						serviceConfigs...,
					)...,
				),
				jen.ID("AuditLog").MapAssign().Qual(proj.AuditServicePackage(), "Config").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("Enabled").MapAssign().ID("true"),
				),
			),
			jen.Line(),
			jen.Return().ID("encryptAndSaveConfig").Call(
				jen.ID("ctx"),
				jen.ID("filePath"),
				jen.ID("cfg"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("frontendTestsConfig").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("filePath").ID("string")).Params(jen.ID("error")).Body(
			jen.ID("cfg").Op(":=").Op("&").Qual(proj.ConfigPackage(), "InstanceConfig").Valuesln(
				jen.ID("Meta").MapAssign().Qual(proj.ConfigPackage(), "MetaSettings").Valuesln(
					jen.ID("Debug").MapAssign().ID("false"),
					jen.ID("RunMode").MapAssign().ID("developmentEnv"),
				),
				jen.ID("Encoding").MapAssign().Qual(proj.EncodingPackage(), "Config").Valuesln(
					jen.ID("ContentType").MapAssign().ID("contentTypeJSON"),
				),
				jen.ID("Server").MapAssign().ID("localServer"),
				jen.ID("Database").MapAssign().Qual(proj.DatabasePackage("config"), "Config").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("RunMigrations").MapAssign().ID("true"),
					jen.ID("Provider").MapAssign().ID("postgres"),
					jen.ID("ConnectionDetails").MapAssign().ID("devPostgresDBConnDetails"),
					jen.ID("MaxPingAttempts").MapAssign().ID("maxAttempts"),
					jen.ID("MetricsCollectionInterval").MapAssign().Qual("time", "Second"),
				),
				jen.ID("Observability").MapAssign().Qual(proj.ObservabilityPackage(), "Config").Valuesln(
					jen.ID("Metrics").MapAssign().Qual(proj.MetricsPackage(), "Config").Valuesln(
						jen.ID("Provider").MapAssign().Lit("prometheus"),
						jen.ID("RouteToken").MapAssign().Lit(""),
						jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Qual("time", "Second"),
					),
					jen.ID("Tracing").MapAssign().ID("noopTracingConfig"),
				),
				jen.ID("Uploads").MapAssign().Qual(proj.UploadsPackage(), "Config").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("Storage").MapAssign().Qual(proj.StoragePackage(), "Config").Valuesln(
						jen.ID("UploadFilenameKey").MapAssign().Lit("avatar"),
						jen.ID("Provider").MapAssign().Lit("memory"),
						jen.ID("BucketName").MapAssign().Lit("avatars"),
					),
				),
				jen.ID("Search").MapAssign().Qual(proj.InternalSearchPackage(), "Config").Valuesln(
					jen.ID("Provider").MapAssign().Lit("bleve"),
				),
				jen.ID("Services").MapAssign().Qual(proj.ConfigPackage(), "ServicesConfigurations").Valuesln(
					append([]jen.Code{
						jen.ID("Auth").MapAssign().Qual(proj.AuthServicePackage(), "Config").Valuesln(
							jen.ID("PASETO").MapAssign().Qual(proj.AuthServicePackage(), "PASETOConfig").Valuesln(
								jen.ID("Issuer").MapAssign().Lit("todo_service"),
								jen.ID("Lifetime").MapAssign().ID("defaultPASETOLifetime"),
								jen.ID("LocalModeKey").MapAssign().ID("examplePASETOKey"),
							),
							jen.ID("Cookies").MapAssign().ID("localCookies"),
							jen.ID("Debug").MapAssign().ID("true"),
							jen.ID("EnableUserSignup").MapAssign().ID("true"),
							jen.ID("MinimumUsernameLength").MapAssign().Lit(4),
							jen.ID("MinimumPasswordLength").MapAssign().Lit(8),
						),
						jen.ID("Frontend").MapAssign().ID("buildLocalFrontendServiceConfig").Call(),
						jen.ID("Webhooks").MapAssign().Qual(proj.WebhooksServicePackage(), "Config").Valuesln(
							jen.ID("Debug").MapAssign().ID("true"),
							jen.ID("Enabled").MapAssign().ID("false"),
						),
					}, serviceConfigs...)...,
				),
				jen.ID("AuditLog").MapAssign().Qual(proj.AuditServicePackage(), "Config").Valuesln(
					jen.ID("Debug").MapAssign().ID("true"),
					jen.ID("Enabled").MapAssign().ID("true"),
				),
			),
			jen.Line(),
			jen.Return().ID("encryptAndSaveConfig").Call(
				jen.ID("ctx"),
				jen.ID("filePath"),
				jen.ID("cfg"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildIntegrationTestForDBImplementation").Params(jen.List(jen.ID("dbVendor"), jen.ID("dbDetails")).ID("string")).Params(jen.ID("configFunc")).Body(
			jen.Return().Func().Params(jen.ID("ctx").Qual("context", "Context"),
				jen.ID("filePath").ID("string")).Params(jen.ID("error")).Body(
				jen.ID("startupDeadline").Op(":=").Qual("time", "Minute"),
				jen.If(jen.ID("dbVendor").Op("==").ID("mariadb")).Body(
					jen.ID("startupDeadline").Equals().Lit(5).Op("*").Qual("time", "Minute"),
				),
				jen.Line(),
				jen.ID("cfg").Op(":=").Op("&").Qual(proj.ConfigPackage(), "InstanceConfig").Valuesln(
					jen.ID("Meta").MapAssign().Qual(proj.ConfigPackage(), "MetaSettings").Valuesln(
						jen.ID("Debug").MapAssign().ID("false"),
						jen.ID("RunMode").MapAssign().ID("testingEnv"),
					),
					jen.ID("Encoding").MapAssign().Qual(proj.EncodingPackage(), "Config").Valuesln(
						jen.ID("ContentType").MapAssign().ID("contentTypeJSON"),
					),
					jen.ID("Server").MapAssign().Qual(proj.HTTPServerPackage(), "Config").Valuesln(
						jen.ID("Debug").MapAssign().ID("false"),
						jen.ID("HTTPPort").MapAssign().ID("defaultPort"),
						jen.ID("StartupDeadline").MapAssign().ID("startupDeadline"),
					),
					jen.ID("Database").MapAssign().Qual(proj.DatabasePackage("config"), "Config").Valuesln(
						jen.ID("Debug").MapAssign().ID("false"),
						jen.ID("RunMigrations").MapAssign().ID("true"),
						jen.ID("Provider").MapAssign().ID("dbVendor"),
						jen.ID("MaxPingAttempts").MapAssign().ID("maxAttempts"),
						jen.ID("MetricsCollectionInterval").MapAssign().Lit(2).Op("*").Qual("time", "Second"),
						jen.ID("ConnectionDetails").MapAssign().Qual(proj.DatabasePackage(), "ConnectionDetails").Call(jen.ID("dbDetails")),
						jen.ID("CreateTestUser").MapAssign().Op("&").Qual(proj.TypesPackage(), "TestUserCreationConfig").Valuesln(
							jen.ID("Username").MapAssign().Lit("exampleUser"),
							jen.ID("Password").MapAssign().Lit("integration-tests-are-cool"),
							jen.ID("HashedPassword").MapAssign().ID("mustHashPass").Call(jen.Lit("integration-tests-are-cool")),
							jen.ID("IsServiceAdmin").MapAssign().ID("true"),
						),
					),
					jen.ID("Observability").MapAssign().Qual(proj.ObservabilityPackage(), "Config").Valuesln(
						jen.ID("Metrics").MapAssign().Qual(proj.MetricsPackage(), "Config").Valuesln(
							jen.ID("Provider").MapAssign().Lit(""),
							jen.ID("RouteToken").MapAssign().Lit(""),
							jen.ID("RuntimeMetricsCollectionInterval").MapAssign().Qual("time", "Second"),
						),
						jen.ID("Tracing").MapAssign().ID("localTracingConfig"),
					),
					jen.ID("Uploads").MapAssign().Qual(proj.UploadsPackage(), "Config").Valuesln(
						jen.ID("Debug").MapAssign().ID("false"),
						jen.ID("Storage").MapAssign().Qual(proj.StoragePackage(), "Config").Valuesln(
							jen.ID("Provider").MapAssign().Lit("memory"),
							jen.ID("BucketName").MapAssign().Lit("avatars"),
							jen.ID("AzureConfig").MapAssign().ID("nil"),
							jen.ID("GCSConfig").MapAssign().ID("nil"),
							jen.ID("S3Config").MapAssign().ID("nil"),
						),
					),
					jen.ID("Search").MapAssign().Qual(proj.InternalSearchPackage(), "Config").Valuesln(
						jen.ID("Provider").MapAssign().Lit("bleve"),
					),
					jen.ID("Services").MapAssign().Qual(proj.ConfigPackage(), "ServicesConfigurations").Valuesln(
						append([]jen.Code{
							jen.ID("Auth").MapAssign().Qual(proj.AuthServicePackage(), "Config").Valuesln(
								jen.ID("PASETO").MapAssign().Qual(proj.AuthServicePackage(), "PASETOConfig").Valuesln(
									jen.ID("Issuer").MapAssign().Lit("todo_service"),
									jen.ID("Lifetime").MapAssign().ID("defaultPASETOLifetime"),
									jen.ID("LocalModeKey").MapAssign().ID("examplePASETOKey"),
								),
								jen.ID("Cookies").MapAssign().Qual(proj.AuthServicePackage(), "CookieConfig").Valuesln(
									jen.ID("Name").MapAssign().ID("defaultCookieName"),
									jen.ID("Domain").MapAssign().ID("defaultCookieDomain"),
									jen.ID("SigningKey").MapAssign().ID("debugCookieSecret"),
									jen.ID("Lifetime").MapAssign().Qual(proj.AuthServicePackage(), "DefaultCookieLifetime"),
									jen.ID("SecureOnly").MapAssign().ID("false"),
								),
								jen.ID("Debug").MapAssign().ID("false"),
								jen.ID("EnableUserSignup").MapAssign().ID("true"),
								jen.ID("MinimumUsernameLength").MapAssign().Lit(4),
								jen.ID("MinimumPasswordLength").MapAssign().Lit(8),
							),
							jen.ID("Frontend").MapAssign().ID("buildLocalFrontendServiceConfig").Call(),
							jen.ID("Webhooks").MapAssign().Qual(proj.WebhooksServicePackage(), "Config").Valuesln(
								jen.ID("Debug").MapAssign().ID("true"),
								jen.ID("Enabled").MapAssign().ID("false"),
							),
						}, serviceConfigs...)...,
					),
					jen.ID("AuditLog").MapAssign().Qual(proj.AuditServicePackage(), "Config").Valuesln(
						jen.ID("Debug").MapAssign().ID("false"),
						jen.ID("Enabled").MapAssign().ID("true"),
					),
				),
				jen.Line(),
				jen.Return().ID("encryptAndSaveConfig").Call(
					jen.ID("ctx"),
					jen.ID("filePath"),
					jen.ID("cfg"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("main").Params().Body(
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.Line(),
			jen.For(jen.List(jen.ID("filePath"), jen.ID("fun")).Op(":=").Range().ID("files")).Body(
				jen.If(jen.ID("err").Op(":=").ID("fun").Call(
					jen.ID("ctx"),
					jen.ID("filePath"),
				),
					jen.ID("err").DoesNotEqual().ID("nil")).Body(
					jen.Qual("log", "Fatalf").Call(
						jen.Lit("error rendering %s: %v"),
						jen.ID("filePath"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
