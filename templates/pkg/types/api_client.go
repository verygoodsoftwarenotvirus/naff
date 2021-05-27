package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("APIClientKey").ID("ContextKey").Op("=").Lit("api_client"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("APIClient").Struct(
				jen.ID("LastUpdatedOn").Op("*").ID("uint64"),
				jen.ID("ArchivedOn").Op("*").ID("uint64"),
				jen.ID("Name").ID("string"),
				jen.ID("ClientID").ID("string"),
				jen.ID("ExternalID").ID("string"),
				jen.ID("ClientSecret").Index().ID("byte"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("ID").ID("uint64"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
			jen.ID("APIClientList").Struct(
				jen.ID("Clients").Index().Op("*").ID("APIClient"),
				jen.ID("Pagination"),
			),
			jen.ID("APIClientCreationInput").Struct(
				jen.ID("UserLoginInput"),
				jen.ID("Name").ID("string"),
				jen.ID("ClientID").ID("string"),
				jen.ID("ClientSecret").Index().ID("byte"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
			jen.ID("APIClientCreationResponse").Struct(
				jen.ID("ClientID").ID("string"),
				jen.ID("ClientSecret").ID("string"),
				jen.ID("ID").ID("uint64"),
			),
			jen.ID("APIClientDataManager").Interface(
				jen.ID("GetAPIClientByClientID").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("string")).Params(jen.Op("*").ID("APIClient"), jen.ID("error")),
				jen.ID("GetAPIClientByDatabaseID").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("ownerUserID")).ID("uint64")).Params(jen.Op("*").ID("APIClient"), jen.ID("error")),
				jen.ID("GetAllAPIClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("resultChannel").Chan().Index().Op("*").ID("APIClient"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")),
				jen.ID("GetTotalAPIClientCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetAPIClients").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("ownerUserID").ID("uint64"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("APIClientList"), jen.ID("error")),
				jen.ID("CreateAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("APIClientCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("APIClient"), jen.ID("error")),
				jen.ID("ArchiveAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("clientID"), jen.ID("ownerUserID"), jen.ID("archivedByUser")).ID("uint64")).Params(jen.ID("error")),
				jen.ID("GetAuditLogEntriesForAPIClient").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("clientID").ID("uint64")).Params(jen.Index().Op("*").ID("AuditLogEntry"), jen.ID("error")),
			),
			jen.ID("APIClientDataService").Interface(
				jen.ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates an APICreationInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("APIClientCreationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("minUsernameLength"), jen.ID("minPasswordLength")).ID("uint8")).Params(jen.ID("error")).Body(
			jen.If(jen.ID("err").Op(":=").ID("x").Dot("UserLoginInput").Dot("ValidateWithContext").Call(
				jen.ID("ctx"),
				jen.ID("minUsernameLength"),
				jen.ID("minPasswordLength"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err")),
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
