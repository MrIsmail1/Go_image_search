package hash

/*
===== MOYENNE DES COEFFICIENTS DCT =====

À QUOI ÇA SERT :
Calcule la moyenne des coefficients DCT pour créer un seuil de binarisation
dans l'algorithme du hash perceptuel (pHash). Cette moyenne sert de "ligne de partage"
pour transformer les coefficients en bits 0 ou 1.

CONTEXTE - QU'EST-CE QUE LA DCT :
- DCT = Transformée en Cosinus Discrète (comme dans JPEG)
- Transforme une image en coefficients de fréquences
- Les basses fréquences (coin sup-gauche) = structure générale
- Les hautes fréquences = détails fins

PRINCIPE DU SEUILLAGE :
- Coefficient > moyenne → bit = 1
- Coefficient ≤ moyenne → bit = 0
- Résultat : hash binaire de 64 bits représentant la structure

Paramètre :
- dct : matrice 32x32 des coefficients DCT calculée précédemment

Retour :
- Moyenne des 63 coefficients (excluant le terme DC)
*/
func averageDCT(dct [][]float64) float64 {
	total := 0.0

	// Parcours des 64 premiers coefficients (8x8 du coin supérieur gauche)
	// POURQUOI seulement 8x8 sur une matrice 32x32 ?
	// - Les basses fréquences (8x8) contiennent l'info structurelle principale
	// - 64 coefficients = 64 bits de hash (taille standard du pHash)
	// - Les coefficients plus loin sont des détails fins moins importants
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {

			// Exclusion du coefficient DC en position [0][0]
			// CONDITION : !(x == 0 && y == 0) = "PAS (première position)"
			// EFFET : Traite tous les coefficients sauf celui du coin sup-gauche
			if !(x == 0 && y == 0) {
				total += dct[y][x] // Accumulation de ce coefficient
			}
		}
	}

	// Division par 63 pour obtenir la moyenne
	// POURQUOI 63 ? 8x8 = 64 coefficients total, moins 1 (le DC) = 63
	// RÉSULTAT : Valeur de seuil pour binariser les coefficients DCT
	return total / 63
}
