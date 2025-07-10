package compare_utils

import "fmt"

/*
===== DISTANCE DE HAMMING ENTRE DEUX pHash =====

À QUOI ÇA SERT :
Compte le nombre de bits différents entre deux hashes perceptuels.
C'est LA métrique pour comparer des pHash - plus efficace que comparer caractère par caractère.

PRINCIPE DE LA DISTANCE DE HAMMING :
- Compare deux séquences binaires bit par bit
- Compte combien de positions ont des bits différents
- Résultat = nombre total de différences

Paramètres :
- hash1, hash2 : hashes sous forme hexadécimale (ex: "a1b2c3d4e5f67890")

Retour :
- Nombre de bits différents (0-64)
*/
func HammingDistance(hash1, hash2 string) int {

	var v1, v2 uint64

	// Conversion des chaînes hexadécimales vers entiers 64 bits
	// FONCTION : fmt.Sscanf avec format "%x" pour hexadécimal
	fmt.Sscanf(hash1, "%x", &v1) // Premier hash → entier v1
	fmt.Sscanf(hash2, "%x", &v2) // Deuxième hash → entier v2

	// XOR binaire pour marquer les positions où les hashes diffèrent
	xor := v1 ^ v2

	dist := 0 // Compteur de bits différents

	// Boucle optimisée de comptage (algorithme de Brian Kernighan)
	// PRINCIPE : À chaque itération, on élimine le bit le plus à droite
	// OPÉRATION MAGIQUE : xor &= (xor - 1)
	//
	for xor != 0 {
		dist++         // Un bit différent de plus trouvé
		xor &= xor - 1 // Élimine le bit le plus à droite (technique de Kernighan)
	}

	return dist // Retour du nombre total de bits différents
}
