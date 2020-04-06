package queriers

import (
	"errors"
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

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, vendor := range []string{postgres /*, sqlite, mariadb*/} { // UNDOME
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
