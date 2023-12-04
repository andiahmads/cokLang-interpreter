package token

// definisikan token type
const (
	ILEGAL = "ILEGAL"
	EOF    = "EOF"

	// identifiers + literal
	IDENT = "IDENT" // add, foobar, x,y ......
	INT   = "INT"   // 1234567

	// operator
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// memungkinkan untuk menggunakan banyak nilai yang berbeda dan membedakan berbagai jenis token
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent memeriksa tabel kata kunci untuk melihat apakah pengenal yang diberikan adalah kata kunci.
// Jika ya, maka akan mengembalikan konstanta TokenType dari kata kunci tersebut.
// Jika tidak, kita hanya mendapatkan kembali token.IDENT yang merupakan TokenType untuk semua pengenal yang ditentukan pengguna.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
