package main

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/new_templates"
)

func main() {
	todoProject := &models.Project{
		Name: models.Name{},
		DataTypes: []models.DataType{
			{
				Name: models.Name{
					Singular:                "Item",
					Plural:                  "Items",
					RouteName:               "items",
					PluralRouteName:         "item",
					UnexportedVarName:       "item",
					PluralUnexportedVarName: "items",
				},
			},
		},
	}
	project.RenderProject(todoProject)
}
