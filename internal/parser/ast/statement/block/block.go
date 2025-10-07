package block

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/parser/ast/statement"
	"github.com/w-h-a/interpreter/internal/token"
)

type Block struct {
	Token      token.Token
	Statements []statement.Statement
}

func (s *Block) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Block) String() string {
	var out strings.Builder

	for _, stmt := range s.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (s *Block) StatementNode() {}
