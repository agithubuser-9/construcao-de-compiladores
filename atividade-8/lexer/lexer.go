// lexer.go
package lexer

import "unicode"

type AnalisadorLexico struct {
	Texto string
	Pos   int
}

func (a *AnalisadorLexico) ProximoToken() Token {
	// Ignora espaços em branco.
	for a.Pos < len(a.Texto) && unicode.IsSpace(rune(a.Texto[a.Pos])) {
		a.Pos++
	}
	if a.Pos >= len(a.Texto) {
		return Token{Tipo: TOKEN_FIM}
	}
	ch := a.Texto[a.Pos]
	// Se for letra, consome o identificador.
	if unicode.IsLetter(rune(ch)) {
		inicio := a.Pos
		for a.Pos < len(a.Texto) && (unicode.IsLetter(rune(a.Texto[a.Pos])) || unicode.IsDigit(rune(a.Texto[a.Pos]))) {
			a.Pos++
		}
		return Token{Tipo: TOKEN_IDENT, Valor: a.Texto[inicio:a.Pos]}
	}
	// Se for dígito, consome o número inteiro.
	if unicode.IsDigit(rune(ch)) {
		inicio := a.Pos
		for a.Pos < len(a.Texto) && unicode.IsDigit(rune(a.Texto[a.Pos])) {
			a.Pos++
		}
		return Token{Tipo: TOKEN_INTEIRO, Valor: a.Texto[inicio:a.Pos]}
	}
	// Reconhecimento dos símbolos.
	switch ch {
	case '+':
		a.Pos++
		return Token{Tipo: TOKEN_SOMA, Valor: string(ch)}
	case '-':
		a.Pos++
		return Token{Tipo: TOKEN_SUB, Valor: string(ch)}
	case '*':
		a.Pos++
		return Token{Tipo: TOKEN_MULT, Valor: string(ch)}
	case '/':
		a.Pos++
		return Token{Tipo: TOKEN_DIV, Valor: string(ch)}
	case '(':
		a.Pos++
		return Token{Tipo: TOKEN_ABREPAR, Valor: string(ch)}
	case ')':
		a.Pos++
		return Token{Tipo: TOKEN_FECHAPAR, Valor: string(ch)}
	case '=':
		a.Pos++
		return Token{Tipo: TOKEN_IGUAL, Valor: string(ch)}
	case ';':
		a.Pos++
		return Token{Tipo: TOKEN_PONTOEVIRGULA, Valor: string(ch)}
	default:
		a.Pos++
		return Token{Tipo: TOKEN_INVALIDO, Valor: string(ch)}
	}
}

