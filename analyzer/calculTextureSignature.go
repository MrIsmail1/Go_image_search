package analyzer

import (
	"image"
	"image/draw"
)

/*
	===== CALCUL SIGNATURE DE TEXTURE =====

Mesure la variation locale des intensités (texture) dans une image
convertie en niveaux de gris.

Retour :
- Moyenne des variations de texture
*/
func computeTextureSignature(img image.Image) float64 {
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.Point{}, draw.Src)

	var variations float64
	bounds := gray.Bounds()

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			c := gray.GrayAt(x, y).Y
			right := gray.GrayAt(x+1, y).Y
			bottom := gray.GrayAt(x, y+1).Y
			variations += float64(absDiff(c, right) + absDiff(c, bottom))
		}
	}
	return variations / float64((bounds.Dx()-2)*(bounds.Dy()-2))
}
