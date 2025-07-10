package main

import (
	"fmt"
	"github.com/MrIsmail1/Golang_images_matcher/analyzer"
	"github.com/MrIsmail1/Golang_images_matcher/model"

	"os"
	"path/filepath"
	"strings"
)

func main() {
	imageName := "13.jpg"
	imagePath := "banque/images/" + imageName
	jsonTarget := "banque/json/" + strings.TrimSuffix(imageName, filepath.Ext(imageName)) + ".json"

	var desc *model.FullImageDescriptor
	if _, err := os.Stat(jsonTarget); os.IsNotExist(err) {
		fmt.Println("🔧 Descripteur non trouvé, génération en cours...")
		d, err := analyzer.AnalyzeImage(imagePath)
		if err != nil {
			fmt.Println("Erreur analyse:", err)
			return
		}
		model.SaveDescriptor(d, jsonTarget)
		desc = d
		fmt.Println("✅ Descripteur généré :", jsonTarget)
	} else {
		fmt.Println("✅ Descripteur existant :", jsonTarget)
		desc = model.LoadDescriptor(jsonTarget)
	}

	bankFiles, _ := filepath.Glob("banque/json/*.json")

	bestMatch := ""
	bestScore := 0.0

	for _, file := range bankFiles {
		if filepath.Base(file) == filepath.Base(jsonTarget) {
			continue
		}

		bankDesc := model.LoadDescriptor(file)
		score := compare.CompareDescriptors(desc, bankDesc)
		fmt.Printf("🔹 %s : %.2f%% de similarité\n", bankDesc.ImageName, score)

		if score > bestScore {
			bestScore = score
			bestMatch = bankDesc.ImageName
		}
	}

	if bestMatch != "" {
		fmt.Printf("\n🏆 Meilleure correspondance : %s avec %.2f%%\n", bestMatch, bestScore)
	} else {
		fmt.Println("❌ Aucune correspondance trouvée.")
	}
}
