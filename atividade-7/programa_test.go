// programa_test.go
package main

import (
	"strings"
	"testing"
)

func TestExpressionsEC2(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"333", 333},
		{"6*7", 42},
		{"3+(4+(11+7))", 25},
		{"33+(912*11)", 10065},
		{"(427/7)+(11*(231+5))", 2657},
		{"33+(912*11)", 10065},
		{"10+(10+10)", 30},
		{"(10+10)+10", 30},
		{"10+(10+10)+10", 40},
		{"(10+10)+(10+10)", 40},
		{"10", 10},
		// Casos adicionais para demonstrar a análise sem parênteses excessivos.
		{"7+5*3", 22},  // 5*3=15, 7+15=22
		{"7+5+3", 15},  // (7+5)+3=15
		{"10-8-2", 0},  // (10-8)-2=0
		{"10-8*2", -6}, // 8*2=16, 10-16 = -6
	}

	for _, tc := range tests {
		input := strings.ReplaceAll(tc.input, " ", "")
		lexer := &AnalisadorLexico{texto: input}
		parser := NovoParser(lexer)
		expressao, err := parser.Analisar()
		if err != nil {
			t.Errorf("Erro ao analisar '%s': %v", tc.input, err)
			continue
		}
		result, err := Avaliar(expressao)
		if err != nil {
			t.Errorf("Erro ao avaliar '%s': %v", tc.input, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("Para '%s', esperado %d, mas obteve %d", tc.input, tc.expected, result)
		}
	}
}
