package querybuilding

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"log"
	"path/filepath"
	"strings"
)

const (
	postgres = string(models.Postgres)
	sqlite   = string(models.Sqlite)
	mariadb  = string(models.MariaDB)

	basePackagePath = "internal/database/querybuilding"

	existencePrefix = "SELECT EXISTS ("
	existenceSuffix = ")"

	whateverValue = "fart"

	// countQuery is a generic counter query used in a few query builders
	countQuery = "COUNT(%s.id)"
)

func isPostgres(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.Singular() == "Postgres"
}

func isSqlite(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.Singular() == "Sqlite"
}

func isMariaDB(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.Singular() == "MariaDB" || dbvendor.RouteName() == "maria_db"
}

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, vendor := range proj.EnabledDatabases() {
		if err := renderDatabasePackage(proj, vendor); err != nil {
			return err
		}
	}
	return nil
}

// GetDatabasePalabra gets a given DB's superpalabra

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

	switch vendor {
	case postgres:
		dbDesc = "Postgres instances"
		vendorWord = wordsmith.FromSingularPascalCase("Postgres")
	case sqlite:
		dbDesc = "sqlite files"
		vendorWord = wordsmith.FromSingularPascalCase("Sqlite")
	case mariadb:
		dbDesc = "MariaDB instances"
		vendorWord = buildMariaDBWord()
	default:
		return errors.New("wtf")
	}

	pn := vendorWord.SingularPackageName()

	files := map[string]*jen.File{
		fmt.Sprintf("%s/oauth2_clients.go", vendor):      oauth2ClientsDotGo(proj, vendorWord),
		fmt.Sprintf("%s/%s.go", pn, pn):                  databaseDotGo(proj, vendorWord),
		fmt.Sprintf("%s/webhooks.go", vendor):            webhooksDotGo(proj, vendorWord),
		fmt.Sprintf("%s/wire.go", vendor):                wireDotGo(proj, vendorWord),
		fmt.Sprintf("%s/doc.go", vendor):                 docDotGo(pn, dbDesc),
		fmt.Sprintf("%s/%s_test.go", pn, pn):             databaseTestDotGo(proj, vendorWord),
		fmt.Sprintf("%s/users.go", vendor):               usersDotGo(proj, vendorWord),
		fmt.Sprintf("%s/users_test.go", vendor):          usersTestDotGo(proj, vendorWord),
		fmt.Sprintf("%s/webhooks_test.go", vendor):       webhooksTestDotGo(proj, vendorWord),
		fmt.Sprintf("%s/migrations.go", vendor):          migrationsDotGo(proj, vendorWord),
		fmt.Sprintf("%s/oauth2_clients_test.go", vendor): oauth2ClientsTestDotGo(proj, vendorWord),
	}

	if isMariaDB(vendorWord) || isSqlite(vendorWord) {
		files[fmt.Sprintf("%s/time_teller.go", vendor)] = timeTellerDotGo(proj, vendorWord)
		files[fmt.Sprintf("%s/time_teller_test.go", vendor)] = timeTellerTestDotGo(proj, vendorWord)
	}

	for _, typ := range proj.DataTypes {
		files[fmt.Sprintf("%s/%s.go", vendor, typ.Name.PluralRouteName())] = iterablesDotGo(proj, vendorWord, typ)
		files[fmt.Sprintf("%s/%s_test.go", vendor, typ.Name.PluralRouteName())] = iterablesTestDotGo(proj, vendorWord, typ)
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
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

type sqlBuilder interface {
	ToSql() (string, []interface{}, error)
}

func convertArgsToCode(args []interface{}) (code []jen.Code) {
	for _, arg := range args {
		if c, ok := arg.(models.Coder); ok {
			code = append(code, c.Code())
		}
	}
	return
}

func buildQueryTest(
	dbvendor wordsmith.SuperPalabra,
	queryName string,
	queryBuilder sqlBuilder,
	expectedArgs,
	callArgs,
	preQueryLines []jen.Code,
) []jen.Code {
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

	dbn := dbvendor.Singular()
	dbi := dbvendor.LowercaseAbbreviation()
	dbfl := dbi[0]

	expectedQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		log.Panicf("error building %q: %v", queryName, err)
	}

	coderArgs := convertArgsToCode(args)
	if len(coderArgs) == len(args) {
		expectedArgs = coderArgs
	}

	block := append(
		[]jen.Code{
			jen.List(jen.ID(string(dbfl)), jen.Underscore()).Assign().ID("buildTestService").Call(jen.ID("t")),
			jen.Line(),
		},
		preQueryLines...,
	)
	block = append(block,
		jen.Line(),
		jen.ID(expectedQueryVarName).Assign().Lit(expectedQuery),
	)

	expectArguments := len(expectedArgs) > 0
	if expectArguments {
		x := jen.ID(expectedArgsVarName).Assign().Index().Interface()
		block = append(block, x.Valuesln(expectedArgs...))
	}

	returnedVars := []jen.Code{jen.ID(actualQueryVarName)}
	if expectArguments {
		returnedVars = append(returnedVars, jen.ID(actualArgsVarName))
	}

	block = append(
		block,
		jen.List(returnedVars...).Assign().ID(dbi).Dotf("build%sQuery", queryName).Call(callArgs...),
		jen.Line(),
		jen.ID("ensureArgCountMatchesQuery").Call(
			jen.ID("t"),
			jen.ID(actualQueryVarName),
			func() jen.Code {
				if expectArguments {
					return jen.ID(actualArgsVarName)
				}
				return jen.Index().Interface().Values()
			}(),
		),
		utils.AssertEqual(jen.ID(expectedQueryVarName), jen.ID(actualQueryVarName), nil),
		func() jen.Code {
			if expectArguments {
				return utils.AssertEqual(jen.ID(expectedArgsVarName), jen.ID(actualArgsVarName), nil)
			}
			return jen.Null()
		}(),
	)

	lines := []jen.Code{
		jen.Func().IDf("Test%s_build%sQuery", dbn, queryName).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext("happy path", block...),
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
		panic(fmt.Sprintf("invalid database type! %q", db.LowercaseAbbreviation()))
	}
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
		panic(fmt.Sprintf("invalid database type! %q", db.LowercaseAbbreviation()))
	}
}
