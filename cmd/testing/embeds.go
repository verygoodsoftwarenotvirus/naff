package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	assets "gitlab.com/verygoodsoftwarenotvirus/naff/embedded"
)

const (
	defaultFileExtension  = ".tmpl"
	baseTemplateDirectory = "template/base_repository/"
	iterableDirectory     = "template/iterables/"
)

func main() {
	thisPackage := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/naff")

	files, err := assets.WalkDirs("template/base_repository", false)
	if err != nil {
		panic(err)
	}

	for _, path := range files {

		renderPath := strings.Replace(
			path,
			filepath.Join(thisPackage, baseTemplateDirectory),
			filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/slef"),
			1,
		)
		renderPath = strings.ReplaceAll(renderPath, ".tmpl", "")

		f, err := assets.FS.OpenFile(context.Background(), path, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		s := string(bs)
		_ = s
		println()
	}

	fmt.Println(files)
}
