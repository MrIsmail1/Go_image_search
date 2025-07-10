package analyzer

import (
	"image"
	"math"
)

/*
	===== DIFFERENCE ABSOLUE ENTRE 2 VALEURS UINT8 =====

Retour :
- Valeur absolue de la différence
*/
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

/*
	===== TRANSFORMÉE EN COSINUS 2D (DCT) =====

Applique une DCT bidimensionnelle à une image en niveaux de gris.
Utilisée dans la génération du pHash.

Retour :
- Matrice des coefficients DCT
*/
func dct2D(img *image.Gray) [][]float64 {
	N := 32
	dct := make([][]float64, N)
	for i := range dct {
		dct[i] = make([]float64, N)
	}

	for u := 0; u < N; u++ {
		for v := 0; v < N; v++ {
			var sum float64
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					pixel := float64(img.GrayAt(x, y).Y)
					sum += pixel *
						math.Cos((2*float64(x)+1)*float64(u)*math.Pi/float64(2*N)) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/float64(2*N))
				}
			}
			cu, cv := 1.0, 1.0
			if u == 0 {
				cu = 1 / math.Sqrt(2)
			}
			if v == 0 {
				cv = 1 / math.Sqrt(2)
			}
			dct[u][v] = 0.25 * cu * cv * sum
		}
	}
	return dct
}
