package color

import (
	"github.com/MrIsmail1/Golang_images_matcher/config"
	"image"
)

/*
===== CALCUL HISTOGRAMME HSV =====

À QUOI ÇA SERT :
Calcule l'histogramme des couleurs dans l'espace HSV plutôt que RGB.
HSV sépare mieux les aspects "couleur pure" et "luminosité" d'une image.

POURQUOI HSV C'EST MIEUX QUE RGB :
- HSV = Hue (teinte), Saturation (pureté), Value (luminosité)
- Plus proche de comment l'œil humain perçoit les couleurs
- Robuste aux variations d'éclairage (ombre vs plein soleil)
- Permet de séparer la couleur de sa brillance

Paramètre :
- img : image à analyser

Retour :
- Dictionnaire avec 3 histogrammes : {"h": [...], "s": [...], "v": [...]}
- Chaque histogramme a config.Bins valeurs (64 par défaut)
*/
func ComputeHistogramHSV(img image.Image) map[string][]int {

	// Création des 3 histogrammes vides avec la taille configurée
	// POURQUOI 3 histogrammes séparés ?
	// - H : distribution des teintes (rouge, vert, bleu, etc.)
	// - S : distribution des saturations (couleur pure vs délavée)
	// - V : distribution des luminosités (sombre vs clair)
	hHist, sHist, vHist := make([]int, config.Bins), make([]int, config.Bins), make([]int, config.Bins)
	bounds := img.Bounds()

	// Parcours de tous les pixels comme d'habitude
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			// Lecture du pixel et conversion RGB → HSV
			r, g, b, _ := img.At(x, y).RGBA()
			h, s, v := rgbToHsv(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			// Calcul des index pour chaque composante HSV
			// FORMULES DE MAPPING :
			// - Teinte : 0-360° → 0 à (Bins-1)
			// - Saturation : 0-1 → 0 à (Bins-1)
			// - Valeur : 0-1 → 0 à (Bins-1)
			hIdx := int(h * float64(config.Bins) / 360) // h va de 0 à 360°
			sIdx := int(s * float64(config.Bins))       // s va de 0 à 1
			vIdx := int(v * float64(config.Bins))       // v va de 0 à 1

			if hIdx >= config.Bins {
				hIdx = config.Bins - 1
			}
			if sIdx >= config.Bins {
				sIdx = config.Bins - 1
			}
			if vIdx >= config.Bins {
				vIdx = config.Bins - 1
			}

			// Incrémentation des compteurs pour ce pixel
			// Chaque pixel contribue +1 au bin correspondant de chaque histogramme
			hHist[hIdx]++ // Compteur teinte
			sHist[sIdx]++ // Compteur saturation
			vHist[vIdx]++ // Compteur valeur/luminosité
		}
	}

	// Retour des 3 histogrammes dans un dictionnaire
	// STRUCTURE : {"h": [64 valeurs], "s": [64 valeurs], "v": [64 valeurs]}
	// UTILISATION : Permet d'accéder facilement à chaque histogramme par clé
	return map[string][]int{"h": hHist, "s": sHist, "v": vHist}
}
