package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// generateMermaidFromMatrix prend une matrice de transition et génère un diagramme d'état Mermaid.
func generateMermaidFromMatrix(matrix [][]string) string {
	if len(matrix) < 2 || len(matrix[0]) < 2 {
		return ""
	}

	// Mapping des noms d'événements aux identifiants de nœuds
	eventMap := make(map[string]string)
	nodeCount := 0
	for _, event := range matrix[0][1:] {
		nodeCount++
		eventMap[event] = fmt.Sprintf("e%d", nodeCount-1)
	}

	var builder strings.Builder
	builder.WriteString("graph LR\n")
	for i := 1; i < len(matrix); i++ {
		from := matrix[i][0]
		for j := 1; j < len(matrix[i]); j++ {
			to := matrix[0][j]
			weight := matrix[i][j]
			if weight != "0" {
				fromNode := eventMap[from]
				toNode := eventMap[to]
				// Utilisez la nouvelle syntaxe avec des étiquettes pour les nœuds
				builder.WriteString(fmt.Sprintf("%s[%s] -- %s --> %s[%s]\n", fromNode, from, weight, toNode, to))
			}
		}
	}
	return builder.String()
}

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'
	data, err := reader.ReadAll()
	return data, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <filename.csv>")
		return
	}
	filename := os.Args[1]

	matrix, err := readCSV(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	diagram := generateMermaidFromMatrix(matrix)
	fmt.Println(diagram)
}
