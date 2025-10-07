package identifier

import (
	"github.com/w-h-a/interpreter/internal/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (e *Identifier) TokenLiteral() string {
	return e.Token.Literal
}

func (e *Identifier) String() string {
	return e.Value
}

func (e *Identifier) ExpressionNode() {}
