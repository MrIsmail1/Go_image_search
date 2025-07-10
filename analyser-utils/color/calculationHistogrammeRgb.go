package color

import (
	"github.com/MrIsmail1/Golang_images_matcher/config"
	"image"
)

/*
===== CALCUL HISTOGRAMME RGB =====

À QUOI ÇA SERT :
Calcule la distribution statistique des couleurs Rouge, Vert et Bleu dans l'image.
C'est comme trier tous les pixels par couleur et compter combien il y en a de chaque type.

PRINCIPE D'UN HISTOGRAMME :
- Divise l'espace des couleurs en "seaux" (bins)
- Compte combien de pixels tombent dans chaque seau
- Donne une signature compacte de la répartition des couleurs

Paramètre :
- img : image à analyser

Retour :
- Dictionnaire {"r": [...], "g": [...], "b": [...]} avec 3 histogrammes
- Chaque histogramme a config.Bins valeurs (64 par défaut)
*/
func ComputeHistogramRGB(img image.Image) map[string][]int {

	// Création des 3 histogrammes vides pour Rouge, Vert, Bleu
	// POURQUOI 3 histogrammes séparés ?
	// - Chaque couleur a sa propre distribution
	// - Permet d'analyser les teintes dominantes séparément
	// - Plus informatif qu'un seul histogramme global
	rHist, gHist, bHist := make([]int, config.Bins), make([]int, config.Bins), make([]int, config.Bins)
	bounds := img.Bounds()

	// Parcours classique de tous les pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			// Lecture du pixel et classification directe dans les bins
			r, g, b, _ := img.At(x, y).RGBA()

			// FORMULE DE BINNING : (valeur_8_bits * Bins) / 256
			// EXPLICATION :
			// - Conversion 16→8 bits : r>>8 (0 à 255)
			// - Mapping vers bins : valeur * 64 / 256 (0 à 63)
			// - Exemple : rouge=200 → bin 200*64/256 = bin 50
			rHist[int(r>>8)*config.Bins/256]++ // Classification + comptage rouge
			gHist[int(g>>8)*config.Bins/256]++ // Classification + comptage vert
			bHist[int(b>>8)*config.Bins/256]++ // Classification + comptage bleu
		}
	}

	// Retour des 3 histogrammes dans un dictionnaire
	// STRUCTURE : {"r": [64 valeurs], "g": [64 valeurs], "b": [64 valeurs]}
	// UTILISATION : Facile d'accéder à chaque couleur individuellement
	return map[string][]int{"r": rHist, "g": gHist, "b": bHist}
}
