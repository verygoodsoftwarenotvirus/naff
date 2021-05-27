package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("numericIDPattern").Op("=").Lit("{%s:[0-9]+}"),
			jen.ID("unauthorizedRedirectResponseCode").Op("=").Qual("net/http", "StatusSeeOther"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SetupRoutes sets up the routes."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("SetupRoutes").Params(jen.ID("router").ID("routing").Dot("Router")).Body(
			jen.ID("router").Op("=").ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("UserAttributionMiddleware")),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/"),
				jen.ID("s").Dot("homepage"),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/dashboard"),
				jen.ID("s").Dot("homepage"),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/favicon.svg"),
				jen.ID("s").Dot("favicon"),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/login"),
				jen.ID("s").Dot("buildLoginView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/components/login_prompt"),
				jen.ID("s").Dot("buildLoginView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/auth/submit_login"),
				jen.ID("s").Dot("handleLoginSubmission"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/logout"),
				jen.ID("s").Dot("handleLogoutSubmission"),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/register"),
				jen.ID("s").Dot("registrationView"),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/components/registration_prompt"),
				jen.ID("s").Dot("registrationComponent"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/auth/submit_registration"),
				jen.ID("s").Dot("handleRegistrationSubmission"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/auth/verify_two_factor_secret"),
				jen.ID("s").Dot("handleTOTPVerificationSubmission"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/billing/checkout/begin"),
				jen.ID("s").Dot("handleCheckoutSessionStart"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/billing/checkout/success"),
				jen.ID("s").Dot("handleCheckoutSuccess"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/billing/checkout/cancel"),
				jen.ID("s").Dot("handleCheckoutCancel"),
			),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/billing/checkout/failures"),
				jen.ID("s").Dot("handleCheckoutFailure"),
			),
			jen.ID("singleAccountPattern").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("numericIDPattern"),
				jen.ID("accountIDURLParamKey"),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/accounts"),
				jen.ID("s").Dot("buildAccountsTableView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/accounts/%s"),
					jen.ID("singleAccountPattern"),
				),
				jen.ID("s").Dot("buildAccountEditorView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/dashboard_pages/accounts"),
				jen.ID("s").Dot("buildAccountsTableView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/accounts/%s"),
					jen.ID("singleAccountPattern"),
				),
				jen.ID("s").Dot("buildAccountEditorView").Call(jen.ID("false")),
			),
			jen.ID("singleAPIClientPattern").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("numericIDPattern"),
				jen.ID("apiClientIDURLParamKey"),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadAPIClientsPermission"))).Dot("Get").Call(
				jen.Lit("/api_clients"),
				jen.ID("s").Dot("buildAPIClientsTableView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadAPIClientsPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/api_clients"),
				jen.ID("s").Dot("buildAPIClientsTableView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadAPIClientsPermission"))).Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/api_clients/%s"),
					jen.ID("singleAPIClientPattern"),
				),
				jen.ID("s").Dot("buildAPIClientEditorView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadAPIClientsPermission"))).Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/api_clients/%s"),
					jen.ID("singleAPIClientPattern"),
				),
				jen.ID("s").Dot("buildAPIClientEditorView").Call(jen.ID("false")),
			),
			jen.ID("singleWebhookPattern").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("numericIDPattern"),
				jen.ID("webhookIDURLParamKey"),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadWebhooksPermission"))).Dot("Get").Call(
				jen.Lit("/account/webhooks"),
				jen.ID("s").Dot("buildWebhooksTableView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadWebhooksPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/account/webhooks"),
				jen.ID("s").Dot("buildWebhooksTableView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateWebhooksPermission"))).Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/account/webhooks/%s"),
					jen.ID("singleWebhookPattern"),
				),
				jen.ID("s").Dot("buildWebhookEditorView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateWebhooksPermission"))).Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/account/webhooks/%s"),
					jen.ID("singleWebhookPattern"),
				),
				jen.ID("s").Dot("buildWebhookEditorView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/user/settings"),
				jen.ID("s").Dot("buildUserSettingsView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("Get").Call(
				jen.Lit("/dashboard_pages/user/settings"),
				jen.ID("s").Dot("buildUserSettingsView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateAccountPermission"))).Dot("Get").Call(
				jen.Lit("/account/settings"),
				jen.ID("s").Dot("buildAccountSettingsView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateAccountPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/account/settings"),
				jen.ID("s").Dot("buildAccountSettingsView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("SearchUserPermission"))).Dot("Get").Call(
				jen.Lit("/admin/users/search"),
				jen.ID("s").Dot("buildUsersTableView").Call(
					jen.ID("true"),
					jen.ID("true"),
				),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("SearchUserPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/admin/users/search"),
				jen.ID("s").Dot("buildUsersTableView").Call(
					jen.ID("false"),
					jen.ID("true"),
				),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadUserPermission"))).Dot("Get").Call(
				jen.Lit("/admin/users"),
				jen.ID("s").Dot("buildUsersTableView").Call(
					jen.ID("true"),
					jen.ID("false"),
				),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadUserPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/admin/users"),
				jen.ID("s").Dot("buildUsersTableView").Call(
					jen.ID("false"),
					jen.ID("false"),
				),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("ServiceAdminMiddleware")).Dot("Get").Call(
				jen.Lit("/admin/settings"),
				jen.ID("s").Dot("buildAdminSettingsView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("ServiceAdminMiddleware")).Dot("Get").Call(
				jen.Lit("/dashboard_pages/admin/settings"),
				jen.ID("s").Dot("buildAdminSettingsView").Call(jen.ID("false")),
			),
			jen.ID("singleItemPattern").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("numericIDPattern"),
				jen.ID("itemIDURLParamKey"),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadItemsPermission"))).Dot("Get").Call(
				jen.Lit("/items"),
				jen.ID("s").Dot("buildItemsTableView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ReadItemsPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/items"),
				jen.ID("s").Dot("buildItemsTableView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("CreateItemsPermission"))).Dot("Get").Call(
				jen.Lit("/items/new"),
				jen.ID("s").Dot("buildItemCreatorView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("CreateItemsPermission"))).Dot("Post").Call(
				jen.Lit("/items/new/submit"),
				jen.ID("s").Dot("handleItemCreationRequest"),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ArchiveItemsPermission"))).Dot("Delete").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/items/%s"),
					jen.ID("singleItemPattern"),
				),
				jen.ID("s").Dot("handleItemDeletionRequest"),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("ArchiveItemsPermission"))).Dot("Get").Call(
				jen.Lit("/dashboard_pages/items/new"),
				jen.ID("s").Dot("buildItemCreatorView").Call(jen.ID("false")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateItemsPermission"))).Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/items/%s"),
					jen.ID("singleItemPattern"),
				),
				jen.ID("s").Dot("buildItemEditorView").Call(jen.ID("true")),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateItemsPermission"))).Dot("Put").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/items/%s"),
					jen.ID("singleItemPattern"),
				),
				jen.ID("s").Dot("handleItemUpdateRequest"),
			),
			jen.ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.ID("authorization").Dot("UpdateItemsPermission"))).Dot("Get").Call(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("/dashboard_pages/items/%s"),
					jen.ID("singleItemPattern"),
				),
				jen.ID("s").Dot("buildItemEditorView").Call(jen.ID("false")),
			),
		),
		jen.Line(),
	)

	return code
}
