package function

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast/expression/identifier"
	"github.com/w-h-a/interpreter/internal/ast/statement/block"
	"github.com/w-h-a/interpreter/internal/token"
)

type Function struct {
	Token      token.Token
	Parameters []*identifier.Identifier
	Body       *block.Block
}

func (e *Function) TokenLiteral() string {
	return e.Token.Literal
}

func (e *Function) String() string {
	var out strings.Builder

	params := []string{}

	for _, p := range e.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(e.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(e.Body.String())

	return out.String()
}

func (e *Function) ExpressionNode() {}
