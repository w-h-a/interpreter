package block

import (
	"strings"

	"github.com/w-h-a/interpreter/internal/ast"
	"github.com/w-h-a/interpreter/internal/ast/statement"
)

type Block struct {
	Token      ast.Token
	Statements []statement.Statement
}

func (s *Block) TokenLiteral() string {
	return s.Token.Literal()
}

func (s *Block) String() string {
	var out strings.Builder

	for _, stmt := range s.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (s *Block) StatementNode() {}
