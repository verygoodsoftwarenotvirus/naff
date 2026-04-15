//
//  RecipePerformanceContentView+Helpers.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

// swiftlint:disable file_length

import SwiftProtobuf
import SwiftUI

// MARK: - Helper Types

enum DiscreteQuantityScaling {
  static func scaled(
    _ quantity: Common_Uint32RangeWithOptionalMax,
    scale: Float
  ) -> Common_Uint32RangeWithOptionalMax {
    var scaledQuantity = quantity
    scaledQuantity.min = quantity.min
    if quantity.hasMax {
      scaledQuantity.max = scale > 1 ? scaledMax(quantity.max, scale: scale) : quantity.max
    }
    return scaledQuantity
  }

  static func scaled(
    _ quantity: Common_Uint16RangeWithOptionalMax,
    scale: Float
  ) -> Common_Uint16RangeWithOptionalMax {
    var scaledQuantity = quantity
    scaledQuantity.min = quantity.min
    if quantity.hasMax {
      scaledQuantity.max = scale > 1 ? scaledMax(quantity.max, scale: scale) : quantity.max
    }
    return scaledQuantity
  }

  static func scaledMax(_ value: UInt32, scale: Float) -> UInt32 {
    let scaled = ceil(Double(value) * Double(scale))
    if scaled >= Double(UInt32.max) {
      return UInt32.max
    }
    return UInt32(scaled)
  }

  static func scaledMax(_ value: UInt16, scale: Float) -> UInt16 {
    let scaled = ceil(Double(value) * Double(scale))
    if scaled >= Double(UInt16.max) {
      return UInt16.max
    }
    return UInt16(scaled)
  }
}

struct StepItem {
  let name: String
  let isProduct: Bool
  /// Display label for prerequisite step (e.g. "Drain", "Add (2nd)") for "(from X)" text.
  let prerequisiteStepLabel: String?
  /// Step ID for tappable navigation to prerequisite step.
  let prerequisiteStepID: String?
  let prerequisiteCompleted: Bool
}

struct InstrumentVesselItem: Identifiable {
  let id: String
  let name: String
  let type: ItemType

  enum ItemType {
    case instrument
    case vessel
  }
}

struct AggregatedInstrumentVessel: Identifiable {
  let itemID: String
  let name: String
  let type: InstrumentVesselItem.ItemType
  let sourceRecipeID: String?
  let sourceRecipeName: String?

  private var totalMin: UInt32 = 0
  private var totalMax: UInt32?
  private var hasAnyQuantity: Bool = false

  init(
    itemID: String, name: String, type: InstrumentVesselItem.ItemType,
    sourceRecipeID: String? = nil,
    sourceRecipeName: String? = nil
  ) {
    self.itemID = itemID
    self.name = name
    self.type = type
    self.sourceRecipeID = sourceRecipeID
    self.sourceRecipeName = sourceRecipeName
  }

  mutating func addQuantity(_ quantity: Common_Uint32RangeWithOptionalMax) {
    hasAnyQuantity = true
    totalMin += quantity.min

    if quantity.hasMax {
      if let currentMax = totalMax {
        totalMax = currentMax + quantity.max
      } else {
        totalMax = quantity.max
      }
    } else {
      // If any quantity doesn't have a max, the total doesn't have a max
      totalMax = nil
    }
  }

  mutating func addQuantity(_ quantity: Common_Uint16RangeWithOptionalMax) {
    hasAnyQuantity = true
    totalMin += quantity.min

    if quantity.hasMax {
      if let currentMax = totalMax {
        totalMax = currentMax + quantity.max
      } else {
        totalMax = quantity.max
      }
    } else {
      // If any quantity doesn't have a max, the total doesn't have a max
      totalMax = nil
    }
  }

  var quantityText: String? {
    quantityText(scale: 1.0)
  }

  func quantityText(scale: Float) -> String? {
    guard hasAnyQuantity else { return nil }

    // Instruments/vessels are discrete items. Keep the baseline minimum required quantity
    // and only expand the possible upper bound when scaling up.
    let scaledMin = totalMin
    let scaledMax = totalMax.map { maxQuantity -> UInt32 in
      guard scale > 1 else { return maxQuantity }
      return DiscreteQuantityScaling.scaledMax(maxQuantity, scale: scale)
    }

    if let max = scaledMax {
      if scaledMin == max {
        return "\(scaledMin)"
      } else {
        return "\(scaledMin) - \(max)"
      }
    } else {
      // No max: treat min as the authoritative single value (not unbounded)
      return "\(scaledMin)"
    }
  }

  var id: String {
    itemID
  }
}

