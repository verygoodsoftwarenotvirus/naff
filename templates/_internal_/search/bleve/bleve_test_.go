package bleve

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bleveTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("exampleType").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("ID").ID("uint64"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
			jen.Newline(),
			jen.ID("exampleTypeWithStringID").Struct(
				jen.ID("ID").ID("string"),
				jen.ID("Name").ID("string"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(constants.TestSuitePackage, "AfterTest").Op("=").Parens(jen.Op("*").ID("bleveIndexManagerTestSuite")).Call(jen.ID("nil")),
			jen.ID("_").Qual(constants.TestSuitePackage, "BeforeTest").Op("=").Parens(jen.Op("*").ID("bleveIndexManagerTestSuite")).Call(jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().ID("bleveIndexManagerTestSuite").Struct(
			jen.Qual(constants.TestSuitePackage, "Suite"),
			jen.Newline(),
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("indexPath").ID("string"),
			jen.ID("exampleAccountID").ID("uint64"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("createTmpIndexPath").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.ID("string")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Newline(),
			jen.List(jen.ID("tmpIndexPath"), jen.ID("err")).Op(":=").Qual("os", "MkdirTemp").Call(
				jen.Lit(""),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("bleve-testidx-%d"),
					jen.Qual("time", "Now").Call().Dot("Unix").Call(),
				),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Newline(),
			jen.Return().ID("tmpIndexPath"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("BeforeTest").Params(jen.List(jen.ID("_"), jen.ID("_")).ID("string")).Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.ID("s").Dot("indexPath").Op("=").ID("createTmpIndexPath").Call(jen.ID("t")),
			jen.Newline(),
			jen.ID("err").Op(":=").Qual("os", "MkdirAll").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Octal(700),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Newline(),
			jen.ID("s").Dot("ctx").Op("=").Qual("context", "Background").Call(),
			jen.ID("s").Dot("exampleAccountID").Op("=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call().Dot("ID"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("AfterTest").Params(jen.List(jen.ID("_"), jen.ID("_")).ID("string")).Body(
			jen.ID("s").Dot("Require").Call().Dot("NoError").Call(jen.Qual("os", "RemoveAll").Call(jen.ID("s").Dot("indexPath")))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestNewBleveIndexManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Qual(constants.TestSuitePackage, "Run").Call(
				jen.ID("T"),
				jen.ID("new").Call(jen.ID("bleveIndexManagerTestSuite")),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestNewBleveIndexManagerWithTestIndex").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("constructor_test_happy_path_test.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	for _, typ := range proj.DataTypes {
		code.Add(
			jen.Newline(),
			jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).IDf("TestNewBleveIndexManagerWith%sIndex", typ.Name.Plural()).Params().Body(
				jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
				jen.Newline(),
				jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
					jen.ID("s").Dot("indexPath"),
					jen.Litf("constructor_test_happy_path_%s.bleve", typ.Name.PluralRouteName()),
				)),
				jen.Newline(),
				jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
					jen.ID("exampleIndexPath"),
					jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sSearchIndexName", typ.Name.Plural())),
					jen.Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.ID("err"),
				),
			),
			jen.Newline(),
		)
	}

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestNewBleveIndexManagerWithInvalidName").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Lit("constructor_test_invalid_name.bleve")),
			jen.Newline(),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.Lit("invalid"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("Error").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestIndex").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.Const().ID("exampleQuery").Op("=").Lit("_test"),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("_test_obligatory.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.ID("x").Op(":=").Op("&").ID("exampleType").Valuesln(
				jen.ID("ID").Op(":").Lit(123),
				jen.ID("Name").Op(":").ID("exampleQuery"), jen.ID("BelongsToUser").Op(":").ID("s").Dot("exampleAccountID"),
			),
			jen.Newline(),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Dot("Index").Call(
					jen.ID("s").Dot("ctx"),
					jen.ID("x").Dot("ID"),
					jen.ID("x"),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestSearch").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.Const().ID("exampleQuery").Op("=").Lit("search_test"),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("search_test_obligatory.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.ID("x").Op(":=").ID("exampleType").Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").ID("exampleQuery"), jen.ID("BelongsToUser").Op(":").ID("s").Dot("exampleAccountID")),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Dot("Index").Call(
					jen.ID("s").Dot("ctx"),
					jen.ID("x").Dot("ID"),
					jen.Op("&").ID("x"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("results"), jen.ID("err")).Op(":=").ID("im").Dot("Search").Call(
				jen.ID("s").Dot("ctx"),
				jen.ID("x").Dot("Name"),
				jen.ID("s").Dot("exampleAccountID"),
			),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("results"),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestSearchWithInvalidQuery").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("search_test_invalid_query.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.List(jen.ID("results"), jen.ID("err")).Op(":=").ID("im").Dot("Search").Call(
				jen.ID("s").Dot("ctx"),
				jen.Lit(""),
				jen.ID("s").Dot("exampleAccountID"),
			),
			jen.ID("assert").Dot("Empty").Call(
				jen.ID("t"),
				jen.ID("results"),
			),
			jen.ID("assert").Dot("Error").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestSearchWithEmptyIndexAndSearch").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("search_test_empty_index.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.List(jen.ID("results"), jen.ID("err")).Op(":=").ID("im").Dot("Search").Call(
				jen.ID("s").Dot("ctx"),
				jen.Lit("example"),
				jen.ID("s").Dot("exampleAccountID"),
			),
			jen.ID("assert").Dot("Empty").Call(
				jen.ID("t"),
				jen.ID("results"),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestSearchWithClosedIndex").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.Const().ID("exampleQuery").Op("=").Lit("search_test"),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("search_test_closed_index.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.ID("x").Op(":=").Op("&").ID("exampleType").Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").ID("exampleQuery"), jen.ID("BelongsToUser").Op(":").ID("s").Dot("exampleAccountID")),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Dot("Index").Call(
					jen.ID("s").Dot("ctx"),
					jen.ID("x").Dot("ID"),
					jen.ID("x"),
				),
			),
			jen.Newline(),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Assert(jen.Op("*").ID("bleveIndexManager")).Dot("index").Dot("Close").Call(),
			),
			jen.Newline(),
			jen.List(jen.ID("results"), jen.ID("err")).Op(":=").ID("im").Dot("Search").Call(
				jen.ID("s").Dot("ctx"),
				jen.ID("x").Dot("Name"),
				jen.ID("s").Dot("exampleAccountID"),
			),
			jen.ID("assert").Dot("Empty").Call(
				jen.ID("t"),
				jen.ID("results"),
			),
			jen.ID("assert").Dot("Error").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestSearchWithInvalidID").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.Const().ID("exampleQuery").Op("=").Lit("search_test"),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("search_test_invalid_id.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.ID("x").Op(":=").Op("&").ID("exampleTypeWithStringID").Valuesln(jen.ID("ID").Op(":").Lit("whatever"), jen.ID("Name").Op(":").ID("exampleQuery"), jen.ID("BelongsToUser").Op(":").ID("s").Dot("exampleAccountID")),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Assert(jen.Op("*").ID("bleveIndexManager")).Dot("index").Dot("Index").Call(
					jen.ID("x").Dot("ID"),
					jen.ID("x"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("results"), jen.ID("err")).Op(":=").ID("im").Dot("Search").Call(
				jen.ID("s").Dot("ctx"),
				jen.ID("x").Dot("Name"),
				jen.ID("s").Dot("exampleAccountID"),
			),
			jen.ID("assert").Dot("Empty").Call(
				jen.ID("t"),
				jen.ID("results"),
			),
			jen.ID("assert").Dot("Error").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestSearchForAdmin").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.Const().ID("exampleQuery").Op("=").Lit("search_test"),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("search_test_obligatory.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.ID("x").Op(":=").ID("exampleType").Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").ID("exampleQuery"), jen.ID("BelongsToUser").Op(":").ID("s").Dot("exampleAccountID")),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Dot("Index").Call(
					jen.ID("s").Dot("ctx"),
					jen.ID("x").Dot("ID"),
					jen.Op("&").ID("x"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("results"), jen.ID("err")).Op(":=").ID("im").Dot("SearchForAdmin").Call(
				jen.ID("s").Dot("ctx"),
				jen.ID("x").Dot("Name"),
			),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("results"),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("bleveIndexManagerTestSuite")).ID("TestDelete").Params().Body(
			jen.ID("t").Op(":=").ID("s").Dot("T").Call(),
			jen.Newline(),
			jen.Const().ID("exampleQuery").Op("=").Lit("delete_test"),
			jen.ID("exampleIndexPath").Op(":=").ID("search").Dot("IndexPath").Call(jen.Qual("path/filepath", "Join").Call(
				jen.ID("s").Dot("indexPath"),
				jen.Lit("delete_test.bleve"),
			)),
			jen.Newline(),
			jen.List(jen.ID("im"), jen.ID("err")).Op(":=").ID("NewBleveIndexManager").Call(
				jen.ID("exampleIndexPath"),
				jen.ID("testingSearchIndexName"),
				jen.ID("logging").Dot("NewNoopLogger").Call(),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("im"),
			),
			jen.Newline(),
			jen.ID("x").Op(":=").Op("&").ID("exampleType").Valuesln(jen.ID("ID").Op(":").Lit(123), jen.ID("Name").Op(":").ID("exampleQuery"), jen.ID("BelongsToUser").Op(":").ID("s").Dot("exampleAccountID")),
			jen.Newline(),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Dot("Index").Call(
					jen.ID("s").Dot("ctx"),
					jen.ID("x").Dot("ID"),
					jen.ID("x"),
				),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("im").Dot("Delete").Call(
					jen.ID("s").Dot("ctx"),
					jen.ID("x").Dot("ID"),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
