package jen

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type tokenType string

const (
	packageToken          tokenType = "package"
	identifierToken       tokenType = "identifier"
	qualifiedToken        tokenType = "qualified"
	keywordToken          tokenType = "keyword"
	operatorToken         tokenType = "operator"
	delimiterToken        tokenType = "delimiter"
	literalToken          tokenType = "literal"
	literalRawStringToken tokenType = "literal_raw_string"
	literalRuneToken      tokenType = "literal_rune"
	literalByteToken      tokenType = "literal_byte"
	nullToken             tokenType = "null"
	layoutToken           tokenType = "layout"
)

type token struct {
	typ     tokenType
	content interface{}
}

func (t token) isNull(f *File) bool {
	if t.typ == packageToken {
		// package token is null if the path is a dot-import or the local package path
		return f.isDotImport(t.content.(string)) || f.isLocal(t.content.(string))
	}
	return t.typ == nullToken
}

func (t token) render(f *File, w io.Writer, s *Statement) error {
	switch t.typ {
	case literalToken:
		var out string
		switch t.content.(type) {
		case string:
			ts := fmt.Sprintf("%s", t.content)
			if strings.Contains(ts, `"`) && !strings.Contains(ts, `\"`) {
				out = fmt.Sprintf("`%s`", t.content)
			} else {
				out = fmt.Sprintf("%q", t.content)
			}
		case bool, int, complex128:
			// default constant types can be left bare
			out = fmt.Sprintf("%#v", t.content)
		case float64:
			out = fmt.Sprintf("%#v", t.content)
			if !strings.Contains(out, ".") && !strings.Contains(out, "e") {
				// If the formatted value is not in scientific notation, and does not have a dot, then
				// we add ".0". Otherwise it will be interpreted as an int.
				// See:
				// https://github.com/dave/jennifer/issues/39
				// https://github.com/golang/go/issues/26363
				out += ".0"
			}
		case float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
			// other built-in types need specific type info
			out = fmt.Sprintf("%T(%#v)", t.content, t.content)
		case complex64:
			// fmt package already renders parenthesis for complex64
			out = fmt.Sprintf("%T%#v", t.content, t.content)
		default:
			panic(fmt.Sprintf("unsupported type for literal: %T", t.content))
		}
		if _, err := w.Write([]byte(out)); err != nil {
			return err
		}
	case literalRawStringToken:
		var out string
		switch t.content.(type) {
		case string:
			out = fmt.Sprintf("`%s`", t.content)
		}
		if _, err := w.Write([]byte(out)); err != nil {
			return err
		}
	case literalRuneToken:
		if _, err := w.Write([]byte(strconv.QuoteRune(t.content.(rune)))); err != nil {
			return err
		}
	case literalByteToken:
		if _, err := w.Write([]byte(fmt.Sprintf("byte(%#v)", t.content))); err != nil {
			return err
		}
	case keywordToken, operatorToken, layoutToken, delimiterToken:
		if _, err := w.Write([]byte(fmt.Sprintf("%s", t.content))); err != nil {
			return err
		}
		if t.content.(string) == "default" {
			// Special case for Default, which must always be followed by a colon
			if _, err := w.Write([]byte(":")); err != nil {
				return err
			}
		}
	case packageToken:
		path := t.content.(string)
		alias := f.register(path)
		if _, err := w.Write([]byte(alias)); err != nil {
			return err
		}
	case identifierToken:
		if _, err := w.Write([]byte(t.content.(string))); err != nil {
			return err
		}
	case nullToken: // notest
		// do nothing (should never render a null token)
	}
	return nil
}

// Null adds a null item. Null items render nothing and are not followed by a
// separator in lists.
func Null() *Statement {
	return newStatement().Null()
}

// Null adds a null item. Null items render nothing and are not followed by a
// separator in lists.
func (g *Group) Null() *Statement {
	s := Null()
	g.items = append(g.items, s)
	return s
}

// Null adds a null item. Null items render nothing and are not followed by a
// separator in lists.
func (s *Statement) Null() *Statement {
	t := token{
		typ: nullToken,
	}
	*s = append(*s, t)
	return s
}

// Empty adds an empty item. Empty items render nothing but are followed by a
// separator in lists.
func Empty() *Statement {
	return newStatement().Empty()
}

// Empty adds an empty item. Empty items render nothing but are followed by a
// separator in lists.
func (g *Group) Empty() *Statement {
	s := Empty()
	g.items = append(g.items, s)
	return s
}

