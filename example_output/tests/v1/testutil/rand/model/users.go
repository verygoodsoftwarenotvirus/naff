package randmodel

import (
	"time"

	"github.com/icrowley/fake"
	"github.com/pquerna/otp/totp"
	models "gitlab.com/verygoodsoftwarenotvirus/todo/models/v1"
)

func init() {
	fake.Seed(time.Now().UnixNano())
}

func mustBuildCode(totpSecret string) string {
	code, err := totp.GenerateCode(totpSecret, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return code
}

// RandomUserInput creates a random UserInput
func RandomUserInput() *models.UserInput {
	// I had difficulty ensuring these values were unique, even when fake.Seed was called. Could've been fake's fault,
	// could've been docker's fault. In either case, it wasn't worth the time to investigate and determine the culprit.
	username := fake.UserName() + fake.HexColor() + fake.Country()
	x := &models.UserInput{
		Username: username,
		Password: fake.Password(64, 128, true, true, true),
	}
	return x
}
