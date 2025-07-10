# 🔍 Go Image Search - Système de Recherche d'Images par Similarité

## 📋 Table des Matières

- [Vue d'ensemble](#vue-densemble)
- [Architecture du projet](#architecture-du-projet)
- [Technologies utilisées](#technologies-utilisées)
- [Algorithmes et méthodes](#algorithmes-et-méthodes)
- [Structure du code](#structure-du-code)
- [Installation et configuration](#installation-et-configuration)
- [Utilisation](#utilisation)
- [Exemples de résultats](#exemples-de-résultats)
- [Détails techniques](#détails-techniques)
- [Optimisations futures](#optimisations-futures)

## 🎯 Vue d'ensemble

Ce projet est un **système de recherche d'images par similarité** développé en **Go (Golang)**. Nous avons conçu cette solution pour permettre de retrouver des images similaires dans une banque d'images en utilisant plusieurs techniques avancées d'analyse visuelle.

### Objectifs du projet

Nous avons voulu créer un système capable de :

- 🔍 **Analyser automatiquement** les caractéristiques visuelles d'une image
- 📊 **Comparer efficacement** plusieurs images entre elles
- 🏆 **Identifier les correspondances** les plus pertinantes
- 💾 **Optimiser les performances** avec un système de cache JSON intelligent

## 🏗️ Architecture du projet

Nous avons refactorisé l'architecture pour une **séparation claire des responsabilités** et **éviter les dépendances circulaires** :

```
Go_image_search/
├── 📁 config/              # Configuration centralisée
│   └── constants.go        # Paramètres du système
├── 📁 analyser-utils/      # Modules d'analyse spécialisés
│   ├── 🎨 color/          # Histogrammes RGB/HSV, couleurs moyennes
│   ├── 🔢 hash/           # Hash perceptuel (pHash) et DCT
│   ├── 🧮 math/           # Fonctions mathématiques utilitaires
│   ├── 🔺 shape/          # Détection de contours et formes
│   └── 🌫️ texture/        # Analyse de texture et rugosité
├── 📁 analyzer/            # Orchestrateur principal d'analyse
├── 📁 compare-utils/       # Métriques de comparaison
├── 📁 compare/             # Moteur de comparaison principal
├── 📁 model/               # Structures de données et persistance
├── 📁 banque/              # Base de données d'images
│   ├── 🖼️ images/         # Images de référence (JPG, PNG)
│   └── 📄 json/           # Descripteurs pré-calculés (cache)
└── 📄 main.go             # Point d'entrée et démonstration
```

## 🛠️ Technologies utilisées

Nous avons sélectionné **Go** comme langage principal pour ses avantages :

- ⚡ **Performance élevée** pour le traitement d'images
- 🔧 **Simplicité de déploiement** (binaire autonome)
- 📚 **Écosystème riche** pour le traitement d'images
- 🧩 **Modularité** excellente pour l'architecture refactorisée

### Dépendances principales

```go
golang.org/x/image v0.26.0  // Extensions pour le traitement d'images
```

## 🧠 Algorithmes et méthodes

### 1. Analyse multi-niveaux

Notre **approche hybride** combine :

#### 🌍 **Analyse globale de l'image**

- **Histogrammes RGB et HSV** : distribution des couleurs avec 64 bins
- **Couleur moyenne** : teinte dominante de l'image entière
- **Hash perceptuel (pHash)** : signature binaire 64 bits robuste
- **Signature de texture** : mesure de rugosité via variations locales
- **Signature de forme** : densité de contours avec filtre Sobel

#### 🔍 **Analyse locale par tuiles**

Division en **81 tuiles** (9×9) pour :
- Capturer les **détails régionaux** non visibles globalement
- Robustesse face aux **occlusions partielles**
- **Correspondance fine** entre zones similaires

### 2. Algorithmes de traitement optimisés

#### **Configuration centralisée**

```go
// config/constants.go
const (
    StandardSize = 256  // Taille standard (256×256 pixels)
    TilesPerRow  = 9    // Grille 9×9 = 81 tuiles
    Bins         = 64   // Résolution des histogrammes
)
```

#### **Hash perceptuel robuste**

```go
// analyser-utils/hash/generatePHash.go
func GeneratePHash(img image.Image) string
```

- **DCT 2D** pour transformation fréquentielle
- **Binarisation adaptative** selon moyenne des coefficients
- **Résistance** aux modifications légères (compression, redimensionnement)

#### **Analyse couleur avancée**

```go
// analyser-utils/color/
- ComputeHistogramRGB()  // Distribution RGB avec 64 bins
- ComputeHistogramHSV()  // HSV plus robuste aux variations d'éclairage
- ComputeMeanColor()     // Couleur dominante globale
- RgbToHsv()            // Conversion d'espace colorimétrique
```

### 3. Métriques de similarité intelligentes

Score composite avec **pondération optimisée** :

```go
finalScore := (globalScore*0.65 + avgTileScore*0.35) * 100
```

**Pondération des caractéristiques :**
- 🎯 **Hash perceptuel** : 25% (distance de Hamming)
- 🔺 **Formes/contours** : 25% (signature Sobel)
- 🎨 **Couleur moyenne** : 15% (distance euclidienne)
- 🌫️ **Texture** : 15% (variations d'intensité)
- 📊 **Histogrammes RGB/HSV** : 20% (distance Manhattan)

## 📂 Structure du code refactorisée

### Module `config/`
**Configuration centralisée** pour éviter les dépendances circulaires :
```go
// constants.go - Point unique de vérité
const StandardSize = 256  // Optimisé performance/qualité
const TilesPerRow = 9     // Équilibre détail/vitesse  
const Bins = 64          // Résolution histogrammes
```

### Module `analyser-utils/`
**Fonctions spécialisées** par domaine d'analyse :

#### `color/` - Analyse couleur
```go
func ComputeHistogramRGB(img image.Image) map[string][]int
func ComputeHistogramHSV(img image.Image) map[string][]int  
func ComputeMeanColor(img image.Image) [3]float64
func RgbToHsv(r, g, b uint8) (float64, float64, float64)
```

#### `hash/` - Hash perceptuel
```go
func GeneratePHash(img image.Image) string
func dct2D(img *image.Gray) [][]float64
func averageDCT(dct [][]float64) float64
```

#### `texture/` - Analyse texture
```go
func ComputeTextureSignature(img image.Image) float64
```

#### `shape/` - Détection formes
```go
func ComputeShapeSignature(img image.Image) float64
```

#### `math/` - Utilitaires mathématiques
```go
func AbsDiff(a, b uint8) uint8
func Dct2D(img *image.Gray) [][]float64
```

### Module `analyzer/`
**Orchestrateur principal** simplifié :
```go
func AnalyzeImage(imagePath string) (*model.FullImageDescriptor, error)
```
- Coordonne tous les analyseurs spécialisés
- Gère le redimensionnement et la standardisation
- Assemble le descripteur final multi-niveaux

### Module `compare-utils/`
**Métriques de comparaison** spécialisées :
```go
func EuclideanDistance(c1, c2 [3]float64) float64
func HammingDistance(hash1, hash2 string) int  
func CompareHistograms(h1, h2 map[string][]int) float64
```

### Module `compare/`
**Moteur de comparaison intelligent** :
```go
func CompareDescriptors(desc1, desc2 *model.FullImageDescriptor) float64
```

## 🚀 Installation et configuration

### Prérequis
- **Go 1.24.1** ou supérieur
- **Git** pour cloner le projet

### Installation
```bash
git clone https://github.com/MrIsmail1/Golang_images_matcher.git
cd Go_image_search
go mod tidy
```

## 📖 Utilisation

### Utilisation basique
1. **Placer vos images** dans `banque/images/`
2. **Modifier l'image cible** dans `main.go` :
```go
imageName := "votre_image.jpg"
```
3. **Exécuter** :
```bash
go run main.go
```

### Exemple de sortie
```
🔧 Descripteur non trouvé, génération en cours...
✅ Descripteur généré : banque/json/chien8.json
🔹 cala1.jpg : 23.45% de similarité
🔹 cala2.jpg : 67.89% de similarité
🔹 chien1.png : 12.34% de similarité

🏆 Meilleure correspondance : cala2.jpg avec 67.89%
```

## 📊 Exemples de résultats

### Seuils de qualité recommandés
- **90-100%** : Images quasi-identiques
- **70-89%** : Images très similaires
- **50-69%** : Images similaires
- **30-49%** : Similarité faible
- **0-29%** : Images différentes

## 🔬 Détails techniques

### Avantages de l'architecture refactorisée
✅ **Modulaire** : Chaque domaine dans son package  
✅ **Évolutive** : Facile d'ajouter de nouveaux descripteurs  
✅ **Testable** : Modules indépendants  
✅ **Maintenable** : Responsabilités claires  
✅ **Performante** : Cache JSON intelligent

### Système de cache optimisé
- **Premier run** : Analyse complète + sauvegarde JSON
- **Runs suivants** : Chargement ultra-rapide (100-1000× plus rapide)
- **Gestion automatique** : Détection et régénération si nécessaire

## 🚀 Optimisations futures

### Performance
- 🔄 **Parallélisation** avec goroutines
- 🗄️ **Base de données** pour gros volumes
- ⚡ **Index LSH** pour recherche approximative

### Algorithmes
- 🤖 **Deep learning** avec CNNs pré-entraînés
- 🎯 **SIFT/SURF** pour points d'intérêt
- 📐 **Invariance géométrique** avancée

### Interface
- 🌐 **API REST** pour intégration web
- 🖥️ **Interface graphique** moderne
- 🔧 **Configuration dynamique** des paramètres

---

**Développé avec ❤️ en Go** - Architecture modulaire pour la recherche d'images par similarité