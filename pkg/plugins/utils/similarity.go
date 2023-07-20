package utils

import (
	"strings"
    "github.com/texttheater/golang-levenshtein/levenshtein"
)

// Função para calcular a similaridade entre dois textos
func TextSimilarity(text1, text2 string) float64 {
    // Converter os textos para letras minúsculas para uma comparação não sensível a maiúsculas e minúsculas
    text1 = strings.ToLower(text1)
    text2 = strings.ToLower(text2)

    // Calcular o Levenshtein distance entre os dois textos
	distance := levenshtein.DistanceForStrings([]rune(text1), []rune(text2), levenshtein.DefaultOptions)
    
    // Calcular o tamanho da maior string
	maxLength := len(text1)
	if len(text2) > maxLength {
		maxLength = len(text2)
	}

	// Calcular a similaridade como uma porcentagem
	similarity := 1.0 - float64(distance)/float64(maxLength)
    
    return similarity
}