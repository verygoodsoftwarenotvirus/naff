package testprojects

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

var (
	everyType = models.DataType{
		Name:            wordsmith.FromSingularPascalCase("EveryType"),
		BelongsToNobody: true,
		Fields: []models.DataField{
			{
				Name:                  wordsmith.FromSingularPascalCase("String"),
				Type:                  "string",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToString"),
				Type:                  "string",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Bool"),
				Type:                  "bool",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToBool"),
				Type:                  "bool",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Int"),
				Type:                  "int",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToInt"),
				Type:                  "int",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Int8"),
				Type:                  "int8",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToInt8"),
				Type:                  "int8",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Int16"),
				Type:                  "int16",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToInt16"),
				Type:                  "int16",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Int32"),
				Type:                  "int32",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToInt32"),
				Type:                  "int32",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Int64"),
				Type:                  "int64",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToInt64"),
				Type:                  "int64",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Uint"),
				Type:                  "uint",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToUint"),
				Type:                  "uint",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Uint8"),
				Type:                  "uint8",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToUint8"),
				Type:                  "uint8",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Uint16"),
				Type:                  "uint16",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToUint16"),
				Type:                  "uint16",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Uint32"),
				Type:                  "uint32",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToUint32"),
				Type:                  "uint32",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Uint64"),
				Type:                  "uint64",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToUint64"),
				Type:                  "uint64",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Float32"),
				Type:                  "float32",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToFloat32"),
				Type:                  "float32",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("Float64"),
				Type:                  "float64",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               false,
			},
			{
				Name:                  wordsmith.FromSingularPascalCase("PointerToFloat64"),
				Type:                  "float64",
				ValidForCreationInput: true,
				ValidForUpdateInput:   true,
				Pointer:               true,
			},
		},
	}

	forumDataTypes = []models.DataType{
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
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Subforum"),
			BelongsToUser:    true,
			RestrictedToUser: false,
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
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Thread"),
			BelongsToUser:    true,
			RestrictedToUser: false,
		},
		{
			Name: wordsmith.FromSingularPascalCase("ReactionIcon"),
			Fields: []models.DataField{
				{
					Name:                  wordsmith.FromSingularPascalCase("Icon"),
					Type:                  "string",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToNobody: true,
		},
		{
			Name: wordsmith.FromSingularPascalCase("PostReaction"),
			Fields: []models.DataField{
				{
					Name:                  nil,
					Type:                  "uint64",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
					BelongsToEnumeration:  wordsmith.FromSingularPascalCase("ReactionIcon"),
				},
			},
			BelongsToStruct:  wordsmith.FromSingularPascalCase("Post"),
			BelongsToUser:    true,
			RestrictedToUser: false,
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
			BelongsToUser:    true,
			RestrictedToUser: true,
		},
	}

	todoDataTypes = []models.DataType{
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
					DefaultValue:          "''",
					Pointer:               false,
					ValidForCreationInput: true,
					ValidForUpdateInput:   true,
				},
			},
			BelongsToUser:    true,
			RestrictedToUser: true,
			SearchEnabled:    true,
		},
	}

	TodoApp = &models.Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
		Name:       wordsmith.FromSingularPascalCase("Todo"),
		DataTypes:  todoDataTypes,
	}

	ForumsApp = &models.Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
		Name:       wordsmith.FromSingularPascalCase("Discussion"),
		DataTypes:  forumDataTypes,
	}

	EveryTypeApp = &models.Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
		Name:       wordsmith.FromSingularPascalCase("Gamut"),
		//DataTypes:  append(forumDataTypes, everyType),
		DataTypes: append(append(forumDataTypes, everyType),
			models.DataType{
				Name:             wordsmith.FromSingularPascalCase("Recipe"),
				BelongsToUser:    true,
				RestrictedToUser: true,
				BelongsToStruct:  nil,
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Title"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
			models.DataType{
				Name:             wordsmith.FromSingularPascalCase("RecipeStep"),
				BelongsToUser:    false,
				RestrictedToUser: false,
				BelongsToStruct:  wordsmith.FromSingularPascalCase("Recipe"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Action"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
			models.DataType{
				Name:             wordsmith.FromSingularPascalCase("RecipeStepIngredient"),
				BelongsToUser:    false,
				RestrictedToUser: false,
				BelongsToStruct:  wordsmith.FromSingularPascalCase("RecipeStep"),
				Fields: []models.DataField{
					{
						Name:                  wordsmith.FromSingularPascalCase("Name"),
						Type:                  "string",
						Pointer:               false,
						ValidForCreationInput: true,
						ValidForUpdateInput:   true,
					},
				},
			},
		),
	}
)

func BuildTodoApp() *models.Project {
	p := &models.Project{
		OutputPath: "gitlab.com/verygoodsoftwarenotvirus/naff/example_output",
		Name:       wordsmith.FromSingularPascalCase("Todo"),
		DataTypes:  todoDataTypes,
	}

	p.EnableDatabase(models.Postgres)
	p.EnableDatabase(models.Sqlite)
	p.EnableDatabase(models.MariaDB)

	return p
}
