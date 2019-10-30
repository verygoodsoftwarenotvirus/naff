package model

import (
	"fmt"
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(rootPkg string, typ models.DataType) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(ret)

	ret.Add(
		jen.Comment("RandomItemCreationInput creates a random ItemInput"),
		jen.Line(),
		jen.Func().ID("RandomItemCreationInput").Params().Params(jen.Op("*").Qual(filepath.Join(rootPkg, "models/v1"), "ItemCreationInput")).Block(
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(rootPkg, "models/v1"), "ItemCreationInput").Valuesln(buildFakeCalls(typ.Fields)...),
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
				jen.ID(sn).Op(":").Qual("github.com/icrowley/fake", "Word").Call(),
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