// Empty adds an empty item. Empty items render nothing but are followed by a
// separator in lists.
func (s *Statement) Empty() *Statement {
	t := token{
		typ:     operatorToken,
		content: "",
	}
	*s = append(*s, t)
	return s
}

// Op renders the provided operator / token.
func Op(op string) *Statement {
	return newStatement().Op(op)
}

// Op renders the provided operator / token.
func (g *Group) Op(op string) *Statement {
	s := Op(op)
	g.items = append(g.items, s)
	return s
}

// Op renders the provided operator / token.
func (s *Statement) Op(op string) *Statement {
	t := token{
		typ:     operatorToken,
		content: op,
	}
	*s = append(*s, t)
	return s
}

// Spread renders the provided operator / token.
func Spread() *Statement {
	return newStatement().Spread()
}

// Spread renders the provided operator / token.
func (g *Group) Spread() *Statement {
	s := Spread()
	g.items = append(g.items, s)
	return s
}

// Spread renders the provided operator / token.
func (s *Statement) Spread() *Statement {
	t := token{
		typ:     operatorToken,
		content: "...",
	}
	*s = append(*s, t)
	return s
}

// BitwiseXOR renders the provided operator / token.
func BitwiseXOR() *Statement {
	return newStatement().BitwiseXOR()
}

// BitwiseXOR renders the provided operator / token.
func (g *Group) BitwiseXOR() *Statement {
	s := BitwiseXOR()
	g.items = append(g.items, s)
	return s
}

// BitwiseXOR renders the provided operator / token.
func (s *Statement) BitwiseXOR() *Statement {
	t := token{
		typ:     operatorToken,
		content: "^",
	}
	*s = append(*s, t)
	return s
}

// BitwiseOR renders the provided operator / token.
func BitwiseOR() *Statement {
	return newStatement().BitwiseOR()
}

// BitwiseOR renders the provided operator / token.
func (g *Group) BitwiseOR() *Statement {
	s := BitwiseOR()
	g.items = append(g.items, s)
	return s
}

// BitwiseOR renders the provided operator / token.
func (s *Statement) BitwiseOR() *Statement {
	t := token{
		typ:     operatorToken,
		content: "|",
	}
	*s = append(*s, t)
	return s
}

// LeftShift renders the provided operator / token.
func LeftShift() *Statement {
	return newStatement().LeftShift()
}

// LeftShift renders the provided operator / token.
func (g *Group) LeftShift() *Statement {
	s := LeftShift()
	g.items = append(g.items, s)
	return s
}

// LeftShift renders the provided operator / token.
func (s *Statement) LeftShift() *Statement {
	t := token{
		typ:     operatorToken,
		content: "<<",
	}
	*s = append(*s, t)
	return s
}

// RightShift renders the provided operator / token.
func RightShift() *Statement {
	return newStatement().RightShift()
}

// RightShift renders the provided operator / token.
func (g *Group) RightShift() *Statement {
	s := RightShift()
	g.items = append(g.items, s)
	return s
}

// RightShift renders the provided operator / token.
func (s *Statement) RightShift() *Statement {
	t := token{
		typ:     operatorToken,
		content: ">>",
	}
	*s = append(*s, t)
	return s
}

// Plus renders the provided operator / token.
func Plus() *Statement {
	return newStatement().Plus()
}

// Plus renders the provided operator / token.
func (g *Group) Plus() *Statement {
	s := Plus()
	g.items = append(g.items, s)
	return s
}

// Plus renders the provided operator / token.
func (s *Statement) Plus() *Statement {
	t := token{
		typ:     operatorToken,
		content: "+",
	}
	*s = append(*s, t)
	return s
}

// Minus renders the provided operator / token.
func Minus() *Statement {
	return newStatement().Minus()
}

// Minus renders the provided operator / token.
func (g *Group) Minus() *Statement {
	s := Minus()
	g.items = append(g.items, s)
	return s
}

// Minus renders the provided operator / token.
func (s *Statement) Minus() *Statement {
	t := token{
		typ:     operatorToken,
		content: "-",
	}
	*s = append(*s, t)
	return s
}

// And renders the provided operator / token.
func And() *Statement {
	return newStatement().And()
}

// And renders the provided operator / token.
func (g *Group) And() *Statement {
	s := And()
	g.items = append(g.items, s)
	return s
}

// And renders the provided operator / token.
func (s *Statement) And() *Statement {
	t := token{
		typ:     operatorToken,
		content: "&&",
	}
	*s = append(*s, t)
	return s
}

