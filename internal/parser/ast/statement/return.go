package statement

import (
	"github.com/w-h-a/interpreter/internal/parser/ast"
	"github.com/w-h-a/interpreter/internal/token"
)

type Return struct {
	Token token.Token
	Value ast.Expression
}

func (s *Return) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Return) StatementNode() {}
