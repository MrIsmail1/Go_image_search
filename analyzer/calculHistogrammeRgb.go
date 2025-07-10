package analyzer

import "image"

/*
	===== CALCUL HISTOGRAMME RGB =====

Calcule l'histogramme des composantes R, G et B pour une image.
Chaque couleur est répartie sur `Bins` intervalles.

Retour :
- Un dictionnaire contenant trois histogrammes (r, g, b)
*/
func computeHistogramRGB(img image.Image) map[string][]int {
	rHist, gHist, bHist := make([]int, Bins), make([]int, Bins), make([]int, Bins)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rHist[int(r>>8)*Bins/256]++
			gHist[int(g>>8)*Bins/256]++
			bHist[int(b>>8)*Bins/256]++
		}
	}
	return map[string][]int{"r": rHist, "g": gHist, "b": bHist}
}