struct AggregatedIngredient: Identifiable {
  let ingredientID: String
  let name: String
  let quantityNotes: String
  let measurementUnit: Mealplanning_ValidMeasurementUnit?
  let sourceRecipeID: String?
  let sourceRecipeName: String?

  private var totalMin: Float = 0
  private var totalMax: Float?
  private var hasAnyQuantity: Bool = false

  init(
    ingredientID: String, name: String, quantityNotes: String,
    measurementUnit: Mealplanning_ValidMeasurementUnit?,
    sourceRecipeID: String? = nil,
    sourceRecipeName: String? = nil
  ) {
    self.ingredientID = ingredientID
    self.name = name
    self.quantityNotes = quantityNotes
    self.measurementUnit = measurementUnit
    self.sourceRecipeID = sourceRecipeID
    self.sourceRecipeName = sourceRecipeName
  }

  mutating func addQuantity(_ quantity: Common_Float32RangeWithOptionalMax) {
    hasAnyQuantity = true
    totalMin += quantity.min

    if quantity.hasMax {
      if let currentMax = totalMax {
        totalMax = currentMax + quantity.max
      } else {
        totalMax = quantity.max
      }
    } else {
      // If any quantity doesn't have a max, the total doesn't have a max
      totalMax = nil
    }
  }

  var quantityText: String? {
    quantityText(scale: 1.0)
  }

  func quantityText(scale: Float) -> String? {
    guard hasAnyQuantity else { return nil }

    let scaledMin = totalMin * scale
    let scaledMax = totalMax.map { $0 * scale }

    let unitName = MeasurementUnitFormatter.displayName(
      for: scaledMin,
      unit: measurementUnit
    )
    let unit = unitName.isEmpty ? "" : " \(unitName)"

    // Format numbers - use fewer decimals for whole numbers
    let formatMin = scaledMin.truncatingRemainder(dividingBy: 1) == 0 ? "%.0f" : "%.2f"
    let formatMax =
      scaledMax.map { $0.truncatingRemainder(dividingBy: 1) == 0 ? "%.0f" : "%.2f" } ?? "%.2f"

    if let max = scaledMax {
      if scaledMin == max {
        return String(format: "\(formatMin)%@", scaledMin, unit).trimmingCharacters(
          in: .whitespaces)
      } else {
        return String(format: "\(formatMin) - \(formatMax)%@", scaledMin, max, unit)
          .trimmingCharacters(in: .whitespaces)
      }
    } else {
      // No max: treat min as the authoritative single value (not unbounded)
      return String(format: "\(formatMin)%@", scaledMin, unit).trimmingCharacters(in: .whitespaces)
    }
  }

  var id: String {
    ingredientID
  }
}

// MARK: - Option Group Types

struct IngredientOption: Identifiable {
  let id: String  // Combination of stepID, index, and optionIndex
  let ingredient: Mealplanning_RecipeStepIngredient
  let optionIndex: UInt32
  let aggregated: AggregatedIngredient
}

struct InstrumentOption: Identifiable {
  let id: String  // Combination of stepID, index, and optionIndex
  let instrument: Mealplanning_RecipeStepInstrument
  let optionIndex: UInt32
  let aggregated: AggregatedInstrumentVessel
}

struct VesselOption: Identifiable {
  let id: String  // Combination of stepID, index, and optionIndex
  let vessel: Mealplanning_RecipeStepVessel
  let optionIndex: UInt32
  let aggregated: AggregatedInstrumentVessel
}

struct OptionGroupAggregate: Identifiable {
  let id: String  // Combination of recipeID, stepID, and index
  let recipeID: String
  let stepID: String
  let stepIndex: Int
  let index: UInt32  // The option group index
  let options: [IngredientOption]
  var selectedOptionIndex: UInt32?  // User's selection (nil = not selected, uses default)
  let sourceRecipeID: String?
  let sourceRecipeName: String?
}

struct InstrumentOptionGroupAggregate: Identifiable {
  let id: String
  let recipeID: String
  let stepID: String
  let stepIndex: Int
  let index: UInt32
  let options: [InstrumentOption]
  var selectedOptionIndex: UInt32?
  let sourceRecipeID: String?
  let sourceRecipeName: String?
}

struct VesselOptionGroupAggregate: Identifiable {
  let id: String
  let recipeID: String
  let stepID: String
  let stepIndex: Int
  let index: UInt32
  let options: [VesselOption]
  var selectedOptionIndex: UInt32?
  let sourceRecipeID: String?
  let sourceRecipeName: String?
}

// MARK: - Helper Functions

