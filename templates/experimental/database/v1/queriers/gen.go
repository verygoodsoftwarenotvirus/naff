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
)

// DatabasePackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	for _, vendor := range []string{postgres, sqlite} {
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
		vendorWord *wordsmith.SuperPalabra
	)

	switch vendor {
	case postgres:
		dbDesc = "Postgres instances"
		vendorWord = wordsmith.FromSingularPascalCase("Postgres")
	case sqlite:
		dbDesc = "sqlite files"
		vendorWord = wordsmith.FromSingularPascalCase("Sqlite")
	default:
		panic("wtf")
	}

	files := map[string]*jen.File{
		fmt.Sprintf("database/v1/queriers/%s/oauth2_clients.go", vendor):      oauth2ClientsDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/%s.go", vendor, vendor):          databaseDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/webhooks.go", vendor):            webhooksDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/wire.go", vendor):                wireDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/doc.go", vendor):                 docDotGo(vendorWord.SingularPackageName(), dbDesc),
		fmt.Sprintf("database/v1/queriers/%s/%s_test.go", vendor, vendor):     databaseTestDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/users.go", vendor):               usersDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/users_test.go", vendor):          usersTestDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/webhooks_test.go", vendor):       webhooksTestDotGo(vendorWord),
		fmt.Sprintf("database/v1/queriers/%s/migrations.go", vendor):          migrationsDotGo(vendorWord, types),
		fmt.Sprintf("database/v1/queriers/%s/oauth2_clients_test.go", vendor): oauth2ClientsTestDotGo(vendorWord),
	}

	for _, typ := range types {
		files[fmt.Sprintf("database/v1/queriers/%s/%s.go", vendor, typ.Name.PluralRouteName())] = iterablesDotGo(vendorWord, typ)
		files[fmt.Sprintf("database/v1/queriers/%s/%s_test.go", vendor, typ.Name.PluralRouteName())] = iterablesTestDotGo(vendorWord, typ)
	}

	for path, file := range files {
		if err := utils.RenderFile(pkgRoot, path, file); err != nil {
			return err
		}
	}

	return nil
}
