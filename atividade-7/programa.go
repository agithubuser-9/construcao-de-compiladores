package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// Definição dos tipos de token.
type TipoToken string

const (
	TOKEN_INTEIRO  TipoToken = "INT"
	TOKEN_SOMA     TipoToken = "+"
	TOKEN_SUB      TipoToken = "-"
	TOKEN_MULT     TipoToken = "*"
	TOKEN_DIV      TipoToken = "/"
	TOKEN_ABREPAR  TipoToken = "("
	TOKEN_FECHAPAR TipoToken = ")"
	TOKEN_FIM      TipoToken = "EOF"
	TOKEN_INVALIDO TipoToken = "ILLEGAL"
)

// Estrutura dos tokens.
type Token struct {
	Tipo  TipoToken
	Valor string
}

// Lexer para análise léxica.
type AnalisadorLexico struct {
	texto string
	pos   int
}

func (a *AnalisadorLexico) proximoToken() Token {
	//Remove os espaços da empressão de entrada (foi mantido pra centralizar a 'limpeza' da entrada no Lexer)
	for a.pos < len(a.texto) && unicode.IsSpace(rune(a.texto[a.pos])) {
		a.pos++
	}
	if a.pos >= len(a.texto) {
		return Token{Tipo: TOKEN_FIM, Valor: ""}
	}
	caractere := a.texto[a.pos]
	if unicode.IsDigit(rune(caractere)) {
		inicio := a.pos
		for a.pos < len(a.texto) && unicode.IsDigit(rune(a.texto[a.pos])) {
			a.pos++
		}
		return Token{Tipo: TOKEN_INTEIRO, Valor: a.texto[inicio:a.pos]}
	}
	switch caractere {
	case '+':
		a.pos++
		return Token{Tipo: TOKEN_SOMA, Valor: string(caractere)}
	case '-':
		a.pos++
		return Token{Tipo: TOKEN_SUB, Valor: string(caractere)}
	case '*':
		a.pos++
		return Token{Tipo: TOKEN_MULT, Valor: string(caractere)}
	case '/':
		a.pos++
		return Token{Tipo: TOKEN_DIV, Valor: string(caractere)}
	case '(':
		a.pos++
		return Token{Tipo: TOKEN_ABREPAR, Valor: string(caractere)}
	case ')':
		a.pos++
		return Token{Tipo: TOKEN_FECHAPAR, Valor: string(caractere)}
	default:
		a.pos++
		return Token{Tipo: TOKEN_INVALIDO, Valor: string(caractere)}
	}
}

// Interface para nós da árvore de sintaxe abstrata.
type Expressao interface{}

// Nó para constante inteira.
type Constante struct {
	Valor int
}

// Nó para operação binária.
type OperacaoBinaria struct {
	Operador string
	Esquerda Expressao
	Direita  Expressao
}

// Parser para a gramática EC2.
type AnalisadorSintatico struct {
	lex        *AnalisadorLexico
	tokenAtual Token
}

// mantem o token atual salvo e só chama o proximoToken quando quer realmente avançar
// resume o uso de uma função 'olhaProximoToken'
func NovoParser(lex *AnalisadorLexico) *AnalisadorSintatico {
	p := &AnalisadorSintatico{lex: lex}
	p.tokenAtual = p.lex.proximoToken()
	return p
}

func (p *AnalisadorSintatico) avancaToken(tipoEsperado TipoToken) error {
	if p.tokenAtual.Tipo == tipoEsperado {
		p.tokenAtual = p.lex.proximoToken()
		return nil
	}
	return fmt.Errorf("esperado token %s, mas encontrou %s", tipoEsperado, p.tokenAtual.Tipo)
}

