//
//  RecipePerformanceContentView+Views.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

// swiftlint:disable file_length

import SwiftProtobuf
import SwiftUI

// MARK: - Timer Condition Row

struct TimerConditionRow: View {
  let viewModel: PerformRecipeViewModel
  let recipeID: String
  let stepID: String
  let minSeconds: UInt32
  let maxSeconds: UInt32?
  var canStartStep: Bool = true
  var timerElapsedSeconds: TimeInterval?
  var onSkipTimerOverride: (() -> Void)?
  var canSkipTimerOverride: Bool?
  var timerMinSeconds: UInt32?
  var timerMaxSeconds: UInt32?

  private var isTimerCompleted: Bool {
    viewModel.isTimerConditionCompleted(recipeID: recipeID, stepID: stepID)
  }

  private var isTimerActive: Bool {
    viewModel.isStepTimerActive(recipeID: recipeID, stepID: stepID)
  }

  private var elapsed: TimeInterval {
    timerElapsedSeconds ?? viewModel.stepTimerElapsedSeconds(recipeID: recipeID, stepID: stepID)
      ?? 0
  }

  private var effectiveMin: UInt32 {
    timerMinSeconds ?? viewModel.stepTimerMinSeconds(recipeID: recipeID, stepID: stepID)
      ?? minSeconds
  }

  private var effectiveMax: UInt32? {
    timerMaxSeconds ?? viewModel.stepTimerMaxSeconds(recipeID: recipeID, stepID: stepID)
      ?? maxSeconds
  }

  private var canMarkDone: Bool {
    canSkipTimerOverride ?? viewModel.canSkipStepTimer(recipeID: recipeID, stepID: stepID)
  }

  private var durationLabel: String {
    if let max = maxSeconds, max > minSeconds {
      return RecipeTimeEstimation.format(minSeconds: minSeconds, maxSeconds: max)
    }
    return RecipeTimeEstimation.format(minSeconds: minSeconds, maxSeconds: minSeconds)
  }

  var body: some View {
    if isTimerCompleted {
      completedState
    } else if isTimerActive {
      runningState
    } else {
      notStartedState
    }
  }

  private var notStartedState: some View {
    HStack(alignment: .center, spacing: 12) {
      Button {
        guard canStartStep else { return }
        viewModel.startStepTimer(recipeID: recipeID, stepID: stepID)
      } label: {
        HStack(spacing: 6) {
          Image(systemName: "play.fill")
            .font(.subheadline)
          Text("Start")
            .font(.subheadline)
            .fontWeight(.semibold)
        }
        .foregroundColor(canStartStep ? .white : .gray)
        .padding(.horizontal, 14)
        .padding(.vertical, 8)
        .background(canStartStep ? Color.orange : Color.gray.opacity(0.5))
        .cornerRadius(8)
      }
      .buttonStyle(.plain)
      .disabled(!canStartStep)

      Text(durationLabel)
        .font(.caption)
        .foregroundColor(.secondary)

      Spacer()
    }
  }

  private var runningState: some View {
    TimelineView(.periodic(from: .now, by: 1.0)) { _ in
      let displayText: String = {
        RecipeTimeEstimation.formatElapsedWithRange(
          elapsedSeconds: elapsed,
          minSeconds: effectiveMin,
          maxSeconds: effectiveMax
        )
      }()
      let timerColor: Color = canMarkDone ? .green : .orange

      HStack(alignment: .center, spacing: 12) {
        HStack(spacing: 8) {
          Image(systemName: "clock.fill")
            .font(.title3)
            .foregroundColor(timerColor)
          Text(displayText)
            .font(.subheadline)
            .fontWeight(.semibold)
            .foregroundColor(timerColor)
            .monospacedDigit()
        }
        .padding(.horizontal, 10)
        .padding(.vertical, 6)
        .background(timerColor.opacity(0.15))
        .cornerRadius(8)

        Spacer()

        if canMarkDone {
          Button {
            viewModel.completeTimerCondition(recipeID: recipeID, stepID: stepID)
          } label: {
            Text("Done")
              .font(.subheadline)
              .fontWeight(.semibold)
              .foregroundColor(.white)
              .padding(.horizontal, 14)
              .padding(.vertical, 8)
              .background(Color.green)
              .cornerRadius(8)
          }
          .buttonStyle(.plain)
        }

        Button {
          if let onSkip = onSkipTimerOverride {
            onSkip()
          } else {
            viewModel.skipStepTimer(recipeID: recipeID, stepID: stepID)
          }
        } label: {
          Text("Skip")
            .font(.subheadline)
            .fontWeight(.medium)
            .foregroundColor(.secondary)
        }
        .buttonStyle(.plain)

        Button {
          viewModel.cancelStepTimer(recipeID: recipeID, stepID: stepID)
        } label: {
          Text("Cancel")
            .font(.subheadline)
            .fontWeight(.medium)
            .foregroundColor(.secondary)
        }
        .buttonStyle(.plain)
      }
    }
  }

  private var completedState: some View {
    HStack(alignment: .top, spacing: 12) {
      Button {
        viewModel.toggleTimerCondition(recipeID: recipeID, stepID: stepID)
      } label: {
        Image(systemName: "checkmark.circle.fill")
          .font(.title3)
          .foregroundColor(.green)
      }
      .buttonStyle(.plain)

      Text(durationLabel)
        .font(.caption)
        .foregroundColor(.secondary)
        .strikethrough(true)

      Spacer()
    }
  }
}

// MARK: - Step Card View

