package client

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
)

var (
	// Files are all the available files to generate
	Files = map[string]*jen.File{
		"client/v1/http/main.go":                mainDotGo(),
		"client/v1/http/main_test.go":           mainTestDotGo(),
		"client/v1/http/helpers.go":             helpersDotGo(),
		"client/v1/http/helpers_test.go":        helpersTestDotGo(),
		"client/v1/http/users.go":               usersDotGo(),
		"client/v1/http/users_test.go":          usersTestDotGo(),
		"client/v1/http/roundtripper.go":        roundtripperDotGo(),
		"client/v1/http/webhooks.go":            webhooksDotGo(),
		"client/v1/http/webhooks_test.go":       webhooksTestDotGo(),
		"client/v1/http/oauth2_clients.go":      oauth2ClientsDotGo(),
		"client/v1/http/oauth2_clients_test.go": oauth2ClientsTestDotGo(),
	}
)

const (
	// T is the big T
	T  = "T"
	t  = "t"
	v1 = "V1Client"

	// packages
	coreOAuth2Pkg  = "golang.org/x/oauth2"
	loggingPkg     = "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	noopLoggingPkg = "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	assertPkg      = "github.com/stretchr/testify/assert"
	mustAssertPkg  = "github.com/stretchr/testify/require"
	mockPkg        = "github.com/stretchr/testify/mock"
	modelsPkg      = "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"

	// routes
	itemRoute     = "/api/v1/items/%d"
	itemListRoute = "/api/v1/items"
)

func addImports(file *jen.File) {
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/mock", "mockutil")

	file.ImportNames(map[string]string{
		"context":           "context",
		"fmt":               "fmt",
		"net/http":          "http",
		"net/http/httputil": "httputil",
		"errors":            "errors",
		"net/url":           "url",
		"path":              "path",
		"strings":           "strings",
		"time":              "time",
		"bytes":             "bytes",
		"encoding/json":     "json",
		"io":                "io",
		"io/ioutil":         "ioutil",
		"reflect":           "reflect",
		//
		loggingPkg:     "logging",
		noopLoggingPkg: "noop",
		modelsPkg:      "models",
		//
		assertPkg:                               "assert",
		mustAssertPkg:                           "require",
		mockPkg:                                 "mock",
		coreOAuth2Pkg:                           "oauth2",
		"github.com/moul/http2curl":             "http2curl",
		"go.opencensus.io/plugin/ochttp":        "ochttp",
		"golang.org/x/oauth2/clientcredentials": "clientcredentials",
	})
	file.Add(jen.Line())
}

func itemsDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)
	ret.Add(jen.Var().ID("itemsBasePath").Op("=").Lit("items"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetItemRequest").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("itemsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("GetItem").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID("item").Op("*").Qual(modelsPkg, "Item"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					)),
			),
			jen.If(jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("item"),
			),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("retrieveErr"),
				),
			),
			jen.Return().List(
				jen.ID("item"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildGetItemsRequest").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID("itemsBasePath"),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("GetItems").Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.ID("items").Op("*").Qual(modelsPkg, "ItemList"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildGetItemsRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.If(
				jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("items"),
				),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(jen.ID("nil"),
					jen.ID("retrieveErr"),
				),
			),
			jen.Return().List(
				jen.ID("items"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildCreateItemRequest").Params(
			ctxParam(),
			jen.ID("body").Op("*").Qual(modelsPkg, "ItemCreationInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("itemsBasePath"),
			),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("CreateItem").Params(
			ctxParam(),
			jen.ID("input").Op("*").Qual(modelsPkg, "ItemCreationInput"),
		).Params(
			jen.ID("item").Op("*").Qual(modelsPkg, "Item"),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildCreateItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("item"),
			),
			jen.Return().List(jen.ID("item"),
				jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildUpdateItemRequest").Params(
			ctxParam(),
			jen.ID("updated").Op("*").Qual(modelsPkg, "Item"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("itemsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("updated").Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("updated"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("UpdateItem").Params(
			ctxParam(),
			jen.ID("updated").Op("*").Qual(modelsPkg, "Item"),
		).Params(
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildUpdateItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("updated"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("BuildArchiveItemRequest").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID("itemsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		newClientMethod("ArchiveItem").Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot("BuildArchiveItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
	)
	return ret
}

func itemsTestDotGo() *jen.File {
	ret := jen.NewFile("client")

	addImports(ret)

	ret.Add(
		testFunc("V1Client_BuildGetItemRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.List(jen.ID("actual"),
					jen.ID("err")).Op(":=").ID("c").Dot("BuildGetItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_GetItem").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "Item").ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				buildTestServer(
					"ts",
					assertTrue(
						jen.Qual("strings", "HasSuffix").Call(
							jen.ID("req").Dot("URL").Dot("String").Call(),
							jen.Qual("strconv", "Itoa").Call(
								jen.ID("int").Call(
									jen.ID("expected").Dot("ID"),
								),
							),
						),
						nil,
					),
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(itemRoute),
							jen.ID("expected").Dot("ID"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID(t),
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetItem").Call(
					jen.ID("ctx"),
					jen.ID("expected").Dot("ID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildGetItemsRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildGetItemsRequest").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_GetItems").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "ItemList").ValuesLn(
					jen.ID("Items").Op(":").Index().Qual(modelsPkg, "Item").ValuesLn(
						jen.ID("ID").Op(":").Lit(1),
						jen.ID("Name").Op(":").Lit("example"),
						jen.ID("Details").Op(":").Lit("blah"),
					),
				),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(itemListRoute),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID(t),
						jen.Qual("encoding/json", "NewEncoder").Call(
							jen.ID("res"),
						).Dot("Encode").Call(
							jen.ID("expected"),
						),
					),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("GetItems").Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildCreateItemRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "ItemCreationInput").ValuesLn(
					jen.ID("Name").Op(":").Lit("expected name"),
					jen.ID("Details").Op(":").Lit("expected details"),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildCreateItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_CreateItem").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "Item").ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "ItemCreationInput").ValuesLn(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("Details").Op(":").ID("expected").Dot("Details"),
				),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(itemListRoute),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Var().ID("x").Op("*").Qual(modelsPkg, "ItemCreationInput"),
					jen.ID("require").Dot("NoError").Call(
						jen.ID(t),
						jen.Qual("encoding/json", "NewDecoder").Call(
							jen.ID("req").Dot("Body"),
						).Dot("Decode").Call(
							jen.Op("&").ID("x"),
						),
					),
					assertEqual(
						jen.ID("exampleInput"),
						jen.ID("x"),
						nil,
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID(t),
						jen.Qual("encoding/json", "NewEncoder").Call(
							jen.ID("res"),
						).Dot("Encode").Call(
							jen.ID("expected"),
						),
					),
					writeHeader("StatusOK"),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("CreateItem").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildUpdateItemRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPut"),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, "Item").ValuesLn(
					jen.ID("Name").Op(":").Lit("changed name"),
					jen.ID("Details").Op(":").Lit("changed details"),
				),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildUpdateItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				requireNotNil(jen.ID("actual"), nil),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_UpdateItem").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, "Item").ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(itemRoute),
							jen.ID("expected").Dot("ID"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPut"),
						nil,
					),
					writeHeader("StatusOK"),
				),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot("UpdateItem").Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildArchiveItemRequest").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot("BuildArchiveItemRequest").Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				requireNotNil(jen.ID("actual"), nil),
				jen.ID("require").Dot("NotNil").Call(
					jen.ID(t),
					jen.ID("actual").Dot("URL"),
				),
				assertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.ID("actual").Dot("Method"),
					jen.ID("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.ID("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_ArchiveItem").Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/api/v1/items/%d"),
							jen.ID("expected"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
					writeHeader("StatusOK"),
				),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot("ArchiveItem").Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				assertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
