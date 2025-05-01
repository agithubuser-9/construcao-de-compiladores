package parser

import (
	"atividade_9/lexer"
	"fmt"
	"strconv"
)

type Expressao interface{}

type Variavel struct {
	Nome string
}

type Declaracao struct {
	Nome string
	Expr Expressao
}

type Programa struct {
	Declaracoes []Declaracao
	ExprFinal   Expressao
}

type Constante struct {
	Valor int
}

type OperacaoBinaria struct {
	Operador string
	Esquerda Expressao
	Direita  Expressao
}

type AnalisadorSintatico struct {
	Lex        *lexer.AnalisadorLexico
	TokenAtual lexer.Token
}

// inicializa o parser mantendo o token atual como referência,
// ajustado pra pervinir o avanço automatico. apenas verifica qual o proximo token
// sem avançar de fato, comportamento de "olharProximoToken" visto em sala
func NovoParser(lex *lexer.AnalisadorLexico) *AnalisadorSintatico {
	p := &AnalisadorSintatico{Lex: lex}
	p.TokenAtual = p.Lex.ProximoToken()
	return p
}

// função para consumir o token atual só se ele corresponder ao tipo esperado.
// caso contrário, retorna um erro. IMPORTANTE garantir que a análise siga a ordem correta
func (p *AnalisadorSintatico) avancaToken(tipoEsperado lexer.TipoToken) error {

	if p.TokenAtual.Tipo == lexer.TOKEN_INVALIDO {
		return fmt.Errorf("erro léxico: caractere inválido '%s'", p.TokenAtual.Valor)
	}

	if p.TokenAtual.Tipo == tipoEsperado {
		p.TokenAtual = p.Lex.ProximoToken()
		if p.TokenAtual.Tipo == lexer.TOKEN_INVALIDO {
			return fmt.Errorf("erro léxico: caractere inválido '%s'", p.TokenAtual.Valor)
		}
		return nil
	}
	return fmt.Errorf("esperado token %s, mas encontrou %s", tipoEsperado, p.TokenAtual.Tipo)
}

// <prim> ::= <num> | '(' <exp_a> ')'
func (p *AnalisadorSintatico) analisarPrimaria(tabela map[string]bool) (Expressao, error) {
	tok := p.TokenAtual
	switch tok.Tipo {
	case lexer.TOKEN_INTEIRO:
		valor, _ := strconv.Atoi(tok.Valor)
		p.avancaToken(lexer.TOKEN_INTEIRO)
		return &Constante{Valor: valor}, nil

	case lexer.TOKEN_IDENT:
		if _, ok := tabela[tok.Valor]; !ok {
			return nil, fmt.Errorf("variável '%s' não declarada", tok.Valor)
		}
		p.avancaToken(lexer.TOKEN_IDENT)
		return &Variavel{Nome: tok.Valor}, nil

	case lexer.TOKEN_ABREPAR:
		p.avancaToken(lexer.TOKEN_ABREPAR)
		expr, err := p.analisarExpA(tabela)
		if err != nil {
			return nil, err
		}
		if p.TokenAtual.Tipo != lexer.TOKEN_FECHAPAR {
			return nil, fmt.Errorf("esperado ')', mas encontrou '%s'", p.TokenAtual.Valor)
		}
		p.avancaToken(lexer.TOKEN_FECHAPAR)
		return expr, nil

	default:
		return nil, fmt.Errorf("token inesperado: %s", tok.Valor)
	}
}

// <exp_m> ::= <prim> (( '*' | '/' ) <prim>)*
func (p *AnalisadorSintatico) analisarExpM(tabela map[string]bool) (Expressao, error) {
	esq, err := p.analisarPrimaria(tabela)
	if err != nil {
		return nil, err
	}
	for p.TokenAtual.Tipo == lexer.TOKEN_MULT || p.TokenAtual.Tipo == lexer.TOKEN_DIV {
		op := p.TokenAtual.Valor
		p.avancaToken(p.TokenAtual.Tipo)
		dir, err := p.analisarPrimaria(tabela)
		if err != nil {
			return nil, err
		}
		esq = &OperacaoBinaria{Operador: op, Esquerda: esq, Direita: dir}
	}
	return esq, nil
}

// <exp_a> ::= <exp_m> (( '+' | '-' ) <exp_m>)*
func (p *AnalisadorSintatico) analisarExpA(tabela map[string]bool) (Expressao, error) {
	esq, err := p.analisarExpM(tabela)
	if err != nil {
		return nil, err
	}
	for p.TokenAtual.Tipo == lexer.TOKEN_SOMA || p.TokenAtual.Tipo == lexer.TOKEN_SUB {
		op := p.TokenAtual.Valor
		p.avancaToken(p.TokenAtual.Tipo)
		dir, err := p.analisarExpM(tabela)
		if err != nil {
			return nil, err
		}
		esq = &OperacaoBinaria{Operador: op, Esquerda: esq, Direita: dir}
	}
	return esq, nil
}

// Inicia a análise sintática a partir de exp_a.
func (p *AnalisadorSintatico) Analisar() (*Programa, error) {
	var declaracoes []Declaracao
	tabelaSimbolos := make(map[string]bool)

	for p.TokenAtual.Tipo == lexer.TOKEN_IDENT {
		decl, err := p.analisarDeclaracao(tabelaSimbolos)
		if err != nil {
			return nil, err
		}
		declaracoes = append(declaracoes, *decl)
	}

	if p.TokenAtual.Tipo != lexer.TOKEN_RETURN {
		return nil, fmt.Errorf("esperado 'return' para iniciar a expressão final, mas encontrou '%s'", p.TokenAtual.Valor)
	}
	p.avancaToken(lexer.TOKEN_RETURN)

	exprFinal, err := p.analisarExpA(tabelaSimbolos)
	if err != nil {
		return nil, err
	}

	return &Programa{
		Declaracoes: declaracoes,
		ExprFinal:   exprFinal,
	}, nil
}
func (p *AnalisadorSintatico) analisarDeclaracao(tabela map[string]bool) (*Declaracao, error) {
	if p.TokenAtual.Tipo != lexer.TOKEN_IDENT {
		return nil, fmt.Errorf("esperado identificador, mas encontrou '%s'", p.TokenAtual.Valor)
	}
	nomeVar := p.TokenAtual.Valor

	if _, existe := tabela[nomeVar]; existe {
		return nil, fmt.Errorf("variável '%s' já foi declarada", nomeVar)
	}

	tabela[nomeVar] = true

	p.avancaToken(lexer.TOKEN_IDENT)

	if err := p.avancaToken(lexer.TOKEN_IGUAL); err != nil {
		return nil, err
	}

	expr, err := p.analisarExpA(tabela)
	if err != nil {
		return nil, err
	}

	if err := p.avancaToken(lexer.TOKEN_PONTOVIRG); err != nil {
		return nil, err
	}

	return &Declaracao{Nome: nomeVar, Expr: expr}, nil
}
