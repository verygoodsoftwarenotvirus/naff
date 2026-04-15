//
//  MealComponentTypeFormatter.swift
//  ios
//
//  Shared formatter for MealComponentType display strings.
//

import SwiftProtobuf

enum MealComponentTypeFormatter {
  /// Returns a human-readable string for the given component type, or empty string for unspecified.
  static func format(_ type: Mealplanning_MealComponentType) -> String {
    switch type {
    case .amuseBouche:
      return "Amuse Bouche"
    case .appetizer:
      return "Appetizer"
    case .soup:
      return "Soup"
    case .main:
      return "Main"
    case .salad:
      return "Salad"
    case .beverage:
      return "Beverage"
    case .side:
      return "Side"
    case .dessert:
      return "Dessert"
    default:
      return ""
    }
  }
}
