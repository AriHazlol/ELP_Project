package main

import (
	"image"
	"image/color"
)

// GaussianBlur3x3 applique un flou gaussien simple (approximation) sur une image en niveaux de gris.
//
// Noyau (kernel) 3x3 utilisé :
//   1 2 1
//   2 4 2   le tout divisé par 16
//   1 2 1
//
// C’est un flou “léger” très standard : parfait avant Sobel pour réduire le bruit.
func GaussianBlur3x3(src *image.Gray) *image.Gray {
	b := src.Bounds()
	dst := image.NewGray(b)

	// Poids du noyau gaussien 3x3 (approximation).
	k := [3][3]int{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	const norm = 16 // somme des poids du noyau

	// On évite les bords (car on lit x±1, y±1).
	// Les pixels sur le contour seront laissés à 0 par défaut (noir),
	// ou tu peux les recopier depuis src si tu préfères.
	for y := b.Min.Y + 1; y < b.Max.Y-1; y++ {
		for x := b.Min.X + 1; x < b.Max.X-1; x++ {

			sum := 0

			// Convolution 3x3 centrée en (x, y)
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					// Pixel voisin (valeur de luminance 0..255)
					p := int(src.GrayAt(x+kx, y+ky).Y)
					// Poids correspondant dans le noyau (ky+1, kx+1)
					w := k[ky+1][kx+1]
					sum += p * w
				}
			}

			// Normalisation : moyenne pondérée
			blurred := sum / norm
			dst.SetGray(x, y, color.Gray{Y: uint8(blurred)})
		}
	}

	// Optionnel : gérer les bords en les recopiant depuis src
	// (souvent préférable pour éviter un contour noir).
	copyBorderGray(dst, src)

	return dst
}

// copyBorderGray recopie les pixels du bord (première/dernière ligne et colonne)
// de src vers dst. Cela évite d’avoir un cadre noir dû au fait qu’on ne calcule
// pas la convolution sur les bords.
func copyBorderGray(dst, src *image.Gray) {
	b := src.Bounds()

	// Lignes du haut et du bas
	for x := b.Min.X; x < b.Max.X; x++ {
		dst.SetGray(x, b.Min.Y, src.GrayAt(x, b.Min.Y))
		dst.SetGray(x, b.Max.Y-1, src.GrayAt(x, b.Max.Y-1))
	}

	// Colonnes de gauche et de droite
	for y := b.Min.Y; y < b.Max.Y; y++ {
		dst.SetGray(b.Min.X, y, src.GrayAt(b.Min.X, y))
		dst.SetGray(b.Max.X-1, y, src.GrayAt(b.Max.X-1, y))
	}
}
