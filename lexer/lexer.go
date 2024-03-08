package lexer

import (
	"QuonkScript/token"
	"QuonkScript/utils"
)

const (
	eqSym = '='

	leftParen          = '('
	rightParen         = ')'
	leftSquareBracket  = '['
	rightSquareBracket = ']'
	leftCurlyBracket   = '{'
	rightCurlyBracket  = '}'

	semi  = ';'
	comma = ','
	colon = ':'
	dot   = '.'
	quote = '"'

	plus   = '+'
	star   = '*'
	slash  = '/'
	minus  = '-'
	modulo = '%'

	greaterThan = '>'
	lessThan    = '<'
	bang        = '!'
	ampersand   = '&'
	pipe        = '|'
)

type Lexer struct {
	source       string
	position     int
	readPosition int
	char         byte
	line         int
}

func New(source string) *Lexer {
	lexer := &Lexer{source: source, line: 1} // Start our lexer at line 1
	lexer.readChar()                         // set up lexer
	return lexer
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.source) {
		l.char = 0
	} else {
		l.char = l.source[l.readPosition]
	}

	// if we read in a newline or carriage return, we should reset the column counter
	if l.char == '\n' || l.char == '\r' {
		l.line += 1
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifer() string {
	position := l.position

	for utils.IsAlpha(string(l.char)) {
		l.readChar() // Advances the position pointer
	}
	return l.source[position:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	encounteredDecimal := false
	for utils.IsNumeric(string(l.char)) || (!encounteredDecimal && l.char == dot) {
		if l.char == dot {
			encounteredDecimal = true
		}
		l.readChar() // This just advances the position pointer
	}
	return l.source[position:l.position], encounteredDecimal
}

func (l *Lexer) readString() string {
	position := l.position + 1 // advance past ""

	for {
		l.readChar()
		if l.char == quote || l.char == 0 {
			break
		}
	}
	return l.source[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.source) {
		return 0
	} else {
		return l.source[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.char {
	// grouping
	case leftParen:
		tok = token.MakeToken(token.LeftParen, l.char, l.line)
	case rightParen:
		tok = token.MakeToken(token.RightParen, l.char, l.line)
	case leftCurlyBracket:
		tok = token.MakeToken(token.LeftCurlyBracket, l.char, l.line)
	case rightCurlyBracket:
		tok = token.MakeToken(token.RightCurlyBracket, l.char, l.line)
	case leftSquareBracket:
		tok = token.MakeToken(token.LeftSquareBracket, l.char, l.line)
	case rightSquareBracket:
		tok = token.MakeToken(token.RightSquareBracket, l.char, l.line)

	// Punctuation
	case semi:
		tok = token.MakeToken(token.Semicolon, l.char, l.line)
	case comma:
		tok = token.MakeToken(token.Comma, l.char, l.line)
	case colon:
		tok = token.MakeToken(token.Colon, l.char, l.line)
	case dot:
		tok = token.MakeToken(token.Dot, l.char, l.line)
	case quote:
		tok.Type = token.String
		tok.Line = l.line
		tok.Literal = l.readString()
	// Symbols
	case eqSym:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.EqualTo, Literal: literal, Line: l.line}
		} else {
			tok = token.MakeToken(token.Assign, l.char, l.line)
		}
	case plus:
		tok = token.MakeToken(token.Plus, l.char, l.line)
	case minus:
		tok = token.MakeToken(token.Minus, l.char, l.line)
	case star:
		tok = token.MakeToken(token.Star, l.char, l.line)
	case slash:
		tok = token.MakeToken(token.Slash, l.char, l.line)
	case modulo:
		tok = token.MakeToken(token.Modulo, l.char, l.line)
	case greaterThan:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.GreaterThanEqualTo, Literal: literal, Line: l.line}
		} else {
			tok = token.MakeToken(token.GreaterThan, l.char, l.line)
		}
	case lessThan:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.LessThanEqualTo, Literal: literal, Line: l.line}
		} else {
			tok = token.MakeToken(token.LessThan, l.char, l.line)
		}
	case bang:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.NotEqualTo, Literal: literal, Line: l.line}
		} else {
			tok = token.MakeToken(token.Bang, l.char, l.line)
		}
	case ampersand:
		if l.peekChar() == ampersand {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.And, Literal: literal, Line: l.line}
		} else {
			// Single & is an illegal char
			tok = token.MakeToken(token.Illegal, l.char, l.line)
		}
	case pipe:
		if l.peekChar() == pipe {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.Or, Literal: literal, Line: l.line}
		} else {
			// Single & is an illegal char
			tok = token.MakeToken(token.Illegal, l.char, l.line)
		}
	case 0:
		tok.Literal = ""
		tok.Type = "EOF"

	default:
		if utils.IsAlpha(string(l.char)) {
			tok.Literal = l.readIdentifer()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = l.line
			return tok // This is to avoid the l.readChar() call before this functions return
		} else if utils.IsNumeric(string(l.char)) {

			literal, decimal := l.readNumber()
			if decimal {
				tok.Type = token.Float
			} else {
				tok.Type = token.Integer
			}
			tok.Literal = literal
			tok.Line = l.line
			return tok // This is to avoid the l.readChar() call before this functions return
		} else {
			tok = token.MakeToken(token.Illegal, l.char, l.line)
		}
	}

	l.readChar()
	return tok
}

var keywords = map[string]token.TokenType{
	"mut":    token.Mut,
	"const":  token.Const,
	"null":   token.Null,
	"true":   token.True,
	"false":  token.False,
	"if":     token.If,
	"else":   token.Else,
	"elseif": token.Elseif,
	"func":   token.Func,
	"return": token.Return,
	"for":    token.For,
	"macro":  token.Macro,
}

func LookupIdent(ident string) token.TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return token.Identifier
}
