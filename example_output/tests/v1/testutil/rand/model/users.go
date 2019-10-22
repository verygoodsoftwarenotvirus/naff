package randmodel

import (
	"time"

	"github.com/icrowley/fake"
	"github.com/pquerna/otp/totp"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
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
	username := fake.UserName() + fake.HexColor() + fake.Country()
	x := &models.UserInput{
		Username: username,
		Password: fake.Password(64, 128, true, true, true),
	}
	return x
}
