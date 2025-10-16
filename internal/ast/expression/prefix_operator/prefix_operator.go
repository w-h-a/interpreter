package prefixoperator

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast/expression"
	"github.com/w-h-a/interpreter/internal/token"
)

type PrefixOperator struct {
	Token    token.Token
	Operator string
	Right    expression.Expression
}

func (e *PrefixOperator) TokenLiteral() string {
	return e.Token.Literal
}

func (e *PrefixOperator) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(e.Operator)
	out.WriteString(e.Right.String())
	out.WriteString(")")

	return out.String()
}

func (e *PrefixOperator) ExpressionNode() {}
