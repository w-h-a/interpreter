package token

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if tk, ok := keywords[ident]; ok {
		return tk
	}
	return Ident
}
