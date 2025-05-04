package parser

import (
	"atividade_9/lexer"
	"fmt"
	"os"
	"strconv"
)

type ParserAtt struct {
	Lexer     *lexer.AnalisadorLexico
	Atual     lexer.Token
	Variaveis map[string]int // Tabela de variáveis com seus valores
}

func (p *ParserAtt) Avancar() {
	p.Atual = p.Lexer.ProximoToken()
}

func (p *ParserAtt) Esperar(tipo lexer.TipoToken) {
	if p.Atual.Tipo == tipo {
		p.Avancar()
	} else {
		p.ErroSintatico(tipo)
	}
}

func (p *ParserAtt) ErroSintatico(esperado lexer.TipoToken) {
	fmt.Fprintf(os.Stderr, "Erro sintático: esperado token '%s', mas encontrou '%s'\n", esperado, p.Atual.Tipo)
	os.Exit(1) // interrompe o programa imediatamente
}

func (p *ParserAtt) Programa() {
	for p.Atual.Tipo == lexer.TOKEN_IDENT {
		p.Declaracao()
	}

	p.Esperar(lexer.TOKEN_ABRECHAVE)

	for p.Atual.Tipo != lexer.TOKEN_RETURN && p.Atual.Tipo != lexer.TOKEN_FECHCHAVE {
		p.Comando()
	}

	if p.Atual.Tipo == lexer.TOKEN_RETURN {
		p.ReturnComando()
	}

	p.Esperar(lexer.TOKEN_FECHCHAVE)
}

// --------------------- DECLARAÇÕES ---------------------

func (p *ParserAtt) Declaracoes() {
	for p.Atual.Tipo == lexer.TOKEN_IDENT {
		p.Declaracao()
	}
}

func (p *ParserAtt) Declaracao() {
	nomeVar := p.Atual.Valor
	p.Esperar(lexer.TOKEN_IDENT)
	p.Esperar(lexer.TOKEN_ATRIBUICAO)
	valor := p.Expressoes()
	p.Variaveis[nomeVar] = valor
	p.Esperar(lexer.TOKEN_PONTOVIRG)
}

// --------------------- COMANDOS ---------------------

func (p *ParserAtt) Comandos() {
	for p.Atual.Tipo == lexer.TOKEN_IF || p.Atual.Tipo == lexer.TOKEN_WHILE || p.Atual.Tipo == lexer.TOKEN_IDENT || p.Atual.Tipo == lexer.TOKEN_RETURN {
		p.Comando()
	}
}

func (p *ParserAtt) Comando() {
	switch p.Atual.Tipo {
	case lexer.TOKEN_IF:
		p.IfComando()
	case lexer.TOKEN_WHILE:
		p.WhileComando()
	case lexer.TOKEN_IDENT:
		p.AtribuicaoComando()
	case lexer.TOKEN_RETURN:
		p.ReturnComando()
	case lexer.TOKEN_ABRECHAVE:
		p.Esperar(lexer.TOKEN_ABRECHAVE)
		p.Comandos()
		p.Esperar(lexer.TOKEN_FECHCHAVE)
	}
}

func (p *ParserAtt) AtribuicaoComando() {
	nomeVar := p.Atual.Valor
	p.Esperar(lexer.TOKEN_IDENT)
	p.Esperar(lexer.TOKEN_ATRIBUICAO)
	valor := p.Expressoes()
	p.Variaveis[nomeVar] = valor
	p.Esperar(lexer.TOKEN_PONTOVIRG)
}

func (p *ParserAtt) ReturnComando() {

	//fmt.Println(">>> ENTROU EM EXPRESSION COM TOKEN:", p.Atual.Tipo, "VALOR:", p.Atual.Valor)

	p.Esperar(lexer.TOKEN_RETURN)
	valor := p.Expressoes()
	fmt.Println(valor)
	p.Esperar(lexer.TOKEN_PONTOVIRG)

}

func (p *ParserAtt) IfComando() {
	p.Esperar(lexer.TOKEN_IF)
	p.Esperar(lexer.TOKEN_ABREPAR)
	condicao := p.Expressoes()
	p.Esperar(lexer.TOKEN_FECHAPAR)

	if condicao != 0 {
		// Executa o bloco do 'if'
		p.Comando()
	} else {
		// Pula o bloco do 'if'
		p.PularComando()
	}

	// Verifica se tem 'else'
	if p.Atual.Tipo == lexer.TOKEN_ELSE {
		p.Esperar(lexer.TOKEN_ELSE)

		if condicao == 0 {
			// Executa o bloco do 'else' se a condição era falsa
			p.Comando()
		} else {
			// Pula o bloco do 'else' se a condição era verdadeira
			p.PularComando()
		}
	}
}

func (p *ParserAtt) PularComando() {
	if p.Atual.Tipo == lexer.TOKEN_ABRECHAVE {
		p.Esperar(lexer.TOKEN_ABRECHAVE)
		for p.Atual.Tipo != lexer.TOKEN_FECHCHAVE && p.Atual.Tipo != lexer.TOKEN_FIM {
			p.Avancar()
		}
		p.Esperar(lexer.TOKEN_FECHCHAVE)
	} else {
		p.Avancar()
	}
}

