package queriers

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	postgres = "postgres"
	sqlite   = "sqlite"
	mariadb  = "mariadb"
)

// DatabasePackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	for _, vendor := range []string{postgres, sqlite, mariadb} {
		if err := renderDatabasePackage(pkgRoot, vendor, types); err != nil {
			return err
		}
	}
	return nil
}

// renderDatabasePackage renders the package
func renderDatabasePackage(pkgRoot, vendor string, types []models.DataType) error {
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
	pn := vendorWord.SingularPackageName()

	files := map[string]*jen.File{
		fmt.Sprintf("database/v1/queriers/%s/oauth2_clients.go", vendor):      oauth2ClientsDotGo(pkgRoot, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/%s.go", pn, pn):                  databaseDotGo(pkgRoot, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/webhooks.go", vendor):            webhooksDotGo(pkgRoot, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/wire.go", vendor):                wireDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/doc.go", vendor):                 docDotGo(pn, dbDesc),
		fmt.Sprintf("database/v1/queriers/%s/%s_test.go", pn, pn):             databaseTestDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/users.go", vendor):               usersDotGo(pkgRoot, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/users_test.go", vendor):          usersTestDotGo(pkgRoot, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/webhooks_test.go", vendor):       webhooksTestDotGo(pkgRoot, vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/migrations.go", vendor):          migrationsDotGo(vendorWord, types),
		fmt.Sprintf("database/v1/queriers/%s/oauth2_clients_test.go", vendor): oauth2ClientsTestDotGo(pkgRoot, vendorWord),
	}

	for _, typ := range types {
		files[fmt.Sprintf("database/v1/queriers/%s/%s.go", vendor, typ.Name.PluralRouteName())] = iterablesDotGo(pkgRoot, vendorWord, typ)
		files[fmt.Sprintf("database/v1/queriers/%s/%s_test.go", vendor, typ.Name.PluralRouteName())] = iterablesTestDotGo(pkgRoot, vendorWord, typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(pkgRoot, path, file); err != nil {
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
