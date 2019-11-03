package main

import (
	"log"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/templates"
)

func main() {
	log.Println("building example output...")
	todoProject := &models.Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/todopartdeux",
		Name:       wordsmith.FromSingularPascalCase("Todo"),
		DataTypes: []models.DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Eman"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Info"),
						Type:                  "string",
						Pointer:               false,
						DefaultAllowed:        true,
						DefaultValue:          "''",
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
		},
	}

	if err := project.RenderProject(todoProject); err != nil {
		log.Fatal(err)
	}
}
