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
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
		Name:       wordsmith.FromSingularPascalCase("Todo"),
		DataTypes: []models.DataType{
			{
				Name: wordsmith.FromSingularPascalCase("Item"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
					{
						Name:                  wordsmith.FromSingularPascalCase("Details"),
						Type:                  "string",
						Pointer:               false,
						DefaultAllowed:        true,
						DefaultValue:          "''",
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToUser: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Forum"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToUser: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Thread"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Title"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
			},
			{
				Name: wordsmith.FromSingularPascalCase("Comment"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Content"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct: wordsmith.FromSingularPascalCase("Thread"),
			},
			{
				Name: wordsmith.FromSingularPascalCase("Tag"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Key"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToNobody: true,
			},
		},
	}

	// clear; rm -rf ../todopartdeux; make example_output
	if err := project.RenderProject(todoProject); err != nil {
		log.Fatal(err)
	}
}