func (p *ParserAtt) WhileComando() {
	// Marca a posição do token 'while' para reiniciar o bloco completo do laço
	inicioWhile := p.Lexer.Pos - len(p.Atual.Valor)

	// Consome o 'while' e a condição
	p.Esperar(lexer.TOKEN_WHILE)
	p.Esperar(lexer.TOKEN_ABREPAR)
	cond := p.Expressoes()
	p.Esperar(lexer.TOKEN_FECHAPAR)

	// Aguarda e entra no bloco
	p.Esperar(lexer.TOKEN_ABRECHAVE)

	for cond != 0 {
		// Executa o corpo do laço
		p.Comandos()

		// Reposiciona para o início do laço para nova avaliação
		p.Lexer.Pos = inicioWhile
		p.Avancar() // sincroniza p.Atual com novo Pos

		// Reconsome toda a estrutura: while (cond) {
		p.Esperar(lexer.TOKEN_WHILE)
		p.Esperar(lexer.TOKEN_ABREPAR)
		cond = p.Expressoes()
		p.Esperar(lexer.TOKEN_FECHAPAR)
		p.Esperar(lexer.TOKEN_ABRECHAVE)
	}

	// Saiu do laço: pula o corpo restante até encontrar a chave de fechamento
	abertas := 1
	for abertas > 0 && p.Atual.Tipo != lexer.TOKEN_FIM {
		if p.Atual.Tipo == lexer.TOKEN_ABRECHAVE {
			abertas++
		} else if p.Atual.Tipo == lexer.TOKEN_FECHCHAVE {
			abertas--
		}
		p.Avancar()
	}
}

// --------------------- EXPRESSÕES ---------------------

func (p *ParserAtt) Expressoes() int {

	//fmt.Println(">>> ENTROU EM EXPRESSION COM TOKEN:", p.Atual.Tipo, "VALOR:", p.Atual.Valor)

	left := p.ExpAritmetica()

	for {
		//fmt.Println("TOKEN NO INÍCIO DO LOOP:", p.Atual.Tipo, "VALOR:", p.Atual.Valor)
		switch p.Atual.Tipo {
		case lexer.TOKEN_MENORQ, lexer.TOKEN_MAIORQ, lexer.TOKEN_IGUALDADE, lexer.TOKEN_MENORIGUAL, lexer.TOKEN_MAIORIGUAL:
			op := p.Atual.Tipo
			p.Avancar()
			right := p.ExpAritmetica()
			//fmt.Println("TOKEN ATUAL:", p.Atual.Tipo, "VALOR:", p.Atual.Valor)

			switch op {
			case lexer.TOKEN_MENORQ:
				if left < right {
					left = 1
				} else {
					left = 0
				}
			case lexer.TOKEN_MAIORQ:
				if left > right {
					left = 1
				} else {
					left = 0
				}
			case lexer.TOKEN_IGUALDADE:
				if left == right {
					left = 1
				} else {
					left = 0
				}
			case lexer.TOKEN_MENORIGUAL:
				if left <= right {
					left = 1
				} else {
					left = 0
				}
			case lexer.TOKEN_MAIORIGUAL:
				if left >= right {
					left = 1
				} else {
					left = 0
				}
			}
		default:
			return left
		}

	}
}

func (p *ParserAtt) ExpAritmetica() int {
	result := p.Termo()
	for p.Atual.Tipo == lexer.TOKEN_SOMA || p.Atual.Tipo == lexer.TOKEN_SUB {
		if p.Atual.Tipo == lexer.TOKEN_SOMA {
			p.Avancar()
			result += p.Termo()
		} else if p.Atual.Tipo == lexer.TOKEN_SUB {
			p.Avancar()
			result -= p.Termo()
		}
	}
	return result
}

func (p *ParserAtt) Termo() int {
	result := p.Fator()
	for p.Atual.Tipo == lexer.TOKEN_MULT || p.Atual.Tipo == lexer.TOKEN_DIV {
		if p.Atual.Tipo == lexer.TOKEN_MULT {
			p.Avancar()
			result *= p.Fator()
		} else if p.Atual.Tipo == lexer.TOKEN_DIV {
			p.Avancar()
			result /= p.Fator()
		}
	}
	return result
}

func (p *ParserAtt) Fator() int {
	switch p.Atual.Tipo {
	case lexer.TOKEN_INTEIRO:
		valor, _ := strconv.Atoi(p.Atual.Valor)
		p.Esperar(lexer.TOKEN_INTEIRO)
		return valor
	case lexer.TOKEN_IDENT:
		val, ok := p.Variaveis[p.Atual.Valor]
		if !ok {
			panic(fmt.Sprintf("Variável '%s' não declarada", p.Atual.Valor))
		}
		p.Esperar(lexer.TOKEN_IDENT)
		return val
	case lexer.TOKEN_ABREPAR:
		p.Esperar(lexer.TOKEN_ABREPAR)
		valor := p.Expressoes()
		p.Esperar(lexer.TOKEN_FECHAPAR)
		return valor
	default:
		p.ErroSintatico("INTEIRO | IDENT | (")
		return 0 //garantia que vai fechar se encontrat o erro
	}
}
