package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/templates"
)

const (
	projectDiscussion = "discussion"
	projectTodo       = "todo"
	projectGamut      = "gamut"
)

var (
	projects = map[string]*models.Project{
		projectTodo:       testprojects.TodoApp,
		projectDiscussion: testprojects.ForumsApp,
		projectGamut:      testprojects.EveryTypeApp,
	}
)

func init() {
	projects[projectGamut].EnableDatabase(models.Postgres)

	projects[projectDiscussion].EnableDatabase(models.Postgres)

	projects[projectTodo].EnableDatabase(models.Postgres)
	projects[projectTodo].EnableDatabase(models.Sqlite)
	projects[projectTodo].EnableDatabase(models.MariaDB)
}

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

		if err := project.RenderProject(chosenProject); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("no project selected")
	}
}
