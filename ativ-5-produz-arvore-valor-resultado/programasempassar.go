package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Tipos de token usados na linguagem EC1.
type TokenType string

const (
	TOKEN_INT    TokenType = "INT"
	TOKEN_PLUS   TokenType = "+"
	TOKEN_MINUS  TokenType = "-"
	TOKEN_MUL    TokenType = "*"
	TOKEN_DIV    TokenType = "/"
	TOKEN_LPAREN TokenType = "("
	TOKEN_RPAREN TokenType = ")"
	TOKEN_EOF    TokenType = "EOF"
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
)

// Estrutura que representa um token (unidade léxica).
type Token struct {
	Type  TokenType
	Value string
}

// Lexer responsável pela análise léxica da entrada.
type Lexer struct {
	text string
	pos  int
}

// Retorna o próximo token da entrada.
func (l *Lexer) getNextToken() Token {
	// Ignora espaços em branco.
	for l.pos < len(l.text) && unicode.IsSpace(rune(l.text[l.pos])) {
		l.pos++
	}
	if l.pos >= len(l.text) {
		return Token{Type: TOKEN_EOF, Value: ""}
	}

	ch := l.text[l.pos]

	// Se for dígito, consome todos os dígitos formando um literal inteiro.
	if unicode.IsDigit(rune(ch)) {
		start := l.pos
		for l.pos < len(l.text) && unicode.IsDigit(rune(l.text[l.pos])) {
			l.pos++
		}
		numStr := l.text[start:l.pos]
		return Token{Type: TOKEN_INT, Value: numStr}
	}

	// Reconhecimento dos símbolos.
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

// Expressão (nó da AST) é definida por uma interface.
type Expr interface{}

// Nó que representa uma constante inteira.
type Const struct {
	Value int
}

// Nó que representa uma operação binária.
type BinOp struct {
	Op    string
	Left  Expr
	Right Expr
}

// Parser implementa a análise sintática descendente recursiva.
type Parser struct {
	lexer        *Lexer
	currentToken Token
}

// Cria um novo parser a partir de um lexer.
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.currentToken = p.lexer.getNextToken()
	return p
}

// Função que analisa uma expressão conforme a gramática EC1.
// Se o token atual é '(', processa uma operação binária; 
// se é um literal inteiro, cria um nó Const.
func (p *Parser) ParseExpression() (Expr, error) {
	token := p.currentToken

	if token.Type == TOKEN_LPAREN {
		// Consome '('
		p.currentToken = p.lexer.getNextToken()
		// Analisa a expressão à esquerda.
		left, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		// O próximo token deve ser um operador.
		opToken := p.currentToken
		if opToken.Type != TOKEN_PLUS && opToken.Type != TOKEN_MINUS &&
			opToken.Type != TOKEN_MUL && opToken.Type != TOKEN_DIV {
			return nil, fmt.Errorf("esperado operador, encontrado '%s'", opToken.Value)
		}
		op := opToken.Value
		p.currentToken = p.lexer.getNextToken()
		// Analisa a expressão à direita.
		right, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		// Verifica o fechamento com ')'.
		if p.currentToken.Type != TOKEN_RPAREN {
			return nil, fmt.Errorf("esperado ')', encontrado '%s'", p.currentToken.Value)
		}
		// Consome ')'
		p.currentToken = p.lexer.getNextToken()
		return &BinOp{Op: op, Left: left, Right: right}, nil

	} else if token.Type == TOKEN_INT {
		// Cria um nó de constante.
		value, err := strconv.Atoi(token.Value)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter número: %s", token.Value)
		}
		p.currentToken = p.lexer.getNextToken()
		return &Const{Value: value}, nil
	}

	return nil, fmt.Errorf("token inesperado: '%s'", token.Value)
}

// Função que avalia a AST e retorna o valor da expressão.
func Evaluate(expr Expr) (int, error) {
	switch e := expr.(type) {
	case *Const:
		return e.Value, nil
	case *BinOp:
		leftVal, err := Evaluate(e.Left)
		if err != nil {
			return 0, err
		}
		rightVal, err := Evaluate(e.Right)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case "+":
			return leftVal + rightVal, nil
		case "-":
			return leftVal - rightVal, nil
		case "*":
			return leftVal * rightVal, nil
		case "/":
			if rightVal == 0 {
				return 0, fmt.Errorf("divisão por zero")
			}
			return leftVal / rightVal, nil
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", e.Op)
		}
	default:
		return 0, fmt.Errorf("expressão inválida")
	}
}

// Função para imprimir a AST de forma que se assemelhe à expressão original.
func PrintExpr(expr Expr) string {
	switch e := expr.(type) {
	case *Const:
		return fmt.Sprintf("%d", e.Value)
	case *BinOp:
		// Impressão entre parênteses com o operador no meio.
		return fmt.Sprintf("(%s %s %s)", PrintExpr(e.Left), e.Op, PrintExpr(e.Right))
	default:
		return ""
	}
}

// Função auxiliar para limpar espaços em excesso na entrada.
func cleanInput(input string) string {
	return strings.TrimSpace(input)
}

// Função principal com exemplos de teste.
// São testadas expressões válidas e um exemplo de erro sintático.
func main() {
	tests := []string{
		"333",
		"(6 * 7)",
		"(3 + (4 + (11 + 7)))",
		"(33 + (912 * 11))",
		"((427 / 7) + (11 * (231 + 5)))",
		"(3 + )", // Exemplo de erro sintático.
	}

	for _, exprStr := range tests {
		fmt.Println("Expressão:", exprStr)
		lexer := &Lexer{text: cleanInput(exprStr)}
		parser := NewParser(lexer)
		expr, err := parser.ParseExpression()
		if err != nil {
			fmt.Println("Erro na análise sintática:", err)
		} else {
			fmt.Println("Árvore:", PrintExpr(expr))
			result, err := Evaluate(expr)
			if err != nil {
				fmt.Println("Erro na interpretação:", err)
			} else {
				fmt.Println("Resultado:", result)
			}
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}

