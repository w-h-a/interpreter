package parser

import "github.com/w-h-a/interpreter/internal/token"

const (
	LOWEST int = iota
	EQUALITY
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var (
	precedences = map[token.TokenType]int{
		token.Identical:    EQUALITY,
		token.NotIdentical: EQUALITY,
		token.LessThan:     LESSGREATER,
		token.GreaterThan:  LESSGREATER,
		token.Plus:         SUM,
		token.Minus:        SUM,
		token.Asterisk:     PRODUCT,
		token.Slash:        PRODUCT,
		token.ParenLeft:    CALL,
	}
)
