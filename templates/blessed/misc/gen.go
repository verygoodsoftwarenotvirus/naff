package misc

import (
	// "encoding/json"
	// "io/ioutil"
	// "log"

	// "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, types []models.DataType) error {
	return nil

	// miscFiles := map[string]string{
	// 	"": "",
	// }

	// for filename, file := range miscFiles {
	// 	fn := utils.BuildTemplatePath(filename)

	// 	f, _ := json.MarshalIndent(file, "", " ")
	// 	if err := ioutil.WriteFile(fn, f, 0644); err != nil {
	// 		log.Printf("error rendering %q: %v\n", filename, err)
	// 	}
	// }

	// return nil
}
