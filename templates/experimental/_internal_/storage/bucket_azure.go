package storage

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bucketAzureDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("AzureProvider").Op("=").Lit("azure"),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AzureRetryConfig").Struct(
				jen.ID("RetryReadsFromSecondaryHost").ID("string"),
				jen.ID("TryTimeout").Qual("time", "Duration"),
				jen.ID("RetryDelay").Qual("time", "Duration"),
				jen.ID("MaxRetryDelay").Qual("time", "Duration"),
				jen.ID("MaxTries").ID("int32"),
			),
			jen.ID("AzureConfig").Struct(
				jen.ID("AuthMethod").ID("string"),
				jen.ID("AccountName").ID("string"),
				jen.ID("BucketName").ID("string"),
				jen.ID("Retrying").Op("*").ID("AzureRetryConfig"),
				jen.ID("TokenCredentialsInitialToken").ID("string"),
				jen.ID("SharedKeyAccountKey").ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("cfg").Op("*").ID("AzureRetryConfig")).ID("buildRetryOptions").Params().Params(jen.ID("azblob").Dot("RetryOptions")).Body(
			jen.Return().ID("azblob").Dot("RetryOptions").Valuesln(
				jen.ID("Policy").Op(":").ID("azblob").Dot("RetryPolicyExponential"), jen.ID("MaxTries").Op(":").ID("cfg").Dot("MaxTries"), jen.ID("TryTimeout").Op(":").ID("cfg").Dot("TryTimeout"), jen.ID("RetryDelay").Op(":").ID("cfg").Dot("RetryDelay"), jen.ID("MaxRetryDelay").Op(":").ID("cfg").Dot("MaxRetryDelay"), jen.ID("RetryReadsFromSecondaryHost").Op(":").ID("cfg").Dot("RetryReadsFromSecondaryHost"))),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("azureSharedKeyAuthMethod1").Op("=").Lit("sharedkey"),
			jen.ID("azureSharedKeyAuthMethod2").Op("=").Lit("shared-key"),
			jen.ID("azureSharedKeyAuthMethod3").Op("=").Lit("shared_key"),
			jen.ID("azureSharedKeyAuthMethod4").Op("=").Lit("shared"),
			jen.ID("azureTokenAuthMethod").Op("=").Lit("token"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("AzureConfig")).ID("authMethodIsSharedKey").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("c").Dot("AuthMethod").Op("==").ID("azureSharedKeyAuthMethod1").Op("||").ID("c").Dot("AuthMethod").Op("==").ID("azureSharedKeyAuthMethod2").Op("||").ID("c").Dot("AuthMethod").Op("==").ID("azureSharedKeyAuthMethod3").Op("||").ID("c").Dot("AuthMethod").Op("==").ID("azureSharedKeyAuthMethod4")),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("AzureConfig")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the AzureConfig."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("AzureConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("AuthMethod"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("AccountName"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("BucketName"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("Retrying"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("Retrying").Op("!=").ID("nil"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("SharedKeyAccountKey"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("authMethodIsSharedKey").Call(),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Empty")),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("TokenCredentialsInitialToken"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("AuthMethod").Op("==").ID("azureTokenAuthMethod"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Empty")),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("provideAzureBucket").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("AzureConfig"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.Op("*").Qual("gocloud.dev/blob", "Bucket"), jen.ID("error")).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.Var().ID("cred").ID("azblob").Dot("Credential"),
			jen.ID("bucket").Op("*").Qual("gocloud.dev/blob", "Bucket"),
			jen.ID("err").ID("error"),
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("cfg").Dot("AuthMethod")))).Body(
				jen.Case(jen.ID("azureSharedKeyAuthMethod1"), jen.ID("azureSharedKeyAuthMethod2"), jen.ID("azureSharedKeyAuthMethod3"), jen.ID("azureSharedKeyAuthMethod4")).Body(
					jen.If(jen.ID("cfg").Dot("SharedKeyAccountKey").Op("==").Lit("")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidConfiguration"))), jen.If(jen.List(jen.ID("cred"), jen.ID("err")).Op("=").ID("azblob").Dot("NewSharedKeyCredential").Call(
						jen.ID("cfg").Dot("AccountName"),
						jen.ID("cfg").Dot("SharedKeyAccountKey"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("reading shared key credential: %w"),
							jen.ID("err"),
						)))),
				jen.Case(jen.ID("azureTokenAuthMethod")).Body(
					jen.If(jen.ID("cfg").Dot("TokenCredentialsInitialToken").Op("==").Lit("")).Body(
						jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidConfiguration"))), jen.ID("cred").Op("=").ID("azblob").Dot("NewTokenCredential").Call(
						jen.ID("cfg").Dot("TokenCredentialsInitialToken"),
						jen.ID("nil"),
					)),
				jen.Default().Body(
					jen.ID("cred").Op("=").ID("azblob").Dot("NewAnonymousCredential").Call()),
			),
			jen.If(jen.List(jen.ID("bucket"), jen.ID("err")).Op("=").ID("azureblob").Dot("OpenBucket").Call(
				jen.ID("ctx"),
				jen.ID("azureblob").Dot("NewPipeline").Call(
					jen.ID("cred"),
					jen.ID("buildPipelineOptions").Call(
						jen.ID("logger"),
						jen.ID("cfg").Dot("Retrying"),
					),
				),
				jen.ID("azureblob").Dot("AccountName").Call(jen.ID("cfg").Dot("AccountName")),
				jen.ID("cfg").Dot("BucketName"),
				jen.Op("&").ID("azureblob").Dot("Options").Valuesln(
					jen.ID("Protocol").Op(":").Lit("https")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("initializing azure bucket: %w"),
					jen.ID("err"),
				))),
			jen.Return().List(jen.ID("bucket"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildPipelineLogFunc").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("pipeline").Dot("LogLevel"), jen.ID("string")).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.Return().Func().Params(jen.ID("level").ID("pipeline").Dot("LogLevel"), jen.ID("message").ID("string")).Body(
				jen.Switch(jen.ID("level")).Body(
					jen.Case(jen.ID("pipeline").Dot("LogNone")).Body(),
					jen.Case(jen.ID("pipeline").Dot("LogPanic"), jen.ID("pipeline").Dot("LogFatal"), jen.ID("pipeline").Dot("LogError")).Body(
						jen.ID("logger").Dot("Error").Call(
							jen.ID("nil"),
							jen.ID("message"),
						)),
					jen.Case(jen.ID("pipeline").Dot("LogWarning")).Body(
						jen.ID("logger").Dot("Debug").Call(jen.ID("message"))),
					jen.Case(jen.ID("pipeline").Dot("LogInfo")).Body(
						jen.ID("logger").Dot("Info").Call(jen.ID("message"))),
					jen.Case(jen.ID("pipeline").Dot("LogDebug")).Body(
						jen.ID("logger").Dot("Debug").Call(jen.ID("message"))),
					jen.Default().Body(
						jen.ID("logger").Dot("Debug").Call(jen.ID("message"))),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildPipelineOptions").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("retrying").Op("*").ID("AzureRetryConfig")).Params(jen.ID("azblob").Dot("PipelineOptions")).Body(
			jen.ID("logger").Op("=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.ID("options").Op(":=").ID("azblob").Dot("PipelineOptions").Valuesln(
				jen.ID("Log").Op(":").ID("pipeline").Dot("LogOptions").Valuesln(
					jen.ID("Log").Op(":").ID("buildPipelineLogFunc").Call(jen.ID("logger")), jen.ID("ShouldLog").Op(":").Func().Params(jen.ID("level").ID("pipeline").Dot("LogLevel")).Params(jen.ID("bool")).Body(
						jen.Return().ID("level").Op("!=").ID("pipeline").Dot("LogNone")))),
			jen.If(jen.ID("retrying").Op("!=").ID("nil")).Body(
				jen.ID("options").Dot("Retry").Op("=").ID("retrying").Dot("buildRetryOptions").Call()),
			jen.Return().ID("options"),
		),
		jen.Line(),
	)

	return code
}
