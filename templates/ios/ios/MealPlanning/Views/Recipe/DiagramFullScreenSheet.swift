//
//  DiagramFullScreenSheet.swift
//  ios
//
//  Full-screen sheet for viewing Mermaid diagrams with zoom and pan.
//

import SwiftUI

/// A full-screen sheet that displays a Mermaid diagram with zoom/pan and a close button.
struct DiagramFullScreenSheet: View {
  let mermaidSource: String
  let title: String
  let onDismiss: () -> Void

  var body: some View {
    NavigationStack {
      MermaidDiagramView(source: mermaidSource)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(Color(.systemBackground))
        .navigationTitle(title)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
          ToolbarItem(placement: .cancellationAction) {
            Button("Done") {
              onDismiss()
            }
          }
        }
    }
  }
}
