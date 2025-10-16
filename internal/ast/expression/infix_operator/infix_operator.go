package infixoperator

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast"
	"github.com/w-h-a/interpreter/internal/ast/expression"
)

type InfixOperator struct {
	Token    ast.Token
	Operator string
	Left     expression.Expression
	Right    expression.Expression
}

func (e *InfixOperator) TokenLiteral() string {
	return e.Token.Literal()
}

func (e *InfixOperator) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(e.Left.String())
	out.WriteString(" ")
	out.WriteString(e.Operator)
	out.WriteString(" ")
	out.WriteString(e.Right.String())
	out.WriteString(")")

	return out.String()
}

func (e *InfixOperator) ExpressionNode() {}
