package main

import (
	"atividade_7_Revisada/avaliador"
	"atividade_7_Revisada/lexer"
	"atividade_7_Revisada/parser"
	"fmt"
)

func main() {
	entrada := "25*(1)+65/(30+1)"
	//analisadorLexico := &AnalisadorLexico{texto: strings.ReplaceAll(entrada, " ", "")}
	lex := &lexer.AnalisadorLexico{Texto: entrada} //chamada sem limpeza de espaços para centralizar essa logica no lexer
	sintatico := parser.NovoParser(lex)

	expressao, err := sintatico.Analisar()
	if err != nil {
		fmt.Println("Erro na análise:", err)
		return
	}

	fmt.Println("Árvore sintática:", avaliador.Imprimir(expressao))

	resultado, err := avaliador.Avaliar(expressao)
	if err != nil {
		fmt.Println("Erro na avaliação:", err)
		return
	}
	fmt.Println("Resultado:", resultado)
}
