package misc

import (
	"log"
	"os"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(pkgRoot string, projectName wordsmith.SuperPalabra, types []models.DataType) error {
	composeFiles := map[string]func() []byte{
		"development/badges.json": badgesDotJSON,
		".dockerignore":           dockerIgnore,
		".gitignore":              gitIgnore,
	}

	for filename, file := range composeFiles {
		fname := utils.BuildTemplatePath(filename)

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("error opening file: %v", err)
			return err
		}

		bytes := file()
		if _, err := f.Write(bytes); err != nil {
			log.Printf("error writing to file: %v", err)
			return err
		}
	}

	return nil
}

func badgesDotJSON() []byte {
	return []byte(`{
    "badges": [
        {
            "name": "godoc",
            "gitlab": {
                "link": "https://godoc.org/gitlab.com/%{project_path}",
                "badge_image_url": "https://godoc.org/gitlab.com/%{project_path}?status.svg"
            }
        },
        {
            "name": "ci",
            "gitlab": {
                "link": "https://gitlab.com/%{project_path}/commits/%{default_branch}",
                "badge_image_url": "https://gitlab.com/%{project_path}/badges/%{default_branch}/pipeline.svg"
            }
        },
        {
            "name": "coverage",
            "gitlab": {
                "link": "https://gitlab.com/%{project_path}",
                "badge_image_url": "https://gitlab.com/%{project_path}/badges/%{default_branch}/coverage.svg"
            }
        },
        {
            "name": "docker",
            "gitlab": {
                "link": "https://hub.docker.com/r/%{project_path}",
                "badge_image_url": "https://img.shields.io/docker/automated/%{project_path}.svg"
            }
        }
    ]
}`)
}

func dockerIgnore() []byte {
	return []byte(`**/node_modules
**/dist
`)
}

func gitIgnore() []byte {
	return []byte(`# Binaries for programs and plugins
*.exe
*.dll
*.so
*.dylib

# Test binary, build with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# pretty good IDEs
.idea
.vscode/

# Vim
*.swp

# Sqlite databases
*.db

# OSX
.DS_Store

# # Go
# vendor

# Python
.env
.mypy_cache
__pycache__
artifacts

# frontend things
node_modules

# Log files
npm-debug.log*
yarn-debug.log*
yarn-error.log*

frontend/v1/public/bundle.*

*.coverprofile
*.profile`)
}
