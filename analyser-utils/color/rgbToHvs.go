package color

import "math"

/*
===== CONVERSION RGB → HSV =====

À QUOI ÇA SERT :
Convertit une couleur du format RGB (écrans/capteurs) vers HSV (perception humaine).
Sépare la couleur pure (H) de sa luminosité (V) et de sa pureté (S).

POURQUOI CETTE CONVERSION :
- RGB mélange couleur et luminosité (difficile à analyser)
- HSV sépare ces aspects (plus intuitif)
- Permet de chercher "du rouge" peu importe sa luminosité
- Robuste aux variations d'éclairage

QU'EST-CE QUE HSV :
- H (Hue/Teinte) : Position sur la roue des couleurs (0°=rouge, 120°=vert, 240°=bleu)
- S (Saturation) : Pureté de la couleur (0=gris, 1=couleur pure)
- V (Value/Valeur) : Luminosité (0=noir, 1=couleur maximale)

Paramètres :
- r, g, b : composantes RGB (0-255)

Retour :
- h : teinte en degrés (0-360)
- s : saturation (0-1)
- v : valeur/luminosité (0-1)
*/
func rgbToHsv(r, g, b uint8) (float64, float64, float64) {

	// Normalisation des valeurs RGB de 0-255 vers 0-1
	// POURQUOI cette normalisation ?
	// - Les formules HSV utilisent des valeurs 0-1
	// - Plus précis de travailler en float pour les calculs
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	// Recherche des valeurs min et max parmi R, G, B
	// CES VALEURS SERVENT À :
	// - max : détermine la luminosité finale (V)
	// - min : aide à calculer la saturation
	// - delta (max-min) : mesure l'étendue des couleurs
	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	delta := max - min

	var h float64
	// Si delta = 0, alors R=G=B → couleur grise → teinte indéfinie
	if delta != 0 {
		// La teinte dépend de quelle couleur domine (R, G ou B)
		// PRINCIPE : On est dans quel "secteur" de la roue des couleurs ?
		switch max {
		case rf: // Rouge dominant → secteur rouge-jaune ou rouge-magenta
			h = math.Mod(((gf - bf) / delta), 6)
		case gf: // Vert dominant → secteur vert-cyan ou vert-jaune
			h = ((bf - rf) / delta) + 2
		case bf: // Bleu dominant → secteur bleu-magenta ou bleu-cyan
			h = ((rf - gf) / delta) + 4
		}

		// Conversion en degrés (0-360°)
		// POURQUOI * 60 ? La roue HSV est divisée en 6 secteurs de 60° chacun
		h *= 60

		// Gestion des valeurs négatives (ramener dans [0-360°])
		if h < 0 {
			h += 360
		}
	}

	s := 0.0
	// Saturation = delta / max (sauf si max = 0 pour éviter division par zéro)
	// INTERPRÉTATION :
	// - Si max = delta (une couleur = max, les autres = 0) → S = 1 (couleur pure)
	// - Si delta = 0 (toutes égales) → S = 0 (gris)
	if max != 0 {
		s = delta / max
	}

	// La valeur V est simplement la composante maximale
	// PRINCIPE : La luminosité d'une couleur RGB = sa composante la plus brillante
	v := max

	return h, s, v
}
