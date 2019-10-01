package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	// T is the big T
	T  = "T"
	t  = "t"
	v1 = "V1Client"
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

	utils.AddImports(ret)
	ret.Add(jen.Const().Defs(
		jen.ID(basePath).Op("=").Lit(prn)),
	)

	ret.Add(
		jen.Comment(fmt.Sprintf("BuildGet%sRequest builds an HTTP request for fetching an %s", ts, ls)),
		jen.Line(),
		newClientMethod(fmt.Sprintf("BuildGet%sRequest", ts)).Params(
			utils.CtxParam(),
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
			utils.CtxParam(),
			jen.ID("id").ID("uint64"),
		).Params(
			jen.ID(vn).Op("*").Qual(utils.ModelsPkg, ts),
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
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
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
			utils.CtxParam(),
			jen.ID("filter").Op("*").Qual(utils.ModelsPkg, "QueryFilter"),
		).Params(
			jen.ID(pvn).Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sList", ts)),
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
			utils.CtxParam(),
			jen.ID("body").Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)),
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
			utils.CtxParam(),
			jen.ID("input").Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)),
		).Params(
			jen.ID(vn).Op("*").Qual(utils.ModelsPkg, ts),
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
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(utils.ModelsPkg, ts),
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
			utils.CtxParam(),
			jen.ID("updated").Op("*").Qual(utils.ModelsPkg, ts),
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
			utils.CtxParam(),
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
			utils.CtxParam(),
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

	utils.AddImports(ret)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(utils.ModelsPkg, ts).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertTrue(
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
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(modelRoute),
							jen.ID("expected").Dot("ID"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", tp)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Get%s", tp)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(utils.ModelsPkg, fmt.Sprintf("%sList", ts)).Valuesln(
					jen.ID(fmt.Sprintf(tp)).Op(":").Index().Qual(utils.ModelsPkg, ts).Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(1),
							jen.ID("Name").Op(":").Lit("example"),
							jen.ID("Details").Op(":").Lit("blah"),
						),
					),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(modelListRoute),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodGet"),
						nil,
					),
					utils.RequireNoError(
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildCreate%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)).Valuesln(
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Create%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(utils.ModelsPkg, ts).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)).Valuesln(
					jen.ID("Name").Op(":").ID("expected").Dot("Name"),
					jen.ID("Details").Op(":").ID("expected").Dot("Details"),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Lit(modelListRoute),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPost"),
						nil,
					),
					jen.Line(),
					jen.Var().ID("x").Op("*").Qual(utils.ModelsPkg, fmt.Sprintf("%sCreationInput", ts)),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewDecoder").Call(
							jen.ID("req").Dot("Body"),
						).Dot("Decode").Call(
							jen.Op("&").ID("x"),
						),
						nil,
					),
					utils.AssertEqual(
						jen.ID("exampleInput"),
						jen.ID("x"),
						nil,
					),
					jen.Line(),
					utils.RequireNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")),
						nil,
					),
					utils.WriteHeader("StatusOK"),
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(jen.ID("expected"),
					jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildUpdate%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPut"),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(utils.ModelsPkg, ts).Valuesln(
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Update%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(utils.ModelsPkg, ts).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("Name").Op(":").Lit("example"),
					jen.ID("Details").Op(":").Lit("blah"),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(modelRoute),
							jen.ID("expected").Dot("ID"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodPut"),
						nil,
					),
					utils.WriteHeader("StatusOK"),
				),
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot(fmt.Sprintf("Update%s", ts)).Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildArchive%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodDelete"),
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
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(
					jen.ID("actual").Dot("URL"),
					nil,
				),
				utils.AssertTrue(
					jen.Qual("strings", "HasSuffix").Call(
						jen.ID("actual").Dot("URL").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%d"),
							jen.ID("expectedID"),
						),
					),
					nil,
				),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
				utils.AssertEqual(
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
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Archive%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(1)),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit(modelRoute),
							jen.ID("expected"),
						),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(
						jen.ID("req").Dot("Method"),
						jen.Qual("net/http", "MethodDelete"),
						nil,
					),
					utils.WriteHeader("StatusOK"),
				),
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(
					jen.ID(t),
					jen.ID("ts"),
				).Dot(fmt.Sprintf("Archive%s", ts)).Call(
					jen.ID("ctx"),
					jen.ID("expected"),
				),
				utils.AssertNoError(
					jen.ID("err"),
					jen.Lit("no error should be returned"),
				),
			),
		),
		jen.Line(),
	)

	return ret
}
