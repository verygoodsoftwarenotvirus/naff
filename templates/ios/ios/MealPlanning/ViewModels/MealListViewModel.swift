//
//  MealListViewModel.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import SwiftProtobuf
import SwiftUI

@Observable
@MainActor
class MealListViewModel {
  var meals: [Mealplanning_Meal] = []
  var searchResults: [Mealplanning_Meal] = []
  var isLoading = false
  var isSearching = false
  var errorMessage: String?
  var searchError: String?

  private let authManager: AuthenticationManager
  private var searchTask: Task<Void, Never>?

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  var displayedMeals: [Mealplanning_Meal] {
    // If we have search results, show those; otherwise show all meals
    return searchResults.isEmpty ? meals : searchResults
  }

  var isInSearchMode: Bool {
    return !searchResults.isEmpty
  }

  func loadMeals() async {
    isLoading = true
    errorMessage = nil

    if APIConfiguration.currentEnvironment.isOffline {
      let provider = await LocalDataProvider.shared
      provider.loadIfNeeded()
      self.meals = provider.getMeals()
      isLoading = false
      return
    }

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "MealListViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      // Get OAuth2 token (will refresh if needed)
      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "MealListViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create request
      let request = Mealplanning_GetMealsRequest()

      // Execute request
      let response = try await clientManager.client.mealPlanning.getMeals(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      self.meals = response.results
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = "Failed to load meals: \(error.localizedDescription)"
      print("❌ Error loading meals: \(error)")
    }

    isLoading = false
  }

  func searchMeals(query: String) {
    // Cancel any existing search task
    searchTask?.cancel()

    let trimmedQuery = query.trimmingCharacters(in: .whitespacesAndNewlines)

    // If query is empty, clear search results
    if trimmedQuery.isEmpty {
      searchResults = []
      searchError = nil
      isSearching = false
      return
    }

    // Debounce: wait 500ms before executing search
    searchTask = Task {
      try? await Task.sleep(nanoseconds: 500_000_000)  // 500ms

      // Check if task was cancelled
      guard !Task.isCancelled else { return }

      await performSearch(query: trimmedQuery)
    }
  }

  private func performSearch(query: String) async {
    isSearching = true
    searchError = nil

    if APIConfiguration.currentEnvironment.isOffline {
      let provider = await LocalDataProvider.shared
      searchResults = provider.searchMeals(query: query)
      isSearching = false
      return
    }

    do {
      guard let clientManager = try? authManager.getClientManager() else {
        throw NSError(
          domain: "MealListViewModel", code: 1,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
      }

      guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
        throw NSError(
          domain: "MealListViewModel", code: 2,
          userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
      }

      let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

      // Create search request
      var request = Mealplanning_SearchForMealsRequest()
      request.query = query
      request.useSearchService = APIConfiguration.useSearchService

      // Execute search
      let response = try await clientManager.client.mealPlanning.searchForMeals(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

      searchResults = response.results
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      searchError = "Failed to search meals: \(error.localizedDescription)"
      print("❌ Error searching for meals: \(error)")
      searchResults = []
    }

    isSearching = false
  }
}
