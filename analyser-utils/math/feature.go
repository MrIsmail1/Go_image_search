package math

import (
	"image"
	"math"
)

/*
===== DIFFÉRENCE ABSOLUE ENTRE 2 VALEURS UINT8 =====

À QUOI ÇA SERT :
Calcule |a - b| sans risque de débordement pour des valeurs uint8.
Fonction utilitaire simple mais essentielle pour les calculs de variations.

POURQUOI CETTE FONCTION :
- Go n'a pas de fonction abs() pour les entiers dans la lib standard
- Évite les valeurs négatives dans les calculs de distance/variation
- Plus rapide qu'une conversion vers float64 puis math.Abs()
- Spécialisée pour uint8 (valeurs de pixels 0-255)

USAGE TYPIQUE :
- Calcul de variations de texture entre pixels voisins
- Mesure de différences d'intensité en niveaux de gris
- Détection de contours simples

Paramètres :
- a, b : deux valeurs uint8 à comparer (typiquement intensités de pixels)

Retour :
- Différence absolue |a - b| toujours positive
*/
func AbsDiff(a, b uint8) uint8 {
	// Test simple pour garantir un résultat positif
	// PRINCIPE : |a - b| = (a - b) si a > b, sinon (b - a)
	if a > b {
		return a - b // a plus grand → différence positive
	}
	return b - a // b plus grand → on inverse pour rester positif
}

/*
===== TRANSFORMÉE EN COSINUS DISCRÈTE 2D (DCT) =====

À QUOI ÇA SERT :
Transforme une image du domaine spatial vers le domaine fréquentiel.
C'est la même technologie que JPEG ! Sépare structure générale des détails fins.

PRINCIPE DE LA DCT :
- Décompose l'image en somme de motifs cosinus de différentes fréquences
- Basses fréquences (coin sup-gauche) = formes générales, contours principaux
- Hautes fréquences (coin inf-droit) = détails fins, textures, bruit
- Concentre l'énergie dans les premiers coefficients

Paramètre :
- img : image en niveaux de gris 32×32 pixels

Retour :
- Matrice 32×32 des coefficients DCT
*/
func Dct2D(img *image.Gray) [][]float64 {
	N := 32 // Taille de la matrice (32×32 standard pour pHash)

	// Initialisation de la matrice de coefficients DCT
	dct := make([][]float64, N)
	for i := range dct {
		dct[i] = make([]float64, N)
	}

	// CALCUL DE LA DCT 2D COMPLÈTE
	// Double boucle sur les fréquences de sortie (u,v)
	// (u,v) = coordonnées dans l'espace fréquentiel
	for u := 0; u < N; u++ { // Pour chaque fréquence horizontale
		for v := 0; v < N; v++ { // Pour chaque fréquence verticale

			var sum float64

			// Double boucle sur les pixels d'entrée (x,y)
			// (x,y) = coordonnées dans l'espace spatial (image originale)
			for x := 0; x < N; x++ { // Pour chaque colonne de l'image
				for y := 0; y < N; y++ { // Pour chaque ligne de l'image

					// Lecture de l'intensité du pixel
					pixel := float64(img.GrayAt(x, y).Y)

					// Application de la formule DCT 2D
					// FORMULE : DCT(u,v) = Σ Σ f(x,y) × cos(...) × cos(...)
					// Les cosinus créent les "bases de fréquence" 2D
					sum += pixel *
						math.Cos((2*float64(x)+1)*float64(u)*math.Pi/float64(2*N)) * // Cosinus horizontal
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/float64(2*N)) // Cosinus vertical
				}
			}

			cu, cv := 1.0, 1.0
			if u == 0 {
				cu = 1 / math.Sqrt(2) // Facteur spécial pour fréquence nulle horizontale
			}
			if v == 0 {
				cv = 1 / math.Sqrt(2) // Facteur spécial pour fréquence nulle verticale
			}

			// Coefficient DCT final avec normalisation
			// FACTEUR 0.25 = 1/(2×2) pour normaliser la DCT 2D
			dct[u][v] = 0.25 * cu * cv * sum
		}
	}

	return dct
}
