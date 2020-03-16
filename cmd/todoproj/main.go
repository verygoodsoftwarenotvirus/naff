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
				BelongsToNobody: true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Subforum"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct:      wordsmith.FromSingularPascalCase("Forum"),
				ReadRestrictedToUser: false,
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
				ReadRestrictedToUser: false,
				BelongsToStruct:      wordsmith.FromSingularPascalCase("Subforum"),
				BelongsToUser:        true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Post"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Content"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				ReadRestrictedToUser: false,
				BelongsToStruct:      wordsmith.FromSingularPascalCase("Thread"),
				BelongsToUser:        true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Message"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Text"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				ReadRestrictedToUser: true,
			},
		},
	}

	if err := project.RenderProject(todoProject); err != nil {
		log.Fatal(err)
	}
}
