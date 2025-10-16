package let

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast"
	"github.com/w-h-a/interpreter/internal/ast/expression"
	"github.com/w-h-a/interpreter/internal/ast/expression/identifier"
)

type Let struct {
	Token ast.Token
	Name  *identifier.Identifier
	Value expression.Expression
}

func (s *Let) TokenLiteral() string {
	return s.Token.Literal()
}

func (s *Let) String() string {
	var out strings.Builder

	out.WriteString(s.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (s *Let) StatementNode() {}
