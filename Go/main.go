package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw" // Nécessaire pour dessiner les lignes
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"time"
)

// Helper pour convertir RGBA vers Gray (nécessaire car vos filtres utilisent des types différents)
func convertToGrayStruct(img *image.RGBA) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// On récupère juste le canal Y (luminance) car l'image est déjà visuellement grise
			c := img.RGBAAt(x, y)
			gray.SetGray(x, y, color.Gray{Y: c.R})
		}
	}
	return gray
}

func drawLines(img *image.RGBA, lines []Line) {
	if len(lines) == 0 {
		fmt.Println("Aucune ligne détectée.")
		return
	}

	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255} // Bleu pour la ligne du milieu
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	halfHeight := height / 2

	var leftLines, rightLines []Line

	// 1. Séparation des lignes par angle (Theta)
	// En Hough, Theta est en radians.
	// Environ < 1.5 rad (85°) = penche à droite, > 1.6 rad (95°) = penche à gauche
	for _, line := range lines {
		if line.Theta < 1.5 {
			leftLines = append(leftLines, line)
		} else if line.Theta > 1.6 {
			rightLines = append(rightLines, line)
		}
	}

	// 2. Fonction interne pour calculer la moyenne et dessiner
	renderAverage := func(group []Line, col color.RGBA) {
		if len(group) == 0 {
			return
		}

		var sRho, sTheta float64
		for _, l := range group {
			sRho += l.Rho
			sTheta += l.Theta
		}
		avgR := sRho / float64(len(group))
		avgT := sTheta / float64(len(group))

		cosT, sinT := math.Cos(avgT), math.Sin(avgT)

		for x := 0; x < width; x++ {
			y := (avgR - float64(x)*cosT) / sinT
			if int(y) >= halfHeight && int(y) < height {
				// Dessin avec épaisseur
				for dy := -1; dy <= 1; dy++ {
					if int(y)+dy < height && int(y)+dy >= 0 {
						img.SetRGBA(x, int(y)+dy, col)
					}
				}
			}
		}
	}

	// 3. Tracer les deux moyennes (Milieu en Bleu, Droite en Rouge)
	renderAverage(leftLines, blue)
	renderAverage(rightLines, red)
}

func main() {
	fmt.Println("=== Démarrage du Pipeline de Vision ===")

	// 1. Chargement de l'image
	imgPath := "road.png"
	f, err := os.Open(imgPath)
	if err != nil {
		// si png n'existe pas
		imgPath = "road.jpg"
		f, err = os.Open(imgPath)
		if err != nil {
			fmt.Println("Erreur: Impossible d'ouvrir road.png ou road.jpg")
			return
		}
	}
	defer f.Close()

	var src image.Image
	if filepath.Ext(imgPath) == ".png" {
		src, _ = png.Decode(f)
	} else {
		src, _ = jpeg.Decode(f)
	}

	// Conversion en RGBA
	bounds := src.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, src, bounds.Min, draw.Src)

	// Paramètres
	numGoroutines := 4 // Peut être ajusté
	fmt.Printf("Traitement de %s (%dx%d) avec %d goroutines\n", imgPath, bounds.Dx(), bounds.Dy(), numGoroutines)

	startTotal := time.Now()

	// ÉTAPE 1 : Niveaux de Gris (utilise grey_filter.go)

	fmt.Println("- Étape 1 : Grayscale...")
	if numGoroutines == 1 {
		ConcurrentGrayscale(img, 1)
	} else {
		ConcurrentGrayscale(img, numGoroutines)
	}

	// ÉTAPE 2 : Region of Interest (ROI) (utilise programme_roi.go)
	fmt.Println("- Étape 2 : ROI (Triangle)...")
	applyFixedTriangleROI(img, numGoroutines)

	// --- Conversion RGBA -> Gray pour la suite des algos ---
	grayImg := convertToGrayStruct(img)

	// ÉTAPE 3 : Flou Gaussien
	fmt.Println("- Étape 3 : Flou Gaussien...")
	blurredImg := GaussianBlurParallel(grayImg)

	// ÉTAPE 4 : Détection de contours Canny

	fmt.Println("- Étape 4 : Canny Edge Detection...")

	edgeImg := CannyEdgeDetection(blurredImg, 50.0, numGoroutines)

	// ÉTAPE 5 : Transformée de Hough*

	fmt.Println("- Étape 5 : Hough Transform...")

	// Paramètres : edgeThresh, thetaSteps, voteThresh, maxLines
	// edgeThresh doit être bas car Canny a déjà "nettoyé"
	lines := HoughLines(edgeImg, 1, 180, 100, 10)
	fmt.Printf("  > %d lignes détectées\n", len(lines))

	// ÉTAPE 6 : Résultat final

	// On recharge l'image couleur originale pour dessiner les lignes (bordures de la route) dessus
	fOriginal, _ := os.Open(imgPath)
	srcOriginal, _, err := image.Decode(fOriginal)
	fOriginal.Close()

	finalImg := image.NewRGBA(bounds)
	draw.Draw(finalImg, bounds, srcOriginal, bounds.Min, draw.Src)

	drawLines(finalImg, lines)

	elapsed := time.Since(startTotal)
	fmt.Printf("Terminé en %v\n", elapsed)

	outFile, _ := os.Create("resultat_final.png")
	defer outFile.Close()
	png.Encode(outFile, finalImg)
	fmt.Println("Image sauvegardée : resultat_final.png")
}
