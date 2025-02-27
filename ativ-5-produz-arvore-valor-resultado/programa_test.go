package main

import (
	"strings"
	"testing"
)

func TestEvaluateExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"333", 333},
		{"(6 * 7)", 42},
		{"(3 + (4 + (11 + 7)))", 25},
		{"(33 + (912 * 11))", 10065},
		{"((427 / 7) + (11 * (231 + 5)))", 2657},
		{"( 33 + ( 912 * 11 ) )", 10065},
		{"(10+(10+10))", 30},
		{"((10+10)+10)", 30},
		{"(10+(10+10)+10)", 40},
		{"((10+10)+(10+10))", 40},
		{"10", 10},
	}

	for _, tc := range tests {
		// Limpa espaços desnecessários na entrada.
		lexer := &Lexer{text: strings.TrimSpace(tc.input)}
		parser := NewParser(lexer)
		expr, err := parser.ParseExpression()
		if err != nil {
			t.Errorf("Erro ao analisar a expressão '%s': %v", tc.input, err)
			continue
		}
		result, err := Evaluate(expr)
		if err != nil {
			t.Errorf("Erro ao avaliar a expressão '%s': %v", tc.input, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("Para a expressão '%s', esperado %d, mas obteve %d", tc.input, tc.expected, result)
		}
	}
}
