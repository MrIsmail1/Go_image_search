package analyzer

import (
	"github.com/MrIsmail1/Golang_images_matcher/analyser-utils/color"
	"github.com/MrIsmail1/Golang_images_matcher/analyser-utils/hash"
	"github.com/MrIsmail1/Golang_images_matcher/analyser-utils/shape"
	"github.com/MrIsmail1/Golang_images_matcher/analyser-utils/texture"
	"github.com/MrIsmail1/Golang_images_matcher/config"
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

/*
===== ANALYSE COMPLÈTE D'UNE IMAGE - CHEF D'ORCHESTRE =====

À QUOI ÇA SERT :
C'est le "chef d'orchestre" du système ! Coordonne tous les analyseurs spécialisés
pour transformer une image brute en descripteur complet multi-niveaux.

STRATÉGIE MULTI-ÉCHELLES :
1. ANALYSE GLOBALE : Caractéristiques de l'image entière
  - Couleurs dominantes, texture générale, structure globale

2. ANALYSE LOCALE : Division en tuiles (9×9 = 81 zones)
  - Détecte les variations régionales dans l'image
  - Complète la vision globale avec des détails fins

Paramètre :
- imagePath : chemin vers l'image à analyser

Retour :
- Descripteur complet avec toutes les caractéristiques
- Erreur si problème de lecture/analyse
*/
func AnalyzeImage(imagePath string) (*model.FullImageDescriptor, error) {

	// Ouverture du fichier avec gestion d'erreur
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err // Fichier inexistant, permissions insuffisantes, etc.
	}
	defer file.Close() // Fermeture automatique même en cas d'erreur

	// Décodage automatique du format (JPEG/PNG/GIF grâce aux imports _)
	// MAGIE : image.Decode détecte automatiquement le format depuis les premiers bytes
	srcImg, _, err := image.Decode(file)
	if err != nil {
		return nil, err // Image corrompue, format non supporté, etc.
	}

	// Redimensionnement à la taille standard configurée (256×256 par défaut)
	// POURQUOI STANDARDISER ?
	// - Base de comparaison uniforme entre toutes les images
	// - Performance prévisible (même temps de traitement)
	// - Descripteurs comparables (même échelle)
	resized := image.NewRGBA(image.Rect(0, 0, config.StandardSize, config.StandardSize))
	drawx.ApproxBiLinear.Scale(resized, resized.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	// Délégation aux modules spécialisés pour chaque type d'analyse
	// AVANTAGE : Chaque module fait ce qu'il sait le mieux faire

	globalRGB := color.ComputeHistogramRGB(resized)           // Distribution des couleurs RGB
	globalHSV := color.ComputeHistogramHSV(resized)           // Distribution des couleurs HSV (complémentaire)
	globalPHash := hash.GeneratePHash(resized)                // Signature binaire robuste 64 bits
	globalMean := color.ComputeMeanColor(resized)             // Couleur dominante simple
	globalTexture := texture.ComputeTextureSignature(resized) // Rugosité/finesse globale
	globalShape := shape.ComputeShapeSignature(resized)       // Densité de contours/formes

	// Calcul de la taille d'une tuile individuelle
	// EXEMPLE : 256 pixels ÷ 9 tuiles = ~28 pixels par tuile
	tileSize := config.StandardSize / config.TilesPerRow
	var tiles []model.TileDescriptor

	// Double boucle pour créer la grille 9×9 de tuiles
	for ty := 0; ty < config.TilesPerRow; ty++ { // Pour chaque ligne de tuiles
		for tx := 0; tx < config.TilesPerRow; tx++ { // Pour chaque colonne de tuiles

			// Extraction de la portion d'image correspondant à cette tuile
			// COORDONNÉES : (tx*tileSize, ty*tileSize) vers ((tx+1)*tileSize, (ty+1)*tileSize)
			tileImg := resized.SubImage(image.Rect(
				tx*tileSize, ty*tileSize, // Coin supérieur gauche
				(tx+1)*tileSize, (ty+1)*tileSize, // Coin inférieur droit
			))

			// Application des MÊMES analyses que pour l'image globale
			// MAIS seulement sur cette petite zone
			// AVANTAGE : Détecte les variations locales ignorées dans l'analyse globale
			tileDesc := model.TileDescriptor{
				HistogramRGB:     color.ComputeHistogramRGB(tileImg),       // Couleurs locales
				HistogramHSV:     color.ComputeHistogramHSV(tileImg),       // HSV local
				PHash:            hash.GeneratePHash(tileImg),              // Signature locale
				MeanColor:        color.ComputeMeanColor(tileImg),          // Couleur dominante locale
				TextureSignature: texture.ComputeTextureSignature(tileImg), // Rugosité locale
				ShapeSignature:   shape.ComputeShapeSignature(tileImg),     // Contours locaux
			}

			// Ajout de cette tuile analysée à la collection
			tiles = append(tiles, tileDesc)
		}
	}

	// Construction de la structure finale qui contient TOUT
	// ORGANISATION :
	// - Métadonnées : nom de fichier
	// - Niveau global : caractéristiques de l'image entière
	// - Niveau local : 81 tuiles avec leurs caractéristiques individuelles
	desc := &model.FullImageDescriptor{
		ImageName:       filepath.Base(imagePath), // Nom du fichier seulement (sans chemin)
		GlobalRGB:       globalRGB,                // Couleurs globales RGB
		GlobalHSV:       globalHSV,                // Couleurs globales HSV
		GlobalPHash:     globalPHash,              // Signature structurelle globale
		GlobalMeanColor: globalMean,               // Teinte dominante globale
		GlobalTexture:   globalTexture,            // Rugosité globale
		GlobalShape:     globalShape,              // Richesse en formes globale
		Tiles:           tiles,                    // Collection des 81 tuiles analysées
	}

	return desc, nil // Mission accomplie ! Descripteur complet prêt à l'emploi
}
