package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

var (
	// files are all the available files to generate
	files = map[string]*jen.File{
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
)

// RenderPackage renders the package
func RenderPackage(types []models.DataType) {
	for path, file := range files {
		renderFile(path, file)
	}

	for _, typ := range types {
		renderFile(fmt.Sprintf("client/v1/http/%s.go", typ.Name.PluralRouteName), itemsDotGo(typ))
		renderFile(fmt.Sprintf("client/v1/http/%s_test.go", typ.Name.PluralRouteName), itemsDotGo(typ))
	}
}

func renderFile(path string, file *jen.File) {
	fp := fmt.Sprintf("/home/jeffrey/src/gitlab.com/verygoodsoftwarenotvirus/naff/templates/%s", path)
	_ = os.Remove(fp)

	var b bytes.Buffer
	if err := file.Render(&b); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fp, b.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

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

func itemsDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

	vn := typ.Name.UnexportedVarName
	pvn := typ.Name.PluralUnexportedVarName
	prn := typ.Name.PluralRouteName
	lp := strings.ToLower(typ.Name.Plural)   // lower plural
	tp := typ.Name.Plural                    // title plural
	ls := strings.ToLower(typ.Name.Singular) // lower singular
	ts := typ.Name.Singular                  // title singular

	basePath := fmt.Sprintf("%sBasePath", pvn)

	addImports(ret)
	ret.Add(jen.Const().Defs(
		jen.ID(basePath).Op("=").Lit(prn)),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", ts)).Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Get%s retrieves an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", ts)).Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID(vn).Op("*").Qual(modelsPkg, ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(
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
			jen.Line(),
			jen.If(jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(vn),
			),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(
					jen.ID("nil"),
					jen.ID("retrieveErr"),
				),
			),
			jen.Line(),
			jen.Return().List(
				jen.ID(vn),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching %s", tp, lp)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", tp)).Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("filter").Dot("ToValues").Call(),
				jen.ID(basePath),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Get%s retrieves a list of %s", tp, lp)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Get%s", tp)).Params(
			ctxParam(),
			jen.ID("filter").Op("*").Qual(modelsPkg, "QueryFilter"),
		).Params(
			jen.ID(pvn).Op("*").Qual(modelsPkg, fmt.Sprintf("%sList", ts)),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
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
			jen.Line(),
			jen.If(
				jen.ID("retrieveErr").Op(":=").ID("c").Dot("retrieve").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID(pvn),
				),
				jen.ID("retrieveErr").Op("!=").ID("nil"),
			).Block(
				jen.Return().List(jen.ID("nil"),
					jen.ID("retrieveErr"),
				),
			),
			jen.Line(),
			jen.Return().List(
				jen.ID(pvn),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildCreate%sRequest builds an HTTP request for creating an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildCreate%sRequest", ts)).Params(
			ctxParam(),
			jen.ID("body").Op("*").Qual(modelsPkg, fmt.Sprintf("%sCreationInput", ts)),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("body"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Create%s creates an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Create%s", ts)).Params(
			ctxParam(),
			jen.ID("input").Op("*").Qual(modelsPkg, fmt.Sprintf("%sCreationInput", ts)),
		).Params(
			jen.ID(vn).Op("*").Qual(modelsPkg, ts),
			jen.ID("err").ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
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
			jen.Line(),
			jen.ID("err").Op("=").ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID(vn),
			),
			jen.Return().List(
				jen.ID(vn),
				jen.ID("err"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildUpdate%sRequest builds an HTTP request for updating an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildUpdate%sRequest", ts)).Params(
			ctxParam(),
			jen.ID("updated").Op("*").Qual(modelsPkg, ts),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("updated").Dot("ID"),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("buildDataRequest").Call(
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("updated"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Update%s updates an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Update%s", ts)).Params(
			ctxParam(),
			jen.ID("updated").Op("*").Qual(modelsPkg, ts),
		).Params(
			jen.ID("error"),
		).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("updated"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildArchive%sRequest builds an HTTP request for updating an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildArchive%sRequest", ts)).Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.Op("*").Qual("net/http", "Request"),
			jen.ID("error"),
		).Block(
			jen.ID("uri").Op(":=").ID("c").Dot("BuildURL").Call(
				jen.ID("nil"),
				jen.ID(basePath),
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.Lit(10),
				),
			),
			jen.Line(),
			jen.Return().Qual("net/http", "NewRequest").Call(
				jen.Qual("net/http", "MethodDelete"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("Archive%s archives an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("Archive%s", ts)).Params(
			ctxParam(),
			jen.ID("id").ID("uint64"),
		).Params(jen.ID("error")).Block(
			jen.List(
				jen.ID("req"),
				jen.ID("err"),
			).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
				jen.ID("ctx"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("building request: %w"),
					jen.ID("err"),
				),
			),
			jen.Line(),
			jen.Return().ID("c").Dot("executeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
		),
	)
	return ret
}

func itemsTestDotGo(typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

	// vn := typ.Name.UnexportedVarName
	// pvn := typ.Name.PluralUnexportedVarName
	prn := typ.Name.PluralRouteName
	// lp := strings.ToLower(typ.Name.Plural)   // lower plural
	tp := typ.Name.Plural // title plural
	// ls := strings.ToLower(typ.Name.Singular) // lower singular
	ts := typ.Name.Singular // title singular

	// routes
	modelRoute := fmt.Sprintf("/api/v1/%s/", prn) + "%d"
	modelListRoute := fmt.Sprintf("/api/v1/%s", prn)

	addImports(ret)

	ret.Add(
		testFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(
					jen.Lit(1),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_Get%s", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, ts).ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.Line(),
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
							jen.Lit(modelRoute),
							jen.ID("expected").Dot("ID"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					requireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(
					jen.ID("ctx"),
					jen.ID("expected").Dot("ID"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", tp)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_Get%s", tp)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, fmt.Sprintf("%sList", ts)).ValuesLn(
					jen.ID(fmt.Sprintf(tp)).Op(":").Index().Qual(modelsPkg, ts).ValuesLn(
						jen.ValuesLn(
							jen.ID("ID").Op(":").Lit(1),
							jen.ID("Name").Op(":").Lit("example"),
							jen.ID("Details").Op(":").Lit("blah"),
						),
					),
				),
				jen.Line(),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(modelListRoute),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					requireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
						nil,
					),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", tp)).Call(
					jen.ID("ctx"),
					jen.ID("nil"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_BuildCreate%sRequest", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, fmt.Sprintf("%sCreationInput", ts)).ValuesLn(
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
				).Op(":=").ID("c").Dot(fmt.Sprintf("BuildCreate%sRequest", ts)).Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_Create%s", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, ts).ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, fmt.Sprintf("%sCreationInput", ts)).ValuesLn(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("Details").Op(":").ID("expected").Dot("Details"),
				),
				jen.Line(),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(modelListRoute),
						jen.Lit("expected and actual path don't match"),
					),
					assertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Var().ID("x").Op("*").Qual(modelsPkg, fmt.Sprintf("%sCreationInput", ts)),
					requireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(
							jen.ID("req").Dot("Body"),
						).Dot("Decode").Call(
							jen.Op("&").ID("x"),
						),
						nil,
					),
					assertEqual(
						jen.ID("exampleInput"),
						jen.ID("x"),
						nil,
					),
					jen.Line(),
					requireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
						nil,
					),
					writeHeader("StatusOK"),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot(fmt.Sprintf("Create%s", ts)).Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_BuildUpdate%sRequest", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodPut"),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(modelsPkg, ts).ValuesLn(
					jen.ID("Name").Op(":").Lit("changed name"),
					jen.ID("Details").Op(":").Lit("changed details"),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				),
				jen.List(
					jen.ID("actual"),
					jen.ID("err"),
				).Op(":=").ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(
					jen.ID("ctx"),
					jen.ID("exampleInput"),
				),
				jen.Line(),
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
		testFunc(fmt.Sprintf("V1Client_Update%s", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(modelsPkg, ts).ValuesLn(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.Line(),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(modelRoute),
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
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot(fmt.Sprintf("Update%s", ts)).Call(
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
		testFunc(fmt.Sprintf("V1Client_BuildArchive%sRequest", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				expectMethod("expectedMethod", "MethodDelete"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
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
				).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(
					jen.ID("ctx"),
					jen.ID("expectedID"),
				),
				jen.Line(),
				requireNotNil(jen.ID("actual"), nil),
				requireNotNil(
					jen.ID("actual").Dot("URL"),
					nil,
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
		testFunc(fmt.Sprintf("V1Client_Archive%s", ts)).Block(
			parallelTest(nil),
			jen.Line(),
			buildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(1)),
				jen.Line(),
				buildTestServer(
					"ts",
					assertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(modelRoute),
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
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot(fmt.Sprintf("Archive%s", ts)).Call(
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
