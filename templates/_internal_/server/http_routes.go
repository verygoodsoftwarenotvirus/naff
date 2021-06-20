package server

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildIterableRoutes(proj *models.Project) []jen.Code {
	out := []jen.Code{}

	for _, typ := range proj.DataTypes {
		uvn := typ.Name.UnexportedVarName()
		puvn := typ.Name.PluralUnexportedVarName()
		sn := typ.Name.Singular()
		pn := typ.Name.Plural()
		pcn := typ.Name.PluralCommonName()

		out = append(out,
			jen.Newline(),
			jen.Comment(pn),
			jen.IDf("%sPath", uvn).Op(":=").Lit(pcn),
			jen.IDf("%sRouteWithPrefix", puvn).Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("/%s"),
				jen.IDf("%sPath", uvn),
			),
			jen.IDf("%sIDRouteParam", uvn).Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.ServicePackage(typ.Name.PackageName()), fmt.Sprintf("%sIDURIParamKey", sn))),
			jen.ID("v1Router").Dot("Route").Call(
				jen.IDf("%sRouteWithPrefix", puvn),
				jen.Func().Params(jen.IDf("%sRouter", puvn).ID("routing").Dot("Router")).Body(
					jen.IDf("%sRouter", puvn).
						Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
						Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Create%sPermission", pn)))).
						Dotln("Post").Call(jen.ID("root"), jen.ID("s").Dotf("%sService", puvn).Dot("CreateHandler")),
					jen.IDf("%sRouter", puvn).
						Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
						Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Read%sPermission", pn)))).
						Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dotf("%sService", puvn).Dot("ListHandler")),
					func() jen.Code {
						if typ.SearchEnabled {
							return jen.IDf("%sRouter", puvn).
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
								Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Read%sPermission", pn)))).
								Dotln("Get").Call(jen.ID("searchRoot"), jen.ID("s").Dotf("%sService", puvn).Dot("SearchHandler"))
						}
						return jen.Null()
					}(),
					jen.Newline(),
					jen.IDf("%sRouter", puvn).Dot("Route").Call(
						jen.IDf("%sIDRouteParam", uvn),
						jen.Func().Params(jen.IDf("single%sRouter", sn).ID("routing").Dot("Router")).Body(
							jen.IDf("single%sRouter", sn).
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
								Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Read%sPermission", pn)))).
								Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dotf("%sService", puvn).Dot("ReadHandler")),
							jen.IDf("single%sRouter", sn).
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
								Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Read%sPermission", pn)))).
								Dotln("Head").Call(jen.ID("root"), jen.ID("s").Dotf("%sService", puvn).Dot("ExistenceHandler")),
							jen.IDf("single%sRouter", sn).
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
								Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Archive%sPermission", pn)))).
								Dotln("Delete").Call(jen.ID("root"), jen.ID("s").Dotf("%sService", puvn).Dot("ArchiveHandler")),
							jen.IDf("single%sRouter", sn).
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
								Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Update%sPermission", pn)))).
								Dotln("Put").Call(jen.ID("root"), jen.ID("s").Dotf("%sService", puvn).Dot("UpdateHandler")),
							jen.IDf("single%sRouter", sn).
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").
								Call(jen.Qual(proj.InternalAuthorizationPackage(), fmt.Sprintf("Read%sAuditLogEntriesPermission", pn)))).
								Dotln("Get").Call(jen.ID("auditRoute"), jen.ID("s").Dotf("%sService", puvn).Dot("AuditEntryHandler")),
						),
					),
				),
			),
		)
	}

	return out
}

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("root").Op("=").Lit("/"),
			jen.ID("auditRoute").Op("=").Lit("/audit"),
			jen.ID("searchRoot").Op("=").Lit("/search"),
			jen.ID("numericIDPattern").Op("=").Lit("{%s:[0-9]+}"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("buildNumericIDURLChunk").Params(jen.ID("key").ID("string")).Params(jen.ID("string")).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.ID("root").Op("+").ID("numericIDPattern"),
				jen.ID("key"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("HTTPServer")).ID("setupRouter").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("router").ID("routing").Dot("Router"), jen.ID("metricsHandler").Qual(proj.MetricsPackage(), "Handler")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("router").Dot("Route").Call(
				jen.Lit("/_meta_"),
				jen.Func().Params(jen.ID("metaRouter").ID("routing").Dot("Router")).Body(
					jen.ID("health").Op(":=").Qual("github.com/heptiolabs/healthcheck", "NewHandler").Call(),
					jen.Comment("Expose a liveness check on /live"),
					jen.ID("metaRouter").Dot("Get").Call(
						jen.Lit("/live"),
						jen.ID("health").Dot("LiveEndpoint"),
					),
					jen.Comment("Expose a readiness check on /ready"),
					jen.ID("metaRouter").Dot("Get").Call(
						jen.Lit("/ready"),
						jen.ID("health").Dot("ReadyEndpoint"),
					),
				),
			),
			jen.Newline(),
			jen.If(jen.ID("metricsHandler").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("establishing metrics handler")),
				jen.ID("router").Dot("HandleFunc").Call(
					jen.Lit("/metrics"),
					jen.ID("metricsHandler").Dot("ServeHTTP"),
				),
			),
			jen.Newline(),
			jen.Comment("Frontend routes."),
			jen.ID("s").Dot("frontendService").Dot("SetupRoutes").Call(jen.ID("router")),
			jen.Newline(),
			jen.ID("router").Dot("Post").Call(
				jen.Lit("/paseto"),
				jen.ID("s").Dot("authService").Dot("PASETOHandler"),
			),
			jen.Newline(),
			jen.ID("authenticatedRouter").Op(":=").ID("router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("UserAttributionMiddleware")),
			jen.ID("authenticatedRouter").Dot("Get").Call(
				jen.Lit("/auth/status"),
				jen.ID("s").Dot("authService").Dot("StatusHandler"),
			),
			jen.Newline(),
			jen.ID("router").Dot("Route").Call(
				jen.Lit("/users"),
				jen.Func().Params(jen.ID("userRouter").ID("routing").Dot("Router")).Body(
					jen.ID("userRouter").Dot("Post").Call(
						jen.Lit("/login"),
						jen.ID("s").Dot("authService").Dot("BeginSessionHandler"),
					),
					jen.ID("userRouter").Dot("WithMiddleware").Call(
						jen.ID("s").Dot("authService").Dot("UserAttributionMiddleware"),
						jen.ID("s").Dot("authService").Dot("CookieRequirementMiddleware"),
					).Dot("Post").Call(
						jen.Lit("/logout"),
						jen.ID("s").Dot("authService").Dot("EndSessionHandler"),
					),
					jen.ID("userRouter").Dot("Post").Call(
						jen.ID("root"),
						jen.ID("s").Dot("usersService").Dot("CreateHandler"),
					),
					jen.ID("userRouter").Dot("Post").Call(
						jen.Lit("/totp_secret/verify"),
						jen.ID("s").Dot("usersService").Dot("TOTPSecretVerificationHandler"),
					),
					jen.Newline(),
					jen.Comment("need credentials beyond this point"),
					jen.ID("authedRouter").Op(":=").ID("userRouter").Dot("WithMiddleware").Call(
						jen.ID("s").Dot("authService").Dot("UserAttributionMiddleware"),
						jen.ID("s").Dot("authService").Dot("AuthorizationMiddleware"),
					),
					jen.ID("authedRouter").Dot("Post").Call(
						jen.Lit("/account/select"),
						jen.ID("s").Dot("authService").Dot("ChangeActiveAccountHandler"),
					),
					jen.ID("authedRouter").Dot("Post").Call(
						jen.Lit("/totp_secret/new"),
						jen.ID("s").Dot("usersService").Dot("NewTOTPSecretHandler"),
					),
					jen.ID("authedRouter").Dot("Put").Call(
						jen.Lit("/password/new"),
						jen.ID("s").Dot("usersService").Dot("UpdatePasswordHandler"),
					),
				),
			),
			jen.Newline(),
			jen.ID("authenticatedRouter").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("AuthorizationMiddleware")).Dot("Route").Call(jen.Lit("/api/v1"), jen.Func().Params(jen.ID("v1Router").ID("routing").Dot("Router")).Body(
				append([]jen.Code{
					jen.ID("adminRouter").Op(":=").ID("v1Router").Dot("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("ServiceAdminMiddleware")),
					jen.Newline(),
					jen.Comment("Admin"),
					jen.ID("adminRouter").Dot("Route").Call(jen.Lit("/admin"), jen.Func().Params(jen.ID("adminRouter").ID("routing").Dot("Router")).Body(
						jen.ID("adminRouter").
							Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "CycleCookieSecretPermission"))).
							Dotln("Post").Call(jen.Lit("/cycle_cookie_secret"), jen.ID("s").Dot("authService").Dot("CycleCookieSecretHandler")),
						jen.ID("adminRouter").
							Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "UpdateUserStatusPermission"))).
							Dotln("Post").Call(jen.Lit("/users/status"),
							jen.ID("s").Dot("adminService").Dot("UserReputationChangeHandler"),
						),
						jen.Newline(),
						jen.ID("adminRouter").Dot("Route").Call(jen.Lit("/audit_log"),
							jen.Func().Params(jen.ID("auditRouter").ID("routing").Dot("Router")).Body(
								jen.ID("entryIDRouteParam").Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.AuditServicePackage(), "LogEntryURIParamKey")),
								jen.ID("auditRouter").
									Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadAllAuditLogEntriesPermission"))).
									Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("auditService").Dot("ListHandler")),
								jen.ID("auditRouter").Dot("Route").Call(
									jen.ID("entryIDRouteParam"),
									jen.Func().Params(jen.ID("singleEntryRouter").ID("routing").Dot("Router")).Body(
										jen.ID("singleEntryRouter").
											Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadAllAuditLogEntriesPermission"))).
											Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("auditService").Dot("ReadHandler")),
									),
								),
							),
						),
					),
					),
					jen.Newline(),
					jen.Comment("Users"),
					jen.ID("v1Router").Dot("Route").Call(
						jen.Lit("/users"),
						jen.Func().Params(jen.ID("usersRouter").ID("routing").Dot("Router")).Body(
							jen.ID("usersRouter").
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadUserPermission"))).
								Dotln("Get").Call(
								jen.ID("root"),
								jen.ID("s").Dot("usersService").Dot("ListHandler"),
							),
							jen.ID("usersRouter").
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "SearchUserPermission"))).
								Dotln("Get").Call(
								jen.Lit("/search"),
								jen.ID("s").Dot("usersService").Dot("UsernameSearchHandler"),
							),
							jen.ID("usersRouter").Dot("Post").Call(
								jen.Lit("/avatar/upload"),
								jen.ID("s").Dot("usersService").Dot("AvatarUploadHandler"),
							),
							jen.ID("usersRouter").Dot("Get").Call(
								jen.Lit("/self"),
								jen.ID("s").Dot("usersService").Dot("SelfHandler"),
							),
							jen.Newline(),
							jen.ID("singleUserRoute").Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.UsersServicePackage(), "UserIDURIParamKey")),
							jen.ID("usersRouter").Dot("Route").Call(
								jen.ID("singleUserRoute"),
								jen.Func().Params(jen.ID("singleUserRouter").ID("routing").Dot("Router")).Body(
									jen.ID("singleUserRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadUserPermission"))).
										Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("usersService").Dot("ReadHandler")),
									jen.ID("singleUserRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadUserAuditLogEntriesPermission"))).
										Dotln("Get").Call(jen.ID("auditRoute"), jen.ID("s").Dot("usersService").Dot("AuditEntryHandler")),
									jen.Newline(),
									jen.ID("singleUserRouter").Dot("Delete").Call(
										jen.ID("root"),
										jen.ID("s").Dot("usersService").Dot("ArchiveHandler"),
									),
								),
							),
						),
					),
					jen.Newline(),
					jen.Comment("Accounts"),
					jen.ID("v1Router").Dot("Route").Call(
						jen.Lit("/accounts"),
						jen.Func().Params(jen.ID("accountsRouter").ID("routing").Dot("Router")).Body(
							jen.ID("accountsRouter").Dot("Post").Call(
								jen.ID("root"),
								jen.ID("s").Dot("accountsService").Dot("CreateHandler"),
							),
							jen.ID("accountsRouter").Dot("Get").Call(
								jen.ID("root"),
								jen.ID("s").Dot("accountsService").Dot("ListHandler"),
							),
							jen.Newline(),
							jen.ID("singleUserRoute").Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.AccountsServicePackage(), "UserIDURIParamKey")),
							jen.ID("singleAccountRoute").Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.AccountsServicePackage(), "AccountIDURIParamKey")),
							jen.ID("accountsRouter").Dot("Route").Call(
								jen.ID("singleAccountRoute"),
								jen.Func().Params(jen.ID("singleAccountRouter").ID("routing").Dot("Router")).Body(
									jen.ID("singleAccountRouter").Dot("Get").Call(
										jen.ID("root"),
										jen.ID("s").Dot("accountsService").Dot("ReadHandler"),
									),
									jen.ID("singleAccountRouter").Dot("Put").Call(
										jen.ID("root"),
										jen.ID("s").Dot("accountsService").Dot("UpdateHandler"),
									),
									jen.ID("singleAccountRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ArchiveAccountPermission"))).
										Dotln("Delete").Call(jen.ID("root"), jen.ID("s").Dot("accountsService").Dot("ArchiveHandler")),
									jen.ID("singleAccountRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadAccountAuditLogEntriesPermission"))).
										Dotln("Get").Call(jen.ID("auditRoute"), jen.ID("s").Dot("accountsService").Dot("AuditEntryHandler")),
									jen.Newline(),
									jen.ID("singleAccountRouter").Dot("Post").Call(
										jen.Lit("/default"),
										jen.ID("s").Dot("accountsService").Dot("MarkAsDefaultAccountHandler"),
									),
									jen.ID("singleAccountRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "RemoveMemberAccountPermission"))).
										Dotln("Delete").Call(jen.Lit("/members").Op("+").ID("singleUserRoute"), jen.ID("s").Dot("accountsService").Dot("RemoveMemberHandler")),
									jen.ID("singleAccountRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "AddMemberAccountPermission"))).
										Dotln("Post").Call(jen.Lit("/member"), jen.ID("s").Dot("accountsService").Dot("AddMemberHandler")),
									jen.ID("singleAccountRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ModifyMemberPermissionsForAccountPermission"))).
										Dotln("Patch").Call(jen.Lit("/members").Op("+").ID("singleUserRoute").Op("+").Lit("/permissions"), jen.ID("s").Dot("accountsService").Dot("ModifyMemberPermissionsHandler")),
									jen.ID("singleAccountRouter").Dot("Post").Call(jen.Lit("/transfer"), jen.ID("s").Dot("accountsService").Dot("TransferAccountOwnershipHandler")),
								),
							),
						),
					),
					jen.Newline(),
					jen.Comment("API Clients"),
					jen.ID("v1Router").Dot("Route").Call(
						jen.Lit("/api_clients"),
						jen.Func().Params(jen.ID("clientRouter").ID("routing").Dot("Router")).Body(
							jen.ID("clientRouter").
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadAPIClientsPermission"))).
								Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("apiClientsService").Dot("ListHandler")),
							jen.ID("clientRouter").
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "CreateAPIClientsPermission"))).
								Dotln("Post").Call(jen.ID("root"), jen.ID("s").Dot("apiClientsService").Dot("CreateHandler")),
							jen.Newline(),
							jen.ID("singleClientRoute").Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.APIClientsServicePackage(), "APIClientIDURIParamKey")),
							jen.ID("clientRouter").Dot("Route").Call(
								jen.ID("singleClientRoute"),
								jen.Func().Params(jen.ID("singleClientRouter").ID("routing").Dot("Router")).Body(
									jen.ID("singleClientRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadAPIClientsPermission"))).
										Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("apiClientsService").Dot("ReadHandler")),
									jen.ID("singleClientRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ArchiveAPIClientsPermission"))).
										Dotln("Delete").Call(jen.ID("root"), jen.ID("s").Dot("apiClientsService").Dot("ArchiveHandler")),
									jen.ID("singleClientRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadAPIClientAuditLogEntriesPermission"))).
										Dotln("Get").Call(jen.ID("auditRoute"), jen.ID("s").Dot("apiClientsService").Dot("AuditEntryHandler")),
								),
							),
						),
					),
					jen.Newline(),
					jen.Comment("Webhooks"),
					jen.ID("v1Router").Dot("Route").Call(
						jen.Lit("/webhooks"),
						jen.Func().Params(jen.ID("webhookRouter").ID("routing").Dot("Router")).Body(
							jen.ID("singleWebhookRoute").Op(":=").ID("buildNumericIDURLChunk").Call(jen.Qual(proj.WebhooksServicePackage(), "WebhookIDURIParamKey")),
							jen.ID("webhookRouter").
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadWebhooksPermission"))).
								Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("webhooksService").Dot("ListHandler")),
							jen.ID("webhookRouter").
								Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "CreateWebhooksPermission"))).
								Dotln("Post").Call(jen.ID("root"), jen.ID("s").Dot("webhooksService").Dot("CreateHandler")),
							jen.ID("webhookRouter").Dot("Route").Call(
								jen.ID("singleWebhookRoute"),
								jen.Func().Params(jen.ID("singleWebhookRouter").ID("routing").Dot("Router")).Body(
									jen.ID("singleWebhookRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadWebhooksPermission"))).
										Dotln("Get").Call(jen.ID("root"), jen.ID("s").Dot("webhooksService").Dot("ReadHandler")),
									jen.ID("singleWebhookRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ArchiveWebhooksPermission"))).
										Dotln("Delete").Call(jen.ID("root"), jen.ID("s").Dot("webhooksService").Dot("ArchiveHandler")),
									jen.ID("singleWebhookRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "UpdateWebhooksPermission"))).
										Dotln("Put").Call(jen.ID("root"), jen.ID("s").Dot("webhooksService").Dot("UpdateHandler")),
									jen.ID("singleWebhookRouter").
										Dotln("WithMiddleware").Call(jen.ID("s").Dot("authService").Dot("PermissionFilterMiddleware").Call(jen.Qual(proj.InternalAuthorizationPackage(), "ReadWebhooksAuditLogEntriesPermission"))).
										Dotln("Get").Call(jen.ID("auditRoute"), jen.ID("s").Dot("webhooksService").Dot("AuditEntryHandler")),
								),
							),
						),
					),
				},
					buildIterableRoutes(proj)...)...,
			)),
			jen.Newline(),
			jen.ID("s").Dot("router").Op("=").ID("router"),
		),
		jen.Newline(),
	)

	return code
}
