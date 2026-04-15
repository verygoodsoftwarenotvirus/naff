//
//  OptionSelectionStepView.swift
//  ios
//
//  Created by Auto on 3/7/25.
//

import SwiftProtobuf
import SwiftUI

/// Dedicated wizard step for selecting ingredient options before creating a meal plan.
/// Ensures users are explicitly prompted to choose options rather than skipping via a dismissible sheet.
struct OptionSelectionStepView: View {
  @Bindable var viewModel: CreateMealPlanViewModel
  let onDismiss: () -> Void
  @Environment(EventReporterService.self) private var eventReporterService

  @State private var selections: [String: [String: [UInt32: UInt32]]] = [:]
  @State private var isCreating = false

  private var recipes: [Mealplanning_Recipe] {
    viewModel.getAllRecipes(from: viewModel.allSelectedMeals)
  }

  var body: some View {
    VStack(alignment: .leading, spacing: 24) {
      VStack(alignment: .leading, spacing: 8) {
        Text("Choose Ingredient Options")
          .font(.title2)
          .fontWeight(.semibold)

        Text(
          "Some recipes have ingredient choices (e.g., chicken vs tofu). Select your preferences below."
        )
        .font(.subheadline)
        .foregroundColor(.secondary)
      }

      RecipeOptionSelectionContent(
        recipes: recipes,
        selections: $selections
      )

      Spacer(minLength: 24)

      HStack(spacing: 12) {
        Button {
          eventReporterService.reporter.track(event: "meal_plan_option_back", properties: [:])
          viewModel.wizardStep = .mealAssignment
        } label: {
          HStack(spacing: 6) {
            Image(systemName: "chevron.left")
            Text("Back")
          }
          .font(.subheadline.weight(.semibold))
          .frame(maxWidth: .infinity)
          .padding()
          .background(Color(.systemGray6))
          .foregroundColor(.blue)
          .cornerRadius(10)
        }

        Button {
          createMealPlan()
        } label: {
          HStack {
            if isCreating {
              ProgressView()
                .progressViewStyle(CircularProgressViewStyle(tint: .white))
            }
            Text(isCreating ? "Creating..." : "Create Meal Plan")
              .fontWeight(.semibold)
          }
          .frame(maxWidth: .infinity)
          .padding()
          .background(isCreating ? Color.gray : Color.blue)
          .foregroundColor(.white)
          .cornerRadius(10)
        }
        .disabled(isCreating)
      }
    }
    .frame(maxWidth: .infinity, alignment: .leading)
  }

  private func createMealPlan() {
    var finalSelections = selections
    if finalSelections.isEmpty {
      for recipe in recipes {
        let defaults = viewModel.getDefaultOptionSelections(for: recipe)
        if !defaults.isEmpty {
          finalSelections[recipe.id] = defaults
        }
      }
    }
    viewModel.setOptionSelections(ingredientSelections: finalSelections)
    isCreating = true
    Task {
      let success = await viewModel.createMealPlan()
      isCreating = false
      if success {
        eventReporterService.reporter.track(
          event: "meal_plan_created",
          properties: [
            "events_count": viewModel.selectedDates.count,
            "meals_count": viewModel.selectedDates.count,
          ])
        NotificationCenter.default.post(name: .mealPlanCreated, object: nil)
        onDismiss()
      }
    }
  }
}
