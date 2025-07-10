package color

import "image"

/*
===== CALCUL DE LA COULEUR MOYENNE D'UNE IMAGE =====

À QUOI ÇA SERT :
Calcule la couleur "représentative" d'une image en faisant la moyenne arithmétique
de tous les pixels. C'est comme si on prenait tous les pixels de l'image,
qu'on les mélangeait ensemble comme de la peinture, et qu'on regardait
la couleur finale obtenue.

Paramètre :
  - img : Image Go standard (interface image.Image)
    Peut être n'importe quel format supporté (JPEG, PNG, GIF, etc.)

Retour :
- Array de 3 float64 : [Rouge_moyen, Vert_moyen, Bleu_moyen]
- Valeurs entre 0.0 et 255.0
- Précision float64 pour éviter les erreurs d'arrondi
*/
func ComputeMeanColor(img image.Image) [3]float64 {
	var rSum, gSum, bSum float64

	bounds := img.Bounds() // Bounds méthode qui récupere les images (implémenter de base sur Go)

	// Dx() = Delta X = largeur (bounds.Max.X - bounds.Min.X)
	// Dy() = Delta Y = hauteur (bounds.Max.Y - bounds.Min.Y)
	total := float64(bounds.Dx() * bounds.Dy())

	// OPTIMISATION : Parcours séquentiel optimal pour la cache mémoire du CPU
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ { // Pour chaque ligne (rangée de pixels)
		for x := bounds.Min.X; x < bounds.Max.X; x++ { // Pour chaque colonne (pixel dans la rangée)

			// Récupération des valeurs RGBA du pixel aux coordonnées (x, y)
			// - img.At(x,y) : Méthode standard Go pour lire un pixel
			// - .RGBA() : Extrait Rouge, Vert, Bleu, Alpha (transparence)
			// - Valeurs 16 bits : Go utilise uint16 (0-65535) pour plus de précision
			r, g, b, _ := img.At(x, y).RGBA()

			rSum += float64(r >> 8) // Conversion 16→8 bits + accumulation Rouge
			gSum += float64(g >> 8) // Conversion 16→8 bits + accumulation Vert
			bSum += float64(b >> 8) // Conversion 16→8 bits + accumulation Bleu
		}
	}

	// Calcul des moyennes finales en divisant les sommes par le nombre total de pixels
	// RÉSULTAT : Couleur moyenne de l'image sous forme [Rouge, Vert, Bleu]
	return [3]float64{rSum / total, gSum / total, bSum / total}
}
