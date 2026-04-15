//
//  MealDetailView+StepMerging.swift
//  ios
//
//  Merges common steps across meal components (e.g. "grind 3g peppercorns" + "grind 2g peppercorns"
//  becomes "grind 5g peppercorns").
//

import SwiftProtobuf
import SwiftUI

// MARK: - Source Step (for merged steps)

struct UnifiedMealStepSource {
  let componentID: String
  let componentIndex: Int
  let componentName: String
  let step: Mealplanning_RecipeStep
  let stepIndex: Int
  let recipeID: String
  let scale: Float
  let viewModel: PerformRecipeViewModel
  let isAssociatedRecipeStep: Bool
  let associatedRecipeName: String?
}

// MARK: - Unified Meal Step (single or merged)

struct UnifiedMealStep: Identifiable {
  enum Category {
    case upNext
    case forLater
    case done
  }

  /// For single steps: one source. For merged steps: multiple sources.
  let sources: [UnifiedMealStepSource]
  /// The step to display (either original or merged with combined quantities)
  let step: Mealplanning_RecipeStep
  let stepIndex: Int  // Display index
  let category: Category

  /// Primary view model (first source) - used for StepDetailsView which needs recipe context
  var viewModel: PerformRecipeViewModel { sources[0].viewModel }
  var recipeID: String { sources[0].recipeID }
  var scale: Float { sources[0].scale }
  var componentIndex: Int { sources[0].componentIndex }
  var componentName: String { sources[0].componentName }

  var isMerged: Bool { sources.count > 1 }

  var isAssociatedRecipeStep: Bool { sources[0].isAssociatedRecipeStep }
  var associatedRecipeName: String? { sources[0].associatedRecipeName }

  var id: String {
    sources.map { "\($0.componentID):\($0.recipeID):\($0.step.id)" }.joined(separator: "|")
  }

  var componentNamesForTag: String {
    if isMerged {
      let recipeNames = sources.map { source in
        if source.isAssociatedRecipeStep, let name = source.associatedRecipeName, !name.isEmpty {
          return name
        }
        return source.componentName
      }
      var seen: Set<String> = []
      let uniqueNames = recipeNames.filter { seen.insert($0).inserted }
      let namesList = uniqueNames.joined(separator: ", ")
      return "Combined (\(namesList))"
    }
    if isAssociatedRecipeStep, let name = associatedRecipeName, !name.isEmpty {
      return name
    }
    return componentName
  }

  /// For merged steps with multiple sources contributing to an ingredient, returns a breakdown
  /// string per ingredient key, e.g. "2g from Recipe A, 3g from Recipe B".
  var ingredientBreakdownBySource: [String: String]? {
    guard isMerged else { return nil }
    var breakdown: [String: [MergedIngredientSourcePart]] = [:]

    for source in sources {
      let sourceName: String = {
        if source.isAssociatedRecipeStep, let name = source.associatedRecipeName, !name.isEmpty {
          return name
        }
        return source.componentName
      }()
      for ing in source.step.ingredients where ing.hasIngredient {
        let key = "\(ing.ingredient.id)|\(ing.hasMeasurementUnit ? ing.measurementUnit.id : "")"
        let unit = ing.hasMeasurementUnit ? ing.measurementUnit : nil
        if ing.hasQuantity {
          let effectiveScale = source.scale * (ing.scaleFactor > 0 ? ing.scaleFactor : 1.0)
          let min = ing.quantity.min * effectiveScale
          let max = ing.quantity.hasMax ? ing.quantity.max * effectiveScale : nil
          breakdown[key, default: []].append(
            MergedIngredientSourcePart(sourceName: sourceName, min: min, max: max, unit: unit)
          )
        }
      }
    }

    var formatted: [String: String] = [:]
    for (key, parts) in breakdown where parts.count > 1 {
      let partsStr = parts.map { part in
        let qtyStr = formatMergedIngredientQuantity(min: part.min, max: part.max, unit: part.unit)
        return "\(qtyStr) from \(part.sourceName)"
      }.joined(separator: ", ")
      formatted[key] = partsStr
    }
    return formatted.isEmpty ? nil : formatted
  }
}