struct StepCardView: View {
  let step: Mealplanning_RecipeStep
  let index: Int
  let viewModel: PerformRecipeViewModel
  let formatStepTitle: (Mealplanning_RecipeStep, PerformRecipeViewModel) -> String
  let recipeID: String
  var mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]?
  var isAssociatedRecipeStep: Bool = false
  var associatedRecipeName: String?
  var highlightedStepIDs: Set<String>?
  var selectedIngredientOptions: [String: UInt32] = [:]
  var selectedInstrumentOptions: [String: UInt32] = [:]
  var selectedVesselOptions: [String: UInt32] = [:]
  var scale: Float = 1.0
  /// Override for merged meal steps (multiple recipes)
  var isCompletedOverride: Bool?
  var canCheckOverride: Bool?
  var onToggleOverride: (() -> Void)?
  /// Called when user completes a step (for dismissing wash-hands reminder)
  var onStepCompleted: (() -> Void)?
  /// For merged meal steps: ingredient key -> "2g from Recipe A, 3g from Recipe B"
  var ingredientBreakdownBySource: [String: String]?
  /// When set, show elapsed timer (for merged steps where parent computes from multiple sources)
  var timerElapsedSeconds: TimeInterval?
  var timerDurationSeconds: UInt32?
  /// When set, call to skip timer (for merged steps - skips all sources with active timers)
  var onSkipTimerOverride: (() -> Void)?
  /// When set, whether skip is allowed (for merged steps - any source can skip)
  var canSkipTimerOverride: Bool?
  /// When set, min/max for range display (for merged steps)
  var timerMinSeconds: UInt32?
  var timerMaxSeconds: UInt32?

  @State private var showCompletionConditionsHint = false

  private var isHighlighted: Bool {
    guard let highlightedStepIDs = highlightedStepIDs else { return true }
    return highlightedStepIDs.contains(step.id)
  }

  private var isDimmed: Bool {
    highlightedStepIDs != nil && !isHighlighted
  }

  var body: some View {
    // Use overrides for merged steps, otherwise use viewModel
    let isCompleted: Bool
    let prerequisiteStepKeys: [String]
    let isTimerActive: Bool

    let canCheckStep: Bool
    if let overrideCompleted = isCompletedOverride, let overrideCanCheck = canCheckOverride {
      isCompleted = overrideCompleted
      canCheckStep = overrideCanCheck
      prerequisiteStepKeys = []
      isTimerActive = timerElapsedSeconds != nil
    } else {
      isCompleted = viewModel.isStepCompleted(recipeID: recipeID, stepID: step.id)
      let canStartStep = viewModel.canStartStep(recipeID: recipeID, stepID: step.id)
      canCheckStep = viewModel.canCheckStep(recipeID: recipeID, stepID: step.id)
      isTimerActive = viewModel.isStepTimerActive(recipeID: recipeID, stepID: step.id)
      prerequisiteStepKeys = viewModel.getPrerequisiteStepKeys(recipeID: recipeID, stepID: step.id)
    }

    let hasPrerequisites = !prerequisiteStepKeys.isEmpty
    let allPrerequisitesCompleted = prerequisiteStepKeys.allSatisfy { stepKey in
      let parts = stepKey.split(separator: ":", maxSplits: 1)
      guard parts.count == 2,
        let rid = String(parts[0]) as String?,
        let sid = String(parts[1]) as String?
      else { return true }
      return viewModel.isStepCompleted(recipeID: rid, stepID: sid)
    }
    let completionConditions = step.completionConditions
    let hasCompletionConditions = !completionConditions.isEmpty
    let canSkipTimer =
      isTimerActive
      ? (canSkipTimerOverride
        ?? viewModel.canSkipStepTimer(
          recipeID: recipeID, stepID: step.id))
      : false
    let timerIconColor: Color = canSkipTimer ? .green : .orange

    let hasTimerCondition =
      step.estimatedTimeInSeconds.hasMin && step.estimatedTimeInSeconds.min > 0
    let hasConditionsOrTimer = hasCompletionConditions || hasTimerCondition

    return VStack(alignment: .leading, spacing: 12) {
      // Step header: step number badge, optional checkbox (simple steps only), title
      HStack(alignment: .top, spacing: 12) {
        if !isAssociatedRecipeStep {
          Text("\(index + 1)")
            .font(.caption)
            .fontWeight(.bold)
            .foregroundColor(.secondary)
            .frame(width: 24, height: 24)
            .background(Color(.systemGray5))
            .clipShape(Circle())
        }

        if isCompleted {
          Image(systemName: "checkmark.circle.fill")
            .font(.title2)
            .foregroundColor(.green)
        }

        // Step title with preparation and ingredients
        VStack(alignment: .leading, spacing: 4) {
          HStack {
            Text(formatStepTitle(step, viewModel))
              .font(.headline)
              .foregroundColor(isCompleted ? .secondary : .primary)
              .italic(isCompleted)

            if step.optional {
              Text("(Optional)")
                .font(.caption)
                .foregroundColor(.secondary)
            }

            if !isTimerActive,
              let stepTime = RecipeTimeEstimation.formatStepTime(
                step.estimatedTimeInSeconds)
            {
              Label(stepTime, systemImage: "clock")
                .font(.caption)
                .foregroundColor(.secondary)
            }
          }

          if !step.explicitInstructions.isEmpty {
            Text(step.explicitInstructions)
              .font(.body)
              .foregroundColor(isCompleted ? .secondary : .primary)
              .strikethrough(isCompleted)
          }

          // Prerequisites warning (only show if wash hands is done)
          if viewModel.washHandsCompleted && hasPrerequisites && !allPrerequisitesCompleted {
            HStack(spacing: 4) {
              Image(systemName: "exclamationmark.triangle.fill")
                .font(.caption)
                .foregroundColor(.orange)
              Text(
                "Complete \(prerequisiteStepKeys.compactMap { viewModel.getStepDisplayLabelForStepKey($0) }.joined(separator: ", ")) first"
              )
              .font(.caption)
              .foregroundColor(.orange)
            }
            .padding(.top, 4)
          }
        }

        Spacer()
      }

      // Step details (ingredients, instruments, vessels)
      if !isCompleted || true {  // Show details even when completed
        StepDetailsView(
          step: step,
          viewModel: viewModel,
          stepIndex: index,
          recipeID: recipeID,
          mealPlanSelections: mealPlanSelections,
          selectedIngredientOptions: selectedIngredientOptions,
          selectedInstrumentOptions: selectedInstrumentOptions,
          selectedVesselOptions: selectedVesselOptions,
          scale: scale,
          ingredientBreakdownBySource: ingredientBreakdownBySource
        )
      }

      if hasCompletionConditions || hasTimerCondition {
        completionConditionsSection(
          completionConditions: completionConditions,
          hasTimerCondition: hasTimerCondition,
          timerMinSeconds: hasTimerCondition
            ? (timerMinSeconds ?? viewModel.stepTimerMinSeconds(recipeID: recipeID, stepID: step.id))
            : nil,
          timerMaxSeconds: hasTimerCondition
            ? (timerMaxSeconds ?? viewModel.stepTimerMaxSeconds(recipeID: recipeID, stepID: step.id))
            : nil,
          isHighlighted: showCompletionConditionsHint
        )
      }

      if !(isCompletedOverride ?? false) {
        if isCompleted {
          Button {
            if let onToggle = onToggleOverride {
              onToggle()
            } else {
              viewModel.toggleStep(recipeID: recipeID, stepID: step.id)
            }
          } label: {
            Text("Undo")
              .font(.subheadline)
              .fontWeight(.medium)
              .foregroundColor(.secondary)
          }
          .buttonStyle(.plain)
        } else {
          let canComplete =
            canCheckOverride ?? viewModel.canCompleteStep(recipeID: recipeID, stepID: step.id)
          Button {
            if canComplete {
              if let onToggle = onToggleOverride {
                onToggle()
              } else {
                viewModel.markStepComplete(recipeID: recipeID, stepID: step.id)
              }
              onStepCompleted?()
            } else if hasConditionsOrTimer {
              UIImpactFeedbackGenerator(style: .light).impactOccurred()
              showCompletionConditionsHint = true
              Task { @MainActor in
                try? await Task.sleep(nanoseconds: 2_500_000_000)
                showCompletionConditionsHint = false
              }
            }
          } label: {
            Text("Complete Step")
              .font(.headline)
              .fontWeight(.semibold)
              .foregroundColor(.white)
              .frame(maxWidth: .infinity)
              .padding(.vertical, 12)
              .background(canComplete ? Color.green : Color.gray)
              .cornerRadius(10)
          }
          .buttonStyle(.plain)
          .disabled(!canComplete)
        }
      }
    }
    .padding()
    .background(
      isAssociatedRecipeStep
        ? Color.purple.opacity(0.05)
        : (isTimerActive
          ? (canSkipTimer ? Color.green.opacity(0.08) : Color.orange.opacity(0.08))
          : (isCompleted ? Color(.systemGray6) : Color(.systemBackground)))
    )
    .cornerRadius(12)
    .overlay(
      RoundedRectangle(cornerRadius: 12)
        .stroke(
          stepBorderColor(
            isCompleted: isCompleted,
            isTimerActive: isTimerActive,
            canSkipTimer: canSkipTimer,
            canCompleteStep: canCheckOverride
              ?? viewModel.canCompleteStep(recipeID: recipeID, stepID: step.id),
            hasConditionsOrTimer: hasConditionsOrTimer
          ),
          lineWidth: isHighlighted && highlightedStepIDs != nil ? 2.5 : 2
        )
    )
    .opacity(isDimmed ? 0.4 : 1.0)
  }

  private func stepTimerSection(
    elapsed: TimeInterval,
    minSeconds: UInt32?,
    maxSeconds: UInt32?,
    canSkip: Bool,
    onSkip: @escaping () -> Void
  ) -> some View {
    let displayText: String = {
      guard let min = minSeconds else {
        return RecipeTimeEstimation.formatElapsed(
          elapsedSeconds: elapsed, totalSeconds: maxSeconds ?? 0)
      }
      return RecipeTimeEstimation.formatElapsedWithRange(
        elapsedSeconds: elapsed, minSeconds: min, maxSeconds: maxSeconds)
    }()
    let timerColor: Color = canSkip ? .green : .orange
    return HStack(spacing: 12) {
      HStack(spacing: 8) {
        Image(systemName: "clock.fill")
          .font(.title3)
          .foregroundColor(timerColor)
        Text(displayText)
          .font(.title3)
          .fontWeight(.semibold)
          .foregroundColor(timerColor)
          .monospacedDigit()
      }
      .padding(.horizontal, 12)
      .padding(.vertical, 8)
      .background(timerColor.opacity(0.15))
      .cornerRadius(8)

      Spacer()

      Button {
        if canSkip { onSkip() }
      } label: {
        Text(canSkip ? "Proceed" : "Skip Timer")
          .font(.subheadline)
          .fontWeight(.semibold)
          .foregroundColor(.white)
          .padding(.horizontal, 16)
          .padding(.vertical, 10)
          .background(canSkip ? Color.green : Color.red)
          .cornerRadius(8)
      }
      .buttonStyle(.plain)
      .allowsHitTesting(canSkip)
    }
    .padding(.leading, 44)
  }

  private func stepBorderColor(
    isCompleted: Bool,
    isTimerActive: Bool,
    canSkipTimer: Bool,
    canCompleteStep: Bool,
    hasConditionsOrTimer: Bool
  ) -> Color {
    if isHighlighted && highlightedStepIDs != nil {
      return .blue.opacity(0.6)
    }
    if isAssociatedRecipeStep {
      return Color.purple.opacity(0.2)
    }
    if isTimerActive {
      return canSkipTimer ? Color.green.opacity(0.4) : Color.orange.opacity(0.4)
    }
    if isCompleted {
      return Color.green.opacity(0.3)
    }
    if hasConditionsOrTimer && canCompleteStep {
      return Color.green.opacity(0.4)
    }
    if hasConditionsOrTimer {
      return Color.blue.opacity(0.3)
    }
    return Color.clear
  }

  @ViewBuilder
  private func completionConditionsSection(
    completionConditions: [Mealplanning_RecipeStepCompletionCondition],
    hasTimerCondition: Bool,
    timerMinSeconds: UInt32?,
    timerMaxSeconds: UInt32?,
    isHighlighted: Bool = false
  ) -> some View {
    VStack(alignment: .leading, spacing: 8) {
      Text("Complete when")
        .font(.caption)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)
        .padding(.leading, 44)

      if hasTimerCondition, let minSec = timerMinSeconds {
        TimerConditionRow(
          viewModel: viewModel,
          recipeID: recipeID,
          stepID: step.id,
          minSeconds: minSec,
          maxSeconds: timerMaxSeconds,
          canStartStep: canCheckOverride
            ?? viewModel.canStartStep(recipeID: recipeID, stepID: step.id),
          timerElapsedSeconds: timerElapsedSeconds,
          onSkipTimerOverride: onSkipTimerOverride,
          canSkipTimerOverride: canSkipTimerOverride,
          timerMinSeconds: timerMinSeconds,
          timerMaxSeconds: timerMaxSeconds
        )
        .padding(.leading, 44)
      }

      ForEach(Array(completionConditions.enumerated()), id: \.offset) { conditionIndex, condition in
        let conditionIdentifier = viewModel.stepCompletionConditionIdentifier(
          condition: condition,
          index: conditionIndex
        )
        let isConditionCompleted = viewModel.isStepCompletionConditionCompleted(
          recipeID: recipeID,
          stepID: step.id,
          conditionIdentifier: conditionIdentifier
        )

        HStack(alignment: .top, spacing: 12) {
          Button(
            action: {
              viewModel.toggleStepCompletionCondition(
                recipeID: recipeID,
                stepID: step.id,
                conditionIdentifier: conditionIdentifier
              )
            },
            label: {
              Image(systemName: isConditionCompleted ? "checkmark.circle.fill" : "circle")
                .font(.title3)
                .foregroundColor(isConditionCompleted ? .green : .secondary)
            }
          )
          .buttonStyle(.plain)

          VStack(alignment: .leading, spacing: 4) {
            Text(completionConditionLabel(condition, position: conditionIndex))
              .font(.caption)
              .foregroundColor(isConditionCompleted ? .secondary : .primary)
              .strikethrough(isConditionCompleted)

            if condition.optional {
              Text("Optional")
                .font(.caption2)
                .foregroundColor(.secondary)
            }
          }

          Spacer()
        }
        .padding(.leading, 44)
      }
    }
    .padding(12)
    .overlay(
      RoundedRectangle(cornerRadius: 8)
        .stroke(isHighlighted ? Color.orange : Color.clear, lineWidth: isHighlighted ? 2.5 : 0)
    )
    .animation(.easeInOut(duration: 0.2), value: isHighlighted)
  }

  private func completionConditionLabel(
    _ condition: Mealplanning_RecipeStepCompletionCondition,
    position: Int
  ) -> String {
    if !condition.notes.isEmpty {
      return condition.notes
    }

    var ingredientNamesByID: [String: String] = [:]
    for ingredient in step.ingredients where !ingredient.id.isEmpty {
      ingredientNamesByID[ingredient.id] = ingredient.name
    }
    let conditionIngredientNames = condition.ingredients.compactMap { conditionIngredient in
      ingredientNamesByID[conditionIngredient.recipeStepIngredient]
    }

    if !condition.ingredientState.name.isEmpty && !conditionIngredientNames.isEmpty {
      return
        "\(condition.ingredientState.name): \(conditionIngredientNames.joined(separator: ", "))"
    }

    if !condition.ingredientState.name.isEmpty {
      return condition.ingredientState.name
    }

    if !conditionIngredientNames.isEmpty {
      return conditionIngredientNames.joined(separator: ", ")
    }

    return "Condition \(position + 1)"
  }
}

