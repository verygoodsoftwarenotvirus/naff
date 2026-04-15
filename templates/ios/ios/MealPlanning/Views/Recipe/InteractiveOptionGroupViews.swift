//
//  InteractiveOptionGroupViews.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftProtobuf
import SwiftUI

// MARK: - Interactive Ingredient Option Group View

struct InteractiveIngredientOptionGroupView: View {
  let group: OptionGroupAggregate
  @Binding var selectedOptionIndex: UInt32
  var scale: Float = 1.0

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      HStack {
        Text("one of:")
          .font(.caption)
          .foregroundColor(.secondary)

        if let sourceRecipeName = group.sourceRecipeName {
          Text("(from \(sourceRecipeName))")
            .font(.caption2)
            .foregroundColor(.secondary)
        }
      }
      .padding(.leading, 16)
      .padding(.horizontal)

      ForEach(group.options) { option in
        Button(
          action: {
            selectedOptionIndex = option.optionIndex
          },
          label: {
            HStack(alignment: .top, spacing: 6) {
              Image(
                systemName: selectedOptionIndex == option.optionIndex
                  && selectedOptionIndex != UInt32.max
                  ? "checkmark.circle.fill" : "circle"
              )
              .font(.caption)
              .foregroundColor(
                selectedOptionIndex == option.optionIndex && selectedOptionIndex != UInt32.max
                  ? .green : .gray)

              if let quantityText = option.aggregated.quantityText(scale: scale) {
                Text(quantityText)
                  .font(.caption)
                  .foregroundColor(.secondary)
              }

              VStack(alignment: .leading, spacing: 2) {
                Text(option.ingredient.name)
                  .font(.caption)
                  .foregroundColor(.secondary)
                if !option.aggregated.quantityNotes.isEmpty {
                  Text(option.aggregated.quantityNotes)
                    .font(.caption2)
                    .foregroundColor(.secondary)
                }
              }
            }
            .padding(.leading, 16)
            .padding(.horizontal)
          }
        )
        .buttonStyle(.plain)
      }
    }
    .padding(.vertical, 4)
  }
}

// MARK: - Interactive Instrument Option Group View

struct InteractiveInstrumentOptionGroupView: View {
  let group: InstrumentOptionGroupAggregate
  @Binding var selectedOptionIndex: UInt32
  var scale: Float = 1.0

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      HStack {
        Text("one of:")
          .font(.caption)
          .foregroundColor(.secondary)

        if let sourceRecipeName = group.sourceRecipeName {
          Text("(from \(sourceRecipeName))")
            .font(.caption2)
            .foregroundColor(.secondary)
        }
      }
      .padding(.leading, 16)
      .padding(.horizontal)

      ForEach(group.options) { option in
        Button(
          action: {
            selectedOptionIndex = option.optionIndex
          },
          label: {
            HStack(spacing: 6) {
              Image(
                systemName: selectedOptionIndex == option.optionIndex
                  && selectedOptionIndex != UInt32.max
                  ? "checkmark.circle.fill" : "circle"
              )
              .font(.caption)
              .foregroundColor(
                selectedOptionIndex == option.optionIndex && selectedOptionIndex != UInt32.max
                  ? .green : .gray)

              Text(option.instrument.name)
                .font(.caption)
                .foregroundColor(.secondary)
            }
            .padding(.leading, 16)
            .padding(.horizontal)
          }
        )
        .buttonStyle(.plain)
      }
    }
    .padding(.vertical, 4)
  }
}

// MARK: - Interactive Vessel Option Group View

struct InteractiveVesselOptionGroupView: View {
  let group: VesselOptionGroupAggregate
  @Binding var selectedOptionIndex: UInt32
  var scale: Float = 1.0

  var body: some View {
    VStack(alignment: .leading, spacing: 4) {
      HStack {
        Text("one of:")
          .font(.caption)
          .foregroundColor(.secondary)

        if let sourceRecipeName = group.sourceRecipeName {
          Text("(from \(sourceRecipeName))")
            .font(.caption2)
            .foregroundColor(.secondary)
        }
      }
      .padding(.leading, 16)
      .padding(.horizontal)

      ForEach(group.options) { option in
        Button(
          action: {
            selectedOptionIndex = option.optionIndex
          },
          label: {
            HStack(spacing: 6) {
              Image(
                systemName: selectedOptionIndex == option.optionIndex
                  && selectedOptionIndex != UInt32.max
                  ? "checkmark.circle.fill" : "circle"
              )
              .font(.caption)
              .foregroundColor(
                selectedOptionIndex == option.optionIndex && selectedOptionIndex != UInt32.max
                  ? .green : .gray)

              Text(option.vessel.name)
                .font(.caption)
                .foregroundColor(.secondary)
            }
            .padding(.leading, 16)
            .padding(.horizontal)
          }
        )
        .buttonStyle(.plain)
      }
    }
    .padding(.vertical, 4)
  }
}
