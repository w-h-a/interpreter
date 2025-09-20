package statement

import (
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

func (s *Let) StatementNode() {}
