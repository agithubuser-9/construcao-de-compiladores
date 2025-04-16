// teste_programa_ev_test.go
package main

import (
	"strings"
	"testing"

	"atividade_8/lexer"
	"atividade_8/parser"
	"atividade_8/avaliador"
)

func TestProgramsEV(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		// Programa sem declaração: expressão final simples.
		{"=10", 10},
		// Programa com uma declaração.
		{"x=7+3;=x", 10},
		// Programa com duas declarações.
		{"x=5;y=x+3;=x+y", 13},
		// Programa do PDF.
		{"x=(7+4)*12;y=x*3+11;=(x*y)+(x*11)+(y*13)", 60467},
		// Casos adicionais.
		{"a=10;b=a*2;=a+b", 30},
		{"z=100;=z-50", 50},
	}

	for _, tc := range tests {
		input := strings.ReplaceAll(tc.input, " ", "")
		lex := &lexer.AnalisadorLexico{Texto: input}
		p := parser.NovoParser(lex)
		prog, err := p.AnalisarPrograma()
		if err != nil {
			t.Errorf("Erro ao analisar '%s': %v", tc.input, err)
			continue
		}
		result, err := avaliador.AvaliarPrograma(prog)
		if err != nil {
			t.Errorf("Erro ao avaliar '%s': %v", tc.input, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("Para '%s', esperado %d, mas obteve %d", tc.input, tc.expected, result)
		}
	}
}

