//
//  MealListView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

struct MealListView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(EventReporterService.self) private var eventReporterService
  @State private var viewModel: MealListViewModel?
  @State private var searchQuery: String = ""

  var body: some View {
    NavigationStack {
      Group {
        if let viewModel = viewModel {
          DSContentState(
            isLoading: viewModel.isLoading,
            loadingMessage: "Loading meals...",
            error: viewModel.errorMessage,
            onRetry: { await viewModel.loadMeals() },
            content: {
              let displayedMeals = viewModel.displayedMeals
              let isSearching = viewModel.isSearching
              let hasSearchQuery = !searchQuery.trimmingCharacters(in: .whitespacesAndNewlines)
                .isEmpty

              if isSearching {
                DSLoadingView("Searching meals...")
              } else if displayedMeals.isEmpty {
                if hasSearchQuery {
                  VStack(spacing: DSTheme.Spacing.lg) {
                    DSEmptyState(
                      icon: "magnifyingglass",
                      title: "No Results",
                      message: "No meals found matching \"\(searchQuery)\""
                    )

                    if let searchError = viewModel.searchError {
                      Text(searchError)
                        .font(DSTheme.Typography.caption)
                        .foregroundColor(DSTheme.Colors.error)
                        .padding(.horizontal)
                    }
                  }
                } else {
                  DSEmptyState(
                    icon: "fork.knife",
                    title: "No Meals",
                    message: "No meals found. Create some meals to get started.",
                    size: .large
                  )
                }
              } else {
                ScrollView {
                  LazyVStack(spacing: DSTheme.Spacing.md) {
                    ForEach(displayedMeals, id: \.id) { meal in
                      NavigationLink(destination: MealDetailView(mealID: meal.id)) {
                        MealCard(meal: meal)
                      }
                      .simultaneousGesture(
                        TapGesture().onEnded {
                          eventReporterService.reporter.track(
                            event: "meal_card_tapped",
                            properties: ["meal_id": meal.id])
                        })
                    }
                  }
                  .dsScreenPadding()
                }
              }
            })
        } else {
          DSInitializingView()
        }
      }
      .navigationTitle("Meals")
      .navigationBarTitleDisplayMode(.large)
      .toolbar {
        ToolbarItem(placement: .primaryAction) {
          NavigationLink {
            CreateMealView()
          } label: {
            Image(systemName: "plus")
          }
          .simultaneousGesture(
            TapGesture().onEnded {
              eventReporterService.reporter.track(event: "create_meal_tapped", properties: [:])
            })
        }
      }
      .searchable(text: $searchQuery, prompt: "Search meals...")
      .onChange(of: searchQuery) { _, newValue in
        if let viewModel = viewModel {
          viewModel.searchMeals(query: newValue)
        }
      }
      .refreshable {
        if let viewModel = viewModel {
          searchQuery = ""
          viewModel.searchResults = []
          await viewModel.loadMeals()
        }
      }
      .onAppear {
        if viewModel == nil {
          viewModel = MealListViewModel(authManager: authManager)
          eventReporterService.reporter.track(event: "meals_list_viewed", properties: [:])
        }
        if let viewModel = viewModel {
          Task {
            await viewModel.loadMeals()
          }
        }
      }
    }
  }
}

// MARK: - Meal Card

struct MealCard: View {
  let meal: Mealplanning_Meal

  var body: some View {
    DSCard(style: .outlined) {
      VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
        Text(meal.name.isEmpty ? "Unnamed Meal" : meal.name)
          .font(DSTheme.Typography.label)
          .foregroundColor(DSTheme.Colors.textPrimary)
          .frame(maxWidth: .infinity, alignment: .leading)

        if !meal.description_p.isEmpty {
          Text(meal.description_p)
            .font(DSTheme.Typography.body)
            .foregroundColor(DSTheme.Colors.textSecondary)
            .lineLimit(2)
            .frame(maxWidth: .infinity, alignment: .leading)
        }

        // Meal metadata - time, components, servings
        HStack(spacing: DSTheme.Spacing.md) {
          if !meal.components.isEmpty, let totalTime = totalEstimatedTime(for: meal.components) {
            Label(
              RecipeTimeEstimation.format(
                minSeconds: totalTime.minSeconds, maxSeconds: totalTime.maxSeconds),
              systemImage: "clock"
            )
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
          }

          if !meal.components.isEmpty {
            Label(
              "\(meal.components.count) component\(meal.components.count == 1 ? "" : "s")",
              systemImage: "square.stack.3d.up"
            )
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
          }

          if meal.hasEstimatedPortions {
            Label(
              "\(PortionsFormatter.format(meal.estimatedPortions)) servings",
              systemImage: "person.2"
            )
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
          }
        }

        // Component type breakdown (when available)
        if !meal.components.isEmpty {
          let componentTypeSummary = formatComponentTypeSummary(meal.components)
          if !componentTypeSummary.isEmpty {
            HStack(spacing: DSTheme.Spacing.xs) {
              ForEach(componentTypeLabels(from: meal.components), id: \.self) { label in
                Text(label)
                  .font(DSTheme.Typography.caption)
                  .foregroundColor(DSTheme.Colors.textSecondary)
                  .padding(.horizontal, DSTheme.Spacing.xs)
                  .padding(.vertical, 2)
                  .background(DSTheme.Colors.borderSubtle.opacity(0.5))
                  .clipShape(RoundedRectangle(cornerRadius: 4))
              }
            }
          }
        }
      }
      .frame(maxWidth: .infinity, alignment: .leading)
    }
  }

  /// Returns labels for component types: sides count (assume one main, don't show), beverage only if present.
  private func componentTypeLabels(from components: [Mealplanning_MealComponent]) -> [String] {
    var labels: [String] = []
    var sideCount = 0
    var hasBeverage = false

    for component in components {
      switch component.componentType {
      case .main:
        break  // Assume one main, don't indicate
      case .side:
        sideCount += 1
      case .beverage:
        hasBeverage = true
      default:
        break  // Omit other types for now
      }
    }

    if sideCount > 0 {
      labels.append(sideCount == 1 ? "1 side" : "\(sideCount) sides")
    }
    if hasBeverage {
      labels.append("Beverage")
    }
    return labels
  }

  private func formatComponentTypeSummary(_ components: [Mealplanning_MealComponent]) -> String {
    componentTypeLabels(from: components).joined(separator: ", ")
  }

  private func totalEstimatedTime(for components: [Mealplanning_MealComponent])
    -> RecipeTimeEstimate?
  {
    var totalMin: UInt32 = 0
    var totalMax: UInt32 = 0
    var hasAny = false
    for component in components {
      guard let estimate = RecipeTimeEstimation.estimate(steps: component.recipe.steps) else {
        continue
      }
      totalMin = totalMin &+ estimate.minSeconds
      totalMax = totalMax &+ estimate.maxSeconds
      hasAny = true
    }
    return hasAny ? RecipeTimeEstimate(minSeconds: totalMin, maxSeconds: totalMax) : nil
  }
}

// MARK: - Preview

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "Test User"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return MealListView()
    .environment(authManager)
}
