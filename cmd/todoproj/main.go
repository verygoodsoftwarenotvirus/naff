package main

import (
	"log"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/templates"
)

var (
	projects = map[string]*models.Project{
		"todo": {
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
							DefaultAllowed:        true,
							DefaultValue:          "''",
							Pointer:               false,
							ValidForCreationInput: true,
							ValidForUpdateInput:   true,
						},
					},
					BelongsToUser: true,
				},
			},
		},
		"discussion": {
			OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
			Name:       wordsmith.FromSingularPascalCase("Discussion"),
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
					Name: wordsmith.FromSingularPascalCase("Notification"),
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
		},
	}
)

func main() {
	const chosenProject = "todo"

	if err := project.RenderProject(projects[chosenProject]); err != nil {
		log.Fatal(err)
	}
}
