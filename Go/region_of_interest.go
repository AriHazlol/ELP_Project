package main

import (
	"image"
	"image/color"
	"sync"
)

func applyFixedTriangleROI(img *image.RGBA, numGoroutines int) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// coord ajustables par rapport à chaque img
	bottomLeftX := 350
	bottomRightX := 1800
	topCenterX := 1100
	topCenterY := 500

	var wg sync.WaitGroup
	blockHeight := height / numGoroutines

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		startY := g * blockHeight
		endY := startY + blockHeight
		if g == numGoroutines-1 {
			endY = height
		}

		go func(startY, endY int) {
			defer wg.Done()

			for y := startY; y < endY; y++ {
				var inTriangle bool
				var leftEdge, rightEdge int

				if y >= topCenterY {
					// Dans le triangle (y entre topCenterY et bas)
					t := float64(y-topCenterY) / float64(height-topCenterY)

					leftEdge = topCenterX + int(t*float64(bottomLeftX-topCenterX))
					rightEdge = topCenterX + int(t*float64(bottomRightX-topCenterX))

					// Assurer leftEdge < rightEdge
					if leftEdge > rightEdge {
						leftEdge, rightEdge = rightEdge, leftEdge
					}

					inTriangle = true
				}

				// Traiter chaque pixel de cette ligne
				for x := 0; x < width; x++ {
					if inTriangle && x >= leftEdge && x <= rightEdge {
						// Dans le triangle → garder tel quel
						continue
					} else {
						// Hors triangle → noir
						img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
					}
				}
			}
		}(startY, endY)
	}

	wg.Wait()
}
