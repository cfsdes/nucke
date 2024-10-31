package utils

import (
	"math"
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

// Função para calcular a similaridade entre dois response bodies com base no comprimento
func ResponseSimilarity(resBody1, resBody2 string) float64 {
	length1 := len(resBody1)
	length2 := len(resBody2)

	if length1 == 0 && length2 == 0 {
		return 1.0 // 100% similar se ambos forem vazios
	}

	// Calcula a diferença de comprimento absoluto
	lengthDifference := math.Abs(float64(length1 - length2))

	// Encontra o comprimento da maior string
	maxLength := float64(length1)
	if length2 > length1 {
		maxLength = float64(length2)
	}

	// Calcula a similaridade como uma porcentagem inversa da diferença de comprimento
	similarity := 1.0 - (lengthDifference / maxLength)

	return similarity
}
