package lexer

import "QuonkScript/utils"

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

func CreateLexer(source string) *Lexer {
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

func (l *Lexer) readNumber() string {
	position := l.position
	for utils.IsNumeric(string(l.char)) {
		l.readChar() // This just advances the position pointer
	}
	return l.source[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.char {
	// Assignment
	case eqSym:
		tok = token(Assign, l.char, l.line)

	// grouping
	case leftParen:
		tok = token(LeftParen, l.char, l.line)
	case rightParen:
		tok = token(RightParen, l.char, l.line)
	case leftCurlyBracket:
		tok = token(LeftCurlyBracket, l.char, l.line)
	case rightCurlyBracket:
		tok = token(RightCurlyBracket, l.char, l.line)
	case leftSquareBracket:
		tok = token(RightSquareBracket, l.char, l.line)
	case rightSquareBracket:
		tok = token(RightSquareBracket, l.char, l.line)

	// Punctuation
	case semi:
		tok = token(Semicolon, l.char, l.line)
	case comma:
		tok = token(Comma, l.char, l.line)
	case colon:
		tok = token(Colon, l.char, l.line)
	case dot:
		tok = token(Dot, l.char, l.line)

	// Symbols
	case plus:
		tok = token(Plus, l.char, l.line)
	case minus:
		tok = token(Minus, l.char, l.line)
	case star:
		tok = token(Star, l.char, l.line)
	case slash:
		tok = token(Slash, l.char, l.line)
	case modulo:
		tok = token(Modulo, l.char, l.line)
	case greaterThan:
		tok = token(GreaterThan, l.char, l.line)
	case lessThan:
		tok = token(LessThan, l.char, l.line)
	case bang:
		tok = token(Bang, l.char, l.line)
		// TODO: Logical operators
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
			tok.Type = Int
			tok.Literal = l.readNumber()
			tok.Line = l.line
			return tok // This is to avoid the l.readChar() call before this functions return
		} else {
			tok = token(Illegal, l.char, l.line)
		}
	}

	l.readChar()
	return tok
}

var keywords = map[string]TokenType{
	"mut":    Mut,
	"const":  Const,
	"null":   Null,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"elseif": Elseif,
	"func":   Func,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Identifier
}