// AddressOf renders the provided operator / token.
func AddressOf() *Statement {
	return newStatement().AddressOf()
}

// AddressOf renders the provided operator / token.
func (g *Group) AddressOf() *Statement {
	s := AddressOf()
	g.items = append(g.items, s)
	return s
}

// AddressOf renders the provided operator / token.
func (s *Statement) AddressOf() *Statement {
	t := token{
		typ:     operatorToken,
		content: "&",
	}
	*s = append(*s, t)
	return s
}

// PointerTo renders the provided operator / token.
func PointerTo() *Statement {
	return newStatement().PointerTo()
}

// PointerTo renders the provided operator / token.
func (g *Group) PointerTo() *Statement {
	s := PointerTo()
	g.items = append(g.items, s)
	return s
}

// PointerTo renders the provided operator / token.
func (s *Statement) PointerTo() *Statement {
	t := token{
		typ:     operatorToken,
		content: "*",
	}
	*s = append(*s, t)
	return s
}

// Times renders the provided operator / token.
func Times() *Statement {
	return newStatement().Times()
}

// Times renders the provided operator / token.
func (g *Group) Times() *Statement {
	s := Times()
	g.items = append(g.items, s)
	return s
}

// Times renders the provided operator / token.
func (s *Statement) Times() *Statement {
	t := token{
		typ:     operatorToken,
		content: "*",
	}
	*s = append(*s, t)
	return s
}

// GreaterThan renders the provided operator / token.
func GreaterThan() *Statement {
	return newStatement().GreaterThan()
}

// GreaterThan renders the provided operator / token.
func (g *Group) GreaterThan() *Statement {
	s := GreaterThan()
	g.items = append(g.items, s)
	return s
}

// GreaterThan renders the provided operator / token.
func (s *Statement) GreaterThan() *Statement {
	t := token{
		typ:     operatorToken,
		content: ">",
	}
	*s = append(*s, t)
	return s
}

// GreaterThanOrEqual renders the provided operator / token.
func GreaterThanOrEqual() *Statement {
	return newStatement().GreaterThanOrEqual()
}

// GreaterThanOrEqual renders the provided operator / token.
func (g *Group) GreaterThanOrEqual() *Statement {
	s := GreaterThanOrEqual()
	g.items = append(g.items, s)
	return s
}

// GreaterThanOrEqual renders the provided operator / token.
func (s *Statement) GreaterThanOrEqual() *Statement {
	t := token{
		typ:     operatorToken,
		content: ">=",
	}
	*s = append(*s, t)
	return s
}

// LessThan renders the provided operator / token.
func LessThan() *Statement {
	return newStatement().LessThan()
}

// LessThan renders the provided operator / token.
func (g *Group) LessThan() *Statement {
	s := LessThan()
	g.items = append(g.items, s)
	return s
}

// LessThan renders the provided operator / token.
func (s *Statement) LessThan() *Statement {
	t := token{
		typ:     operatorToken,
		content: "<",
	}
	*s = append(*s, t)
	return s
}

// LessThanOrEqual renders the provided operator / token.
func LessThanOrEqual() *Statement {
	return newStatement().LessThanOrEqual()
}

// LessThanOrEqual renders the provided operator / token.
func (g *Group) LessThanOrEqual() *Statement {
	s := LessThanOrEqual()
	g.items = append(g.items, s)
	return s
}

// LessThanOrEqual renders the provided operator / token.
func (s *Statement) LessThanOrEqual() *Statement {
	t := token{
		typ:     operatorToken,
		content: "<=",
	}
	*s = append(*s, t)
	return s
}

// MapAssign renders the provided operator / token.
func MapAssign() *Statement {
	return newStatement().MapAssign()
}

// MapAssign renders the provided operator / token.
func (g *Group) MapAssign() *Statement {
	s := MapAssign()
	g.items = append(g.items, s)
	return s
}

// MapAssign renders the provided operator / token.
func (s *Statement) MapAssign() *Statement {
	t := token{
		typ:     operatorToken,
		content: ":",
	}
	*s = append(*s, t)
	return s
}

// Equals renders the provided operator / token.
func Equals() *Statement {
	return newStatement().Equals()
}

// Equals renders the provided operator / token.
func (g *Group) Equals() *Statement {
	s := Equals()
	g.items = append(g.items, s)
	return s
}

// Equals renders the provided operator / token.
func (s *Statement) Equals() *Statement {
	t := token{
		typ:     operatorToken,
		content: "=",
	}
	*s = append(*s, t)
	return s
}

