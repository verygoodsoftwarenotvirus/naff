# Recipe Validation Improvements Plan

Based on recipe graph fixes applied to Sous Vide Chicken Breast, Soy Sauce Braised Chicken Thighs, Roasted Brussels Sprouts, and Carne Asada, this document proposes changes to recipe validation to catch similar issues earlier.

---

## Summary of Issues Fixed

| Issue                     | Recipe(s)                                                  | Root Cause                                                                                                                         | Fix Applied                                                                       |
|---------------------------|------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------|
| Unused preheated vessel   | Sous Vide Chicken Breast                                   | Step 5 vessel lacked `ProductOfRecipeStepProductIndex` and used wrong `ValidPreparationVesselID` (heat vs sous-vide)               | Corrected vessel bridge and product reference                                     |
| Preparation mismatch      | Sous Vide Chicken Breast                                   | Step 5 used `heatSousVideCookerVPI` (heat prep) instead of `sousVideCookerVPI` (sous-vide prep)                                    | Use preparation-specific instrument bridge                                        |
| Orphan step outputs       | Soy Sauce Braised Chicken Thighs                           | Step 2 (dry) and Step 3 (season) produced outputs not consumed; Step 15 (transfer) produced "seared chicken on plate" not consumed | Added `ProductOfRecipeStepIndex` references so downstream steps consume them      |
| Preheated vessel fan-out  | Roasted Brussels Sprouts, Soy Sauce Braised Chicken Thighs | Step 4/6 preheated vessel consumed directly by many steps instead of chaining                                                      | Chained vessel: preheat → remove → place → return → roast → stir → rotate → roast |
| Missing MeasurementUnitID | Carne Asada                                                | Grind product lacked `MeasurementUnitID`; graph resolution may require it for ingredient products                                  | Added `MeasurementUnitID` to grind product                                        |
| Orphan prep task outputs  | Carne Asada                                                | Toast→grind produced spices not linked to blend; unrefrigerate produced marinade container not linked to slice                     | Added product references and slice+sealed container vessel                        |

---

## Current Validation Scope

The `RecipeValidator` (`internal/domain/mealplanning/recipevalidator/recipe_validator.go`) currently:

1. **Validates bridge table IDs** – VIP, VIMU, VPI, VPV must exist and match the step's preparation
2. **Validates option grouping** – Index/option_index combinations must be sequential
3. **Populates derived fields** – IngredientID, MeasurementUnitID, InstrumentID, VesselID from bridge tables

**Gaps:** No validation of product references, graph structure, orphan detection, or vessel chaining.

---

## Proposed Validation Enhancements

### 1. Product Reference Validation (High Priority)

**When:** During `ValidateAndPopulate`, for ingredients/instruments/vessels with `ProductOfRecipeStepIndex` set.

**Checks:**

- Referenced step index exists (`< len(input.Steps)`)
- Referenced step has a product at `ProductOfRecipeStepProductIndex`
- For vessels: `ValidPreparationVesselID` is set and matches the consuming step's preparation (recipe step product vessels currently skip bridge validation)

**Rationale:** Prevents invalid references that would cause graph resolution to fail or produce wrong edges.

**Implementation note:** Validator receives `RecipeDatabaseCreationInput` which has steps. Products are defined per-step. We can validate references without needing the full recipe to be persisted.

---

### 2. Orphan Product Detection (Medium Priority – Optional/Warning)

**When:** After building the consumption graph from the recipe (post-creation or in a lint/analysis pass).

**Check:** Flag steps that produce products (ingredient or vessel type) which are never consumed by any downstream step.

**Rationale:** Orphan products indicate missing `ProductOfRecipeStepIndex` references and lead to disconnected graph nodes.

**Considerations:**

- Final step's output is often intentionally unconsumed (e.g., "sliced carne asada")
- Optional steps may produce outputs that are conditionally consumed
- Could be a **warning** rather than an error, or a separate analysis/lint tool

---

### 3. Vessel Fan-Out Warning (Low Priority – Optional)

**When:** Graph analysis pass.

**Check:** When a step's vessel product is consumed by more than N steps (e.g., 3+), flag as potential chaining opportunity.

