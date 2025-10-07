package expressionstatement

import (
	"github.com/w-h-a/interpreter/internal/parser/ast/expression"
	"github.com/w-h-a/interpreter/internal/token"
)

type Expression struct {
	Token      token.Token
	Expression expression.Expression
}

func (s *Expression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Expression) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}

func (s *Expression) StatementNode() {}
