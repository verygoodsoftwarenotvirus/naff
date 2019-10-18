package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	sourcePath := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http")

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			fp := filepath.Join(sourcePath, file.Name())
			toCompare := strings.Replace(fp, "verygoodsoftwarenotvirus/todo", "verygoodsoftwarenotvirus/naff/example_output", 1)

			ogB, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			newB, err := ioutil.ReadFile(toCompare)
			if err != nil {
				log.Fatal(err)
			}

			ogSum := fmt.Sprintf("%x", sha256.Sum256(ogB))
			newSum := fmt.Sprintf("%x", sha256.Sum256(newB))

			if ogSum != newSum {
				fmt.Printf("diff --side-by-side --brief %q %q\n", fp, toCompare)
			}
		}
	}
}
