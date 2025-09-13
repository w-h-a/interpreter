package token

type TokenType string

const (
	// Special tokens
	TokenIllegal TokenType = "ILLEGAL"
	TokenEOF     TokenType = "EOF"

	// Identifiers + literals
	TokenIdent TokenType = "IDENT"
	TokenInt   TokenType = "INT"

	// Operators
	TokenAssign TokenType = "="
	TokenPlus   TokenType = "+"

	// Delimiters
	TokenComma      TokenType = ","
	TokenSemicolon  TokenType = ";"
	TokenParenLeft  TokenType = "("
	TokenParenRight TokenType = ")"
	TokenBraceLeft  TokenType = "{"
	TokenBraceRight TokenType = "}"

	// Keywords
	TokenFunction TokenType = "FUNCTION"
	TokenLet      TokenType = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}
