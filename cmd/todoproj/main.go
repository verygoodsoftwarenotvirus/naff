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
			},
		},
	}

	// taskProject := &models.Project{
	// 	OutputPath: "gitlab.com/verygoodsoftwarenotvirus/todopartdeux",
	// 	Name:       wordsmith.FromSingularPascalCase("Task"),
	// 	DataTypes: []models.DataType{
	// 		{
	// 			Name: wordsmith.FromSingularPascalCase("Entry"),
	// 			Fields: []models.DataField{
	// 				{
	// 					Name:                  wordsmith.FromSingularPascalCase("Eman"),
	// 					Type:                  "string",
	// 					Pointer:               false,
	// 					ValidForCreationInput: true,
	// 					ValidForUpdateInput:   true,
	// 				},
	// 				{
	// 					Name:                  wordsmith.FromSingularPascalCase("Deets"),
	// 					Type:                  "string",
	// 					Pointer:               false,
	// 					DefaultAllowed:        true,
	// 					DefaultValue:          "''",
	// 					ValidForCreationInput: true,
	// 					ValidForUpdateInput:   true,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	//////////////////////////////////////////////////////////////////

	// allFields := []models.DataField{}
	// for _, typ := range []string{
	// 	"bool",
	// 	"string",
	// 	"float32",
	// 	"float64",
	// 	"uint",
	// 	"uint8",
	// 	"uint16",
	// 	"uint32",
	// 	"uint64",
	// 	"int",
	// 	"int8",
	// 	"int16",
	// 	"int32",
	// 	"int64",
	// } {
	// 	allFields = append(
	// 		allFields,
	// 		models.DataField{
	// 			Name:                  wordsmith.FromSingularPascalCase(fmt.Sprintf("field_%s", strings.ToTitle(typ))),
	// 			Type:                  typ,
	// 			Pointer:               false,
	// 			ValidForCreationInput: true,
	// 			ValidForUpdateInput:   true,
	// 		},
	// 		models.DataField{
	// 			Name:                  wordsmith.FromSingularPascalCase(fmt.Sprintf("field_%s", strings.ToTitle(typ)+"AsPointer")),
	// 			Type:                  typ,
	// 			Pointer:               true,
	// 			ValidForCreationInput: true,
	// 			ValidForUpdateInput:   true,
	// 		},
	// 	)
	// }

	// typeProject := &models.Project{
	// 	OutputPath: "gitlab.com/verygoodsoftwarenotvirus/todopartdeux",
	// 	Name:       wordsmith.FromSingularPascalCase("Task"),
	// 	DataTypes: []models.DataType{
	// 		{
	// 			Name:   wordsmith.FromSingularPascalCase("Entry"),
	// 			Fields: allFields,
	// 		},
	// 	},
	// }

	// clear; rm -rf ../todopartdeux; make example_output
	if err := project.RenderProject(todoProject); err != nil {
		log.Fatal(err)
	}
}
