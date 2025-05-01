package avaliador

import (
	"atividade_9/parser"
	"fmt"
)

func Avaliar(expr parser.Expressao, ctx map[string]int) (int, error) {
	switch e := expr.(type) {

	case *parser.Constante:
		return e.Valor, nil

	case *parser.Variavel:
		val, ok := ctx[e.Nome]
		if !ok {
			return 0, fmt.Errorf("variável '%s' não inicializada", e.Nome)
		}
		return val, nil

	case *parser.OperacaoBinaria:
		esq, err := Avaliar(e.Esquerda, ctx)
		if err != nil {
			return 0, err
		}
		dir, err := Avaliar(e.Direita, ctx)
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
		}
	}
	return 0, fmt.Errorf("expressão inválida")
}

func AvaliarPrograma(prog *parser.Programa) (int, error) {
	ctx := make(map[string]int)

	for _, decl := range prog.Declaracoes {
		val, err := Avaliar(decl.Expr, ctx)
		if err != nil {
			return 0, fmt.Errorf("erro ao avaliar '%s': %v", decl.Nome, err)
		}
		ctx[decl.Nome] = val
	}

	// Avaliar a expressão final
	return Avaliar(prog.ExprFinal, ctx)
}

// Função para imprimir a árvore sintática (em formato "parenthesizado" ).
// DEPRECATED após adição da linguagem EV
func Imprimir(expr parser.Expressao) string {
	switch e := expr.(type) {
	case *parser.Constante:
		return fmt.Sprintf("%d", e.Valor)
	case *parser.Variavel:
		return e.Nome
	case *parser.OperacaoBinaria:
		return fmt.Sprintf("(%s %s %s)", Imprimir(e.Esquerda), e.Operador, Imprimir(e.Direita))
	default:
		return ""
	}
}
