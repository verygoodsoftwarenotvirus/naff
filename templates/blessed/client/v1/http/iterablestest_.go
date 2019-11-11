package client

import (
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesTestDotGo(pkgRoot string, typ models.DataType) *jen.File {
	ret := jen.NewFile("client")

	prn := typ.Name.PluralRouteName()
	tp := typ.Name.Plural()   // title plural
	ts := typ.Name.Singular() // title singular

	// routes
	modelRoute := fmt.Sprintf("/api/v1/%s/", prn) + "%d"
	modelListRoute := fmt.Sprintf("/api/v1/%s", prn)

	utils.AddImports(pkgRoot, []models.DataType{typ}, ret)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildGet%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodGet"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.ID("expectedID").Op(":=").ID("uint64").Call(jen.Lit(1)),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", ts)).Call(jen.ID("ctx"), jen.ID("expectedID")),
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
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), ts).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
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
					utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodGet"), nil),
					utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")), nil),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("Get%s", ts)).Call(jen.ID("ctx"), jen.ID("expected").Dot("ID")),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildGet%sRequest", tp)).Call(jen.ID("ctx"), jen.ID("nil")),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
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
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sList", ts)).Valuesln(
					jen.ID(tp).Op(":").Index().Qual(filepath.Join(pkgRoot, "models/v1"), ts).Valuesln(
						jen.Valuesln(
							jen.ID("ID").Op(":").Lit(1),
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
					jen.ID("t"),
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

	creationFields := func() []jen.Code {
		lines := []jen.Code{
			jen.ID("ID").Op(":").Lit(1),
		}

		for _, field := range typ.Fields {
			lines = append(lines, jen.ID(field.Name.Singular()).Op(":").Add(utils.ExampleValueForField(field)))
		}
		return lines
	}
	cfs := creationFields()

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildCreate%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				utils.ExpectMethod("expectedMethod", "MethodPost"),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.Line(),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", ts)).Valuesln(
					cfs[1:]...,
				),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(
					jen.ID("t"),
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

	createCreationInputLines := func() []jen.Code {
		var lines []jen.Code
		for _, field := range typ.Fields {
			if field.ValidForCreationInput {
				lines = append(lines, jen.ID(field.Name.Singular()).Op(":").ID("expected").Dot(field.Name.Singular()))
			}
		}
		return lines
	}

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_Create%s", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), ts).Valuesln(
					cfs...,
				),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", ts)).Valuesln(
					createCreationInputLines()...,
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(jen.ID("req").Dot("URL").Dot("Path"), jen.Lit(modelListRoute), jen.Lit("expected and actual path don't match")),
					utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPost"), nil),
					jen.Line(),
					jen.Var().ID("x").Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", ts)),
					utils.RequireNoError(jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("x")), nil),
					utils.AssertEqual(jen.ID("exampleInput"), jen.ID("x"), nil),
					jen.Line(),
					utils.RequireNoError(jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("expected")), nil),
				),
				jen.Line(),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("Create%s", ts)).Call(jen.ID("ctx"), jen.ID("exampleInput")),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		utils.OuterTestFunc(fmt.Sprintf("V1Client_BuildUpdate%sRequest", ts)).Block(
			utils.ParallelTest(nil),
			jen.Line(),
			utils.BuildSubTest("happy path", utils.ExpectMethod("expectedMethod", "MethodPut"),
				jen.ID("exampleInput").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), ts).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
				),
				jen.Line(),
				jen.ID("ts").Op(":=").Qual("net/http/httptest", "NewTLSServer").Call(jen.ID("nil")),
				jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildUpdate%sRequest", ts)).Call(jen.ID("ctx"), jen.ID("exampleInput")),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
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
			utils.BuildSubTest("happy path",
				jen.ID("expected").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), ts).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
				),
				jen.Line(),
				utils.BuildTestServer(
					"ts",
					utils.AssertEqual(
						jen.ID("req").Dot("URL").Dot("Path"),
						jen.Qual("fmt", "Sprintf").Call(jen.Lit(modelRoute), jen.ID("expected").Dot("ID")),
						jen.Lit("expected and actual path don't match"),
					),
					utils.AssertEqual(jen.ID("req").Dot("Method"), jen.Qual("net/http", "MethodPut"), nil),
					utils.AssertNoError(
						jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), ts).Values()),
						nil,
					),
				),
				jen.Line(),
				jen.ID("err").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Update%s", ts)).Call(jen.ID("ctx"), jen.ID("expected")),
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
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
				jen.ID("c").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(fmt.Sprintf("BuildArchive%sRequest", ts)).Call(jen.ID("ctx"), jen.ID("expectedID")),
				jen.Line(),
				utils.RequireNotNil(jen.ID("actual"), nil),
				utils.RequireNotNil(jen.ID("actual").Dot("URL"), nil),
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
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
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
						jen.Qual("fmt", "Sprintf").Call(jen.Lit(modelRoute), jen.ID("expected")),
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
				jen.ID("err").Op(":=").ID("buildTestClient").Call(jen.ID("t"), jen.ID("ts")).Dot(fmt.Sprintf("Archive%s", ts)).Call(jen.ID("ctx"), jen.ID("expected")),
				utils.AssertNoError(jen.ID("err"), jen.Lit("no error should be returned")),
			),
		),
		jen.Line(),
	)

	return ret
}