private struct MergedIngredientSourcePart {
  let sourceName: String
  let min: Float
  let max: Float?
  let unit: Mealplanning_ValidMeasurementUnit?
}

private func formatMergedIngredientQuantity(
  min: Float, max: Float?, unit: Mealplanning_ValidMeasurementUnit?
) -> String {
  let unitName = MeasurementUnitFormatter.displayName(for: min, unit: unit)
  let unitSuffix = unitName.isEmpty ? "" : " \(unitName)"
  let formatMin = min.truncatingRemainder(dividingBy: 1) == 0 ? "%.0f" : "%.2f"
  if let maxVal = max {
    if min == maxVal {
      return (String(format: formatMin, min) + unitSuffix).trimmingCharacters(in: .whitespaces)
    }
    let formatMax = maxVal.truncatingRemainder(dividingBy: 1) == 0 ? "%.0f" : "%.2f"
    return (String(format: "\(formatMin) - \(formatMax)", min, maxVal) + unitSuffix)
      .trimmingCharacters(in: .whitespaces)
  }
  // No max: treat min as the authoritative single value (not unbounded)
  return (String(format: formatMin, min) + unitSuffix).trimmingCharacters(in: .whitespaces)
}

// MARK: - Step Merge Key

private func stepMergeKey(step: Mealplanning_RecipeStep) -> String? {
  // Only merge steps that use base ingredients (no product references)
  let hasProductRefs =
    step.ingredients.contains { $0.hasRecipeStepProductID }
    || step.instruments.contains { $0.hasRecipeStepProductID }
    || step.vessels.contains { $0.hasRecipeStepProductID }
  if hasProductRefs {
    return nil
  }

  let prepID = step.hasPreparation ? step.preparation.id : ""
  let ingredientIDs = step.ingredients
    .filter { $0.hasIngredient && !$0.ingredient.id.isEmpty }
    .map { $0.ingredient.id }
    .sorted()
  let instrumentIDs = step.instruments
    .filter { $0.hasInstrument && !$0.instrument.id.isEmpty }
    .map { $0.instrument.id }
    .sorted()
  let vesselIDs = step.vessels
    .filter { $0.hasVessel && !$0.vessel.id.isEmpty }
    .map { $0.vessel.id }
    .sorted()

  // Require at least preparation or ingredients for a meaningful merge
  if prepID.isEmpty && ingredientIDs.isEmpty {
    return nil
  }

  return [
    "p:\(prepID)",
    "i:\(ingredientIDs.joined(separator: ","))",
    "n:\(instrumentIDs.joined(separator: ","))",
    "v:\(vesselIDs.joined(separator: ","))",
  ].joined(separator: "|")
}

// MARK: - Merge Steps

private func mergeStepIngredients(
  _ sources: [UnifiedMealStepSource]
) -> [Mealplanning_RecipeStepIngredient] {
  // Group ingredients by (ingredientID, measurementUnitID) - same ingredient + unit
  var mergedByKey: [String: Mealplanning_RecipeStepIngredient] = [:]
  var totals: [String: (min: Float, max: Float?)] = [:]

  for source in sources {
    let scale = source.scale
    for ing in source.step.ingredients where ing.hasIngredient {
      let key = "\(ing.ingredient.id)|\(ing.hasMeasurementUnit ? ing.measurementUnit.id : "")"
      if mergedByKey[key] == nil {
        var copy = ing
        if ing.hasQuantity {
          let effectiveScale = scale * (ing.scaleFactor > 0 ? ing.scaleFactor : 1.0)
          var scaledQuantity = ing.quantity
          scaledQuantity.min *= effectiveScale
          if scaledQuantity.hasMax {
            scaledQuantity.max *= effectiveScale
          }
          copy.quantity = scaledQuantity
          totals[key] = (scaledQuantity.min, scaledQuantity.hasMax ? scaledQuantity.max : nil)
        }
        mergedByKey[key] = copy
      } else if ing.hasQuantity {
        let effectiveScale = scale * (ing.scaleFactor > 0 ? ing.scaleFactor : 1.0)
        var (totalMin, totalMax) = totals[key] ?? (0, nil)
        totalMin += ing.quantity.min * effectiveScale
        if ing.quantity.hasMax {
          totalMax = (totalMax ?? totalMin) + ing.quantity.max * effectiveScale
        } else {
          totalMax = nil
        }
        totals[key] = (totalMin, totalMax)
        guard var existing = mergedByKey[key] else { continue }
        var newQuantity = Common_Float32RangeWithOptionalMax()
        newQuantity.min = totalMin
        if let maxVal = totalMax {
          newQuantity.max = maxVal
        }
        existing.quantity = newQuantity
        mergedByKey[key] = existing
      }
    }
  }

  return mergedByKey.values.sorted { ($0.ingredient.id, $0.name) < ($1.ingredient.id, $1.name) }
}

