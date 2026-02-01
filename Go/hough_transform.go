package main

import (
	"image"
	"math"
)

type Line struct {
	Rho   float64 // distance à l'origine (dans le repère image utilisé)
	Theta float64 // angle (en radians) dans [0, pi)
	Votes int     // nombre de votes (force de la détection)
}

func HoughLines(edges *image.Gray, edgeThresh uint8, thetaSteps int, voteThresh int, maxLines int) []Line {
	b := edges.Bounds()
	w := b.Dx()
	h := b.Dy()

	// rho correspond à une distance, son amplitude max est approximativement la diagonale de l'image.
	// On prend rhoMax = ceil(sqrt(w^2 + h^2))
	rhoMax := int(math.Ceil(math.Hypot(float64(w), float64(h))))

	// rho peut être négatif ou positif selon la représentation.
	// Pour indexer un tableau, on décale : rhoIndex = round(rho) + rhoMax
	rhoRange := 2*rhoMax + 1 // rhoIndex ∈ [0, 2*rhoMax]

	// - thetaIndex ∈ [0, thetaSteps-1]
	// - rhoIndex   ∈ [0, rhoRange-1]
	acc := make([][]int, thetaSteps)
	for t := 0; t < thetaSteps; t++ {
		acc[t] = make([]int, rhoRange)
	}

	cosT := make([]float64, thetaSteps)
	sinT := make([]float64, thetaSteps)
	for t := 0; t < thetaSteps; t++ {
		theta := float64(t) * math.Pi / float64(thetaSteps) // theta ∈ [0, pi)
		cosT[t] = math.Cos(theta)
		sinT[t] = math.Sin(theta)
	}

	// Remarque : on travaille dans un repère local à l'image : x,y ∈ [0..w-1], [0..h-1]
	// Donc on soustrait b.Min.X / b.Min.Y.
	for y := b.Min.Y; y < b.Max.Y; y++ {
		yf := float64(y - b.Min.Y)
		for x := b.Min.X; x < b.Max.X; x++ {

			// Ignorer les pixels non-bords
			if edges.GrayAt(x, y).Y < edgeThresh {
				continue
			}

			xf := float64(x - b.Min.X)

			// Pour ce pixel (x,y), on parcourt tous les angles theta
			for t := 0; t < thetaSteps; t++ {
				// Formule de Hough pour une droite :
				// rho = x*cos(theta) + y*sin(theta)
				rho := xf*cosT[t] + yf*sinT[t]

				// Conversion rho (float) -> index entier dans l'accumulateur
				rIdx := int(math.Round(rho)) + rhoMax

				// Sécurité (normalement rIdx est dans [0, rhoRange))
				if rIdx >= 0 && rIdx < rhoRange {
					acc[t][rIdx]++
				}
			}
		}
	}

	candidates := make([]Line, 0)
	for t := 0; t < thetaSteps; t++ {
		theta := float64(t) * math.Pi / float64(thetaSteps)
		for rIdx := 0; rIdx < rhoRange; rIdx++ {
			v := acc[t][rIdx]
			if v < voteThresh {
				continue
			}
			rho := float64(rIdx - rhoMax)
			candidates = append(candidates, Line{
				Rho:   rho,
				Theta: theta,
				Votes: v,
			})
		}
	}

	for i := 0; i < len(candidates); i++ {
		best := i
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].Votes > candidates[best].Votes {
				best = j
			}
		}
		candidates[i], candidates[best] = candidates[best], candidates[i]
	}

	seuilRho := 10.0                  // ~10 pixels
	seuilTheta := 3.0 * math.Pi / 180 // ~3 degrés

	out := make([]Line, 0, maxLines)
	for _, ln := range candidates {
		tooClose := false
		for _, kept := range out {
			if math.Abs(ln.Rho-kept.Rho) < seuilRho && math.Abs(ln.Theta-kept.Theta) < seuilTheta {
				tooClose = true
				break
			}
		}
		if tooClose {
			continue
		}
		out = append(out, ln)
		if len(out) >= maxLines {
			break
		}
	}

	return out
}
