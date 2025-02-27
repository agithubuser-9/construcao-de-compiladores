package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Definição dos lexemas aceitos
var lexemas = map[string]string{
	"ParenEsq":  "(",
	"ParenDir":  ")",
	"Soma":      "+",
	"Mult":      "*",
}

// Regex para número (sequência de dígitos)
var numeroRegex = regexp.MustCompile(`^\d+$`)

// Função para classificar um caractere ou sequência de caracteres
func classificar(lexema string) (string, bool) {
	for classif, val := range lexemas {
		if lexema == val {
			return classif, true
		}
	}
	if numeroRegex.MatchString(lexema) {
		return "Numero", true
	}
	return "", false
}

func processarArquivo(nomeArquivo string) {
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	linhaNum := 0

	for scanner.Scan() {
		linhaNum++
		linha := scanner.Text()
		// linha = strings.TrimSpace(linha) // Remove espaços extras

		lexemaAtual := ""
		posicaoInicial := 0
		var coluna_sem_0 int
		
		fmt.Printf("\nLinha %s\n", linha)
		
		for coluna_com_0, char := range linha {
			coluna_sem_0 = coluna_com_0 + 1
			if char == ' ' {
				continue // Ignora espaços em branco
			}

			strChar := string(char)

			// Se o caractere faz parte de um número, acumula no lexemaAtual
			if numeroRegex.MatchString(strChar) {
				if lexemaAtual == "" {
					posicaoInicial = coluna_sem_0
				}
				lexemaAtual += strChar
				continue
			}

			// Se estava formando um número, classifica ele antes de continuar
			if lexemaAtual != "" {
				fmt.Printf("<%s, \"%s\", %d>\n", "Numero", lexemaAtual, posicaoInicial)
				lexemaAtual = "" // Reseta para processar o próximo token
			}

			// Verifica se o caractere tem classificação
			classificacao, encontrado := classificar(strChar)
			if encontrado {
				fmt.Printf("<%s, \"%s\", %d>\n", classificacao, strChar, coluna_sem_0)
			} else {
				fmt.Printf("Erro: Caracter '%s' inesperado encontrado na linha %d, coluna %d\n", strChar, linhaNum, coluna_sem_0)
				os.Exit(1)
			}
		}

		// Se ainda houver um número acumulado ao final da linha, imprime ele
		if lexemaAtual != "" {
			fmt.Printf("<%s, \"%s\", %d>\n", "Numero", lexemaAtual, posicaoInicial)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run programa.go <arquivo.txt>")
		return
	}

	processarArquivo(os.Args[1])
}
