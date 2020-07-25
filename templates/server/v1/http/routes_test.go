package httpserver

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildIterableAPIRoutesBlock(T *testing.T) {
	T.Parallel()

	T.Run("normal operation", func(t *testing.T) {
		apple := models.DataType{
			Name: wordsmith.FromSingularPascalCase("Apple"),
		}
		banana := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Banana"),
			BelongsToStruct: apple.Name,
		}
		cherry := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Cherry"),
			BelongsToStruct: banana.Name,
		}

		proj := &models.Project{
			DataTypes: []models.DataType{apple, banana, cherry},
		}

		code := jen.NewFile("farts")

		code.Add(
			jen.Func().ID("doSomething").Params().Block(
				buildIterableAPIRoutes(proj),
			),
		)

		var b bytes.Buffer
		err := code.Render(&b)
		require.NoError(t, err)

		expected := `package farts

import (
	"fmt"
	chi "github.com/go-chi/chi"
	apples "services/v1/apples"
	bananas "services/v1/bananas"
	cherries "services/v1/cherries"
)

func doSomething() {
	// Apples
	v1Router.Route("/apples", func(applesRouter chi.Router) {
		appleRoute := fmt.Sprintf(numericIDPattern, apples.URIParamKey)
		applesRouter.With(s.applesService.CreationInputMiddleware).Post("/", s.applesService.CreateHandler())
		applesRouter.Get(appleRoute, s.applesService.ReadHandler())
		applesRouter.With(s.applesService.UpdateInputMiddleware).Put(appleRoute, s.applesService.UpdateHandler())
		applesRouter.Delete(appleRoute, s.applesService.ArchiveHandler())
		applesRouter.Get("/", s.applesService.ListHandler())

		// Bananas

		appleRouter.Route("/bananas", func(bananasRouter chi.Router) {
			bananaRoute := fmt.Sprintf(numericIDPattern, bananas.URIParamKey)
			bananasRouter.With(s.bananasService.CreationInputMiddleware).Post("/", s.bananasService.CreateHandler())
			bananasRouter.Get(bananaRoute, s.bananasService.ReadHandler())
			bananasRouter.With(s.bananasService.UpdateInputMiddleware).Put(bananaRoute, s.bananasService.UpdateHandler())
			bananasRouter.Delete(bananaRoute, s.bananasService.ArchiveHandler())
			bananasRouter.Get("/", s.bananasService.ListHandler())

			// Cherries

			bananaRouter.Route("/cherries", func(cherriesRouter chi.Router) {
				cherryRoute := fmt.Sprintf(numericIDPattern, cherries.URIParamKey)
				cherriesRouter.With(s.cherriesService.CreationInputMiddleware).Post("/", s.cherriesService.CreateHandler())
				cherriesRouter.Get(cherryRoute, s.cherriesService.ReadHandler())
				cherriesRouter.With(s.cherriesService.UpdateInputMiddleware).Put(cherryRoute, s.cherriesService.UpdateHandler())
				cherriesRouter.Delete(cherryRoute, s.cherriesService.ArchiveHandler())
				cherriesRouter.Get("/", s.cherriesService.ListHandler())
			})

		})

	})

}
`
		actual := b.String()

		assert.Equal(t, expected, actual)
	})
}