// MARK: - Step Details View

// swiftlint:disable:next type_body_length
struct StepDetailsView: View {
  @Environment(\.horizontalSizeClass) private var horizontalSizeClass

  let step: Mealplanning_RecipeStep
  let viewModel: PerformRecipeViewModel
  let stepIndex: Int
  let recipeID: String
  var mealPlanSelections: [Mealplanning_MealPlanRecipeOptionSelection]?
  var selectedIngredientOptions: [String: UInt32] = [:]
  var selectedInstrumentOptions: [String: UInt32] = [:]
  var selectedVesselOptions: [String: UInt32] = [:]
  var scale: Float = 1.0
  /// For merged meal steps: ingredient key -> "2g from Recipe A, 3g from Recipe B"
  var ingredientBreakdownBySource: [String: String]?

  private struct InstrumentVesselSectionData {
    let items: [StepItem]
    let instrumentOptionGroups: [InstrumentOptionGroupAggregate]
    let vesselOptionGroups: [VesselOptionGroupAggregate]
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 20) {
      if horizontalSizeClass == .regular {
        HStack(alignment: .top, spacing: 32) {
          ingredientSection
            .frame(maxWidth: .infinity, alignment: .topLeading)
          instrumentsAndVesselsSection
            .frame(maxWidth: .infinity, alignment: .topLeading)
        }
      } else {
        VStack(alignment: .leading, spacing: 20) {
          ingredientSection
          instrumentsAndVesselsSection
        }
      }

      notesSection
    }
    .padding(16)
    .padding(.leading, 28)  // Align with step content (44 - 16)
    .frame(maxWidth: .infinity, alignment: .leading)
    .background(Color(.systemGray6).opacity(0.6))
    .cornerRadius(10)
  }

  @ViewBuilder
  private var ingredientSection: some View {
    let (regular, optionGroups) = groupIngredientsForStep(
      step.ingredients,
      stepID: step.id,
      stepIndex: stepIndex,
      recipeID: recipeID
    )

    StepItemsSectionView(
      title: "Ingredients",
      recipeID: recipeID,
      items: regular.map { ingredient in
        let isProduct = ingredient.hasRecipeStepProductID
        let productID = isProduct ? ingredient.recipeStepProductID : nil
        let displayInfo = productID.flatMap { viewModel.getStepDisplayInfoForProductID($0) }
        let prerequisiteCompleted =
          displayInfo.map { viewModel.isStepCompleted(recipeID: recipeID, stepID: $0.stepID) }
          ?? true
        let ingredientKey =
          ingredient.hasIngredient
          ? "\(ingredient.ingredient.id)|\(ingredient.hasMeasurementUnit ? ingredient.measurementUnit.id : "")"
          : ""
        let breakdownSuffix = ingredientBreakdownBySource?[ingredientKey]

        return StepItem(
          name: formatStepIngredientDisplay(
            ingredient, scale: scale, breakdownSuffix: breakdownSuffix),
          isProduct: isProduct,
          prerequisiteStepLabel: displayInfo?.label,
          prerequisiteStepID: displayInfo?.stepID,
          prerequisiteCompleted: prerequisiteCompleted
        )
      },
      ingredientOptionGroups: filterOptionGroups(
        optionGroups, for: .ingredient, selectedOptions: selectedIngredientOptions),
      scale: scale
    )
  }

  @ViewBuilder
  private var instrumentsAndVesselsSection: some View {
    let data = instrumentVesselSectionData()
    StepItemsSectionView(
      title: "Equipment",
      recipeID: recipeID,
      items: data.items,
      instrumentOptionGroups: data.instrumentOptionGroups,
      vesselOptionGroups: data.vesselOptionGroups,
      scale: scale
    )
  }

  @ViewBuilder
  private var notesSection: some View {
    if !step.notes.isEmpty {
      Text(step.notes)
        .font(.subheadline)
        .foregroundColor(.secondary)
        .italic()
        .padding(.top, 4)
    }
  }

  private func instrumentVesselSectionData() -> InstrumentVesselSectionData {
    let (regularInstruments, instrumentGroups) = groupInstrumentsForStep(
      step.instruments,
      stepID: step.id,
      stepIndex: stepIndex,
      recipeID: recipeID
    )
    let (regularVessels, vesselGroups) = groupVesselsForStep(
      step.vessels,
      stepID: step.id,
      stepIndex: stepIndex,
      recipeID: recipeID
    )

    let mappedInstruments: [StepItem] = regularInstruments.map { instrument in
      let isProduct = instrument.hasRecipeStepProductID
      let productID = isProduct ? instrument.recipeStepProductID : nil
      let displayInfo = productID.flatMap { viewModel.getStepDisplayInfoForProductID($0) }
      let prerequisiteCompleted =
        displayInfo.map { viewModel.isStepCompleted(recipeID: recipeID, stepID: $0.stepID) } ?? true

      return StepItem(
        name: instrument.name,
        isProduct: isProduct,
        prerequisiteStepLabel: displayInfo?.label,
        prerequisiteStepID: displayInfo?.stepID,
        prerequisiteCompleted: prerequisiteCompleted
      )
    }

    let mappedVessels: [StepItem] = regularVessels.map { vessel in
      let isProduct = vessel.hasRecipeStepProductID
      let productID = isProduct ? vessel.recipeStepProductID : nil
      let displayInfo = productID.flatMap { viewModel.getStepDisplayInfoForProductID($0) }
      let prerequisiteCompleted =
        displayInfo.map { viewModel.isStepCompleted(recipeID: recipeID, stepID: $0.stepID) } ?? true

      return StepItem(
        name: vessel.name,
        isProduct: isProduct,
        prerequisiteStepLabel: displayInfo?.label,
        prerequisiteStepID: displayInfo?.stepID,
        prerequisiteCompleted: prerequisiteCompleted
      )
    }

    return InstrumentVesselSectionData(
      items: mappedInstruments + mappedVessels,
      instrumentOptionGroups: filterInstrumentOptionGroups(
        instrumentGroups, selectedOptions: selectedInstrumentOptions),
      vesselOptionGroups: filterVesselOptionGroups(
        vesselGroups, selectedOptions: selectedVesselOptions)
    )
  }

  // Helper functions to group items for a single step
  private func groupIngredientsForStep(
    _ ingredients: [Mealplanning_RecipeStepIngredient],
    stepID: String,
    stepIndex: Int,
    recipeID: String
  ) -> (
    regular: [Mealplanning_RecipeStepIngredient],
    optionGroups: [OptionGroupAggregate]
  ) {
    // Use the extension function from RecipePerformanceContentView
    // We need to access it through a temporary view extension
    var regular: [Mealplanning_RecipeStepIngredient] = []
    var optionGroupsByIndex: [UInt32: [Mealplanning_RecipeStepIngredient]] = [:]

    for ingredient in ingredients {
      let index = ingredient.index
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

    var optionGroups: [OptionGroupAggregate] = []
    for (index, groupIngredients) in optionGroupsByIndex {
      let sorted = groupIngredients.sorted { lhs, rhs in
        let lhsIndex = lhs.optionIndex
        let rhsIndex = rhs.optionIndex
        return lhsIndex < rhsIndex
      }

      var options: [IngredientOption] = []
      for ingredient in sorted {
        let optionIndex = ingredient.optionIndex
        let optionID = "\(stepID)-\(index)-\(optionIndex)"

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
          selectedOptionIndex: nil,
          sourceRecipeID: nil,
          sourceRecipeName: nil
        )
      )
    }

    optionGroups.sort { $0.index < $1.index }
    return (regular: regular, optionGroups: optionGroups)
  }

  private func groupInstrumentsForStep(
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

  private func groupVesselsForStep(
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

  // Filter option groups based on meal plan selections or user selections
  private func filterOptionGroups(
    _ groups: [OptionGroupAggregate],
    for selectionType: Mealplanning_MealPlanRecipeOptionSelectionType,
    selectedOptions: [String: UInt32]
  ) -> [OptionGroupAggregate] {
    return groups.map { group in
      // Check meal plan selections first, then user selections
      let selectedIndex: UInt32?
      if let selections = mealPlanSelections,
        let selection = selections.first(where: { sel in
          sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
            && sel.ingredientIndex == group.index && sel.selectionType == selectionType
        })
      {
        selectedIndex = selection.selectedOptionIndex
      } else if let userSelection = selectedOptions[group.id] {
        selectedIndex = userSelection
      } else {
        selectedIndex = nil  // No selection - show all options
      }

      // If selection exists, show only that option; otherwise show all options
      let filteredOptions: [IngredientOption]
      if let selectedIndex = selectedIndex {
        filteredOptions = group.options.filter { $0.optionIndex == selectedIndex }
      } else {
        filteredOptions = group.options  // Show all options when no selection
      }

      return OptionGroupAggregate(
        id: group.id,
        recipeID: group.recipeID,
        stepID: group.stepID,
        stepIndex: group.stepIndex,
        index: group.index,
        options: filteredOptions,
        selectedOptionIndex: selectedIndex,
        sourceRecipeID: group.sourceRecipeID,
        sourceRecipeName: group.sourceRecipeName
      )
    }
  }

  private func filterInstrumentOptionGroups(
    _ groups: [InstrumentOptionGroupAggregate],
    selectedOptions: [String: UInt32]
  ) -> [InstrumentOptionGroupAggregate] {
    return groups.map { group in
      // Check meal plan selections first, then user selections
      let selectedIndex: UInt32?
      if let selections = mealPlanSelections,
        let selection = selections.first(where: { sel in
          sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
            && sel.selectionType == .instrument
        })
      {
        selectedIndex = selection.selectedOptionIndex
      } else if let userSelection = selectedOptions[group.id] {
        selectedIndex = userSelection
      } else {
        selectedIndex = nil  // No selection - show all options
      }

      // If selection exists, show only that option; otherwise show all options
      let filteredOptions: [InstrumentOption]
      if let selectedIndex = selectedIndex {
        filteredOptions = group.options.filter { $0.optionIndex == selectedIndex }
      } else {
        filteredOptions = group.options  // Show all options when no selection
      }

      return InstrumentOptionGroupAggregate(
        id: group.id,
        recipeID: group.recipeID,
        stepID: group.stepID,
        stepIndex: group.stepIndex,
        index: group.index,
        options: filteredOptions,
        selectedOptionIndex: selectedIndex,
        sourceRecipeID: group.sourceRecipeID,
        sourceRecipeName: group.sourceRecipeName
      )
    }
  }

  private func filterVesselOptionGroups(
    _ groups: [VesselOptionGroupAggregate],
    selectedOptions: [String: UInt32]
  ) -> [VesselOptionGroupAggregate] {
    return groups.map { group in
      // Check meal plan selections first, then user selections
      let selectedIndex: UInt32?
      if let selections = mealPlanSelections,
        let selection = selections.first(where: { sel in
          sel.recipeID == group.recipeID && sel.recipeStepID == group.stepID
            && sel.selectionType == .vessel
        })
      {
        selectedIndex = selection.selectedOptionIndex
      } else if let userSelection = selectedOptions[group.id] {
        selectedIndex = userSelection
      } else {
        selectedIndex = nil  // No selection - show all options
      }

      // If selection exists, show only that option; otherwise show all options
      let filteredOptions: [VesselOption]
      if let selectedIndex = selectedIndex {
        filteredOptions = group.options.filter { $0.optionIndex == selectedIndex }
      } else {
        filteredOptions = group.options  // Show all options when no selection
      }

      return VesselOptionGroupAggregate(
        id: group.id,
        recipeID: group.recipeID,
        stepID: group.stepID,
        stepIndex: group.stepIndex,
        index: group.index,
        options: filteredOptions,
        selectedOptionIndex: selectedIndex,
        sourceRecipeID: group.sourceRecipeID,
        sourceRecipeName: group.sourceRecipeName
      )
    }
  }
}

// MARK: - Step Items Section View

struct StepItemsSectionView: View {
  @Environment(\.horizontalSizeClass) private var horizontalSizeClass
  @Environment(AuthenticationManager.self) private var authManager

  let title: String
  let recipeID: String
  let items: [StepItem]
  var ingredientOptionGroups: [OptionGroupAggregate] = []
  var instrumentOptionGroups: [InstrumentOptionGroupAggregate] = []
  var vesselOptionGroups: [VesselOptionGroupAggregate] = []
  var scale: Float = 1.0

  var body: some View {
    VStack(alignment: .leading, spacing: 12) {
      Text(title)
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)

      if items.isEmpty
        && ingredientOptionGroups.isEmpty
        && instrumentOptionGroups.isEmpty
        && vesselOptionGroups.isEmpty
      {
        Text("none")
          .font(.subheadline)
          .foregroundColor(.secondary)
      }

      // Regular items
      VStack(alignment: .leading, spacing: 8) {
        ForEach(Array(items.enumerated()), id: \.offset) { _, item in
          itemRow(item)
        }

        // Ingredient option groups
        ForEach(ingredientOptionGroups) { group in
          OptionGroupView(group: group, scale: scale)
        }

        // Instrument option groups
        ForEach(instrumentOptionGroups) { group in
          InstrumentOptionGroupView(group: group, scale: scale)
        }

        // Vessel option groups
        ForEach(vesselOptionGroups) { group in
          VesselOptionGroupView(group: group, scale: scale)
        }
      }
    }
  }

  @ViewBuilder
  private func itemRow(_ item: StepItem) -> some View {
    HStack(alignment: .top, spacing: 8) {
      if item.isProduct && !item.prerequisiteCompleted {
        Image(systemName: "clock.fill")
          .font(.subheadline)
          .foregroundColor(.orange)
      }
      Text(item.name)
        .font(.subheadline)
        .foregroundColor(
          (item.isProduct && !item.prerequisiteCompleted) ? .orange : .primary
        )
        .lineLimit(3)
      if let label = item.prerequisiteStepLabel {
        let referenceText = "(from \(label))"
        if let stepID = item.prerequisiteStepID {
          NavigationLink {
            PerformRecipeView(
              recipeID: recipeID,
              highlightedStepIDs: [stepID]
            )
            .environment(authManager)
          } label: {
            Text(referenceText)
              .font(.caption2)
              .foregroundColor(.secondary)
          }
          .buttonStyle(.plain)
        } else {
          Text(referenceText)
            .font(.caption2)
            .foregroundColor(.secondary)
        }
      }
    }
  }
}

