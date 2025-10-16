package boolean

import (
	"github.com/w-h-a/interpreter/internal/ast"
)

type Boolean struct {
	Token ast.Token
	Value bool
}

func (e *Boolean) TokenLiteral() string {
	return e.Token.Literal()
}

func (e *Boolean) String() string {
	return e.Token.Literal()
}

func (e *Boolean) ExpressionNode() {}
