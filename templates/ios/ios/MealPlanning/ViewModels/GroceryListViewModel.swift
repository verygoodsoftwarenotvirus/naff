//
//  GroceryListViewModel.swift
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
class GroceryListViewModel {
  // Data
  var items: [Mealplanning_MealPlanGroceryListItem] = []
  var mealPlan: Mealplanning_MealPlan

  // Loading states
  var isLoading = false
  var isUpdating = false
  var errorMessage: String?

  private let authManager: AuthenticationManager

  init(
    mealPlan: Mealplanning_MealPlan, items: [Mealplanning_MealPlanGroceryListItem],
    authManager: AuthenticationManager
  ) {
    self.mealPlan = mealPlan
    self.items = items
    self.authManager = authManager
  }

  func loadItems() async {
    isLoading = true
    errorMessage = nil

    do {
      let fetchedItems = try await fetchGroceryListItems()
      self.items = fetchedItems
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = "Failed to load grocery list: \(error.localizedDescription)"
      print("❌ Error loading grocery list: \(error)")
    }

    isLoading = false
  }

  private func fetchGroceryListItems() async throws -> [Mealplanning_MealPlanGroceryListItem] {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "GroceryListViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "GroceryListViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    var request = Mealplanning_GetMealPlanGroceryListItemsForMealPlanRequest()
    request.mealPlanID = mealPlan.id

    let response = try await clientManager.client.mealPlanning
      .getMealPlanGroceryListItemsForMealPlan(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )

    return response.results
  }

  func updateItem(
    _ item: Mealplanning_MealPlanGroceryListItem,
    status: Mealplanning_MealPlanGroceryListItemStatus? = nil,
    quantityNeededMin: Float? = nil,
    quantityNeededMax: Float? = nil,
    quantityPurchased: Float? = nil
  ) async {
    isUpdating = true
    errorMessage = nil

    do {
      var updateInput = Mealplanning_MealPlanGroceryListItemUpdateRequestInput()

      // Update status if provided
      if let status = status {
        updateInput.status = status
      }

      // Update quantity needed if provided
      if let min = quantityNeededMin {
        var quantityNeeded = Common_Float32RangeWithOptionalMaxUpdateRequestInput()
        quantityNeeded.min = min
        if let max = quantityNeededMax {
          quantityNeeded.max = max
        }
        // When max is nil, we omit it so backend clears max (single-value edit)
        updateInput.quantityNeeded = quantityNeeded
      }

      // Update quantity purchased if provided
      if let purchased = quantityPurchased {
        updateInput.quantityPurchased = purchased
      }

      let updated = try await updateGroceryListItem(itemID: item.id, input: updateInput)

      // Optimistically update local state instead of full reload
      if let index = items.firstIndex(where: { $0.id == item.id }) {
        items[index] = updated
      }
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = "Failed to update item: \(error.localizedDescription)"
      print("❌ Error updating grocery list item: \(error)")
    }

    isUpdating = false
  }

  func markAsAcquired(_ item: Mealplanning_MealPlanGroceryListItem) async {
    await updateItem(item, status: .acquired)
  }

  func markAsAlreadyOwned(_ item: Mealplanning_MealPlanGroceryListItem) async {
    await updateItem(item, status: .alreadyOwned)
  }

  func markAsNeeds(_ item: Mealplanning_MealPlanGroceryListItem) async {
    await updateItem(item, status: .needs)
  }

  func markAsUnavailable(_ item: Mealplanning_MealPlanGroceryListItem) async {
    await updateItem(item, status: .unavailable)
  }

  func updateQuantityNeeded(
    _ item: Mealplanning_MealPlanGroceryListItem,
    min: Float,
    max: Float?
  ) async {
    await updateItem(
      item,
      quantityNeededMin: min,
      quantityNeededMax: max
    )
  }

  func updateQuantityPurchased(
    _ item: Mealplanning_MealPlanGroceryListItem,
    quantity: Float
  ) async {
    await updateItem(
      item,
      quantityPurchased: quantity
    )
  }

  private func updateGroceryListItem(
    itemID: String,
    input: Mealplanning_MealPlanGroceryListItemUpdateRequestInput
  ) async throws -> Mealplanning_MealPlanGroceryListItem {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "GroceryListViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "GroceryListViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)

    var request = Mealplanning_UpdateMealPlanGroceryListItemRequest()
    request.mealPlanID = mealPlan.id
    request.mealPlanGroceryListItemID = itemID
    request.input = input

    let response = try await clientManager.client.mealPlanning.updateMealPlanGroceryListItem(
      request,
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )
    return response.updated
  }

  // Computed properties for filtering
  var itemsByStatus:
    [Mealplanning_MealPlanGroceryListItemStatus: [Mealplanning_MealPlanGroceryListItem]]
  {
    Dictionary(grouping: items) { $0.status }
  }

  var needsItems: [Mealplanning_MealPlanGroceryListItem] {
    items.filter { $0.status == .needs }
  }

  var alreadyOwnedItems: [Mealplanning_MealPlanGroceryListItem] {
    items.filter { $0.status == .alreadyOwned }
  }

  var acquiredItems: [Mealplanning_MealPlanGroceryListItem] {
    items.filter { $0.status == .acquired }
  }

  var unavailableItems: [Mealplanning_MealPlanGroceryListItem] {
    items.filter { $0.status == .unavailable }
  }

  /// Combined acquired + already owned for the "Have" section.
  var haveItems: [Mealplanning_MealPlanGroceryListItem] {
    items.filter { $0.status == .acquired || $0.status == .alreadyOwned }
  }
}
