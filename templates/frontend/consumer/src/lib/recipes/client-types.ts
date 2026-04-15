/**
 * Client-safe types for recipe creation. NO imports from mealplanning or generated code—
 * those pull in @grpc/grpc-js (Node-only) and break in the browser.
 */

import type { MealComponentType, RecipeStepProductType } from './client-enums';

// Range types (minimal)
export interface Float32Range {
  min: number;
  max?: number;
}
export interface Uint16Range {
  min: number;
  max?: number;
}
export interface Uint32Range {
  min: number;
  max?: number;
}

// API response types (from search endpoints)
export interface ValidPreparation {
  id: string;
  name?: string;
  label?: string;
  [key: string]: unknown;
}
export interface ValidIngredient {
  id: string;
  name?: string;
  [key: string]: unknown;
}
export interface ValidIngredientPreparation {
  id: string;
  [key: string]: unknown;
}
export interface ValidIngredientMeasurementUnit {
  id: string;
  measurementUnit?: { id?: string; name?: string };
  [key: string]: unknown;
}
export interface ValidPreparationInstrument {
  id: string;
  [key: string]: unknown;
}
export interface ValidPreparationVessel {
  id: string;
  [key: string]: unknown;
}
export interface ValidMeasurementUnit {
  id: string;
  [key: string]: unknown;
}
export interface ValidIngredientState {
  id: string;
  [key: string]: unknown;
}
export interface ValidVessel {
  id: string;
  [key: string]: unknown;
}

// Recipe creation request types
export interface RecipeStepCompletionConditionCreationRequestInput {
  ingredientStateId: string;
  belongsToRecipeStep: string;
  notes: string;
  ingredients: number[];
  optional: boolean;
}

export interface RecipeStepIngredientCreationRequestInput {
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  vesselIndex?: number;
  productPercentageToUse?: number;
  recipeStepProductRecipeId?: string;
  recipeStepProductRecipeSlug?: string;
  ingredientNotes: string;
  name: string;
  quantityNotes: string;
  quantity: Float32Range | undefined;
  optionIndex: number;
  optional: boolean;
  toTaste: boolean;
  validIngredientPreparationId?: string;
  validIngredientMeasurementUnitId?: string;
  index?: number;
}

export interface RecipeStepInstrumentCreationRequestInput {
  recipeStepProductId?: string;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  name: string;
  notes: string;
  quantity: Float32Range | undefined;
  optionIndex: number;
  optional: boolean;
  preferenceRank: number;
}

export interface RecipeStepProductCreationRequestInput {
  storageTemperatureInCelsius?: Float32Range;
  storageDurationInSeconds?: Uint32Range;
  measurementQuantity?: Float32Range;
  itemQuantity?: Float32Range;
  measurementUnitId?: string;
  containedInVesselIndex?: number;
  quantityNotes: string;
  name: string;
  storageInstructions: string;
  type: RecipeStepProductType;
  index: number;
  compostable: boolean;
  isLiquid: boolean;
  isWaste: boolean;
}

export interface RecipeStepVesselCreationRequestInput {
  recipeStepProductId?: string;
  productOfRecipeStepIndex?: number;
  productOfRecipeStepProductIndex?: number;
  quantity: Uint16Range | undefined;
  name: string;
  notes: string;
  vesselPreposition: string;
  unavailableAfterStep: boolean;
  validPreparationVesselId?: string;
  index?: number;
  optionIndex: number;
}

export interface RecipeStepCreationRequestInput {
  estimatedTimeInSeconds?: Uint32Range;
  temperatureInCelsius?: Float32Range;
  preparationId: string;
  notes: string;
  conditionExpression: string;
  explicitInstructions: string;
  instruments: RecipeStepInstrumentCreationRequestInput[];
  vessels: RecipeStepVesselCreationRequestInput[];
  products: RecipeStepProductCreationRequestInput[];
  ingredients: RecipeStepIngredientCreationRequestInput[];
  completionConditions: RecipeStepCompletionConditionCreationRequestInput[];
  index: number;
  optional: boolean;
  startTimerAutomatically: boolean;
}

export interface RecipeMediaCreationRequestInput {
  belongsToRecipe: string;
  belongsToRecipeStep: string;
  mimeType: string;
  internalPath: string;
  externalPath: string;
  index: number;
}

export interface RecipePrepTaskStepWithinRecipeCreationRequestInput {
  belongsToRecipeStepIndex: number;
  satisfiesRecipeStep: boolean;
}

export interface RecipePrepTaskWithinRecipeCreationRequestInput {
  storageTemperatureInCelsius?: Float32Range;
  timeBufferBeforeRecipeInSeconds?: Uint32Range;
  storageType: string;
  name: string;
  description: string;
  explicitStorageInstructions: string;
  notes: string;
  belongsToRecipe: string;
  recipeSteps: RecipePrepTaskStepWithinRecipeCreationRequestInput[];
  optional: boolean;
}

export interface RecipeCreationRequestInput {
  inspiredByRecipeId?: string;
  name: string;
  source: string;
  sourceIsbn: string;
  description: string;
  pluralPortionName: string;
  portionName: string;
  slug: string;
  yieldsComponentType: MealComponentType;
  estimatedPortions: Float32Range | undefined;
  prepTasks: RecipePrepTaskWithinRecipeCreationRequestInput[];
  steps: RecipeStepCreationRequestInput[];
  alsoCreateMeal: boolean;
  eligibleForMeals: boolean;
  media: RecipeMediaCreationRequestInput[];
}
