package token

func Factory(t TokenType, char string) (Token, error) {
	return Token{
		Type:    t,
		Literal: char,
	}, nil
}
