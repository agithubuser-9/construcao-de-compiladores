package lexer

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

type Token struct {
	Tipo  TipoToken
	Valor string
}
