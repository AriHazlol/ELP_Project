package main

import (
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "image/png"
    "os"
    "sync"
    "time"
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

func main() {
    
    // Essaie d'abord road.png, puis road.jpg
    var filename string
    var file *os.File
    var err error
    
    // Essaie road.png
    fmt.Println("Recherche de road.png...")
    file, err = os.Open("road.png")
    if err == nil {
        filename = "road.png"
        fmt.Println("✓ road.png trouvé")
    } else {
        // Essaie road.jpg
        file, err = os.Open("road.jpg")
        if err == nil {
            filename = "road.jpg"
            fmt.Println("✓ road.jpg trouvé")
        } else {
            fmt.Println("✗ Aucun fichier trouvé!")
            return
        }
    }
    defer file.Close()
    
    // Décoder l'image
    fmt.Println("Décodage de l'image...")
    var img image.Image
    
    // Devine le format par l'extension
    if len(filename) > 4 && filename[len(filename)-4:] == ".png" {
        img, err = png.Decode(file)
        if err != nil {
            fmt.Printf("Erreur PNG: %v\n", err)
            return
        }
    } else {
        img, err = jpeg.Decode(file)
        if err != nil {
            fmt.Printf("Erreur JPEG: %v\n", err)
            return
        }
    }
    
    // Convertir en RGBA
    fmt.Println("Conversion en RGBA...")
    bounds := img.Bounds()
    rgbaImg := image.NewRGBA(bounds)
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            rgbaImg.Set(x, y, img.At(x, y))
        }
    }
    
    fmt.Printf("Image: %d x %d pixels\n", bounds.Dx(), bounds.Dy())
    
    // Demander combien de goroutines
    var numGoroutines int
    fmt.Print("\nCombien de goroutines utiliser ? (1, 2, 4, 8): ")
    fmt.Scan(&numGoroutines)
    
    if numGoroutines < 1 {
        numGoroutines = 1
    }
    
    // Appliquer le filtre
    fmt.Printf("Application du filtre avec %d goroutine(s)...\n", numGoroutines)
    start := time.Now()
    
    if numGoroutines == 1 {
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                c := rgbaImg.RGBAAt(x, y)
                gray := uint8(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B))
                rgbaImg.SetRGBA(x, y, color.RGBA{gray, gray, gray, c.A})
            }
        }
    } else {
        ConcurrentGrayscale(rgbaImg, numGoroutines)
    }
    
    elapsed := time.Since(start)
    
    // Sauvegarder
    outputFile := "road_gris.png"
    fmt.Printf("Sauvegarde de %s...\n", outputFile)
    
    outFile, err := os.Create(outputFile)
    if err != nil {
        fmt.Printf("Erreur: %v\n", err)
        return
    }
    defer outFile.Close()
    
    png.Encode(outFile, rgbaImg)
    
    // Résultat
    fmt.Printf("Image traitée: %d x %d pixels\n", bounds.Dx(), bounds.Dy())
    fmt.Printf("Temps de traitement: %v\n", elapsed)
    fmt.Printf("Fichier créé: %s\n", outputFile)
}
