//
//  WeekSelectionStepView.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import SwiftUI

struct WeekSelectionStepView: View {
  @Bindable var viewModel: CreateMealPlanViewModel
  @Environment(EventReporterService.self) private var eventReporterService

  private let shortDateFormatter: DateFormatter = {
    let formatter = DateFormatter()
    formatter.dateFormat = "M/d"
    return formatter
  }()

  var body: some View {
    VStack(alignment: .leading, spacing: 24) {
      Text("Select the days to plan")
        .font(.title2)
        .fontWeight(.semibold)

      HStack {
        Button(
          action: {
            eventReporterService.reporter.track(event: "meal_plan_week_prev", properties: [:])
            viewModel.goToPreviousWeek()
          },
          label: {
            Image(systemName: "chevron.left")
              .font(.title3)
          }
        )
        .disabled(viewModel.selectedWeekOffset == 0)

        Spacer()

        Text(weekRangeTitle)
          .font(.headline)

        Spacer()

        Button(
          action: {
            eventReporterService.reporter.track(event: "meal_plan_week_next", properties: [:])
            viewModel.goToNextWeek()
          },
          label: {
            Image(systemName: "chevron.right")
              .font(.title3)
          }
        )
        .disabled(viewModel.selectedWeekOffset >= 4)
      }
      .padding(.horizontal, 8)

      weekDayGrid

      Text("Tap to select or deselect. Drag across days to select a range.")
        .font(.caption)
        .foregroundColor(.secondary)

      Spacer(minLength: 24)

      Button(
        action: {
          eventReporterService.reporter.track(
            event: "meal_plan_week_continue",
            properties: ["selected_days_count": viewModel.selectedDates.count])
          viewModel.wizardStep = .mealAssignment
        },
        label: {
          Text("Continue to Assign Meals")
            .fontWeight(.semibold)
            .frame(maxWidth: .infinity)
            .padding()
            .background(viewModel.selectedDates.isEmpty ? Color.gray : Color.blue)
            .foregroundColor(.white)
            .cornerRadius(10)
        }
      )
      .disabled(viewModel.selectedDates.isEmpty)
    }
    .frame(maxWidth: .infinity, alignment: .leading)
  }

  private var weekRangeTitle: String {
    let days = viewModel.displayedWeekDays
    guard let first = days.first, let last = days.last else { return "Week" }
    let formatter = DateFormatter()
    formatter.dateFormat = "MMM d"
    return "\(formatter.string(from: first)) – \(formatter.string(from: last))"
  }

  private var weekDayGrid: some View {
    GeometryReader { geo in
      let width = geo.size.width
      let spacing: CGFloat = 6
      let totalSpacing = spacing * 6  // 6 gaps between 7 items
      let itemWidth = (width - totalSpacing) / 7

      HStack(spacing: spacing) {
        ForEach(Array(viewModel.displayedWeekDays.enumerated()), id: \.element) { _, date in
          dayChip(date: date)
            .frame(width: itemWidth)
        }
      }
      .contentShape(Rectangle())
      .highPriorityGesture(
        DragGesture(minimumDistance: 8)
          .onChanged { value in
            let startCol = dayIndexForX(
              value.startLocation.x, totalWidth: width, spacing: spacing, itemWidth: itemWidth)
            let currCol = dayIndexForX(
              value.location.x, totalWidth: width, spacing: spacing, itemWidth: itemWidth)
            let days = viewModel.displayedWeekDays
            guard !days.isEmpty, startCol >= 0, startCol < days.count, currCol >= 0,
              currCol < days.count
            else {
              return
            }
            let lowIndex = min(startCol, currCol)
            let highIndex = max(startCol, currCol)
            viewModel.setDateRangeSelection(from: days[lowIndex], to: days[highIndex])
          }
      )
    }
    .frame(height: 56)
  }

  private func dayIndexForX(
    _ locationX: CGFloat, totalWidth: CGFloat, spacing: CGFloat, itemWidth: CGFloat
  ) -> Int {
    guard totalWidth > 0 else { return 0 }
    // Each slot: itemWidth + spacing, except last has no trailing spacing
    let slotWidth = itemWidth + spacing
    let index = Int(locationX / slotWidth)
    return min(6, max(0, index))
  }

  private struct DayChipStyle {
    let backgroundColor: Color
    let foregroundColor: Color
    let strokeColor: Color
    let opacity: Double
  }

  private func dayChipStyle(isPlanable: Bool, isSelected: Bool, occupancy: DateOccupancy?)
    -> DayChipStyle
  {
    if isSelected {
      return DayChipStyle(
        backgroundColor: Color.blue.opacity(0.2),
        foregroundColor: .blue,
        strokeColor: .blue,
        opacity: 1
      )
    }
    switch occupancy {
    case .accepted:
      return DayChipStyle(
        backgroundColor: Color.red.opacity(0.2),
        foregroundColor: .red,
        strokeColor: .red.opacity(0.6),
        opacity: 1
      )
    case .proposed:
      return DayChipStyle(
        backgroundColor: Color.yellow.opacity(0.3),
        foregroundColor: .orange,
        strokeColor: Color.clear,
        opacity: 1
      )
    case .none:
      return DayChipStyle(
        backgroundColor: Color(.systemGray6),
        foregroundColor: isPlanable ? .primary : .secondary,
        strokeColor: Color.clear,
        opacity: isPlanable ? 1 : 0.5
      )
    }
  }

  private func dayChip(date: Date) -> some View {
    let isPlanable = viewModel.isDatePlanable(date)
    let isSelected = viewModel.isDateSelected(date)
    let occupancy = viewModel.dateOccupancy(for: date)
    // Allow tap when planable (select/deselect) or when selected (deselect only, e.g. after 6PM)
    let isInteractive = isPlanable || isSelected
    let label = compactWeekdayLabel(for: date)
    let dayNum = shortDateFormatter.string(from: date)
    let style = dayChipStyle(isPlanable: isPlanable, isSelected: isSelected, occupancy: occupancy)

    return Button(
      action: {
        viewModel.toggleDateSelection(date)
      },
      label: {
        VStack(spacing: 2) {
          Text(label)
            .font(.caption)
            .fontWeight(.semibold)
          Text(dayNum)
            .font(.system(size: 10))
            .fontWeight(.medium)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, 8)
        .background(style.backgroundColor)
        .foregroundColor(style.foregroundColor)
        .opacity(style.opacity)
        .cornerRadius(8)
        .overlay(
          RoundedRectangle(cornerRadius: 8)
            .stroke(style.strokeColor, lineWidth: 2)
        )
      }
    )
    .buttonStyle(.plain)
    .disabled(!isInteractive)
  }

  private func compactWeekdayLabel(for date: Date) -> String {
    let weekday = Calendar.current.component(.weekday, from: date)
    switch weekday {
    case 1: return "Su"
    case 2: return "M"
    case 3: return "T"
    case 4: return "W"
    case 5: return "Th"
    case 6: return "F"
    case 7: return "Sa"
    default: return "?"
    }
  }
}
