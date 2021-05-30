package secrets

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("ProviderLocal").Op("=").Lit("local"),
			jen.ID("ProviderGCP").Op("=").Lit("gcp_kms"),
			jen.ID("ProviderAWS").Op("=").Lit("aws_kms"),
			jen.ID("ProviderAzureKeyVault").Op("=").Lit("azure_keyvault"),
			jen.ID("ProviderHashicorpVault").Op("=").Lit("vault"),
			jen.ID("expectedLocalKeyLength").Op("=").Lit(32),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidProvider").Op("=").Qual("errors", "New").Call(jen.Lit("invalid provider")),
			jen.ID("errNilConfig").Op("=").Qual("errors", "New").Call(jen.Lit("nil config provided")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Config").Struct(
				jen.ID("Provider").ID("string"),
				jen.ID("Key").ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideSecretKeeper provides a new secret keeper."),
		jen.Line(),
		jen.Func().ID("ProvideSecretKeeper").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.Op("*").Qual("gocloud.dev/secrets", "Keeper"), jen.ID("error")).Body(
			jen.If(jen.ID("cfg").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errNilConfig"))),
			jen.Switch(jen.ID("cfg").Dot("Provider")).Body(
				jen.Case(jen.ID("ProviderGCP")).Body(
					jen.List(jen.ID("client"), jen.ID("_"), jen.ID("err")).Op(":=").ID("gcpkms").Dot("Dial").Call(
						jen.ID("ctx"),
						jen.ID("nil"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("connecting to GCP KMS: %w"),
							jen.ID("err"),
						))), jen.ID("keeper").Op(":=").ID("gcpkms").Dot("OpenKeeper").Call(
						jen.ID("client"),
						jen.ID("cfg").Dot("Key"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("keeper"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderAWS")).Body(
					jen.List(jen.ID("sess"), jen.ID("err")).Op(":=").ID("session").Dot("NewSession").Call(jen.ID("nil")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("doing: %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("awskms").Dot("Dial").Call(jen.ID("sess")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("doing: %w"),
							jen.ID("err"),
						))), jen.ID("keeper").Op(":=").ID("awskms").Dot("OpenKeeper").Call(
						jen.ID("client"),
						jen.ID("cfg").Dot("Key"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("keeper"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderAzureKeyVault")).Body(
					jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("azurekeyvault").Dot("Dial").Call(), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("doing: %w"),
							jen.ID("err"),
						))), jen.List(jen.ID("keeper"), jen.ID("err")).Op(":=").ID("azurekeyvault").Dot("OpenKeeper").Call(
						jen.ID("client"),
						jen.ID("cfg").Dot("Key"),
						jen.ID("nil"),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("doing: %w"),
							jen.ID("err"),
						))), jen.Return().List(jen.ID("keeper"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderHashicorpVault")).Body(
					jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("hashivault").Dot("Dial").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("hashivault").Dot("Config").Valuesln(jen.ID("Token").Op(":").Lit(""), jen.ID("APIConfig").Op(":").ID("api").Dot("Config").Valuesln(jen.ID("Address").Op(":").Lit(""))),
					), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("doing: %w"),
							jen.ID("err"),
						))), jen.ID("keeper").Op(":=").ID("hashivault").Dot("OpenKeeper").Call(
						jen.ID("client"),
						jen.ID("cfg").Dot("Key"),
						jen.ID("nil"),
					), jen.Return().List(jen.ID("keeper"), jen.ID("nil"))),
				jen.Case(jen.ID("ProviderLocal")).Body(
					jen.List(jen.ID("key"), jen.ID("err")).Op(":=").ID("localsecrets").Dot("Base64Key").Call(jen.ID("cfg").Dot("Key")), jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("doing: %w"),
							jen.ID("err"),
						))), jen.ID("keeper").Op(":=").ID("localsecrets").Dot("NewKeeper").Call(jen.ID("key")), jen.Return().List(jen.ID("keeper"), jen.ID("nil"))),
				jen.Default().Body(
					jen.Return().List(jen.ID("nil"), jen.ID("errInvalidProvider"))),
			),
		),
		jen.Line(),
	)

	return code
}
