package geradorAssembly

import (
	"atividade_9/parser"
	"fmt"
	"strings"
)

// GL
var rotuloCount int = 0

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
	for _, cmd := range prog.Comandos {
		sb.WriteString(gerarCodigoComando(cmd))
	}

	//sb.WriteString("\n.include \"runtime.s\"\n")

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

		var sb strings.Builder
		sb.WriteString(esq)
		sb.WriteString("push %rax\n") // salva lado esquerdo
		sb.WriteString(dir)
		sb.WriteString("pop %rbx\n")       // lado esquerdo → %rbx
		sb.WriteString("cmp %rax, %rbx\n") // compara esq (em %rbx) com dir (em %rax)

		switch e.Operador {
		case "+":
			sb.WriteString("add %rax, %rbx\nmov %rbx, %rax\n")
		case "-":
			sb.WriteString("sub %rax, %rbx\nmov %rbx, %rax\n")
		case "*":
			sb.WriteString("imul %rax, %rbx\nmov %rbx, %rax\n")
		case "/":
			sb.WriteString("mov %rbx, %rax\ncqo\nidiv %rax\n") // opcional: ajustar se necessário

		case ">":
			sb.WriteString("setg %al\nmovzx %al, %rax\n")
		case "<":
			sb.WriteString("setl %al\nmovzx %al, %rax\n")
		case "==":
			sb.WriteString("sete %al\nmovzx %al, %rax\n")
		case ">=":
			sb.WriteString("setge %al\nmovzx %al, %rax\n")
		case "<=":
			sb.WriteString("setle %al\nmovzx %al, %rax\n")
		default:
			sb.WriteString("# operador desconhecido\n")
		}

		return sb.String()

	default:
		return "# erro: tipo de expressão desconhecido\n"
	}
}

func gerarCodigoComando(cmd parser.Comando) string {
	switch c := cmd.(type) {

	case *parser.Atribuicao:
		codigo := gerarCodigoExpr(c.Expr)
		return codigo + fmt.Sprintf("mov %%rax, %s\n", c.Nome)

	case *parser.Retorno:
		codigo := gerarCodigoExpr(c.Expr)
		return codigo + "call imprime_num\ncall sair\n"

	case *parser.Sequencia:
		var sb strings.Builder
		for _, sub := range c.Comandos {
			sb.WriteString(gerarCodigoComando(sub))
		}
		return sb.String()

	case *parser.Se:
		id := rotuloCount
		rotuloCount++

		codigo := gerarCodigoExpr(c.Condicao)
		codigo += "cmp $0, %rax\n"
		codigo += fmt.Sprintf("je else_%d\n", id)
		codigo += gerarCodigoComando(c.Entao)
		codigo += fmt.Sprintf("jmp fim_if_%d\n", id)
		codigo += fmt.Sprintf("else_%d:\n", id)
		if c.Senao != nil {
			codigo += gerarCodigoComando(c.Senao)
		}
		codigo += fmt.Sprintf("fim_if_%d:\n", id)
		return codigo

	case *parser.Enquanto:
		id := rotuloCount
		rotuloCount++

		codigo := fmt.Sprintf("inicio_while_%d:\n", id)
		codigo += gerarCodigoExpr(c.Condicao)
		codigo += "cmp $0, %rax\n"
		codigo += fmt.Sprintf("je fim_while_%d\n", id)
		codigo += gerarCodigoComando(c.Corpo)
		codigo += fmt.Sprintf("jmp inicio_while_%d\n", id)
		codigo += fmt.Sprintf("fim_while_%d:\n", id)
		return codigo

	default:
		return "# comando desconhecido\n"
	}
}
