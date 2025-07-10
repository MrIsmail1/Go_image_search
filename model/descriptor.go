package model

import (
	"encoding/json" // Pour l'encodage JSON
	"fmt"
	"os" // Pour les opérations sur les fichiers
)

/*
	===== STRUCTURE D'UNE TUILE (Tile) =====

Chaque image est divisée en petites tuiles pour une analyse plus fine.
Chaque tuile possède plusieurs descripteurs (couleur, texture, etc.).
*/
type TileDescriptor struct {
	HistogramRGB     map[string][]int `json:"histogram_rgb"`     // Histogramme des couleurs RGB (rouge, vert, bleu)
	HistogramHSV     map[string][]int `json:"histogram_hsv"`     // Histogramme des couleurs HSV (teinte, saturation, valeur)
	PHash            string           `json:"phash"`             // Hash perceptuel (pour comparer visuellement les images)
	MeanColor        [3]float64       `json:"mean_color"`        // Moyenne des couleurs [R, G, B] pour la tuile
	TextureSignature float64          `json:"texture_signature"` // Signature de texture (mesure de variation locale)
	ShapeSignature   float64          `json:"shape_signature"`   // Signature des formes (contours)
}

/*
	===== STRUCTURE D'UNE IMAGE ENTIÈRE =====

Cette structure regroupe les descripteurs globaux et ceux des tuiles.
*/
type FullImageDescriptor struct {
	ImageName       string           `json:"image_name"`        // Nom de l'image (souvent le nom de fichier)
	GlobalRGB       map[string][]int `json:"global_rgb"`        // Histogramme global RGB
	GlobalHSV       map[string][]int `json:"global_hsv"`        // Histogramme global HSV
	GlobalPHash     string           `json:"global_phash"`      // Hash perceptuel global
	GlobalMeanColor [3]float64       `json:"global_mean_color"` // Moyenne des couleurs globales [R, G, B]
	GlobalTexture   float64          `json:"global_texture"`    // Texture globale de l'image
	Tiles           []TileDescriptor `json:"tiles"`             // Liste des descripteurs pour chaque tuile
	ShapeSignature  float64          `json:"shape_signature"`   // Signature des formes (contours)
	GlobalShape     float64          `json:"global_shape"`      // 👈 AJOUTE CETTE LIGNE
}

/*
	===== SAUVEGARDE D'UN DESCRIPTEUR DANS UN FICHIER JSON =====

Cette fonction permet d'enregistrer un descripteur d'image complet (FullImageDescriptor)
dans un fichier au format JSON, afin de pouvoir le réutiliser plus tard sans devoir
réanalyser l'image.

Paramètres :
- desc : le descripteur à sauvegarder (pointeur vers FullImageDescriptor)
- outputPath : chemin complet du fichier .json dans lequel sauvegarder les données

Retour :
- Une erreur si la sauvegarde échoue, sinon nil
*/
func saveDescriptor(desc *FullImageDescriptor, outputPath string) error {
	// Création (ou écrasement) du fichier de sortie
	file, err := os.Create(outputPath)
	if err != nil {
		return err // Retourne l'erreur si le fichier ne peut pas être créé
	}
	defer file.Close() // S'assure que le fichier sera fermé même en cas d'erreur plus bas

	// Préparation de l'encodeur JSON avec indentation pour une meilleure lisibilité
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // JSON "joli" avec indentation

	// Encodage du descripteur dans le fichier
	return encoder.Encode(desc)
}

func SaveDescriptor(desc *FullImageDescriptor, outputPath string) error {
	return saveDescriptor(desc, outputPath)
}

func LoadDescriptor(inputPath string) *FullImageDescriptor {
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Erreur ouverture fichier JSON:", err)
		return nil
	}
	defer file.Close()

	var desc FullImageDescriptor
	if err := json.NewDecoder(file).Decode(&desc); err != nil {
		fmt.Println("Erreur décodage JSON:", err)
		return nil
	}
	return &desc
}
