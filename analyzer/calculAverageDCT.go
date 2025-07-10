package analyzer

/*
	===== MOYENNE DES COEFFICIENTS DCT =====

Calcule la moyenne des 64 premiers coefficients DCT,
en ignorant le premier coefficient (DC term).

Retour :
- Moyenne des coefficients DCT
*/
func averageDCT(dct [][]float64) float64 {
	total := 0.0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if !(x == 0 && y == 0) {
				total += dct[y][x]
			}
		}
	}
	return total / 63
}
