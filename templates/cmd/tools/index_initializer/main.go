package indexinitializer

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"path/filepath"
	"strings"
)

const (
	packageName = "main"

	basePackagePath = "cmd/tools/index_initializer"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"main.go": mainDotGo(proj),
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	return nil
}

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.PackageComment(`Command index_initializer is a CLI that takes in some data via flags about your 
database and the type you want to index, and hydrates a Bleve index full of that type.
This tool is to be used in the event of some data corruption that takes the search index
out of commission.`)

	code.ImportName(constants.FlagParsingLibrary, "flag")
	code.ImportAlias(constants.FlagParsingLibrary, "flag")

	validTypes := []jen.Code{}
	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			validTypes = append(validTypes, jen.Lit(typ.Name.RouteName()).MapAssign().Values())
		}
	}

	code.Add(buildVarDeclarations(proj)...)
	code.Add(buildConstDeclarations()...)
	code.Add(buildInit()...)
	code.Add(buildMain(proj)...)

	return code
}

func buildVarDeclarations(proj *models.Project) []jen.Code {
	validTypes := []jen.Code{}
	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			validTypes = append(validTypes, jen.Lit(typ.Name.RouteName()).MapAssign().Values())
		}
	}

	return []jen.Code{
		jen.Var().Defs(
			jen.ID("indexOutputPath").String(),
			jen.ID("typeName").String(),
			jen.Line(),
			jen.ID("dbConnectionDetails").String(),
			jen.ID("databaseType").String(),
			jen.Line(),
			jen.ID("deadline").Qual("time", "Duration"),
			jen.Line(),
			jen.ID("validTypeNames").Equals().Map(jen.String()).Struct().Valuesln(
				validTypes...,
			),
			jen.Line(),
			jen.ID("validDatabaseTypes").Equals().Map(jen.String()).Struct().Valuesln(
				jen.Qual(proj.InternalConfigPackage(), "PostgresProviderKey").MapAssign().Values(),
				jen.Qual(proj.InternalConfigPackage(), "MariaDBProviderKey").MapAssign().Values(),
				jen.Qual(proj.InternalConfigPackage(), "SqliteProviderKey").MapAssign().Values(),
			),
		),
		jen.Line(),
	}
}

func buildConstDeclarations() []jen.Code {
	return []jen.Code{
		jen.Const().Defs(
			jen.ID("outputPathVerboseFlagName").Equals().Lit("output"),
			jen.ID("dbConnectionVerboseFlagName").Equals().Lit("db_connection"),
			jen.ID("dbTypeVerboseFlagName").Equals().Lit("db_type"),
		),
		jen.Line(),
	}
}

func buildInit() []jen.Code {
	return []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.Qual(constants.FlagParsingLibrary, "StringVarP").Call(
				jen.AddressOf().ID("indexOutputPath"),
				jen.ID("outputPathVerboseFlagName"),
				jen.Lit("o"),
				jen.EmptyString(),
				jen.Lit("output path for bleve index"),
			),
			jen.Qual(constants.FlagParsingLibrary, "StringVarP").Call(
				jen.AddressOf().ID("typeName"),
				jen.Lit("type"),
				jen.Lit("t"),
				jen.EmptyString(),
				jen.Lit("which type to create bleve index for"),
			),
			jen.Line(),

			jen.Qual(constants.FlagParsingLibrary, "StringVarP").Call(
				jen.AddressOf().ID("dbConnectionDetails"),
				jen.ID("dbConnectionVerboseFlagName"),
				jen.Lit("c"),
				jen.EmptyString(),
				jen.Lit("connection string for the relevant database"),
			),
			jen.Qual(constants.FlagParsingLibrary, "StringVarP").Call(
				jen.AddressOf().ID("databaseType"),
				jen.ID("dbTypeVerboseFlagName"),
				jen.Lit("b"),
				jen.EmptyString(),
				jen.Lit("which type of database to connect to"),
			),
			jen.Line(),
			jen.Qual(constants.FlagParsingLibrary, "DurationVarP").Call(
				jen.AddressOf().ID("deadline"),
				jen.Lit("deadline"),
				jen.Lit("d"),
				jen.Qual("time", "Minute"),
				jen.Lit("amount of time to spend adding to the index"),
			),
		),
		jen.Line(),
	}
}

