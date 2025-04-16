// main.go
package main

import (
	"atividade_8/avaliador"
	"atividade_8/lexer"
	"atividade_8/parser"
	"atividade_8/gerador"
	"fmt"
)

func main() {
	// Exemplo de programa EV:
	// x=(7+4)*12; y=x*3+11; = (x*y)+(x*11)+(y*13)
	entrada := "x=(7+4)*12;y=x*3+11;=(x*y)+(x*11)+(y*13)"
	lex := &lexer.AnalisadorLexico{Texto: entrada}
	sintatico := parser.NovoParser(lex)

	prog, err := sintatico.AnalisarPrograma()
	if err != nil {
		fmt.Println("Erro na análise:", err)
		return
	}

	fmt.Println("Árvore sintática (impressa):", avaliador.Imprimir(prog))
	
	resultado, err := avaliador.AvaliarPrograma(prog)
	if err != nil {
		fmt.Println("Erro na avaliação:", err)
		return
	}
	fmt.Println("Resultado da avaliação:", resultado)

	codigo := gerador.GerarCodigo(prog)
	fmt.Println("\nCódigo Assembly Gerado:\n", codigo)
}

