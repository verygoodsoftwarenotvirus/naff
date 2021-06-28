package utils

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
)

// FakeSeedFunc builds a consistent fake library seed init function
func FakeSeedFunc() jen.Code {
	return jen.Func().ID("init").Params().Body(
		InlineFakeSeedFunc(),
	)
}

// InlineFakeSeedFunc builds a consistent fake library seed init function
func InlineFakeSeedFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call())
}

func FakeFuncForType(typ string, isPointer bool) func() jen.Code {
	switch typ {
	case "string":
		if isPointer {
			return ValuePointerWrapper(typ, FakeStringFunc())
		}
		return FakeStringFunc
	case "bool":
		if isPointer {
			return ValuePointerWrapper(typ, FakeBoolFunc())
		}
		return FakeBoolFunc
	case "int":
		if isPointer {
			return ValuePointerWrapper(typ, FakeIntFunc())
		}
		return FakeIntFunc
	case "int8":
		if isPointer {
			return ValuePointerWrapper(typ, FakeInt8Func())
		}
		return FakeInt8Func
	case "int16":
		if isPointer {
			return ValuePointerWrapper(typ, FakeInt16Func())
		}
		return FakeInt16Func
	case "int32":
		if isPointer {
			return ValuePointerWrapper(typ, FakeInt32Func())
		}
		return FakeInt32Func
	case "int64":
		if isPointer {
			return ValuePointerWrapper(typ, FakeInt64Func())
		}
		return FakeInt64Func
	case "uint":
		if isPointer {
			return ValuePointerWrapper(typ, FakeUintFunc())
		}
		return FakeUintFunc
	case "uint8":
		if isPointer {
			return ValuePointerWrapper(typ, FakeUint8Func())
		}
		return FakeUint8Func
	case "uint16":
		if isPointer {
			return ValuePointerWrapper(typ, FakeUint16Func())
		}
		return FakeUint16Func
	case "uint32":
		if isPointer {
			return ValuePointerWrapper(typ, FakeUint32Func())
		}
		return FakeUint32Func
	case "uint64":
		if isPointer {
			return ValuePointerWrapper(typ, FakeUint64Func())
		}
		return FakeUint64Func
	case "float32":
		if isPointer {
			return ValuePointerWrapper(typ, FakeFloat32Func())
		}
		return FakeFloat32Func
	case "float64":
		if isPointer {
			return ValuePointerWrapper(typ, FakeFloat64Func())
		}
		return FakeFloat64Func
	default:
		panic(fmt.Sprintf("unknown type! %q", typ))
	}
}

func ValuePointerWrapper(typ string, c jen.Code) func() jen.Code {
	return func() jen.Code {
		return jen.Func().Params(jen.ID("x").ID(typ)).Params(jen.PointerTo().ID(typ)).SingleLineBody(jen.Return(jen.AddressOf().ID("x"))).Call(c)
	}
}

func FakeStringFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Word").Call()
}

func FakeContentTypeFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "FileMimeType").Call()
}

func FakeUUIDFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "UUID").Call()
}

func FakeURLFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "URL").Call()
}

func FakeHTTPMethodFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "HTTPMethod").Call()
}

func FakeBoolFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Bool").Call()
}

func FakeIntFunc() jen.Code {
	return jen.Int().Call(jen.Qual(constants.FakeLibrary, "Int32").Call())
}

func FakeInt8Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Int8").Call()
}

func FakeInt16Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Int16").Call()
}

func FakeInt32Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Int32").Call()
}

func FakeInt64Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Int64").Call()
}

func FakeInt64WhichIsReallyAnInt32Func() jen.Code {
	return jen.Int64().Call(jen.Qual(constants.FakeLibrary, "Int32").Call())
}

func FakeUintFunc() jen.Code {
	return jen.Uint().Call(jen.Qual(constants.FakeLibrary, "Uint32").Call())
}

func FakeUint8Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Uint8").Call()
}

func FakeUint16Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Uint16").Call()
}

func FakeUint32Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Uint32").Call()
}

func FakeUint64WhichIsActuallyAUint32Func() jen.Code {
	return jen.Uint64().Call(jen.Qual(constants.FakeLibrary, "Uint32").Call())
}

func FakeUint64Func() jen.Code {
	return jen.Uint64().Call(jen.Qual(constants.FakeLibrary, "Uint32").Call())
}

func FakeFloat32Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Float32").Call()
}

func FakeFloat64WhichIsActuallyAFloat32Func() jen.Code {
	return jen.Float64().Call(jen.Qual(constants.FakeLibrary, "Float32").Call())
}

func FakeFloat64Func() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Float64").Call()
}

func FakeUsernameFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Username").Call()
}

func FakeUnixTimeFunc() jen.Code {
	return jen.Uint64().Call(jen.Uint32().Call(jen.Qual(constants.FakeLibrary, "Date").Call().Dot("Unix").Call()))
}

func FakePasswordFunc() jen.Code {
	return jen.Qual(constants.FakeLibrary, "Password").Call(
		jen.True(),
		jen.True(),
		jen.True(),
		jen.True(),
		jen.True(),
		jen.Lit(32),
	)
}
