package shape

import (
	"image"
	"image/draw"
	"math"
)

/*
===== DÉTECTION DE CONTOURS (SIGNATURE DE FORME) =====

À QUOI ÇA SERT :
Détecte les contours (bordures des objets) et calcule une signature basée sur leur densité.
Révèle la richesse en formes et structures géométriques de l'image.

PRINCIPE DES CONTOURS :
- Contour = transition brutale d'intensité lumineuse
- Bordure entre objet sombre et fond clair
- Arêtes de bâtiments, silhouettes, changements de texture

QU'EST-CE QUE LE FILTRE DE SOBEL :
- Opérateur de détection de contours très populaire
- Calcule les gradients (variations) horizontaux et verticaux
- Combine les deux pour obtenir l'intensité totale du contour
- Plus robuste au bruit que des méthodes simples

Paramètre :
- img : image à analyser

Retour :
- Ratio 0-1 : densité de contours (0=uniforme, 1=très détaillé)
*/
func ComputeShapeSignature(img image.Image) float64 {

	// Conversion nécessaire car les contours dépendent de l'intensité, pas de la couleur
	// AVANTAGE : Plus simple (1 valeur/pixel), plus rapide, plus robuste
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.Point{}, draw.Src)

	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	var contourPixels int

	// Masque de Sobel X (détecte contours verticaux)
	// PRINCIPE : Différence pondérée entre colonnes gauche et droite
	// EFFET : Réagit aux transitions horizontales (objets verticaux)
	kernelX := [3][3]int{
		{-1, 0, 1}, // Ligne haut : -gauche + droite
		{-2, 0, 2}, // Ligne milieu : -2×gauche + 2×droite (poids double)
		{-1, 0, 1}, // Ligne bas : -gauche + droite
	}

	// Masque de Sobel Y (détecte contours horizontaux)
	// PRINCIPE : Différence pondérée entre lignes haut et bas
	// EFFET : Réagit aux transitions verticales (objets horizontaux)
	kernelY := [3][3]int{
		{-1, -2, -1}, // Soustrait ligne du haut (poids double au centre)
		{0, 0, 0},    // Ligne centrale neutre
		{1, 2, 1},    // Ajoute ligne du bas (poids double au centre)
	}

	// Parcours en évitant les bords (masque 3×3 a besoin de voisins complets)
	for y := 1; y < height-1; y++ { // Évite première et dernière ligne
		for x := 1; x < width-1; x++ { // Évite première et dernière colonne

			var gx, gy int // Gradients horizontal et vertical

			// Convolution 3×3 avec les masques de Sobel
			for ky := -1; ky <= 1; ky++ { // Voisinage vertical (-1, 0, +1)
				for kx := -1; kx <= 1; kx++ { // Voisinage horizontal (-1, 0, +1)

					// Lecture de l'intensité du pixel voisin
					p := gray.GrayAt(x+kx, y+ky).Y

					// Accumulation pondérée pour chaque direction
					// +1 car les indices kernel vont de 0 à 2, mais k va de -1 à +1
					gx += int(p) * kernelX[ky+1][kx+1] // Contribution au gradient horizontal
					gy += int(p) * kernelY[ky+1][kx+1] // Contribution au gradient vertical
				}
			}

			// Magnitude du gradient : combine horizontal et vertical
			// FORMULE : |G| = √(Gx² + Gy²)
			// PRINCIPE : Intensité du contour indépendamment de sa direction
			// EXEMPLE : Contour diagonal aura Gx et Gy non nuls
			gradient := math.Sqrt(float64(gx*gx + gy*gy))

			// Test de seuillage pour détecter un contour significatif
			// SEUIL 100 : Valeur empirique (sur échelle 0-255×√8 ≈ 0-720)
			// EFFET : Ignore le bruit, garde les vrais contours
			if gradient > 100 {
				contourPixels++ // Ce pixel compte comme un contour
			}
		}
	}

	// Calcul du nombre total de pixels analysés (sans les bordures)
	totalPixels := (width - 2) * (height - 2)

	// Ratio de densité de contours
	return float64(contourPixels) / float64(totalPixels)
}
