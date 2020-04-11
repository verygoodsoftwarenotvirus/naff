package queriers

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"log"
	"strings"
)

const (
	postgres = "postgres"
	sqlite   = "sqlite"
	mariadb  = "mariadb"

	existencePrefix = "SELECT EXISTS ("
	existenceSuffix = ")"

	whateverValue = "fart"

	// countQuery is a generic counter query used in a few query builders
	countQuery = "COUNT(%s.id)"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, vendor := range []string{postgres, sqlite, mariadb} {
		if err := renderDatabasePackage(proj, vendor); err != nil {
			return err
		}
	}
	return nil
}

// GetDatabasePalabra gets a given DB's superpalabra
func GetDatabasePalabra(vendor string) wordsmith.SuperPalabra {
	switch vendor {
	case postgres:
		return wordsmith.FromSingularPascalCase("Postgres")
	case sqlite:
		return wordsmith.FromSingularPascalCase("Sqlite")
	case mariadb:
		return &wordsmith.ManualWord{
			SingularStr:                           "MariaDB",
			PluralStr:                             "MariaDBs",
			RouteNameStr:                          "mariadb",
			KebabNameStr:                          "mariadb",
			PluralRouteNameStr:                    "mariadbs",
			UnexportedVarNameStr:                  "mariaDB",
			PluralUnexportedVarNameStr:            "mariaDBs",
			PackageNameStr:                        "mariadbs",
			SingularPackageNameStr:                "mariadb",
			SingularCommonNameStr:                 "maria DB",
			ProperSingularCommonNameWithPrefixStr: "a Maria DB",
			PluralCommonNameStr:                   "maria DBs",
			SingularCommonNameWithPrefixStr:       "maria DB",
			PluralCommonNameWithPrefixStr:         "maria DBs",
		}
	default:
		panic(fmt.Sprintf("unknown vendor: %q", vendor))
	}
}

// GetOAuth2ClientPalabra gets a given DB's superpalabra
func GetOAuth2ClientPalabra() wordsmith.SuperPalabra {
	return &wordsmith.ManualWord{
		SingularStr:                           "OAuth2Client",
		PluralStr:                             "OAuth2Clients",
		RouteNameStr:                          "oauth2_client",
		KebabNameStr:                          "oauth2-client",
		PluralRouteNameStr:                    "oauth2_clients",
		UnexportedVarNameStr:                  "oauth2Client",
		PluralUnexportedVarNameStr:            "oauth2Clients",
		PackageNameStr:                        "oauth2clients",
		SingularPackageNameStr:                "oauth2client",
		SingularCommonNameStr:                 "OAuth2 client",
		ProperSingularCommonNameWithPrefixStr: "an OAuth2 client",
		PluralCommonNameStr:                   "OAuth2 clients",
		SingularCommonNameWithPrefixStr:       "OAuth2 client",
		PluralCommonNameWithPrefixStr:         "OAuth2 clients",
	}
}

// GetUserPalabra gets a given DB's superpalabra
func GetUserPalabra() wordsmith.SuperPalabra {
	return &wordsmith.ManualWord{
		SingularStr:                           "User",
		PluralStr:                             "Users",
		RouteNameStr:                          "user",
		KebabNameStr:                          "user",
		PluralRouteNameStr:                    "users",
		UnexportedVarNameStr:                  "user",
		PluralUnexportedVarNameStr:            "users",
		PackageNameStr:                        "users",
		SingularPackageNameStr:                "user",
		SingularCommonNameStr:                 "user",
		ProperSingularCommonNameWithPrefixStr: "a User",
		PluralCommonNameStr:                   "users",
		SingularCommonNameWithPrefixStr:       "a user",
		PluralCommonNameWithPrefixStr:         "users",
	}
}

