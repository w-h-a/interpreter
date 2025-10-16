package identifier

import (
	"github.com/w-h-a/interpreter/internal/ast"
)

type Identifier struct {
	Token ast.Token
	Value string
}

func (e *Identifier) TokenLiteral() string {
	return e.Token.Literal()
}

func (e *Identifier) String() string {
	return e.Value
}

func (e *Identifier) ExpressionNode() {}
