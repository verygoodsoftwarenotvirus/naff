package storage

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bucketAzureTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestAzureConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AuthMethod").Op(":").ID("azureTokenAuthMethod"), jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("TokenCredentialsInitialToken").Op(":").ID("t").Dot("Name").Call()),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAzureConfig_authMethodIsSharedKey").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.Parens(jen.Op("&").ID("AzureConfig").Valuesln(
							jen.ID("AuthMethod").Op(":").ID("azureSharedKeyAuthMethod4"))).Dot("authMethodIsSharedKey").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAzureRetryConfig_buildRetryOptions").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.Parens(jen.Op("&").ID("AzureRetryConfig").Values()).Dot("buildRetryOptions").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_provideAzureBucket").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with anonymous credential"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with shared key"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AuthMethod").Op(":").ID("azureSharedKeyAuthMethod4"), jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("SharedKeyAccountKey").Op(":").Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()))),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with shared key and not shared key"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AuthMethod").Op(":").ID("azureSharedKeyAuthMethod4"), jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("SharedKeyAccountKey").Op(":").Lit("")),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with shared key and invalid key"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AuthMethod").Op(":").ID("azureSharedKeyAuthMethod4"), jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("SharedKeyAccountKey").Op(":").Lit("        lol not valid base64       ")),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AuthMethod").Op(":").ID("azureTokenAuthMethod"), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("TokenCredentialsInitialToken").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with token auth and no token"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("AzureConfig").Valuesln(
						jen.ID("AuthMethod").Op(":").ID("azureTokenAuthMethod"), jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("AccountName").Op(":").ID("t").Dot("Name").Call(), jen.ID("TokenCredentialsInitialToken").Op(":").Lit("")),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg"),
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildPipelineLogFunc").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").ID("buildPipelineLogFunc").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
					jen.For(jen.List(jen.ID("_"), jen.ID("level")).Op(":=").Range().Index().ID("pipeline").Dot("LogLevel").Valuesln(
						jen.ID("pipeline").Dot("LogNone"), jen.ID("pipeline").Dot("LogError"), jen.ID("pipeline").Dot("LogWarning"), jen.ID("pipeline").Dot("LogInfo"), jen.ID("pipeline").Dot("LogDebug"), jen.ID("pipeline").Dot("LogLevel").Call(jen.Qual("math", "MaxUint32")))).Body(
						jen.ID("x").Call(
							jen.ID("level"),
							jen.ID("t").Dot("Name").Call(),
						)),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildPipelineOptions").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("buildPipelineOptions").Call(
							jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
							jen.Op("&").ID("AzureRetryConfig").Values(),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