func buildMain(proj *models.Project) []jen.Code {
	searchTypeNames := []string{}
	switchCases := []jen.Code{}

	for _, typ := range proj.DataTypes {
		pn := typ.Name.Plural()
		if typ.SearchEnabled {
			searchTypeNames = append(searchTypeNames, fmt.Sprintf("'%s'", typ.Name.RouteName()))

			switchCases = append(switchCases,
				jen.Case(jen.Lit(typ.Name.RouteName())).Body(
					jen.ID("outputChan").Assign().Make(jen.Chan().Index().Qual(proj.TypesPackage(), typ.Name.Singular())),
					jen.If(
						jen.ID("queryErr").Assign().ID("dbClient").Dotf("GetAll%s", pn).Call(
							constants.CtxVar(),
							jen.ID("outputChan"),
						),
						jen.ID("queryErr").DoesNotEqual().Nil(),
					).Body(
						jen.Qual("log", "Fatalf").Call(jen.Lit("error fetching "+typ.Name.PluralCommonName()+" from database: %v"), jen.Err()),
					),
					jen.Line(),
					jen.For().Body(
						jen.Select().Body(
							jen.Case(jen.ID(typ.Name.PluralUnexportedVarName()).Assign().ReceiveFromChannel().ID("outputChan")).Body(
								jen.For(jen.List(jen.Underscore(), jen.ID("x").Assign().Range().ID(typ.Name.PluralUnexportedVarName()))).Body(
									jen.If(
										jen.ID("searchIndexErr").Assign().ID("im").Dot("Index").Call(
											constants.CtxVar(),
											jen.ID("x").Dot("ID"),
											jen.ID("x"),
										),
										jen.ID("searchIndexErr").DoesNotEqual().Nil(),
									).Body(
										jen.ID(constants.LoggerVarName).
											Dot("WithValue").Call(jen.Lit("id"), jen.ID("x").Dot("ID")).
											Dot("Error").Call(jen.ID("searchIndexErr"), jen.Lit("error adding to search index")),
									),
								),
							),
							jen.Case(jen.ReceiveFromChannel().Qual("time", "After").Call(jen.ID("deadline"))).Body(
								jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("terminating")),
								jen.Return(),
							),
						),
					),
				),
			)
		}
	}

	switchCases = append(switchCases,
		jen.Default().Body(
			jen.Qual("log", "Fatal").Call(jen.Lit("this should never occur")),
		),
	)

	return []jen.Code{
		jen.Func().ID("main").Params().Body(
			jen.Qual(constants.FlagParsingLibrary, "Parse").Call(),
			jen.ID(constants.LoggerVarName).Assign().Qual(filepath.Join(proj.InternalLoggingPackage(), "zerolog"), "NewZeroLogger").Call().Dot("WithName").Call(jen.Lit("search_index_initializer")),
			constants.CreateCtx(),
			jen.Line(),
			jen.If(jen.ID("indexOutputPath").IsEqualTo().EmptyString()).Body(
				jen.Qual("log", "Fatalf").Call(jen.Lit("No output path specified, please provide one via the --%s flag"), jen.ID("outputPathVerboseFlagName")),
				jen.Return(),
			).Else().If(jen.List(jen.Underscore(), jen.ID("ok")).Assign().ID("validTypeNames").Index(jen.ID("typeName")), jen.Not().ID("ok")).Body(
				jen.Qual("log", "Fatalf").Call(
					jen.Lit("Invalid type name %q specified, one of [ "+strings.Join(searchTypeNames, ", ")+" ] expected"),
					jen.ID("typeName"),
				),
				jen.Return(),
			).Else().If(jen.ID("dbConnectionDetails").IsEqualTo().EmptyString()).Body(
				jen.Qual("log", "Fatalf").Call(
					jen.Lit("No database connection details %q specified, please provide one via the --%s flag"),
					jen.ID("dbConnectionDetails"),
					jen.ID("dbConnectionVerboseFlagName"),
				),
				jen.Return(),
			).Else().If(jen.List(jen.Underscore(), jen.ID("ok")).Assign().ID("validDatabaseTypes").Index(jen.ID("databaseType")), jen.Not().ID("ok")).Body(
				jen.Qual("log", "Fatalf").Call(
					jen.Lit("Invalid database type %q specified, please provide one via the --%s flag"),
					jen.ID("databaseType"),
					jen.ID("dbTypeVerboseFlagName"),
				),
				jen.Return(),
			),
			jen.Line(),
			jen.List(jen.ID("im"), jen.Err().Assign().Qual(proj.InternalSearchPackage("bleve"), "NewBleveIndexManager").Call(
				jen.Qual(proj.InternalSearchPackage(), "IndexPath").Call(jen.ID("indexOutputPath")),
				jen.Qual(proj.InternalSearchPackage(), "IndexName").Call(jen.ID("typeName")),
				jen.ID(constants.LoggerVarName),
			)),
			jen.If(jen.Err().DoesNotEqual().Nil()).Body(
				jen.Qual("log", "Fatal").Call(jen.Err()),
			),
			jen.Line(),
			jen.ID("cfg").Assign().AddressOf().Qual(proj.InternalConfigPackage(), "ServerConfig").Valuesln(
				jen.ID("Database").MapAssign().Qual(proj.InternalConfigPackage(), "DatabaseSettings").Valuesln(
					jen.ID("Provider").MapAssign().ID("databaseType"),
					jen.ID("ConnectionDetails").MapAssign().Qual(proj.DatabasePackage(), "ConnectionDetails").Call(jen.ID("dbConnectionDetails")),
				),
				jen.ID("Metrics").MapAssign().Qual(proj.InternalConfigPackage(), "MetricsSettings").Valuesln(
					jen.ID("DBMetricsCollectionInterval").MapAssign().Qual("time", "Second"),
				),
			),
			jen.Line(),
			jen.Comment("connect to our database."),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("connecting to database")),
			jen.List(jen.ID("rawDB"), jen.Err()).Assign().ID("cfg").Dot("ProvideDatabaseConnection").Call(
				jen.ID(constants.LoggerVarName),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Body(
				jen.Qual("log", "Fatalf").Call(jen.Lit("error establishing connection to database: %v"), jen.Err()),
			),
			jen.Line(),
			jen.Comment("establish the database client."),
			jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("setting up database client")),
			jen.List(jen.ID("dbClient"), jen.Err()).Assign().ID("cfg").Dot("ProvideDatabaseClient").Call(
				constants.CtxVar(),
				jen.ID(constants.LoggerVarName),
				jen.ID("rawDB"),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Body(
				jen.Qual("log", "Fatalf").Call(jen.Lit("error initializing database client: %v"), jen.Err()),
			),
			jen.Line(),
			jen.Switch(jen.ID("typeName")).Body(
				switchCases...,
			),
		),
	}
}
