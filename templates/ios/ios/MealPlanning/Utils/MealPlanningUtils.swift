//
//  MealPlanningUtils.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import Foundation
import SwiftProtobuf

/// Utility functions for meal planning
enum MealPlanningUtils {
  /// Convert MealPlanEventName enum to a user-friendly display string
  /// - Parameter mealName: The meal plan event name enum
  /// - Returns: A formatted display string (e.g., "Breakfast", "Second Breakfast")
  static func formatMealName(_ mealName: Mealplanning_MealPlanEventName) -> String {
    switch mealName {
    case .breakfast:
      return "Breakfast"
    case .secondBreakfast:
      return "Second Breakfast"
    case .brunch:
      return "Brunch"
    case .lunch:
      return "Lunch"
    case .supper:
      return "Supper"
    case .dinner:
      return "Dinner"
    case .UNRECOGNIZED:
      return "Meal"
    }
  }

  /// Convert MealPlanEventName enum to API string format
  /// - Parameter mealName: The meal plan event name enum
  /// - Returns: A string in the format expected by the API (e.g., "breakfast", "second_breakfast")
  static func mealPlanEventNameToString(_ mealName: Mealplanning_MealPlanEventName) -> String {
    switch mealName {
    case .breakfast:
      return "breakfast"
    case .secondBreakfast:
      return "second_breakfast"
    case .brunch:
      return "brunch"
    case .lunch:
      return "lunch"
    case .supper:
      return "supper"
    case .dinner:
      return "dinner"
    case .UNRECOGNIZED:
      return "dinner"  // Default fallback
    }
  }
}
