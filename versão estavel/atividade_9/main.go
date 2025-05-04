package main

import (
	"atividade_9/lexer"
	"atividade_9/parser"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//////////////////////////////////
// Alunos: 						//
// Yuri Gabriel da Silva Alves	//
// Cássio Andrêzza de Almeida	//
//////////////////////////////////

func main() {

	//caminhoArquivo := "codigo.yzy"

	// Lê o conteúdo do arquivo
	//codigoYzy, err := os.ReadFile(caminhoArquivo)
	//if err != nil {
	//	log.Fatalf("Erro ao ler o arquivo '%s': %v", caminhoArquivo, err)
	//}

	//lex := &lexer.AnalisadorLexico{Texto: string(codigoYzy)}

	arquivos, _ := filepath.Glob("testes/*.yzy")

	for _, caminho := range arquivos {
		nomeBase := strings.TrimSuffix(filepath.Base(caminho), ".yzy")
		esperadoPath := filepath.Join("resultados", nomeBase+".txt")

		codigo, _ := os.ReadFile(caminho)
		esperado, _ := os.ReadFile(esperadoPath)

		fmt.Printf("Exeutando %s... ", nomeBase)

		lex := &lexer.AnalisadorLexico{Texto: string(codigo)}
		p := parser.ParserAtt{
			Lexer:     lex,
			Variaveis: make(map[string]int),
		}
		p.Avancar()

		r, w, _ := os.Pipe()
		orig := os.Stdout
		os.Stdout = w

		p.Programa()

		w.Close()
		os.Stdout = orig

		out, _ := io.ReadAll(r)
		saida := strings.TrimSpace(string(out))
		esperadoStr := strings.TrimSpace(string(esperado))

		if saida == esperadoStr {
			fmt.Println("✅ OK")
		} else {
			fmt.Println("❌ ERRO")
			fmt.Printf("   Esperado: %s\n", esperadoStr)
			fmt.Printf("   Obtido  : %s\n", saida)
		}
	}
}
