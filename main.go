package main

import (
	"fmt"

	"github.com/MrIsmail1/Golang_images_matcher/analyzer"
	"github.com/MrIsmail1/Golang_images_matcher/compare"
	"github.com/MrIsmail1/Golang_images_matcher/model"

	"os"
	"path/filepath"
	"strings"
)

/*
===== PROGRAMME PRINCIPAL DE RECONNAISSANCE D'IMAGES =====

À QUOI ÇA SERT :
Démonstrateur complet du système ! Trouve l'image la plus similaire dans une base
de données à partir d'une image de requête.
*/
func main() {

	// Image cible à analyser (pourrait venir d'arguments CLI)
	imageName := "chien13.png"
	imagePath := "banque/images/" + imageName

	// Construction du chemin du cache JSON correspondant
	// strings.TrimSuffix enlève l'extension, puis on ajoute .json
	jsonTarget := "banque/json/" + strings.TrimSuffix(imageName, filepath.Ext(imageName)) + ".json"

	var desc *model.FullImageDescriptor

	// Vérification de l'existence du cache JSON
	// os.Stat retourne une erreur si le fichier n'existe pas
	if _, err := os.Stat(jsonTarget); os.IsNotExist(err) {

		fmt.Println("🔧 Descripteur non trouvé, génération en cours...")

		// Analyse complète de l'image (opération coûteuse 1-3 secondes)
		d, err := analyzer.AnalyzeImage(imagePath)
		if err != nil {
			fmt.Println("Erreur analyse:", err)
			return // Arrêt du programme en cas d'erreur critique
		}

		// Sauvegarde du descripteur en cache pour les prochaines fois
		model.SaveDescriptor(d, jsonTarget)
		desc = d
		fmt.Println("✅ Descripteur généré :", jsonTarget)

	} else {

		fmt.Println("✅ Descripteur existant :", jsonTarget)

		// Chargement ultra-rapide depuis le JSON (quelques millisecondes)
		desc = model.LoadDescriptor(jsonTarget)
	}

	// Recherche de tous les fichiers JSON dans la banque
	bankFiles, _ := filepath.Glob("banque/json/*.json")

	// Variables pour tracker le meilleur match
	bestMatch := ""
	bestScore := 0.0

	// Parcours de tous les descripteurs de la banque
	for _, file := range bankFiles {

		// Auto-exclusion : évite de comparer l'image avec elle-même
		// PRINCIPE : Si même nom de fichier → skip cette comparaison
		if filepath.Base(file) == filepath.Base(jsonTarget) {
			continue // Passe au fichier suivant
		}

		// Chargement du descripteur de cette image de la banque
		bankDesc := model.LoadDescriptor(file)

		// Calcul du score de similarité (cœur du système !)
		score := compare.CompareDescriptors(desc, bankDesc)

		// Affichage du résultat de cette comparaison
		fmt.Printf("🔹 %s : %.2f%% de similarité\n", bankDesc.ImageName, score)

		// Mise à jour du meilleur match si ce score est supérieur
		if score > bestScore {
			bestScore = score
			bestMatch = bankDesc.ImageName
		}
	}

	// Annonce du gagnant ou d'échec
	if bestMatch != "" {
		fmt.Printf("\n🏆 Meilleure correspondance : %s avec %.2f%%\n", bestMatch, bestScore)
	} else {
		fmt.Println("❌ Aucune correspondance trouvée.")
	}
}