// DoubleEquals renders the provided operator / token.
func DoubleEquals() *Statement {
	return newStatement().DoubleEquals()
}

// DoubleEquals renders the provided operator / token.
func (g *Group) DoubleEquals() *Statement {
	s := DoubleEquals()
	g.items = append(g.items, s)
	return s
}

// DoubleEquals renders the provided operator / token.
func (s *Statement) DoubleEquals() *Statement {
	t := token{
		typ:     operatorToken,
		content: "==",
	}
	*s = append(*s, t)
	return s
}

// DoesNotEqual renders the provided operator / token.
func DoesNotEqual() *Statement {
	return newStatement().DoesNotEqual()
}

// DoesNotEqual renders the provided operator / token.
func (g *Group) DoesNotEqual() *Statement {
	s := DoesNotEqual()
	g.items = append(g.items, s)
	return s
}

// DoesNotEqual renders the provided operator / token.
func (s *Statement) DoesNotEqual() *Statement {
	t := token{
		typ:     operatorToken,
		content: "!=",
	}
	*s = append(*s, t)
	return s
}

// Opln renders the provided operator / token.
func Opln(op string) *Statement {
	return newStatement().Opln(op)
}

// Opln renders the provided operator / token.
func (g *Group) Opln(op string) *Statement {
	s := Opln(op)
	g.items = append(g.items, s)
	return s
}

// Opln renders the provided operator / token.
func (s *Statement) Opln(op string) *Statement {
	t := token{
		typ:     operatorToken,
		content: fmt.Sprintf("%s\n", op),
	}
	*s = append(*s, t)
	return s
}

// Assign renders the colon equals ( `:=` ) operator
func Assign() *Statement {
	return newStatement().Op(":=")
}

// Assign renders the colon equals ( `:=` ) operator
func (g *Group) Assign() *Statement {
	s := Op(":=")
	g.items = append(g.items, s)
	return s
}

// Assign renders the colon equals ( `:=` ) operator
func (s *Statement) Assign() *Statement {
	t := token{
		typ:     operatorToken,
		content: ":=",
	}
	*s = append(*s, t)
	return s
}

// Dot renders a period followed by an identifier. Use for fields and selectors.
func Dot(name string) *Statement {
	// notest
	// don't think this can be used in valid code?
	return newStatement().Dot(name)
}

// Dot renders a period followed by an identifier. Use for fields and selectors.
func (g *Group) Dot(name string) *Statement {
	// notest
	// don't think this can be used in valid code?
	s := Dot(name)
	g.items = append(g.items, s)
	return s
}

// Dot renders a period followed by an identifier. Use for fields and selectors.
func (s *Statement) Dot(name string) *Statement {
	d := token{
		typ:     delimiterToken,
		content: ".",
	}
	t := token{
		typ:     identifierToken,
		content: name,
	}
	*s = append(*s, d, t)
	return s
}

// Dotln renders a period followed by an identifier. Use for fields and selectors.
func Dotln(name string) *Statement {
	// notest
	// don't think this can be used in valid code?
	return newStatement().Dotln(name)
}

// Dotln renders a period followed by an identifier. Use for fields and selectors.
func (g *Group) Dotln(name string) *Statement {
	// notest
	// don't think this can be used in valid code?
	s := Dotln(name)
	g.items = append(g.items, s)
	return s
}

// Dotln renders a period followed by an identifier. Use for fields and selectors.
func (s *Statement) Dotln(name string) *Statement {
	d := token{
		typ:     delimiterToken,
		content: ".\n\t",
	}
	t := token{
		typ:     identifierToken,
		content: name,
	}
	*s = append(*s, d, t)
	return s
}

// Dotf renders a period followed by an identifier. Use for fields and selectors.
func Dotf(name string, args ...interface{}) *Statement {
	// notest
	// don't think this can be used in valid code?
	return newStatement().Dotf(name, args...)
}

// Dotf renders a period followed by an identifier. Use for fields and selectors.
func (g *Group) Dotf(name string, args ...interface{}) *Statement {
	// notest
	// don't think this can be used in valid code?
	s := Dotf(name, args...)
	g.items = append(g.items, s)
	return s
}

// Dotf renders a period followed by an identifier. Use for fields and selectors.
func (s *Statement) Dotf(name string, args ...interface{}) *Statement {
	d := token{
		typ:     delimiterToken,
		content: ".",
	}
	t := token{
		typ:     identifierToken,
		content: fmt.Sprintf(name, args...),
	}
	*s = append(*s, d, t)
	return s
}

