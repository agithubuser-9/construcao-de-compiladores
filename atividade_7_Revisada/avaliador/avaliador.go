package avaliador

import (
	"atividade_7_Revisada/parser"
	"fmt"
)

// Função que avalia a árvore sintática, verifica se a expressao está valida.
func Avaliar(expr parser.Expressao) (int, error) {
	switch e := expr.(type) {
	case *parser.Constante:
		return e.Valor, nil
	case *parser.OperacaoBinaria:
		esq, _ := Avaliar(e.Esquerda)
		dir, _ := Avaliar(e.Direita)
		switch e.Operador {
		case "+":
			return esq + dir, nil
		case "-":
			return esq - dir, nil
		case "*":
			return esq * dir, nil
		case "/":
			//tratamento pra caso a expressão tenha uma divisao por zero
			if dir == 0 {
				return 0, fmt.Errorf("divisão por zero")
			}
			return esq / dir, nil
		}
	}
	return 0, fmt.Errorf("expressão inválida")
}

// Função para imprimir a árvore sintática (em formato "parenthesizado" ).
func Imprimir(expr parser.Expressao) string {
	switch e := expr.(type) {
	case *parser.Constante:
		return fmt.Sprintf("%d", e.Valor)
	case *parser.OperacaoBinaria:
		return fmt.Sprintf("(%s %s %s)", Imprimir(e.Esquerda), e.Operador, Imprimir(e.Direita))
	default:
		return ""
	}
}
