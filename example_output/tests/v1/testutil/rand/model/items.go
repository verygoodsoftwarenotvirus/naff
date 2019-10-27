package randmodel

import (
	"github.com/icrowley/fake"
	models "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

// RandomItemCreationInput creates a random ItemInput
func RandomItemCreationInput() *models.ItemCreationInput {
	x := &models.ItemCreationInput{
		Name:    fake.Word(),
		Details: fake.Word(),
	}

	return x
}
