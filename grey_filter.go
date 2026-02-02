package main

import (
	"image"
	"image/color"
	"sync"
)

func processBlock(img *image.RGBA, startY, endY int, wg *sync.WaitGroup) {
	defer wg.Done()

	bounds := img.Bounds()
	for y := startY; y < endY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			gray := uint8(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B))
			img.SetRGBA(x, y, color.RGBA{gray, gray, gray, c.A})
		}
	}
}

func ConcurrentGrayscale(img *image.RGBA, numGoroutines int) *image.RGBA {
	bounds := img.Bounds()
	height := bounds.Max.Y - bounds.Min.Y
	blockHeight := height / numGoroutines

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		startY := bounds.Min.Y + i*blockHeight
		endY := startY + blockHeight
		if i == numGoroutines-1 {
			endY = bounds.Max.Y
		}
		go processBlock(img, startY, endY, &wg)
	}

	wg.Wait()
	return img
}
