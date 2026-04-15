//
//  MeasurementUnitFormatter.swift
//  ios
//
//  Formats measurement unit names with correct singular/plural form based on quantity.
//  Uses ValidMeasurementUnit.pluralName when quantity != 1; otherwise uses .name.
//

import SwiftProtobuf

enum MeasurementUnitFormatter {
  /// Returns the display name for a measurement unit based on quantity.
  /// - quantity == 1: returns singular form (unit.name)
  /// - quantity != 1: returns plural form (unit.pluralName if available, else unit.name)
  static func displayName(for quantity: Float, unit: Mealplanning_ValidMeasurementUnit?) -> String {
    guard let unit = unit, !unit.name.isEmpty else { return "" }
    if abs(quantity - 1) < 0.001 {
      return unit.name
    }
    if !unit.pluralName.isEmpty {
      return unit.pluralName
    }
    return unit.name
  }
}
