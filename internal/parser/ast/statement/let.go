package statement

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/token"
)

type Let struct {
	Token token.Token
	Name  *expression.Identifier
	Value ast.Expression
}

func (s *Let) TokenLiteral() string {
	return s.Token.Literal
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
