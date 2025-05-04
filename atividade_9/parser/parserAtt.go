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

// Programa representa um conjunto de declarações e comandos
type Programa struct {
	Declaracoes []Declaracao
	Comandos    []Comando
}

// Declaracao representa uma variável inicializada
type Declaracao struct {
	Nome string
	Expr Expressao
}

// Comando representa uma instrução (atribuição, if, while, retorno, etc.)
type Comando interface{}

// Expressao representa qualquer valor computável (número, variável, operação)
type Expressao interface{}

// ===== EXPRESSÕES =====
type Constante struct {
	Valor int
}

type Variavel struct {
	Nome string
}

type OperacaoBinaria struct {
	Operador          string
	Esquerda, Direita Expressao
}

// ===== COMANDOS =====
type Atribuicao struct {
	Nome string
	Expr Expressao
}

type Se struct {
	Condicao     Expressao
	Entao, Senao Comando // Senao pode ser nil
}

type Enquanto struct {
	Condicao Expressao
	Corpo    Comando
}

type Retorno struct {
	Expr Expressao
}

type Sequencia struct {
	Comandos []Comando
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

func (p *ParserAtt) Programa() *Programa {
	var declaracoes []Declaracao

	// Processa declarações do tipo: nome = expressao;
	for p.Atual.Tipo == lexer.TOKEN_IDENT {
		decl := p.ParseDeclaracao()
		declaracoes = append(declaracoes, decl)
	}

	p.Esperar(lexer.TOKEN_ABRECHAVE)

	var comandos []Comando
	for p.Atual.Tipo != lexer.TOKEN_FECHCHAVE && p.Atual.Tipo != lexer.TOKEN_FIM {
		comando := p.ParseComando()
		comandos = append(comandos, comando)
	}

	p.Esperar(lexer.TOKEN_FECHCHAVE)

	return &Programa{
		Declaracoes: declaracoes,
		Comandos:    comandos,
	}
}

// --------------------- DECLARAÇÕES ---------------------

func (p *ParserAtt) Declaracoes() {
	for p.Atual.Tipo == lexer.TOKEN_IDENT {
		p.ParseDeclaracao()
	}
}

func (p *ParserAtt) ParseDeclaracao() Declaracao {
	nomeVar := p.Atual.Valor
	p.Esperar(lexer.TOKEN_IDENT)
	p.Esperar(lexer.TOKEN_ATRIBUICAO)
	valor := p.ParseExpressao()
	p.Esperar(lexer.TOKEN_PONTOVIRG)

	return Declaracao{
		Nome: nomeVar,
		Expr: valor,
	}
}

// --------------------- COMANDOS ---------------------

func (p *ParserAtt) Comandos() {
	for p.Atual.Tipo == lexer.TOKEN_IF || p.Atual.Tipo == lexer.TOKEN_WHILE || p.Atual.Tipo == lexer.TOKEN_IDENT || p.Atual.Tipo == lexer.TOKEN_RETURN {
		p.ParseComando()
	}
}

func (p *ParserAtt) ParseComando() Comando {
	switch p.Atual.Tipo {
	case lexer.TOKEN_IF:
		return p.ParseIf()
	case lexer.TOKEN_WHILE:
		return p.ParseWhile()
	case lexer.TOKEN_RETURN:
		return p.ParseReturn()
	case lexer.TOKEN_IDENT:
		return p.ParseAtribuicao()
	case lexer.TOKEN_ABRECHAVE:
		return p.ParseBloco()
	default:
		p.ErroSintatico("comando válido")
		return nil
	}
}

func (p *ParserAtt) ParseAtribuicao() Comando {
	nome := p.Atual.Valor
	p.Esperar(lexer.TOKEN_IDENT)
	p.Esperar(lexer.TOKEN_ATRIBUICAO)
	expr := p.ParseExpressao()
	p.Esperar(lexer.TOKEN_PONTOVIRG)

	return &Atribuicao{
		Nome: nome,
		Expr: expr,
	}
}

func (p *ParserAtt) ParseReturn() Comando {
	p.Esperar(lexer.TOKEN_RETURN)
	expr := p.ParseExpressao()
	p.Esperar(lexer.TOKEN_PONTOVIRG)

	return &Retorno{Expr: expr}
}

func (p *ParserAtt) ParseIf() Comando {
	p.Esperar(lexer.TOKEN_IF)
	p.Esperar(lexer.TOKEN_ABREPAR)
	cond := p.ParseExpressao()
	p.Esperar(lexer.TOKEN_FECHAPAR)

	entao := p.ParseComando()

	var senao Comando = nil
	if p.Atual.Tipo == lexer.TOKEN_ELSE {
		p.Esperar(lexer.TOKEN_ELSE)
		senao = p.ParseComando()
	}

	return &Se{
		Condicao: cond,
		Entao:    entao,
		Senao:    senao,
	}
}

func (p *ParserAtt) ParseBloco() Comando {
	p.Esperar(lexer.TOKEN_ABRECHAVE)

	var comandos []Comando
	for p.Atual.Tipo != lexer.TOKEN_FECHCHAVE && p.Atual.Tipo != lexer.TOKEN_FIM {
		comandos = append(comandos, p.ParseComando())
	}

	p.Esperar(lexer.TOKEN_FECHCHAVE)

	return &Sequencia{Comandos: comandos}
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

func (p *ParserAtt) ParseWhile() Comando {
	p.Esperar(lexer.TOKEN_WHILE)
	p.Esperar(lexer.TOKEN_ABREPAR)
	cond := p.ParseExpressao()
	p.Esperar(lexer.TOKEN_FECHAPAR)

	corpo := p.ParseComando()

	return &Enquanto{
		Condicao: cond,
		Corpo:    corpo,
	}
}

// --------------------- EXPRESSÕES ---------------------

func (p *ParserAtt) ParseExpressao() Expressao {
	left := p.ParseExpAritmetica()

	for {
		switch p.Atual.Tipo {
		case lexer.TOKEN_MENORQ, lexer.TOKEN_MAIORQ, lexer.TOKEN_IGUALDADE, lexer.TOKEN_MENORIGUAL, lexer.TOKEN_MAIORIGUAL:
			op := p.Atual.Valor
			p.Avancar()
			right := p.ParseExpAritmetica()
			left = &OperacaoBinaria{
				Operador: op,
				Esquerda: left,
				Direita:  right,
			}
		default:
			return left
		}
	}
}

func (p *ParserAtt) ParseExpAritmetica() Expressao {
	expr := p.ParseTermo()

	for p.Atual.Tipo == lexer.TOKEN_SOMA || p.Atual.Tipo == lexer.TOKEN_SUB {
		op := p.Atual.Valor
		p.Avancar()
		direita := p.ParseTermo()

		expr = &OperacaoBinaria{
			Operador: op,
			Esquerda: expr,
			Direita:  direita,
		}
	}

	return expr
}

func (p *ParserAtt) ParseTermo() Expressao {
	expr := p.ParseFator()

	for p.Atual.Tipo == lexer.TOKEN_MULT || p.Atual.Tipo == lexer.TOKEN_DIV {
		op := p.Atual.Valor
		p.Avancar()
		direita := p.ParseFator()

		expr = &OperacaoBinaria{
			Operador: op,
			Esquerda: expr,
			Direita:  direita,
		}
	}

	return expr
}

func (p *ParserAtt) ParseFator() Expressao {
	switch p.Atual.Tipo {
	case lexer.TOKEN_INTEIRO:
		valor, _ := strconv.Atoi(p.Atual.Valor)
		p.Esperar(lexer.TOKEN_INTEIRO)
		return &Constante{Valor: valor}

	case lexer.TOKEN_IDENT:
		nome := p.Atual.Valor
		p.Esperar(lexer.TOKEN_IDENT)
		return &Variavel{Nome: nome}

	case lexer.TOKEN_ABREPAR:
		p.Esperar(lexer.TOKEN_ABREPAR)
		expr := p.ParseExpressao()
		p.Esperar(lexer.TOKEN_FECHAPAR)
		return expr

	default:
		p.ErroSintatico("INTEIRO | IDENT | (")
		return nil
	}
}
