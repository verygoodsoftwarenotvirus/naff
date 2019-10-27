package randmodel

import (
	"github.com/icrowley/fake"
	models "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

// RandomWebhookInput creates a random WebhookCreationInput
func RandomWebhookInput() *models.WebhookCreationInput {
	x := &models.WebhookCreationInput{
		Name:        fake.Word(),
		URL:         fake.DomainName(),
		ContentType: "application/json",
		Method:      "POST",
	}
	return x
}
