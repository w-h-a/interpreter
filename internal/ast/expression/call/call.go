package call

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast"
	"github.com/w-h-a/interpreter/internal/ast/expression"
)

type Call struct {
	Token     ast.Token
	Function  expression.Expression
	Arguments []expression.Expression
}

func (e *Call) TokenLiteral() string {
	return e.Token.Literal()
}

func (e *Call) String() string {
	var out strings.Builder

	args := []string{}

	for _, a := range e.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(e.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (e *Call) ExpressionNode() {

}
