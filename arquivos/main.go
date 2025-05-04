package main

import (
	"atividade_9/avaliador"
	"atividade_9/lexer"
	"atividade_9/parser"
	"fmt"
	"log"
	"os"
)

//////////////////////////////////
// Alunos: 						//
// Yuri Gabriel da Silva Alves	//
// Cássio Andrêzza de Almeida	//
//////////////////////////////////

func main() {
	// Caminho do arquivo contendo o código YZY
	caminhoArquivo := "codigo.yzy"

	// Lê o conteúdo do arquivo
	codigoYzy, err := os.ReadFile(caminhoArquivo)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo '%s': %v", caminhoArquivo, err)
	}

	// Inicializa o analisador léxico com o código YZY lido do arquivo
	lex := &lexer.AnalisadorLexico{Texto: string(codigoYzy)}

	// Inicializa o parser com o analisador léxico
	sintatico := parser.NovoParser(lex)

	// Realiza a análise sintática
	programa, err := sintatico.Analisar()
	if err != nil {
		log.Fatalf("Erro na análise sintática: %v", err)
	}

	// Realiza a avaliação semântica
	resultado, err := avaliador.AvaliarPrograma(programa)
	if err != nil {
		log.Fatalf("Erro na avaliação semântica: %v", err)
	}

	// Exibe o resultado da execução
	fmt.Printf("Resultado da execução: %d\n", resultado)
}
