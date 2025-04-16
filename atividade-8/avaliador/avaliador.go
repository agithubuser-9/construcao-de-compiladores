// avaliador.go
package avaliador

import (
	"atividade_8/parser"
	"fmt"
)

// AvaliarPrograma processa declarações e avalia a expressão final, retornando o resultado.
func AvaliarPrograma(prog *parser.Programa) (int, error) {
	env := make(map[string]int)
	// Executa as declarações.
	for _, decl := range prog.Declaracoes {
		val, err := avaliarExpressao(decl.Expressao, env)
		if err != nil {
			return 0, err
		}
		env[decl.Nome] = val
	}
	// Avalia a expressão final.
	return avaliarExpressao(prog.Expressao, env)
}

func avaliarExpressao(expr parser.Expressao, env map[string]int) (int, error) {
	switch e := expr.(type) {
	case *parser.Constante:
		return e.Valor, nil
	case *parser.OperacaoBinaria:
		esq, err := avaliarExpressao(e.Esquerda, env)
		if err != nil {
			return 0, err
		}
		dir, err := avaliarExpressao(e.Direita, env)
		if err != nil {
			return 0, err
		}
		switch e.Operador {
		case "+":
			return esq + dir, nil
		case "-":
			return esq - dir, nil
		case "*":
			return esq * dir, nil
		case "/":
			if dir == 0 {
				return 0, fmt.Errorf("divisão por zero")
			}
			return esq / dir, nil
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", e.Operador)
		}
	case *parser.Variavel:
		if val, ok := env[e.Nome]; ok {
			return val, nil
		}
		return 0, fmt.Errorf("variável '%s' não declarada", e.Nome)
	}
	return 0, fmt.Errorf("expressão inválida")
}

// Imprimir converte a AST em uma forma parenthesizada.
func Imprimir(expr parser.Expressao) string {
	switch e := expr.(type) {
	case *parser.Constante:
		return fmt.Sprintf("%d", e.Valor)
	case *parser.Variavel:
		return e.Nome
	case *parser.OperacaoBinaria:
		return fmt.Sprintf("(%s %s %s)", Imprimir(e.Esquerda), e.Operador, Imprimir(e.Direita))
	}
	return ""
}

