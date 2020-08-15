package main

import (
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/templates"
)

const (
	projectForums = "forums"
	projectTodo   = "todo"
	projectGamut  = "every_type"
)

var (
	projects = map[string]*models.Project{
		projectTodo:   testprojects.BuildTodoApp(),
		projectForums: testprojects.BuildForumsApp(),
		projectGamut:  testprojects.BuildEveryTypeApp(),
	}
)

const (
	this = "gitlab.com/verygoodsoftwarenotvirus/naff"
)

func main() {
	if err := os.RemoveAll(filepath.Join(os.Getenv("GOPATH"), "src", this, "example_output")); err != nil {
		log.Printf("error removing output dir: %v", err)
	}

	if chosenProjectKey := os.Getenv("PROJECT"); chosenProjectKey != "" {
		chosenProject := projects[chosenProjectKey]

		if outputDir := os.Getenv("OUTPUT_DIR"); outputDir != "" {
			chosenProject.OutputPath = filepath.Join(this, outputDir)
		}

		project.RenderProject(chosenProject)
	} else {
		log.Fatal("no project selected")
	}
}