// MARK: - Option Group Views

struct OptionGroupView: View {
  let group: OptionGroupAggregate
  var scale: Float = 1.0

  var body: some View {
    // If a selection has been made (selectedOptionIndex is not nil), show only that option without indentation
    if group.selectedOptionIndex != nil,
      group.options.count == 1,
      let option = group.options.first
    {
      // Selected option - show like a regular ingredient (no indentation, no "one of:" label)
      HStack(alignment: .top, spacing: 8) {
        if let quantityText = option.aggregated.quantityText(scale: scale) {
          Text(quantityText)
            .font(.subheadline)
            .foregroundColor(.secondary)
        }

        VStack(alignment: .leading, spacing: 2) {
          Text(option.ingredient.name)
            .font(.subheadline)
            .foregroundColor(.primary)
          if !option.aggregated.quantityNotes.isEmpty {
            Text(option.aggregated.quantityNotes)
              .font(.caption)
              .foregroundColor(.secondary)
          }
        }
      }
    } else {
      // No selection - show all options with "one of:" label (indented)
      VStack(alignment: .leading, spacing: 8) {
        Text("one of:")
          .font(.subheadline)
          .foregroundColor(.secondary)
          .padding(.leading, 16)

        ForEach(group.options) { option in
          HStack(alignment: .top, spacing: 8) {
            if let quantityText = option.aggregated.quantityText(scale: scale) {
              Text(quantityText)
                .font(.subheadline)
                .foregroundColor(.secondary)
            }

            VStack(alignment: .leading, spacing: 2) {
              Text(option.ingredient.name)
                .font(.subheadline)
                .foregroundColor(.primary)
              if !option.aggregated.quantityNotes.isEmpty {
                Text(option.aggregated.quantityNotes)
                  .font(.caption)
                  .foregroundColor(.secondary)
              }
            }
          }
          .padding(.leading, 16)
        }
      }
    }
  }
}

