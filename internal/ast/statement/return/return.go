package returnstatement

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast"
	"github.com/w-h-a/interpreter/internal/ast/expression"
)

type Return struct {
	Token ast.Token
	Value expression.Expression
}

func (s *Return) TokenLiteral() string {
	return s.Token.Literal()
}

func (s *Return) String() string {
	var out strings.Builder

	out.WriteString(s.TokenLiteral())
	out.WriteString(" ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (s *Return) StatementNode() {}
