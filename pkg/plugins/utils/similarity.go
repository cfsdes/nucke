package utils

import (
	"strings"
)

// Função para calcular a similaridade entre dois textos
func TextSimilarity(text1, text2 string) float64 {
    // Converter os textos para letras minúsculas para uma comparação não sensível a maiúsculas e minúsculas
    text1 = strings.ToLower(text1)
    text2 = strings.ToLower(text2)

    // Calcular a distância de edição entre os textos
    distance := levenshteinDistance(text1, text2)

    // Calcular a similaridade de Jaccard (1 - distância/len(união dos caracteres))
    unionLength := float64(len(text1) + len(text2) - distance)
    jaccardSimilarity := unionLength / float64(len(text1)+len(text2))

    return jaccardSimilarity
}

// Função para calcular a distância de edição entre dois textos
func levenshteinDistance(s, t string) int {
    m := len(s)
    n := len(t)

    // Criar uma matriz para armazenar os resultados dos subproblemas da distância de edição
    dp := make([][]int, m+1)
    for i := 0; i <= m; i++ {
        dp[i] = make([]int, n+1)
        dp[i][0] = i
    }
    for j := 0; j <= n; j++ {
        dp[0][j] = j
    }

    // Preencher a matriz usando a abordagem bottom-up
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            cost := 1
            if s[i-1] == t[j-1] {
                cost = 0
            }
            dp[i][j] = min(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+cost)
        }
    }

    return dp[m][n]
}

func min(a, b, c int) int {
    if a < b {
        if a < c {
            return a
        }
    } else if b < c {
        return b
    }
    return c
}