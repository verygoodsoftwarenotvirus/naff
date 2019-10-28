package randmodel

import (
	models "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"

	"github.com/icrowley/fake"
)

// RandomItemCreationInput creates a random ItemInput
func RandomItemCreationInput() *models.ItemCreationInput {
	x := &models.ItemCreationInput{
		Name:    fake.Word(),
		Details: fake.Word(),
	}

	return x
}