private func createMergedStep(from sources: [UnifiedMealStepSource]) -> Mealplanning_RecipeStep {
  let first = sources[0].step
  var merged = Mealplanning_RecipeStep()
  merged.id = "merged-\(sources.map { $0.step.id }.joined(separator: "-"))"
  merged.preparation = first.preparation
  merged.explicitInstructions = first.explicitInstructions
  merged.notes = first.notes
  merged.optional = sources.contains { $0.step.optional }
  merged.ingredients = mergeStepIngredients(sources)
  merged.instruments = first.instruments  // Same instruments across merged steps
  merged.vessels = first.vessels  // Same vessels
  merged.completionConditions = first.completionConditions
  merged.index = first.index
  return merged
}

// MARK: - Collect and Merge

@MainActor
func collectUnifiedMealStepsWithMerging(
  meal: Mealplanning_Meal,
  componentViewModels: [String: PerformRecipeViewModel],
  loadedRecipes: [String: (recipe: Mealplanning_Recipe, scale: Float)],
  baseComponentScales: [String: Float],
  mealScale: Float,
  formatComponentType: (Mealplanning_MealComponentType) -> String
) -> [UnifiedMealStep] {
  var rawSteps: [UnifiedMealStepSource] = []

  for (componentIndex, component) in meal.components.enumerated() {
    let componentID = component.recipe.id
    guard let viewModel = componentViewModels[componentID], let recipe = viewModel.recipe else {
      continue
    }

    let scale =
      loadedRecipes[componentID]?.scale ?? (baseComponentScales[componentID] ?? 1.0) * mealScale
    let componentName =
      recipe.name.isEmpty ? formatComponentType(component.componentType) : recipe.name

    // Include associated (prerequisite) recipe steps first, then main recipe steps
    for associatedRecipe in recipe.associatedRecipes {
      for (index, step) in associatedRecipe.steps.enumerated() {
        rawSteps.append(
          UnifiedMealStepSource(
            componentID: componentID,
            componentIndex: componentIndex,
            componentName: componentName,
            step: step,
            stepIndex: index,
            recipeID: associatedRecipe.id,
            scale: scale,
            viewModel: viewModel,
            isAssociatedRecipeStep: true,
            associatedRecipeName: associatedRecipe.name.isEmpty ? nil : associatedRecipe.name
          )
        )
      }
    }

    for (index, step) in recipe.steps.enumerated() {
      rawSteps.append(
        UnifiedMealStepSource(
          componentID: componentID,
          componentIndex: componentIndex,
          componentName: componentName,
          step: step,
          stepIndex: index,
          recipeID: recipe.id,
          scale: scale,
          viewModel: viewModel,
          isAssociatedRecipeStep: false,
          associatedRecipeName: nil
        )
      )
    }
  }

  // Group by merge key (preserving order within group by first appearance)
  var groups: [String: [UnifiedMealStepSource]] = [:]
  for source in rawSteps {
    if let key = stepMergeKey(step: source.step) {
      groups[key, default: []].append(source)
    }
  }

  // Build result in original step order. When we hit a step that's in a multi-group,
  // partition by completion status: only merge steps that are all done or all not done.
  // Mixed completion (e.g. one done from prep, one not) must stay separate.
  var result: [UnifiedMealStep] = []
  var displayIndex = 0
  var consumedInGroup: Set<String> = []  // "recipeID:stepID" for steps we've output
  let sourceToRawIndex: [String: Int] = Dictionary(
    uniqueKeysWithValues: rawSteps.enumerated().map { ($1.recipeID + ":" + $1.step.id, $0) }
  )

  for source in rawSteps {
    let sourceKey = "\(source.recipeID):\(source.step.id)"
    if consumedInGroup.contains(sourceKey) {
      continue
    }

    if let key = stepMergeKey(step: source.step), let group = groups[key] {
      if group.count > 1 {
        let partitions = partitionGroupByCompletionStatus(group)
        let orderedPartitions = partitions.sorted { partition1, partition2 in
          let min1 =
            partition1.compactMap { sourceToRawIndex["\($0.recipeID):\($0.step.id)"] }.min()
            ?? Int.max
          let min2 =
            partition2.compactMap { sourceToRawIndex["\($0.recipeID):\($0.step.id)"] }.min()
            ?? Int.max
          return min1 < min2
        }
        for partition in orderedPartitions {
          if partition.count > 1 {
            let category = categorizeMergedStep(sources: partition)
            let mergedStep = createMergedStep(from: partition)
            result.append(
              UnifiedMealStep(
                sources: partition,
                step: mergedStep,
                stepIndex: displayIndex,
                category: category
              )
            )
          } else {
            let category = categorizeSingleStep(source: partition[0])
            result.append(
              UnifiedMealStep(
                sources: partition,
                step: partition[0].step,
                stepIndex: displayIndex,
                category: category
              )
            )
          }
          for sourceStep in partition {
            consumedInGroup.insert("\(sourceStep.recipeID):\(sourceStep.step.id)")
          }
          displayIndex += 1
        }
      } else {
        let category = categorizeSingleStep(source: source)
        result.append(
          UnifiedMealStep(
            sources: [source],
            step: source.step,
            stepIndex: displayIndex,
            category: category
          )
        )
        displayIndex += 1
      }
    } else {
      let category = categorizeSingleStep(source: source)
      result.append(
        UnifiedMealStep(
          sources: [source],
          step: source.step,
          stepIndex: displayIndex,
          category: category
        )
      )
      displayIndex += 1
    }
  }

  return result
}