// GetWebhookPalabra gets a given DB's superpalabra
func GetWebhookPalabra() wordsmith.SuperPalabra {
	return &wordsmith.ManualWord{
		SingularStr:                           "Webhook",
		PluralStr:                             "Webhooks",
		RouteNameStr:                          "webhook",
		KebabNameStr:                          "webhook",
		PluralRouteNameStr:                    "webhooks",
		UnexportedVarNameStr:                  "webhook",
		PluralUnexportedVarNameStr:            "webhooks",
		PackageNameStr:                        "webhooks",
		SingularPackageNameStr:                "webhook",
		SingularCommonNameStr:                 "webhook",
		ProperSingularCommonNameWithPrefixStr: "a Webhook",
		PluralCommonNameStr:                   "webhooks",
		SingularCommonNameWithPrefixStr:       "a webhook",
		PluralCommonNameWithPrefixStr:         "webhooks",
	}
}

// renderDatabasePackage renders the package
func renderDatabasePackage(proj *models.Project, vendor string) error {
	var (
		dbDesc     string
		vendorWord wordsmith.SuperPalabra
	)

	if vendor == postgres {
		dbDesc = "Postgres instances"
		vendorWord = wordsmith.FromSingularPascalCase("Postgres")
	} else if vendor == sqlite {
		dbDesc = "sqlite files"
		vendorWord = wordsmith.FromSingularPascalCase("Sqlite")
	} else if vendor == mariadb {
		dbDesc = "MariaDB instances"
		vendorWord = buildMariaDBWord()
	}

	if vendorWord == nil {
		return errors.New("wtf")
	}

	pn := vendorWord.SingularPackageName()

	files := map[string]*jen.File{
		fmt.Sprintf("database/v1/queriers/%s/oauth2_clients.go", vendor):      oauth2ClientsDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/%s.go", pn, pn):                  databaseDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/webhooks.go", vendor):            webhooksDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/wire.go", vendor):                wireDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/doc.go", vendor):                 docDotGo(pn, dbDesc),
		fmt.Sprintf("database/v1/queriers/%s/%s_test.go", pn, pn):             databaseTestDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/users.go", vendor):               usersDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/users_test.go", vendor):          usersTestDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/webhooks_test.go", vendor):       webhooksTestDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/migrations.go", vendor):          migrationsDotGo(proj, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/oauth2_clients_test.go", vendor): oauth2ClientsTestDotGo(proj, vendorWord),
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("database/v1/queriers/%s/%s.go", vendor, typ.Name.PluralRouteName())] = iterablesDotGo(proj, vendorWord, typ)
		files[fmt.Sprintf("database/v1/queriers/%s/%s_test.go", vendor, typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, vendorWord, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, path, file); err != nil {
			return err
		}
	}

	return nil
}

func buildMariaDBWord() wordsmith.SuperPalabra {
	return &wordsmith.ManualWord{
		SingularStr:                           "MariaDB",
		PluralStr:                             "MariaDBs",
		RouteNameStr:                          "mariadb",
		KebabNameStr:                          "mariadb",
		AbbreviationStr:                       "M",
		LowercaseAbbreviationStr:              "m",
		PluralRouteNameStr:                    "mariadbs",
		UnexportedVarNameStr:                  "mariaDB",
		PluralUnexportedVarNameStr:            "mariaDBs",
		PackageNameStr:                        "mariadbs",
		SingularPackageNameStr:                "mariadb",
		SingularCommonNameStr:                 "maria DB",
		ProperSingularCommonNameWithPrefixStr: "a Maria DB",
		PluralCommonNameStr:                   "maria DBs",
		SingularCommonNameWithPrefixStr:       "maria DB",
		PluralCommonNameWithPrefixStr:         "maria DBs",
	}
}

type SqlBuilder interface {
	ToSql() (string, []interface{}, error)
}

