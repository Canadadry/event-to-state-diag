package main

import (
	"testing"
)

func TestGenerateMermaidFromMatrix(t *testing.T) {
	tests := map[string]struct {
		in  [][]string
		out string
	}{
		"basic transition": {
			in: [][]string{
				{"From/To", "Event1", "Event2", "Event3"},
				{"Event1", "0", "2", "3"},
				{"Event2", "1", "0", "0"},
				{"Event3", "0", "0", "0"},
			},
			out: `graph LR
e0[Event1] -- 2 --> e1[Event2]
e0[Event1] -- 3 --> e2[Event3]
e1[Event2] -- 1 --> e0[Event1]
`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := generateMermaidFromMatrix(tt.in)
			if got != tt.out {
				t.Errorf("generateMermaidFromMatrix() got = %v, want %v", got, tt.out)
			}
		})
	}
}
