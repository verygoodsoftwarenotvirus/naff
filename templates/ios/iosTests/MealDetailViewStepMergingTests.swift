//
//  MealDetailViewStepMergingTests.swift
//  iosTests
//
//  Unit tests for unified meal step merging (e.g. "grind 3g peppercorns" + "grind 2g peppercorns"
//  becomes "grind 5g peppercorns").
//

import Foundation
import SwiftProtobuf
@testable import ios
import Testing

// MARK: - Test Helpers

private func makeGrindPeppercornsStep(
  stepID: String,
  quantityGrams: Float
) -> Mealplanning_RecipeStep {
  var prep = Mealplanning_ValidPreparation()
  prep.id = "prep-grind"
  prep.name = "Grind"

  var ingredient = Mealplanning_ValidIngredient()
  ingredient.id = "ingredient-peppercorns"
  ingredient.name = "peppercorns"

  var unit = Mealplanning_ValidMeasurementUnit()
  unit.id = "unit-grams"
  unit.name = "grams"

  var stepIngredient = Mealplanning_RecipeStepIngredient()
  stepIngredient.id = "rsi-\(stepID)"
  stepIngredient.name = "peppercorns"
  stepIngredient.ingredient = ingredient
  var quantity = Common_Float32RangeWithOptionalMax()
  quantity.min = quantityGrams
  quantity.max = quantityGrams
  stepIngredient.quantity = quantity
  stepIngredient.measurementUnit = unit

  var instrument = Mealplanning_ValidInstrument()
  instrument.id = "instrument-mortar"
  instrument.name = "mortar and pestle"

  var stepInstrument = Mealplanning_RecipeStepInstrument()
  stepInstrument.instrument = instrument

  var step = Mealplanning_RecipeStep()
  step.id = stepID
  step.preparation = prep
  step.explicitInstructions = "Grind in mortar and pestle"
  step.ingredients = [stepIngredient]
  step.instruments = [stepInstrument]
  step.index = 0
  return step
}

private func makeChopOnionStep(stepID: String) -> Mealplanning_RecipeStep {
  var prep = Mealplanning_ValidPreparation()
  prep.id = "prep-chop"
  prep.name = "Chop"

  var ingredient = Mealplanning_ValidIngredient()
  ingredient.id = "ingredient-onion"
  ingredient.name = "onion"

  var stepIngredient = Mealplanning_RecipeStepIngredient()
  stepIngredient.id = "rsi-\(stepID)"
  stepIngredient.name = "onion"
  stepIngredient.ingredient = ingredient

  var step = Mealplanning_RecipeStep()
  step.id = stepID
  step.preparation = prep
  step.ingredients = [stepIngredient]
  step.index = 0
  return step
}

private func makeRecipe(
  id: String,
  name: String,
  steps: [Mealplanning_RecipeStep]
) -> Mealplanning_Recipe {
  var recipe = Mealplanning_Recipe()
  recipe.id = id
  recipe.name = name
  recipe.steps = steps
  return recipe
}

private func makeMealComponent(
  recipe: Mealplanning_Recipe,
  componentType: Mealplanning_MealComponentType = .main,
  recipeScale: Float = 1.0
) -> Mealplanning_MealComponent {
  var component = Mealplanning_MealComponent()
  component.recipe = recipe
  component.componentType = componentType
  component.recipeScale = recipeScale
  return component
}

private func makeMeal(components: [Mealplanning_MealComponent]) -> Mealplanning_Meal {
  var meal = Mealplanning_Meal()
  meal.id = "meal-1"
  meal.name = "Test Meal"
  meal.components = components
  return meal
}

private func formatComponentType(_ type: Mealplanning_MealComponentType) -> String {
  switch type {
  case .main: return "Main"
  case .side: return "Side"
  default: return "Component"
  }
}

// MARK: - Step Merging Tests

