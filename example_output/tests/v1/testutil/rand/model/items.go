package randmodel

import (
	"github.com/icrowley/fake"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// RandomItemCreationInput creates a random ItemInput
func RandomItemCreationInput() *models.ItemCreationInput {
	x := &models.ItemCreationInput{
		Name:    fake.Word(),
		Details: fake.Sentence(),
	}
	return x
}
