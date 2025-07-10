package analyzer

import "image"

/*
	===== CALCUL MOYENNE DES COULEURS =====

Calcule la moyenne des composantes rouge, vert et bleu sur l'image.

Retour :
- Vecteur [R, G, B] en float64
*/
func computeMeanColor(img image.Image) [3]float64 {
	var rSum, gSum, bSum float64
	bounds := img.Bounds()
	total := float64(bounds.Dx() * bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rSum += float64(r >> 8)
			gSum += float64(g >> 8)
			bSum += float64(b >> 8)
		}
	}
	return [3]float64{rSum / total, gSum / total, bSum / total}
}
