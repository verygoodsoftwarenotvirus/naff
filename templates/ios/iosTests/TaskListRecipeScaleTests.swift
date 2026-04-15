//
//  TaskListRecipeScaleTests.swift
//  iosTests
//
//  Unit tests for recipe scale preservation when navigating from meal plan tasks.
//  Validates that component.recipeScale * option.mealScale is correctly computed.
//

import Foundation
import SwiftProtobuf
@testable import ios
import Testing

// MARK: - Helpers

private func makeRecipe(id: String, name: String = "Test Recipe") -> Mealplanning_Recipe {
  var recipe = Mealplanning_Recipe()
  recipe.id = id
  recipe.name = name
  return recipe
}

private func makeMealComponent(
  recipe: Mealplanning_Recipe,
  recipeScale: Float
) -> Mealplanning_MealComponent {
  var component = Mealplanning_MealComponent()
  component.recipe = recipe
  component.componentType = .main
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

private func makeMealPlanOption(
  meal: Mealplanning_Meal,
  mealScale: Float
) -> Mealplanning_MealPlanOption {
  var option = Mealplanning_MealPlanOption()
  option.id = "option-1"
  option.meal = meal
  option.mealScale = mealScale
  return option
}

private func makeTask(mealPlanOption: Mealplanning_MealPlanOption) -> Mealplanning_MealPlanTask {
  var task = Mealplanning_MealPlanTask()
  task.id = "task-1"
  task.mealPlanOption = mealPlanOption
  return task
}

private func makeTaskWithoutMealOption() -> Mealplanning_MealPlanTask {
  var task = Mealplanning_MealPlanTask()
  task.id = "task-2"
  return task
}

// MARK: - Recipe Scale from Meal Plan Tests

@Suite(.serialized)
struct TaskListRecipeScaleTests {
  @Test("Preserves 0.5x scale when meal has recipe at 0.5x and option at 1x")
  func testPreservesHalfScale() {
    let recipe = makeRecipe(id: "recipe-chicken")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipe, recipeScale: 0.5),
    ])
    let option = makeMealPlanOption(meal: meal, mealScale: 1.0)
    let task = makeTask(mealPlanOption: option)

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-chicken")

    #expect(scale == 0.5)
  }

  @Test("Preserves 0.5x scale when component is 1x and option mealScale is 0.5x")
  func testMealScaleMultiplier() {
    let recipe = makeRecipe(id: "recipe-soup")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipe, recipeScale: 1.0),
    ])
    let option = makeMealPlanOption(meal: meal, mealScale: 0.5)
    let task = makeTask(mealPlanOption: option)

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-soup")

    #expect(scale == 0.5)
  }

  @Test("Treats component recipeScale 0 as 1.0")
  func testComponentScaleZeroTreatedAsOne() {
    let recipe = makeRecipe(id: "recipe-default")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipe, recipeScale: 0),  // protobuf default
    ])
    let option = makeMealPlanOption(meal: meal, mealScale: 0.5)
    let task = makeTask(mealPlanOption: option)

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-default")

    #expect(scale == 0.5)  // 1.0 * 0.5
  }

  @Test("Treats option mealScale 0 as 1.0")
  func testOptionMealScaleZeroTreatedAsOne() {
    let recipe = makeRecipe(id: "recipe-main")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipe, recipeScale: 0.5),
    ])
    var option = Mealplanning_MealPlanOption()
    option.id = "option-1"
    option.meal = meal
    option.mealScale = 0  // protobuf default
    let task = makeTask(mealPlanOption: option)

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-main")

    #expect(scale == 0.5)  // 0.5 * 1.0
  }

  @Test("Returns nil when task has no meal plan option")
  func testReturnsNilWithoutMealOption() {
    let task = makeTaskWithoutMealOption()

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-any")

    #expect(scale == nil)
  }

  @Test("Returns nil when recipe not found in meal components")
  func testReturnsNilWhenRecipeNotFound() {
    let recipe = makeRecipe(id: "recipe-a")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipe, recipeScale: 0.5),
    ])
    let option = makeMealPlanOption(meal: meal, mealScale: 1.0)
    let task = makeTask(mealPlanOption: option)

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-nonexistent")

    #expect(scale == nil)
  }

  @Test("Correctly scales for specific recipe in multi-component meal")
  func testMultiComponentMealScale() {
    let recipeA = makeRecipe(id: "recipe-a")
    let recipeB = makeRecipe(id: "recipe-b")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipeA, recipeScale: 0.5),
      makeMealComponent(recipe: recipeB, recipeScale: 2.0),
    ])
    let option = makeMealPlanOption(meal: meal, mealScale: 1.0)
    let task = makeTask(mealPlanOption: option)

    let scaleA = task.recipeScaleForMealPlan(recipeID: "recipe-a")
    let scaleB = task.recipeScaleForMealPlan(recipeID: "recipe-b")

    #expect(scaleA == 0.5)
    #expect(scaleB == 2.0)
  }

  @Test("Combines component and meal scales (0.5 * 0.5 = 0.25)")
  func testCombinedScales() {
    let recipe = makeRecipe(id: "recipe-combined")
    let meal = makeMeal(components: [
      makeMealComponent(recipe: recipe, recipeScale: 0.5),
    ])
    let option = makeMealPlanOption(meal: meal, mealScale: 0.5)
    let task = makeTask(mealPlanOption: option)

    let scale = task.recipeScaleForMealPlan(recipeID: "recipe-combined")

    #expect(scale == 0.25)
  }
}
