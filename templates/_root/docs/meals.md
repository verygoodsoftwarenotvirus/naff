# Meals

A `Meal` represents a complete dining experience composed of multiple recipes with specific scaling and component types. Meals are the building blocks of meal planning and can be used to create meal plan options that users can vote on.

**Key Context**: In this system, a `Recipe` is a detailed cooking instruction (ingredients, steps, prep time, etc.) that serves a specific number of people. A `Meal` combines one or more recipes into a cohesive dining experience, with each recipe scaled appropriately for the meal's target serving size.

## Core Fields

### `EstimatedPortions` (Float32RangeWithOptionalMax)

Defines the serving size range for this meal:

- `Min` (float32, required): Minimum number of portions the meal serves
- `Max` (*float32, optional): Maximum number of portions the meal serves

This range helps users understand how many people the meal can feed and is used in meal planning calculations.

### `Components` ([]*MealComponent)

An array of recipes that make up this meal, each with specific metadata:

#### `MealComponent` Structure

- **`ComponentType`** (string): The role this recipe plays in the meal. Must be one of:
  - `"unspecified"` - Default/unknown type
  - `"amuse-bouche"` - Small appetizer
  - `"appetizer"` - Starter course
  - `"soup"` - Soup course
  - `"main"` - Main course (at least one required)
  - `"salad"` - Salad course
  - `"beverage"` - Drink component
  - `"side"` - Side dish
  - `"dessert"` - Dessert course

- **`Recipe`** (Recipe): The actual recipe object containing ingredients, steps, etc.
- **`RecipeScale`** (float32): Multiplier to adjust the recipe's portion size for this specific meal context

### `EligibleForMealPlans` (bool)

- **Purpose**: Whether this meal can be included in new meal plans
- **Primary Use Case**: Intended as a soft deletion mechanism for meals currently in active meal plans. Allows admins to prevent new usage of a meal while keeping existing meal plans functional, giving time to replace the meal in active plans
- **Note**: **This functionality is largely unimplemented** - the system doesn't currently enforce this flag when creating new meal plans or meal plan options

## Business Logic

### Meal Creation Requirements

- At least one component must have `ComponentType` of `"main"`
- All components must have valid `ComponentType` values
- `Name` and `Components` are required fields
- No restrictions on the number of components (you could have 1 main + 99 amuse-bouches)

### Meals vs Recipes

While a single-component meal might seem redundant compared to using a recipe directly, meals serve several purposes:

- **Complete Dining Experience**: Meals represent what you'd actually serve to guests, not just individual recipes
- **Proportional Scaling**: Components can be scaled independently to create balanced multi-course meals
- **Meal Planning Integration**: Meals are the smallest unit that can be selected during meal planning (users can't vote on individual recipes, only complete meals)
- **Future UI Organization**: Component types provide implied ordering and may affect display in future versions

For simple cases (like eating just a salmon filet), using a recipe directly is perfectly fine. Meals are designed for more complex, multi-component dining experiences.

### Portion Scaling

The `EstimatedPortions` field works in conjunction with `RecipeScale` in components to enable flexible portion scaling:

- **Meal's `EstimatedPortions`**: Defines the target serving size for the entire meal
- **Component's `RecipeScale`**: Multiplies the base recipe portions to match the meal's target size
- **Cumulative Scaling**: When a meal is used in meal planning, additional scaling can be applied to the entire meal

**Example**: A mashed potatoes recipe serves 10 people, but you want a 4-person meal:

- Set `RecipeScale: 0.4` for the mashed potatoes component
- Set `EstimatedPortions: {min: 4, max: 4}` for the meal
- Another user can then scale the entire meal to 6 people by setting the meal scale to 1.5 (adding `.5`), and the mashed potatoes will automatically scale to `.4 + (.5*.4) = 0.6`, and feed 6 people.

This system allows users to create properly proportioned multi-component meals without needing to manually adjust each recipe's portion calculations.

## Current Limitations & Future Considerations

### Search & Discovery

- Current search is limited to meal names and component types
- Search is powered by a search index (similar to Algolia) that indexes meal names and component types
- Future plans include vector search capabilities for more sophisticated meal discovery

### Editing & Lifecycle

- Meals can currently be edited even when used in active meal plans (no validation prevents this)
- Only the creating user can archive (soft-delete) a meal (`CreatedByUser` field enforces this)
- All changes are tracked in the database audit log for compliance and debugging
- Future versions may add validation to prevent editing meals in active meal plans

### Portion Precision

- Backend stores precise float32 values (e.g., 1.7348411 onions)
- Frontend clients handle formatting for display
- Future UI may include fractional shortcuts (1/8, 3/8, etc.) that convert to float equivalents

## Usage Examples

### Simple Single-Course Meal

```json
{
  "name": "Grilled Salmon",
  "description": "Perfectly grilled salmon with herbs",
  "estimatedPortions": { "min": 2, "max": 4 },
  "components": [
    {
      "componentType": "main",
      "recipe": { "id": "salmon-recipe-id" },
      "recipeScale": 1.0
    }
  ],
  "eligibleForMealPlans": true
}
```

### Multi-Course Meal with Scaling

```json
{
  "name": "Thanksgiving Dinner",
  "description": "Complete holiday feast",
  "estimatedPortions": { "min": 6, "max": 8 },
  "components": [
    {
      "componentType": "appetizer",
      "recipe": { "id": "cheese-plate-recipe" },
      "recipeScale": 0.5
    },
    {
      "componentType": "main",
      "recipe": { "id": "turkey-recipe" },
      "recipeScale": 1.0
    },
    {
      "componentType": "side",
      "recipe": { "id": "stuffing-recipe" },
      "recipeScale": 1.2
    },
    {
      "componentType": "dessert",
      "recipe": { "id": "pie-recipe" },
      "recipeScale": 0.8
    }
  ],
  "eligibleForMealPlans": true
}
```

**Scaling Explanation**:

- The turkey recipe serves 8 people, so `recipeScale: 1.0` means it serves the full 6-8 portions
- The cheese plate recipe serves 12 people, so `recipeScale: 0.5` scales it down to serve 6 people
- The stuffing recipe serves 5 people, so `recipeScale: 1.2` scales it up to serve 6 people
- The pie recipe serves 10 people, so `recipeScale: 0.8` scales it down to serve 8 people

When this meal is used in meal planning and scaled to 10 people, all components will be proportionally scaled up.
