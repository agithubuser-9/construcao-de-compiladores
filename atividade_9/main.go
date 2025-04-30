package main

import (
	"atividade_7_Revisada/avaliador"
	"atividade_7_Revisada/geradorAssembly" // certifique-se de que o nome está certo
	"atividade_7_Revisada/lexer"
	"atividade_7_Revisada/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Alunos:
// Yuri Gabriel da Silva Alves
// Cássio Andrêzza de Almeida

func main() {
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
}
