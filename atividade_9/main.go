package main

import (

	// certifique-se de que o nome está certo

	"atividade_9/lexer"
	"atividade_9/parser"
	"fmt"
)

// Alunos:
// Yuri Gabriel da Silva Alves
// Cássio Andrêzza de Almeida

func main() {
	// Exemplo de entrada para testar
	entrada := `{
		a = 1;
		b = 2;
		c = 3;
		delta = b * b - 4 * a * c;
		if delta < 0 {
			delta = 0 - delta;
		} else {
			delta = delta;
		}
		return delta;
	}`

	// Cria o analisador léxico
	lex := &lexer.AnalisadorLexico{Texto: entrada}

	// Cria o parser
	sintatico := parser.NovoParser(lex)

	// Realiza a análise sintática
	programa, err := sintatico.Analisar()
	if err != nil {
		// Em caso de erro na análise sintática
		fmt.Println("Erro na análise sintática:", err)
		return
	}

	// Em caso de sucesso, imprime a árvore sintática

	fmt.Println("Árvore Sintática do Programa:")
	imprimirPrograma(programa)

}

/*
	testesDir := "testes"
	resultadosDir := "resultados"

	// Garante que a pasta de saída exista
	os.MkdirAll("assembly", os.ModePerm)

	err := filepath.Walk(testesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".yzy") {
			return nil
		}

		nomeArquivo := filepath.Base(path)
		nomeTeste := strings.TrimSuffix(nomeArquivo, ".yzy")
		caminhoResultado := filepath.Join(resultadosDir, nomeTeste+".txt")

		fmt.Printf("🧪 Rodando teste: %s\n", nomeTeste)

		// Ler o código-fonte
		conteudo, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Erro ao ler arquivo:", err)
			return nil
		}
		entrada := string(conteudo)

		// Análise e avaliação
		lex := &lexer.AnalisadorLexico{Texto: entrada}
		sintatico := parser.NovoParser(lex)

		programa, err := sintatico.Analisar()
		var saida string
		if err != nil {
			saida = err.Error()
			fmt.Println("📛 Assembly não gerado: erro na análise sintática")
		} else {
			resultado, err := avaliador.AvaliarPrograma(programa)
			if err != nil {
				saida = err.Error()
				fmt.Println("📛 Assembly não gerado: erro na avaliação semântica")
			} else {
				saida = fmt.Sprintf("%d", resultado)

				// geração do assembly
				codigoAssembly := geradorAssembly.GerarCodigoPrograma(programa)
				nomeSaida := nomeTeste + ".s"
				caminhoSaida := filepath.Join("assembly", nomeSaida)
				os.WriteFile(caminhoSaida, []byte(codigoAssembly), 0644)
				fmt.Println("💾 Assembly gerado em:", caminhoSaida)
			}
		}

		// Ler resultado esperado
		esperadoBytes, err := os.ReadFile(caminhoResultado)
		if err != nil {
			fmt.Println("❗ Erro ao ler resultado esperado:", err)
			return nil
		}
		esperado := strings.TrimSpace(string(esperadoBytes))

		// Comparação
		if saida == esperado {
			fmt.Println("✅ Sucesso!")
		} else {
			fmt.Println("❌ Falhou!")
			fmt.Printf("Esperado: %s\n", esperado)
			fmt.Printf("Obtido : %s\n", saida)
		}

		fmt.Println("--------------------------------------------------")
		return nil
	})

	if err != nil {
		fmt.Println("Erro geral:", err)
	}
*/
