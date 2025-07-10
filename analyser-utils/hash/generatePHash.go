package hash

import (
	"fmt"
	"github.com/MrIsmail1/Golang_images_matcher/analyser-utils/math"
	drawx "golang.org/x/image/draw"
	"image"
	"image/draw"
)

/*
===== GÉNÉRATION DU pHash (Perceptual Hash) =====

À QUOI ÇA SERT :
Crée une "empreinte digitale" de 64 bits qui capture l'essence visuelle d'une image.
Contrairement aux hash classiques (MD5/SHA), deux images similaires auront des pHash similaires.

RÉVOLUTION DU pHash :
- Hash classique : 1 pixel différent = hash totalement différent
- pHash : Images similaires = hash similaires (même avec petites modifs)
- Résistant à : compression, redimensionnement, légères rotations, changements de luminosité

POURQUOI CETTE MÉTHODE FONCTIONNE :
1. 32x32 → Élimine les détails fins, garde la structure
2. Niveaux de gris → Ignore les variations de couleur mineures
3. DCT → Extrait les patterns visuels dominants (comme JPEG)
4. Binarisation → Signature compacte de 64 bits

Paramètre :
- img : image à traiter (toute taille, tout format)

Retour :
- Hash hexadécimal de 16 caractères (64 bits)
*/
func GeneratePHash(img image.Image) string {

	// ============================================================================================
	// ÉTAPE 1 : STANDARDISATION (32x32 niveaux de gris)
	// ============================================================================================

	// Taille standard pour le pHash
	size := 32 //  Standard établi pour les algorithmes pHash

	// Création d'une image en niveaux de gris 32x32
	gray := image.NewGray(image.Rect(0, 0, size, size))
	drawx.ApproxBiLinear.Scale(gray, gray.Bounds(), img, img.Bounds(), draw.Over, nil)
	dctVals := math.Dct2D(gray)

	avg := averageDCT(dctVals)

	// Variable pour construire le hash binaire 64 bits
	var hash uint64
	index := 0 // Position du bit (0 à 63)

	// Parcours des 64 premiers coefficients DCT (8x8 du coin supérieur gauche)
	// CES COEFFICIENTS contiennent l'info structurelle principale de l'image
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {

			// Test de binarisation : coefficient > seuil ?
			if dctVals[y][x] > avg {
				// Si oui → place un bit 1 à la position 'index'
				// OPÉRATION : hash |= (1 << index) met le bit à 1
				hash |= 1 << index
			}
			// Si non → le bit reste à 0 (valeur par défaut)

			index++ // Passe au bit suivant
		}
	}

	// FORMAT : %016x = hex sur exactement 16 caractères avec zéros initiaux
	return fmt.Sprintf("%016x", hash)
}
