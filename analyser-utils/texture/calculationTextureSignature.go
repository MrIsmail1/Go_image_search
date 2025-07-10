package texture

import (
	"github.com/MrIsmail1/Golang_images_matcher/analyser-utils/math"
	"image"
	"image/draw"
)

/*
===== CALCUL SIGNATURE DE TEXTURE =====

À QUOI ÇA SERT :
Mesure la "rugosité" visuelle d'une image en calculant les variations d'intensité
entre pixels voisins. Permet de distinguer surfaces lisses des surfaces rugueuses.

CONCEPT DE TEXTURE :
- Texture lisse : variations faibles entre voisins (ciel, eau calme, mur peint)
- Texture rugueuse : variations importantes entre voisins (herbe, fourrure, bois)
- Signature = moyenne de toutes ces variations locales

Paramètre :
- img : image à analyser

Retour :
- Valeur représentant l'intensité moyenne de texture
- Plus élevé = plus rugueux, plus faible = plus lisse
*/
func ComputeTextureSignature(img image.Image) float64 {

	// Conversion nécessaire car la texture concerne la structure, pas la couleur
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.Point{}, draw.Src)

	var variations float64
	bounds := gray.Bounds()

	// Parcours en évitant les bords (on a besoin des voisins droite et bas)
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ { // Évite première et dernière ligne
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ { // Évite première et dernière colonne

			// Lecture du pixel central et de ses voisins directs
			c := gray.GrayAt(x, y).Y        // Pixel central
			right := gray.GrayAt(x+1, y).Y  // Voisin de droite (variation horizontale)
			bottom := gray.GrayAt(x, y+1).Y // Voisin du dessous (variation verticale)

			// Calcul et accumulation des variations dans les 2 directions
			// AbsDiff() garantit une valeur positive (on s'intéresse à l'ampleur, pas au sens)
			variations += float64(math.AbsDiff(c, right) + math.AbsDiff(c, bottom))
		}
	}

	// Calcul de la moyenne pour normaliser selon la taille de l'image
	// RÉSULTAT : Signature indépendante de la taille, comparable entre images
	return variations / float64((bounds.Dx()-2)*(bounds.Dy()-2))
}
