package lexer

import "github.com/w-h-a/interpreter/internal/token"

type Lexer interface {
	NextToken() token.Token
}
