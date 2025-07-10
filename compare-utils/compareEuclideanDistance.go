package compare_utils

import "math"

/*
===== DISTANCE EUCLIDIENNE ENTRE DEUX COULEURS =====

À QUOI ÇA SERT :
Calcule la "distance géométrique" entre deux couleurs dans l'espace RGB 3D.
Mesure à quel point deux couleurs sont différentes de manière intuitive.

PRINCIPE GÉOMÉTRIQUE :
- Chaque couleur = point dans un cube 3D avec coordonnées (R, G, B)
- Distance = longueur de la ligne droite entre les deux points
- Formule classique : √[(R1-R2)² + (G1-G2)² + (B1-B2)²]

UTILISATION DANS LE SYSTÈME :
- Comparaison des couleurs moyennes globales entre images
- Comparaison des couleurs moyennes entre tuiles correspondantes
- Contribue au score de similarité final (pondéré à 15%)

Paramètres :
- c1, c2 : couleurs sous forme [R,G,B] avec valeurs 0-255

Retour :
- Distance euclidienne (0 = identiques, ~441 = opposées complètement)
*/
func EuclideanDistance(c1, c2 [3]float64) float64 {

	// Application directe de la formule de distance euclidienne 3D
	// DÉCOMPOSITION :
	// - (c1[0]-c2[0])² : Carré de la différence des rouges
	// - (c1[1]-c2[1])² : Carré de la différence des verts
	// - (c1[2]-c2[2])² : Carré de la différence des bleus
	// - √(somme) : Racine carrée pour obtenir la distance réelle

	return math.Sqrt(
		math.Pow(c1[0]-c2[0], 2) + // Contribution du rouge à la distance
			math.Pow(c1[1]-c2[1], 2) + // Contribution du vert à la distance
			math.Pow(c1[2]-c2[2], 2), // Contribution du bleu à la distance
	)
}
