package analyzer

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/MrIsmail1/Golang_images_matcher/model"
	drawx "golang.org/x/image/draw"
)

const (
	StandardSize = 256
	TilesPerRow  = 9
	Bins         = 64
)

/*
	===== ANALYSE D'UNE IMAGE =====

Cette fonction lit une image depuis le disque, la redimensionne à une taille standard,
puis en extrait un ensemble de descripteurs globaux et locaux (par tuiles).

Elle retourne une structure complète `FullImageDescriptor` contenant :
- des histogrammes de couleur (RGB et HSV),
- un hash perceptuel (pHash),
- la moyenne des couleurs,
- une signature de texture,
- et des informations similaires pour chaque tuile.

Paramètre :
- imagePath : chemin complet vers le fichier image à analyser

Retour :
- Pointeur vers un `FullImageDescriptor` rempli
- Une erreur si le chargement ou l'analyse échoue
*/
func AnalyzeImage(imagePath string) (*model.FullImageDescriptor, error) {
	// === Chargement de l'image ===
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err // Échec si le fichier ne peut pas être ouvert
	}
	defer file.Close()

	// Décodage de l'image (JPEG, PNG, GIF, etc.)
	srcImg, _, err := image.Decode(file)
	if err != nil {
		return nil, err // Échec si l'image n'est pas décodable
	}

	// === Redimensionnement de l'image à une taille standard (256x256) ===
	resized := image.NewRGBA(image.Rect(0, 0, StandardSize, StandardSize))
	drawx.ApproxBiLinear.Scale(resized, resized.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	// === Extraction des descripteurs globaux ===
	globalRGB := computeHistogramRGB(resized)         // Histogramme RGB de l'image entière
	globalHSV := computeHistogramHSV(resized)         // Histogramme HSV de l'image entière
	globalPHash := generatePHash(resized)             // Hash perceptuel global
	globalMean := computeMeanColor(resized)           // Moyenne des couleurs globales
	globalTexture := computeTextureSignature(resized) // Texture globale de l'image
	globalShape := computeShapeSignature(resized)

	// === Découpage en tuiles et analyse de chaque tuile ===
	tileSize := StandardSize / TilesPerRow // Taille d'une tuile (ex. 256/9 = 28 px)
	var tiles []model.TileDescriptor

	for ty := 0; ty < TilesPerRow; ty++ {
		for tx := 0; tx < TilesPerRow; tx++ {
			// Définir les coordonnées de la tuile
			tileImg := resized.SubImage(image.Rect(tx*tileSize, ty*tileSize, (tx+1)*tileSize, (ty+1)*tileSize))

			// Analyser chaque tuile individuellement
			tileDesc := model.TileDescriptor{
				HistogramRGB:     computeHistogramRGB(tileImg),
				HistogramHSV:     computeHistogramHSV(tileImg),
				PHash:            generatePHash(tileImg),
				MeanColor:        computeMeanColor(tileImg),
				TextureSignature: computeTextureSignature(tileImg),
				ShapeSignature:   computeShapeSignature(tileImg),
			}

			// Ajouter la tuile au tableau
			tiles = append(tiles, tileDesc)
		}
	}

	// === Assemblage final du descripteur d'image ===
	desc := &model.FullImageDescriptor{
		ImageName:       filepath.Base(imagePath), // Nom du fichier uniquement
		GlobalRGB:       globalRGB,
		GlobalHSV:       globalHSV,
		GlobalPHash:     globalPHash,
		GlobalMeanColor: globalMean,
		GlobalTexture:   globalTexture,
		GlobalShape:     globalShape,
		Tiles:           tiles,
	}

	return desc, nil // Retour du descripteur complet
}