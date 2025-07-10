package analyzer

import (
	"fmt"
	drawx "golang.org/x/image/draw"
	"image"
	"image/draw"
	"math"
)

/*
	===== CALCUL HISTOGRAMME RGB =====

Calcule l'histogramme des composantes R, G et B pour une image.
Chaque couleur est répartie sur `Bins` intervalles.

Retour :
- Un dictionnaire contenant trois histogrammes (r, g, b)
*/
func computeHistogramRGB(img image.Image) map[string][]int {
	rHist, gHist, bHist := make([]int, Bins), make([]int, Bins), make([]int, Bins)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rHist[int(r>>8)*Bins/256]++
			gHist[int(g>>8)*Bins/256]++
			bHist[int(b>>8)*Bins/256]++
		}
	}
	return map[string][]int{"r": rHist, "g": gHist, "b": bHist}
}

/*
	===== CALCUL HISTOGRAMME HSV =====

Calcule l'histogramme pour la teinte (H), saturation (S) et valeur (V).
Les valeurs sont converties depuis RGB.

Retour :
- Un dictionnaire contenant trois histogrammes (h, s, v)
*/
func computeHistogramHSV(img image.Image) map[string][]int {
	hHist, sHist, vHist := make([]int, Bins), make([]int, Bins), make([]int, Bins)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			h, s, v := rgbToHsv(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			hIdx := int(h * float64(Bins) / 360)
			sIdx := int(s * float64(Bins))
			vIdx := int(v * float64(Bins))

			// Sécuriser les index
			if hIdx >= Bins {
				hIdx = Bins - 1
			}
			if sIdx >= Bins {
				sIdx = Bins - 1
			}
			if vIdx >= Bins {
				vIdx = Bins - 1
			}

			hHist[hIdx]++
			sHist[sIdx]++
			vHist[vIdx]++
		}
	}
	return map[string][]int{"h": hHist, "s": sHist, "v": vHist}
}

/*
	===== CALCUL MOYENNE DES COULEURS =====

Calcule la moyenne des composantes rouge, vert et bleu sur l'image.

Retour :
- Vecteur [R, G, B] en float64
*/
func computeMeanColor(img image.Image) [3]float64 {
	var rSum, gSum, bSum float64
	bounds := img.Bounds()
	total := float64(bounds.Dx() * bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rSum += float64(r >> 8)
			gSum += float64(g >> 8)
			bSum += float64(b >> 8)
		}
	}
	return [3]float64{rSum / total, gSum / total, bSum / total}
}

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

/*
	===== CONVERSION RGB -> HSV =====

Convertit une couleur RGB en HSV :
- H (teinte) ∈ [0, 360)
- S (saturation) ∈ [0, 1]
- V (valeur) ∈ [0, 1]

Retour :
- h, s, v en float64
*/
func rgbToHsv(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	delta := max - min

	var h float64
	if delta != 0 {
		switch max {
		case rf:
			h = math.Mod(((gf - bf) / delta), 6)
		case gf:
			h = ((bf - rf) / delta) + 2
		case bf:
			h = ((rf - gf) / delta) + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	s := 0.0
	if max != 0 {
		s = delta / max
	}
	v := max
	return h, s, v
}

/*
	===== DIFFERENCE ABSOLUE ENTRE 2 VALEURS UINT8 =====

Retour :
- Valeur absolue de la différence
*/
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

/*
	===== GÉNÉRATION DU pHash (Perceptual Hash) =====

Produit un hash compact (64 bits) représentant l'aspect visuel global de l'image.

Méthode :
1. Redimensionnement en 32x32
2. Conversion en niveaux de gris
3. Application d'une DCT (Discrete Cosine Transform)
4. Comparaison aux coefficients moyens

Retour :
- Hash sous forme d'une chaîne hexadécimale (16 caractères)
*/
func generatePHash(img image.Image) string {
	size := 32
	gray := image.NewGray(image.Rect(0, 0, size, size))
	drawx.ApproxBiLinear.Scale(gray, gray.Bounds(), img, img.Bounds(), draw.Over, nil)

	dctVals := dct2D(gray)
	avg := averageDCT(dctVals)

	var hash uint64
	index := 0

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if dctVals[y][x] > avg {
				hash |= 1 << index
			}
			index++
		}
	}
	return fmt.Sprintf("%016x", hash)
}

/*
	===== TRANSFORMÉE EN COSINUS 2D (DCT) =====

Applique une DCT bidimensionnelle à une image en niveaux de gris.
Utilisée dans la génération du pHash.

Retour :
- Matrice des coefficients DCT
*/
func dct2D(img *image.Gray) [][]float64 {
	N := 32
	dct := make([][]float64, N)
	for i := range dct {
		dct[i] = make([]float64, N)
	}

	for u := 0; u < N; u++ {
		for v := 0; v < N; v++ {
			var sum float64
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					pixel := float64(img.GrayAt(x, y).Y)
					sum += pixel *
						math.Cos((2*float64(x)+1)*float64(u)*math.Pi/float64(2*N)) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/float64(2*N))
				}
			}
			cu, cv := 1.0, 1.0
			if u == 0 {
				cu = 1 / math.Sqrt(2)
			}
			if v == 0 {
				cv = 1 / math.Sqrt(2)
			}
			dct[u][v] = 0.25 * cu * cv * sum
		}
	}
	return dct
}

/*
	===== MOYENNE DES COEFFICIENTS DCT =====

Calcule la moyenne des 64 premiers coefficients DCT,
en ignorant le premier coefficient (DC term).

Retour :
- Moyenne des coefficients DCT
*/
func averageDCT(dct [][]float64) float64 {
	total := 0.0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if !(x == 0 && y == 0) {
				total += dct[y][x]
			}
		}
	}
	return total / 63
}
