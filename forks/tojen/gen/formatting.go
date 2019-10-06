package gen

import (
	"go/format"
)

func goFormat(file []byte) ([]byte, error) {
	fmtBytes, err := format.Source(file)
	if err != nil {
		return file, err
	}
	return fmtBytes, nil
}
