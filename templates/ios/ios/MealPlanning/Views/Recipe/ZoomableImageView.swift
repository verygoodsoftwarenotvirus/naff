//
//  ZoomableImageView.swift
//  ios
//
//  A SwiftUI wrapper around UIScrollView that provides pinch-to-zoom and pan
//  for viewing images (e.g. Mermaid diagrams) in detail.
//

import SwiftUI
import UIKit

/// A view that displays an image with pinch-to-zoom and pan gestures.
struct ZoomableImageView: UIViewRepresentable {
  let image: UIImage
  var minimumZoomScale: CGFloat = 1.0
  var maximumZoomScale: CGFloat = 4.0

  func makeUIView(context: Context) -> UIScrollView {
    let scrollView = UIScrollView()
    scrollView.delegate = context.coordinator
    scrollView.minimumZoomScale = minimumZoomScale
    scrollView.maximumZoomScale = maximumZoomScale
    scrollView.showsHorizontalScrollIndicator = true
    scrollView.showsVerticalScrollIndicator = true
    scrollView.bouncesZoom = true

    let imageView = UIImageView(image: image)
    imageView.contentMode = .scaleAspectFit
    imageView.isUserInteractionEnabled = true
    imageView.tag = 1
    imageView.frame = CGRect(origin: .zero, size: image.size)

    scrollView.addSubview(imageView)
    scrollView.contentSize = image.size

    return scrollView
  }

  func updateUIView(_ scrollView: UIScrollView, context: Context) {
    guard let imageView = scrollView.viewWithTag(1) as? UIImageView else { return }
    imageView.image = image
    imageView.frame = CGRect(origin: .zero, size: image.size)
    scrollView.contentSize = image.size
  }

  func makeCoordinator() -> Coordinator {
    Coordinator()
  }

  class Coordinator: NSObject, UIScrollViewDelegate {
    func viewForZooming(in scrollView: UIScrollView) -> UIView? {
      scrollView.viewWithTag(1)
    }
  }
}
