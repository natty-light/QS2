package lexer

import (
	"quonk/token"
	"quonk/utils"
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
}

func New(source string) *Lexer {
	lexer := &Lexer{source: source} // Start our lexer at line 1
	lexer.readChar()                // set up lexer
	return lexer
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.source) {
		l.char = 0
	} else {
		l.char = l.source[l.readPosition]
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
		tok = token.MakeToken(token.LeftParen, l.char, l.position, l.position)
	case rightParen:
		tok = token.MakeToken(token.RightParen, l.char, l.position, l.position)
	case leftCurlyBracket:
		tok = token.MakeToken(token.LeftCurlyBracket, l.char, l.position, l.position)
	case rightCurlyBracket:
		tok = token.MakeToken(token.RightCurlyBracket, l.char, l.position, l.position)
	case leftSquareBracket:
		tok = token.MakeToken(token.LeftSquareBracket, l.char, l.position, l.position)
	case rightSquareBracket:
		tok = token.MakeToken(token.RightSquareBracket, l.char, l.position, l.position)

	// Punctuation
	case semi:
		tok = token.MakeToken(token.Semicolon, l.char, l.position, l.position)
	case comma:
		tok = token.MakeToken(token.Comma, l.char, l.position, l.position)
	case colon:
		tok = token.MakeToken(token.Colon, l.char, l.position, l.position)
	case dot:
		tok = token.MakeToken(token.Dot, l.char, l.position, l.position)
	case quote:
		tok.Type = token.String
		tok.StartPos = l.position
		tok.Literal = l.readString()
		tok.EndPos = l.position
	// Symbols
	case eqSym:
		if l.peekChar() == eqSym {
			char := l.char
			start := l.position
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.EqualTo, Literal: literal, StartPos: start, EndPos: l.position}
		} else {
			tok = token.MakeToken(token.Assign, l.char, l.position, l.position)
		}
	case plus:
		tok = token.MakeToken(token.Plus, l.char, l.position, l.position)
	case minus:
		tok = token.MakeToken(token.Minus, l.char, l.position, l.position)
	case star:
		tok = token.MakeToken(token.Star, l.char, l.position, l.position)
	case slash:
		tok = token.MakeToken(token.Slash, l.char, l.position, l.position)
	case modulo:
		tok = token.MakeToken(token.Modulo, l.char, l.position, l.position)
	case greaterThan:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.GreaterThanEqualTo, Literal: literal, StartPos: l.position, EndPos: l.position + 1}
		} else {
			tok = token.MakeToken(token.GreaterThan, l.char, l.position, l.position)
		}
	case lessThan:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.LessThanEqualTo, Literal: literal, StartPos: l.position, EndPos: l.position + 1}
		} else {
			tok = token.MakeToken(token.LessThan, l.char, l.position, l.position)
		}
	case bang:
		if l.peekChar() == eqSym {
			char := l.char
			l.readChar() // advance past first equals
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.NotEqualTo, Literal: literal, StartPos: l.position, EndPos: l.position + 1}
		} else {
			tok = token.MakeToken(token.Bang, l.char, l.position, l.position)
		}
	case ampersand:
		if l.peekChar() == ampersand {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.And, Literal: literal, StartPos: l.position, EndPos: l.position + 1}
		} else {
			// Single & is an illegal char
			tok = token.MakeToken(token.Illegal, l.char, l.position, l.position)
		}
	case pipe:
		if l.peekChar() == pipe {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)
			tok = token.Token{Type: token.Or, Literal: literal, StartPos: l.position, EndPos: l.position + 1}
		} else {
			// Single & is an illegal char
			tok = token.MakeToken(token.Illegal, l.char, l.position, l.position)
		}
	case 0:
		tok.Literal = ""
		tok.Type = "EOF"

	default:
		if utils.IsAlpha(string(l.char)) {
			tok.StartPos = l.position
			tok.Literal = l.readIdentifer()
			tok.EndPos = l.position
			tok.Type = LookupIdent(tok.Literal)
			return tok // This is to avoid the l.readChar() call before this functions return
		} else if utils.IsNumeric(string(l.char)) {
			tok.StartPos = l.position
			literal, decimal := l.readNumber()
			tok.EndPos = l.position
			if decimal {
				tok.Type = token.Float
			} else {
				tok.Type = token.Integer
			}
			tok.Literal = literal

			return tok // This is to avoid the l.readChar() call before this functions return
		} else {
			tok = token.MakeToken(token.Illegal, l.char, l.position, l.position)
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
