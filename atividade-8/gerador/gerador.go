// gerador.go
package gerador

import (
	"atividade_8/parser"
	"fmt"
	"strings"
)

// GerarCodigo gera o código assembly do programa EV.
func GerarCodigo(prog *parser.Programa) string {
	var bss strings.Builder
	var text strings.Builder

	// Seção BSS: declara as variáveis.
	bss.WriteString(".section .bss\n")
	for _, decl := range prog.Declaracoes {
		bss.WriteString(fmt.Sprintf(".lcomm %s, 8\n", decl.Nome))
	}

	// Seção TEXT.
	text.WriteString(".section .text\n")
	text.WriteString(".globl _start\n")
	text.WriteString("_start:\n")

	// Gera código para as declarações.
	for _, decl := range prog.Declaracoes {
		text.WriteString(generateExpr(decl.Expressao))
		text.WriteString(fmt.Sprintf("mov %%rax, %s\n", decl.Nome))
	}

	// Código para a expressão final.
	text.WriteString("; Expressão final\n")
	text.WriteString(generateExpr(prog.Expressao))
	text.WriteString("call imprime_num\n")
	text.WriteString("call sair\n")

	// Combina as seções.
	return bss.String() + "\n" + text.String()
}

func generateExpr(expr parser.Expressao) string {
	switch e := expr.(type) {
	case *parser.Constante:
		return fmt.Sprintf("mov $%d, %%rax\n", e.Valor)
	case *parser.Variavel:
		return fmt.Sprintf("mov %s, %%rax\n", e.Nome)
	case *parser.OperacaoBinaria:
		code := generateExpr(e.Esquerda)
		code += "push %%rax\n"
		code += generateExpr(e.Direita)
		code += "pop %%rbx\n"
		switch e.Operador {
		case "+":
			code += "add %%rbx, %%rax\n"
		case "-":
			code += "sub %%rax, %%rbx\n"
			code += "mov %%rbx, %%rax\n"
		case "*":
			code += "imul %%rbx, %%rax\n"
		case "/":
			code += "cqto\n"
			code += "idiv %%rbx\n"
		}
		return code
	default:
		return ""
	}
}