func formatStepIngredientDisplay(
  _ ingredient: Mealplanning_RecipeStepIngredient,
  scale: Float,
  breakdownSuffix: String? = nil
) -> String {
  var aggregated = AggregatedIngredient(
    ingredientID: ingredient.hasIngredient ? ingredient.ingredient.id : ingredient.id,
    name: ingredient.name,
    quantityNotes: ingredient.quantityNotes,
    measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil
  )

  if ingredient.hasQuantity {
    aggregated.addQuantity(ingredient.quantity)
  }

  var result: String
  if let quantityText = aggregated.quantityText(scale: scale) {
    result = "\(quantityText) \(ingredient.name)"
  } else {
    result = ingredient.name
  }

  if !ingredient.quantityNotes.isEmpty {
    result += " (\(ingredient.quantityNotes))"
  }

  if let suffix = breakdownSuffix, !suffix.isEmpty {
    result += " (\(suffix))"
  }
  return result
}

extension RecipePerformanceContentView {
  // swiftlint:disable:next cyclomatic_complexity function_body_length
  func getAggregatedInstrumentsAndVessels(
    from recipe: Mealplanning_Recipe,
    selectedInstrumentOptions: [String: UInt32] = [:],
    selectedVesselOptions: [String: UInt32] = [:],
    mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]? = nil,
    scale: Float = 1.0
      // swiftlint:disable:next large_tuple
  ) -> (
    regular: [AggregatedInstrumentVessel],
    instrumentOptionGroups: [InstrumentOptionGroupAggregate],
    vesselOptionGroups: [VesselOptionGroupAggregate]
  ) {
    var regularAggregated: [String: AggregatedInstrumentVessel] = [:]
    var instrumentOptionGroups: [InstrumentOptionGroupAggregate] = []
    var vesselOptionGroups: [VesselOptionGroupAggregate] = []

    // Process main recipe
    for (stepIndex, step) in recipe.steps.enumerated() {
      // Process instruments
      let (regularInstruments, instrumentGroups) = groupInstrumentsByOptions(
        Array(step.instruments),
        stepID: step.id,
        stepIndex: stepIndex,
        recipeID: recipe.id
      )

      for instrument in regularInstruments where instrument.hasInstrument {
        let validInstrument = instrument.instrument
        if validInstrument.displayInSummaryLists {
          let itemID = validInstrument.id
          if !itemID.isEmpty {
            if regularAggregated[itemID] == nil {
              regularAggregated[itemID] = AggregatedInstrumentVessel(
                itemID: itemID,
                name: instrument.name,
                type: .instrument
              )
            }

            if instrument.hasQuantity, var current = regularAggregated[itemID] {
              current.addQuantity(instrument.quantity)
              regularAggregated[itemID] = current
            }
          }
        }
      }

      // Filter instrument option groups to only include those with displayInSummaryLists
      let filteredInstrumentGroups = instrumentGroups.filter { group in
        group.options.contains { option in
          option.instrument.hasInstrument && option.instrument.instrument.displayInSummaryLists
        }
      }

      // Add selected options from instrument groups to regular aggregated list
      for group in filteredInstrumentGroups {
        // Check meal plan selections first, then user selections
        // Only add to regular aggregated if a selection has been made
        let selectedIndex: UInt32?
        if let selections = mealPlanSelections,
          let selection = selections.first(where: { sel in
            sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
              && sel.selectionType == .instrument
          })
        {
          selectedIndex = selection.selectedOptionIndex
        } else {
          // Check user selections - use sentinel value if none
          let userSelection = selectedInstrumentOptions[group.id]
          selectedIndex = userSelection == UInt32.max ? nil : userSelection
        }

        // Only add to regular aggregated if a selection has been made
        if let selectedIndex = selectedIndex,
          let selectedOption = group.options.first(where: { $0.optionIndex == selectedIndex }),
          selectedOption.instrument.hasInstrument,
          selectedOption.instrument.instrument.displayInSummaryLists
        {
          let itemID = selectedOption.instrument.instrument.id
          if !itemID.isEmpty {
            if regularAggregated[itemID] == nil {
              regularAggregated[itemID] = AggregatedInstrumentVessel(
                itemID: itemID,
                name: selectedOption.instrument.name,
                type: .instrument
              )
            }
            if selectedOption.instrument.hasQuantity, var current = regularAggregated[itemID] {
              current.addQuantity(selectedOption.instrument.quantity)
              regularAggregated[itemID] = current
            }
          }
        }
      }

      instrumentOptionGroups.append(contentsOf: filteredInstrumentGroups)

      // Process vessels
      let (regularVessels, vesselGroups) = groupVesselsByOptions(
        Array(step.vessels),
        stepID: step.id,
        stepIndex: stepIndex,
        recipeID: recipe.id
      )

      for vessel in regularVessels where vessel.hasVessel {
        let validVessel = vessel.vessel
        if validVessel.displayInSummaryLists {
          let itemID = validVessel.id
          if !itemID.isEmpty {
            if regularAggregated[itemID] == nil {
              regularAggregated[itemID] = AggregatedInstrumentVessel(
                itemID: itemID,
                name: vessel.name,
                type: .vessel
              )
            }

            if vessel.hasQuantity, var current = regularAggregated[itemID] {
              current.addQuantity(vessel.quantity)
              regularAggregated[itemID] = current
            }
          }
        }
      }

      // Filter vessel option groups to only include those with displayInSummaryLists
      let filteredVesselGroups = vesselGroups.filter { group in
        group.options.contains { option in
          option.vessel.hasVessel && option.vessel.vessel.displayInSummaryLists
        }
      }

      // Add selected options from vessel groups to regular aggregated list
      for group in filteredVesselGroups {
        // Check meal plan selections first, then user selections
        // Only add to regular aggregated if a selection has been made
        let selectedIndex: UInt32?
        if let selections = mealPlanSelections,
          let selection = selections.first(where: { sel in
            sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
              && sel.selectionType == .vessel
          })
        {
          selectedIndex = selection.selectedOptionIndex
        } else {
          // Check user selections - use sentinel value if none
          let userSelection = selectedVesselOptions[group.id]
          selectedIndex = userSelection == UInt32.max ? nil : userSelection
        }

        // Only add to regular aggregated if a selection has been made
        if let selectedIndex = selectedIndex,
          let selectedOption = group.options.first(where: { $0.optionIndex == selectedIndex }),
          selectedOption.vessel.hasVessel,
          selectedOption.vessel.vessel.displayInSummaryLists
        {
          let itemID = selectedOption.vessel.vessel.id
          if !itemID.isEmpty {
            if regularAggregated[itemID] == nil {
              regularAggregated[itemID] = AggregatedInstrumentVessel(
                itemID: itemID,
                name: selectedOption.vessel.name,
                type: .vessel
              )
            }
            if selectedOption.vessel.hasQuantity, var current = regularAggregated[itemID] {
              current.addQuantity(selectedOption.vessel.quantity)
              regularAggregated[itemID] = current
            }
          }
        }
      }

      vesselOptionGroups.append(contentsOf: filteredVesselGroups)
    }

    // Process associated recipes
    for associatedRecipe in recipe.associatedRecipes {
      for (stepIndex, step) in associatedRecipe.steps.enumerated() {
        let (regularInstruments, instrumentGroups) = groupInstrumentsByOptions(
          Array(step.instruments),
          stepID: step.id,
          stepIndex: stepIndex,
          recipeID: associatedRecipe.id
        )

        for instrument in regularInstruments where instrument.hasInstrument {
          let validInstrument = instrument.instrument
          if validInstrument.displayInSummaryLists {
            let itemID = validInstrument.id
            if !itemID.isEmpty {
              if regularAggregated[itemID] == nil {
                regularAggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: instrument.name,
                  type: .instrument,
                  sourceRecipeID: associatedRecipe.id,
                  sourceRecipeName: associatedRecipe.name
                )
              }

              if instrument.hasQuantity, var current = regularAggregated[itemID] {
                current.addQuantity(instrument.quantity)
                regularAggregated[itemID] = current
              }
            }
          }
        }

        var instrumentGroupsWithSource = instrumentGroups
        for index in 0..<instrumentGroupsWithSource.count {
          instrumentGroupsWithSource[index] = InstrumentOptionGroupAggregate(
            id: instrumentGroupsWithSource[index].id,
            recipeID: instrumentGroupsWithSource[index].recipeID,
            stepID: instrumentGroupsWithSource[index].stepID,
            stepIndex: instrumentGroupsWithSource[index].stepIndex,
            index: instrumentGroupsWithSource[index].index,
            options: instrumentGroupsWithSource[index].options,
            selectedOptionIndex: instrumentGroupsWithSource[index].selectedOptionIndex,
            sourceRecipeID: associatedRecipe.id,
            sourceRecipeName: associatedRecipe.name
          )
        }

        // Add selected options from instrument groups to regular aggregated list
        for group in instrumentGroupsWithSource {
          // Check meal plan selections first, then user selections, then default to optionIndex 0
          let selectedIndex: UInt32
          if let selections = mealPlanSelections,
            let selection = selections.first(where: { sel in
              sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
                && sel.selectionType == .instrument
            })
          {
            selectedIndex = selection.selectedOptionIndex
          } else {
            selectedIndex =
              selectedInstrumentOptions[group.id]
              ?? (group.options.first(where: { $0.optionIndex == 0 })?.optionIndex ?? group.options
                .first?.optionIndex ?? 0)
          }
          if let selectedOption = group.options.first(where: { $0.optionIndex == selectedIndex }),
            selectedOption.instrument.hasInstrument,
            selectedOption.instrument.instrument.displayInSummaryLists
          {
            let itemID = selectedOption.instrument.instrument.id
            if !itemID.isEmpty {
              if regularAggregated[itemID] == nil {
                regularAggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: selectedOption.instrument.name,
                  type: .instrument,
                  sourceRecipeID: associatedRecipe.id,
                  sourceRecipeName: associatedRecipe.name
                )
              }
              if selectedOption.instrument.hasQuantity, var current = regularAggregated[itemID] {
                current.addQuantity(selectedOption.instrument.quantity)
                regularAggregated[itemID] = current
              }
            }
          }
        }

        instrumentOptionGroups.append(contentsOf: instrumentGroupsWithSource)

        let (regularVessels, vesselGroups) = groupVesselsByOptions(
          Array(step.vessels),
          stepID: step.id,
          stepIndex: stepIndex,
          recipeID: associatedRecipe.id
        )

        for vessel in regularVessels where vessel.hasVessel {
          let validVessel = vessel.vessel
          if validVessel.displayInSummaryLists {
            let itemID = validVessel.id
            if !itemID.isEmpty {
              if regularAggregated[itemID] == nil {
                regularAggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: vessel.name,
                  type: .vessel,
                  sourceRecipeID: associatedRecipe.id,
                  sourceRecipeName: associatedRecipe.name
                )
              }

              if vessel.hasQuantity, var current = regularAggregated[itemID] {
                current.addQuantity(vessel.quantity)
                regularAggregated[itemID] = current
              }
            }
          }
        }

        var vesselGroupsWithSource = vesselGroups
        for index in 0..<vesselGroupsWithSource.count {
          vesselGroupsWithSource[index] = VesselOptionGroupAggregate(
            id: vesselGroupsWithSource[index].id,
            recipeID: vesselGroupsWithSource[index].recipeID,
            stepID: vesselGroupsWithSource[index].stepID,
            stepIndex: vesselGroupsWithSource[index].stepIndex,
            index: vesselGroupsWithSource[index].index,
            options: vesselGroupsWithSource[index].options,
            selectedOptionIndex: vesselGroupsWithSource[index].selectedOptionIndex,
            sourceRecipeID: associatedRecipe.id,
            sourceRecipeName: associatedRecipe.name
          )
        }

        // Add selected options from vessel groups to regular aggregated list
        for group in vesselGroupsWithSource {
          // Check meal plan selections first, then user selections, then default to optionIndex 0
          let selectedIndex: UInt32
          if let selections = mealPlanSelections,
            let selection = selections.first(where: { sel in
              sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
                && sel.selectionType == .vessel
            })
          {
            selectedIndex = selection.selectedOptionIndex
          } else {
            selectedIndex =
              selectedVesselOptions[group.id]
              ?? (group.options.first(where: { $0.optionIndex == 0 })?.optionIndex ?? group.options
                .first?.optionIndex ?? 0)
          }
          if let selectedOption = group.options.first(where: { $0.optionIndex == selectedIndex }),
            selectedOption.vessel.hasVessel,
            selectedOption.vessel.vessel.displayInSummaryLists
          {
            let itemID = selectedOption.vessel.vessel.id
            if !itemID.isEmpty {
              if regularAggregated[itemID] == nil {
                regularAggregated[itemID] = AggregatedInstrumentVessel(
                  itemID: itemID,
                  name: selectedOption.vessel.name,
                  type: .vessel,
                  sourceRecipeID: associatedRecipe.id,
                  sourceRecipeName: associatedRecipe.name
                )
              }
              if selectedOption.vessel.hasQuantity, var current = regularAggregated[itemID] {
                current.addQuantity(selectedOption.vessel.quantity)
                regularAggregated[itemID] = current
              }
            }
          }
        }

        vesselOptionGroups.append(contentsOf: vesselGroupsWithSource)
      }
    }

    return (
      regular: Array(regularAggregated.values).sorted { $0.name < $1.name },
      instrumentOptionGroups: instrumentOptionGroups,
      vesselOptionGroups: vesselOptionGroups
    )
  }

  // swiftlint:disable:next cyclomatic_complexity
  func getAggregatedIngredients(
    from recipe: Mealplanning_Recipe,
    selectedIngredientOptions: [String: UInt32] = [:],
    mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]? = nil,
    scale: Float = 1.0
  ) -> (
    regular: [AggregatedIngredient],
    optionGroups: [OptionGroupAggregate]
  ) {
    var regularAggregated: [String: AggregatedIngredient] = [:]
    var optionGroups: [OptionGroupAggregate] = []

    // Process main recipe
    for (stepIndex, step) in recipe.steps.enumerated() {
      let (regular, groups) = groupIngredientsByOptions(
        Array(step.ingredients),
        stepID: step.id,
        stepIndex: stepIndex,
        recipeID: recipe.id
      )

      // Process regular ingredients
      for ingredient in regular where ingredient.hasIngredient {
        let validIngredient = ingredient.ingredient
        let key = validIngredient.id
        if !key.isEmpty {
          if regularAggregated[key] == nil {
            regularAggregated[key] = AggregatedIngredient(
              ingredientID: key,
              name: ingredient.name,
              quantityNotes: ingredient.quantityNotes,
              measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil
            )
          }

          if ingredient.hasQuantity, var current = regularAggregated[key] {
            current.addQuantity(ingredient.quantity)
            regularAggregated[key] = current
          }
        }
      }

      // Add selected options from ingredient groups to regular aggregated list
      for group in groups {
        // Check meal plan selections first, then user selections, then default to optionIndex 0
        let selectedIndex: UInt32
        if let selections = mealPlanSelections,
          let selection = selections.first(where: { sel in
            sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
              && sel.ingredientIndex == group.index && sel.selectionType == .ingredient
          })
        {
          selectedIndex = selection.selectedOptionIndex
        } else {
          selectedIndex =
            selectedIngredientOptions[group.id]
            ?? (group.options.first(where: { $0.optionIndex == 0 })?.optionIndex ?? group.options
              .first?.optionIndex ?? 0)
        }
        if let selectedOption = group.options.first(where: { $0.optionIndex == selectedIndex }),
          selectedOption.ingredient.hasIngredient
        {
          let key = selectedOption.ingredient.ingredient.id
          if !key.isEmpty {
            if regularAggregated[key] == nil {
              regularAggregated[key] = AggregatedIngredient(
                ingredientID: key,
                name: selectedOption.ingredient.name,
                quantityNotes: selectedOption.ingredient.quantityNotes,
                measurementUnit: selectedOption.ingredient.hasMeasurementUnit
                  ? selectedOption.ingredient.measurementUnit : nil
              )
            }
            if selectedOption.ingredient.hasQuantity, var current = regularAggregated[key] {
              current.addQuantity(selectedOption.ingredient.quantity)
              regularAggregated[key] = current
            }
          }
        }
      }

      // Add option groups
      optionGroups.append(contentsOf: groups)
    }

    // Process associated recipes
    for associatedRecipe in recipe.associatedRecipes {
      for (stepIndex, step) in associatedRecipe.steps.enumerated() {
        let (regular, groups) = groupIngredientsByOptions(
          Array(step.ingredients),
          stepID: step.id,
          stepIndex: stepIndex,
          recipeID: associatedRecipe.id
        )

        // Process regular ingredients from associated recipe
        for ingredient in regular where ingredient.hasIngredient {
          let validIngredient = ingredient.ingredient
          let key = validIngredient.id
          if !key.isEmpty {
            if regularAggregated[key] == nil {
              regularAggregated[key] = AggregatedIngredient(
                ingredientID: key,
                name: ingredient.name,
                quantityNotes: ingredient.quantityNotes,
                measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil,
                sourceRecipeID: associatedRecipe.id,
                sourceRecipeName: associatedRecipe.name
              )
            }

            if ingredient.hasQuantity, var current = regularAggregated[key] {
              current.addQuantity(ingredient.quantity)
              regularAggregated[key] = current
            }
          }
        }

        // Add option groups from associated recipe with source info
        var groupsWithSource = groups
        for index in 0..<groupsWithSource.count {
          groupsWithSource[index] = OptionGroupAggregate(
            id: groupsWithSource[index].id,
            recipeID: groupsWithSource[index].recipeID,
            stepID: groupsWithSource[index].stepID,
            stepIndex: groupsWithSource[index].stepIndex,
            index: groupsWithSource[index].index,
            options: groupsWithSource[index].options,
            selectedOptionIndex: groupsWithSource[index].selectedOptionIndex,
            sourceRecipeID: associatedRecipe.id,
            sourceRecipeName: associatedRecipe.name
          )
        }

        // Add selected options from ingredient groups to regular aggregated list
        for group in groupsWithSource {
          // Check meal plan selections first, then user selections, then default to optionIndex 0
          let selectedIndex: UInt32
          if let selections = mealPlanSelections,
            let selection = selections.first(where: { sel in
              sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
                && sel.ingredientIndex == group.index && sel.selectionType == .ingredient
            })
          {
            selectedIndex = selection.selectedOptionIndex
          } else {
            selectedIndex =
              selectedIngredientOptions[group.id]
              ?? (group.options.first(where: { $0.optionIndex == 0 })?.optionIndex ?? group.options
                .first?.optionIndex ?? 0)
          }
          if let selectedOption = group.options.first(where: { $0.optionIndex == selectedIndex }),
            selectedOption.ingredient.hasIngredient
          {
            let key = selectedOption.ingredient.ingredient.id
            if !key.isEmpty {
              if regularAggregated[key] == nil {
                regularAggregated[key] = AggregatedIngredient(
                  ingredientID: key,
                  name: selectedOption.ingredient.name,
                  quantityNotes: selectedOption.ingredient.quantityNotes,
                  measurementUnit: selectedOption.ingredient.hasMeasurementUnit
                    ? selectedOption.ingredient.measurementUnit : nil,
                  sourceRecipeID: associatedRecipe.id,
                  sourceRecipeName: associatedRecipe.name
                )
              }
              if selectedOption.ingredient.hasQuantity, var current = regularAggregated[key] {
                current.addQuantity(selectedOption.ingredient.quantity)
                regularAggregated[key] = current
              }
            }
          }
        }

        optionGroups.append(contentsOf: groupsWithSource)
      }
    }

    return (
      regular: Array(regularAggregated.values).sorted { $0.name < $1.name },
      optionGroups: optionGroups
    )
  }

  func formatStepTitle(step: Mealplanning_RecipeStep, viewModel: PerformRecipeViewModel) -> String {
    if step.hasPreparation && !step.preparation.name.isEmpty {
      return step.preparation.name
    }
    return "Step \(Int(step.index) + 1)"
  }

  // MARK: - Option Grouping Functions

  /// Groups ingredients by option groups. Returns regular ingredients and option groups separately.
  func groupIngredientsByOptions(
    _ ingredients: [Mealplanning_RecipeStepIngredient],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepIngredient],
    optionGroups: [OptionGroupAggregate]
  ) {
    var regular: [Mealplanning_RecipeStepIngredient] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]

    for ingredient in ingredients {
      let index = ingredient.index
      // Check if there are other ingredients with the same index (option group)
      let hasOptions = ingredients.contains { other in
        other.id != ingredient.id && other.index == index
      }

      if hasOptions {
        if optionGroupsByIndex[index] == nil {
          optionGroupsByIndex[index] = []
        }
        optionGroupsByIndex[index]?.append(ingredient)
      } else {
        regular.append(ingredient)
      }
    }

    // Convert option groups to OptionGroupAggregate
    var optionGroups: [OptionGroupAggregate] = []
    for (index, groupIngredients) in optionGroupsByIndex {
      // Sort by optionIndex
      let sorted = groupIngredients.sorted { lhs, rhs in
        let lhsIndex = lhs.optionIndex
        let rhsIndex = rhs.optionIndex
        return lhsIndex < rhsIndex
      }

      // Create aggregated ingredients for each option
      var options: [IngredientOption] = []
      for ingredient in sorted {
        let optionIndex = ingredient.optionIndex
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

        // Create aggregated ingredient for this option
        var aggregated = AggregatedIngredient(
          ingredientID: ingredient.hasIngredient ? ingredient.ingredient.id : ingredient.id,
          name: ingredient.name,
          quantityNotes: ingredient.quantityNotes,
          measurementUnit: ingredient.hasMeasurementUnit ? ingredient.measurementUnit : nil
        )

        if ingredient.hasQuantity {
          aggregated.addQuantity(ingredient.quantity)
        }

        options.append(
          IngredientOption(
            id: optionID,
            ingredient: ingredient,
            optionIndex: optionIndex,
            aggregated: aggregated
          )
        )
      }

      let groupID = "\(recipeID)-\(stepID)-\(index)"
      optionGroups.append(
        OptionGroupAggregate(
          id: groupID,
          recipeID: recipeID,
          stepID: stepID,
          stepIndex: stepIndex,
          index: index,
          options: options,
          selectedOptionIndex: nil,  // Will be set by view state
          sourceRecipeID: nil,
          sourceRecipeName: nil
        )
      )
    }

    // Sort option groups by index
    optionGroups.sort { $0.index < $1.index }

    return (regular: regular, optionGroups: optionGroups)
  }

  /// Groups instruments by option groups
  func groupInstrumentsByOptions(
    _ instruments: [Mealplanning_RecipeStepInstrument],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepInstrument],
    optionGroups: [InstrumentOptionGroupAggregate]
  ) {
    var regular: [Mealplanning_RecipeStepInstrument] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepInstrument]] = [:]

    for instrument in instruments {
      // Index 0 typically means not in an option group
      if instrument.index != 0 {
        let index = instrument.index
        let hasOptions = instruments.contains { other in
          other.id != instrument.id && other.index != 0 && other.index == index
        }

        if hasOptions {
          if optionGroupsByIndex[index] == nil {
            optionGroupsByIndex[index] = []
          }
          optionGroupsByIndex[index]?.append(instrument)
        } else {
          regular.append(instrument)
        }
      } else {
        regular.append(instrument)
      }
    }

    var optionGroups: [InstrumentOptionGroupAggregate] = []
    for (index, groupInstruments) in optionGroupsByIndex {
      let sorted = groupInstruments.sorted { lhs, rhs in
        let lhsIndex = lhs.optionIndex
        let rhsIndex = rhs.optionIndex
        return lhsIndex < rhsIndex
      }

      var options: [InstrumentOption] = []
      for instrument in sorted {
        let optionIndex = instrument.optionIndex
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

        var aggregated = AggregatedInstrumentVessel(
          itemID: instrument.hasInstrument ? instrument.instrument.id : instrument.id,
          name: instrument.name,
          type: .instrument
        )

        if instrument.hasQuantity {
          aggregated.addQuantity(instrument.quantity)
        }

        options.append(
          InstrumentOption(
            id: optionID,
            instrument: instrument,
            optionIndex: optionIndex,
            aggregated: aggregated
          )
        )
      }

      let groupID = "\(recipeID)-\(stepID)-\(index)"
      optionGroups.append(
        InstrumentOptionGroupAggregate(
          id: groupID,
          recipeID: recipeID,
          stepID: stepID,
          stepIndex: stepIndex,
          index: index,
          options: options,
          selectedOptionIndex: nil,
          sourceRecipeID: nil,
          sourceRecipeName: nil
        )
      )
    }

    optionGroups.sort { $0.index < $1.index }

    return (regular: regular, optionGroups: optionGroups)
  }

  /// Groups vessels by option groups
  func groupVesselsByOptions(
    _ vessels: [Mealplanning_RecipeStepVessel],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepVessel],
    optionGroups: [VesselOptionGroupAggregate]
  ) {
    var regular: [Mealplanning_RecipeStepVessel] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepVessel]] = [:]

    for vessel in vessels {
      // Index 0 typically means not in an option group
      if vessel.index != 0 {
        let index = vessel.index
        let hasOptions = vessels.contains { other in
          other.id != vessel.id && other.index != 0 && other.index == index
        }

        if hasOptions {
          if optionGroupsByIndex[index] == nil {
            optionGroupsByIndex[index] = []
          }
          optionGroupsByIndex[index]?.append(vessel)
        } else {
          regular.append(vessel)
        }
      } else {
        regular.append(vessel)
      }
    }

    var optionGroups: [VesselOptionGroupAggregate] = []
    for (index, groupVessels) in optionGroupsByIndex {
      let sorted = groupVessels.sorted { lhs, rhs in
        let lhsIndex = lhs.optionIndex
        let rhsIndex = rhs.optionIndex
        return lhsIndex < rhsIndex
      }

      var options: [VesselOption] = []
      for vessel in sorted {
        let optionIndex = vessel.optionIndex
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

        var aggregated = AggregatedInstrumentVessel(
          itemID: vessel.hasVessel ? vessel.vessel.id : vessel.id,
          name: vessel.name,
          type: .vessel
        )

        if vessel.hasQuantity {
          aggregated.addQuantity(vessel.quantity)
        }

        options.append(
          VesselOption(
            id: optionID,
            vessel: vessel,
            optionIndex: optionIndex,
            aggregated: aggregated
          )
        )
      }

      let groupID = "\(recipeID)-\(stepID)-\(index)"
      optionGroups.append(
        VesselOptionGroupAggregate(
          id: groupID,
          recipeID: recipeID,
          stepID: stepID,
          stepIndex: stepIndex,
          index: index,
          options: options,
          selectedOptionIndex: nil,
          sourceRecipeID: nil,
          sourceRecipeName: nil
        )
      )
    }

    optionGroups.sort { $0.index < $1.index }

    return (regular: regular, optionGroups: optionGroups)
  }
}
