package main

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// On fait ici :
// 1. Calcul des gradients X et Y (Sobel)
// 2. Calcul de la magnitude
// 3. Seuillage (Threshold) pour ne garder que les bords forts
func CannyEdgeDetection(img *image.Gray, threshold float64, numGoroutines int) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	out := image.NewGray(bounds)

	// Kernel de Sobel
	gx := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	gy := [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	var wg sync.WaitGroup
	rowsPerGoroutine := height / numGoroutines
	if rowsPerGoroutine == 0 {
		rowsPerGoroutine = 1
	}

	// Lancement des workers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		startY := i * rowsPerGoroutine
		endY := startY + rowsPerGoroutine
		if i == numGoroutines-1 {
			endY = height // Le dernier worker prend le reste
		}

		go func(yStart, yEnd int) {
			defer wg.Done()

			// On évite les bords de l'image (1 pixel) car le noyau fait 3x3
			if yStart == 0 {
				yStart = 1
			}
			if yEnd == height {
				yEnd = height - 1
			}

			for y := yStart; y < yEnd; y++ {
				for x := 1; x < width-1; x++ {
					var sumX, sumY int

					// Convolution
					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {
							pixelVal := int(img.GrayAt(x+kx, y+ky).Y)
							sumX += pixelVal * gx[ky+1][kx+1]
							sumY += pixelVal * gy[ky+1][kx+1]
						}
					}

					// Magnitude du gradient : G = sqrt(Gx² + Gy²)
					magnitude := math.Sqrt(float64(sumX*sumX + sumY*sumY))

					// Seuillage simple (remplace l'hystérésis complet pour la vitesse)
					val := uint8(0)
					if magnitude > threshold {
						val = 255
					}
					out.SetGray(x, y, color.Gray{Y: val})
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return out
}
