package jen

import "fmt"

// RawString renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func RawString(v string) *Statement {
	return newStatement().RawString(v)
}

// RawString renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (g *Group) RawString(v string) *Statement {
	s := RawString(v)
	g.items = append(g.items, s)
	return s
}

// RawString renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (s *Statement) RawString(v string) *Statement {
	t := token{
		typ:     literalRawStringToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}

// RawStringf renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func RawStringf(msg string, args ...interface{}) *Statement {
	return newStatement().RawStringf(msg, args...)
}

// RawStringf renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (g *Group) RawStringf(msg string, args ...interface{}) *Statement {
	s := RawStringf(msg, args...)
	g.items = append(g.items, s)
	return s
}

// RawStringf renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (s *Statement) RawStringf(msg string, args ...interface{}) *Statement {
	t := token{
		typ:     literalRawStringToken,
		content: fmt.Sprintf(msg, args...),
	}
	*s = append(*s, t)
	return s
}

// Octal renders a literal. Octal supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func Octal(v interface{}) *Statement {
	return newStatement().Octal(v)
}

// Octal renders a literal. Octal supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (g *Group) Octal(v interface{}) *Statement {
	s := Octal(v)
	g.items = append(g.items, s)
	return s
}

// Octal renders a literal. Octal supports only built-in types (int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64).
// Passing any other type will panic.
func (s *Statement) Octal(v interface{}) *Statement {
	t := token{
		typ:     octalToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}

// Lit renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func Lit(v interface{}) *Statement {
	return newStatement().Lit(v)
}

// Lit renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (g *Group) Lit(v interface{}) *Statement {
	s := Lit(v)
	g.items = append(g.items, s)
	return s
}

// Lit renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (s *Statement) Lit(v interface{}) *Statement {
	t := token{
		typ:     literalToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}

// EmptyString renders an empty string literal.
func EmptyString() *Statement {
	return newStatement().EmptyString()
}

// EmptyString renders an empty string literal.
func (g *Group) EmptyString() *Statement {
	s := EmptyString()
	g.items = append(g.items, s)
	return s
}

// EmptyString renders an empty string literal.
func (s *Statement) EmptyString() *Statement {
	t := token{
		typ:     literalToken,
		content: "",
	}
	*s = append(*s, t)
	return s
}

// Zero renders an empty string literal.
func Zero() *Statement {
	return newStatement().Zero()
}

// Zero renders an empty string literal.
func (g *Group) Zero() *Statement {
	s := Zero()
	g.items = append(g.items, s)
	return s
}

// Zero renders an empty string literal.
func (s *Statement) Zero() *Statement {
	t := token{
		typ:     literalToken,
		content: 0,
	}
	*s = append(*s, t)
	return s
}

// One renders an empty string literal.
func One() *Statement {
	return newStatement().One()
}

// One renders an empty string literal.
func (g *Group) One() *Statement {
	s := One()
	g.items = append(g.items, s)
	return s
}

// One renders an empty string literal.
func (s *Statement) One() *Statement {
	t := token{
		typ:     literalToken,
		content: 1,
	}
	*s = append(*s, t)
	return s
}

// Litf renders a literal. Litf supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func Litf(format string, args ...interface{}) *Statement {
	return newStatement().Litf(format, args...)
}

// Litf renders a literal. Litf supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (g *Group) Litf(format string, args ...interface{}) *Statement {
	s := Litf(format, args...)
	g.items = append(g.items, s)
	return s
}

// Litf renders a literal. Litf supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func (s *Statement) Litf(format string, args ...interface{}) *Statement {
	v := fmt.Sprintf(format, args...)
	t := token{
		typ:     literalToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}

// LitFunc renders a literal. LitFunc generates the value to render by executing the provided
// function. LitFunc supports only built-in types (bool, string, int, complex128, float64, float32,
// int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Returning any other type will panic.
func LitFunc(f func() interface{}) *Statement {
	return newStatement().LitFunc(f)
}

// LitFunc renders a literal. LitFunc generates the value to render by executing the provided
// function. LitFunc supports only built-in types (bool, string, int, complex128, float64, float32,
// int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Returning any other type will panic.
func (g *Group) LitFunc(f func() interface{}) *Statement {
	s := LitFunc(f)
	g.items = append(g.items, s)
	return s
}

// LitFunc renders a literal. LitFunc generates the value to render by executing the provided
// function. LitFunc supports only built-in types (bool, string, int, complex128, float64, float32,
// int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Returning any other type will panic.
func (s *Statement) LitFunc(f func() interface{}) *Statement {
	t := token{
		typ:     literalToken,
		content: f(),
	}
	*s = append(*s, t)
	return s
}

// LitRune renders a rune literal.
func LitRune(v rune) *Statement {
	return newStatement().LitRune(v)
}

// LitRune renders a rune literal.
func (g *Group) LitRune(v rune) *Statement {
	s := LitRune(v)
	g.items = append(g.items, s)
	return s
}

// LitRune renders a rune literal.
func (s *Statement) LitRune(v rune) *Statement {
	t := token{
		typ:     literalRuneToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}

// LitRuneFunc renders a rune literal. LitRuneFunc generates the value to
// render by executing the provided function.
func LitRuneFunc(f func() rune) *Statement {
	return newStatement().LitRuneFunc(f)
}

// LitRuneFunc renders a rune literal. LitRuneFunc generates the value to
// render by executing the provided function.
func (g *Group) LitRuneFunc(f func() rune) *Statement {
	s := LitRuneFunc(f)
	g.items = append(g.items, s)
	return s
}

// LitRuneFunc renders a rune literal. LitRuneFunc generates the value to
// render by executing the provided function.
func (s *Statement) LitRuneFunc(f func() rune) *Statement {
	t := token{
		typ:     literalRuneToken,
		content: f(),
	}
	*s = append(*s, t)
	return s
}

// LitByte renders a byte literal.
func LitByte(v byte) *Statement {
	return newStatement().LitByte(v)
}

// LitByte renders a byte literal.
func (g *Group) LitByte(v byte) *Statement {
	s := LitByte(v)
	g.items = append(g.items, s)
	return s
}

// LitByte renders a byte literal.
func (s *Statement) LitByte(v byte) *Statement {
	t := token{
		typ:     literalByteToken,
		content: v,
	}
	*s = append(*s, t)
	return s
}

// LitByteFunc renders a byte literal. LitByteFunc generates the value to
// render by executing the provided function.
func LitByteFunc(f func() byte) *Statement {
	return newStatement().LitByteFunc(f)
}

// LitByteFunc renders a byte literal. LitByteFunc generates the value to
// render by executing the provided function.
func (g *Group) LitByteFunc(f func() byte) *Statement {
	s := LitByteFunc(f)
	g.items = append(g.items, s)
	return s
}

// LitByteFunc renders a byte literal. LitByteFunc generates the value to
// render by executing the provided function.
func (s *Statement) LitByteFunc(f func() byte) *Statement {
	t := token{
		typ:     literalByteToken,
		content: f(),
	}
	*s = append(*s, t)
	return s
}
