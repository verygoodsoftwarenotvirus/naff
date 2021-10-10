package data_scaffolder

import (
	"bytes"
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	packageName = "data_scaffolder"

	basePackagePath = "cmd/tools/data_scaffolder"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	stringFiles := map[string]string{
		"main.go":        mainDotGoString(proj),
		"exiter.go":      exiterDotGoString(proj),
		"exiter_test.go": exiterTestDotGoString(proj),
	}

	for path, file := range stringFiles {
		if err := utils.RenderStringFile(proj, filepath.Join(basePackagePath, path), file, true); err != nil {
			return err
		}
	}

	return nil
}

//go:embed main.gotpl
var mainTemplate string

func mainDotGoString(proj *models.Project) string {
	typeInitializers := buildInitializers(proj)

	var b bytes.Buffer
	if err := jen.Null().Add(typeInitializers...).RenderWithoutFormatting(&b); err != nil {
		panic(err)
	}

	generated := map[string]string{
		"typeInitializers": b.String(),
	}

	return models.RenderCodeFile(proj, mainTemplate, generated)
}

//go:embed exiter.gotpl
var exiterTemplate string

func exiterDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, exiterTemplate, nil)
}

//go:embed exiter_test.gotpl
var exiterTestTemplate string

func exiterTestDotGoString(proj *models.Project) string {
	return models.RenderCodeFile(proj, exiterTestTemplate, nil)
}

func buildInitializers(proj *models.Project) []jen.Code {
	initializers := []jen.Code{}
	/*
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for j := 0; j < int(dataCount); j++ {
				iterationLogger := userLogger.WithValue("creating", "items").WithValue("iteration", j)

				// create item
				createdItem, itemCreationErr := userClient.CreateItem(ctx, fakes.BuildFakeItemCreationInput())
				if itemCreationErr != nil {
					quitter.ComplainAndQuit(fmt.Errorf("creating item #%d: %w", j, itemCreationErr))
				}

				iterationLogger.WithValue(keys.WebhookIDKey, createdItem.ID).Debug("created item")
			}
			wg.Done()
		}(wg)
	*/

	for _, typ := range proj.DataTypes {
		if len(proj.FindOwnerTypeChain(typ)) == 0 {
			initializers = append(initializers, jen.Newline(), jen.Newline(),
				jen.ID("wg").Dot("Add").Call(jen.One()),
				jen.Newline(),
				jen.Go().Func().Params(jen.ID("wg").PointerTo().Qual("sync", "WaitGroup")).Body(
					jen.For(jen.ID("j").Assign().Zero(), jen.ID("j").LessThan().Int().Call(jen.ID("dataCount")), jen.ID("j").Increment()).Body(
						jen.ID("iterationLogger").Assign().ID("userLogger").Dot("WithValue").Call(jen.Lit("creating"), jen.Lit(typ.Name.PluralCommonName())).Dot("WithValue").Call(jen.Lit("iteration"), jen.ID("j")),
						jen.Newline(),
						jen.Commentf("create %s", typ.Name.SingularCommonName()),
						jen.List(jen.IDf("created%s", typ.Name.Singular()), jen.IDf("%sCreationErr", typ.Name.UnexportedVarName())).Assign().ID("userClient").Dotf("Create%s", typ.Name.Singular()).Call(constants.CtxVar(), jen.Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInput", typ.Name.Singular())).Call()),
						jen.If(jen.IDf("%sCreationErr", typ.Name.UnexportedVarName()).DoesNotEqual().Nil()).Body(
							jen.ID("quitter").Dot("ComplainAndQuit").Call(jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("creating %s ", typ.Name.SingularCommonName())+"#%d: %w"), jen.ID("j"), jen.IDf("%sCreationErr", typ.Name.UnexportedVarName()))),
						),
						jen.Newline(),
						jen.ID("iterationLogger").Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("created%s", typ.Name.Singular()).Dot("ID")).Dot("Debug").Call(jen.Litf("created %s", typ.Name.SingularCommonName())),
					),
					jen.ID("wg").Dot("Done").Call(),
				).Call(jen.ID("wg")),
			)
		}
	}

	return initializers
}