struct InstrumentOptionGroupView: View {
  let group: InstrumentOptionGroupAggregate
  var scale: Float = 1.0

  var body: some View {
    // If a selection has been made (selectedOptionIndex is not nil), show only that option without indentation
    if group.selectedOptionIndex != nil,
      group.options.count == 1,
      let option = group.options.first
    {
      // Selected option - show like a regular instrument (no indentation, no "one of:" label)
      HStack(spacing: 8) {
        Text(option.instrument.name)
          .font(.subheadline)
          .foregroundColor(.primary)
      }
    } else {
      // No selection - show all options with "one of:" label (indented)
      VStack(alignment: .leading, spacing: 8) {
        Text("one of:")
          .font(.subheadline)
          .foregroundColor(.secondary)
          .padding(.leading, 16)

        ForEach(group.options) { option in
          HStack(spacing: 8) {
            Text(option.instrument.name)
              .font(.subheadline)
              .foregroundColor(.primary)
          }
          .padding(.leading, 16)
        }
      }
    }
  }
}

struct VesselOptionGroupView: View {
  let group: VesselOptionGroupAggregate
  var scale: Float = 1.0

  var body: some View {
    // If a selection has been made (selectedOptionIndex is not nil), show only that option without indentation
    if group.selectedOptionIndex != nil,
      group.options.count == 1,
      let option = group.options.first
    {
      // Selected option - show like a regular vessel (no indentation, no "one of:" label)
      HStack(spacing: 8) {
        Text(option.vessel.name)
          .font(.subheadline)
          .foregroundColor(.primary)
      }
    } else {
      // No selection - show all options with "one of:" label (indented)
      VStack(alignment: .leading, spacing: 8) {
        Text("one of:")
          .font(.subheadline)
          .foregroundColor(.secondary)
          .padding(.leading, 16)

        ForEach(group.options) { option in
          HStack(spacing: 8) {
            Text(option.vessel.name)
              .font(.subheadline)
              .foregroundColor(.primary)
          }
          .padding(.leading, 16)
        }
      }
    }
  }
}

