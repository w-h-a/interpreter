package token

var keywords = map[string]TokenType{
	"fn":  Function,
	"let": Let,
}

func LookupIdent(ident string) TokenType {
	if tk, ok := keywords[ident]; ok {
		return tk
	}
	return Ident
}
