package token

type TokenType string

const (
	// Special tokens
	Illegal TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	Ident TokenType = "IDENT"
	Int   TokenType = "INT"

	// Operators
	Assign TokenType = "="
	Plus   TokenType = "+"

	// Delimiters
	Comma      TokenType = ","
	Semicolon  TokenType = ";"
	ParenLeft  TokenType = "("
	ParenRight TokenType = ")"
	BraceLeft  TokenType = "{"
	BraceRight TokenType = "}"

	// Keywords
	Function TokenType = "FUNCTION"
	Let      TokenType = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}
