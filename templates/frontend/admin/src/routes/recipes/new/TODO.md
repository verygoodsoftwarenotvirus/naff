# Recipe Creation Page TODO

## Ingredient / step inputs

- **Recipe-as-ingredient**: Ingredients cannot yet be products of _established_ recipes. You can reference outputs of
  earlier steps in the same recipe, but you can't add "Caesar Dressing" (another recipe) as an ingredient in a salad.
  Needs `recipeStepProductRecipeId` (and possibly recipe search + slug resolution) in the ingredient flow.

## Recipe-level fields

- **Yields component type**: No UI to set `yieldsComponentType` (appetizer, main, side, dessert, etc.). State defaults
  to "main"; users cannot choose.
- **Estimated portions range**: Only a single "Est. Portions" number is shown. The model supports optional
  `estimatedPortions.max` (e.g. "4–6 servings"); need UI for optional max.
- **Prep tasks**: `prepTasks` exist in the type and state but there is no UI to add or edit advance-prep tasks (name,
  description, storage, which steps they satisfy).
- **Eligible for meals**: Field exists in state (default true) but is not exposed in the UI (optional; matters for
  ingredient-only recipes).
- **Inspired by / clone**: No way to set `inspiredByRecipeId` when creating from an existing recipe (clone workflow
  could be a separate flow).
- **Media**: No UI to attach recipe-level or step-level media.

## Step products (outputs of each step)

- **Product quantity/unit**: Step products only have name and type in the UI. The model requires proper product
  definitions for scaling:
  - **Continuous**: `measurementUnitId` + `measurementQuantity` (total amount, e.g. "2 cups sauce").
  - **Discrete**: `itemQuantity` (count) + `measurementQuantity` (per item) + `measurementUnitId` (e.g. "4 patties, 4 oz
    each"). Need UI for unit, quantity, and discrete vs continuous (item count).

## Option groups

- **Alternative ingredients/instruments/vessels**: No UI for option groups (same `index`, different `optionIndex`) to
  express "butter OR margarine", "stand mixer OR hand mixer", etc.

## Validation and UX

- **Minimum 2 steps**: Backend requires at least 2 steps; the form allows a single step. Either start with 2 steps, or
  validate before submit and show a clear error.
- **Step notes vs instructions**: Step has both `notes` and `explicitInstructions`; confirm both are captured if the
  backend uses `notes` (e.g. for internal/structured notes).
