package main

import (
	"atividade_9/avaliador"
	"atividade_9/geradorAssembly"
	"atividade_9/lexer"
	"atividade_9/parser"
	"fmt"
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
	//testes durante o desenvolvimento
	//caminhoArquivo := "codigo.yzy"
	/*
		// Lê o conteúdo do arquivo
		codigoYzy, err := os.ReadFile(caminhoArquivo)
		if err != nil {
			log.Fatalf("Erro ao ler o arquivo '%s': %v", caminhoArquivo, err)
		}

		lex := &lexer.AnalisadorLexico{Texto: string(codigoYzy)}

		parser := &parser.ParserAtt{
			Lexer:     lex,
			Variaveis: make(map[string]int),
		}
		parser.Avancar()
		resultado, err := avaliador.AvaliarPrograma(parser.Programa())

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resultado)
		}
	*/

	files, _ := filepath.Glob("testes/*.yzy")
	os.MkdirAll("assembly", 0755)

	for _, caminho := range files {
		nome := strings.TrimSuffix(filepath.Base(caminho), ".yzy")
		arquivoAsm := filepath.Join("assembly", nome+".s")

		// Lê o arquivo-fonte .yzy
		codigo, _ := os.ReadFile(caminho)

		// Parser → AST
		lex := &lexer.AnalisadorLexico{Texto: string(codigo)}
		p := &parser.ParserAtt{
			Lexer:     lex,
			Variaveis: make(map[string]int),
		}
		p.Avancar()
		prog := p.Programa()

		// RETORNO SUCESSO X FALHA
		fmt.Printf("🧪 %s: ", nome)
		if res, err := avaliador.AvaliarPrograma(prog); err != nil {
			fmt.Printf("❌ Erro de execução: %v\n", err)
		} else {
			fmt.Printf("✅ Resultado: %d\n", res)
		}

		// Geração de Assembly
		asm := geradorAssembly.GerarCodigoPrograma(prog)
		os.WriteFile(arquivoAsm, []byte(asm), 0644)
	}

}
