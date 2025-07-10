package analyzer

import "math"

/*
	===== CONVERSION RGB -> HSV =====

Convertit une couleur RGB en HSV :
- H (teinte) ∈ [0, 360)
- S (saturation) ∈ [0, 1]
- V (valeur) ∈ [0, 1]

Retour :
- h, s, v en float64
*/
func rgbToHsv(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	delta := max - min

	var h float64
	if delta != 0 {
		switch max {
		case rf:
			h = math.Mod(((gf - bf) / delta), 6)
		case gf:
			h = ((bf - rf) / delta) + 2
		case bf:
			h = ((rf - gf) / delta) + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	s := 0.0
	if max != 0 {
		s = delta / max
	}
	v := max
	return h, s, v
}
