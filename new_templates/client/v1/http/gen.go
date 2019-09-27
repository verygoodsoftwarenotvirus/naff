package client

import (
	jen "github.com/dave/jennifer/jen"
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
	v1             = "V1Client"
	coreOAuth2Pkg  = "golang.org/x/oauth2"
	loggingPkg     = "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	noopLoggingPkg = "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	assertPkg      = "github.com/stretchr/testify/assert"
	mustAssertPkg  = "github.com/stretchr/testify/require"
	mockPkg        = "github.com/stretchr/testify/mock"
)

func addImports(file *jen.File) {
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
	ret.Add(jen.Var().Id("itemsBasePath").Op("=").Lit("items"))

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildGetItemRequest").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("itemsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("id"),
					jen.Lit(10),
				),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetItem").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Id("item").Op("*").Id("models").Dot("Item"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetItemRequest").Call(
				jen.Id("ctx"),
				jen.Id("id"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				)),
			),
			jen.If(jen.Id("retrieveErr").Op(":=").Id("c").Dot("retrieve").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Op("&").Id("item"),
			),
				jen.Id("retrieveErr").Op("!=").Id("nil"),
			).Block(
				jen.Return().List(
					jen.Id("nil"),
					jen.Id("retrieveErr"),
				),
			),
			jen.Return().List(
				jen.Id("item"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildGetItemsRequest").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("filter").Dot("ToValues").Call(),
				jen.Id("itemsBasePath"),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("GetItems").Params(
			ctxParam(),
			jen.Id("filter").Op("*").Id("models").Dot("QueryFilter"),
		).Params(
			jen.Id("items").Op("*").Id("models").Dot("ItemList"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildGetItemsRequest").Call(
				jen.Id("ctx"),
				jen.Id("filter"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.If(
				jen.Id("retrieveErr").Op(":=").Id("c").Dot("retrieve").Call(
					jen.Id("ctx"),
					jen.Id("req"),
					jen.Op("&").Id("items"),
				),
				jen.Id("retrieveErr").Op("!=").Id("nil"),
			).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Id("retrieveErr"),
				),
			), jen.Return().List(
				jen.Id("items"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildCreateItemRequest").Params(
			ctxParam(),
			jen.Id("body").Op("*").Id("models").Dot("ItemCreationInput"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("itemsBasePath"),
			), jen.Return().Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.Id("uri"),
				jen.Id("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("CreateItem").Params(
			ctxParam(),
			jen.Id("input").Op("*").Id("models").Dot("ItemCreationInput"),
		).Params(
			jen.Id("item").Op("*").Id("models").Dot("Item"),
			jen.Id("err").Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildCreateItemRequest").Call(
				jen.Id("ctx"),
				jen.Id("input"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().List(jen.Id("nil"),
					jen.Qual("fmt", "Errorf").Call(
						jen.Lit("building request: %w"),
						jen.Id("err"),
					),
				),
			),
			jen.Id("err").Op("=").Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"), jen.Op("&").Id("item"),
			), jen.Return().List(jen.Id("item"),
				jen.Id("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildUpdateItemRequest").Params(
			ctxParam(),
			jen.Id("updated").Op("*").Id("models").Dot("Item"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("itemsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("updated").Dot("ID"), jen.Lit(10),
				),
			), jen.Return().Id("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.Id("uri"),
				jen.Id("updated"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("UpdateItem").Params(
			ctxParam(),
			jen.Id("updated").Op("*").Id("models").Dot("Item"),
		).Params(
			jen.Id("error"),
		).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildUpdateItemRequest").Call(
				jen.Id("ctx"),
				jen.Id("updated"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				),
			), jen.Return().Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"), jen.Op("&").Id("updated"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(
			jen.Id("c").Op("*").Id(v1),
		).Id("BuildArchiveItemRequest").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.Id("error"),
		).Block(
			jen.Id("uri").Op(":=").Id("c").Dot("BuildURL").Call(
				jen.Id("nil"),
				jen.Id("itemsBasePath"),
				jen.Qual("strconv", "FormatUint").Call(
					jen.Id("id"),
					jen.Lit(10),
				),
			), jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.Id("uri"),
				jen.Id("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(""),
		jen.Line(),
		jen.Func().Params(jen.Id("c").Op("*").Id(v1)).Id("ArchiveItem").Params(
			ctxParam(),
			jen.Id("id").Id("uint64"),
		).Params(jen.Id("error")).Block(
			jen.List(
				jen.Id("req"),
				jen.Id("err"),
			).Op(":=").Id("c").Dot("BuildArchiveItemRequest").Call(
				jen.Id("ctx"),
				jen.Id("id"),
			),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.Id("err"),
				),
			), jen.Return().Id("c").Dot("executeRequest").Call(
				jen.Id("ctx"),
				jen.Id("req"),
				jen.Id("nil"),
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
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				jen.List(jen.Id("actual"), jen.Id("err")).Op(":=").Id("c").Dot("BuildGetItemRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertTrue(
					jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(
						jen.Id("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.Id("expectedID"),
						),
					),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_GetItem").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Item").Values(jen.Dict{
					jen.Id("ID"):      jen.Lit(1),
					jen.Id("Name"):    jen.Lit("example"),
					jen.Id("Details"): jen.Lit("blah"),
				}),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(
					jen.Qual("net/http", "HandlerFunc").Call(
						jen.Func().Params(
							jen.Id("res").Qual("net/http", "ResponseWriter"),
							jen.Id("req").Op("*").Qual("net/http", "Request"),
						).Block(
							assertTrue(
								jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(
									jen.Id("req").Dot("URL").Dot("String").Call(),
									jen.Qual("strconv", "Itoa").Call(
										jen.Id("int").Call(
											jen.Id("expected").Dot("ID"),
										),
									),
								),
							),
							assertEqual(
								jen.Id("t"),
								jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("/api/v1/items/%d"),
									jen.Id("expected").Dot("ID"),
								), jen.Lit("expected and actual path don't match"),
							),
							assertEqual(
								jen.Id("t"),
								jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet"),
							),
							jen.Id("require").Dot("NoError").Call(
								jen.Id("t"),
								jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("res")).Dot("Encode").Call(jen.Id("expected")),
							),
						),
					),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetItem").Call(
					jen.Id("ctx"),
					jen.Id("expected").Dot("ID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildGetItemsRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildGetItemsRequest").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_GetItems").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("ItemList").Values(jen.Dict{
					jen.Id("Items"): jen.Index().Id("models").Dot("Item").Values(jen.Dict{
						jen.Id("ID"):      jen.Lit(1),
						jen.Id("Name"):    jen.Lit("example"),
						jen.Id("Details"): jen.Lit("blah"),
					}),
				}),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/api/v1/items"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodGet"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("GetItems").Call(
					jen.Id("ctx"),
					jen.Id("nil"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildCreateItemRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("ItemCreationInput").Values(
					jen.Id("Name").Op(":").Lit("expected name"),
					jen.Id("Details").Op(":").Lit("expected details")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildCreateItemRequest").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_CreateItem").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Item").Values(
					jen.Id("ID").Op(":").Lit(1),
					jen.Id("Name").Op(":").Lit("example"),
					jen.Id("Details").Op(":").Lit("blah"),
				),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("ItemCreationInput").Values(
					jen.Id("Name").Op(":").Id("expected").Dot("Name"),
					jen.Id("Details").Op(":").Id("expected").Dot("Details"),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"),
							jen.Lit("/api/v1/items"),
							jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPost"),
						), jen.Var().Id("x").Op("*").Id("models").Dot("ItemCreationInput"),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewDecoder").Call(
								jen.Id("req").Dot("Body"),
							).Dot("Decode").Call(
								jen.Op("&").Id("x"),
							),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("exampleInput"),
							jen.Id("x"),
						),
						jen.Id("require").Dot("NoError").Call(
							jen.Id("t"), jen.Qual("encoding/json", "NewEncoder").Call(
								jen.Id("res"),
							).Dot("Encode").Call(
								jen.Id("expected"),
							),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("CreateItem").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(jen.Id("expected"), jen.Id("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildUpdateItemRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPut"),
				createCtx(),
				jen.Id("exampleInput").Op(":=").Op("&").Id("models").Dot("Item").Values(
					jen.Id("Name").Op(":").Lit("changed name"),
					jen.Id("Details").Op(":").Lit("changed details"),
				),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildUpdateItemRequest").Call(
					jen.Id("ctx"),
					jen.Id("exampleInput"),
				),
				requireNotNil(jen.Id("actual"), nil),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_UpdateItem").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Op("&").Id("models").Dot("Item").Values(
					jen.Id("ID").Op(":").Lit(1),
					jen.Id("Name").Op(":").Lit("example"),
					jen.Id("Details").Op(":").Lit("blah"),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/items/%d"),
								jen.Id("expected").Dot("ID"),
							), jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodPut"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("err").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				).Dot("UpdateItem").Call(
					jen.Id("ctx"),
					jen.Id("expected"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_BuildArchiveItemRequest").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Id("nil")),
				jen.Id("expectedID").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				jen.Id("c").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				),
				jen.List(
					jen.Id("actual"),
					jen.Id("err"),
				).Op(":=").Id("c").Dot("BuildArchiveItemRequest").Call(
					jen.Id("ctx"),
					jen.Id("expectedID"),
				),
				requireNotNil(jen.Id("actual"), nil),
				jen.Id("require").Dot("NotNil").Call(
					jen.Id("t"),
					jen.Id("actual").Dot("URL"),
				),
				assertTrue(
					jen.Id("t"), jen.Qual("strings", "HasSuffix").Call(
						jen.Id("actual").Dot("URL").Dot("String").Call(), jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.Id("expectedID"),
						),
					),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
				assertEqual(
					jen.Id("t"),
					jen.Id("actual").Dot("Method"),
					jen.Id("expectedMethod"),
					jen.Lit("request should be a %s request"),
					jen.Id("expectedMethod"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		testFunc("V1Client_ArchiveItem").Block(
			jen.Id("T").Dot("Parallel").Call(),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.Id("expected").Op(":=").Id("uint64").Call(
					jen.Lit(1),
				),
				createCtx(),
				jen.Id("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(
						jen.Id("res").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).Block(
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("URL").Dot("Path"), jen.Qual("fmt", "Sprintf").Call(
								jen.Lit("/api/v1/items/%d"),
								jen.Id("expected"),
							), jen.Lit("expected and actual path don't match"),
						),
						assertEqual(
							jen.Id("t"),
							jen.Id("req").Dot("Method"), jen.Qual("net/http", "MethodDelete"),
						),
						writeHeader("StatusOK"),
					),
				),
				),
				jen.Id("err").Op(":=").Id("buildTestClient").Call(
					jen.Id("t"),
					jen.Id("ts"),
				).Dot("ArchiveItem").Call(
					jen.Id("ctx"),
					jen.Id("expected"),
				),
				assertNoError(
					jen.Id("t"),
					jen.Id("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
