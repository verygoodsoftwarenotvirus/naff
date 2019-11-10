package utils

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func ExampleValueForField(field models.DataField) jen.Code {
	switch field.Type {
	case "bool":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(false))
		}
		return jen.Lit(false)
	case "string":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit("example"))
		}
		return jen.Lit("example")
	case "float32":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1.23))
		}
		return jen.Lit(1.23)
	case "float64":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1.23))
		}
		return jen.Lit(1.23)
	case "uint":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "uint8":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "uint16":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "uint32":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "uint64":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "int":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "int8":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "int16":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "int32":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	case "int64":
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(jen.Lit(1))
		}
		return jen.Lit(1)
	default:
		panic(fmt.Sprintf("unknown type: %q", field.Type))
	}
}

func FakeCallForField(pkgRoot string, field models.DataField) jen.Code {
	switch field.Type {
	case "bool":
		x := jen.Qual(FakeLibrary, "Bool").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "string":
		x := jen.Qual(FakeLibrary, "Word").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "float32":
		x := jen.Qual(FakeLibrary, "Float32").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "float64":
		x := jen.Qual(FakeLibrary, "Float64").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint":
		x := jen.ID("uint").Call(jen.Qual(FakeLibrary, "Uint32").Call())
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint8":
		x := jen.Qual(FakeLibrary, "Uint8").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint16":
		x := jen.Qual(FakeLibrary, "Uint16").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint32":
		x := jen.Qual(FakeLibrary, "Uint32").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "uint64":
		x := jen.ID("uint64").Call(jen.Qual(FakeLibrary, "Uint32").Call())
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int":
		x := jen.ID("int").Call(jen.Qual(FakeLibrary, "Int32").Call())
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int8":
		x := jen.Qual(FakeLibrary, "Int8").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int16":
		x := jen.Qual(FakeLibrary, "Int16").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int32":
		x := jen.Qual(FakeLibrary, "Int32").Call()
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	case "int64":
		x := jen.ID("int64").Call(jen.Qual(FakeLibrary, "Int32").Call())
		if field.Pointer {
			return jen.Func().Params(jen.ID("x").ID(field.Type)).Params(jen.Op("*").ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID("x"))).Call(x)
		}
		return x
	default:
		return jen.Null()
	}
}
