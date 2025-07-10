package analyzer

import (
	"fmt"
	drawx "golang.org/x/image/draw"
	"image"
	"image/draw"
)

/*
	===== GÉNÉRATION DU pHash (Perceptual Hash) =====

Produit un hash compact (64 bits) représentant l'aspect visuel global de l'image.

Méthode :
1. Redimensionnement en 32x32
2. Conversion en niveaux de gris
3. Application d'une DCT (Discrete Cosine Transform)
4. Comparaison aux coefficients moyens

Retour :
- Hash sous forme d'une chaîne hexadécimale (16 caractères)
*/
func generatePHash(img image.Image) string {
	size := 32
	gray := image.NewGray(image.Rect(0, 0, size, size))
	drawx.ApproxBiLinear.Scale(gray, gray.Bounds(), img, img.Bounds(), draw.Over, nil)

	dctVals := dct2D(gray)
	avg := averageDCT(dctVals)

	var hash uint64
	index := 0

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if dctVals[y][x] > avg {
				hash |= 1 << index
			}
			index++
		}
	}
	return fmt.Sprintf("%016x", hash)
}
