package expressionstatement

import (
	"github.com/w-h-a/interpreter/internal/ast"
	"github.com/w-h-a/interpreter/internal/ast/expression"
)

type Expression struct {
	Token      ast.Token
	Expression expression.Expression
}

func (s *Expression) TokenLiteral() string {
	return s.Token.Literal()
}

func (s *Expression) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}

func (s *Expression) StatementNode() {}
