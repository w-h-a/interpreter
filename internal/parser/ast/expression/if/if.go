package ifexpression

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/parser/ast/statement/block"
	"github.com/w-h-a/interpreter/internal/token"
)

type If struct {
	Token       token.Token
	Condition   expression.Expression
	Consequence *block.Block
	Alternative *block.Block
}

func (e *If) TokenLiteral() string {
	return e.Token.Literal
}

func (e *If) String() string {
	var out strings.Builder

	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString(e.Condition.String())
	out.WriteString(" ")
	out.WriteString(e.Consequence.String())

	if e.Alternative != nil {
		out.WriteString(" ")
		out.WriteString("else")
		out.WriteString(" ")
		out.WriteString(e.Alternative.String())
	}

	return out.String()
}

func (e *If) ExpressionNode() {}
