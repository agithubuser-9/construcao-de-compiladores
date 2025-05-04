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
		case ">":
			if esq > dir {
				return 1, nil
			}
			return 0, nil
		case "<":
			if esq < dir {
				return 1, nil
			}
			return 0, nil
		case "==":
			if esq == dir {
				return 1, nil
			}
			return 0, nil
		case ">=":
			if esq >= dir {
				return 1, nil
			}
			return 0, nil
		case "<=":
			if esq <= dir {
				return 1, nil
			}
			return 0, nil
		}

	}
	return 0, fmt.Errorf("expressão inválida")
}

func AvaliarPrograma(prog *parser.Programa) (int, error) {
	ctx := make(map[string]int)

	// Avaliar declarações
	for _, decl := range prog.Declaracoes {
		val, err := Avaliar(decl.Expr, ctx)
		if err != nil {
			return 0, fmt.Errorf("erro ao avaliar '%s': %v", decl.Nome, err)
		}
		ctx[decl.Nome] = val
	}

	// Executar comandos
	for _, cmd := range prog.Comandos {
		val, retorno, err := ExecutarComando(cmd, ctx)
		if err != nil {
			return 0, err
		}
		if retorno {
			return val, nil
		}
	}

	return 0, fmt.Errorf("programa não possui comando return")
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

func ExecutarComando(cmd parser.Comando, ctx map[string]int) (int, bool, error) {
	switch c := cmd.(type) {

	case *parser.Atribuicao:
		val, err := Avaliar(c.Expr, ctx)
		if err != nil {
			return 0, false, err
		}
		ctx[c.Nome] = val
		return 0, false, nil

	case *parser.Se:
		cond, err := Avaliar(c.Condicao, ctx)
		if err != nil {
			return 0, false, err
		}
		if cond != 0 {
			return ExecutarComando(c.Entao, ctx)
		} else if c.Senao != nil {
			return ExecutarComando(c.Senao, ctx)
		}
		return 0, false, nil

	case *parser.Enquanto:
		for {
			cond, err := Avaliar(c.Condicao, ctx)
			if err != nil {
				return 0, false, err
			}
			if cond == 0 {
				break
			}
			val, retorno, err := ExecutarComando(c.Corpo, ctx)
			if err != nil {
				return 0, false, err
			}
			if retorno {
				return val, true, nil
			}
		}
		return 0, false, nil

	case *parser.Sequencia:
		for _, sub := range c.Comandos {
			val, retorno, err := ExecutarComando(sub, ctx)
			if err != nil {
				return 0, false, err
			}
			if retorno {
				return val, true, nil
			}
		}
		return 0, false, nil

	case *parser.Retorno:
		val, err := Avaliar(c.Expr, ctx)
		if err != nil {
			return 0, false, err
		}
		return val, true, nil

	default:
		return 0, false, fmt.Errorf("comando desconhecido")
	}
}