// Underscore renders an identifier.
func Underscore() *Statement {
	return newStatement().Underscore()
}

// Underscore renders an identifier.
func (g *Group) Underscore() *Statement {
	s := Underscore()
	g.items = append(g.items, s)
	return s
}

// Underscore renders an identifier.
func (s *Statement) Underscore() *Statement {
	t := token{
		typ:     identifierToken,
		content: "_",
	}
	*s = append(*s, t)
	return s
}

// Or renders an identifier.
func Or() *Statement {
	return newStatement().Or()
}

// Or renders an identifier.
func (g *Group) Or() *Statement {
	s := Or()
	g.items = append(g.items, s)
	return s
}

// Or renders an identifier.
func (s *Statement) Or() *Statement {
	t := token{
		typ:     identifierToken,
		content: "||",
	}
	*s = append(*s, t)
	return s
}

// ID renders an identifier.
func ID(name string) *Statement {
	return newStatement().ID(name)
}

// ID renders an identifier.
func (g *Group) ID(name string) *Statement {
	s := ID(name)
	g.items = append(g.items, s)
	return s
}

// ID renders an identifier.
func (s *Statement) ID(name string) *Statement {
	t := token{
		typ:     identifierToken,
		content: name,
	}
	*s = append(*s, t)
	return s
}

// IDf renders an identifier.
func IDf(name string, args ...interface{}) *Statement {
	return newStatement().IDf(name, args...)
}

// IDf renders an identifier.
func (g *Group) IDf(name string, args ...interface{}) *Statement {
	s := IDf(name, args...)
	g.items = append(g.items, s)
	return s
}

//f ID renders an identifier.
func (s *Statement) IDf(name string, args ...interface{}) *Statement {
	t := token{
		typ:     identifierToken,
		content: fmt.Sprintf(name, args...),
	}
	*s = append(*s, t)
	return s
}

// Qual renders a qualified identifier. Imports are automatically added when
// used with a File. If the path matches the local path, the package name is
// omitted. If package names conflict they are automatically renamed. Note that
// it is not possible to reliably determine the package name given an arbitrary
// package path, so a sensible name is guessed from the path and added as an
// alias. The names of all standard library packages are known so these do not
// need to be aliased. If more control is needed of the aliases, see
// [File.ImportName](#importname) or [File.ImportAlias](#importalias).
func Qual(path, name string) *Statement {
	return newStatement().Qual(path, name)
}

// Qual renders a qualified identifier. Imports are automatically added when
// used with a File. If the path matches the local path, the package name is
// omitted. If package names conflict they are automatically renamed. Note that
// it is not possible to reliably determine the package name given an arbitrary
// package path, so a sensible name is guessed from the path and added as an
// alias. The names of all standard library packages are known so these do not
// need to be aliased. If more control is needed of the aliases, see
// [File.ImportName](#importname) or [File.ImportAlias](#importalias).
func (g *Group) Qual(path, name string) *Statement {
	s := Qual(path, name)
	g.items = append(g.items, s)
	return s
}

// Qual renders a qualified identifier. Imports are automatically added when
// used with a File. If the path matches the local path, the package name is
// omitted. If package names conflict they are automatically renamed. Note that
// it is not possible to reliably determine the package name given an arbitrary
// package path, so a sensible name is guessed from the path and added as an
// alias. The names of all standard library packages are known so these do not
// need to be aliased. If more control is needed of the aliases, see
// [File.ImportName](#importname) or [File.ImportAlias](#importalias).
func (s *Statement) Qual(path, name string) *Statement {
	g := &Group{
		close: " ",
		items: []Code{
			token{
				typ:     packageToken,
				content: path,
			},
			token{
				typ:     identifierToken,
				content: name,
			},
		},
		name:      "qual",
		open:      " ",
		separator: ".",
	}
	*s = append(*s, g)
	return s
}

// Line inserts a blank line.
func Line() *Statement {
	return newStatement().Line()
}

// Line inserts a blank line.
func (g *Group) Line() *Statement {
	s := Line()
	g.items = append(g.items, s)
	return s
}

// Line inserts a blank line.
func (s *Statement) Line() *Statement {
	t := token{
		typ:     layoutToken,
		content: "\n",
	}
	*s = append(*s, t)
	return s
}

func IsLine(t token) bool {
	if x, ok := t.content.(string); ok {
		return x == "\n" && t.typ == layoutToken
	}
	return false
}
