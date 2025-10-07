package integer

import "github.com/w-h-a/interpreter/internal/token"

type Integer struct {
	Token token.Token
	Value int64
}

func (e *Integer) TokenLiteral() string {
	return e.Token.Literal
}

func (e *Integer) String() string {
	return e.Token.Literal
}

func (e *Integer) ExpressionNode() {}
