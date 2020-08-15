package composefiles

import (
	"fmt"
	"os"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	for _, folder := range []string{
		utils.BuildTemplatePath(proj.OutputPath, "environments/local/config_files"),
		utils.BuildTemplatePath(proj.OutputPath, "environments/testing/config_files"),
	} {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			if mkdirErr := os.MkdirAll(folder, os.ModePerm); mkdirErr != nil {
				return fmt.Errorf("error making necessary directory: %w", mkdirErr)
			}
		} else if err != nil {
			return fmt.Errorf("error creating folders: %w", err)
		}
	}

	return nil
}
