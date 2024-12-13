package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

/*
Alunos:
Cássio Andrêzza de Almeida
Yuri Gabriel Da Silva
*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Por favor, forneça um arquivo de texto como argumento.")
		os.Exit(1)
	}

	// Verifica se o arquivo fornecido é um arquivo de texto
	filename := os.Args[1]

	// Lê o conteúdo do arquivo
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		os.Exit(1)
	}

	digitsString := string(content)

	// Remove quebras de linha no final do conteúdo
	digitsString = strings.TrimRight(digitsString, "\n")

	// Verifica se o conteúdo é composto apenas por dígitos
	if !isAllDigits(digitsString) {
		fmt.Println("O conteúdo do arquivo não é composto apenas por dígitos.")
		os.Exit(1)
	}

	// Declaração do slice de strings com instruções assembly
	assemblyGerado := []string{
		"#",
		"# modelo de saida para o compilador",
		"#",
		".section .text",
		".globl _start",
		"_start:",
		"## pega string e coloca num registrador",
		"call imprime_num",
		"call sair",
		".include \"runtime.s\"",
	}

	// Adiciona a string de dígitos ao slice de assembly
	assemblyGerado[6] = "mov $" + digitsString + ", %rax"

	// Exibe o conteúdo do slice de strings para verificação
	for _, line := range assemblyGerado {
		fmt.Println(line)
	}

	// Cria o arquivo "modelo.s"
	file, err := os.Create("modelo.s")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()

	// Escreve as strings do slice no arquivo
	writer := bufio.NewWriter(file)
	for _, line := range assemblyGerado {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}
	}
	// Garante que os dados sejam escritos no arquivo
	writer.Flush()

	fmt.Println("Arquivo 'modelo.s' criado com sucesso!")
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
