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
			//{
			//	Name: wordsmith.FromSingularPascalCase("Item"),
			//	Fields: []models.DataField{
			//		{
			//			Name:                  wordsmith.FromSingularPascalCase("Name"),
			//			Type:                  "string",
			//			Pointer:               false,
			//			ValidForCreationInput: true,
			//			ValidForUpdateInput:   true,
			//		},
			//		{
			//			Name:                  wordsmith.FromSingularPascalCase("Details"),
			//			Type:                  "string",
			//			Pointer:               false,
			//			DefaultAllowed:        true,
			//			DefaultValue:          "''",
			//			ValidForCreationInput: true,
			//			ValidForUpdateInput:   true,
			//		},
			//	},
			//	BelongsToUser: true,
			//},
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
				BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
				BelongsToNobody: true,
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
				BelongsToStruct: wordsmith.FromSingularPascalCase("Subforum"),
				BelongsToUser:   true,
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
				BelongsToStruct: wordsmith.FromSingularPascalCase("Thread"),
				BelongsToUser:   true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("PostRating"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Rating"),
						Type:                  "uint8",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToStruct: wordsmith.FromSingularPascalCase("Post"),
				BelongsToUser:   true,
			},
			{
				Name: wordsmith.FromSingularPascalCase("Signature"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Text"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
				BelongsToUser: true,
			},
		},
	}

	if err := project.RenderProject(todoProject); err != nil {
		log.Fatal(err)
	}
}
