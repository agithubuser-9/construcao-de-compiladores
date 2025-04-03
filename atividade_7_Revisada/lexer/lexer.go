package lexer

import "unicode"

type AnalisadorLexico struct {
	Texto string
	Pos   int
}

func (a *AnalisadorLexico) ProximoToken() Token {
	for a.Pos < len(a.Texto) && unicode.IsSpace(rune(a.Texto[a.Pos])) {
		a.Pos++
	}
	if a.Pos >= len(a.Texto) {
		return Token{Tipo: TOKEN_FIM}
	}
	ch := a.Texto[a.Pos]
	switch {
	case unicode.IsDigit(rune(ch)):
		inicio := a.Pos
		for a.Pos < len(a.Texto) && unicode.IsDigit(rune(a.Texto[a.Pos])) {
			a.Pos++
		}
		return Token{Tipo: TOKEN_INTEIRO, Valor: a.Texto[inicio:a.Pos]}
	case ch == '+':
		a.Pos++
		return Token{Tipo: TOKEN_SOMA, Valor: string(ch)}
	case ch == '-':
		a.Pos++
		return Token{Tipo: TOKEN_SUB, Valor: string(ch)}
	case ch == '*':
		a.Pos++
		return Token{Tipo: TOKEN_MULT, Valor: string(ch)}
	case ch == '/':
		a.Pos++
		return Token{Tipo: TOKEN_DIV, Valor: string(ch)}
	case ch == '(':
		a.Pos++
		return Token{Tipo: TOKEN_ABREPAR, Valor: string(ch)}
	case ch == ')':
		a.Pos++
		return Token{Tipo: TOKEN_FECHAPAR, Valor: string(ch)}
	default:
		a.Pos++
		return Token{Tipo: TOKEN_INVALIDO, Valor: string(ch)}
	}
}
