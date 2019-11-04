package model

import (
	"fmt"
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(pkgRoot string, typ models.DataType) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(pkgRoot, []models.DataType{typ}, ret)
	sn := typ.Name.Singular()

	ret.Add(
		jen.Commentf("Random%sCreationInput creates a random %sInput", sn, sn),
		jen.Line(),
		jen.Func().IDf("Random%sCreationInput", sn).Params().Params(jen.Op("*").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", sn))).Block(
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkgRoot, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(buildFakeCalls(typ.Fields)...),
			jen.Line(),
			jen.Return().ID("x"),
		),
		jen.Line(),
	)
	return ret
}

func buildFakeCalls(fields []models.DataField) []jen.Code {
	var out []jen.Code

	for _, field := range fields {
		sn := field.Name.Singular()
		typ := strings.ToLower(field.Type)
		switch typ {
		case "string":
			out = append(
				out,
				jen.ID(sn).Op(":").Qual(utils.FakeLibrary, "Word").Call(),
			)
		case "float32", "float64", "int", "int32", "int64", "uint32", "uint64":
			numberMethods := map[string]string{
				"float32": "Float32",
				"float64": "Float64",
				"int":     "Int",
				"int32":   "Int31",
				"int64":   "Int63",
				"uint32":  "Uint32",
				"uint64":  "Uint64",
			}

			out = append(
				out,
				jen.ID(sn).Op(":").Qual("math/rand", numberMethods[typ]).Call(),
			)
		default:
			panic(fmt.Sprintf("unaccounted for type!: %q", strings.ToLower(field.Type)))
		}
	}

	return out
}
