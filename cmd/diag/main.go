package diag

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Name = "diag"

// generateMermaidFromMatrix prend une matrice de transition et génère un diagramme d'état Mermaid.
func generateMermaidFromMatrix(matrix [][]string, filter int) string {
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
			weight, _ := strconv.Atoi(matrix[i][j])
			if weight > filter {
				fromNode := eventMap[from]
				toNode := eventMap[to]
				builder.WriteString(fmt.Sprintf("%s[%s] -- %d --> %s[%s]\n", fromNode, from, weight, toNode, to))
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

func Run(args []string) error {
	filename := ""
	filter := 0
	fs := flag.NewFlagSet(Name, flag.ContinueOnError)
	fs.StringVar(&filename, "in", "", "csv file to convert")
	fs.IntVar(&filter, "min", 0, "min transition value to keep")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("cannot read app args : %w", err)
	}

	matrix, err := readCSV(filename)
	if err != nil {
		return fmt.Errorf("cannot read input file : %w", err)
	}
	diagram := generateMermaidFromMatrix(matrix, filter)
	fmt.Println(diagram)
	return nil
}
