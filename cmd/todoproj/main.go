package main

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	project "gitlab.com/verygoodsoftwarenotvirus/naff/templates"
	"log"
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

	if err := project.RenderProject(todoProject); err != nil {
		log.Fatal(err)
	}
}
