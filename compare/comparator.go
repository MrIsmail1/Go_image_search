package compare

import (
	"github.com/MrIsmail1/Golang_images_matcher/compare-utils"
	"github.com/MrIsmail1/Golang_images_matcher/config"
	"github.com/MrIsmail1/Golang_images_matcher/model"
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
	rgbDist := compare_utils.CompareHistograms(desc1.GlobalRGB, desc2.GlobalRGB)               // Distance des histogrammes RGB
	hsvDist := compare_utils.CompareHistograms(desc1.GlobalHSV, desc2.GlobalHSV)               // Distance des histogrammes HSV
	colorDist := compare_utils.EuclideanDistance(desc1.GlobalMeanColor, desc2.GlobalMeanColor) // Distance euclidienne des moyennes de couleurs
	textureDist := math.Abs(desc1.GlobalTexture - desc2.GlobalTexture)                         // Différence absolue de texture
	phashDist := compare_utils.HammingDistance(desc1.GlobalPHash, desc2.GlobalPHash)           // Distance de Hamming entre les pHash

	// --- Normalisation des distances ---
	normRGB := rgbDist / float64(config.Bins*3*255)
	normHSV := hsvDist / float64(config.Bins*3*255)
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
		rgbDist := compare_utils.CompareHistograms(t1.HistogramRGB, t2.HistogramRGB)
		hsvDist := compare_utils.CompareHistograms(t1.HistogramHSV, t2.HistogramHSV)
		colorDist := compare_utils.EuclideanDistance(t1.MeanColor, t2.MeanColor)
		textureDist := math.Abs(t1.TextureSignature - t2.TextureSignature)
		phashDist := compare_utils.HammingDistance(t1.PHash, t2.PHash)
		shapeDist := math.Abs(t1.ShapeSignature - t2.ShapeSignature)

		// Normalisation
		normRGB := rgbDist / float64(config.Bins*3*255)
		normHSV := hsvDist / float64(config.Bins*3*255)
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
