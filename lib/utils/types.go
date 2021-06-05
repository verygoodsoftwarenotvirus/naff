package utils

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func FakeCallForField(pkgRoot string, field models.DataField) jen.Code {
	const varName = "x"

	switch field.Type {
	case "bool":
		x := jen.Qual(constants.FakeLibrary, "Bool").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "string":
		x := jen.Qual(constants.FakeLibrary, "Word").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "float32":
		x := jen.Qual(constants.FakeLibrary, "Float32").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "float64":
		x := jen.Qual(constants.FakeLibrary, "Float64").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "uint":
		x := jen.Uint().Call(jen.Qual(constants.FakeLibrary, "Uint32").Call())
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "uint8":
		x := jen.Qual(constants.FakeLibrary, "Uint8").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "uint16":
		x := jen.Qual(constants.FakeLibrary, "Uint16").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "uint32":
		x := jen.Qual(constants.FakeLibrary, "Uint32").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "uint64":
		x := jen.Uint64().Call(jen.Qual(constants.FakeLibrary, "Uint32").Call())
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "int":
		x := jen.Int().Call(jen.Qual(constants.FakeLibrary, "Int32").Call())
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "int8":
		x := jen.Qual(constants.FakeLibrary, "Int8").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "int16":
		x := jen.Qual(constants.FakeLibrary, "Int16").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "int32":
		x := jen.Qual(constants.FakeLibrary, "Int32").Call()
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	case "int64":
		x := jen.Int64().Call(jen.Qual(constants.FakeLibrary, "Int32").Call())
		if field.IsPointer {
			return jen.Func().Params(jen.ID(varName).ID(field.Type)).Params(jen.PointerTo().ID(field.Type)).SingleLineBlock(jen.Return(jen.Op("&").ID(varName))).Call(x)
		}
		return x
	default:
		return jen.Null()
	}
}
