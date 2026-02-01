package main

import (
	"image"
	"image/color"
	"runtime"
	"sync"
)

var kernel5x5 = [][]float64{
	{1, 4, 7, 4, 1},
	{4, 16, 26, 16, 4},
	{7, 26, 41, 26, 7},
	{4, 16, 26, 16, 4},
	{1, 4, 7, 4, 1},
}

func GaussianBlurParallel(src *image.Gray) *image.Gray {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	dest := image.NewGray(bounds)

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		startY := i * (height / numWorkers)
		endY := (i + 1) * (height / numWorkers)
		if i == numWorkers-1 {
			endY = height
		}

		wg.Add(1)
		go func(sY, eY int) {
			defer wg.Done()
			for y := sY; y < eY; y++ {
				for x := 0; x < width; x++ {
					// Ignorer les bords pour simplifier (2 pixels de marge pour kernel 5x5)
					if x < 2 || x > width-3 || y < 2 || y > height-3 {
						dest.SetGray(x, y, src.GrayAt(x, y))
						continue
					}

					var sum float64
					for ky := -2; ky <= 2; ky++ {
						for kx := -2; kx <= 2; kx++ {
							val := src.GrayAt(x+kx, y+ky).Y
							sum += float64(val) * kernel5x5[ky+2][kx+2]
						}
					}
					dest.SetGray(x, y, color.Gray{Y: uint8(sum / 273.0)})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return dest
}