// MARK: - Step Products Section View

struct StepProductsSectionView: View {
  let products: [Mealplanning_RecipeStepProduct]
  var scale: Float = 1.0

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      Text("Products")
        .font(.subheadline)
        .fontWeight(.semibold)
        .foregroundColor(.secondary)

      ForEach(Array(products.enumerated()), id: \.offset) { _, product in
        HStack(spacing: 6) {
          Text(formatProductQuantity(product, scale: scale))
            .font(.caption)
            .foregroundColor(.secondary)
        }
      }
    }
  }

  private func formatProductQuantity(_ product: Mealplanning_RecipeStepProduct, scale: Float)
    -> String
  {
    // Check if product is discrete (has ItemQuantity set)
    let isDiscrete =
      product.hasItemQuantity && (product.itemQuantity.hasMin || product.itemQuantity.hasMax)

    if isDiscrete {
      // Discrete product: item quantity scales, per-item measurement stays constant
      var itemQtyStr = ""
      if product.itemQuantity.hasMin {
        let min = Float(product.itemQuantity.min) * scale
        if product.itemQuantity.hasMax {
          let max = Float(product.itemQuantity.max) * scale
          if min == max {
            itemQtyStr = formatQuantity(min)
          } else {
            itemQtyStr = "\(formatQuantity(min))-\(formatQuantity(max))"
          }
        } else {
          itemQtyStr = formatQuantity(min)
        }
      }

      // Per-item measurement quantity does NOT scale (stays constant)
      var measurementQtyStr = ""
      if product.hasMeasurementQuantity && product.measurementQuantity.hasMin {
        let min = product.measurementQuantity.min  // Not scaled
        if product.measurementQuantity.hasMax {
          let max = product.measurementQuantity.max  // Not scaled
          if min == max {
            measurementQtyStr = formatQuantity(min)
          } else {
            measurementQtyStr = "\(formatQuantity(min))-\(formatQuantity(max))"
          }
        } else {
          measurementQtyStr = formatQuantity(min)
        }
      }

      let unitName =
        product.hasMeasurementUnit
        ? MeasurementUnitFormatter.displayName(
          for: product.measurementQuantity.min,
          unit: product.measurementUnit
        )
        : ""

      if !itemQtyStr.isEmpty && !measurementQtyStr.isEmpty && !unitName.isEmpty {
        // Format: "4 patties (4 oz each)"
        return "\(itemQtyStr) \(product.name) (\(measurementQtyStr) \(unitName) each)"
      } else if !itemQtyStr.isEmpty {
        // Fallback: just show count and name if measurement is missing
        return "\(itemQtyStr) \(product.name)"
      } else {
        // Fallback: just show name if quantities are missing
        return product.name
      }
    } else if product.hasMeasurementQuantity && product.measurementQuantity.hasMin {
      // Continuous product: total measurement quantity scales
      let min = product.measurementQuantity.min * scale
      var qtyStr = formatQuantity(min)

      if product.measurementQuantity.hasMax {
        let max = product.measurementQuantity.max * scale
        if min != max {
          qtyStr = "\(qtyStr)-\(formatQuantity(max))"
        }
      }

      let unitName =
        product.hasMeasurementUnit
        ? MeasurementUnitFormatter.displayName(for: min, unit: product.measurementUnit)
        : ""
      if !unitName.isEmpty {
        return "\(product.name): \(qtyStr) \(unitName)"
      } else {
        return "\(product.name): \(qtyStr)"
      }
    }

    // Fallback: just show name if no quantities
    return product.name
  }

  private func formatQuantity(_ qty: Float) -> String {
    // Format numbers - use fewer decimals for whole numbers
    if qty.truncatingRemainder(dividingBy: 1) == 0 {
      return String(format: "%.0f", qty)
    } else {
      return String(format: "%.2f", qty)
    }
  }
}
