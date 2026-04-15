package templates

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/pipeline"
)

func TestFullPipelineIntegration(t *testing.T) {
	t.Parallel()

	// Load the integration test config.
	cfg, err := config.LoadFromFile("../testdata/integration_project.yaml")
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	if err = cfg.Validate(); err != nil {
		t.Fatalf("validating config: %v", err)
	}

	// Generate into a temp directory.
	outputDir := t.TempDir()

	pipeline.TemplateFS = FS
	p := pipeline.New(cfg, outputDir, false)

	if err = p.Run(context.Background()); err != nil {
		t.Fatalf("pipeline run failed: %v", err)
	}

	// Verify files were generated.
	var generatedFiles []string
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(outputDir, path)
			generatedFiles = append(generatedFiles, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking output: %v", err)
	}

	if len(generatedFiles) == 0 {
		t.Fatal("no files were generated")
	}

	t.Logf("Generated %d files", len(generatedFiles))

	// --- Verify user domain files exist ---

	expectedFiles := []string{
		// Recipe domain layer
		"internal/domain/recipes/recipe.generated.go",
		"internal/domain/recipes/repository.generated.go",
		"internal/domain/recipes/keys/keys.generated.go",
		"internal/domain/recipes/fakes/fake.generated.go",
		"internal/domain/recipes/converters/recipe.generated.go",
		"internal/domain/recipes/mock/repository.generated.go",
		// Recipe manager
		"internal/domain/recipes/manager/interface.generated.go",
		"internal/domain/recipes/manager/manager.generated.go",
		"internal/domain/recipes/manager/do.generated.go",
		"internal/domain/recipes/manager/mock/manager.generated.go",
		// Recipe repository
		"internal/repositories/postgres/recipes/client.generated.go",
		"internal/repositories/postgres/recipes/recipe.generated.go",
		"internal/repositories/postgres/recipes/do.generated.go",
		// Recipe codegen
		"cmd/tools/codegen/queries/recipes_recipes.generated.go",
		// Recipe migration
		"internal/repositories/postgres/migrations/migration_files/recipes_recipe.generated.sql",
		// Recipe proto
		"proto/recipes/recipe_messages.generated.proto",
		"proto/recipes/recipe_service.generated.proto",
		// Recipe gRPC
		"internal/services/recipes/grpc/service.generated.go",
		"internal/services/recipes/grpc/recipe.generated.go",
		"internal/services/recipes/grpc/converters/converters.generated.go",
		"internal/services/recipes/grpc/permissions.generated.go",
		"internal/services/recipes/grpc/do.generated.go",
		// Recipe authorization
		"internal/authorization/recipes_permissions.generated.go",

		// RecipeStep files (same domain, second entity)
		"internal/domain/recipes/recipestep.generated.go",
		"internal/domain/recipes/converters/recipestep.generated.go",
		"internal/repositories/postgres/recipes/recipestep.generated.go",
		"cmd/tools/codegen/queries/recipes_recipe_steps.generated.go",
		"internal/services/recipes/grpc/recipestep.generated.go",

		// Built-in domain files
		"internal/domain/webhooks/webhook.generated.go",
		"internal/domain/issuereports/issuereport.generated.go",
		"internal/domain/comments/comment.generated.go",
		"internal/domain/audit/auditlogentry.generated.go",
	}

	fileSet := make(map[string]bool)
	for _, f := range generatedFiles {
		fileSet[f] = true
	}

	for _, expected := range expectedFiles {
		if !fileSet[expected] {
			t.Errorf("expected file not generated: %s", expected)
		}
	}

	// --- Verify all .generated.go files have valid content ---

	var goFileCount, nonGoFileCount int
	for _, f := range generatedFiles {
		fullPath := filepath.Join(outputDir, f)
		content, readErr := os.ReadFile(fullPath)
		if readErr != nil {
			t.Errorf("reading %s: %v", f, readErr)
			continue
		}

		if len(content) == 0 {
			t.Errorf("file %s is empty", f)
			continue
		}

		if strings.HasSuffix(f, ".generated.go") {
			goFileCount++
			// Verify Go files have a package declaration.
			if !strings.Contains(string(content), "package ") {
				t.Errorf("Go file %s missing package declaration", f)
			}
		} else {
			nonGoFileCount++
		}
	}

	t.Logf("Go files: %d, Non-Go files: %d", goFileCount, nonGoFileCount)

	// --- Spot-check specific file contents ---

	// Recipe entity types
	recipeEntity, err := os.ReadFile(filepath.Join(outputDir, "internal/domain/recipes/recipe.generated.go"))
	if err != nil {
		t.Fatalf("reading recipe entity: %v", err)
	}
	recipeContent := string(recipeEntity)

	recipeChecks := []string{
		"package recipes",
		"Recipe struct",
		"RecipeCreationRequestInput struct",
		"RecipeDatabaseCreationInput struct",
		"RecipeUpdateRequestInput struct",
		"RecipeDataManager interface",
		"GetRecipe(",
		"GetRecipes(",
		"GetRecipesForAccount(",
		"CreateRecipe(",
		"UpdateRecipe(",
		"ArchiveRecipe(",
		"BelongsToAccount",
		"CreatedByUser",
		// Field-specific checks
		`json:"name"`,
		`json:"description"`,
		`json:"sealOfApproval"`,
		"MinEstimatedPortions",
		"MaxEstimatedPortions",
		"InspiredByRecipeID",
		// SealOfApproval should NOT be in creation input (creatable: false)
		"func (x *Recipe) Update(",
		"validation.Required",
	}

	for _, check := range recipeChecks {
		if !strings.Contains(recipeContent, check) {
			t.Errorf("recipe entity missing %q", check)
		}
	}

	// SealOfApproval should NOT be in CreationRequestInput
	// Find the creation input struct and check it doesn't have SealOfApproval
	creationIdx := strings.Index(recipeContent, "RecipeCreationRequestInput struct")
	if creationIdx > 0 {
		// Find the closing brace
		nextBrace := strings.Index(recipeContent[creationIdx:], "}")
		if nextBrace > 0 {
			creationBlock := recipeContent[creationIdx : creationIdx+nextBrace]
			if strings.Contains(creationBlock, "SealOfApproval") {
				t.Error("SealOfApproval should NOT be in RecipeCreationRequestInput (creatable: false)")
			}
		}
	}

	// RecipeStep should have BelongsTo
	recipeStepEntity, err := os.ReadFile(filepath.Join(outputDir, "internal/domain/recipes/recipestep.generated.go"))
	if err != nil {
		t.Fatalf("reading recipestep entity: %v", err)
	}
	recipeStepContent := string(recipeStepEntity)

	if !strings.Contains(recipeStepContent, "BelongsToRecipe") {
		t.Error("RecipeStep should have BelongsToRecipe field")
	}

	// Migration SQL
	recipeMigration, err := os.ReadFile(filepath.Join(outputDir, "internal/repositories/postgres/migrations/migration_files/recipes_recipe.generated.sql"))
	if err != nil {
		t.Fatalf("reading recipe migration: %v", err)
	}
	migrationContent := string(recipeMigration)

	migrationChecks := []string{
		"CREATE TABLE IF NOT EXISTS recipes",
		"id TEXT NOT NULL PRIMARY KEY",
		"name TEXT NOT NULL",
		"description TEXT NOT NULL",
		"min_estimated_portions REAL NOT NULL",
		"max_estimated_portions REAL",
		"seal_of_approval BOOLEAN NOT NULL",
		"created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()",
		"belongs_to_account TEXT NOT NULL REFERENCES accounts",
		"created_by_user TEXT NOT NULL REFERENCES users",
	}

	for _, check := range migrationChecks {
		if !strings.Contains(migrationContent, check) {
			t.Errorf("recipe migration missing %q\n\nContent:\n%s", check, migrationContent)
		}
	}

	// Proto messages
	protoMessages, err := os.ReadFile(filepath.Join(outputDir, "proto/recipes/recipe_messages.generated.proto"))
	if err != nil {
		t.Fatalf("reading proto messages: %v", err)
	}
	protoContent := string(protoMessages)

	protoChecks := []string{
		"syntax = \"proto3\"",
		"package recipes",
		"message Recipe",
		"message RecipeCreationRequestInput",
		"message RecipeUpdateRequestInput",
		"string name",
		"float min_estimated_portions",
	}

	for _, check := range protoChecks {
		if !strings.Contains(protoContent, check) {
			t.Errorf("proto messages missing %q", check)
		}
	}

	// Authorization permissions
	authPerms, err := os.ReadFile(filepath.Join(outputDir, "internal/authorization/recipes_permissions.generated.go"))
	if err != nil {
		t.Fatalf("reading auth permissions: %v", err)
	}
	authContent := string(authPerms)

	authChecks := []string{
		"package authorization",
		"CreateRecipesPermission",
		"ReadRecipesPermission",
		"UpdateRecipesPermission",
		"ArchiveRecipesPermission",
		`"create.recipes"`,
		`"read.recipes"`,
	}

	for _, check := range authChecks {
		if !strings.Contains(authContent, check) {
			t.Errorf("authorization missing %q", check)
		}
	}
}

func TestPipelineIdempotency(t *testing.T) {
	t.Parallel()

	cfg, err := config.LoadFromFile("../testdata/integration_project.yaml")
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	outputDir := t.TempDir()
	pipeline.TemplateFS = FS

	// First run.
	p := pipeline.New(cfg, outputDir, false)
	if err = p.Run(context.Background()); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	// Count files from first run.
	var firstRunCount int
	filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			firstRunCount++
		}
		return nil
	})

	// Second run -- should report all files unchanged.
	p2 := pipeline.New(cfg, outputDir, false)
	if err = p2.Run(context.Background()); err != nil {
		t.Fatalf("second run failed: %v", err)
	}

	// Count files from second run (should be same).
	var secondRunCount int
	filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			secondRunCount++
		}
		return nil
	})

	if firstRunCount != secondRunCount {
		t.Errorf("file count changed between runs: %d -> %d", firstRunCount, secondRunCount)
	}
}