**Rationale:** Preheated vessels (oven, water bath, baking sheets) should typically flow: preheat → first consumer → second consumer → … rather than preheat → many consumers.

**Considerations:**

- Some recipes may legitimately have one vessel feed many steps (e.g., shared oven)
- Subjective; better as a **warning** or documentation guideline than a hard error

---

### 4. MeasurementUnitID for Ingredient Products (Low Priority – Optional)

**When:** During validation of recipe step products (in creation input conversion or a separate pass).

**Check:** For `RecipeStepProductIngredientType` products, recommend or require `MeasurementUnitID` when the product may be consumed by downstream steps.

**Rationale:** Graph resolution (`findCreatedRecipeStepProductsForIngredients`) inherits MeasurementUnitID from the product when the consuming ingredient doesn't have it. Some resolution paths may depend on it.

**Considerations:**

- May not be strictly required for all flows
- Could be a **warning** to improve robustness

---

### 5. Vessel Product Must Have ValidPreparationVesselID When Consumed (High Priority)

**When:** During vessel validation.

**Check:** When a vessel has `ProductOfRecipeStepIndex` set (consuming a product from a previous step), it must also have `ValidPreparationVesselID` set, and that VPV must match the consuming step's preparation.

**Rationale:** The current validator skips bridge validation for recipe step product vessels (`isRecipeStepProductVessel` returns true). But consumed vessels still need the correct preparation bridge for the graph and for consistency.

**Implementation:** Extend `validateAndPopulateVessel` to validate `ValidPreparationVesselID` even when `ProductOfRecipeStepIndex` is set, ensuring the VPV matches the step's preparation.

---

### 6. Preparation Match for Recipe Step Product Instruments (Medium Priority)

**When:** During instrument validation.

**Check:** When an instrument has `ProductOfRecipeStepIndex` (consuming an instrument product), ensure any `ValidPreparationInstrumentID` matches the step's preparation. (Recipe step product instruments currently skip validation entirely.)

**Rationale:** Similar to the sous vide case—if we ever have instrument products flowing between steps, the preparation must match.

---

## Implementation Phases

### Phase 1: Critical Fixes (Recommended)

- **1. Product Reference Validation** – Validate step index and product index exist
- **5. Vessel Product ValidPreparationVesselID** – Require VPV when consuming vessel product

### Phase 2: Robustness

- **6. Preparation Match for Recipe Step Product Instruments** – If/when instrument products are used
- **4. MeasurementUnitID for Ingredient Products** – As warning or recommendation

### Phase 3: Analysis / Lint Tool (Optional)

- **2. Orphan Product Detection** – As a separate analysis or lint step
- **3. Vessel Fan-Out Warning** – As a documentation or lint guideline

---

## Files to Modify

| File                                                               | Changes                                                                                                |
|--------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|
| `internal/domain/mealplanning/recipevalidator/recipe_validator.go` | Add product reference validation, extend vessel validation for recipe step products                    |
| `internal/domain/mealplanning/recipe_step_ingredient.go`           | Possibly add `RecipeStepProductIngredientCreationRequestInput` if validation needs creation-time types |
| `internal/domain/mealplanning/recipeanalysis/recipe_analyzer.go`   | Optional: add `AnalyzeRecipeForOrphans` or similar for lint tool                                       |

---

## Testing Strategy

1. **Unit tests** – Add cases to `recipe_validator_test.go` for:
   - Invalid `ProductOfRecipeStepIndex` (out of range)
   - Invalid `ProductOfRecipeStepProductIndex` (product doesn't exist)
   - Vessel with `ProductOfRecipeStepIndex` but missing/mismatched `ValidPreparationVesselID`

2. **Integration tests** – Ensure bootstrap recipes still pass validation after changes

3. **Regression** – Re-run dry init with all bootstrap recipes to confirm no new failures

---

## References

- Conversation fixes: Sous Vide Chicken Breast, Soy Sauce Braised Chicken Thighs, Roasted Brussels Sprouts, Carne Asada
- `findCreatedRecipeStepProductsForIngredients` – `internal/repositories/postgres/mealplanning/recipes.go`
- Graph building – `internal/domain/mealplanning/recipeanalysis/recipe_analyzer.go`
