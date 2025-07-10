package compare_utils

import "math"

/*
===== COMPARAISON D'HISTOGRAMMES =====

À QUOI ÇA SERT :
Compare deux histogrammes de couleur en calculant la "distance" entre leurs distributions.
Plus la distance est faible, plus les images ont des répartitions de couleurs similaires.

PRINCIPE DE COMPARAISON :
- Histogramme = répartition statistique des couleurs (64 bins par couleur)
- Comparaison = somme des différences bin par bin
- Distance faible = distributions similaires = images similaires

UTILISATION DANS LE SYSTÈME :
- Comparaison histogrammes RGB globaux (pondération 10%)
- Comparaison histogrammes HSV globaux (pondération 10%)
- Comparaison tuile par tuile pour analyse locale

Paramètres :
- h1, h2 : histogrammes sous forme map[string][]int

Retour :
- Distance L1 (somme des différences absolues)
*/
func CompareHistograms(h1, h2 map[string][]int) float64 {
	total := 0.0

	// Parcours de toutes les composantes de l'histogramme
	for key, arr1 := range h1 {

		// Récupération de l'array correspondant dans le second histogramme
		arr2 := h2[key]

		// Parcours de tous les bins de cette composante couleur
		// CHAQUE BIN = compteur de pixels dans cet intervalle de couleur
		for i := 0; i < len(arr1); i++ {

			// Calcul de la différence absolue pour ce bin
			// PRINCIPE : |compteur1 - compteur2|
			// EXEMPLE :
			// - Image1 a 100 pixels rouges dans le bin 30
			// - Image2 a 80 pixels rouges dans le bin 30
			// - Différence = |100 - 80| = 20
			total += math.Abs(float64(arr1[i] - arr2[i]))
		}
	}

	// Retour de la somme totale des différences
	return total
}
