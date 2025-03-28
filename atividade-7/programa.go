// programa.go
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Definição dos tipos de token.
type TokenType string

const (
	TOKEN_INT     TokenType = "INT"
	TOKEN_PLUS    TokenType = "+"
	TOKEN_MINUS   TokenType = "-"
	TOKEN_MUL     TokenType = "*"
	TOKEN_DIV     TokenType = "/"
	TOKEN_LPAREN  TokenType = "("
	TOKEN_RPAREN  TokenType = ")"
	TOKEN_EOF     TokenType = "EOF"
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
)

// Estrutura que representa um token.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer para análise léxica.
type Lexer struct {
	text string
	pos  int
}

func (l *Lexer) getNextToken() Token {
	// Ignora espaços.
	for l.pos < len(l.text) && unicode.IsSpace(rune(l.text[l.pos])) {
		l.pos++
	}
	if l.pos >= len(l.text) {
		return Token{Type: TOKEN_EOF, Value: ""}
	}
	ch := l.text[l.pos]
	// Reconhecimento de números.
	if unicode.IsDigit(rune(ch)) {
		start := l.pos
		for l.pos < len(l.text) && unicode.IsDigit(rune(l.text[l.pos])) {
			l.pos++
		}
		return Token{Type: TOKEN_INT, Value: l.text[start:l.pos]}
	}
	// Reconhecimento de símbolos.
	switch ch {
	case '+':
		l.pos++
		return Token{Type: TOKEN_PLUS, Value: string(ch)}
	case '-':
		l.pos++
		return Token{Type: TOKEN_MINUS, Value: string(ch)}
	case '*':
		l.pos++
		return Token{Type: TOKEN_MUL, Value: string(ch)}
	case '/':
		l.pos++
		return Token{Type: TOKEN_DIV, Value: string(ch)}
	case '(':
		l.pos++
		return Token{Type: TOKEN_LPAREN, Value: string(ch)}
	case ')':
		l.pos++
		return Token{Type: TOKEN_RPAREN, Value: string(ch)}
	default:
		l.pos++
		return Token{Type: TOKEN_ILLEGAL, Value: string(ch)}
	}
}

// Interface para nós da árvore de sintaxe abstrata.
type Expr interface{}

// Nó para constante inteira.
type Const struct {
	Value int
}

// Nó para operação binária.
type BinOp struct {
	Op    string
	Left  Expr
	Right Expr
}

// Parser para a gramática EC2.
type Parser struct {
	lexer        *Lexer
	currentToken Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.currentToken = p.lexer.getNextToken()
	return p
}

// Função auxiliar para consumir o token esperado.
func (p *Parser) eat(tokenType TokenType) error {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.getNextToken()
		return nil
	}
	return fmt.Errorf("esperado token %s, mas encontrou %s", tokenType, p.currentToken.Type)
}

// <prim> ::= <num> | '(' <exp_a> ')'
func (p *Parser) parsePrim() (Expr, error) {
	token := p.currentToken
	if token.Type == TOKEN_INT {
		value, err := strconv.Atoi(token.Value)
		if err != nil {
			return nil, err
		}
		p.eat(TOKEN_INT)
		return &Const{Value: value}, nil
	} else if token.Type == TOKEN_LPAREN {
		// Consome '(' e processa uma expressão aditiva.
		p.eat(TOKEN_LPAREN)
		expr, err := p.parseExpA()
		if err != nil {
			return nil, err
		}
		if err := p.eat(TOKEN_RPAREN); err != nil {
			return nil, err
		}
		return expr, nil
	}
	return nil, fmt.Errorf("token inesperado: %s", token.Value)
}

// <exp_m> ::= <prim> (( '*' | '/' ) <prim>)*
func (p *Parser) parseExpM() (Expr, error) {
	left, err := p.parsePrim()
	if err != nil {
		return nil, err
	}
	for p.currentToken.Type == TOKEN_MUL || p.currentToken.Type == TOKEN_DIV {
		op := p.currentToken.Value
		p.eat(p.currentToken.Type)
		right, err := p.parsePrim()
		if err != nil {
			return nil, err
		}
		left = &BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

// <exp_a> ::= <exp_m> (( '+' | '-' ) <exp_m>)*
func (p *Parser) parseExpA() (Expr, error) {
	left, err := p.parseExpM()
	if err != nil {
		return nil, err
	}
	for p.currentToken.Type == TOKEN_PLUS || p.currentToken.Type == TOKEN_MINUS {
		op := p.currentToken.Value
		p.eat(p.currentToken.Type)
		right, err := p.parseExpM()
		if err != nil {
			return nil, err
		}
		left = &BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

// Inicia a análise sintática a partir de exp_a.
func (p *Parser) Parse() (Expr, error) {
	return p.parseExpA()
}

// Função que avalia a árvore sintática.
func Evaluate(expr Expr) (int, error) {
	switch e := expr.(type) {
	case *Const:
		return e.Value, nil
	case *BinOp:
		left, err := Evaluate(e.Left)
		if err != nil {
			return 0, err
		}
		right, err := Evaluate(e.Right)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("divisão por zero")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", e.Op)
		}
	}
	return 0, fmt.Errorf("expressão inválida")
}

// Função para imprimir a árvore sintática (em formato parenthesizado).
func PrintExpr(expr Expr) string {
	switch e := expr.(type) {
	case *Const:
		return fmt.Sprintf("%d", e.Value)
	case *BinOp:
		return fmt.Sprintf("(%s %s %s)", PrintExpr(e.Left), e.Op, PrintExpr(e.Right))
	}
	return ""
}

func main() {
	// Exemplo de uso.
	input := "7+5*3-2"
	lexer := &Lexer{text: strings.ReplaceAll(input, " ", "")}
	parser := NewParser(lexer)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Println("Erro no parser:", err)
		return
	}
	fmt.Println("Árvore sintática:", PrintExpr(expr))
	result, err := Evaluate(expr)
	if err != nil {
		fmt.Println("Erro na avaliação:", err)
		return
	}
	fmt.Println("Resultado:", result)
}