/// Partitions a merge group by completion status. Steps that are already done (e.g. from prep tasks)
/// must not be merged with steps that are still pending—they need separate quantities.
@MainActor
private func partitionGroupByCompletionStatus(
  _ group: [UnifiedMealStepSource]
) -> [[UnifiedMealStepSource]] {
  var done: [UnifiedMealStepSource] = []
  var notDone: [UnifiedMealStepSource] = []
  for source in group {
    let category = source.viewModel.categorizeStep(
      recipeID: source.recipeID, stepID: source.step.id)
    if category == .done {
      done.append(source)
    } else {
      notDone.append(source)
    }
  }
  var result: [[UnifiedMealStepSource]] = []
  if !done.isEmpty { result.append(done) }
  if !notDone.isEmpty { result.append(notDone) }
  return result
}

@MainActor
private func categorizeSingleStep(source: UnifiedMealStepSource) -> UnifiedMealStep.Category {
  switch source.viewModel.categorizeStep(recipeID: source.recipeID, stepID: source.step.id) {
  case .upNext: return .upNext
  case .forLater: return .forLater
  case .done: return .done
  }
}

@MainActor
private func categorizeMergedStep(sources: [UnifiedMealStepSource]) -> UnifiedMealStep.Category {
  let categories = sources.map {
    source in
    source.viewModel.categorizeStep(recipeID: source.recipeID, stepID: source.step.id)
  }
  if categories.allSatisfy({ $0 == .done }) {
    return .done
  }
  if categories.contains(where: { $0 == .upNext }) {
    // If any is upNext, check if all deps are done
    let allDepsDone = sources.allSatisfy { source in
      let prereqs = source.viewModel.getPrerequisiteStepKeys(
        recipeID: source.recipeID, stepID: source.step.id)
      return prereqs.allSatisfy { source.viewModel.completedSteps.contains($0) }
    }
    if allDepsDone {
      return .upNext
    }
  }
  return .forLater
}
