//
//  RecipeTimeEstimationTests.swift
//  iosTests
//
//  Unit tests for recipe time estimation (aggregating step times with 1-3 min fallback).
//

import Foundation
import SwiftProtobuf
@testable import ios
import Testing

// MARK: - Test Helpers

private func makeStep(
  stepID: String,
  estimatedTime: (min: UInt32?, max: UInt32?)
) -> Mealplanning_RecipeStep {
  var prep = Mealplanning_ValidPreparation()
  prep.id = "prep-\(stepID)"
  prep.name = "Step"

  var step = Mealplanning_RecipeStep()
  step.id = stepID
  step.preparation = prep
  step.index = 0

  if let minVal = estimatedTime.min {
    step.estimatedTimeInSeconds.min = minVal
  }
  if let maxVal = estimatedTime.max {
    step.estimatedTimeInSeconds.max = maxVal
  }
  return step
}

// MARK: - RecipeTimeEstimationTests

@Suite(.serialized)
struct RecipeTimeEstimationTests {
  @Test("Returns nil for empty steps")
  func emptyStepsReturnsNil() {
    let result = RecipeTimeEstimation.estimate(steps: [])
    #expect(result == nil)
  }

  @Test("Single step with no time uses 1-3 min fallback")
  func singleStepNoTime() {
    let step = makeStep(stepID: "s1", estimatedTime: (nil, nil))
    let result = RecipeTimeEstimation.estimate(steps: [step])
    #expect(result != nil)
    #expect(result?.minSeconds == 60)
    #expect(result?.maxSeconds == 180)
  }

  @Test("Single step with min only uses min for both bounds")
  func singleStepMinOnly() {
    let step = makeStep(stepID: "s1", estimatedTime: (120, nil))
    let result = RecipeTimeEstimation.estimate(steps: [step])
    #expect(result != nil)
    #expect(result?.minSeconds == 120)
    #expect(result?.maxSeconds == 120)
  }

  @Test("Single step with max only uses max for both bounds")
  func singleStepMaxOnly() {
    let step = makeStep(stepID: "s1", estimatedTime: (nil, 300))
    let result = RecipeTimeEstimation.estimate(steps: [step])
    #expect(result != nil)
    #expect(result?.minSeconds == 300)
    #expect(result?.maxSeconds == 300)
  }

  @Test("Single step with both min and max uses exact values")
  func singleStepBothMinMax() {
    let step = makeStep(stepID: "s1", estimatedTime: (60, 120))
    let result = RecipeTimeEstimation.estimate(steps: [step])
    #expect(result != nil)
    #expect(result?.minSeconds == 60)
    #expect(result?.maxSeconds == 120)
  }

  @Test("Multiple steps with time are summed correctly")
  func multipleStepsWithTime() {
    let steps = [
      makeStep(stepID: "s1", estimatedTime: (60, 120)),
      makeStep(stepID: "s2", estimatedTime: (120, 180)),
      makeStep(stepID: "s3", estimatedTime: (30, 60)),
    ]
    let result = RecipeTimeEstimation.estimate(steps: steps)
    #expect(result != nil)
    #expect(result?.minSeconds == 210)   // 60 + 120 + 30
    #expect(result?.maxSeconds == 360)   // 120 + 180 + 60
  }

  @Test("Multiple steps without time use N × (60-180) sec")
  func multipleStepsNoTime() {
    let steps = [
      makeStep(stepID: "s1", estimatedTime: (nil, nil)),
      makeStep(stepID: "s2", estimatedTime: (nil, nil)),
      makeStep(stepID: "s3", estimatedTime: (nil, nil)),
    ]
    let result = RecipeTimeEstimation.estimate(steps: steps)
    #expect(result != nil)
    #expect(result?.minSeconds == 180)   // 3 × 60
    #expect(result?.maxSeconds == 540)   // 3 × 180
  }

  @Test("Mixed steps aggregate correctly")
  func mixedSteps() {
    let steps = [
      makeStep(stepID: "s1", estimatedTime: (120, 180)),  // has time
      makeStep(stepID: "s2", estimatedTime: (nil, nil)),   // fallback 60-180
      makeStep(stepID: "s3", estimatedTime: (300, nil)),  // min only
    ]
    let result = RecipeTimeEstimation.estimate(steps: steps)
    #expect(result != nil)
    #expect(result?.minSeconds == 480)   // 120 + 60 + 300
    #expect(result?.maxSeconds == 660)   // 180 + 180 + 300
  }

  @Test("Format displays single value when min equals max")
  func formatSingleValue() {
    let formatted = RecipeTimeEstimation.format(minSeconds: 600, maxSeconds: 600)
    #expect(formatted == "10 min")
  }

  @Test("Format displays range when min differs from max")
  func formatRange() {
    let formatted = RecipeTimeEstimation.format(minSeconds: 300, maxSeconds: 900)
    #expect(formatted == "5–15 min")
  }

  @Test("Format rounds down seconds to minutes")
  func formatRoundsDown() {
    let formatted = RecipeTimeEstimation.format(minSeconds: 89, maxSeconds: 119)
    #expect(formatted == "1 min")  // Both round to 1 min, shown as single value
  }

  @Test("Format shows range when seconds round to different minutes")
  func formatRoundsToDifferentMinutes() {
    let formatted = RecipeTimeEstimation.format(minSeconds: 61, maxSeconds: 179)
    #expect(formatted == "1–2 min")
  }

  @Test("Format uses hours for values over 60 minutes")
  func formatUsesHoursForLargeValues() {
    let formatted = RecipeTimeEstimation.format(minSeconds: 4020, maxSeconds: 7200)
    #expect(formatted == "1 hr 7 min – 2 hr")
  }

  @Test("Format caps absurdly large max at 24+ hr")
  func formatCapsLargeMax() {
    let formatted = RecipeTimeEstimation.format(minSeconds: 4020, maxSeconds: 261_540)
    #expect(formatted == "1 hr 7 min – 24+ hr")
  }

  @Test("FormatStepTime returns nil when step has no time")
  func formatStepTimeReturnsNilWhenEmpty() {
    let range = Common_OptionalUint32Range()
    #expect(RecipeTimeEstimation.formatStepTime(range) == nil)
  }

  @Test("FormatStepTime formats step with both min and max")
  func formatStepTimeWithBoth() {
    var range = Common_OptionalUint32Range()
    range.min = 300
    range.max = 600
    #expect(RecipeTimeEstimation.formatStepTime(range) == "5–10 min")
  }
}
