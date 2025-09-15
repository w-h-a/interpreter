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
	Assign       TokenType = "="
	Plus         TokenType = "+"
	Minus        TokenType = "-"
	Bang         TokenType = "!"
	Asterisk     TokenType = "*"
	Slash        TokenType = "/"
	LessThan     TokenType = "<"
	GreaterThan  TokenType = ">"
	Identical    TokenType = "=="
	NotIdentical TokenType = "!="

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
	True     TokenType = "TRUE"
	False    TokenType = "FALSE"
	If       TokenType = "IF"
	Else     TokenType = "ELSE"
	Return   TokenType = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}
