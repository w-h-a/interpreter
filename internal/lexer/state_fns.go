package lexer

import (
	"github.com/w-h-a/interpreter/internal/token"
)

type stateFn func(*Lexer) stateFn

func lex(l *Lexer) stateFn {
	l.skip()

	switch char := l.peek(); {
	case char == 0:
		return lexStop
	case IsLetter(char):
		return lexIdentifier
	case IsDigit(char):
		return lexNumber
	default:
		return lexSymbol
	}
}

func lexIdentifier(l *Lexer) stateFn {
	l.next()

	for l.pos < len(l.input) && IsLetter(l.input[l.pos]) {
		l.pos += 1
	}

	literal := l.input[l.start:l.pos]

	l.emit(token.LookupIdent(literal))

	return lex
}

func lexNumber(l *Lexer) stateFn {
	l.next()

	for l.pos < len(l.input) && IsDigit(l.input[l.pos]) {
		l.pos += 1
	}

	l.emit(token.Int)

	return lex
}

func lexSymbol(l *Lexer) stateFn {
	switch char := l.next(); char {
	case '=':
		return lexEqual
	case '+':
		l.emit(token.Plus)
	case '-':
		l.emit(token.Minus)
	case '!':
		return lexBang
	case '*':
		l.emit(token.Asterisk)
	case '/':
		l.emit(token.Slash)
	case '<':
		l.emit(token.LessThan)
	case '>':
		l.emit(token.GreaterThan)
	case '(':
		l.emit(token.ParenLeft)
	case ')':
		l.emit(token.ParenRight)
	case '{':
		l.emit(token.BraceLeft)
	case '}':
		l.emit(token.BraceRight)
	case ',':
		l.emit(token.Comma)
	case ';':
		l.emit(token.Semicolon)
	default:
		l.emit(token.Illegal)
	}

	return lex
}

func lexEqual(l *Lexer) stateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(token.Identical)
	} else {
		l.emit(token.Assign)
	}

	return lex
}

func lexBang(l *Lexer) stateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(token.NotIdentical)
	} else {
		l.emit(token.Bang)
	}

	return lex
}

func lexStop(l *Lexer) stateFn {
	l.emit(token.EOF)
	return nil
}
