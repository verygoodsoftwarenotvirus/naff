package client

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
// 	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
// )

// func Test_buildPathForType(T *testing.T) {

// 	T.Parallel()

// 	T.Run("happy path", func(t *testing.T) {
// 		apple := models.DataType{
// 			Name: wordsmith.FromSingularPascalCase("Apple"),
// 		}
// 		banana := models.DataType{
// 			Name:            wordsmith.FromSingularPascalCase("Banana"),
// 			BelongsToStruct: apple.Name,
// 		}
// 		cherry := models.DataType{
// 			Name:            wordsmith.FromSingularPascalCase("Cherry"),
// 			BelongsToStruct: banana.Name,
// 		}

// 		proj := &models.Project{
// 			DataTypes: []models.DataType{apple, banana, cherry},
// 		}

// 		expected := ""
// 		actual := buildPathForType(proj, cherry)

// 		assert.Equal(t, expected, actual)
// 	})
// }
