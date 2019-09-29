package gen

import (
	"go/format"
	"regexp"
)

// EXPERIMENTAL

func formatNulls(file []byte) []byte {
	reg := regexp.MustCompile("(jen[.]Null[(][)][.])")
	return reg.ReplaceAll(file, []byte("jen."))
}
func formatStructs(file []byte) []byte {
	reg2 := regexp.MustCompile(`Struct([(]).+[)]`)
	ret := reg2.ReplaceAllFunc(file, func(b []byte) []byte {
		reg := regexp.MustCompile(`(Struct[(])`)
		b2 := reg.ReplaceAll(b, []byte("$0\n"))
		reg2 := regexp.MustCompile(`([)])([)])$`)
		b3 := reg2.ReplaceAll(b2, []byte("$1,\n$2"))
		//reg3 := regexp.MustCompile(",")
		//commas := reg3.ReplaceAll(b3, []byte("$0\n"))
		return b3
	})

	return ret
}
func formatBlock(file []byte) []byte {
	reg2 := regexp.MustCompile(`Block([(]).+[)]`)
	ret := reg2.ReplaceAllFunc(file, func(b []byte) []byte {
		reg := regexp.MustCompile(`(Block[(])`)
		b2 := reg.ReplaceAll(b, []byte("$0\n"))
		reg2 := regexp.MustCompile(`([)])([)])$`)
		b3 := reg2.ReplaceAll(b2, []byte("$1,\n$2"))
		//reg3 := regexp.MustCompile(",")
		//commas := reg3.ReplaceAll(b3, []byte("$0\n"))
		return b3
	})

	return ret
}
func formatParams(file []byte) []byte {
	reg2 := regexp.MustCompile(`(:?Params)[(].+"[)]{2}.`)
	ret := reg2.ReplaceAllFunc(file, func(b []byte) []byte {
		reg := regexp.MustCompile(`(Params[(])`)
		b2 := reg.ReplaceAll(b, []byte("$0\n"))
		reg2 := regexp.MustCompile(`([)])([)])$`)
		b3 := reg2.ReplaceAll(b2, []byte("$1,\n$2"))
		//reg3 := regexp.MustCompile(",")
		//commas := reg3.ReplaceAll(b3, []byte("$0\n"))
		return b3
	})

	return ret
}
func goFormat(file []byte) ([]byte, error) {
	fmtBytes, err := format.Source(file)
	if err != nil {
		return file, err
	}
	return fmtBytes, nil
}
