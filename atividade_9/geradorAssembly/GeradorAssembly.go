package geradorAssembly

import (
	"atividade_7_Revisada/parser"
	"fmt"
	"strings"
)

func GerarCodigoPrograma(prog *parser.Programa) string {
	var sb strings.Builder

	//declaração de variáveis
	sb.WriteString(".section .bss\n")
	declared := make(map[string]bool)
	for _, decl := range prog.Declaracoes {
		if !declared[decl.Nome] {
			sb.WriteString(fmt.Sprintf(".lcomm %s, 8\n", decl.Nome))
			declared[decl.Nome] = true
		}
	}
	sb.WriteString("\n.section .text\n.globl _start\n_start:\n")

	for _, decl := range prog.Declaracoes {
		codigo := gerarCodigoExpr(decl.Expr)
		sb.WriteString(codigo)
		sb.WriteString(fmt.Sprintf("mov %%rax, %s\n\n", decl.Nome))
	}

	// Print Extressao
	sb.WriteString("# return\n")
	sb.WriteString(gerarCodigoExpr(prog.ExprFinal))
	sb.WriteString("call imprime_num\n")
	sb.WriteString("call sair\n")
	sb.WriteString("\n.include \"runtime.s\"\n")

	return sb.String()
}
func gerarCodigoExpr(expr parser.Expressao) string {
	switch e := expr.(type) {
	case *parser.Constante:
		return fmt.Sprintf("mov $%d, %%rax\n", e.Valor)

	case *parser.Variavel:
		return fmt.Sprintf("mov %s, %%rax\n", e.Nome)

	case *parser.OperacaoBinaria:
		esq := gerarCodigoExpr(e.Esquerda)
		dir := gerarCodigoExpr(e.Direita)

		var operacao string
		switch e.Operador {
		case "+":
			operacao = "add %rbx, %rax"
		case "-":
			operacao = "sub %rbx, %rax"
		case "*":
			operacao = "mul %rbx"
		case "/":
			operacao = "mov $0, %rdx\nidiv %rbx"
		}

		return esq +
			"push %rax\n" +
			dir +
			"pop %rbx\n" +
			operacao + "\n"

	default:
		return "# erro: tipo de expressão desconhecido\n"
	}
}
