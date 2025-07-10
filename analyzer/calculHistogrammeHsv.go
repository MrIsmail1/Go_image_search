package analyzer

import "image"

/*
	===== CALCUL HISTOGRAMME HSV =====

Calcule l'histogramme pour la teinte (H), saturation (S) et valeur (V).
Les valeurs sont converties depuis RGB.

Retour :
- Un dictionnaire contenant trois histogrammes (h, s, v)
*/
func computeHistogramHSV(img image.Image) map[string][]int {
	hHist, sHist, vHist := make([]int, Bins), make([]int, Bins), make([]int, Bins)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			h, s, v := rgbToHsv(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			hIdx := int(h * float64(Bins) / 360)
			sIdx := int(s * float64(Bins))
			vIdx := int(v * float64(Bins))

			// Sécuriser les index
			if hIdx >= Bins {
				hIdx = Bins - 1
			}
			if sIdx >= Bins {
				sIdx = Bins - 1
			}
			if vIdx >= Bins {
				vIdx = Bins - 1
			}

			hHist[hIdx]++
			sHist[sIdx]++
			vHist[vIdx]++
		}
	}
	return map[string][]int{"h": hHist, "s": sHist, "v": vHist}
}
