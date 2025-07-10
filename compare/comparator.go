package compare

import (
	"fmt"
	"math"
)

/*
	===== COMPARAISON DE DEUX DESCRIPTEURS D'IMAGE =====

Cette fonction évalue la similarité entre deux images, en comparant :
- leurs caractéristiques globales (histogrammes, couleurs moyennes, textures, hash perceptuel),
- ainsi que les caractéristiques détaillées de leurs tuiles.

Retour :
- Score de similarité final entre 0 et 100 (%)
*/
func CompareDescriptors(desc1, desc2 *model.FullImageDescriptor) float64 {
	// --- Comparaison globale ---
	rgbDist := compareHistograms(desc1.GlobalRGB, desc2.GlobalRGB)               // Distance des histogrammes RGB
	hsvDist := compareHistograms(desc1.GlobalHSV, desc2.GlobalHSV)               // Distance des histogrammes HSV
	colorDist := euclideanDistance(desc1.GlobalMeanColor, desc2.GlobalMeanColor) // Distance euclidienne des moyennes de couleurs
	textureDist := math.Abs(desc1.GlobalTexture - desc2.GlobalTexture)           // Différence absolue de texture
	phashDist := hammingDistance(desc1.GlobalPHash, desc2.GlobalPHash)           // Distance de Hamming entre les pHash

	// --- Normalisation des distances ---
	normRGB := rgbDist / float64(analyzer.Bins*3*255)
	normHSV := hsvDist / float64(analyzer.Bins*3*255)
	normColor := colorDist / (255 * math.Sqrt(3))
	normTexture := textureDist / 500.0
	normPHash := float64(phashDist) / 64.0

	shapeDist := math.Abs(desc1.GlobalShape - desc2.GlobalShape)
	normShape := shapeDist / 1.0 // car ratio ∈ [0,1]

	globalScore := 1 -
		normRGB*0.1 -
		normHSV*0.1 -
		normColor*0.15 -
		normTexture*0.15 -
		normShape*0.25 -
		(1-normPHash)*0.25

	// --- Comparaison tuile par tuile ---
	var tileScoreSum float64
	for i := 0; i < len(desc1.Tiles); i++ {
		t1 := desc1.Tiles[i]
		t2 := desc2.Tiles[i]

		// Comparaison locale
		rgbDist := compareHistograms(t1.HistogramRGB, t2.HistogramRGB)
		hsvDist := compareHistograms(t1.HistogramHSV, t2.HistogramHSV)
		colorDist := euclideanDistance(t1.MeanColor, t2.MeanColor)
		textureDist := math.Abs(t1.TextureSignature - t2.TextureSignature)
		phashDist := hammingDistance(t1.PHash, t2.PHash)
		shapeDist := math.Abs(t1.ShapeSignature - t2.ShapeSignature)

		// Normalisation
		normRGB := rgbDist / float64(analyzer.Bins*3*255)
		normHSV := hsvDist / float64(analyzer.Bins*3*255)
		normColor := colorDist / (255 * math.Sqrt(3))
		normTexture := textureDist / 1000.0
		normPHash := float64(phashDist) / 64.0
		normShape := shapeDist / 1.0

		tileScore := 1 - normRGB*0.1 - normHSV*0.1 - normColor*0.15 - normTexture*0.15 - normShape*0.25 - (1-normPHash)*0.25

		// Correction : tuile très similaire = score parfait
		if tileScore < 0.85 {
			tileScoreSum += tileScore
		} else {
			tileScoreSum += 1.0
		}
	}

	// Moyenne des scores des tuiles
	avgTileScore := tileScoreSum / float64(len(desc1.Tiles))

	// --- Score final ---
	finalScore := (globalScore*0.65 + avgTileScore*0.35) * 100
	if finalScore < 0 {
		finalScore = 0
	}
	return finalScore
}

/*
	===== COMPARAISON D'HISTOGRAMMES =====

Calcule la distance entre deux histogrammes (RGB ou HSV),
en sommant les différences absolues entre chaque bin.

Retour :
- Somme des différences absolues
*/
func compareHistograms(h1, h2 map[string][]int) float64 {
	total := 0.0
	for key, arr1 := range h1 {
		arr2 := h2[key]
		for i := 0; i < len(arr1); i++ {
			total += math.Abs(float64(arr1[i] - arr2[i]))
		}
	}
	return total
}

/*
	===== DISTANCE EUCLIDIENNE ENTRE DEUX COULEURS =====

Calcule la distance euclidienne dans l'espace RGB,
entre deux vecteurs de couleur (chacun ayant 3 composantes).

Retour :
- Distance euclidienne
*/
func euclideanDistance(c1, c2 [3]float64) float64 {
	return math.Sqrt(
		math.Pow(c1[0]-c2[0], 2) +
			math.Pow(c1[1]-c2[1], 2) +
			math.Pow(c1[2]-c2[2], 2),
	)
}

/*
	===== DISTANCE DE HAMMING ENTRE DEUX pHash =====

Calcule le nombre de bits différents entre deux hash perceptuels (pHash),
en utilisant un XOR binaire suivi d'un comptage des bits à 1.

Retour :
- Nombre de bits différents (entier entre 0 et 64)
*/
func hammingDistance(hash1, hash2 string) int {
	var v1, v2 uint64
	fmt.Sscanf(hash1, "%x", &v1) // Conversion du hash hexadécimal en entier
	fmt.Sscanf(hash2, "%x", &v2)

	xor := v1 ^ v2 // XOR binaire
	dist := 0
	for xor != 0 {
		dist++
		xor &= xor - 1 // Technique rapide pour compter les bits à 1
	}
	return dist
}