package token

func Factory(t TokenType, char string) Token {
	return Token{
		Type:    t,
		literal: char,
	}
}
