package analyzer

import (
	"image"
	"image/draw"
	"math"
)

/*
	===== DÉTECTION DE CONTOURS (Signature de forme) =====

Cette fonction applique un filtre Sobel simple pour détecter les contours
dans une image convertie en niveaux de gris, puis calcule une signature
basée sur la densité de ces contours.

Retour :
- Ratio entre le nombre de pixels "de contours" et la surface totale de l'image
*/
func computeShapeSignature(img image.Image) float64 {
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.Point{}, draw.Src)

	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var contourPixels int

	// Matrices de Sobel (X et Y)
	kernelX := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	kernelY := [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	// Calcul du gradient (Sobel)
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var gx, gy int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					p := gray.GrayAt(x+kx, y+ky).Y
					gx += int(p) * kernelX[ky+1][kx+1]
					gy += int(p) * kernelY[ky+1][kx+1]
				}
			}
			gradient := math.Sqrt(float64(gx*gx + gy*gy))
			if gradient > 100 { // seuil simple
				contourPixels++
			}
		}
	}

	totalPixels := (width - 2) * (height - 2) // bords ignorés
	return float64(contourPixels) / float64(totalPixels)
}