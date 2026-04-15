//
//  PortionsFormatter.swift
//  ios
//
//  Formats estimated portions for display.
//  - No max: treat min as authoritative (not "limitless")
//  - Same min and max: single authoritative number
//  - Different min and max: show range
//

import SwiftProtobuf

enum PortionsFormatter {
  /// Formats a portion range for display.
  /// - No max: shows min as the authoritative number (e.g. "4")
  /// - Min equals max: shows single number (e.g. "4")
  /// - Min differs from max: shows range (e.g. "4-6")
  static func format(_ range: Common_Float32RangeWithOptionalMax) -> String {
    if range.hasMax {
      if range.min == range.max {
        return String(format: "%.1f", range.min)
      } else {
        return String(format: "%.1f-%.1f", range.min, range.max)
      }
    } else {
      return String(format: "%.1f", range.min)
    }
  }

  /// Formats a scaled portion range for display (e.g. when scaling a meal).
  static func formatScaled(_ range: Common_Float32RangeWithOptionalMax, scale: Float) -> String {
    let scaledMin = range.min * scale
    if range.hasMax {
      let scaledMax = range.max * scale
      if scaledMin == scaledMax {
        return String(format: "%.1f", scaledMin)
      } else {
        return String(format: "%.1f-%.1f", scaledMin, scaledMax)
      }
    } else {
      return String(format: "%.1f", scaledMin)
    }
  }
}
