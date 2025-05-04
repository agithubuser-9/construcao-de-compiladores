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

	ch := rune(a.Texto[a.Pos])
	// Identificadores (começam com letra, podem ter letras ou dígitos)
	if unicode.IsLetter(ch) {
		inicio := a.Pos
		a.Pos++
		for a.Pos < len(a.Texto) && (unicode.IsLetter(rune(a.Texto[a.Pos])) || unicode.IsDigit(rune(a.Texto[a.Pos]))) {
			a.Pos++
		}
		ident := a.Texto[inicio:a.Pos]
		if ident == "return" {
			return Token{Tipo: TOKEN_RETURN, Valor: ident}
		}
		return Token{Tipo: TOKEN_IDENT, Valor: ident}
	}

	// Números
	if unicode.IsDigit(ch) {
		inicio := a.Pos
		a.Pos++
		for a.Pos < len(a.Texto) && unicode.IsDigit(rune(a.Texto[a.Pos])) {
			a.Pos++
		}
		// Verificação extra: se depois do número tiver uma letra, erro
		if a.Pos < len(a.Texto) && unicode.IsLetter(rune(a.Texto[a.Pos])) {
			return Token{Tipo: TOKEN_INVALIDO, Valor: a.Texto[inicio : a.Pos+1]}
		}
		return Token{Tipo: TOKEN_INTEIRO, Valor: a.Texto[inicio:a.Pos]}
	}

	// Operadores e símbolos
	a.Pos++
	switch ch {
	case '+':
		return Token{Tipo: TOKEN_SOMA, Valor: string(ch)}
	case '-':
		return Token{Tipo: TOKEN_SUB, Valor: string(ch)}
	case '*':
		return Token{Tipo: TOKEN_MULT, Valor: string(ch)}
	case '/':
		return Token{Tipo: TOKEN_DIV, Valor: string(ch)}
	case '(':
		return Token{Tipo: TOKEN_ABREPAR, Valor: string(ch)}
	case ')':
		return Token{Tipo: TOKEN_FECHAPAR, Valor: string(ch)}
	case '{':
		return Token{Tipo: TOKEN_ABRECHAVE, Valor: string(ch)}
	case '}':
		return Token{Tipo: TOKEN_FECHCHAVE, Valor: string(ch)}
	case '=':
		return Token{Tipo: TOKEN_IGUAL, Valor: string(ch)}
	case ';':
		return Token{Tipo: TOKEN_PONTOVIRG, Valor: string(ch)}
	default:
		return Token{Tipo: TOKEN_INVALIDO, Valor: string(ch)}
	}
}
