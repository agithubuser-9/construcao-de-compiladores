package lexer

type TipoToken string

const (
	TOKEN_INTEIRO    TipoToken = "INT"
	TOKEN_SOMA       TipoToken = "+"
	TOKEN_SUB        TipoToken = "-"
	TOKEN_MULT       TipoToken = "*"
	TOKEN_DIV        TipoToken = "/"
	TOKEN_ABREPAR    TipoToken = "("
	TOKEN_FECHAPAR   TipoToken = ")"
	TOKEN_ABRECHAVE  TipoToken = "{"
	TOKEN_FECHCHAVE  TipoToken = "}"
	TOKEN_FIM        TipoToken = "EOF"
	TOKEN_INVALIDO   TipoToken = "ILLEGAL"
	TOKEN_ATRIBUICAO TipoToken = "="
	TOKEN_IGUALDADE  TipoToken = "=="
	TOKEN_MENORQ     TipoToken = "<"
	TOKEN_MAIORQ     TipoToken = ">"
	TOKEN_MENORIGUAL TipoToken = "<="
	TOKEN_MAIORIGUAL TipoToken = ">="
	TOKEN_PONTOVIRG  TipoToken = ";"
	TOKEN_IDENT      TipoToken = "IDENT"
	TOKEN_RETURN     TipoToken = "RETURN"
	TOKEN_IF         TipoToken = "IF"
	TOKEN_ELSE       TipoToken = "ELSE"
	TOKEN_WHILE      TipoToken = "WHILE"
)

type Token struct {
	Tipo  TipoToken
	Valor string
}