func buildQueryTest(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType, queryName string, queryBuilder SqlBuilder, expectedArgs, callArgs []jen.Code, includeExpectedAndActualArgs, listQuery, includeFilter, createUser, createExampleVariable, excludeUserID bool, preQueryLines []jen.Code) []jen.Code {
	const (
		expectedQueryVarName = "expectedQuery"
		expectedArgsVarName  = "expectedArgs"
		actualQueryVarName   = "actualQuery"
		actualArgsVarName    = "actualArgs"
		uhOh                 = "QUERY"
	)

	if strings.HasSuffix(strings.ToUpper(queryName), uhOh) {
		queryName = queryName[:len(queryName)-len(uhOh)]
	}

	if createUser && !excludeUserID {
		callArgs = append(callArgs, jen.ID("exampleUser").Dot("ID"))
	}
	if includeFilter {
		callArgs = append(callArgs, jen.ID(utils.FilterVarName))
	}

	dbn := dbvendor.Singular()
	dbi := dbvendor.LowercaseAbbreviation()
	dbfl := dbi[0]

	expectedQuery, _, err := queryBuilder.ToSql()
	if err != nil {
		log.Panicf("error building %q: %v", queryName, err)
	}

	lines := []jen.Code{
		jen.Func().IDf("Test%s_build%sQuery", dbn, queryName).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.List(jen.ID(string(dbfl)), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
				func() jen.Code {
					if createUser {
						return jen.ID("exampleUser").Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call()
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if createExampleVariable && !listQuery && typ.Name != nil {
						sn := typ.Name.Singular()
						return jen.ID(utils.BuildFakeVarName(sn)).Assign().Qual(proj.FakeModelsPackage(), fmt.Sprintf("BuildFake%s", sn)).Call()
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if includeFilter {
						return jen.ID(utils.FilterVarName).Assign().Qual(proj.FakeModelsPackage(), "BuildFleshedOutQueryFilter").Call()
					}
					return jen.Null()
				}(),
				func() jen.Code {
					g := &jen.Group{}
					if len(preQueryLines) > 0 {
						for i := range preQueryLines {
							g.Add(preQueryLines[i], jen.Line())
						}
					}
					return g
				}(),
				jen.Line(),
				jen.ID(expectedQueryVarName).Assign().Lit(expectedQuery),
				func() jen.Code {
					if includeExpectedAndActualArgs {
						x := jen.ID(expectedArgsVarName).Assign().Index().Interface()
						if len(expectedArgs) > 0 {
							return x.Valuesln(expectedArgs...)
						}
						return x.Values(expectedArgs...)
					}
					return jen.Null()
				}(),
				func() jen.Code {
					returnedVars := []jen.Code{jen.ID(actualQueryVarName)}
					if includeExpectedAndActualArgs {
						returnedVars = append(returnedVars, jen.ID(actualArgsVarName))
					}
					return jen.List(returnedVars...).Assign().ID(dbi).Dotf("build%sQuery", queryName).Call(callArgs...)
				}(),
				jen.Line(),
				jen.ID("ensureArgCountMatchesQuery").Call(
					jen.ID("t"),
					jen.ID(actualQueryVarName),
					func() jen.Code {
						if includeExpectedAndActualArgs {
							return jen.ID(actualArgsVarName)
						}
						return jen.Index().Interface().Values()
					}(),
				),
				utils.AssertEqual(jen.ID(expectedQueryVarName), jen.ID(actualQueryVarName), nil),
				func() jen.Code {
					if includeExpectedAndActualArgs {
						return utils.AssertEqual(jen.ID(expectedArgsVarName), jen.ID(actualArgsVarName), nil)
					}
					return jen.Null()
				}(),
			),
		),
		jen.Line(),
	}

	return lines
}

func unixTimeForDatabase(db wordsmith.SuperPalabra) string {
	switch db.LowercaseAbbreviation() {
	case "m":
		return "UNIX_TIMESTAMP()"
	case "p":
		return "extract(epoch FROM NOW())"
	case "s":
		return "(strftime('%s','now'))"
	default:
		log.Fatalf("invalid database type! %q", db.LowercaseAbbreviation())
	}
	panic("won't get here")
}

func queryBuilderForDatabase(db wordsmith.SuperPalabra) squirrel.StatementBuilderType {
	switch db.LowercaseAbbreviation() {
	case "m":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	case "p":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	case "s":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	default:
		log.Fatalf("invalid database type! %q", db.LowercaseAbbreviation())
	}
	panic("won't get here")
}
