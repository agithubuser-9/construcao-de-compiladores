// parser.go
package parser

import (
	"atividade_8/lexer"
	"fmt"
	"strconv"
)

// Expressao representa um nó genérico da AST.
type Expressao interface{}

// Nós da AST.
type Constante struct {
	Valor int
}

type Variavel struct {
	Nome string
}

type OperacaoBinaria struct {
	Operador string
	Esquerda Expressao
	Direita  Expressao
}

type Declaracao struct {
	Nome      string
	Expressao Expressao
}

type Programa struct {
	Declaracoes []Declaracao
	Expressao   Expressao
}

// AnalisadorSintatico – o parser.
type AnalisadorSintatico struct {
	Lex            *lexer.AnalisadorLexico
	TokenAtual     lexer.Token
	TabelaSimbolos map[string]bool // Variáveis declaradas.
}

func NovoParser(lex *lexer.AnalisadorLexico) *AnalisadorSintatico {
	p := &AnalisadorSintatico{
		Lex:            lex,
		TabelaSimbolos: make(map[string]bool),
	}
	p.TokenAtual = p.Lex.ProximoToken()
	return p
}

func (p *AnalisadorSintatico) avancaToken(tipoEsperado lexer.TipoToken) error {
	if p.TokenAtual.Tipo == tipoEsperado {
		p.TokenAtual = p.Lex.ProximoToken()
		return nil
	}
	return fmt.Errorf("esperado token %s, mas encontrou %s", tipoEsperado, p.TokenAtual.Tipo)
}

// <prim> ::= <num> | <var> | '(' <exp_a> ')'
func (p *AnalisadorSintatico) analisarPrimaria() (Expressao, error) {
	tok := p.TokenAtual
	if tok.Tipo == lexer.TOKEN_INTEIRO {
		valor, err := strconv.Atoi(tok.Valor)
		if err != nil {
			return nil, err
		}
		p.avancaToken(lexer.TOKEN_INTEIRO)
		return &Constante{Valor: valor}, nil
	} else if tok.Tipo == lexer.TOKEN_IDENT {
		nome := tok.Valor
		p.avancaToken(lexer.TOKEN_IDENT)
		// Verifica se a variável foi declarada.
		if _, ok := p.TabelaSimbolos[nome]; !ok {
			return nil, fmt.Errorf("variável '%s' não declarada", nome)
		}
		return &Variavel{Nome: nome}, nil
	} else if tok.Tipo == lexer.TOKEN_ABREPAR {
		p.avancaToken(lexer.TOKEN_ABREPAR)
		expr, err := p.analisarExpA()
		if err != nil {
			return nil, err
		}
		if err := p.avancaToken(lexer.TOKEN_FECHAPAR); err != nil {
			return nil, err
		}
		return expr, nil
	}
	return nil, fmt.Errorf("token inesperado: %s", tok.Valor)
}

// <exp_m> ::= <prim> { ( '*' | '/' ) <prim> }*
func (p *AnalisadorSintatico) analisarExpM() (Expressao, error) {
	esq, err := p.analisarPrimaria()
	if err != nil {
		return nil, err
	}
	for p.TokenAtual.Tipo == lexer.TOKEN_MULT || p.TokenAtual.Tipo == lexer.TOKEN_DIV {
		op := p.TokenAtual.Valor
		p.avancaToken(p.TokenAtual.Tipo)
		dir, err := p.analisarPrimaria()
		if err != nil {
			return nil, err
		}
		esq = &OperacaoBinaria{Operador: op, Esquerda: esq, Direita: dir}
	}
	return esq, nil
}

// <exp_a> ::= <exp_m> { ( '+' | '-' ) <exp_m> }*
func (p *AnalisadorSintatico) analisarExpA() (Expressao, error) {
	esq, err := p.analisarExpM()
	if err != nil {
		return nil, err
	}
	for p.TokenAtual.Tipo == lexer.TOKEN_SOMA || p.TokenAtual.Tipo == lexer.TOKEN_SUB {
		op := p.TokenAtual.Valor
		p.avancaToken(p.TokenAtual.Tipo)
		dir, err := p.analisarExpM()
		if err != nil {
			return nil, err
		}
		esq = &OperacaoBinaria{Operador: op, Esquerda: esq, Direita: dir}
	}
	return esq, nil
}

// <decl> ::= <var> '=' <exp_a> ';'
func (p *AnalisadorSintatico) analisarDeclaracao() (Declaracao, error) {
	if p.TokenAtual.Tipo != lexer.TOKEN_IDENT {
		return Declaracao{}, fmt.Errorf("declaração deve iniciar com identificador, mas encontrou %s", p.TokenAtual.Valor)
	}
	nome := p.TokenAtual.Valor
	p.avancaToken(lexer.TOKEN_IDENT)
	if err := p.avancaToken(lexer.TOKEN_IGUAL); err != nil {
		return Declaracao{}, err
	}
	expr, err := p.analisarExpA()
	if err != nil {
		return Declaracao{}, err
	}
	if err := p.avancaToken(lexer.TOKEN_PONTOEVIRGULA); err != nil {
		return Declaracao{}, err
	}
	// Registra a variável.
	p.TabelaSimbolos[nome] = true
	return Declaracao{Nome: nome, Expressao: expr}, nil
}

// <programa> ::= { <decl> } '=' <exp_a>
func (p *AnalisadorSintatico) AnalisarPrograma() (*Programa, error) {
	var decls []Declaracao
	// Processa declarações enquanto houver identificador.
	for p.TokenAtual.Tipo == lexer.TOKEN_IDENT {
		decl, err := p.analisarDeclaracao()
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)
	}
	// A expressão final deve iniciar com '='.
	if p.TokenAtual.Tipo != lexer.TOKEN_IGUAL {
		return nil, fmt.Errorf("esperado '=' iniciando a expressão final, mas encontrou %s", p.TokenAtual.Valor)
	}
	p.avancaToken(lexer.TOKEN_IGUAL)
	expr, err := p.analisarExpA()
	if err != nil {
		return nil, err
	}
	return &Programa{Declaracoes: decls, Expressao: expr}, nil
}

