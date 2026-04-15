//
//  MealPreviewHelper.swift
//  ios
//
//  Builds grocery list and prep task previews from meal draft components.
//

import Foundation
import SwiftProtobuf

// MARK: - Preview Types

struct MealPreviewGroceryItem: Identifiable {
  let id: String
  let displayText: String
  let name: String
}

struct MealPreviewPrepTask: Identifiable {
  let id: String
  let name: String
  let recipeName: String
  /// How far in advance this can be done (e.g. "Up to 3 hr in advance").
  let advanceText: String?
}

// MARK: - MealPreviewHelper

enum MealPreviewHelper {

  /// Aggregates ingredients from draft components into grocery items.
  /// For option groups (alternatives at same step index), uses optionIndex 0.
  static func groceryItems(from components: [CreateMealDraftComponent]) -> [MealPreviewGroceryItem]
  {
    var aggregated: [String: AggregatedIngredient] = [:]

    for component in components {
      let scale = component.recipeScale
      let recipeName = component.recipe.name.isEmpty ? "Unnamed Recipe" : component.recipe.name

      for step in component.recipe.steps {
        let ingredientsToAdd = ingredientsForPreview(from: step.ingredients)
        for ingredient in ingredientsToAdd {
          // Skip recipe step products (outputs of other steps) - not purchased ingredients
          guard !ingredient.hasRecipeStepProductID, !ingredient.hasRecipeStepProductRecipeID else {
            continue
          }
          addIngredientToAggregated(
            ingredient,
            scale: scale,
            recipeName: recipeName,
            aggregated: &aggregated
          )
        }
      }
    }

    return aggregated.values
      .map { item in
        let displayText: String
        if let qty = item.quantityText(scale: 1.0) {
          displayText = "\(qty) \(item.name)"
        } else {
          displayText = item.name
        }
        return MealPreviewGroceryItem(
          id: item.ingredientID,
          displayText: displayText,
          name: item.name
        )
      }
      .sorted { $0.name.localizedCaseInsensitiveCompare($1.name) == .orderedAscending }
  }

  /// Collects prep tasks from draft components.
  static func prepTasks(from components: [CreateMealDraftComponent]) -> [MealPreviewPrepTask] {
    var tasks: [MealPreviewPrepTask] = []
    var seenIDs: Set<String> = []

    for component in components {
      let recipeName = component.recipe.name.isEmpty ? "Unnamed Recipe" : component.recipe.name
      for prepTask in component.recipe.prepTasks {
        let taskID = "\(component.id)-\(prepTask.id)"
        guard !seenIDs.contains(taskID) else { continue }
        seenIDs.insert(taskID)
        let advanceText =
          prepTask.hasTimeBufferBeforeRecipeInSeconds
          ? RecipeTimeEstimation.formatTimeBufferInAdvance(prepTask.timeBufferBeforeRecipeInSeconds)
          : nil
        tasks.append(
          MealPreviewPrepTask(
            id: taskID,
            name: prepTask.name.isEmpty ? "Prep" : prepTask.name,
            recipeName: recipeName,
            advanceText: advanceText
          )
        )
      }
    }

    return tasks
  }

  // MARK: - Private

  private static func addIngredientToAggregated(
    _ ingredient: Mealplanning_RecipeStepIngredient,
    scale: Float,
    recipeName: String,
    aggregated: inout [String: AggregatedIngredient]
  ) {
    // Dedupe by canonical ingredient+unit when available; otherwise by normalized name+unit.
    let key: String
    if ingredient.hasIngredient {
      let id = ingredient.ingredient.id
      let unitID = ingredient.hasMeasurementUnit ? ingredient.measurementUnit.id : ""
      key = "ingredient:\(id)|\(unitID)"
    } else {
      let normalizedName = ingredient.name
        .trimmingCharacters(in: .whitespaces)
        .lowercased()
      let namePart = normalizedName.isEmpty ? "unnamed-\(ingredient.id)" : normalizedName
      let unitID = ingredient.hasMeasurementUnit ? ingredient.measurementUnit.id : ""
      key = "name:\(namePart)|\(unitID)"
    }
    if key.isEmpty { return }

    var item =
      aggregated[key]
      ?? AggregatedIngredient(
        ingredientID: key,
        name: ingredient.name,
        quantityNotes: ingredient.quantityNotes,
        measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil,
        sourceRecipeName: recipeName
      )

    if ingredient.hasQuantity {
      var scaledQty = ingredient.quantity
      scaledQty.min *= scale
      if scaledQty.hasMax {
        scaledQty.max *= scale
      }
      item.addQuantity(scaledQty)
    }

    aggregated[key] = item
  }

  /// Picks one ingredient per step index. For option groups, uses optionIndex 0.
  private static func ingredientsForPreview(
    from ingredients: [Mealplanning_RecipeStepIngredient]
  ) -> [Mealplanning_RecipeStepIngredient] {
    var byIndex: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]
    for (index, ingredient) in ingredients.enumerated() {
      let idx = UInt32(index)
      byIndex[idx, default: []].append(ingredient)
    }
    return byIndex.values.compactMap { group in
      group.min(by: { $0.optionIndex < $1.optionIndex })
    }
  }
}