@Suite(.serialized)
struct MealDetailViewStepMergingTests {
  @Test("Merges common grind steps from two recipes")
  @MainActor
  func testMergesCommonGrindSteps() async {
    let recipeA = makeRecipe(
      id: "recipe-a",
      name: "Recipe A",
      steps: [makeGrindPeppercornsStep(stepID: "step-a1", quantityGrams: 3)]
    )
    let recipeB = makeRecipe(
      id: "recipe-b",
      name: "Recipe B",
      steps: [makeGrindPeppercornsStep(stepID: "step-b1", quantityGrams: 2)]
    )

    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipeA),
      makeMealComponent(recipe: recipeB),
    ])

    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true

    let viewModelA = PerformRecipeViewModel(recipeID: "recipe-a", authManager: authManager)
    viewModelA.recipe = recipeA
    viewModelA.washHandsCompleted = true

    let viewModelB = PerformRecipeViewModel(recipeID: "recipe-b", authManager: authManager)
    viewModelB.recipe = recipeB
    viewModelB.washHandsCompleted = true

    let componentViewModels = ["recipe-a": viewModelA, "recipe-b": viewModelB]
    let loadedRecipes: [String: (recipe: Mealplanning_Recipe, scale: Float)] = [
      "recipe-a": (recipeA, 1.0),
      "recipe-b": (recipeB, 1.0),
    ]
    let baseComponentScales = ["recipe-a": Float(1.0), "recipe-b": Float(1.0)]

    let steps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: componentViewModels,
      loadedRecipes: loadedRecipes,
      baseComponentScales: baseComponentScales,
      mealScale: 1.0,
      formatComponentType: formatComponentType
    )

    #expect(steps.count == 1)
    #expect(steps[0].isMerged == true)
    #expect(steps[0].componentNamesForTag == "Combined (Recipe A, Recipe B)")
    #expect(steps[0].sources.count == 2)

    let mergedStep = steps[0].step
    #expect(mergedStep.ingredients.count == 1)
    #expect(mergedStep.ingredients[0].quantity.min == 5)
    #expect(mergedStep.ingredients[0].quantity.max == 5)
    #expect(mergedStep.ingredients[0].name == "peppercorns")
  }

  @Test("Does not merge steps with different preparations")
  @MainActor
  func testDoesNotMergeDifferentPreparations() async {
    let grindStep = makeGrindPeppercornsStep(stepID: "step-grind", quantityGrams: 3)
    let chopStep = makeChopOnionStep(stepID: "step-chop")

    let recipeA = makeRecipe(id: "recipe-a", name: "Recipe A", steps: [grindStep])
    let recipeB = makeRecipe(id: "recipe-b", name: "Recipe B", steps: [chopStep])

    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipeA),
      makeMealComponent(recipe: recipeB),
    ])

    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true

    let viewModelA = PerformRecipeViewModel(recipeID: "recipe-a", authManager: authManager)
    viewModelA.recipe = recipeA
    viewModelA.washHandsCompleted = true

    let viewModelB = PerformRecipeViewModel(recipeID: "recipe-b", authManager: authManager)
    viewModelB.recipe = recipeB
    viewModelB.washHandsCompleted = true

    let steps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: ["recipe-a": viewModelA, "recipe-b": viewModelB],
      loadedRecipes: [
        "recipe-a": (recipeA, 1.0),
        "recipe-b": (recipeB, 1.0),
      ],
      baseComponentScales: ["recipe-a": 1.0, "recipe-b": 1.0],
      mealScale: 1.0,
      formatComponentType: formatComponentType
    )

    #expect(steps.count == 2)
    #expect(steps[0].isMerged == false)
    #expect(steps[1].isMerged == false)
  }

  @Test("Applies recipe scale when merging quantities")
  @MainActor
  func testAppliesScaleWhenMerging() async {
    let recipeA = makeRecipe(
      id: "recipe-a",
      name: "Recipe A",
      steps: [makeGrindPeppercornsStep(stepID: "step-a1", quantityGrams: 4)]
    )
    let recipeB = makeRecipe(
      id: "recipe-b",
      name: "Recipe B",
      steps: [makeGrindPeppercornsStep(stepID: "step-b1", quantityGrams: 2)]
    )

    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipeA, recipeScale: 0.5),
      makeMealComponent(recipe: recipeB, recipeScale: 1.0),
    ])

    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true

    let viewModelA = PerformRecipeViewModel(recipeID: "recipe-a", authManager: authManager)
    viewModelA.recipe = recipeA
    viewModelA.washHandsCompleted = true

    let viewModelB = PerformRecipeViewModel(recipeID: "recipe-b", authManager: authManager)
    viewModelB.recipe = recipeB
    viewModelB.washHandsCompleted = true

    let steps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: ["recipe-a": viewModelA, "recipe-b": viewModelB],
      loadedRecipes: [
        "recipe-a": (recipeA, 0.5),
        "recipe-b": (recipeB, 1.0),
      ],
      baseComponentScales: ["recipe-a": 0.5, "recipe-b": 1.0],
      mealScale: 1.0,
      formatComponentType: formatComponentType
    )

    #expect(steps.count == 1)
    #expect(steps[0].isMerged == true)
    // 4 * 0.5 + 2 * 1.0 = 2 + 2 = 4
    #expect(steps[0].step.ingredients[0].quantity.min == 4)
  }

  @Test("Single recipe produces non-merged steps")
  @MainActor
  func testSingleRecipeNoMerging() async {
    let recipe = makeRecipe(
      id: "recipe-solo",
      name: "Solo Recipe",
      steps: [
        makeGrindPeppercornsStep(stepID: "step-1", quantityGrams: 5),
        makeChopOnionStep(stepID: "step-2"),
      ]
    )

    let meal = makeMeal(components: [makeMealComponent(recipe: recipe)])

    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true

    let viewModel = PerformRecipeViewModel(recipeID: "recipe-solo", authManager: authManager)
    viewModel.recipe = recipe
    viewModel.washHandsCompleted = true

    let steps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: ["recipe-solo": viewModel],
      loadedRecipes: ["recipe-solo": (recipe, 1.0)],
      baseComponentScales: ["recipe-solo": 1.0],
      mealScale: 1.0,
      formatComponentType: formatComponentType
    )

    #expect(steps.count == 2)
    #expect(steps[0].isMerged == false)
    #expect(steps[0].componentNamesForTag == "Solo Recipe")
    #expect(steps[1].isMerged == false)
  }

  @Test("UnifiedMealStep componentNamesForTag includes recipe names for merged")
  @MainActor
  func testComponentNamesForTagMerged() async {
    let recipeA = makeRecipe(
      id: "recipe-a",
      name: "Recipe A",
      steps: [makeGrindPeppercornsStep(stepID: "step-a1", quantityGrams: 1)]
    )
    let recipeB = makeRecipe(
      id: "recipe-b",
      name: "Recipe B",
      steps: [makeGrindPeppercornsStep(stepID: "step-b1", quantityGrams: 1)]
    )

    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipeA),
      makeMealComponent(recipe: recipeB),
    ])

    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true

    let viewModelA = PerformRecipeViewModel(recipeID: "recipe-a", authManager: authManager)
    viewModelA.recipe = recipeA
    viewModelA.washHandsCompleted = true

    let viewModelB = PerformRecipeViewModel(recipeID: "recipe-b", authManager: authManager)
    viewModelB.recipe = recipeB
    viewModelB.washHandsCompleted = true

    let steps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: ["recipe-a": viewModelA, "recipe-b": viewModelB],
      loadedRecipes: ["recipe-a": (recipeA, 1.0), "recipe-b": (recipeB, 1.0)],
      baseComponentScales: ["recipe-a": 1.0, "recipe-b": 1.0],
      mealScale: 1.0,
      formatComponentType: formatComponentType
    )

    #expect(steps[0].componentNamesForTag == "Combined (Recipe A, Recipe B)")
  }

  @Test("Does not merge when one step is already done (e.g. from prep) and the other is not")
  @MainActor
  func testDoesNotMergeWhenCompletionStatusDiffers() async {
    let recipeA = makeRecipe(
      id: "recipe-a",
      name: "Chicken",
      steps: [makeGrindPeppercornsStep(stepID: "step-a1", quantityGrams: 3)]
    )
    let recipeB = makeRecipe(
      id: "recipe-b",
      name: "Broccoli",
      steps: [makeGrindPeppercornsStep(stepID: "step-b1", quantityGrams: 2)]
    )

    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipeA),
      makeMealComponent(recipe: recipeB),
    ])

    let authManager = AuthenticationManager()
    authManager.isAuthenticated = true

    let viewModelA = PerformRecipeViewModel(recipeID: "recipe-a", authManager: authManager)
    viewModelA.recipe = recipeA
    viewModelA.washHandsCompleted = true
    // Simulate: user completed chicken's grind step as part of a prep task
    viewModelA.completedSteps.insert("recipe-a:step-a1")

    let viewModelB = PerformRecipeViewModel(recipeID: "recipe-b", authManager: authManager)
    viewModelB.recipe = recipeB
    viewModelB.washHandsCompleted = true

    let steps = collectUnifiedMealStepsWithMerging(
      meal: meal,
      componentViewModels: ["recipe-a": viewModelA, "recipe-b": viewModelB],
      loadedRecipes: ["recipe-a": (recipeA, 1.0), "recipe-b": (recipeB, 1.0)],
      baseComponentScales: ["recipe-a": 1.0, "recipe-b": 1.0],
      mealScale: 1.0,
      formatComponentType: formatComponentType
    )

    // Should NOT merge: chicken step is done, broccoli step still needs to be done
    #expect(steps.count == 2)
    #expect(steps[0].isMerged == false)
    #expect(steps[0].category == .done)
    #expect(steps[0].step.ingredients[0].quantity.min == 3)
    #expect(steps[1].isMerged == false)
    #expect(steps[1].category == .upNext)
    #expect(steps[1].step.ingredients[0].quantity.min == 2)
  }
}
