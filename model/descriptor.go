package model

import (
	"encoding/json" // Pour la sérialisation JSON automatique
	"fmt"
	"os" // Pour les opérations sur fichiers
)

/*
===== STRUCTURE D'UN DESCRIPTEUR DE TUILE =====

À QUOI ÇA SERT :
Stocke tous les descripteurs calculés pour UNE seule tuile (portion 28×28 de l'image).
Chaque tuile est analysée indépendamment avec les mêmes techniques que l'image globale.
*/
type TileDescriptor struct {
	// Distribution des couleurs RGB dans cette tuile
	// FORMAT : {"r": [64 bins], "g": [64 bins], "b": [64 bins]}
	HistogramRGB map[string][]int `json:"histogram_rgb"`

	// Distribution des couleurs HSV dans cette tuile (complémentaire de RGB)
	// FORMAT : {"h": [64 bins], "s": [64 bins], "v": [64 bins]}
	HistogramHSV map[string][]int `json:"histogram_hsv"`

	// Hash perceptuel de cette tuile (signature binaire 64 bits)
	PHash string `json:"phash"`

	// Couleur moyenne de cette tuile
	// FORMAT : [Rouge, Vert, Bleu] avec valeurs 0-255
	MeanColor [3]float64 `json:"mean_color"`

	// Signature de texture de cette tuile
	TextureSignature float64 `json:"texture_signature"`

	// Signature de forme de cette tuile
	ShapeSignature float64 `json:"shape_signature"`
}

/*
===== STRUCTURE COMPLÈTE D'UN DESCRIPTEUR D'IMAGE =====

À QUOI ÇA SERT :
Structure principale qui contient TOUTE l'information extraite d'une image.
C'est la "carte d'identité" complète d'une image pour la reconnaissance.
*/
type FullImageDescriptor struct {

	// Nom du fichier image (sans le chemin complet)
	// UTILITÉ : Identification et affichage des résultats
	ImageName string `json:"image_name"`

	// Histogramme RGB global - Distribution générale des couleurs
	GlobalRGB map[string][]int `json:"global_rgb"`

	// Histogramme HSV global - Vision complémentaire des couleurs
	GlobalHSV map[string][]int `json:"global_hsv"`

	// Hash perceptuel global - Signature structurelle de l'image entière
	GlobalPHash string `json:"global_phash"`

	// Couleur moyenne globale - Teinte dominante de toute l'image
	GlobalMeanColor [3]float64 `json:"global_mean_color"`

	// Signature de texture globale - Rugosité moyenne de l'image
	GlobalTexture float64 `json:"global_texture"`

	// Signature de forme globale - Densité des contours dans l'image
	GlobalShape float64 `json:"global_shape"`

	Tiles []TileDescriptor `json:"tiles"`
}

/*
===== SAUVEGARDE D'UN DESCRIPTEUR AU FORMAT JSON =====

À QUOI ÇA SERT :
Sérialise un descripteur complet dans un fichier JSON pour créer un cache persistant.
Évite de recalculer les descripteurs à chaque utilisation (gain de temps énorme).

Paramètres :
- desc : descripteur à sauvegarder
- outputPath : chemin du fichier JSON de destination

Retour :
- nil si succès, erreur si échec
*/
func saveDescriptor(desc *FullImageDescriptor, outputPath string) error {

	// Création du fichier de sortie (écrase s'il existe)
	file, err := os.Create(outputPath)
	if err != nil {
		return err // Erreur : disque plein, permissions, chemin invalide, etc.
	}
	defer file.Close() // Fermeture garantie même en cas d'erreur

	// Configuration de l'encodeur JSON pour un rendu "joli"
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Indentation de 2 espaces pour lisibilité

	// Sérialisation automatique : Go → JSON
	// MAGIE : Les tags `json:"..."` contrôlent le mapping
	return encoder.Encode(desc)
}

/*
===== FONCTION PUBLIQUE DE SAUVEGARDE =====

Wrapper public pour respecter les conventions Go.
(Fonction publique = première lettre majuscule)
*/
func SaveDescriptor(desc *FullImageDescriptor, outputPath string) error {
	return saveDescriptor(desc, outputPath)
}

/*
===== CHARGEMENT D'UN DESCRIPTEUR DEPUIS JSON =====

À QUOI ÇA SERT :
Désérialise un fichier JSON pour reconstruire un descripteur complet.
Permet de récupérer instantanément les résultats d'analyses précédentes.

Paramètre :
- inputPath : chemin vers le fichier JSON à charger

Retour :
- Pointeur vers descripteur reconstitué, ou nil si échec
*/
func LoadDescriptor(inputPath string) *FullImageDescriptor {

	// Tentative d'ouverture du fichier JSON
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Erreur ouverture fichier JSON:", err)
		return nil // Fichier inexistant, permissions, etc.
	}
	defer file.Close()

	// Structure vide pour recevoir les données désérialisées
	var desc FullImageDescriptor

	// Désérialisation JSON → Go
	if err := json.NewDecoder(file).Decode(&desc); err != nil {
		fmt.Println("Erreur décodage JSON:", err)
		return nil // JSON malformé, champs manquants, etc.
	}

	return &desc // Succès : descripteur reconstitué
}
