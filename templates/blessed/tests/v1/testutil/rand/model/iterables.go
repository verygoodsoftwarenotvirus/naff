package model

import (
	"fmt"
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(pkg *models.Project, typ models.DataType) *jen.File {
	ret := jen.NewFile("randmodel")

	utils.AddImports(pkg, ret)
	sn := typ.Name.Singular()

	ret.Add(
		jen.Commentf("Random%sCreationInput creates a random %sInput", sn, sn),
		jen.Line(),
		jen.Func().IDf("Random%sCreationInput", sn).Params().Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn))).Block(
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)).Valuesln(buildFakeCalls(typ.Fields)...),
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
		case "bool",
			"float32",
			"float64",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"int",
			"int8",
			"int16",
			"int32",
			"int64",
			"string":
			numberMethods := map[string]jen.Code{
				"string":  jen.Qual(utils.FakeLibrary, "Word").Call(),
				"bool":    jen.Qual(utils.FakeLibrary, "Bool").Call(),
				"float32": jen.Qual(utils.FakeLibrary, "Float32").Call(),
				"float64": jen.Qual(utils.FakeLibrary, "Float64").Call(),
				"uint":    jen.ID("uint").Call(jen.Qual(utils.FakeLibrary, "Uint32").Call()),
				"uint8":   jen.Qual(utils.FakeLibrary, "Uint8").Call(),
				"uint16":  jen.Qual(utils.FakeLibrary, "Uint16").Call(),
				"uint32":  jen.Qual(utils.FakeLibrary, "Uint32").Call(),
				"uint64":  jen.Qual(utils.FakeLibrary, "Uint64").Call(),
				"int":     jen.ID("int").Call(jen.Qual(utils.FakeLibrary, "Int32").Call()),
				"int8":    jen.Qual(utils.FakeLibrary, "Int8").Call(),
				"int16":   jen.Qual(utils.FakeLibrary, "Int16").Call(),
				"int32":   jen.Qual(utils.FakeLibrary, "Int32").Call(),
				"int64":   jen.Qual(utils.FakeLibrary, "Int64").Call(),
			}

			if field.Pointer {
				out = append(
					out,
					jen.ID(sn).Op(":").Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(numberMethods[typ]),
				)
			} else {
				out = append(
					out,
					jen.ID(sn).Op(":").Add(numberMethods[typ]),
				)
			}
		default:
			panic(fmt.Sprintf("unaccounted for type!: %q", strings.ToLower(field.Type)))
		}
	}

	return out
}