// <prim> ::= <num> | '(' <exp_a> ')'
func (p *AnalisadorSintatico) analisarPrimaria() (Expressao, error) {
	token := p.tokenAtual
	if token.Tipo == TOKEN_INTEIRO {
		valor, err := strconv.Atoi(token.Valor)
		if err != nil {
			return nil, err
		}
		p.avancaToken(TOKEN_INTEIRO)
		return &Constante{Valor: valor}, nil
	} else if token.Tipo == TOKEN_ABREPAR {
		p.avancaToken(TOKEN_ABREPAR)
		expr, err := p.analisarExpA()
		if err != nil {
			return nil, err
		}
		if err := p.avancaToken(TOKEN_FECHAPAR); err != nil {
			return nil, err
		}
		return expr, nil
	}
	return nil, fmt.Errorf("token inesperado: %s", token.Valor)
}

// <exp_m> ::= <prim> (( '*' | '/' ) <prim>)*
func (p *AnalisadorSintatico) analisarExpM() (Expressao, error) {
	esquerda, err := p.analisarPrimaria()
	if err != nil {
		return nil, err
	}
	for p.tokenAtual.Tipo == TOKEN_MULT || p.tokenAtual.Tipo == TOKEN_DIV {
		operador := p.tokenAtual.Valor
		p.avancaToken(p.tokenAtual.Tipo)
		direita, err := p.analisarPrimaria()
		if err != nil {
			return nil, err
		}
		esquerda = &OperacaoBinaria{Operador: operador, Esquerda: esquerda, Direita: direita}
	}
	return esquerda, nil
}

// <exp_a> ::= <exp_m> (( '+' | '-' ) <exp_m>)*
func (p *AnalisadorSintatico) analisarExpA() (Expressao, error) {
	esquerda, err := p.analisarExpM()
	if err != nil {
		return nil, err
	}
	for p.tokenAtual.Tipo == TOKEN_SOMA || p.tokenAtual.Tipo == TOKEN_SUB {
		operador := p.tokenAtual.Valor
		p.avancaToken(p.tokenAtual.Tipo)
		direita, err := p.analisarExpM()
		if err != nil {
			return nil, err
		}
		esquerda = &OperacaoBinaria{Operador: operador, Esquerda: esquerda, Direita: direita}
	}
	return esquerda, nil
}

// Inicia a análise sintática a partir de exp_a.
func (p *AnalisadorSintatico) Analisar() (Expressao, error) {
	return p.analisarExpA()
}

// Função que avalia a árvore sintática.
func Avaliar(expr Expressao) (int, error) {
	switch e := expr.(type) {
	case *Constante:
		return e.Valor, nil
	case *OperacaoBinaria:
		esq, err := Avaliar(e.Esquerda)
		if err != nil {
			return 0, err
		}
		dir, err := Avaliar(e.Direita)
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
	}
	return 0, fmt.Errorf("expressão inválida")
}

// Função para imprimir a árvore sintática (em formato parenthesizado).
func ImprimirExpressao(expr Expressao) string {
	switch e := expr.(type) {
	case *Constante:
		return fmt.Sprintf("%d", e.Valor)
	case *OperacaoBinaria:
		return fmt.Sprintf("(%s %s %s)", ImprimirExpressao(e.Esquerda), e.Operador, ImprimirExpressao(e.Direita))
	}
	return ""
}

func main() {
	entrada := "25*(1)+65/(30+1)"
	//analisadorLexico := &AnalisadorLexico{texto: strings.ReplaceAll(entrada, " ", "")}
	analisadorLexico := &AnalisadorLexico{texto: entrada} //chamada sem limpeza de espaços para centralizar essa logica no lexer
	parser := NovoParser(analisadorLexico)
	expressao, err := parser.Analisar()
	if err != nil {
		fmt.Println("Erro no parser:", err)
		return
	}
	fmt.Println("Árvore sintática:", ImprimirExpressao(expressao))
	resultado, err := Avaliar(expressao)
	if err != nil {
		fmt.Println("Erro na avaliação:", err)
		return
	}
	fmt.Println("Resultado:", resultado)
}
