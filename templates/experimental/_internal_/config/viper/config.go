package viper

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("maxPASETOLifetime").Op("=").Lit(10).Op("*").Qual("time", "Minute"),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("errNilInput").Op("=").Qual("errors", "New").Call(jen.Lit("nil input provided")),
		jen.ID("errInvalidTestUserRunModeConfiguration").Op("=").Qual("errors", "New").Call(jen.Lit("requested test user in production run mode")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildViperConfig is a constructor function that initializes a viper config."),
		jen.Line(),
		jen.Func().ID("BuildViperConfig").Params().Params(jen.Op("*").ID("viper").Dot("Viper")).Body(
			jen.ID("cfg").Op(":=").ID("viper").Dot("New").Call(),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyMetaRunMode"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config", "DefaultRunMode"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyServerStartupDeadline"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config", "DefaultStartupDeadline"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyEncodingContentType"),
				jen.Lit("application/json"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyAuthCookieDomain"),
				jen.ID("authentication").Dot("DefaultCookieDomain"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyAuthCookieLifetime"),
				jen.ID("authentication").Dot("DefaultCookieLifetime"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyAuthEnableUserSignup"),
				jen.ID("true"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyDatabaseRunMigrations"),
				jen.ID("true"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyAuthMinimumUsernameLength"),
				jen.Lit(4),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyAuthMinimumPasswordLength"),
				jen.Lit(8),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyDatabaseMetricsCollectionInterval"),
				jen.ID("metrics").Dot("DefaultMetricsCollectionInterval"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyMetricsRuntimeCollectionInterval"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/database/config", "DefaultMetricsCollectionInterval"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyObservabilityTracingSpanCollectionProbability"),
				jen.Lit(1),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyAuditLogEnabled"),
				jen.ID("true"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeySearchProvider"),
				jen.ID("search").Dot("BleveProvider"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyWebhooksEnabled"),
				jen.ID("false"),
			),
			jen.ID("cfg").Dot("SetDefault").Call(
				jen.ID("ConfigKeyServerHTTPPort"),
				jen.Lit(80),
			),
			jen.Return().ID("cfg"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("FromConfig returns a viper instance from a config struct."),
		jen.Line(),
		jen.Func().ID("FromConfig").Params(jen.ID("input").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config", "ServerConfig")).Params(jen.Op("*").ID("viper").Dot("Viper"), jen.ID("error")).Body(
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilInput"))),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("cfg").Op(":=").ID("BuildViperConfig").Call(),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyMetaDebug"),
				jen.ID("input").Dot("Meta").Dot("Debug"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyMetaRunMode"),
				jen.ID("string").Call(jen.ID("input").Dot("Meta").Dot("RunMode")),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyServerStartupDeadline"),
				jen.ID("input").Dot("Server").Dot("StartupDeadline"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyServerHTTPPort"),
				jen.ID("input").Dot("Server").Dot("HTTPPort"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyServerDebug"),
				jen.ID("input").Dot("Server").Dot("Debug"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyEncodingContentType"),
				jen.ID("input").Dot("Encoding").Dot("ContentType"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyFrontendUseFakeData"),
				jen.ID("input").Dot("Frontend").Dot("UseFakeData"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthDebug"),
				jen.ID("input").Dot("Auth").Dot("Debug"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthEnableUserSignup"),
				jen.ID("input").Dot("Auth").Dot("EnableUserSignup"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthMinimumUsernameLength"),
				jen.ID("input").Dot("Auth").Dot("MinimumUsernameLength"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthMinimumPasswordLength"),
				jen.ID("input").Dot("Auth").Dot("MinimumPasswordLength"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthCookieName"),
				jen.ID("input").Dot("Auth").Dot("Cookies").Dot("Name"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthCookieDomain"),
				jen.ID("input").Dot("Auth").Dot("Cookies").Dot("Domain"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthCookieHashKey"),
				jen.ID("input").Dot("Auth").Dot("Cookies").Dot("HashKey"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthCookieSigningKey"),
				jen.ID("input").Dot("Auth").Dot("Cookies").Dot("SigningKey"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthCookieLifetime"),
				jen.ID("input").Dot("Auth").Dot("Cookies").Dot("Lifetime"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthSecureCookiesOnly"),
				jen.ID("input").Dot("Auth").Dot("Cookies").Dot("SecureOnly"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyCapitalismEnabled"),
				jen.ID("input").Dot("Capitalism").Dot("Enabled"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyCapitalismProvider"),
				jen.ID("input").Dot("Capitalism").Dot("Provider"),
			),
			jen.If(jen.ID("input").Dot("Capitalism").Dot("Stripe").Op("!=").ID("nil")).Body(
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyCapitalismStripeAPIKey"),
					jen.ID("input").Dot("Capitalism").Dot("Stripe").Dot("APIKey"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyCapitalismStripeSuccessURL"),
					jen.ID("input").Dot("Capitalism").Dot("Stripe").Dot("SuccessURL"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyCapitalismStripeCancelURL"),
					jen.ID("input").Dot("Capitalism").Dot("Stripe").Dot("CancelURL"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyCapitalismStripeWebhookSecret"),
					jen.ID("input").Dot("Capitalism").Dot("Stripe").Dot("WebhookSecret"),
				),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthPASETOListener"),
				jen.ID("input").Dot("Auth").Dot("PASETO").Dot("Issuer"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthPASETOLifetimeKey"),
				jen.Qual("time", "Duration").Call(jen.Qual("math", "Min").Call(
					jen.ID("float64").Call(jen.ID("input").Dot("Auth").Dot("PASETO").Dot("Lifetime")),
					jen.ID("float64").Call(jen.ID("maxPASETOLifetime")),
				)),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuthPASETOLocalModeKey"),
				jen.ID("input").Dot("Auth").Dot("PASETO").Dot("LocalModeKey"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyMetricsProvider"),
				jen.ID("input").Dot("Observability").Dot("Metrics").Dot("Provider"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyObservabilityTracingProvider"),
				jen.ID("input").Dot("Observability").Dot("Tracing").Dot("Provider"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyObservabilityTracingSpanCollectionProbability"),
				jen.ID("input").Dot("Observability").Dot("Tracing").Dot("SpanCollectionProbability"),
			),
			jen.If(jen.ID("input").Dot("Observability").Dot("Tracing").Dot("Jaeger").Op("!=").ID("nil")).Body(
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyObservabilityTracingJaegerCollectorEndpoint"),
					jen.ID("input").Dot("Observability").Dot("Tracing").Dot("Jaeger").Dot("CollectorEndpoint"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyObservabilityTracingJaegerServiceName"),
					jen.ID("input").Dot("Observability").Dot("Tracing").Dot("Jaeger").Dot("ServiceName"),
				),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyMetricsRuntimeCollectionInterval"),
				jen.ID("input").Dot("Observability").Dot("Metrics").Dot("RuntimeMetricsCollectionInterval"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyDatabaseDebug"),
				jen.ID("input").Dot("Database").Dot("Debug"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyDatabaseProvider"),
				jen.ID("input").Dot("Database").Dot("Provider"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyDatabaseMaxPingAttempts"),
				jen.ID("input").Dot("Database").Dot("MaxPingAttempts"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyDatabaseConnectionDetails"),
				jen.ID("string").Call(jen.ID("input").Dot("Database").Dot("ConnectionDetails")),
			),
			jen.If(jen.ID("input").Dot("Database").Dot("CreateTestUser").Op("!=").ID("nil")).Body(
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyDatabaseCreateTestUserUsername"),
					jen.ID("input").Dot("Database").Dot("CreateTestUser").Dot("Username"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyDatabaseCreateTestUserPassword"),
					jen.ID("input").Dot("Database").Dot("CreateTestUser").Dot("Password"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyDatabaseCreateTestUserIsServiceAdmin"),
					jen.ID("input").Dot("Database").Dot("CreateTestUser").Dot("IsServiceAdmin"),
				),
				jen.ID("cfg").Dot("Set").Call(
					jen.ID("ConfigKeyDatabaseCreateTestUserHashedPassword"),
					jen.ID("input").Dot("Database").Dot("CreateTestUser").Dot("HashedPassword"),
				),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyDatabaseRunMigrations"),
				jen.ID("input").Dot("Database").Dot("RunMigrations"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyDatabaseMetricsCollectionInterval"),
				jen.ID("input").Dot("Database").Dot("MetricsCollectionInterval"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeySearchProvider"),
				jen.ID("input").Dot("Search").Dot("Provider"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyItemsSearchIndexPath"),
				jen.ID("string").Call(jen.ID("input").Dot("Search").Dot("ItemsIndexPath")),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyUploaderProvider"),
				jen.ID("input").Dot("Uploads").Dot("Storage").Dot("Provider"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyUploaderDebug"),
				jen.ID("input").Dot("Uploads").Dot("Debug"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyUploaderBucketName"),
				jen.ID("input").Dot("Uploads").Dot("Storage").Dot("BucketName"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyUploaderUploadFilename"),
				jen.ID("input").Dot("Uploads").Dot("Storage").Dot("UploadFilenameKey"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyAuditLogEnabled"),
				jen.ID("input").Dot("AuditLog").Dot("Enabled"),
			),
			jen.ID("cfg").Dot("Set").Call(
				jen.ID("ConfigKeyWebhooksEnabled"),
				jen.ID("input").Dot("Webhooks").Dot("Enabled"),
			),
			jen.Switch().Body(
				jen.Case(jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Op("!=").ID("nil")).Body(
					jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderProvider"),
						jen.Lit("azure"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureAuthMethod"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("AuthMethod"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureAccountName"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("AccountName"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureBucketName"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("BucketName"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureMaxTries"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("Retrying").Dot("MaxTries"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureTryTimeout"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("Retrying").Dot("TryTimeout"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureRetryDelay"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("Retrying").Dot("RetryDelay"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureMaxRetryDelay"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("Retrying").Dot("MaxRetryDelay"),
					), jen.If(jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Op("!=").ID("nil")).Body(
						jen.ID("cfg").Dot("Set").Call(
							jen.ID("ConfigKeyUploaderAzureRetryReadsFromSecondaryHost"),
							jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("Retrying").Dot("RetryReadsFromSecondaryHost"),
						)), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureTokenCredentialsInitialToken"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("TokenCredentialsInitialToken"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderAzureSharedKeyAccountKey"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("AzureConfig").Dot("SharedKeyAccountKey"),
					), jen.Fallthrough()),
				jen.Case(jen.ID("input").Dot("Uploads").Dot("Storage").Dot("GCSConfig").Op("!=").ID("nil")).Body(
					jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderProvider"),
						jen.Lit("gcs"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderGCSAccountKeyFilepath"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("GCSConfig").Dot("ServiceAccountKeyFilepath"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderGCSScopes"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("GCSConfig").Dot("Scopes"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderGCSBucketName"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("GCSConfig").Dot("BucketName"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderGCSGoogleAccessID"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("GCSConfig").Dot("BlobSettings").Dot("GoogleAccessID"),
					), jen.Fallthrough()),
				jen.Case(jen.ID("input").Dot("Uploads").Dot("Storage").Dot("S3Config").Op("!=").ID("nil")).Body(
					jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderProvider"),
						jen.Lit("s3"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderS3BucketName"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("S3Config").Dot("BucketName"),
					), jen.Fallthrough()),
				jen.Case(jen.ID("input").Dot("Uploads").Dot("Storage").Dot("FilesystemConfig").Op("!=").ID("nil")).Body(
					jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderProvider"),
						jen.Lit("filesystem"),
					), jen.ID("cfg").Dot("Set").Call(
						jen.ID("ConfigKeyUploaderFilesystemRootDirectory"),
						jen.ID("input").Dot("Uploads").Dot("Storage").Dot("FilesystemConfig").Dot("RootDirectory"),
					)),
			),
			jen.Return().List(jen.ID("cfg"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ParseConfigFile parses a configuration file."),
		jen.Line(),
		jen.Func().ID("ParseConfigFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("filePath").ID("string")).Params(jen.Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config", "ServerConfig"), jen.ID("error")).Body(
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("filepath"),
				jen.ID("filePath"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("parsing config file")),
			jen.ID("cfg").Op(":=").ID("BuildViperConfig").Call(),
			jen.ID("cfg").Dot("SetConfigFile").Call(jen.ID("filePath")),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("ReadInConfig").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("trying to read the config file: %w"),
					jen.ID("err"),
				))),
			jen.Var().ID("serverConfig").Op("*").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config", "ServerConfig"),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("Unmarshal").Call(jen.Op("&").ID("serverConfig")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("trying to unmarshal the config: %w"),
					jen.ID("err"),
				))),
			jen.If(jen.ID("serverConfig").Dot("Database").Dot("CreateTestUser").Op("!=").ID("nil").Op("&&").ID("serverConfig").Dot("Meta").Dot("RunMode").Op("==").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/config", "ProductionRunMode")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errInvalidTestUserRunModeConfiguration"))),
			jen.If(jen.ID("validationErr").Op(":=").ID("serverConfig").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("validationErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("validationErr"))),
			jen.Return().List(jen.ID("serverConfig"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
