package parser

import (
	"atividade_7_Revisada/lexer"
	"fmt"
	"strconv"
)

type Expressao interface{}

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
	if p.TokenAtual.Tipo == tipoEsperado {
		p.TokenAtual = p.Lex.ProximoToken()
		return nil
	}
	return fmt.Errorf("esperado token %s, mas encontrou %s", tipoEsperado, p.TokenAtual.Tipo)
}

// <prim> ::= <num> | '(' <exp_a> ')'
func (p *AnalisadorSintatico) analisarPrimaria() (Expressao, error) {
	tok := p.TokenAtual
	if tok.Tipo == lexer.TOKEN_INTEIRO {
		valor, _ := strconv.Atoi(tok.Valor)
		p.avancaToken(lexer.TOKEN_INTEIRO)
		return &Constante{Valor: valor}, nil
	} else if tok.Tipo == lexer.TOKEN_ABREPAR {
		p.avancaToken(lexer.TOKEN_ABREPAR)
		expr, err := p.analisarExpA()
		if err != nil {
			return nil, err
		}
		p.avancaToken(lexer.TOKEN_FECHAPAR)
		return expr, nil
	}
	return nil, fmt.Errorf("token inesperado: %s", tok.Valor)
}

// <exp_m> ::= <prim> (( '*' | '/' ) <prim>)*
func (p *AnalisadorSintatico) analisarExpM() (Expressao, error) {
	esq, _ := p.analisarPrimaria()
	for p.TokenAtual.Tipo == lexer.TOKEN_MULT || p.TokenAtual.Tipo == lexer.TOKEN_DIV {
		op := p.TokenAtual.Valor
		p.avancaToken(p.TokenAtual.Tipo)
		dir, _ := p.analisarPrimaria()
		esq = &OperacaoBinaria{Operador: op, Esquerda: esq, Direita: dir}
	}
	return esq, nil
}

// <exp_a> ::= <exp_m> (( '+' | '-' ) <exp_m>)*
func (p *AnalisadorSintatico) analisarExpA() (Expressao, error) {
	esq, _ := p.analisarExpM()
	for p.TokenAtual.Tipo == lexer.TOKEN_SOMA || p.TokenAtual.Tipo == lexer.TOKEN_SUB {
		op := p.TokenAtual.Valor
		p.avancaToken(p.TokenAtual.Tipo)
		dir, _ := p.analisarExpM()
		esq = &OperacaoBinaria{Operador: op, Esquerda: esq, Direita: dir}
	}
	return esq, nil
}

// Inicia a análise sintática a partir de exp_a.
func (p *AnalisadorSintatico) Analisar() (Expressao, error) {
	return p.analisarExpA()
}
